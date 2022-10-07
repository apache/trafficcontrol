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

.. _to-api-v4-server_capabilities:

***********************
``server_capabilities``
***********************

``GET``
=======
Retrieves :term:`Server Capabilities`.

:Auth. Required: Yes
:Roles Required: "read-only"
:Permissions Required: SERVER-CAPABILITY:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+-----------------------------------------------------+
	| Name | Required | Description                                         |
	+======+==========+=====================================================+
	| name | no       | Return the :term:`Server Capability` with this name |
	+------+----------+-----------------------------------------------------+

.. code-block:: http
	:caption: Request Structure

	GET /api/4.0/server_capabilities?name=RAM HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:name:        The name of this :term:`Server Capability`
:lastUpdated: The date and time at which this :term:`Server Capability` was last updated, in ISO-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: EH8jo8OrCu79Tz9xpgT3YRyKJ/p2NcTmbS3huwtqRByHz9H6qZLQjA59RIPaVSq3ZxsU6QhTaox5nBkQ9LPSAA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 21:36:13 GMT
	Content-Length: 68

	{
		"response": [
			{
				"name": "RAM",
				"lastUpdated": "2019-10-07 20:38:24+00"
			}
		]
	}

``POST``
========
Create a new :term:`Server Capability`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER-CAPABILITY:CREATE, SERVER-CAPABILITY:READ
:Response Type:  Object

Request Structure
-----------------
:name: The name of the :term:`Server Capability`

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/server_capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 15
	Content-Type: application/json

	{
		"name": "RAM"
	}

Response Structure
------------------
:name:        The name of this :term:`Server Capability`
:lastUpdated: The date and time at which this :term:`Server Capability` was last updated, in ISO-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: ysdopC//JQI79BRUa61s6M2HzHxYHpo5RdcuauOoqCYxiVOoUhNZfOVydVkv8zDN2qA374XKnym4kWj3VzQIXg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 22:10:00 GMT
	Content-Length: 137

	{
		"alerts": [
			{
				"text": "server capability was created.",
				"level": "success"
			}
		],
		"response": {
			"name": "RAM",
			"lastUpdated": "2019-10-07 22:10:00+00"
		}
	}

``PUT``
========
Update an existing :term:`Server Capability`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER-CAPABILITY:UPDATE, SERVER-CAPABILITY:READ
:Response Type:  Object

Request Structure
-----------------
:name: The name of the :term:`Server Capability`

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/server_capabilities?name=RAM HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 15
	Content-Type: application/json

	{
		"name": "HDD"
	}

Response Structure
------------------
:name:        The name of this :term:`Server Capability`
:lastUpdated: The date and time at which this :term:`Server Capability` was last updated, in ISO-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: ysdopC//JQI79BRUa61s6M2HzHxYHpo5RdcuauOoqCYxiVOoUhNZfOVydVkv8zDN2qA374XKnym4kWj3VzQIXg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 03 March 2021 21:22:08 GMT
	Content-Length: 137

	{
		"alerts": [
			{
				"text": "server capability was updated.",
				"level": "success"
			}
		],
		"response": {
			"name": "HDD",
			"lastUpdated": "2021-03-03 21:22:08+00"
		}
	}

``DELETE``
==========
Deletes a specific :term:`Server Capability`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER-CAPABILITY:DELETE, SERVER-CAPABILITY:READ
:Response Type:  ``undefined``


Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+---------------------------------------------------------+
	| Name | Required | Description                                             |
	+======+==========+=========================================================+
	| name | yes      | The name of the :term:`Server Capability` to be deleted |
	+------+----------+---------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/server_capabilities?name=RAM HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 8zCAATbCzcqiqigGVBy7WF1duDuXu1Wg2DBe9yfqTw/c+yhE2eUk73hFTA/Oqt0kocaN7+1GkbFdPkQPvbnRaA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 20:44:40 GMT
	Content-Length: 72

	{
		"alerts": [
			{
				"text": "server capability was deleted.",
				"level": "success"
			}
		]
	}
