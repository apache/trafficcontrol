#!/usr/bin/env bash
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

# Script for running the Dockerfile for Traffic Ops.
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set, ordinarily by `docker run -e` arguments:
# DB_SERVER
# DB_PORT
# DB_ROOT_PASS
# DB_USER
# DB_USER_PASS
# DB_NAME
# ADMIN_USER
# ADMIN_PASS
# DOMAIN

# TODO:  Unused -- should be removed?  TRAFFIC_VAULT_PASS

# Setting the monitor shell option enables job control, which we need in order
# to bring traffic_ops_golang back to the foreground.
set -o errexit -o monitor -o pipefail -o xtrace;

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS)
for v in $envvars; do
	if [[ -z $$v ]]; then
		echo "$v is unset" >&2;
		exit 1;
	fi
done

export PATH="$PATH:/opt/traffic_ops/go/bin"

/set-dns.sh
/insert-self-into-dns.sh

/set-to-ips-from-dns.sh

# Source to-access functions and FQDN vars
source /to-access.sh

# Create SSSL certificates and trust the shared CA.
source /generate-certs.sh

# copy contents of /ca to /export/ssl
# update the permissions
# TODO: figure out how to do this without all the 'chmod 777's
mkdir -p "$X509_CA_PERSIST_DIR";
chmod 777 "$X509_CA_PERSIST_DIR";
chmod -R a+rw "$X509_CA_PERSIST_DIR";

if [ -r "$X509_CA_PERSIST_ENV_FILE" ]; then
	umask "$X509_CA_UMASK";
	mkdir -p "$X509_CA_DIR";
	chmod 777 "$X509_CA_DIR";
	chmod -R a+rw "$X509_CA_DIR";
	rsync -a "$X509_CA_PERSIST_DIR/" "$X509_CA_DIR/";
	sync;
	echo "PERSIST CERTS FROM $X509_CA_PERSIST_DIR to $X509_CA_DIR";
	sleep 4;
	source "$X509_CA_ENV_FILE";
elif x509v3_init; then
	umask $X509_CA_UMASK;
	x509v3_create_cert "$INFRA_SUBDOMAIN" "$INFRA_FQDN";
	for ds in $DS_HOSTS; do
		x509v3_create_cert "$ds" "$ds.$CDN_FQDN";
	done
	echo "X509_GENERATION_COMPLETE=\"YES\"" >> "$X509_CA_ENV_FILE";
	x509v3_dump_env
	# Save newly generated certs for future restarts.
	rsync -av "$X509_CA_DIR/" "$X509_CA_PERSIST_DIR/";
	chmod 777 "$X509_CA_PERSIST_DIR";
	chmod -R a+rw "$X509_CA_DIR";
	sync;
	echo "GENERATE CERTS FROM $X509_CA_DIR to $X509_CA_PERSIST_DIR";
	sleep 4;
fi

chown -R trafops:trafops "$X509_CA_PERSIST_DIR";
chmod -R a+rw "$X509_CA_PERSIST_DIR";

# Write config files
. /config.sh;

pg_isready=$(rpm -ql postgresql13 | grep bin/pg_isready);
if [[ ! -x $pg_isready ]]; then
	echo "Can't find pg_ready in postgresql13" >&2;
	echo "PATH: $PATH" >&2;
	find / -name "*postgresql*";
	exit 1;
fi

while ! $pg_isready -h "$DB_SERVER" -p "$DB_PORT" -d "$DB_NAME"; do
	echo "waiting for db on $DB_SERVER:$DB_PORT";
	sleep 3;
done

cd /opt/traffic_ops/app;

# Add admin user -- all other users should be created using the API
/adduser.pl "$TO_ADMIN_USER" "$TO_ADMIN_PASSWORD" "admin" "root" | psql -v ON_ERROR_STOP=1 -U "$DB_USER" -h "$DB_SERVER" -d "$DB_NAME";

