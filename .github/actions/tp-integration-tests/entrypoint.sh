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

onFail() {
  echo "Error on line ${1} of ${2}" >&2;
  if ! [[ -d Reports ]]; then
    mkdir Reports;
  fi
  if [[ -f tp.log ]]; then
    mv tp.log Reports/pm2.log
  fi
  if [[ -f "${REPO_DIR}/traffic_ops/traffic_ops_golang/out.log" ]]; then
    mv "${REPO_DIR}/traffic_ops/traffic_ops_golang/out.log" Reports/to.log
  fi
  docker logs $CHROME_CONTAINER > Reports/chrome.log 2>&1;
  docker logs $HUB_CONTAINER > Reports/hub.log 2>&1;
  echo "Detailed logs produced info Reports artifact"
  exit 1
}

trap 'onFail "${LINENO}" "${0}"' ERR
set -o errexit -o nounset -o pipefail

hub_fqdn="http://localhost:4444/wd/hub/status"
to_fqdn="https://localhost:6443"
tp_fqdn="https://172.18.0.1:8443"

if ! curl -Lvsk "${hub_fqdn}" >/dev/null 2>&1; then
  echo "Selenium not started on ${hub_fqdn}" >&2;
  exit 1
fi

export PGUSER="traffic_ops"
export PGPASSWORD="twelve"
export PGHOST="localhost"
export PGDATABASE="traffic_ops"
export PGPORT="5432"

to_admin_username="$(jq -r '.params.login.username' "${GITHUB_WORKSPACE}/traffic_portal/test/integration/config.json")"
to_admin_password="$(jq -r '.params.login.password' "${GITHUB_WORKSPACE}/traffic_portal/test/integration/config.json")"
password_hash="$(<<PYTHON_COMMANDS PYTHONPATH="${GITHUB_WORKSPACE}/traffic_ops/install/bin" python
from _postinstall import hash_pass
print(hash_pass('${to_admin_password}'))
PYTHON_COMMANDS
)"
<<QUERY psql
INSERT INTO tm_user (username, role, tenant_id, local_passwd)
	VALUES ('${to_admin_username}', (
		SELECT id
		FROM "role"
		WHERE "name" = 'admin'
	), (
		SELECT id
		FROM tenant
		WHERE "name" = 'root'
	),
    '${password_hash}'
  );
QUERY

sudo useradd trafops

ciab_dir="${GITHUB_WORKSPACE}/infrastructure/cdn-in-a-box";
openssl rand 32 | base64 | sudo tee /aes.key

sudo apt-get install -y --no-install-recommends gettext curl

sudo npm i -g pm2 grunt sass

CHROME_CONTAINER=$(docker ps -qf name=chrome)
HUB_CONTAINER=$(docker ps -qf name=hub)
CHROME_VER=$(docker exec "$CHROME_CONTAINER" google-chrome --version | grep -Eo '[0-9.]+')

export GOPATH="${HOME}/go"
readonly ORG_DIR="$GOPATH/src/github.com/apache"
readonly REPO_DIR="${ORG_DIR}/trafficcontrol"
resources="$(dirname "$0")"
if [[ ! -e "$REPO_DIR" ]]; then
	mkdir -p "$ORG_DIR"
	cd
	mv "${GITHUB_WORKSPACE}" "${REPO_DIR}/"
	ln -s "$REPO_DIR" "${GITHUB_WORKSPACE}"
fi

pushd "${REPO_DIR}/traffic_ops/traffic_ops_golang"
if  [[ ! -d "${GITHUB_WORKSPACE}/vendor/golang.org" ]]; then
  go mod vendor
fi
go build .

openssl req -new -x509 -nodes -newkey rsa:4096 -out localhost.crt -keyout localhost.key -subj "/CN=tptests";

envsubst <"${resources}/cdn.json" >cdn.conf
cp "${resources}/database.json" database.conf

truncate -s0 out.log
./traffic_ops_golang --cfg ./cdn.conf --dbcfg ./database.conf >out.log 2>&1 &
popd

pushd "${REPO_DIR}/traffic_portal"
npm ci
grunt dist

cp "${resources}/config.js" ./conf/
truncate -s0 out.log
sudo pm2 start server.js -l out.log
popd

cd "${REPO_DIR}/traffic_portal/test/integration"
npm ci

npx webdriver-manager update --gecko false --versions.chrome "LATEST_RELEASE_$CHROME_VER"

jq " .capabilities.chromeOptions.args = [
    \"--headless\",
    \"--no-sandbox\",
    \"--disable-gpu\",
    \"--ignore-certificate-errors\"
  ] | .params.apiUrl = \"${to_fqdn}/api/5.0\" | .params.baseUrl =\"${tp_fqdn}\"
  | .capabilities[\"goog:chromeOptions\"].w3c = false | .capabilities.chromeOptions.w3c = false" \
  config.json > config.json.tmp && mv config.json.tmp config.json

# Wait for tp/to build
timeout 5m bash <<TMOUT
  while ! curl -Lvsk "${tp_fqdn}/api/5.0/ping" >/dev/null 2>&1; do
    echo "waiting for TP/TO server to start on '${tp_fqdn}'"
    sleep 10
  done
TMOUT

npm test -- --params.baseUrl="${tp_fqdn}" --params.apiUrl="${to_fqdn}/api/5.0"
