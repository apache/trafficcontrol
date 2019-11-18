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

.. _to-api-capabilities-name:

*************************
``capabilities/{{name}}``
*************************

``GET``
=======
Get a capability by name.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------+
	| Name |          Description                   |
	+======+========================================+
	| name | The name of the capability of interest |
	+------+----------------------------------------+

Response Structure
------------------
:description: Describes the APIs covered by the capability
:lastUpdated: Date and time of the last update made to this capability, in ISO format
:name:        Name of the capability

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 20:37:17 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 0YBTC5TEAOJ6B8gsaKgOD1ni2hnZ8Kh9u2JhcmExoGIPaMEKpp4Omr4FglkOQZuh/IB90eJjBMNMeCEvZCxWRg==
	Content-Length: 167

	{ "response": [
		{
			"lastUpdated": "2018-11-14 20:33:00.275376+00",
			"name": "test",
			"description": "This is only a test. If this were a real capability, it might do something"
		}
	]}


``PUT``
=======
Edit a capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------+
	| Name |          Description                      |
	+======+===========================================+
	| name | The name of the capability to be modified |
	+------+-------------------------------------------+

:description: Describes the APIs covered by the capability

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/capabilities/test HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 45
	Content-Type: application/json

	{"description": "A much shorter description"}

Response Structure
------------------
:description: Describes the APIs covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in ISO format
:name:        The name of the capability

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 20:40:33 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: +5mLZ/CJnDkJMbnFviXtVdjwt4bu7ykiMIs73zsnuKV/k4q/d025b2pjYDQkSgtfWPJ73FcusAuBM9TCVT3KsA==
	Content-Length: 181

	{ "alerts": [
		{
			"level": "success",
			"text": "Capability was updated."
		}
	],
	"response": {
		"lastUpdated": "2018-11-14 20:33:00.275376+00",
		"name": "test",
		"description": "A much shorter description"
	}}


``DELETE``
==========
Delete a capability.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters**

	+-----------------+----------+------------------------------------------------+
	| Name            | Required | Description                                    |
	+=================+==========+================================================+
	| ``name``        | yes      | Capability name.                               |
	+-----------------+----------+------------------------------------------------+

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 20:45:37 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: IlAiV4ebwTpMIgeYlR5RuwOhwmHsFs8Ekt7AaEDb3v+lXjvjkqU98xFsfNWvpvPbT/iJnotENhtVq8TVdvoPLg==
	Content-Length: 61

	{ "alerts": [
		{
			"level": "success",
			"text": "Capability deleted."
		}
	]}

