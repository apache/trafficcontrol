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
No available parameters

Response Structure
------------------
:name:        Name of the capability
:description: Describes the APIs covered by the capability.
:lastUpdated: Date and time of the last update made to this capability, in ISO format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 20:26:19 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: zmjsQO3Y4r1/xCFOHB+E+8+bbgDyVcvoR0d4gKqqsWTFaUnxp2flIzuFqWjXf+wb4Bbd1e2Ojse4nQKnyIFKGw==
	Transfer-Encoding: chunked

	{ "response": [
		{
			"name": "cdn-read",
			"description": "View CDN configuration",
			"lastUpdated": "2017-04-02 08:22:43"
		},
		{
			"name": "cdn-write",
			"description": "Create, edit or delete CDN configuration",
			"lastUpdated": "2017-04-02 08:22:43"
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
	Content-Length: 109
	Content-Type: application/json

	{
		"name": "test",
		"description": "This is only a test. If this were a real capability, it might do something"
	}

Response Structure
------------------
:description: Describes the APIs covered by the capability.
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
	Date: Wed, 14 Nov 2018 20:33:00 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: HhhQzw3JBLv90lOeeSGj75uknADanz3fUnQt1E266HAKPTFuTjuIJpf8ni9fb9Chv9LN7mt16utcHMbP8MBHZw==
	Content-Length: 183

	{ "alerts": [
		{
			"level": "success",
			"text": "Capability was created."
		}
	],
	"response": {
		"name": "test",
		"description": "This is only a test. If this were a real capability, it might do something"
	}}

