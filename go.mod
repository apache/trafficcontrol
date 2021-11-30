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

replace (
	github.com/fsnotify/fsnotify v1.4.9 => github.com/fsnotify/fsnotify v1.3.0
	github.com/golang/protobuf v1.4.2 => github.com/golang/protobuf v0.0.0-20171021043952-1643683e1b54
	gopkg.in/yaml.v2 v2.3.0 => gopkg.in/yaml.v2 v2.2.8
)

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20180108190415-b31f603f5e1e
	github.com/GehirnInc/crypt v0.0.0-20190301055215-6c0105aabd46
	github.com/asaskevich/govalidator v0.0.0-20180319081651-7d2e70ef918f
	github.com/basho/backoff v0.0.0-20150307023525-2ff7c4694083 // indirect
	github.com/basho/riak-go-client v1.7.1-0.20170327205844-5587c16e0b8b
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cihub/seelog v0.0.0-20170110094445-7bfb7937d106
	github.com/dchest/siphash v1.1.0
	github.com/dgrijalva/jwt-go v3.2.1-0.20190620180102-5e25c22bd5d6+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-ozzo/ozzo-validation v3.0.3-0.20180119232150-44af65fe9adf+incompatible
	github.com/gofrs/flock v0.7.2-0.20190320160742-5135e617513b
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/uuid v1.1.2
	github.com/hydrogen18/stoppableListener v0.0.0-20151210151943-dadc9ccc400c
	github.com/influxdata/influxdb v1.1.1-0.20170104212736-6a94d200c826
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.6-0.20181024152841-05d041de1043
	github.com/kelseyhightower/envconfig v1.3.1-0.20180308190516-b2c5c876e265
	github.com/kylelemons/godebug v1.1.1-0.20201107061927-e693023230a4
	github.com/lib/pq v1.8.0
	github.com/mattn/go-sqlite3 v1.14.5 // indirect
	github.com/miekg/dns v1.0.6-0.20180406150955-01d59357d468
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/pborman/getopt/v2 v2.1.0
	github.com/pkg/errors v0.9.1
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211104170005-ce137452f963
	golang.org/x/sys v0.0.0-20211103235746-7861aae1554b
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

require github.com/go-ldap/ldap/v3 v3.4.1

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.0-20210816181553-5444fa50b93d // indirect
	github.com/goccy/go-json v0.7.10 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.0 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/lestrrat-go/jwx v1.2.12-0.20211125230033-881a53637884
	github.com/nxadm/tail v1.4.4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
