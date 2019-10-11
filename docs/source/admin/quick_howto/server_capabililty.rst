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
Server capabilities are designed to enable system operators to control the flow of delivery service traffic through caches (Edges or Mids) with ONLY the required capabilities. For example, delivery services that serve large binary files should only be routed to caches with sufficient disk cache. Currently, this can be controlled at the Edge tier where system operators can explicitly assign only Edge caches with sufficient disk cache to the delivery service. However, the system operators do not have control of the Mid tier and cannot dictate which Mid caches are qualified to serve these large binary files. This will cause a problem if a Mid cache with insufficient disk cache is asked to serve the large binary files.

A list of the server capabilities can be found under :menuselection:`Configure --> Server Capabilities`. Users with a higher level :term:`Role` ("operations" or "admin") can create or delete server capabilities assuming they are not currently being used by a server or required by a delivery service.

.. figure:: server_capability/server_caps_table.png
	:align: center
	:alt: A screenshot of the Traffic Portal UI depicting an example list of Server Capabilities

	Example Server Capabilities Listing