(
maxtries=10
for ((tries = 0; tries < maxtries; tries++)); do
	if nc -zvw2 "$SMTP_FQDN" "$SMTP_PORT"; then
		echo "${SMTP_FQDN}:${SMTP_PORT} was found.";
		break;
	fi;
	echo "waiting for ${SMTP_FQDN}:${SMTP_PORT}";
	sleep 3;
done
if (( tries == maxtries )); then
	echo "SMTP service was not found at ${SMTP_FQDN}:${SMTP_PORT} after ${maxtries} tries. Skipping...";
fi
)

cd /opt/traffic_ops/app;

CDNCONF=/opt/traffic_ops/app/conf/cdn.conf
DBCONF=/opt/traffic_ops/app/conf/production/database.conf
RIAKCONF=/opt/traffic_ops/app/conf/production/riak.conf
mkdir -p /var/log/traffic_ops
touch "$TO_LOG_ERROR" "$TO_LOG_WARNING" "$TO_LOG_INFO" "$TO_LOG_DEBUG" "$TO_LOG_EVENT"
tail -qf "$TO_LOG_ERROR" "$TO_LOG_WARNING" "$TO_LOG_INFO" "$TO_LOG_DEBUG" "$TO_LOG_EVENT" &

traffic_ops_golang_command=(./bin/traffic_ops_golang -cfg "$CDNCONF" -dbcfg "$DBCONF" -riakcfg "$RIAKCONF");
if [[ "$TO_DEBUG_ENABLE" == true ]]; then
	traffic_ops_golang_command=(dlv '--accept-multiclient' '--continue' '--listen=:2345' '--headless=true' '--api-version=2' exec
		"${traffic_ops_golang_command[0]}" -- "${traffic_ops_golang_command[@]:1}");
fi;
"${traffic_ops_golang_command[@]}" &

until [[ -f "$ENROLLER_DIR/enroller-started" ]]; do
	echo "waiting for enroller";
	sleep 3;
done

/trafficops-init.sh;

# enroll in the background so traffic_ops_golang can run in foreground
TO_USER=$TO_ADMIN_USER TO_PASSWORD=$TO_ADMIN_PASSWORD to-enroll $(hostname -s) &

# Add initial data to traffic ops
/trafficops-init.sh

to-enroll "to" ALL;

while true; do
	echo "Verifying that edge was associated to delivery service...";

	cachegroup="$(to-get "api/${TO_API_VERSION}/servers?hostName=edge" 2>/dev/null | jq -r -c '.response[0]|.cachegroup')"
	ds_name=$(to-get "api/${TO_API_VERSION}/deliveryservices" 2>/dev/null | jq -r -c '.response[] | select(.cdnName == "'"$CDN_NAME"'").xmlId')
	topology=$(to-get "api/${TO_API_VERSION}/deliveryservices" 2>/dev/null | jq -r -c '.response[] | select(.cdnName == "'"$CDN_NAME"'").topology')
	topology_node="$(to-get "/api/${TO_API_VERSION}/topologies?name=${topology}" | jq -r '.response[].nodes[] | select(.cachegroup == "'"$cachegroup"'") | .cachegroup')"

  if [[ -n "$topology_node" ]] ; then
	break
  fi

	sleep 2;
done

### Add SSL keys for demo1 delivery service
until [[ -s "$X509_DEMO1_CERT_FILE" && -s "$X509_DEMO1_REQUEST_FILE" && -s "$X509_DEMO1_KEY_FILE" ]]; do
	echo "Waiting on X509_DEMO1 files to exist";
	sleep 3;
	source "$X509_CA_ENV_FILE";
done
to-add-sslkeys "$CDN_NAME" "$ds_name" "*.demo1.mycdn.ciab.test" "$X509_DEMO1_CERT_FILE" "$X509_DEMO1_REQUEST_FILE" "$X509_DEMO1_KEY_FILE";

### Automatic Queue/Snapshot ###
if [[ "$AUTO_SNAPQUEUE_ENABLED" = true ]]; then
	# AUTO_SNAPQUEUE_SERVERS should be a comma delimited list of expected docker service names to be enrolled - see varibles.env
	to-auto-snapqueue "$AUTO_SNAPQUEUE_SERVERS" "$CDN_NAME";
fi

tail -f /dev/null; # Keeps the container running indefinitely. The container health check (see dockerfile) will report whether Traffic Ops is running.
