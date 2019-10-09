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

.. _to-api-deliveryservices-id-routing:

***********************************
``deliveryservices/{{ID}}/routing``
***********************************

``GET``
=======
Retrieves routing method statistics for a particular :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------+
	| Name | Description                                                                  |
	+======+==============================================================================+
	|  ID  | The integral, unique identifier for the :term:`Delivery Service` of interest |
	+------+------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/deliveryservices/1/routing HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cz:                The percent of requests to the Traffic Router for this :term:`Delivery Service` that were satisfied by a :term:`Coverage Zone File`
:dsr:               The percent of requests to the Traffic Router for this :term:`Delivery Service` that were satisfied by sending the client to an overflow :term:`Delivery Service`
:err:               The percent of requests to the Traffic Router for this :term:`Delivery Service` that resulted in an error
:fed:               The percent of requests to the Traffic Router for this :term:`Delivery Service` that were satisfied by sending the client to a federated CDN
:geo:               The percent of requests to the Traffic Router for this :term:`Delivery Service` that were satisfied using 3rd party geographic IP mapping
:miss:              The percent of requests to the Traffic Router for this :term:`Delivery Service` that could not be satisfied
:regionalAlternate: The percent of requests to the Traffic Router for this :term:`Delivery Service` that were satisfied by sending the client to the alternate, Regional Geo-blocking URL
:regionalDenied:    The percent of Traffic Router requests for this :term:`Delivery Service` that were denied due to geographic location policy
:staticRoute:       The percent of requests to the Traffic Router for this :term:`Delivery Service` that were satisfied with :ref:`ds-static-dns-entries`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Fri, 30 Nov 2018 15:08:07 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Fri, 30 Nov 2018 19:08:07 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: UgPziRC/5u4+CfkZ9xm0EkEzjjJVu6cwBrFd/n3xH/ZmlkaXkQaa1y4+B7DyE46vxFLYE0ODOcQchyn7JkoQOg==
	Content-Length: 132

	{ "response": {
		"staticRoute": 0,
		"geo": 100,
		"err": 0,
		"fed": 0,
		"cz": 0,
		"dsr": 0,
		"regionalAlternate": 0,
		"deepCz": 0,
		"regionalDenied": 0,
		"miss": 0
	}}

.. [#tenancy] Users will only be able to view routing details for the :term:`Delivery Services` their :term:`Tenant` is allowed to see.
