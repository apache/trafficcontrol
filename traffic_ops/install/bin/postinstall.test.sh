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

if [[ ! -x /usr/bin/python3 ]]; then
	echo "Python 3.6+ is required to run - or test - postinstall.py" >&2;
	exit 1;
fi

readonly ROOT_DIR="$(mktemp -d)";

trap 'rm -rf $ROOT_DIR' EXIT;

mkdir -p "$ROOT_DIR/etc/pki/tls/certs";
mkdir "$ROOT_DIR/etc/pki/tls/private";
mkdir -p "$ROOT_DIR/opt/traffic_ops/app/public/routing";
mkdir "$ROOT_DIR/opt/traffic_ops/app/db";
mkdir -p "$ROOT_DIR/opt/traffic_ops/app/conf/production";
cat > "$ROOT_DIR/opt/traffic_ops/app/conf/cdn.conf" <<EOF
{
	"hypnotoad": {
		"listen": [
			"https://[::]:60443?cert=$ROOT_DIR/etc/pki/tls/certs/localhost.crt&key=$ROOT_DIR/etc/pki/tls/private/localhost.key"
		]
	}
}
EOF

mkdir -p "$ROOT_DIR/opt/traffic_ops/install/data/json";
mkdir "$ROOT_DIR/opt/traffic_ops/install/bin";

readonly MY_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )";

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
			"Password for Traffic Ops database user": "twelve",
			"config_var": "password",
			"hidden": true
		}
	],
	"/opt/traffic_ops/app/db/dbconf.yml": [
		{
			"Database server root (admin) user": "postgres",
			"config_var": "pgUser",
			"hidden": false
		},
		{
			"Password for database server admin": "twelve",
			"config_var": "pgPassword",
			"hidden": true
		},
		{
			"Download Maxmind Database?": "no",
			"config_var": "maxmind",
			"hidden": false
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
			"Do you want to set up LDAP?": "no",
			"config_var": "setupLdap",
			"hidden": false
		},
		{
			"LDAP server hostname": "",
			"config_var": "host",
			"hidden": false
		},
		{
			"LDAP Admin DN": "",
			"config_var": "admin_dn",
			"hidden": false
		},
		{
			"LDAP Admin Password": "",
			"config_var": "admin_pass",
			"hidden": true
		},
		{
			"LDAP Search Base": "",
			"config_var": "search_base",
			"hidden": false
		},
		{
			"LDAP Search Query": "",
			"config_var": "search_query",
			"hidden": false
		},
		{
			"LDAP Skip TLS verify": "",
			"config_var": "insecure",
			"hidden": false
		},
		{
			"LDAP Timeout Seconds": "",
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

"$MY_DIR/postinstall.py" --no-root --root-directory="$ROOT_DIR" --no-restart-to --no-database --ops-user="$(whoami)" --ops-group="$(id -gn)" --automatic --cfile="$ROOT_DIR/defaults.json" --debug 2>"$ROOT_DIR/stderr" | tee "$ROOT_DIR/stdout"

if grep -q 'ERROR' $ROOT_DIR/stderr; then
	echo "Errors found in script logs" >&2;
	cat "$ROOT_DIR/stderr";
	cat "$ROOT_DIR/stdout";
	exit 1;
fi

readonly USERS_JSON_FILE="$ROOT_DIR/opt/traffic_ops/install/data/json/users.json";

/usr/bin/python3 <<EOF
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
	\"tm.url\": \"https://localhost\",
	\"cdn_name\": \"kabletown_cdn\",
	\"dns_subdomain\": \"cdn1.kabletown.net\"
}";

readonly PROFILES_JSON_ACTUAL="$(cat $ROOT_DIR/opt/traffic_ops/install/data/json/profiles.json)";
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

/usr/bin/python3 <<EOF
import json
import string
import sys

try:
	with(open('$ROOT_DIR/opt/traffic_ops/app/conf/cdn.conf')) as fd:
		conf = json.load(fd)
except Exception as e:
	print('Error loading cdn.conf file:', e, file=sys.stderr)
	exit(1)

if not isinstance(conf, dict) or len(conf) != 4 or 'hypnotoad' not in conf or 'secrets' not in conf or 'to' not in conf or 'traffic_ops_golang' not in conf:
	print('Malformed cdn.conf file - not an object or missing keys', file=sys.stderr)
	exit(1)

if not isinstance(conf['hypnotoad'], dict) or len(conf['hypnotoad']) != 1 or 'listen' not in conf['hypnotoad'] or not isinstance(conf['hypnotoad']['listen'], list) or len(conf['hypnotoad']['listen']) != 1 or not isinstance(conf['hypnotoad']['listen'][0], str):
	print('Malformed hypnotoad object in cdn.conf:', conf['hypnotoad'], file=sys.stderr)
	exit(1)

listen = 'https://[::]:60443?cert=$ROOT_DIR/etc/pki/tls/certs/localhost.crt&key=$ROOT_DIR/etc/pki/tls/private/localhost.key'
if conf['hypnotoad']['listen'][0] != listen:
	print('Incorrect hypnotoad.listen[0] in cdn.conf, expected:', listen, 'got:', conf['hypnotoad']['listen'][0], file=sys.stderr)
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

if not isinstance(conf['traffic_ops_golang'], dict) or len(conf['traffic_ops_golang']) != 3 or 'port' not in conf['traffic_ops_golang'] or 'log_location_error' not in conf['traffic_ops_golang'] or 'log_location_event' not in conf['traffic_ops_golang']:
	print('Malformed traffic_ops_golang object in cdn.conf:', conf['traffic_ops_golang'], sys.stderr)
	exit(1)

if conf['traffic_ops_golang']['port'] != '443':
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
	"type": "Pg",
	"dbname": "traffic_ops",
	"hostname": "localhost",
	"port": "5432",
	"user": "traffic_ops",
	"password": "twelve",
	"description": "Pg database on localhost:5432"
}';

readonly DATABASE_CONF_ACTUAL="$(cat $ROOT_DIR/opt/traffic_ops/app/conf/production/database.conf)";
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
