module github.com/apache/trafficcontrol/v8

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

go 1.22.0

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/basho/riak-go-client v1.7.1-0.20170327205844-5587c16e0b8b
	github.com/cbroglie/mustache v1.4.0
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/dchest/siphash v1.2.2
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
	github.com/lestrrat-go/jwx v1.2.29
	github.com/lib/pq v1.10.4
	github.com/miekg/dns v1.1.43
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pborman/getopt/v2 v2.1.0
	github.com/pkg/errors v0.9.1
	go.etcd.io/bbolt v1.3.6
	golang.org/x/crypto v0.21.0
	golang.org/x/net v0.21.0
	golang.org/x/sys v0.18.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/basho/backoff v0.0.0-20150307023525-2ff7c4694083 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8
	go.uber.org/atomic v1.6.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

require (
	github.com/Shopify/sarama v1.27.2
	github.com/klauspost/compress v1.13.6 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.5.0 // indirect
	gopkg.in/jcmturner/rpc.v1 v1.1.0 // indirect
)
