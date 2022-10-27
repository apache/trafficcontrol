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

.. _to-api-multiple_servers_capabilities:

*********************************
``multiple_servers_capabilities``
*********************************

``POST``
========
Inserts a list of :term:`Server Capability` names associated to a server and vice versa i.e insert a list of :term:`Server` ids associated to a server capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:READ, SERVER:CREATE, SERVER-CAPABILITY:READ, SERVER-CAPABILITY:CREATE
:Response Type:  Object

Request Structure
-----------------
:serverIds:          List of :term:`Server` ids (integral, unique identifier) associated with a :term:`Server Capability`
:serverCapabilities: List of :term:`Server Capability` names to associate with a :term:`Server` id
:pageType:           To determine which configuration (server or server capabilities) is requesting association. Only two values are permitted: `server` or `sc` (short for server capability). If `server` is chosen, it implies that multiple server capabilities are to be associated with a given server id. If `sc` is chosen, it implies that multiple server ids are to be associated with a given server capability.

.. code-block:: http
	:caption: Request Example1

	POST /api/5.0/multiple_servers_capabilities/ HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
		"serverIds": [1],
		"serverCapabilities": ["test", "disk"]
		"pageType": "server"
	}

.. code-block:: http
	:caption: Request Example2

	POST /api/5.0/multiple_servers_capabilities/ HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
		"serverIds": [2, 3]
		"serverCapabilities": ["eas"],
		"pageType": "sc"
	}

Response Structure
------------------
:serverId:           List of :term:`Server` ids (integral, unique identifier) associated with a server capability.
:serverCapabilities: List of :term:`Server Capability` names to be associated with a :term:`Server` id.

.. code-block:: http
	:caption: Response Example1

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 8 Aug 2022 22:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 08 Aug 2022 16:15:11 GMT
	Content-Length: 157

	{
		"alerts": [{
			"text": "Multiple Server Capabilities assigned to a server",
			"level": "success"
		}],
		"response": {
			"serverIds": [1],
			"serverCapabilities": ["test", "disk"]
		}
	}

.. code-block:: http
	:caption: Response Example2

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 8 Aug 2022 22:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 08 Aug 2022 16:15:11 GMT
	Content-Length: 157

	{
		"alerts": [{
			"text": "Multiple Servers assigned to a capability",
			"level": "success"
		}],
		"response": {
			"serverIds": [2, 3]
			"serverCapabilities": ["eas"],
			"pageType": "sc"
		}
	}

``DELETE``
==========
Deletes a list of :term:`Server Capability` names associated to a server and vice versa i.e. deletes a list of :term:`Server` ids associated to a server capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:READ, SERVER:DELETE, SERVER-CAPABILITY:READ, SERVER-CAPABILITY:DELETE
:Response Type:  Object

Request Structure
-----------------
:serverIds:          List of :term:`Server` ids (integral, unique identifier) associated with a :term:`Server Capability`
:serverCapabilities: List of :term:`Server Capability` names to associate with a :term:`Server` id
:pageType:           To determine which configuration (server or server capabilities) is requesting deletion. Only two values are permitted: `server` or `sc` (short for server capability). If `server` is chosen, it implies that multiple server capabilities are to be deleted for a given server id. If `sc` is chosen, it implies that multiple server ids are to be deleted for a given server capability.

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/multiple_servers_capabilities/ HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
		"serverIds": [2, 3]
		"serverCapabilities": ["eas"],
	}

Response Structure
------------------
:serverId:           List of :term:`Server` ids (integral, unique identifier) associated with a server capability.
:serverCapabilities: List of :term:`Server Capability` names to be associated with a :term:`Server` id.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 8 Aug 2022 22:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 08 Aug 2022 16:15:11 GMT
	Content-Length: 157

	{
		"alerts": [{
			"text": "Removed multiple servers from capabilities or multiple servers to a capability",
			"level": "success"
		}],
		"response": {
			"serverIds": [2, 3]
			"serverCapabilities": ["eas"],
		}
	}

