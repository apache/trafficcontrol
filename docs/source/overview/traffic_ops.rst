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

Traffic Ops
===========
Traffic Ops is the tool for administration (configuration and monitoring) of all components in a Traffic Control CDN. Traffic Portal uses Traffic Ops APIs to manage servers, cache groups, delivery services, etc. In many cases, a configuration change requires propagation to several, or even all, caches and only explicitly after or before the same change propagates to Traffic Router. Traffic Ops takes care of this required consistency between the different components and their configuration. Traffic Ops exposes its data through a series of HTTP APIs. Traffic Ops also provides a legacy user interface that has been superseded by the Traffic Portal UI.

Traffic Ops uses a PostgreSQL database to store the configuration information, and the `Mojolicious framework <http://mojolicio.us/>`_ and `Go <https://golang.org/>`_ to generate APIs used by the Traffic Portal. Not all configuration data is in this database however; for sensitive data like SSL private keys or token based authentication shared secrets, a separate key/value store is used, allowing the operator to harden the server that runs this key/value store better from a security perspective (i.e only allow Traffic Ops access it with a cert). The Traffic Ops server, by design, needs to be accessible from all the other servers in the Traffic Control CDN.

Traffic Ops generates all the application specific configuration files for the caches and other servers. The caches and other servers check in with Traffic Ops at a regular interval (default 15 minutes) to see if updated configuration files require application. This is done by the Traffic Ops Operations Readiness Test (ORT) script.

Traffic Ops also runs a collection of periodic checks to determine the operational readiness of the caches. These periodic checks are customizable by the Traffic Ops administrative user using extensions.

Traffic Ops is in the process of migrating from Perl to Go, and currently runs as two applications. The Go application serves all endpoints which have been rewritten in the Go language, and transparently proxies all other requests to the old Perl application. Both applications are installed by the RPM, and both run as a single service. When the project has fully migrated to Go, the Perl application will be removed, and the RPM and service will consist solely of the Go application.

.. _trops-ext:

Traffic Ops Extension
---------------------
Traffic Ops Extensions are a way to enhance the basic functionality of Traffic Ops in a custom manner. There are two types of extensions:

:ref:`to-check-ext`
	Allow you to add custom checks to the 'Monitor' -> 'Cache Checks' view.

:ref:`to-datasource-ext`
	Allow you to add data sources for the graph views and usage APIs.


.. These are listed as "in beta" as far back as TO 1.0, sooo
.. Configuration Extension
.. 	Allows you to add custom configuration file generators.
