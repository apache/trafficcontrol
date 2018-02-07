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

Administrator's Guide
*********************

Traffic Control is distributed in source form for the developer, but also as a binary package. This guide details how to install and configure a Traffic Control CDN using the binary packages, as well as how to perform common operations running a CDN.

When installing a complete CDN from scratch, a sample recommended order is:

#. Traffic Ops
#. Traffic Vault (Riak)
#. Traffic Monitor
#. Apache Traffic Server Mid-Tier Caches
#. Apache Traffic Server Edge Caches
#. Traffic Router
#. Traffic Stats
#. Traffic Portal

Once everything is installed, you will need to configure the servers to talk to each other. You will also need Origin server(s), which the Mid-Tier Cache(s) get content from. An Origin server is simply an HTTP(S) server which serves the content you wish to cache on the CDN.

.. toctree::
  :maxdepth: 3

  traffic_ops/installation.rst
  traffic_ops/default_profiles.rst
  traffic_ops/migration_from_10_to_20.rst
  traffic_ops/migration_from_20_to_22.rst
  traffic_ops/configuration.rst
  traffic_ops/using.rst
  traffic_ops/extensions.rst
  traffic_portal/installation.rst
  traffic_portal/usingtrafficportal.rst
  traffic_monitor.rst
  traffic_monitor_golang.rst
  traffic_router.rst
  traffic_stats.rst
  traffic_server.rst
  traffic_vault.rst
  quick_howto/index.rst
