..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

*********************
Administrator's Guide
*********************

Traffic Control is distributed in source form for the developer, but also as a binary package. This guide details how to install and configure a Traffic Control CDN using the binary packages, as well as how to perform common operations running a CDN.

When installing a complete CDN from scratch, a sample recommended order is:

#. Traffic Ops DB (PostgreSQL)
#. `InfluxDB [Optional] <https://github.com/influxdata/influxdb>`_
#. Traffic Vault (PostgreSQL or Riak) [Optional]
#. Fake Origin [Optional]
#. Traffic Ops
#. Traffic Portal
#. Initial Traffic Ops Dataset Setup [if applicable]
#. Traffic Monitor
#. Apache Traffic Server Mid-Tier Caches
#. Apache Traffic Server Edge-Tier Caches
#. Grove [Optional]
#. Traffic Router
#. `InfluxDB-relay [Optional] <https://github.com/influxdata/influxdb-relay>`_
#. Traffic Stats [Optional]

Once everything is installed, you will need to configure the servers to talk to each other. You will also need Origin server(s), from which the Mid-Tier Cache(s) will obtain content. An Origin server is simply an HTTP(S) server which serves the content you wish to cache on the CDN.

.. toctree::
	:maxdepth: 3
	:glob:

	traffic_ops.rst
	environment_creation.rst
	traffic_portal/index.rst
	traffic_monitor.rst
	traffic_router.rst
	traffic_stats.rst
	traffic_server.rst
	t3c/index.rst
	traffic_vault.rst
	quick_howto/index.rst
	cdni.rst
