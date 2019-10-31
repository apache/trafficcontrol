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

******************************************
``user/{{ID}}/deliveryservices/available``
******************************************

``GET``
=======
Lists identifying information for all of the :term:`Delivery Services` assigned to a user - **not**, as the name implies, the :term:`Delivery Services` *available* to be assigned to that user.

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                       |
	+======+===================================================================================================+
	|  ID  | The integral, unique identifier of the users whose :term:`Delivery Services` shall be retrieved   |
	+------+---------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/user/2/deliveryservices/available HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:displayName: This :term:`Delivery Service`'s :ref:`ds-display-name`
:id:          The integral, unique identifier of this :term:`Delivery Service`
:xmlId:       The :ref:`ds-xmlid` which (also) uniquely identifies this :term:`Delivery Service`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: A1IUM2qkvJkviD0mcEADoCMiy76AWRO/Xnc70ur3CrOlkySXwqxfrhLc3wKlI1926yW+QrTd3nQaVpbX7Rd9wQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 13 Dec 2018 22:31:44 GMT
	Content-Length: 62

	{ "response": [
		{
			"id": 1,
			"displayName": "Demo 1",
			"xmlId": "demo1"
		}
	]}

.. [#tenancy] Only the :term:`Delivery Services` visible to the requesting user's :term:`Tenant` will appear, regardless of :term:`Role` or actual 'assignment' status.
