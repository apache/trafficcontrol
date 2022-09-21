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

.. _to-api-v4-multiple_servers_per_capability:

***********************************
``multiple_servers_per_capability``
***********************************

``PUT``
========
Associates a list of :term:`Servers` to a server capability. The API call replaces all the servers assigned to a server capability with the ones specified in the servers field.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:READ, SERVER-CAPABILITY:READ, SERVER-CAPABILITY:UPDATE
:Response Type:  Object

Request Structure
-----------------
:serverCapability:  The unique identifier of a server capability to be associated with a :term:`Server`
:serverIds:         List of :term:`Server` ids associated with a server capability

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/multiple_server_capabilities/ HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
		"serverCapability": "eas",
		"serverIds": [2, 3]
	}

Response Structure
------------------
:serverCapability:   The unique identifier of a server capability to be associated with a :term:`Server`
:serverIds:          List of :term:`Server` ids associated with a server capability

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Sat, 24 Sep 2022 22:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Tues, 20 Sep 2022 16:15:11 GMT
	Content-Length: 157

	{
	"alerts": [{
			"text": "Multiple Servers assigned to a capability",
			"level": "success"
		}],
		"response": {
			"serverCapability": "eas",
			"serverIds": [2, 3]
		}
	}
