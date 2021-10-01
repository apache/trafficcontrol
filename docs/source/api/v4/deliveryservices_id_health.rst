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

.. _to-api-deliveryservices-id-health:

**********************************
``deliveryservices/{{ID}}/health``
**********************************

.. seealso:: :ref:`health-proto`

``GET``
=======
Retrieves the health of all :term:`Cache Groups` assigned to a particular :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:READ, CACHE-GROUP:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                                |
	+======+============================================================================================================+
	| ID   | The integral, unique identifier of the Delivery service for which :term:`Cache Groups` will be displayed   |
	+------+------------------------------------------------------------------------------------------------------------+


Response Structure
------------------
:cachegroups: An array of objects that represent the health of each :term:`Cache Group` assigned to this :term:`Delivery Service`

	:name:    A string that is the :ref:`name of the Cache Group <cache-group-name>` represented by this object
	:offline: The number of OFFLINE :term:`cache servers` within this :term:`Cache Group`
	:online:  The number of ONLINE :term:`cache servers` within this :term:`Cache Group`

:totalOffline: Total number of OFFLINE :term:`cache servers` assigned to this :term:`Delivery Service`
:totalOnline:  Total number of ONLINE :term:`cache servers` assigned to this :term:`Delivery Service`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 15 Nov 2018 14:43:43 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: KpXViXeAgch58ueQqdyU8NuINBw1EUedE6Rv2ewcLUajJp6kowdbVynpwW7XiSvAyHdtClIOuT3OkhIimghzSA==
	Content-Length: 115

	{ "response": {
		"totalOffline": 0,
		"totalOnline": 1,
		"cachegroups": [
			{
				"offline": 0,
				"name": "CDN_in_a_Box_Edge",
				"online": 1
			}
		]
	}}

.. [#tenancy] Users will only be able to see :term:`Cache Group` health details for the :term:`Delivery Services` their :term:`Tenant` is allowed to see.
