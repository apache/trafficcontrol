#!/bin/bash
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

fqdn="http://localhost:4444/wd/hub/status"
if ! curl -Lvsk "${fqdn}" >/dev/null 2>&1; then
  echo "Selenium not started on ${fqdn}"
  exit 1
fi

DIVISION="adivision"
REGION="aregion"
PHYS="aloc"
COORD="acoord"
CDN="zcdn"
CG="acg"
export PGUSER="traffic_ops"
export PGPASSWORD="twelve"
export PGHOST="localhost"
export PGDATABASE="traffic_ops"
export PGPORT="5432"

<<QUERY psql
INSERT INTO tm_user (username, local_passwd, role, tenant_id) VALUES ('admin', 'SCRYPT:16384:8:1:vVw4X6mhoEMQXVGB/ENaXJEcF4Hdq34t5N8lapIjDQEAS4hChfMJMzwwmHfXByqUtjmMemapOPsDQXG+BAX/hA==:vORiLhCm1EtEQJULvPFteKbAX2DgxanPhHdrYN8VzhZBNF81NRxxpo7ig720KcrjH1XFO6BUTDAYTSBGU9KO3Q==', 1, 1);
INSERT INTO division(name) VALUES('${DIVISION}');
INSERT INTO region(name, division) VALUES('${REGION}', 1);
INSERT INTO phys_location(name, short_name, region, address, city, state, zip) VALUES('${PHYS}', '${PHYS}', 1, 'some place idk', 'Denver', 'CO', '88888');
INSERT INTO coordinate(name) VALUES('${COORD}');
INSERT INTO cdn(name, domain_name) VALUES('${CDN}', 'infra.ciab.test');

WITH TYPE AS (SELECT id FROM type WHERE name = 'TC_LOC')
INSERT INTO cachegroup(name, short_name, type, coordinate)
SELECT '${CG}', '${CG}', TYPE.id, 1
FROM TYPE;

WITH TYPE AS (SELECT id FROM type WHERE name = 'RIAK'),
PROFILE AS (SELECT id FROM profile WHERE name = 'RIAK_ALL'),
STATUS AS (SELECT id FROM status WHERE name = 'ONLINE'),
PHYS AS (SELECT id FROM phys_location WHERE name = '${PHYS}'),
CDN AS (SELECT id FROM cdn WHERE name = '${CDN}'),
CG AS (SELECT id from cachegroup WHERE name = '${CG}')
INSERT INTO server(host_name, domain_name, cachegroup, type, status, profile, phys_location, cdn_id)
SELECT 'trafficvault', 'infra.ciab.test', CG.ID, TYPE.id, STATUS.id, PROFILE.id, PHYS.id, CDN.id
FROM TYPE
JOIN STATUS ON 1=1
JOIN PROFILE ON 1=1
JOIN PHYS ON 1=1
JOIN CDN ON 1=1
JOIN CG ON 1=1;
QUERY


download_go() {
	. build/functions.sh
	if verify_and_set_go_version; then
		return
	fi
	go_version="$(cat "${GITHUB_WORKSPACE}/GO_VERSION")"
	wget -O go.tar.gz "https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz" --no-verbose
	echo "Extracting Go ${go_version}..."
	<<-'SUDO_COMMANDS' sudo sh
		set -o errexit
    go_dir="$(command -v go | xargs realpath | xargs dirname | xargs dirname)"
		mv "$go_dir" "${go_dir}.unused"
		tar -C /usr/local -xzf go.tar.gz
	SUDO_COMMANDS
	rm go.tar.gz
	go version
}

gray_bg="$(printf '%s%s' $'\x1B' '[100m')";
red_bg="$(printf '%s%s' $'\x1B' '[41m')";
yellow_bg="$(printf '%s%s' $'\x1B' '[43m')";
black_fg="$(printf '%s%s' $'\x1B' '[30m')";
color_and_prefix() {
	color="$1";
	shift;
	prefix="$1";
	normal_bg="$(printf '%s%s' $'\x1B' '[49m')";
	normal_fg="$(printf '%s%s' $'\x1B' '[39m')";
	sed "s/^/${color}${black_fg}${prefix}: /" | sed "s/$/${normal_bg}${normal_fg}/";
}

