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

.. _to-api-v4-steering-id-targets:

***************************
``steering/{{ID}}/targets``
***************************

``GET``
=======
Get all targets for a steering :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: STEERING:READ, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------------------------------------------------------------+
	| Name |                Description                                                                               |
	+======+==========================================================================================================+
	|  ID  | The integral, unique identifier of a steering :term:`Delivery Service` for which targets shall be listed |
	+------+----------------------------------------------------------------------------------------------------------+

.. table:: Request Query Parameters

	+-----------+-------------------------------------------------------------------------------------------------------------------------+
	| Name      | Description                                                                                                             |
	+===========+=========================================================================================================================+
	| target    | Return only the target mappings that target the :term:`Delivery Service` identified by this integral, unique identifier |
	+-----------+-------------------------------------------------------------------------------------------------------------------------+
	| orderby   | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` array     |
	+-----------+-------------------------------------------------------------------------------------------------------------------------+
	| sortOrder | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                |
	+-----------+-------------------------------------------------------------------------------------------------------------------------+
	| limit     | Choose the maximum number of results to return                                                                          |
	+-----------+-------------------------------------------------------------------------------------------------------------------------+
	| offset    | The number of results to skip before beginning to return results. Must use in conjunction with limit                    |
	+-----------+-------------------------------------------------------------------------------------------------------------------------+
	| page      | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the     |
	|           | first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use   |
	|           | of ``page``.                                                                                                            |
	+-----------+-------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Structure

	GET /api/4.0/steering/2/targets?target=1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:deliveryService:   A string that is the :ref:`ds-xmlid` of the steering :term:`Delivery Service`
:deliveryServiceId: An integral, unique identifier for the steering :term:`Delivery Service`
:target:            A string that is the :ref:`ds-xmlid` of this target :term:`Delivery Service`
:targetId:          An integral, unique identifier for this target :term:`Delivery Service`
:type:              The steering type of this target :term:`Delivery Service`. This should be one of ``STEERING_WEIGHT``, ``STEERING_ORDER``, ``STEERING_GEO_ORDER`` or ``STEERING_GEO_WEIGHT``
:typeId:            An integral, unique identifier for the :ref:`steering type <ds-steering>` of this target :term:`Delivery Service`
:value:             The 'weight', 'order', 'geo_order' or 'geo_weight' attributed to this steering target as an integer

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: utlJK4oYS2l6Ff7NzAqRuQeMEtazYn3rM3Nlux2XgTLxvSyslHy0mJrwDExSU05gVMdrgYCLZrZEvPHlENT1nA==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 14:09:23 GMT
	Content-Length: 130

	{ "response": [
		{
			"deliveryService": "test",
			"deliveryServiceId": 2,
			"target": "demo1",
			"targetId": 1,
			"type": "STEERING_ORDER",
			"typeId": 44,
			"value": 100
		}
	]}

``POST``
========
Create a steering target.

:Auth. Required: Yes
:Roles Required: Portal, Steering, Federation, "operations" or "admin"
:Permissions Required: STEERING:CREATE, STEERING:READ, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------------------+
	| Name |                Description                                                                              |
	+======+=========================================================================================================+
	|  ID  | The integral, unique identifier of a steering :term:`Delivery Service` to which a target shall be added |
	+------+---------------------------------------------------------------------------------------------------------+

:targetId: The integral, unique identifier of a :term:`Delivery Service` which shall be a new steering target for the :term:`Delivery Service` identified by the ``ID`` path parameter
:typeId:   The integral, unique identifier of the steering type of the new target :term:`Delivery Service`. This should be corresponding to one of ``STEERING_WEIGHT``, ``STEERING_ORDER``, ``STEERING_GEO_ORDER`` or ``STEERING_GEO_WEIGHT``
:value:    The 'weight', 'order', 'geo_order' or 'geo_weight' which shall be attributed to the new target :term:`Delivery Service`

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/steering/2/targets HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 43
	Content-Type: application/json

	{
		"targetId": 1,
		"value": 100,
		"typeId": 43
	}

Response Structure
------------------
:deliveryService:   A string that is the :ref:`ds-xmlid` of the steering :term:`Delivery Service`
:deliveryServiceId: An integral, unique identifier for the steering :term:`Delivery Service`
:target:            A string that is the :ref:`ds-xmlid` of this target :term:`Delivery Service`
:targetId:          An integral, unique identifier for this target :term:`Delivery Service`
:type:              The steering type of this target :term:`Delivery Service`. This should be one of ``STEERING_WEIGHT``, ``STEERING_ORDER``, ``STEERING_GEO_ORDER`` or ``STEERING_GEO_WEIGHT``
:typeId:            An integral, unique identifier for the :ref:`steering type <ds-steering>` of this target :term:`Delivery Service`
:value:             The 'weight', 'order', 'geo_order' or 'geo_weight' attributed to this steering target as an integer

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: +dTvfzrnOhdwAOMmY28r0+gFV5z+3aABI2FfAMziTYcU+pZrDanrJzMXpKWIL5Q/oCUBZpJDRt9hRCFkT4oGYw==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 21:22:17 GMT
	Content-Length: 196

	{ "alerts": [
		{
			"text": "steeringtarget was created.",
			"level": "success"
		}
	],
	"response": {
		"deliveryService": "test",
		"deliveryServiceId": 2,
		"target": "demo1",
		"targetId": 1,
		"type": "HTTP",
		"typeId": 1,
		"value": 100
	}}
