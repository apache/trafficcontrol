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
set -e;

cd "$(dirname "${BASH_SOURCE[0]}")";
readonly MY_DIR="$(pwd)";

help_string="$(<<-'HELP_STRING' cat
	Usage: ./postinstall.test.h [
	    -b        Explicitly set the path to the Python binary as this value
	    -h, ?     Print this help text and exit
HELP_STRING
)"

while getopts :hb: opt; do
	case "$opt" in
		b) python_bin="$OPTARG";;
		h) echo "$help_string" && exit;;
		?) echo "$help_string" && exit;;
		*) echo "Invalid flag received: ${OPTARG}" >&2 && echo "$help_string" && exit 1;;
	esac;
done;

python_bin="${python_bin:-/usr/bin/python3}";

if [[ ! -x "$python_bin" ]]; then
	echo "Python 3.6+ is required to run - or test - _postinstall.py" >&2;
	exit 1;
fi

readonly TO_PASSWORD=twelve;
readonly ROOT_DIR="$(mktemp -d)";

trap 'rm -rf $ROOT_DIR' EXIT;

"$python_bin" <<EOF;
import importlib
import sys
from os.path import dirname, join
from _postinstall import Scrypt

passwd = '${TO_PASSWORD}'
n = 2 ** 10
r_val = 1
p_val = 1
dklen = 2 ** 4
salt = bytearray([196, 187, 115, 30, 109, 244, 168, 124, 70, 67, 229, 123, 156, 3, 138, 243, 234, 79, 79, 31, 67, 239, 249, 177, 237, 240, 201, 216, 81, 116, 186, 172, 153, 99, 240, 184, 186, 0, 119, 34, 165, 220, 3, 201, 104, 13, 13, 189, 135, 76, 160, 6, 206, 154, 124, 78, 112, 243, 132, 30, 48, 223, 224, 28])
scrypt = Scrypt(password=passwd.encode(), salt=salt, cost_factor=n, block_size_factor=r_val,
				parallelization_factor=p_val, key_length=dklen)

expected_block = [82000861, 2203666842, 4001293736, 627876473, 3101038348, 376175724, 2967675936, 3143524608, 1069098580, 1894075103, 3699786793, 3537442772, 3575269184, 2926196224, 913960627, 2079499993]
actual_block = [2245251288, 1072667772, 4071019211, 2090053191, 2877361598, 1101440729, 1502049634, 3905719376, 3112080378, 1388114151, 3517514506, 1152690600, 2085938056, 2696735995, 3835186347, 283826820]
scrypt.salsa20(actual_block)
if expected_block != actual_block:
	print('Expected {expected_block} for salsa20 result, got {actual_block}'.format(expected_block=expected_block, actual_block=actual_block), file=sys.stderr)
	exit(1)

input_block = [1923378, 729355550, 2408212191, 579221939, 681409774, 1765430015, 3846256959, 831940078, 1480976199, 2878095125, 4245323720, 2776886825, 3332759976, 3497079966, 3107631655, 3763839506, 1283955177, 2851514107, 1743501900, 1888209181, 3387403441, 2898469985, 3685946334, 2122268467, 2234902587, 2934192414, 2528543680, 3247696936, 4144265372, 1687923239, 1573958329, 422403479]
expected_block = [54040099, 3246390556, 3905565410, 4170358448, 2569315507, 3679433373, 2964493607, 3621375783, 318358481, 2014381982, 3240374105, 3569092356, 3150068788, 569153936, 2099549087, 2807540417, 2384835523, 4053238240, 1126008925, 1477842924, 1740405559, 1762470512, 2159908599, 1049875013, 2630682622, 1368095319, 1753173294, 3987760372, 3175003396, 1324304335, 775131569, 2728051478]
actual_block = scrypt.block_mix(input_block)
if expected_block != actual_block:
	print('Expected {expected_block} for block mix result, got {actual_block}'.format(expected_block=expected_block, actual_block=actual_block), file=sys.stderr)
	exit(1)

