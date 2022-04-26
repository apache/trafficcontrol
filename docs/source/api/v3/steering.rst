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

.. _to-api-v3-steering:

************
``steering``
************

``GET``
=======
Gets a list of all :ref:`Steering Targets <steering-qht>` in the Traffic Ops database.

:Auth. Required: Yes
:Roles Required: "Portal", "Steering", "Federation", "operations" or "admin"
:Response Type:  Array

Request Structure
-----------------
No parameters available.

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/steering HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:deliveryService:               A string that is the :ref:`ds-xmlid` of the steering :term:`Delivery Service`
:clientSteering:                Whether this is a :ref:`client steering <ds-client-steering>` Delivery Service.
:targets:                       The delivery services that the :ref:`Steering Delivery Service <tr-steering>` targets.

	:order:                 If this is a :ref:`STEERING_ORDER <ds-steering-order>` target, this is the value of the order. Otherwise, ``0``.
	:weight:                If this is a :ref:`STEERING_WEIGHT <ds-steering-weight>` target, this is the value of the weight. Otherwise, ``0``.
	:deliveryService:       A string that is the :ref:`ds-xmlid` of the steering :term:`Delivery Service`

:filters:                       Filters of type :ref:`STEERING_REGEXP <ds-steering-regexp>` that exist on either of the targets.

	:deliveryService:       A string that is the :ref:`ds-xmlid` of the steering :term:`Delivery Service`
	:pattern:               A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 18:56:57 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: hcJa4xVLDx7bxBmoSjYo5YUwdSBWQr9GlqRYrc6ZU7LeenjiV3go22YlIHt/GtjLcHQjJ5DulKRhdsvFMq7Fng==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 17:56:57 GMT
	Content-Length: 167

	{
		"response": [
			{
				"deliveryService": "steering1",
				"clientSteering": true,
				"targets": [
					{
						"order": 0,
						"weight": 1,
						"deliveryService": "demo1"
					},
					{
						"order": 0,
						"weight": 2,
						"deliveryService": "demo2"
					}
				],
				"filters": [
					{
						"deliveryService": "demo1",
						"pattern": ".*\\.demo1\\..*"
					},
					{
						"deliveryService": "demo2",
						"pattern": ".*\\.demo2*\\..*"
					}
				]
			}
		]
	}
