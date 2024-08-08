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

> **note**
>
> For a more in-depth discussion of the CDN in a Box system, please see [the official documentation](https://traffic-control-cdn.readthedocs.io/en/latest/admin/quick_howto/ciab.html).

## Setup
The containers run on Docker, and require Docker (tested v17.05.0-ce) and Docker
Compose (tested v1.9.0) to build and run. On most 'nix systems these can be installed
via the distribution's package manager under the names `docker-ce` and
`docker compose`, respectively (e.g. `sudo dnf install docker-ce`).

Each container (except the origin) requires an `.rpm` file to install the Traffic Control
component for which it is responsible. You can download these `*.rpm` files from an archive
(e.g. under "Releases"), use the provided [Makefile](./Makefile) to generate them (simply
type `make` while in the `cdn-in-a-box` directory) or create them yourself by using the
[`pkg`](../../pkg) script at the root of the repository. If you choose the latter, copy
the `*.rpm`s without any version/architecture information to their respective component
directories, such that their filenames are as follows:

* `edge/trafficcontrol-cache-config.rpm`
* `mid/trafficcontrol-cache-config.rpm`
* `traffic_monitor/traffic_monitor.rpm`
* `traffic_ops/traffic_ops.rpm`
* `traffic_portal/traffic_portal.rpm`

Finally, run the test CDN using the command:

```bash
docker compose up --build
```

## Readiness
To know if your CDN in a Box has started up successfully and is ready to use,
you can optionally start the "readiness" container which will test your CDN and
exit successfully when your CDN in a Box is ready:

```bash
docker compose -f docker-compose.readiness.yml up --build
```

If the container does not exit successfully after a reasonable amount of time,
something might have gone wrong with the main CDN services. Because the
container continually runs end-to-end CDN requests, it will never exit
successfully if there are issues with the main CDN services that cause the
requests to fail. Check the log output of the main CDN services to see what
might be getting stuck.

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

> <table>
> <colgroup>
> <col width="18%" />
> <col width="34%" />
> <col width="22%" />
> <col width="24%" />
> </colgroup>
> <thead>
> <tr class="header">
> <th align="left">Service</th>
> <th align="left">Ports exposed and their usage</th>
> <th align="left">Username</th>
> <th align="left">Password</th>
> </tr>
> </thead>
> <tbody>
> <tr class="odd">
> <td align="left">DNS</td>
> <td align="left">DNS name resolution on 9353</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="even">
> <td align="left">Edge-Tier Cache</td>
> <td align="left">Apache Trafficserver HTTP caching reverse proxy on port 9000</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="odd">
> <td align="left">Mid-Tier Cache</td>
> <td align="left">Apache Trafficserver HTTP caching forward proxy on port 9100</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="even">
> <td align="left">Second Mid-Tier Cache (parent of the first Mid-Tier Cache)</td>
> <td align="left">Apache Trafficserver HTTP caching forward proxy on port 9100</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="odd">
> <td align="left">Mock Origin Server</td>
> <td align="left">Example web page served on port 9200</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="even">
> <td align="left">Traffic Monitor</td>
> <td align="left">Web interface and API on port 80</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="odd">
> <td align="left">Traffic Ops</td>
> <td align="left">API on port 6443</td>
> <td align="left"><code>TO_ADMIN_USER</code> in variables.env</td>
> <td align="left"><code>TO_ADMIN_PASSWORD</code> in variables.env</td>
> </tr>
> <tr class="even">
> <td align="left">Traffic Ops PostgresQL Database</td>
> <td align="left">PostgresQL connections accepted on port 5432 (database name: <code>DB_NAME</code> in variables.env)</td>
> <td align="left"><code>DB_USER</code> in variables.env</td>
> <td align="left"><code>DB_USER_PASS</code> in variables.env</td>
> </tr>
> <tr class="odd">
> <td align="left">Traffic Portal</td>
> <td align="left">Web interface on 443 (Javascript required)</td>
> <td align="left"><code>TO_ADMIN_USER</code> in variables.env</td>
> <td align="left"><code>TO_ADMIN_PASSWORD</code> in variables.env</td>
> </tr>
> <tr class="even">
> <td align="left">Traffic Router</td>
> <td align="left">Web interfaces on ports 3080 (HTTP) and 3443 (HTTPS), with a DNS service on 53 and an API on 3333</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tbody>
> </table>
>

## Host Ports

By default, `docker-compose.yml` does not expose ports to the host. This allows the host to be running other services on those ports, as well as allowing multiple CDN-in-a-Boxes to run on the same host, without port conflicts.

To expose the ports of each service on the host, add the `docker-compose.expose-ports.yml` file. For example, `docker compose -f docker-compose.yml -f docker-compose.expose-ports.yml up`.

## Varnish

By default, CDN-in-a-Box uses Apache Traffic Server as the cache server.

To run CDN-in-a-Box with Varnish add the `docker-compose.varnish.yml` file.
For example, `docker compose -f docker-compose.yml -f docker-compose.varnish.yml up`

## Common Pitfalls

### Traffic Monitor is stuck waiting for a valid Snapshot

Often times you must take a CDN [Snapshot](https://traffic-control-cdn.readthedocs.io/en/latest/glossary.html#term-snapshot) in order for a valid Snapshot to be generated. This can be done through [Traffic Portal's "CDNs" view](https://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_portal/usingtrafficportal.html#cdns), clicking on the "CDN-in-a-Box" CDN, then pressing the camera button, and finally the "Perform Snapshot" button.

### I'm seeing a failure to open a socket and/or set a socket option

Try disabling SELinux or setting it to 'permissive'. SELinux hates letting containers bind to certain ports. You can also try re-labeling the `docker` executable if you feel comfortable.

### Traffic Vault container exits with cp /usr/local/share/ca-certificates cp: missing destination

Bring all components down, remove the `traffic_ops/ca` directory, and delete the volumes with `docker volume prune`. This will force the regeneration of the certificates.

## Notes for macOS

CDN in a Box should work, without modification, on any architecture that can run Docker. If it does not, that is a bug, please open a new issue for it.

Re/installed docker from the command line to use `--user=[your username]` flag. Link to install info for docker for [Install from command line] (https://docs.docker.com/desktop/install/mac-install/)  and use Install from command line located in the "Mac with Apple Silicon" tab.

Build and run of it:
In the trafficcontrol/infrastructure/cdn-in-a-box directory run the following:

- `make build-builders` ~~> this will create all the rpms and copy each rpms into its own folder in the cdn-in-a-box project. This will also create the dist folder under trafficcontrol folder structure even if you deleted yours.

- `docker compose up` ~~> this will create docker images, and if rebuild is needed, run `docker compose up --build`.

### Docker v4.11 and later of Docker Desktop for Mac

"Privileged configurations are applied during the installation with the --user flag on the install command. In this case, the user is not prompted to grant root privileges on the first run of Docker Desktop. Specifically, the --user flag:

- Uninstalls the previous com.docker.vmnetd if present
- Sets up symlinks for the user
- Ensures that localhost is resolved to 127.0.0.1

The limitation of this approach is that Docker Desktop can only be run by one user account per machine, namely the one specified in the -–user flag."

The above is a direct quote found in [Installing from the commandline] (https://docs.docker.com/desktop/mac/permission-requirements/)

Note: The Install from command line was the only install that resolved my issue of being unable to build the docker images and run the containers in localhost. Homebrew and the usual automatic install when double-clicking in the Docker.dmg did not work.
