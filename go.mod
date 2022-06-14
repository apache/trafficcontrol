module github.com/apache/trafficcontrol

// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

go 1.17

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/basho/riak-go-client v1.7.1-0.20170327205844-5587c16e0b8b
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/dchest/siphash v1.2.2
	github.com/dgrijalva/jwt-go v3.2.1-0.20190620180102-5e25c22bd5d6+incompatible
	github.com/fsnotify/fsnotify v1.5.1
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/gofrs/flock v0.8.1
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/google/uuid v1.3.0
	github.com/hydrogen18/stoppableListener v0.0.0-20161101122645-827d760f0663
	github.com/influxdata/influxdb v1.9.5
	github.com/jmoiron/sqlx v1.3.4
	github.com/json-iterator/go v1.1.12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kylelemons/godebug v1.1.1-0.20201107061927-e693023230a4
	github.com/lestrrat-go/jwx v1.2.12
	github.com/lib/pq v1.10.4
	github.com/miekg/dns v1.1.43
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pborman/getopt/v2 v2.1.0
	github.com/pkg/errors v0.9.1
	go.etcd.io/bbolt v1.3.6
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/net v0.0.0-20220105145211-5b0dc2dfae98
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/basho/backoff v0.0.0-20150307023525-2ff7c4694083 // indirect
	github.com/google/gopacket v1.1.19
)
