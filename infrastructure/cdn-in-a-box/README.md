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
`docker-compose`, respectively (e.g. `sudo yum install docker-ce`).

Each container (except the origin) requires an `.rpm` file to install the Traffic Control
component for which it is responsible. You can download these `*.rpm` files from an archive
(e.g. under "Releases"), use the provided [Makefile](./Makefile) to generate them (simply
type `make` while in the `cdn-in-a-box` directory) or create them yourself by using the
[`pkg`](../../pkg) script at the root of the repository. If you choose the latter, copy
the `*.rpm`s without any version/architecture information to their respective component
directories, such that their filenames are as follows:

* `edge/traffic_ops_ort.rpm`
* `mid/traffic_ops_ort.rpm`
* `traffic_monitor/traffic_monitor.rpm`
* `traffic_ops/traffic_ops.rpm`
* `traffic_portal/traffic_portal.rpm`

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
> <td align="left">Edge Tier Cache</td>
> <td align="left">Apache Trafficserver HTTP caching reverse proxy on port 9000</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="odd">
> <td align="left">Mid Tier Cache</td>
> <td align="left">Apache Trafficserver HTTP caching forward proxy on port 9100</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="even">
> <td align="left">Mock Origin Server</td>
> <td align="left">Example web page served on port 9200</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="odd">
> <td align="left">Traffic Monitor</td>
> <td align="left">Web interface and API on port 80</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="even">
> <td align="left">Traffic Ops</td>
> <td align="left">Main API endpoints on port 6443, with a direct route to the Perl API on port 60443<a href="#fn1" class="footnoteRef" id="fnref1"><sup>1</sup></a></td>
> <td align="left"><code>TO_ADMIN_USER</code> in variables.env</td>
> <td align="left"><code>TO_ADMIN_PASSWORD</code> in variables.env</td>
> </tr>
> <tr class="odd">
> <td align="left">Traffic Ops PostgresQL Database</td>
> <td align="left">PostgresQL connections accepted on port 5432 (database name: <code>DB_NAME</code> in variables.env)</td>
> <td align="left"><code>DB_USER</code> in variables.env</td>
> <td align="left"><code>DB_USER_PASS</code> in variables.env</td>
> </tr>
> <tr class="even">
> <td align="left">Traffic Portal</td>
> <td align="left">Web interface on 443 (Javascript required)</td>
> <td align="left"><code>TO_ADMIN_USER</code> in variables.env</td>
> <td align="left"><code>TO_ADMIN_PASSWORD</code> in variables.env</td>
> </tr>
> <tr class="odd">
> <td align="left">Traffic Router</td>
> <td align="left">Web interfaces on ports 3080 (HTTP) and 3443 (HTTPS), with a DNS service on 53 and an API on 3333</td>
> <td align="left">N/A</td>
> <td align="left">N/A</td>
> </tr>
> <tr class="even">
> <td align="left">Traffic Vault</td>
> <td align="left">Riak key-value store on port 8010</td>
> <td align="left"><code>TV_ADMIN_USER</code> in variables.env</td>
> <td align="left"><code>TV_ADMIN_PASSWORD</code> in variables.env</td>
> </tr>
> </tbody>
> </table>
> <div class="footnotes">
> <hr />
> <ol>
> <li id="fn1"><p>Please do NOT use the Perl endpoints directly. The CDN will only work properly if everything hits the Go API, which will proxy to the Perl endpoints as needed.<a href="#fnref1">↩</a></p></li>
> </ol>
> </div>
>

## Host Ports

By default, `docker-compose.yml` does not expose ports to the host. This allows the host to be running other services on those ports, as well as allowing multiple CDN-in-a-Boxes to run on the same host, without port conflicts.

To expose the ports of each service on the host, add the `docker-compose.expose-ports.yml` file. For example, `docker-compose -f docker-compose.yml -f docker-compose.expose-ports.yml up`.

## Common Pitfalls

### Everything's "waiting for Traffic Ops" forever and nothing seems to be working

If you scroll back through the output ( or use `docker-compose logs trafficops-perl | grep "User defined signal 2"` ) and see a line that says something like `/run.sh: line 79: 118 User defined signal 2 $TO_DIR/local/bin/hypnotoad script/cdn` then you've hit a mysterious known error. We don't know what this is or why it happens, but your best bet is to send up a quick prayer and restart the stack.

### Traffic Monitor is stuck waiting for a valid Snapshot

Often times you must take a CDN [Snapshot](https://traffic-control-cdn.readthedocs.io/en/latest/glossary.html#term-snapshot) in order for a valid Snapshot to be generated. This can be done through [Traffic Portal's "CDNs" view](https://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_portal/usingtrafficportal.html#cdns), clicking on the "CDN-in-a-Box" CDN, then pressing the camera button, and finally the "Perform Snapshot" button.

### I'm seeing a failure to open a socket and/or set a socket option

Try disabling SELinux or setting it to 'permissive'. SELinux hates letting containers bind to certain ports. You can also try re-labeling the `docker` executable if you feel comfortable.

### Traffic Vault container exits with cp /usr/local/share/ca-certificates cp: missing destination

Bring all components down, remove the `traffic_ops/ca` directory, and delete the volumes with `docker volume prune`. This will force the regeneration of the certificates.
