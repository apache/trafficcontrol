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

.. _to-api-v3-servercheck:

***************
``servercheck``
***************

.. seealso:: :ref:`to-check-ext`

``GET``
=======
Fetches identifying and meta information as well as "check" values regarding all servers that have a :term:`Type` with a name beginning with "EDGE" or "MID" (ostensibly this is equivalent to all :term:`cache servers`).

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                        |
	+===========+==========+====================================================================================+
	| id        | no       | Return only :term:`cache servers` with this integral, unique identifier (id)       |
	+-----------+----------+------------------------------------------------------------------------------------+
	| hostName  | no       | Return only :term:`cache servers` with this host_name                              |
	+-----------+----------+------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example with ``hostName`` query param

	GET /api/4.0/servercheck?hostName=edge HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

.. code-block:: http
	:caption: Request Example with ``id`` query param

	GET /api/4.0/servercheck?id=12 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:adminState:   The name of the server's :term:`Status` - called "adminState" for legacy reasons
:cacheGroup:   The name of the :term:`Cache Group` to which the server belongs
:checks:       An optionally present map of the names of "checks" to their values. Only numeric and boolean checks are represented, and boolean checks are represented as integers with ``0`` meaning "false" and ``1`` meaning "true". Will not appear if the server in question has no valued "checks".
:hostName:     The (short) hostname of the server
:id:           The server's integral, unique identifier
:profile:      The name of the :term:`Profile` used by the server
:revalPending: A boolean that indicates whether or not the server has pending revalidations
:type:         The name of the server's :term:`Type`
:updPending:   A boolean that indicates whether or not the server has pending updates

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Thu, 18 Feb 2021 20:00:19 GMT; Max-Age=3600; HttpOnly
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 18 Feb 2021 19:00:19 GMT
	Content-Length: 352

	{ "response": [
		{
			"adminState": "REPORTED",
			"cacheGroup": "CDN_in_a_Box_Edge",
			"id": 12,
			"hostName": "edge",
			"revalPending": false,
			"profile": "ATS_EDGE_TIER_CACHE",
			"type": "EDGE",
			"updPending": false
		}
	]}

``POST``
========
Post a server check result to the "serverchecks" table. Updates the resulting value from running a given check extension on a server.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type: Object

Request Structure
-----------------
The request only requires to have either ``host_name`` or ``id`` defined.

:host_name:              The hostname of the server to which this "servercheck" refers.
:id:                     The id of the server to which this "servercheck" refers.
:servercheck_short_name: The short name of the "servercheck".
:value:                  The value of the "servercheck"

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/servercheck HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 113
	Content-Type: application/json

	{
		"id": 1,
		"host_name": "edge",
		"servercheck_short_name": "test",
		"value": 1
	}

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Server Check was successfully updated."
		}
	]}

.. [1] No roles are required to use this endpoint, however access is controlled by username. Only the reserved user ``extension`` is permitted the use of this endpoint.

