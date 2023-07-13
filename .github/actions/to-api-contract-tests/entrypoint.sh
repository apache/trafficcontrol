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

set -o errexit -o nounset -o pipefail

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

openssl rand 32 | base64 | sudo tee /aes.key

sudo apt-get install -y --no-install-recommends gettext

cd "traffic_ops/traffic_ops_golang"

if  [[ ! -d "${GITHUB_WORKSPACE}/vendor/golang.org" ]]; then
	go mod vendor
fi
go build .

echo "
-----BEGIN CERTIFICATE-----
MIIDtTCCAp2gAwIBAgIJAJgQuE9T48+gMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTcwNTA5MDIyNTI0WhcNMTgwNTA5MDIyNTI0WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEA1OsuQ71qaSZ4ivGnh4YryeQEHMn5wZLX7kWYB637ssyOJnkU1ZGs93JM
XJADzsjmssP6icSDhV2JPgDDYzx1eBJt6y3vHI7L3AdGfQJj+4FFABKR8moqpc1J
WIMGnPzO6DeEc8irf0qxSh+yvuFX0j6oS8oCqiRxz5+HL2wEGWmrgr37JY4/bs7o
4CMY19Ru1dP2Fr292HEIqCEnLTOuaHSWAEWx1Tm93kT9sXbw/SG2JTLQSX80biFL
7foJeoGWLls2reTCYTprzWFaMu3x9I8HLtf4VIN44rtvo5N20KYgjGqvPjFGPljL
yrgB8rXSCpH3M4AbazxD8fZKbdORawIDAQABo4GnMIGkMB0GA1UdDgQWBBT6zEpf
DYbYCI3Bu82+Q5SmI+/7ojB1BgNVHSMEbjBsgBT6zEpfDYbYCI3Bu82+Q5SmI+/7
oqFJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNV
BAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJAJgQuE9T48+gMAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAGLs1NcYNtUgN6FuMb6/UskEWLTKwfno
NBtNdIbcZP3HmJHwruLWCeqj6HIWJC87EqmPTIYPdem3SAN1L20fWpzm7AB7av+2
wTCAPVP0punF/IouSb6fyo8fdG1a104Mge4iy/Sf2uf09NEv08sfVdB4P0tKRRlg
5KChhmspdPP7fmPXyghm4IC0Seknmh6IlVOnALXLU5OoCLHTie5Hjv4Tm8Xu0oBA
dIH/cPu2/w5SAIVq9CtcsdglS0ZsCAv4W2YieuSLPf5xuI0q/5lFZGNoDpIWJldx
Y2IpnoNCrHEAxijP5ctPawsxkSt2PmQ5uNNL7TbMudc3hZzOpTPkGoo=
-----END CERTIFICATE-----
" > localhost.crt

echo "-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1OsuQ71qaSZ4ivGnh4YryeQEHMn5wZLX7kWYB637ssyOJnkU
1ZGs93JMXJADzsjmssP6icSDhV2JPgDDYzx1eBJt6y3vHI7L3AdGfQJj+4FFABKR
8moqpc1JWIMGnPzO6DeEc8irf0qxSh+yvuFX0j6oS8oCqiRxz5+HL2wEGWmrgr37
JY4/bs7o4CMY19Ru1dP2Fr292HEIqCEnLTOuaHSWAEWx1Tm93kT9sXbw/SG2JTLQ
SX80biFL7foJeoGWLls2reTCYTprzWFaMu3x9I8HLtf4VIN44rtvo5N20KYgjGqv
PjFGPljLyrgB8rXSCpH3M4AbazxD8fZKbdORawIDAQABAoIBAQCYOINc/qCDCHQJ
sfa511ya/B9MjcG3eMpTmQG2C9b033WJX+tbPMjSJ68cRgHS5qK4j5AgypPU1yh1
YYpO+jxpWZOoHbDjU9u/NJxaZ0kf2C2CfcRF8U0IOJoFY7doqP0r2/Uf6glh+f6C
JeNewDBPKWictpHtHh0X+M9nQew0VZ7slXnV+IwUxiWYEtiIjwMyzSmfDEnN3ix5
fuVQLvVaq+bbqXj2rMpJWFj7zMsG5HRePQl2kQGtMYLCIalnJIQs5jQn2YsliNyy
fQiwWnU0wkrLlmkhlivlISRDtP35WQgF8ObsoQ3LZXRflB0C7U7zEl7Dj3Vi7WXr
jsRZC4dxAoGBAPwuPdtc9gSNKjn8RnqfEJjdSo1qdLbGvRcSJNy4/kEEFECJXkeO
mV/aklCi39cxAaIjVdTQ1XN67RMxgdekCI2Eg8h4RdvwgB/tAO+C3ExzMSOA1IcZ
tWuwIA2YnaFF9Sla9iJqxgtoGlaqm4VTUM/IdZqlzsP08pfNq7bXPsr9AoGBANgk
tkovf1Y0O4lBHX3eVLnHXForxEZh8bGHwuJJWWzb0ZFcXrrSd8FSycZrR28v4sdQ
WSSVPz3Op95HoTVXVL9EJcZ+MTnHaoCHbYBkrGTlGviu5Fl2V5EbrN7R7CdxJeem
HOU4shTy1acMPgf8sT17ykkXhVeUhSK2Jg6fZn6HAoGAI4/51SeE4htuKwMyhTRN
SOFcFBlBIE1ieRBr9lx4Ln7+xCMbEog/hM7z9z8gxd35Vv4YqoxQrZpWOHCw2NIf
CqX3V5vubhe6WcY4bY5MttM/yLvwPKUZeng57PDqucV9zzkuoKfiCdXCcRpaGDEp
okOooghj4ip204WDg6NTDZkCgYEAwZTfzsGLgmF1kRBIoZqmt1zeUcQxHfhKx32Y
BaM7/EtD/rSEAz7NEtBa9uLOL77rlSdZL3KcGXck0efFckitFkCqtIQBAoaf1E12
vS9tV0/6QBAjZByhgM0Qnt/Uad7k2/vilUmZ9TkoMVy9kdm3xCFCowP14OKb+uK4
YxBQc7ECgYEAm7eVtNlPHYmC54FU2bLucryNMLmu9I8O6zvbK5sxiMdtlh7OjaUB
RQS5iVc0iTacDZTGh7eqNzgGplj76pWGHeZUy0xIj/ZIRu2qOy0v+ffqfX1wCz7p
A22D22wvfs7CE3cUz/8UnvLM3kbTTu1WbbBbrHjAV47sAHjW/ckTqeo=
-----END RSA PRIVATE KEY-----
" > localhost.key

resources="$(dirname "$0")"
envsubst <"${resources}/cdn.json" >cdn.conf
cp "${resources}/database.json" database.conf

truncate --size=0 traffic.ops.log # Removes output from previous API versions and makes sure files exist
./traffic_ops_golang --cfg ./cdn.conf --dbcfg ./database.conf &

cd "../testing/api_contract/v$INPUT_VERSION"

pytest -rA --to-user admin --to-password twelve12 --to-url https://localhost:6443/api/4.0
