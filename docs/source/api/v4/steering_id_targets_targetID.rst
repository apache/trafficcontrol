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

.. _to-api-v4-steering-id-targets-targetID:

****************************************
``steering/{{ID}}/targets/{{targetID}}``
****************************************

``PUT``
=======
Updates a steering target.

:Auth. Required: Yes
:Roles Required: Portal, Steering, Federation, "operations" or "admin"
:Permissions Required: STEERING:UPDATE, STEERING:READ, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	|   Name   |                Description                                                                                                           |
	+==========+======================================================================================================================================+
	|    ID    | The integral, unique identifier of a steering :term:`Delivery Service`                                                               |
	+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| targetID | The integral, unique identifier of a :term:`Delivery Service` which is a target of the :term:`Delivery Service` identified by ``ID`` |
	+----------+--------------------------------------------------------------------------------------------------------------------------------------+

:typeId: The integral, unique identifier of the :ref:`steering type <ds-steering>` of the target :term:`Delivery Service`. This should be corresponding to one of ``STEERING_WEIGHT``, ``STEERING_ORDER``, ``STEERING_GEO_ORDER`` or ``STEERING_GEO_WEIGHT``
:value:  The 'weight', 'order', 'geo_order' or 'geo_weight' which shall be attributed to the target :term:`Delivery Service`

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/steering/2/targets/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 26
	Content-Type: application/json

	{
		"value": 1,
		"typeId": 43
	}

Response Structure
------------------
:deliveryService:   A string that is the :ref:`ds-xmlid` of the steering :term:`Delivery Service`
:deliveryServiceId: An integral, unique identifier for the steering :term:`Delivery Service`
:target:            A string that is the :ref:`ds-xmlid` of this target :term:`Delivery Service`
:targetId:          An integral, unique identifier for this target :term:`Delivery Service`
:type:              The steering type of this target :term:`Delivery Service`
:typeId:            An integral, unique identifier for the :ref:`steering type <ds-steering>` of this target :term:`Delivery Service`
:value:             The 'weight' attributed to this steering target as an integer

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: AfXsIRzdtU3HZYkr93qBMVTZRJ5oTF2u5sKYnd+DSqxZ+RQxY6vXtCupnnXCf9dxMt5QXRW1EFOW/FBG6lFrTg==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 14:34:22 GMT
	Content-Length: 194

	{ "alerts": [
		{
			"text": "steeringtarget was updated.",
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
		"value": 1
	}}

``DELETE``
==========
Removes a specific target mapping from a specific :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: Portal, Steering, Federation, "operations" or "admin"
:Permissions Required: STEERING:DELETE, STEERING:READ, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+----------+----------------------------------------------------------------------------------------------------------------------------------------------------+
	|   Name   |                Description                                                                                                                         |
	+==========+====================================================================================================================================================+
	|    ID    | The integral, unique identifier of a steering :term:`Delivery Service` - a target of which shall be deleted                                        |
	+----------+----------------------------------------------------------------------------------------------------------------------------------------------------+
	| targetID | The integral, unique identifier of a :term:`Delivery Service` which is a target to be removed of the :term:`Delivery Service` identified by ``ID`` |
	+----------+----------------------------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/steering/2/targets/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: N6h8Kl7uQveqpTc3fmKXFDY2yYe5smApNcaTow4ab0DHGFdJfqQh89I4nvvaXvmVNhxVAqX3UE/6blbO8/9Xqg==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 14:42:54 GMT
	Content-Length: 69

	{ "alerts": [
		{
			"text": "steeringtarget was deleted.",
			"level": "success"
		}
	]}
