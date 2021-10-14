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
    mv tp.log Reports/forever.log
  fi
  if [[ -f access.log ]]; then
    mv access.log Reports/tp-access.log
  fi
  if [[ -f out.log ]]; then
    mv out.log Reports/node.log
  fi
  docker logs $CHROMIUM_CONTAINER > Reports/chromium.log 2>&1;
  docker logs $HUB_CONTAINER > Reports/hub.log 2>&1;
  if [[ -f "${REPO_DIR}/traffic_ops/traffic_ops_golang" ]]; then
    cp "${REPO_DIR}/traffic_ops/traffic_ops_golang" Reports/to.log;
  fi
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
import _postinstall
print(_postinstall.hash_pass('${to_admin_password}'))
PYTHON_COMMANDS
)"
<<QUERY psql
INSERT INTO tm_user (username, role, tenant_id, local_passwd)
  VALUES ('${to_admin_username}', 1, 1,
    '${password_hash}'
  );
QUERY

sudo useradd trafops

ciab_dir="${GITHUB_WORKSPACE}/infrastructure/cdn-in-a-box";
openssl rand 32 | base64 | sudo tee /aes.key

sudo apt-get install -y --no-install-recommends gettext \
	ruby ruby-dev libc-dev curl \
	gcc

sudo gem install sass compass
sudo npm i -g forever grunt

CHROMIUM_CONTAINER=$(docker ps -qf name=chromium)
HUB_CONTAINER=$(docker ps -qf name=hub)
CHROMIUM_VER=$(docker exec "$CHROMIUM_CONTAINER" chromium --version | grep -Eo '[0-9.]+')

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

to_build() {
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
}

tp_build() {
  pushd "${REPO_DIR}/traffic_portal"
  npm ci
  grunt dist

  cp "${resources}/config.js" ./conf/
  touch tp.log access.log out.log err.log
  sudo forever --minUptime 5000 --spinSleepTime 2000 -f start server.js >out.log 2>&1 &
  popd
}

to_build
tp_build

cd "${REPO_DIR}/traffic_portal/test/integration"
npm ci

./node_modules/.bin/webdriver-manager update --gecko false --versions.chrome "LATEST_RELEASE_$CHROMIUM_VER"

jq " .capabilities.chromeOptions.args = [
    \"--headless\",
    \"--no-sandbox\",
    \"--disable-gpu\",
    \"--ignore-certificate-errors\"
  ] | .params.apiUrl = \"${to_fqdn}/api/4.0\" | .params.baseUrl =\"${tp_fqdn}\"
  | .capabilities[\"goog:chromeOptions\"].w3c = false | .capabilities.chromeOptions.w3c = false" \
  config.json > config.json.tmp && mv config.json.tmp config.json

npm run build

# Wait for tp/to build
timeout 5m bash <<TMOUT
  while ! curl -Lvsk "${tp_fqdn}/api/4.0/ping" >/dev/null 2>&1; do
    echo "waiting for TP/TO server to start on '${tp_fqdn}'"
    sleep 10
  done
TMOUT

npm test -- --params.baseUrl="${tp_fqdn}" --params.apiUrl="${to_fqdn}/api/4.0"
