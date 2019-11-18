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

.. _to-api-cdns-routing:

****************
``cdns/routing``
****************

``GET``
=======
Retrieves the aggregate routing percentages of Cache Groups assigned to any CDN.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:cz:          Used Coverage Zone geographic IP mapping
:dsr:         Overflow traffic sent to secondary CDN
:err:         Error localizing client IP
:geo:         Used 3rd party geographic IP mapping
:miss:        No location available for client IP
:staticRoute: Used pre-configured DNS entries

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 21:29:32 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 7LjytwKyRzSKM4cRIol4OMIJxApFpTWJaSK73rbtUIQdASZjI64XxLVzZP0OGRU7XeJ22YKUyQ30qbKHDRv7FQ==
	Content-Length: 130

	{ "response": {
		"staticRoute": 0,
		"geo": 20.6251834458468,
		"err": 0,
		"fed": 0.287643087760493,
		"cz": 79.0607572644555,
		"regionalAlternate": 0,
		"dsr": 0,
		"miss": 0.0264162019371881,
		"regionalDenied": 0
	}}

