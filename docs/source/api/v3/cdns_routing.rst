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

.. _to-api-v3-cdns-routing:

****************
``cdns/routing``
****************

``GET``
=======
Retrieves the aggregated routing percentages across all CDNs. This is accomplished by polling stats from all online Traffic Routers via the ``/crs/stats`` route.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:cz:                The percent of requests to online Traffic Routers that were satisfied by a :term:`Coverage Zone File`
:deepCz:            The percent of requests to online Traffic Routers that were satisfied by a :term:`Deep Coverage Zone File`
:dsr:               The percent of requests to online Traffic Routers that were satisfied by sending the client to an overflow :term:`Delivery Service`
:err:               The percent of requests to online Traffic Routers that resulted in an error
:fed:               The percent of requests to online Traffic Routers that were satisfied by sending the client to a federated CDN
:geo:               The percent of requests to online Traffic Routers that were satisfied using 3rd party geographic IP mapping
:miss:              The percent of requests to online Traffic Routers that could not be satisfied
:regionalAlternate: The percent of requests to online Traffic Routers that were satisfied by sending the client to the alternate, Regional Geo-blocking URL
:regionalDenied:    The percent of requests to online Traffic Routers that were denied due to geographic location policy
:staticRoute:       The percent of requests to online Traffic Routers that were satisfied with :ref:`ds-static-dns-entries`

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
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 7LjytwKyRzSKM4cRIol4OMIJxApFpTWJaSK73rbtUIQdASZjI64XxLVzZP0OGRU7XeJ22YKUyQ30qbKHDRv7FQ==
	Content-Length: 130

	{ "response": {
		"cz": 79,
		"deepCz": 0.50,
		"dsr": 0,
		"err": 0,
		"fed": 0.25,
		"geo": 20,
		"miss": 0.25,
		"regionalAlternate": 0,
		"regionalDenied": 0,
		"staticRoute": 0
	}}