expected_digest = bytearray([86, 124, 148, 28, 117, 181, 239, 64, 228, 6, 247, 83, 210, 88, 43, 144])
actual_digest = scrypt.derive()
if expected_digest != actual_digest:
	print('Expected {expected_digest} for scrypt, got {actual_digest}'.format(expected_digest=expected_digest, actual_digest=actual_digest), file=sys.stderr)
	exit(1)
EOF

mkdir -p "$ROOT_DIR/etc/pki/tls/certs";
mkdir "$ROOT_DIR/etc/pki/tls/private";
mkdir -p "$ROOT_DIR/opt/traffic_ops/app/public/routing";
mkdir "$ROOT_DIR/opt/traffic_ops/app/db";
mkdir "$ROOT_DIR/opt/traffic_ops/app/db/trafficvault";
mkdir -p "$ROOT_DIR/opt/traffic_ops/app/conf/production";
cat > "$ROOT_DIR/opt/traffic_ops/app/conf/cdn.conf" <<EOF
{
	"traffic_ops_golang": {
    "cert" : "$ROOT_DIR/etc/pki/tls/certs/localhost.crt",
    "key"  : "$ROOT_DIR/etc/pki/tls/private/localhost.key"
  }
}
EOF

"$python_bin" <<TESTS 2>/dev/null | tee -a "${ROOT_DIR}/stdout";
import subprocess
import sys
import _postinstall
from os.path import dirname, join

download_tool = '/does/not/exist'
root = '${ROOT_DIR}'

_postinstall.exec_psql('N/A', 'N/A', '--version')
TESTS

mkdir -p "$ROOT_DIR/opt/traffic_ops/install/data/json";
mkdir "$ROOT_DIR/opt/traffic_ops/install/bin";

