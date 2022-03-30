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
alias atc-build="docker-compose build";
alias atc-stop="docker-compose kill && docker-compose down -v --remove-orphans";

function atc-restart {
	if ! atc-stop $@; then
		return 1;
	fi
	atc-start $@;
	return $?;
}

function atc-ready {
	local url="https://localhost:6443/api/${API_VERSION}/ping";
	if [[ $# -gt 0 ]]; then
		case "$1" in
			-w|--wait)
				until curl -skL "$url" >/dev/null 2>&1; do
					sleep 1;
				done
				return 0;;
			-d|--deliveryservice)
				local deliveryservice=cdn.dev-ds.ciab.test
				until curl -4sfH "Host: ${deliveryservice}" localhost:3080 &&
								<<<"$(dig +short -4 @localhost -p 3053 "$deliveryservice")" grep -q '^[0-9.]\+$';
				do
					sleep 1;
				done
				return 0;;
			-h|--help)
				echo "Usage: $0 [-h] [-w]";
				echo "";
				echo "-h, --help  print usage information and exit";
				echo "-w, --wait  wait for ATC to be ready, instead of just checking if it is ready";
				echo "-d, --wait  wait for the ATC delivery service to be ready";
				return 0;;
			*)
				echo "Usage: $0 [-h] [-w]" >&2;
				echo "" >&2;
				echo "-h, --help  print usage information and exit" >&2;
				echo "-w, --wait  wait for ATC to be ready, instead of just checking if it is ready" >&2;
				echo "-d, --wait  wait for the ATC delivery service to be ready" >&2;
				return 1;;
		esac
	fi
	curl -skL "$url" >/dev/null 2>&1;
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

export t3cDir="/root/go/src/github.com/apache/trafficcontrol/cache-config";

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

declare -r API_VERSION=4.0
export API_VERSION
export TO_URL="https://localhost:6443"
export TO_USER="admin"
export TO_PASSWORD="twelve12"


# On some shell/system combinations, either or both of these are available as
# shell variables but aren't exported to the execution environment. In others,
# they may just not be set. In any case, trying to set one or both of these to
# certain values - or even at all, on some systems - will fail, so we hope this
# isn't necessary.
if [[ -z "$USER" ]]; then
	USER="$(id -un)";
fi
export USER;

if [[ -z "$UID" ]]; then
	UID="$(id -u)";
fi
export UID;

