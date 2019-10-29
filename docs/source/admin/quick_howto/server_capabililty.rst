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

.. _server_capability:

**************************
Manage Server Capabilities
**************************
Server capabilities are designed to enable users with the "operator" :term:`role` ("operators") to control the flow of :term:`delivery service` traffic through :term:`cache servers` (:term:`Edge` or :term:`Mid`) with ONLY the required capabilities. For example, :term:`delivery services` that serve large binary files should only be routed to :term:`cache servers` with sufficient disk cache. Currently, this can be controlled at the :term:`Edge-tier` where system operators can explicitly assign only :term:`Edge-tier caches` with sufficient disk cache to the :term:`delivery service`. However, the operators do not have control of the :term:`Mid-tier` and cannot dictate which :term:`Mid-tier caches` are qualified to serve these large binary files. This will cause a problem if a :term:`Mid-tier cache` with insufficient disk cache is asked to serve the large binary files.

A list of the server capabilities can be found under :menuselection:`Configure --> Server Capabilities`. Users with a higher-level :term:`role` ("operations" or "admin") can create or delete server capabilities. Server capabilities can only be deleted if they are not currently being used by a :term:`cache server` or required by a :term:`delivery service`.

.. figure:: server_capability/server_caps_table.png
	:align: center
	:alt: A screenshot of the Traffic Portal UI depicting an example list of Server Capabilities

	Example Server Capabilities Listing

Assigning a Server Capability to a Server
=========================================
Users with the Operations or Admin :term:`Role` can associate one or more server capabilities with a server by navigating to a server via :menuselection:`Configure --> Servers` and using the context menu for the server table and selecting :menuselection:`Manage Capabilities` or by navigating to :menuselection:`Configure --> Servers --> Server --> More --> Manage Capabilities`.

.. figure:: server_capability/server_server_caps_table.png
	:align: center
	:alt: A screenshot of the Traffic Portal UI depicting an example list of Server Capabilities attached to a Server

	Example Server Capabilities for a Server Listing


