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

.. _to-overview:

Traffic Ops
===========
Traffic Ops is the tool for administration (configuration and monitoring) of all components in a Traffic Control CDN. :ref:`tp-overview` uses the :ref:`to-api` to manage servers, :term:`Cache Groups`, :term:`Delivery Services`, etc. In many cases, a configuration change requires propagation to several, or even all, :term:`cache servers` and only explicitly after or before the same change propagates to :ref:`tr-overview`. Traffic Ops takes care of this required consistency between the different components and their configuration.

Traffic Ops uses a `PostgreSQL <https://www.postgresql.org/>`_ database to store the configuration information, and a combination of the `Mojolicious framework <http://mojolicio.us/>`_ and `Go <https://golang.org/>`_ to provide the :ref:`to-api`. Not all configuration data is in this database however; for sensitive data like private SSL keys or token-based authentication shared secrets, :ref:`tv-overview` is used as a separate, key/value store, allowing administrators to harden the :ref:`tv-overview` server better from a security perspective (i.e only allow Traffic Ops to access it, verifying authenticity with a certificate). The Traffic Ops server, by design, needs to be accessible from all the other servers in the Traffic Control CDN.

Traffic Ops generates all the application-specific configuration files for the :term:`cache servers` and other servers. The :term:`cache servers` and other servers check in with Traffic Ops at a regular interval to see if updated configuration files require application. On :term:`cache servers` this is done by the :term:`ORT` script.

Traffic Ops also runs a collection of periodic checks to determine the operating state of the :term:`cache servers`. These periodic checks are customizable by the Traffic Ops administrative user using `Traffic Ops Extension`_\ s.

.. _trops-ext:

Traffic Ops Extension
---------------------
Traffic Ops Extensions are a way to enhance the basic functionality of Traffic Ops in a custom manner. There are two types of extensions:

:ref:`to-check-ext`
	Allow you to add custom checks to the :menuselection:`Monitor --> Cache Checks` view in :ref:`tp-overview`.
