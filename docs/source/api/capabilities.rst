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

.. _to-api-capabilities:

****************
``capabilities``
****************

``GET``
=======
Get all capabilities.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                         |
	+===========+==========+=====================================================================================================================+
	| name      | no       | Return only the capability that has this name                                                                       |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` array |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                            |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                      |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with ``limit``            |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the |
	|           |          | first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make   |
	|           |          | use of ``page``.                                                                                                    |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------------+


.. code-block:: http
	:caption: Request Example

	GET /api/1.4/capabilities?name=test HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:name:        Name of the capability
:description: Describes the permissions covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in an ISO-like format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Vary: Accept-Encoding
	Transfer-Encoding: chunked
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: c18+GtX2ZI8PoCSwuAzBhl+6w3vDpKQTa/cDJC0WHxdpguOL378KBxGWW5PCSyZfJUb7wPyOL5qKMn6NNTufhg==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Aug 2019 17:20:20 GMT
	Content-Length: 161

	{ "response": [
		{
			"description": "This is only a test. If this were a real capability, it might do something",
			"lastUpdated": "2019-08-15 17:18:03+00",
			"name": "test"
		}
	]}

``POST``
========
Create a capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object


Request Structure
-----------------
:name:        The name of the capability being created
:description: A description of what the capability allows

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 108
	Content-Type: application/json

	{
		"name": "test",
		"description": "This is only a test. If this were a real capability, it might do something"
	}

Response Structure
------------------
:description: Describes the permissions covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in an ISO-like format
:name:        Name of the capability

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: A1rjpDy+O+oooYeer2j09pCEDpPEFk/nt8/AaJye2sLkfy93MtquCsB/Rlgz7sCYputd/EPOPDyi2WkN8UB1Rw==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Aug 2019 17:18:03 GMT
	Content-Length: 219

	{ "alerts": [
		{
			"text": "Capability created.",
			"level": "success"
		}
	],
	"response": {
		"description": "This is only a test. If this were a real capability, it might do something",
		"lastUpdated": "2019-08-15 17:18:03+00",
		"name": "test"
	}}


``PUT``
=======
.. versionadded:: 1.4

Replace a capability with the one provided.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+---------------------------------------------------+
	| Name | Required | Description                                       |
	+======+==========+===================================================+
	| name | yes      | The (current) name of the capability to be edited |
	+------+----------+---------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/capabilities?name=test HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 109
	Content-Type: application/json

Response Structure
------------------
:description: Describes the permissions covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in an ISO-like format
:name:        Name of the capability

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: eciuE8oKQqBOtMThcQvSrPEIuJ9gUeutB00eW7g4KSscwO/vzplyOg8i/EVgfR9NFhK9VSVvdrKvxHC7HsG2fg==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Aug 2019 17:21:50 GMT
	Content-Length: 224

	{ "alerts": [
		{
			"text": "Capability was updated.",
			"level": "success"
		}
	],
	"response": {
		"description": "This is only a test. If this were a real capability, it might do something",
		"lastUpdated": "2019-08-15 17:21:50+00",
		"name": "quest"
	}}

``DELETE``
==========
.. versionadded:: 1.4

Delete a capability.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+------------------------------------------+
	| Name | Required | Description                              |
	+======+==========+==========================================+
	| name | yes      | The name of the capability to be deleted |
	+------+----------+------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/capabilities?name=quest HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description: Describes the permissions that were covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in an ISO-like format
:name:        Name of the capability

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: 7lWTuaI1BUeXrnTG1fbFeKuvVuqojZJjSQV5MOtT0a++VV1PUAXYSIwe2vUOpoM4uwCKpeAc86J75OJGLgLHdg==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Aug 2019 17:26:00 GMT
	Content-Length: 220

	{ "alerts": [
		{
			"text": "Capability deleted.",
			"level": "success"
		}
	],
	"response": {
		"description": "This is only a test. If this were a real capability, it might do something",
		"lastUpdated": "2019-08-15 17:21:50+00",
		"name": "quest"
	}}
