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

alias atc-start="docker compose up -d --build";
alias atc-build="docker compose build";
alias atc-stop="docker compose kill && docker compose down -v --remove-orphans";

function atc-restart {
	if ! atc-stop $@; then
		return 1;
	fi
	atc-start $@;
	return $?;
}

function atc-ready {
	local usage url="https://localhost:6443/api/4.0/ping";
	usage="$(<<-USAGE cat
		Usage: ${0} [-h] [-w]

		-h, --help
		      print usage information and exit
		-w, --wait
		      wait for ATC to be ready, instead of just checking if it is ready
		-d, --delivery-service
		      wait for the ATC delivery service to be ready
		USAGE
	)"
	if [[ $# -gt 0 ]]; then
		case "$1" in
			-w|--wait)
				until curl -skL "$url" >/dev/null 2>&1; do
					sleep 1;
				done
				return 0;;
			-d|--delivery-service)
				local deliveryservice=cdn.dev-ds.ciab.test
				until curl -4sfH "Host: ${deliveryservice}" localhost:3080 &&
								<<<"$(dig +short -4 @localhost -p 3053 "$deliveryservice")" grep -q '^[0-9.]\+$';
				do
					sleep 1;
				done
				return 0;;
			-h|--help)
				echo "$usage";
				return 0;;
			*)
				echo "$usage" >&2;
				return 1;;
		esac
	fi
	timeout 1s curl -skL "$url" >/dev/null 2>&1;
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

function atc {
	if [[ $# -lt 1 ]]; then
		echo "Usage: atc OPERATION" >&2;
		return 1;
	fi
	local arg="$1";
	shift;
	case "$arg" in
		build)
			atc-build $@;;
		connect)
			atc-connect $@;;
		exec)
			atc-exec $@;;
		ready)
			atc-ready $@;;
		restart)
			atc-restart $@;;
		start)
			atc-start $@;;
		stop)
			atc-stop $@;;
		-h|--help|/\?)
			echo "Usage: atc OPERATION";
			echo "";
			echo "Valid OPERATIONs:";
			echo "  build   Build the images for the environment, but do not start it";
			echo "  connect Connect to a shell session inside a dev container";
			echo "  exec    Run a command in a dev container";
			echo "  ready   Check if the development environment is ready";
			echo "  restart Restart the development environment";
			echo "  start   Start up the development environment";
			echo "  stop    Stop the development environment";
			;;
		*)
			echo "Usage: atc OPERATION" >&2;
			echo "" >&2;
			echo "Valid OPERATIONs:" >&2;
			echo "  build   Build the images for the environment, but do not start it" >&2;
			echo "  connect Connect to a shell session inside a dev container" >&2;
			echo "  exec    Run a command in a dev container" >&2;
			echo "  ready   Check if the development environment is ready" >&2;
			echo "  restart Restart the development environment" >&2;
			echo "  start   Start up the development environment" >&2;
			echo "  stop    Stop the development environment" >&2;
			return 2;;
	esac
	return "$?";
}

export t3cDir="/go/src/github.com/apache/trafficcontrol/cache-config";

function t3c {
	trap 'atc-exec t3c ps | grep dlv | tr -s " " | cut -d " " -f1 | xargs docker exec trafficcontrol_t3c_1 kill' INT;
	local dExec=(docker exec);
	local dlv=();
	if [[ ! -z "$TC_WAIT" ]]; then
		dlv=(dlv --accept-multiclient --listen=:8081 --headless --api-version=2 debug --);
	else
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
