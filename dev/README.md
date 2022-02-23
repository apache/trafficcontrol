<!--
	Licensed to the Apache Software Foundation (ASF) under one
	or more contributor license agreements.  See the NOTICE file
	distributed with this work for additional information
	regarding copyright ownership.  The ASF licenses this file
	to you under the Apache License, Version 2.0 (the
	"License"); you may not use this file except in compliance
	with the License.  You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing,
	software distributed under the License is distributed on an
	"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
	KIND, either express or implied.  See the License for the
	specific language governing permissions and limitations
	under the License.
-->

# Development Environment
The ATC development environment - housed in this directory - can be used to
quickly make changes to ATC components and test them immediately.

To use the development environment, ensure you are in the repository's root and
source the `dev/atc.dev.sh` file. Then, use `atc` to run commands (see
`atc --help` for usage).

## Traffic Ops
Traffic Ops will start with its API bound to local port 6443. The API will use a
self-signed certificate, so `curl` commands to the API will need to use
`-k`/`--insecure` (as will `toget`/`toput`/`topost`/`todelete`). The Delve
debugger for Go listens on port 6444 for connections to debug Traffic Ops.

The login credentials for the "admin" user are the same as those for the user of
the same name in CDN-in-a-Box - password is `twelve12`.

## Traffic Portal
The current version of Traffic Portal serves using HTTPS on port 444. The
certificate it uses is self-signed, so browsers will warn that the site is
insecure.

## Traffic Portal "v2"
The experimental Traffic Portal (`experimental/traffic-portal`) serves using
HTTP on port 443. The certificate it uses is self-signed, so browsers will warn
that the site is insecure.

## Traffic Monitor
The Traffic Monitor API is served locally over HTTP on port 80. The Delve
debugger for Go listens on port 81 for connections to debug Traffic Monitor.

Note that Traffic Monitor will do almost nothing useful if the edge cache server
(the `t3c` service) is not running when it starts.

Traffic Monitor writes its backups for CDN Snapshots and Monitoring Configs in
the `dev/traffic_monitor` directory, so you can see them.

## Database/Traffic Vault
A Postgres database listens on port 5432 (this conflicts with the default port
for running Postgres, so any Postgres servers running on the host machine may
need to be stopped before running ATC) and houses the Traffic Ops database as
`traffic_ops_development`, and the Traffic Vault database as
`traffic_vault_development`. To connect as the Traffic Ops user to the Traffic
Ops database, use the username `traffic_ops` and the password `twelve12`. To
connect as the Traffic Ops Vault user to the Traffic Vault database, use the
username `traffic_vault` and the password `twelve12`.

## T3C
An edge-tier cache server listens for HTTP (HTTPS not supported) connections on
local port 8080. The Delve debugger for Go listens on port 8081 for connections
to debug `t3c` sub-commands.

Note that, while in most production deployments `t3c` runs on a `cron` schedule,
`t3c` is never run in this service container, normally. One must manually trigger
a run, usually by using the `t3c` function provided by `atc.dev.sh`.

## Traffic Router
Traffic Router listens locally for DNS queries on port 3053 (TCP and/or UDP),
HTTP requests from clients to be routed on ports 3080 (HTTP) and 3443 (HTTPS),
HTTP requests to its API on ports 3333 (HTTP) and 2222 (HTTPS), and listens for
JDPA debugging connections on port 5005.

Traffic Router writes its backups for the "Coverage Zone" file, CDN Snapshot,
Federations, cache health (as published by Traffic Monitor), LetsEncrypt data,
and Steering information into `dev/traffic_router/db/` so you can see them.
Generated DNS Zones are written in `dev/traffic_router/var/` and Traffic Router
will use `dev/traffic_router/temp` to create any temporary files it needs (that
will grow without bound, so it may need to be cleaned up every now and again).
