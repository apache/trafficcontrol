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

.. _to-api-multiple_server_capabilities:

******************************
``multiple_server_capabilities``
******************************

``PUT``
========
Associates a list of :term:`Server Capability` to a server.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:UPDATE, SERVER:CREATE, SERVER:READ, SERVER-CAPABILITY:READ
:Response Type:  Object

Request Structure
-----------------
:serverId:         The integral, unique identifier of a server to be associated with a :term:`Server Capability`
:serverCapability: List of :term:`Server Capability`'s name to associate

.. code-block:: http
	:caption: Request Example

	PUT /api/4.1/multiple_server_capabilities/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
        "serverId": 1,
        "serverCapability": ["test", "disk"]
    }

Response Structure
------------------
:serverId:         The integral, unique identifier of the newly associated server
:serverCapability: List of :term:`Server Capability`'s name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 8 Aug 2022 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 02 Aug 2022 22:15:11 GMT
	Content-Length: 157

	{
	"alerts": [{
		"text": "Multiple Server Capabilities assigned to a server:1",
		"level": "success"
	}],
	"response": {
		"serverId": 1,
		"serverCapability": ["test", "disk"]
	}
}
