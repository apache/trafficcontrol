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

.. _to-api-federations-id-deliveryservices:

***************************************
``federations/{{ID}}/deliveryservices``
***************************************

``GET``
=======
Retrieves :term:`Delivery Service`\ s assigned to a federation.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------+
	| Name |                 Description                                        |
	+======+====================================================================+
	|  ID  | The integral, unique identifier for the federation to be inspected |
	+------+--------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/federations/1/deliveryservices HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cdn:   The CDN to which this :term:`Delivery Service` Belongs
:id:    The integral, unique identifier for the Deliver Service
:type:  The routing type used by this :term:`Delivery Service`
:xmlId: The 'xml_id' which uniquely identifies this :term:`Delivery Service`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:44:13 GMT
	server: Mojolicious (Perl)
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:44:13 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: 7Y9Q/qHeXfbjJduvucRCR85wf4VRfyYhlK59sNRkzIJuwnsMhFcEfYfNqrvELwfexOum/VEX2f/1oa+I/edGfw==
	content-length: 74

	{ "response": [
		{
			"xmlId": "demo1",
			"cdn": "CDN-in-a-Box",
			"type": "HTTP",
			"id": 1
		}
	]}

``POST``
========
Assigns one or more :term:`Delivery Service`\ s to a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------+
	| Name |                 Description                                        |
	+======+====================================================================+
	|  ID  | The integral, unique identifier for the federation to be inspected |
	+------+--------------------------------------------------------------------+

:dsIds:   An array of integral, unique identifiers for :term:`Delivery Service`\ s which will be assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, will cause any conflicting assignments already in place to be overridden by this request

	.. note:: If ``replace`` is not given (and/or not ``true``), then any conflicts with existing assignments will cause the entire operation to fail.

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/federations/1/deliveryservices HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 32
	Content-Type: application/json

	{
		"dsIds": [1],
		"replace": true
	}

Response Structure
------------------
:dsIds:   An array of integral, unique identifiers for :term:`Delivery Service`\ s which are now assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, means any conflicting assignments already in place were overridden by this request

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	content-type: application/json
	set-cookie: mojolicious=...; Path=/; HttpOnly
	whole-content-sha512: rVd0nx8G3bRI8ub1zw6FTdmwQ7jer4zoqzOZf5tC1ckrR0HEIOH1Azdcmvv0FVE5I0omcHVnrYbzab7tUtmnog==
	x-server-name: traffic_ops_golang/
	content-length: 137
	date: Wed, 05 Dec 2018 00:34:06 GMT

	{ "alerts": [
		{
			"text": "1 delivery service(s) were assigned to the federation 1",
			"level": "success"
		}
	],
	"response": {
		"dsIds": [
			1
		],
		"replace": true
	}}

