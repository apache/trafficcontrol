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

# CDN In a Box (containerized)
This is intended to simplify the process of creating a "CDN in a box", easing
the barrier to entry for newcomers as well as providing a way to spin up a
minimal CDN for full system testing.

## Implemented Components
As of the time of this writing, Traffic Ops, Traffic Monitor, Traffic Portal and a
database server for Traffic Ops are all fully implemented. An edge-tier cache,
mid-tier cache and simple origin are also all implemented, although with limited
functionality (caches do not respond to update queues/routing, and origin serves
static HTTP content - likely subject to change in favor of proper video streaming in
the future). Other components will follow as well as details on specific parts of the
implementation.

## Setup
The containers run on docker, and require Docker (tested v17.05.0-ce) and Docker
Compose (tested v1.9.0) to build and run. On most 'nix systems these can be installed
via the distribution's package manager under the names `docker-ce` and
`docker-compose`, respectively (e.g. `sudo yum install docker-ce`).

Each container (except the origin) requires an `.rpm` file to install the Traffic Control
component for which it is responsible. You can either download these `*.rpm` files or
create them yourself by using the [`pkg`](../../pkg) script at the root of the
repository. Copy the `*.rpm`s without any version/architecture information to their
respective component directories, such that their filenames are as follows:

* `edge/traffic_ops_ort.rpm`
* `mid/traffic_ops_ort.rpm`
* `traffic_monitor/traffic_monitor.rpm`
* `traffic_ops/traffic_ops.rpm`
* `traffic_portal/traffic_portal.rpm`

Also ensure that the edge- and mid-tier caches have copies of the
[`traffic_ops/to-access.sh`](./traffic_ops/to-access) script in their own directories.
From this directory, you can accomplish this by running:

```bash
cp -f traffic_ops/to-access.sh edge/to-access.sh && cp -f traffic_ops/to-access.sh mid/to-access.sh
```

Finally, run the test CDN using the command:

```bash
docker-compose up --build
```

## Components
> The following assumes that the default configuration provided in
> [`variables.env`](./variables.env) is used.

Once your CDN is running, you should see a cascade of output on your terminal. This is
typically the output of the build, then setup, and finally logging infrastructure
(assuming nothing goes wrong). You can now access the various components of the CDN on
your local machine. For example, opening [`https://localhost`](https://localhost) should
show you the default UI for interacting with the CDN - Traffic Portal.

> Note: You will likely see a warning about an untrusted or invalid certificate for
> components that serve over HTTPS (Traffic Ops & Traffic Portal). If you
> are sure that you are looking at `https://localhost:N` for some integer `N`, these
> warnings may be safely ignored via the e.g. `Add Exception` button (possibly hidden
> behind e.g. `Advanced Options`).

### Traffic Ops
The API and legacy UI for the CDN
* URLs:
	* New Golang endpoints: [`https://localhost:6443`](https://localhost:6443)
	* Limited, Legacy Perl endpoints: [`https://localhost:60443`](https://localhost:60443)
* Login Credentials:
	* Username: `admin`
	* Password: `twelve`

### Traffic Portal
The modern UI for the CDN
* URL: [`https://localhost`](https://localhost)
* Login Credentials:
	* Username: `admin`
	* Password: `twelve`

### Traffic Ops PostgreSQL Database
This holds the configuration information for the entire CDN. It is normally accessed
directly only by Traffic Ops.
* URL: [`postgres://traffic_ops:twelve@localhost:5432/traffic_ops`](postgres://traffic_ops:twelve@localhost:5432/traffic_ops)
* Login Credentials:
	* Username: `traffic_ops`
	* Password: `twelve`
* Port: 5432
* Database: `traffic_ops`

### Traffic Vault
A secure storage server for private keys used by Traffic Ops
* Port: 8010

### Edge-Tier Cache
An edge-tier cache sits at the outermost extremity of the CDN, typically serving content
directly to the user from either its cache, or a parent cache group. The management port
is not exposed locally - however the main content port is.

* URL: [`http://localhost:8080`](http://localhost:8080)

### Mid-Tier Cache
A mid-tier cache serves content internally - within the CDN. Typically, requests are
made from an edge-tier cache and the mid serves content from its own cache, a parent
mid-tier cache or directly from the origin. The management port is not exposed locally -
however the main content port is.

* URL: [`http://localhost:9080`](http://localhost:9080)

### Origin Server
An origin server simply serves HTTP(S) content. The CDN-in-a-box origin server serves up
a very simple page sporting the Traffic Control logo.

* URL: [`http://localhost`](http://localhost)

* URL: [`http://localhost:8081`](http://localhost:8081)

The process creates containers for each component with ports exposed on the host.  The
following should be available once the system is running:

	Traffic Portal: https://localhost
	Traffic Ops (go): https://localhost:6443
	Traffic Ops (perl): https://localhost:60443
	Postgres: `psql -h localhost -p 5432 -U postgres`
