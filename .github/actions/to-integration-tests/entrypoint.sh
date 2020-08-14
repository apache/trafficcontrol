#!/bin/sh -l
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

GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"

cd "$SRCDIR/trafficcontrol/traffic_ops/app/db"

echo 'version: "1.0"
name: dbconf.yml

test:
  driver: postgres
  open: host=postgres port=5432 user=traffic_ops password=twelve dbname=traffic_ops sslmode=disable

' > dbconf.yml

psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./create_tables.sql >/dev/null
goose --env=test --path="$PWD" up
psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./seeds.sql >/dev/null
psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./patches.sql >/dev/null

cd "$SRCDIR/trafficcontrol/traffic_ops/traffic_ops_golang"


/usr/local/go/bin/go get ./...
/usr/local/go/bin/go build .

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

cat <<-EOF >cdn.conf
{
	"hypnotoad" : {
		"listen" : [
			"https://not-a-real-host.test:1?cert=$PWD/localhost.crt&key=$PWD/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED"
		],
		"user" : "trafops",
		"group" : "trafops",
		"heartbeat_timeout" : 20,
		"pid_file" : "/var/run/traffic_ops.pid",
		"workers" : 12
	},
	"use_ims": true,
	"traffic_ops_golang" : {
		"insecure": true,
		"port" : "6443",
		"log_location_error": "stderr,
		"log_location_warning": "stderr",
		"log_location_info": null,
		"log_location_debug": null,
		"log_location_event": null,
		"max_db_connections": 20,
		"db_conn_max_lifetime_seconds": 60,
		"db_query_timeout_seconds": 20,
		"supported_ds_metrics": [ "kbps", "tps_total", "tps_2xx", "tps_3xx", "tps_4xx", "tps_5xx" ]
	},
	"cors" : {
		"access_control_allow_origin" : "*"
	},
	"to" : {
		"base_url" : "https://localhost",
		"email_from" : "no-reply@traffic-ops.test",
		"no_account_found_msg" : "A Traffic Ops user account is required for access. Please contact your Traffic Ops user administrator."
	},
	"portal" : {
		"base_url" : "https://not-a-real-host.test/#!/",
		"email_from" : "no-reply@traffic-portal.test",
		"pass_reset_path" : "user",
		"user_register_path" : "user"
	},
	"secrets" : [
		"blahblah"
	],
	"geniso" : {
		"iso_root_path" : "/opt/traffic_ops/app/public"
	},
	"inactivity_timeout" : 60,
	"smtp" : {
		"enabled" : false
	},
	"InfluxEnabled": false
}
EOF

cat <<-EOF >database.conf
{
	"description": "Local PostgreSQL database on port 5432",
	"dbname": "traffic_ops",
	"hostname": "postgres",
	"user": "traffic_ops",
	"password": "twelve",
	"port": "5432",
	"ssl": false,
	"type": "Pg"
}
EOF

./traffic_ops_golang --cfg ./cdn.conf --dbcfg ./database.conf >out.log 2>err.log &

cd ../testing/api/v1

makeCFG() {
	cat <<-EOF >traffic-ops-test.conf
	{
		"default": {
			"logLocations": {
				"debug": "stderr",
				"error": "stderr",
				"event": "stderr",
				"info": "stderr",
				"warning": "stderr"
			},
			"session": {
				"timeoutInSecs": 60
			}
		},
		"trafficOps": {
			"URL": "https://localhost:6443",
			"password": "twelve",
			"users": {
				"disallowed": "disallowed",
				"operations": "operations",
				"admin": "admin",
				"federation": "federation",
				"portal": "portal",
				"readOnly": "readOnly",
				"extension": "extension"
			}
		},
		"trafficOpsDB": {
			"dbname": "traffic_ops",
			"description": "Test database to_test",
			"hostname": "postgres",
			"password": "twelve",
			"port": "5432",
			"type": "Pg",
			"user": "traffic_ops"
		}
	}
	EOF
}

CODE="0"

makeCFG
/usr/local/go/bin/go test -v --cfg ./traffic-ops-test.conf
if [ "$?" -gt "0" ]; then
	CODE="1"
fi
rm traffic-ops-test.conf

cd ../v2
makeCFG
 /usr/local/go/bin/go test -v --cfg ./traffic-ops-test.conf
if [ "$?" -gt "0" ]; then
	CODE="1"
fi
rm traffic-ops-test.conf

cd ../v3
makeCFG
 /usr/local/go/bin/go test -v --cfg ./traffic-ops-test.conf
if [ "$?" -gt "0" ]; then
	CODE="1"
fi
rm traffic-ops-test

cat ../../../traffic_ops_golang/out.log
cat ../../../traffic_ops_golang/err.log

exit "$CODE"