# defaults.json is used as input into the `--cfile` option of _postinstall.py
# for testing purposes
cat <<- EOF > "$ROOT_DIR/defaults.json"
{
	"/opt/traffic_ops/app/conf/production/database.conf": [
		{
			"Database type": "Pg",
			"config_var": "type",
			"hidden": false
		},
		{
			"Database name": "traffic_ops",
			"config_var": "dbname",
			"hidden": false
		},
		{
			"Database server hostname IP or FQDN": "localhost",
			"config_var": "hostname",
			"hidden": false
		},
		{
			"Database port number": "5432",
			"config_var": "port",
			"hidden": false
		},
		{
			"Traffic Ops database user": "traffic_ops",
			"config_var": "user",
			"hidden": false
		},
		{
			"Password for Traffic Ops database user": "${TO_PASSWORD}",
			"config_var": "password",
			"hidden": true
		}
	],
	"/opt/traffic_ops/app/conf/production/tv.conf": [
		{
			"Database type": "Pg",
			"config_var": "type",
			"hidden": false
		},
		{
			"Database name": "traffic_vault",
			"config_var": "dbname",
			"hidden": false
		},
		{
			"Database server hostname IP or FQDN": "localhost",
			"config_var": "hostname",
			"hidden": false
		},
		{
			"Database port number": "5432",
			"config_var": "port",
			"hidden": false
		},
		{
			"Traffic Ops database user": "traffic_vault",
			"config_var": "user",
			"hidden": false
		},
		{
			"Password for Traffic Ops database user": "${TO_PASSWORD}",
			"config_var": "password",
			"hidden": true
		}
	],
	"/opt/traffic_ops/app/conf/cdn.conf": [
		{
			"Generate a new secret?": "yes",
			"config_var": "genSecret",
			"hidden": false
		},
		{
			"Number of secrets to keep?": "1",
			"config_var": "keepSecrets",
			"hidden": false
		},
		{
			"Port to serve on?": "443",
			"config_var": "port",
			"hidden": false
		},
		{
			"Number of workers?": "12",
			"config_var": "workers",
			"hidden": false
		},
		{
			"Traffic Ops url?": "http://localhost:3000",
			"config_var": "base_url",
			"hidden": false
		},
		{
			"ldap.conf location?": "/opt/traffic_ops/app/conf/ldap.conf",
			"config_var": "ldap_conf_location",
			"hidden": false
		}
	],
	"/opt/traffic_ops/app/conf/ldap.conf": [
		{
			"Do you want to set up LDAP?": "yes",
			"config_var": "setupLdap",
			"hidden": false
		},
		{
			"LDAP server hostname": "ldaps://ad.cdn.site:3269",
			"config_var": "host",
			"hidden": false
		},
		{
			"LDAP Admin DN": "contact@cdn.site",
			"config_var": "admin_dn",
			"hidden": false
		},
		{
			"LDAP Admin Password": "${TO_PASSWORD}",
			"config_var": "admin_pass",
			"hidden": true
		},
		{
			"LDAP Search Base": "dc=cdn,dc=site",
			"config_var": "search_base",
			"hidden": false
		},
		{
			"LDAP Search Query": "(&(objectCategory=person)(objectClass=user)(sAMAccountName=%s))",
			"config_var": "search_query",
			"hidden": false
		},
		{
			"LDAP Skip TLS verify": "True",
			"config_var": "insecure",
			"hidden": false
		},
		{
			"LDAP Timeout Seconds": "120",
			"config_var": "ldap_timeout_secs",
			"hidden": false
		}
	],
	"/opt/traffic_ops/install/data/json/users.json": [
		{
			"Administration username for Traffic Ops": "admin",
			"config_var": "tmAdminUser",
			"hidden": false
		},
		{
			"Password for the admin user": "twelve",
			"config_var": "tmAdminPw",
			"hidden": true
		}
	],
	"/opt/traffic_ops/install/data/profiles/": [
		{
			"Add custom profiles?": "no",
			"config_var": "custom_profiles",
			"hidden": false
		}
	],
	"/opt/traffic_ops/install/data/json/openssl_configuration.json": [
		{
			"Do you want to generate a certificate?": "yes",
			"config_var": "genCert",
			"hidden": false
		},
		{
			"Country Name (2 letter code)": "US",
			"config_var": "country",
			"hidden": false
		},
		{
			"State or Province Name (full name)": "Colorado",
			"config_var": "state",
			"hidden": false
		},
		{
			"Locality Name (eg, city)": "Denver",
			"config_var": "locality",
			"hidden": false
		},
		{
			"Organization Name (eg, company)": "Comcast",
			"config_var": "company",
			"hidden": false
		},
		{
			"Organizational Unit Name (eg, section)": "Viper",
			"config_var": "org_unit",
			"hidden": false
		},
		{
			"Common Name (eg, your name or your server's hostname)": "cdn",
			"config_var": "common_name",
			"hidden": false
		},
		{
			"RSA Passphrase": "testquest",
			"config_var": "rsaPassword",
			"hidden": true
		}
	],
	"/opt/traffic_ops/install/data/json/profiles.json": [
		{
			"Traffic Ops url": "https://localhost",
			"config_var": "tm.url",
			"hidden": false
		},
		{
			"Human-readable CDN Name. (No whitespace, please)": "kabletown_cdn",
			"config_var": "cdn_name",
			"hidden": false
		},
		{
			"DNS sub-domain for which your CDN is authoritative": "cdn1.kabletown.net",
			"config_var": "dns_subdomain",
			"hidden": false
		}
	]
}
EOF

"$python_bin" "$MY_DIR/_postinstall.py" --no-root --root-directory="$ROOT_DIR" --no-restart-to --no-database --ops-user="$(whoami)" --ops-group="$(id -gn)" --automatic --cfile="$ROOT_DIR/defaults.json" --debug > >(tee -a "$ROOT_DIR/stdout") 2> >(tee -a "$ROOT_DIR/stderr" >&2);

if grep -q 'ERROR' $ROOT_DIR/stdout; then
	echo "Errors found in script logs" >&2;
	cat "$ROOT_DIR/stdout" "$ROOT_DIR/stderr";
	exit 1;
fi

readonly USERS_JSON_FILE="$ROOT_DIR/opt/traffic_ops/install/data/json/users.json";

"$python_bin" <<EOF
import json
import sys

try:
	with open('$USERS_JSON_FILE') as fd:
		users_json = json.load(fd)
except Exception as e:
	print('Error loading users.json file:', e, file=sys.stderr)
	exit(1)