ciab_dir="${GITHUB_WORKSPACE}/infrastructure/cdn-in-a-box";
trafficvault=trafficvault;
start_traffic_vault() {
	<<-'/ETC/HOSTS' sudo tee --append /etc/hosts
		172.17.0.1    trafficvault.infra.ciab.test
	/ETC/HOSTS

	<<-'BASH_LINES' cat >infrastructure/cdn-in-a-box/traffic_vault/prestart.d/00-0-standalone-config.sh;
		TV_FQDN="${TV_HOST}.${INFRA_SUBDOMAIN}.${TLD_DOMAIN}" # Also used in 02-add-search-schema.sh
		certs_dir=/etc/ssl/certs;
		X509_INFRA_CERT_FILE="${certs_dir}/trafficvault.crt";
		X509_INFRA_KEY_FILE="${certs_dir}/trafficvault.key";

		# Generate x509 certificate
		openssl req -new -x509 -nodes -newkey rsa:4096 -out "$X509_INFRA_CERT_FILE" -keyout "$X509_INFRA_KEY_FILE" -subj "/CN=${TV_FQDN}";

		# Do not wait for CDN in a Box to generate SSL keys
		sed -i '0,/^update-ca-certificates/d' /etc/riak/prestart.d/00-config.sh;

		# Do not try to source to-access.sh
		sed -i '/to-access\.sh\|^to-enroll/d' /etc/riak/{prestart.d,poststart.d}/*
	BASH_LINES

	DOCKER_BUILDKIT=1 docker build "$ciab_dir" -f "${ciab_dir}/traffic_vault/Dockerfile" -t "$trafficvault" 2>&1 |
		color_and_prefix "$gray_bg" "building Traffic Vault";
	echo 'Starting Traffic Vault...';
	docker run \
		--detach \
		--env-file="${ciab_dir}/variables.env" \
		--hostname="${trafficvault}.infra.ciab.test" \
		--name="$trafficvault" \
		--publish=8087:8087 \
		--rm \
		"$trafficvault" \
		/usr/lib/riak/riak-cluster.sh;
}
start_traffic_vault &

sudo apt-get install -y --no-install-recommends gettext \
	ruby ruby-dev libc-dev curl \
	chromium-chromedriver postgresql-client \
	gcc musl-dev

sudo gem update --system && sudo gem install sass compass
sudo npm i -g protractor@^7.0.0 forever bower grunt selenium-webdriver

CONTAINER=$(docker ps | grep "selenium/node-chrome" | awk '{print $1}')
CHROME_VER=$(docker exec "$CONTAINER" google-chrome --version | sed -E 's/.* ([0-9.]+).*/\1/')
sudo webdriver-manager update --gecko false --standalone false --versions.chrome "LATEST_RELEASE_$CHROME_VER"

GOROOT=/usr/local/go
export GOPATH PATH="${PATH}:${GOROOT}/bin"
download_go
GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"

cd "$SRCDIR/trafficcontrol/traffic_ops/traffic_ops_golang"

/usr/local/go/bin/go get -v golang.org/x/net/publicsuffix\
	golang.org/x/crypto/ed25519 \
	golang.org/x/crypto/scrypt \
	golang.org/x/net/idna \
	golang.org/x/net/ipv4 \
	golang.org/x/net/ipv6 \
	golang.org/x/sys/unix \
	golang.org/x/text/secure/bidirule > /dev/null
/usr/local/go/bin/go build . > /dev/null

openssl req -new -x509 -nodes -newkey rsa:4096 -out localhost.crt -keyout localhost.key -subj "/CN=tptests";

resources="$(dirname "$0")"
envsubst <"${resources}/cdn.json" >cdn.conf
cp "${resources}/database.json" database.conf

export $(<"${ciab_dir}/variables.env" sed '/^#/d') # defines TV_ADMIN_USER/PASSWORD
envsubst <"${resources}/riak.json" >riak.conf

truncate --size=0 warning.log error.log # Removes output from previous API versions and makes sure files exist
./traffic_ops_golang --cfg ./cdn.conf --dbcfg ./database.conf -riakcfg riak.conf &
tail -f warning.log 2>&1 | color_and_prefix "${yellow_bg}" 'Traffic Ops' &
tail -f error.log 2>&1 | color_and_prefix "${red_bg}" 'Traffic Ops' &

cd "../../traffic_portal"
npm ci
bower install
grunt dist

cp "${resources}/config.js" ./conf/
touch tp.log access.log
sudo forever --minUptime 5000 --spinSleepTime 2000 -l ./tp.log start server.js &

fqdn="https://localhost:8443/"
while ! curl -Lvsk "${fqdn}api/3.0/ping" >/dev/null 2>&1; do
  echo "waiting for TP/TO server to start on '${fqdn}'"
  sleep 10
done


cd "test/end_to_end"
cp "${resources}/conf.json" .

sudo protractor ./conf.js
CODE=$?

if [ $CODE -ne 0 ]; then
	docker logs -f "$trafficvault" 2>&1 |
		color_and_prefix "$gray_bg" 'Traffic Vault';
  tail -f tp.log 2>&1 | color_and_prefix "${gray_bg}" 'Forever' &
  tail -f access.log 2>&1 | color_and_prefix "${gray_bg}" 'Traffic Portal' &
fi

exit $CODE
