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

go 1.15

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/basho/backoff v0.0.0-20150307023525-2ff7c4694083 // indirect
	github.com/basho/riak-go-client v1.7.0
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/dchest/siphash v1.2.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/gofrs/flock v0.8.0
	github.com/google/uuid v1.1.2
	github.com/hydrogen18/stoppableListener v0.0.0-20161101122645-827d760f0663
	github.com/influxdata/influxdb v1.8.3
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.10
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kylelemons/godebug v1.1.1-0.20201107061927-e693023230a4
	github.com/lestrrat-go/jwx v1.0.5
	github.com/lestrrat/go-jwx v0.0.0-20180221005942-b7d4802280ae
	github.com/lestrrat/go-pdebug v0.0.0-20180220043741-569c97477ae8 // indirect
	github.com/lib/pq v1.8.0
	github.com/miekg/dns v1.1.35
	github.com/ogier/pflag v0.0.1
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/pborman/getopt/v2 v2.1.0
	github.com/pkg/errors v0.9.1
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sys v0.0.0-20201117222635-ba5294a509c7
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)