if not isinstance(users_json, dict) or len(users_json) != 2 or 'username' not in users_json or 'password' not in users_json:
	print('Malformed users.json file - not an object or incorrect keys', file=sys.stderr)
	exit(1)

username = users_json['username']
if not isinstance(username, str):
	print('Username is not a string in users.json:', username, file=sys.stderr)
	exit(1)

if username != 'admin':
	print('Incorrect username in users.json, expected: admin, got:', username, file=sys.stderr)
	exit(1)

password = users_json['password']
if not isinstance(password, str):
	print('Password is not a string in users.json:', password, file=sys.stderr)
	exit(1)

if not password.startswith('SCRYPT:16384:8:1:') or len(password.split(':')) != 6:
	print('Malformed password field in users.json:', password, file=sys.stderr)
	exit(1)

exit(0)
EOF

readonly POST_INSTALL_JSON="$ROOT_DIR/opt/traffic_ops/install/data/json/post_install.json";
if [[ "$(cat $POST_INSTALL_JSON)" != "{}" ]]; then
	echo "Incorrect post_install.json, expected: {}, got: $(cat $POST_INSTALL_JSON)" >&2;
	exit 1;
fi

readonly PROFILES_JSON_EXPECTED="{
	\"cdn_name\": \"kabletown_cdn\",
	\"dns_subdomain\": \"cdn1.kabletown.net\",
	\"tm.url\": \"https://localhost\"
}";

readonly PROFILES_JSON_ACTUAL="$(<"$ROOT_DIR/opt/traffic_ops/install/data/json/profiles.json" jq -S --tab .)";
if [[ "$PROFILES_JSON_ACTUAL" != "$PROFILES_JSON_EXPECTED" ]]; then
	echo "Incorrect profiles.json, expected: $PROFILES_JSON_EXPECTED, got: $PROFILES_JSON_ACTUAL" >&2;
	exit 1;
fi

readonly DB_CONF_EXPECTED="production:
    driver: postgres
    open: host=localhost port=5432 user=traffic_ops password=twelve dbname=traffic_ops sslmode=disable";

readonly DB_CONF_ACTUAL="$(cat $ROOT_DIR/opt/traffic_ops/app/db/dbconf.yml)";
if [[ "$DB_CONF_ACTUAL" != "$DB_CONF_EXPECTED" ]]; then
	echo "Incorrect dbconf.yml, expected:" >&2;
	echo "$DB_CONF_EXPECTED" >&2;
	echo "got:" >&2;
	echo "$DB_CONF_ACTUAL" >&2;
	exit 1;
fi

"$python_bin" <<EOF
import json
import string
import sys

try:
	with(open('$ROOT_DIR/opt/traffic_ops/app/conf/cdn.conf')) as fd:
		conf = json.load(fd)
except Exception as e:
	print('Error loading cdn.conf file:', e, file=sys.stderr)
	exit(1)

if not isinstance(conf, dict) or len(conf) != 3 or 'secrets' not in conf or 'to' not in conf or 'traffic_ops_golang' not in conf:
	print('Malformed cdn.conf file - not an object or missing keys', file=sys.stderr)
	exit(1)

if not isinstance(conf['secrets'], list) or len(conf['secrets']) != 1 or not isinstance(conf['secrets'][0], str):
	print('Malformed secrets object in cdn.conf:', conf['secrets'], file=sys.stderr)
	exit(1)

if len(conf['secrets'][0]) != 12 or any(True for x in conf['secrets'][0] if x not in string.ascii_letters + string.digits + '_'):
	print('Incorrect secret in cdn.conf, expected 12 word characters, got:', conf['secrets'][0], file=sys.stderr)
	exit(1)

if not isinstance(conf['to'], dict) or 'base_url' not in conf['to'] or len(conf['to']) != 1 or not isinstance(conf['to']['base_url'], str):
	print('Malformed to object in cdn.conf:', conf['to'])
	exit(1)

if conf['to']['base_url'] != 'http://localhost:3000':
	print('Incorrect to.base_url in cdn.conf, expected: http://localhost:3000, got:', conf['to']['base_url'], file=sys.stderr)
	exit(1)

