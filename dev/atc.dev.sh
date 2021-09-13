# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

alias atc-start="docker-compose up -d --build";
alias atc-stop="docker-compose down -v --remove-orphans";
alias atc-restart="atc-stop && atc-start";

function atc-ready {
	curl -skL https://localhost:6443/api/4.0/ping >/dev/null 2>&1;
	return $?;
}

function atc-exec {
	if [[ $# -lt 2 ]]; then
		echo "Usage: atc-exec SERVICE CMD" >&2;
		return 1;
	fi;
	local service="trafficcontrol_$1_1";
	shift;
	docker exec "$service" $@;
	return $?;
}

function atc-connect {
	if [[ $# -ne 1 ]]; then
		echo "Usage: atc-connect SERVICE" >&2;
		return 1;
	fi;
	docker exec -it "trafficcontrol_$1_1" /bin/sh;
	return $?;
}

export t3cDir="/root/go/src/github.com/apache/trafficcontrol/cache-config";

function t3c {
	trap 'atc-exec t3c ps | grep dlv | tr -s " " | cut -d " " -f1 | xargs docker exec trafficcontrol_t3c_1 kill' INT;
	local dExec=(docker exec);
	local dlv=();
	if [[ ! -z "$TC_WAIT" ]]; then
		dlv=(dlv --accept-multiclient --listen=:8081 --headless --api-version=2 debug --);
	else;
		dlv=(dlv --accept-multiclient --continue --listen=:8081 --headless --api-version=2 debug --);
	fi
	if [[ $# -lt 2 ]]; then
		$dExec -w "$t3cDir/t3c" trafficcontrol_t3c_1 $dlv;
		return $?;
	fi;
	local pkg="t3c"
	case "$1" in
		apply|check|check-refs|check-reload|diff|generate|preprocess|request|update)
			pkg="t3c-$1";;
	esac;
	shift;
	$dExec -w "$t3cDir/$pkg" trafficcontrol_t3c_1 $dlv $@;
	return "$?";
}

function tm-health-client {
	trap 'atc-exec t3c ps | grep dlv | tr -s " " | cut -d " " -f1 | xargs docker exec trafficcontrol_t3c_1 kill' INT;
	local dExec=(docker exec -w "$t3cDir/tm-health-client" trafficcontrol_t3c_1);
	local dlv=();
	if [[ ! -z "$TC_WAIT" ]]; then
		dlv=(dlv --accept-multiclient --listen=:8081 --headless --api-version=2 debug --);
	else
		dlv=(dlv --accept-multiclient --continue --listen=:8081 --headless --api-version=2 debug --);
	fi
	$dExec $dlv $@ || atc-exec t3c 'ps | grep dlv | tr -s ' ' | cut -d -f1 | xargs kill';
	return $?;
}

export TO_URL="https://localhost:6443"
export TO_USER="admin"
export TO_PASSWORD="twelve12"
