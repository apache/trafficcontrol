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
  cd "${REPO_DIR}/experimental/traffic-portal"
  if ! [[ -d Reports ]]; then
    mkdir Reports;
  fi
  if [[ -d nightwatch/junit ]]; then
    mv nightwatch/junit Reports
  fi
  if [[ -d nightwatch/screens ]]; then
    mv nightwatch/screens Reports
  fi
  if [[ -d logs ]]; then
    mv logs Reports
  fi
  if [[ -f "${REPO_DIR}/traffic_ops/traffic_ops_golang" ]]; then
    cp "${REPO_DIR}/traffic_ops/traffic_ops_golang" Reports/to.log;
  fi
  echo "Detailed logs produced info Reports artifact"
  exit 1
}

trap 'onFail "${LINENO}" "${0}"' ERR
set -o errexit -o nounset -o pipefail

to_fqdn="https://localhost:6443"
tp_fqdn="http://localhost:4200"

export PGUSER="traffic_ops"
export PGPASSWORD="twelve"
export PGHOST="localhost"
export PGDATABASE="traffic_ops"
export PGPORT="5432"

to_admin_username="admin"
to_admin_password="twelve12"
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

sudo apt-get install -y --no-install-recommends gettext curl

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

cd "${REPO_DIR}/experimental/traffic-portal"
npx ng serve &

# Wait for tp/to build
timeout 15m bash <<TMOUT
  while ! curl -Lvsk "${tp_fqdn}/api/4.0/ping" >/dev/null 2>&1; do
    echo "waiting for TP/TO server to start on '${tp_fqdn}'"
    sleep 30
  done
TMOUT

npm run e2e:ci
