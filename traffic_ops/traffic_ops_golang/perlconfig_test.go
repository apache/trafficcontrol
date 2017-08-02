package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

func TestGetCDNConf(t *testing.T) {
	input := `
{
	hypnotoad => {
		listen => [
			'https://[::]:60443?cert=/etc/pki/tls/certs/localhost.crt&key=/etc/pki/tls/private/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED'
		],
		user     => 'trafops',
		group    => 'trafops',
		heartbeat_timeout => 20,
		pid_file => '/var/run/traffic_ops.pid',
		workers  => 96
	},
	cors => {
		access_control_allow_origin => '*'
	},
	to => {
		base_url   => 'http://localhost:3000',                    # this is where traffic ops app resides
		email_from => 'no-reply@traffic-ops-domain.com'           # traffic ops email address
	},
	portal => {
		base_url   => 'http://localhost:8080',                    # this is where the traffic portal resides (a javascript client that consumes the TO API)
		email_from => 'no-reply@traffic-portal-domain.com'        # traffic portal email address
	},

	# 1st secret is used to generate new signatures. Older one kept around for existing signed cookies.
		#  Remove old one(s) when ready to invalidate old cookies.
		secrets => [ 'walrus' ],
	geniso  => {
		iso_root_path => '/opt/traffic_ops/app/public',          # the location where the iso files will be written
	},
	inactivity_timeout => 60,
	traffic_ops_golang_port => '443'
};
`

	expected := Config{
		HTTPPort: "443",
		TOSecret: "walrus",
		TOURLStr: "https://127.0.0.1:60443",
		CertPath: "/etc/pki/tls/certs/localhost.crt",
		KeyPath:  "/etc/pki/tls/private/localhost.key",
	}
	err := error(nil)
	if expected.TOURL, err = url.Parse(expected.TOURLStr); err != nil {
		t.Errorf("expected URL parse '%+v' err nil actual %+v", expected.TOURLStr, err)
	}

	cfg, err := getCDNConf(input)
	if err != nil {
		t.Errorf("expected nil err actual %v", err)
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("expected %+v actual %+v", expected, cfg)
	}
}

func TestGetPerlConfigsFromStrs(t *testing.T) {
	cdnConfInput := `
{
	hypnotoad => {
		listen => [
			'https://[::]:60443?cert=/etc/pki/tls/certs/localhost.crt&key=/etc/pki/tls/private/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED'
		],
		user     => 'trafops',
		group    => 'trafops',
		heartbeat_timeout => 20,
		pid_file => '/var/run/traffic_ops.pid',
		workers  => 96
	},
	cors => {
		access_control_allow_origin => '*'
	},
	to => {
		base_url   => 'http://localhost:3000',                    # this is where traffic ops app resides
		email_from => 'no-reply@traffic-ops-domain.com'           # traffic ops email address
	},
	portal => {
		base_url   => 'http://localhost:8080',                    # this is where the traffic portal resides (a javascript client that consumes the TO API)
		email_from => 'no-reply@traffic-portal-domain.com'        # traffic portal email address
	},

	# 1st secret is used to generate new signatures. Older one kept around for existing signed cookies.
		#  Remove old one(s) when ready to invalidate old cookies.
		secrets => [ 'walrus' ],
	geniso  => {
		iso_root_path => '/opt/traffic_ops/app/public',          # the location where the iso files will be written
	},
	inactivity_timeout => 60,
	traffic_ops_golang_port => '443'
};
`

	dbConfInput := `
{
   "password" : "thelizard",
   "user" : "bill",
   "type" : "Pg",
   "hostname" : "db.to.example.net",
   "description" : "Postgres database",
   "port" : "5432",
   "dbname" : "to"
}
`

	expected := Config{
		HTTPPort:           "443",
		DBUser:             "bill",
		DBPass:             "thelizard",
		DBServer:           "db.to.example.net:5432",
		DBDB:               "to",
		DBSSL:              false,
		TOSecret:           "walrus",
		TOURLStr:           "https://127.0.0.1:60443",
		CertPath:           "/etc/pki/tls/certs/localhost.crt",
		KeyPath:            "/etc/pki/tls/private/localhost.key",
		MaxDBConnections:   DefaultMaxDBConnections,
		LogLocationError:   NewLogPath,
		LogLocationWarning: NewLogPath,
		LogLocationInfo:    NewLogPath,
		LogLocationEvent:   OldAccessLogPath,
		LogLocationDebug:   log.LogLocationNull,
	}
	err := error(nil)
	if expected.TOURL, err = url.Parse(expected.TOURLStr); err != nil {
		t.Errorf("expected URL parse '%+v' err nil actual %+v", expected.TOURLStr, err)
	}

	cfg, err := getPerlConfigsFromStrs(cdnConfInput, dbConfInput)
	if err != nil {
		t.Errorf("expected nil err actual %v", err)
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("expected %+v actual %+v", expected, cfg)
	}
}
