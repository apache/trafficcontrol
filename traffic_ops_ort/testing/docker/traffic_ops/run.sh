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
# TO_ADMIN_USER
# TO_ADMIN_PASS
# CERT_COUNTRY
# CERT_STATE
# CERT_CITY
# CERT_COMPANY
# TO_DOMAIN
# TRAFFIC_VAULT_PASS

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS TO_ADMIN_USER TO_ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY TO_DOMAIN)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

start() {
	service traffic_ops start
	exec tail -f /var/log/traffic_ops/traffic_ops.log
}

init() {
	local postinstall_input_file="postinstall-input.json"
	cat > "$postinstall_input_file" <<- ENDOFMESSAGE
{
  "/opt/traffic_ops/app/conf/production/database.conf":[
    {
      "Database type":"Pg",
      "config_var":"type"
    },
    {
      "Database name":"$DB_NAME",
      "config_var":"dbname"
    },
    {
      "Database server hostname IP or FQDN":"$DB_SERVER",
      "config_var":"hostname"
    },
    {
      "Database port number":"$DB_PORT",
      "config_var":"port"
    },
    {
      "Traffic Ops database user":"$DB_USER",
      "config_var":"user"
    },
    {
      "Traffic Ops database password":"$DB_USER_PASS",
      "config_var":"password",
      "hidden":"1"
    }
  ],
  "/opt/traffic_ops/app/db/dbconf.yml":[
    {
      "Database server root (admin) user":"postgres",
      "config_var":"pgUser"
    },
    {
      "Database server admin password":"$DB_ROOT_PASS",
      "config_var":"pgPassword",
      "hidden":"1"
    },
    {
      "Download Maxmind Database?":"yes",
      "config_var":"maxmind"
    }
  ],
  "/opt/traffic_ops/app/conf/cdn.conf":[
    {
      "Generate a new secret?":"yes",
      "config_var":"genSecret"
    },
    {
      "Port to serve on?": "443",
      "config_var": "port"
    },
    {
      "Number of workers?": "12",
      "config_var":"workers"
    },
    {
      "Traffic Ops url?": "https://$TO_HOSTNAME",
      "config_var": "base_url"
    },
    {
      "Number of secrets to keep?":"1",
      "config_var":"keepSecrets"
    }
  ],
  "/opt/traffic_ops/app/conf/ldap.conf":[
    {
      "Do you want to set up LDAP?":"no",
      "config_var":"setupLdap"
    },
    {
      "LDAP server hostname":"",
      "config_var":"host"
    },
    {
      "LDAP Admin DN":"",
      "config_var":"admin_dn"
    },
    {
      "LDAP Admin Password":"",
      "config_var":"admin_pass",
      "hidden":"1"
    },
    {
      "LDAP Search Base":"",
      "config_var":"search_base"
    }
  ],
  "/opt/traffic_ops/install/data/json/users.json":[
    {
      "Administration username for Traffic Ops":"$TO_ADMIN_USER",
      "config_var":"tmAdminUser"
    },
    {
      "Password for the admin user":"$TO_ADMIN_PASS",
      "config_var":"tmAdminPw",
      "hidden":"1"
    }
  ],
  "/opt/traffic_ops/install/data/profiles/":[
    {
      "Add custom profiles?":"no",
      "config_var":"custom_profiles"
    }
  ],
  "/opt/traffic_ops/install/data/json/openssl_configuration.json":[
    {
      "Do you want to generate a certificate?":"yes",
      "config_var":"genCert"
    },
    {
      "Country Name (2 letter code)":"$CERT_COUNTRY",
      "config_var":"country"
    },
    {
      "State or Province Name (full name)":"$CERT_STATE",
      "config_var":"state"
    },
    {
      "Locality Name (eg, city)":"$CERT_CITY",
      "config_var":"locality"
    },
    {
      "Organization Name (eg, company)":"$CERT_COMPANY",
      "config_var":"company"
    },
    {
      "Organizational Unit Name (eg, section)":"",
      "config_var":"org_unit"
    },
    {
      "Common Name (eg, your name or your server's hostname)":"$TO_HOSTNAME",
      "config_var":"common_name"
    },
    {
      "RSA Passphrase":"passphrase",
      "config_var":"rsaPassword",
      "hidden":"1"
    }
  ],
  "/opt/traffic_ops/install/data/json/profiles.json":[
    {
      "Traffic Ops url":"https://$TO_HOSTNAME",
      "config_var":"tm.url"
    },
    {
      "Human-readable CDN Name.  (No whitespace, please)":"cdn",
      "config_var":"cdn_name"
    },
    {
      "Health Polling Interval (milliseconds)":"8000",
      "config_var":"health_polling_int"
    },
    {
      "DNS sub-domain for which your CDN is authoritative":"$TO_HOSTNAME.$TO_DOMAIN",
      "config_var":"dns_subdomain"
    },
    {
      "TLD SOA admin":"traffic_ops",
      "config_var":"soa_admin"
    },
    {
      "TrafficServer Drive Prefix":"/dev/ram",
      "config_var":"driver_prefix"
    },
    {
      "TrafficServer RAM Drive Prefix":"/dev/ram",
      "config_var":"ram_drive_prefix"
    },
    {
      "TrafficServer RAM Drive Letters (comma separated)":"1",
      "config_var":"ram_drive_letters"
    },
    {
      "Health Threshold Load Average":"25",
      "config_var":"health_thresh_load_avg"
    },
    {
      "Health Threshold Available Bandwidth in Kbps":"1750000",
      "config_var":"health_thresh_kbps"
    },
    {
      "Traffic Server Health Connection Timeout (milliseconds)":"2000",
      "config_var":"health_connect_timeout"
    }
  ]
}
ENDOFMESSAGE

cat > /opt/traffic_ops/app/conf/production/riak.conf << EOM
{
  "user": "riakuser",
  "password": "$RIAK_USER_PASS",
  "MaxTLSVersion": "1.1",
  "tlsConfig": {
    "insecureSkipVerify": true
  }
}
EOM

	# TODO determine if term, user are necessary
	export TERM=xterm && export USER=root && /opt/traffic_ops/install/bin/postinstall -cfile "$postinstall_input_file"

	# Only listen on IPv4, not IPv6, because Docker doesn't provide a v6 interface by default. See http://mojolicious.org/perldoc/Mojo/Server/Daemon#listen
	sed -i -e 's#https://\[::\]#https://127\.0\.0\.1#' /opt/traffic_ops/app/conf/cdn.conf
	service traffic_ops restart

}

if [ -f /GO_VERSION ]; then
  go_version=$(cat /GO_VERSION) && \
      curl -Lo go.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
        tar -C /usr/local -xvzf go.tar.gz && \
        ln -s /usr/local/go/bin/go /usr/bin/go && \
        rm go.tar.gz
else
  echo "no GO_VERSION file, unable to install go"
  exit 0
fi
/opt/traffic_ops/install/bin/install_goose.sh

(cd /opt/traffic_ops/app && db/admin --env=production reset)
source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