if not isinstance(conf['traffic_ops_golang'], dict) or len(conf['traffic_ops_golang']) != 5 or 'cert' not in conf['traffic_ops_golang'] or 'key' not in conf['traffic_ops_golang'] or 'port' not in conf['traffic_ops_golang'] or 'log_location_error' not in conf['traffic_ops_golang'] or 'log_location_event' not in conf['traffic_ops_golang']:
	print('Malformed traffic_ops_golang object in cdn.conf:', conf['traffic_ops_golang'], sys.stderr)
	exit(1)

cert='$ROOT_DIR/etc/pki/tls/certs/localhost.crt'
if conf['traffic_ops_golang']['cert']!= cert:
	print('Incorrect cert in cdn.conf, expected:', cert, 'got:', conf['traffic_ops_golang']['cert'], file=sys.stderr)
	exit(1)

key='$ROOT_DIR/etc/pki/tls/private/localhost.key'
if conf['traffic_ops_golang']['key']!= key:
	print('Incorrect key in cdn.conf, expected:', key, 'got:', conf['traffic_ops_golang']['key'], file=sys.stderr)
	exit(1)

if conf['traffic_ops_golang']['port'] != "443":
	print('Incorrect traffic_ops_golang.port, expected: 443, got:', conf['traffic_ops_golang']['port'], file=sys.stderr)
	exit(1)

if conf['traffic_ops_golang']['log_location_error'] != '$ROOT_DIR/var/log/traffic_ops/error.log':
	print('Incorrect traffic_ops_golang.log_location_error in cdn.conf, expected: $ROOT_DIR/var/log/traffic_ops/error.log, got:', conf['traffic_ops_golang']['log_location_error'], file=sys.stderr)
	exit(1)

if conf['traffic_ops_golang']['log_location_event'] != '$ROOT_DIR/var/log/traffic_ops/access.log':
	print('Incorrect traffic_ops_golang.log_location_event in cdn.conf, expected: $ROOT_DIR/var/log/traffic_ops/access.log, got:', conf['traffic_ops_golang']['log_location_event'], file=sys.stderr)
	exit(1)

exit(0)
EOF

readonly DATABASE_CONF_EXPECTED='{
	"dbname": "traffic_ops",
	"description": "Pg database on localhost:5432",
	"hostname": "localhost",
	"password": "twelve",
	"port": "5432",
	"type": "Pg",
	"user": "traffic_ops"
}';

readonly DATABASE_CONF_ACTUAL="$(<"$ROOT_DIR/opt/traffic_ops/app/conf/production/database.conf" jq -S --tab .)";
if [[ "$DATABASE_CONF_ACTUAL" != "$DATABASE_CONF_EXPECTED" ]]; then
	echo "Incorrect database.conf, expected: $DATABASE_CONF_EXPECTED, got $DATABASE_CONF_ACTUAL" >&2;
	exit 1;
fi

readonly CSR_FILE="$ROOT_DIR/etc/pki/tls/certs/localhost.csr";
readonly CSR_FILE_TYPE="$(file $CSR_FILE)";
if [[ "$CSR_FILE_TYPE" != "$CSR_FILE: PEM certificate request" ]]; then
	echo "Incorrect csr file, expected a PEM certificate request, got: $CSR_FILE_TYPE" >&2;
	exit 1;
fi

readonly CERT_FILE="$ROOT_DIR/etc/pki/tls/certs/localhost.crt";
readonly CERT_FILE_TYPE="$(file $CERT_FILE)";
if [[ "$CERT_FILE_TYPE" != "$CERT_FILE: PEM certificate" ]]; then
	echo "Incorrect cert file, expected a PEM certificate, got: $CERT_FILE_TYPE" >&2;
	exit 1;
fi

readonly KEY_FILE="$ROOT_DIR/etc/pki/tls/private/localhost.key";
readonly KEY_FILE_TYPE="$(file $KEY_FILE)";
if [[ "$KEY_FILE_TYPE" != "$KEY_FILE: PEM RSA private key" ]]; then
	echo "Incorrect key file, expected PEM RSA private key, got: $KEY_FILE_TYPE" >&2;
	exit 1;
fi
