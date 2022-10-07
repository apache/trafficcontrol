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

.. _to-api-v4-deliveryservices-id-capacity:

************************************
``deliveryservices/{{ID}}/capacity``
************************************

.. seealso:: :ref:`health-proto`

``GET``
=======
Retrieves the usage percentages of a servers associated with a :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------+
	| Name | Description                                                                  |
	+======+==============================================================================+
	| ID   | The integral, unique identifier for the :term:`Delivery Service` of interest |
	+------+------------------------------------------------------------------------------+

Response Structure
------------------
:availablePercent:   The percent of servers assigned to this :term:`Delivery Service` that is available - the allowed traffic level in terms of data per time period for all :term:`cache servers` that remains unused
:unavailablePercent: The percent of servers assigned to this :term:`Delivery Service` that is unavailable - the allowed traffic level in terms of data per time period for all :term:`cache servers` that can't be used because the servers are deemed unhealthy
:utilizedPercent:    The percent of servers assigned to this :term:`Delivery Service` that is currently in use - the allowed traffic level in terms of data per time period that is currently devoted to servicing requests
:maintenancePercent: The percent of servers assigned to this :term:`Delivery Service` that is unavailable due to server maintenance - the allowed traffic level in terms of data per time period that is unavailable because servers have intentionally been marked offline by administrators

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 15 Nov 2018 14:41:27 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: ++dFR9V1c60CHGNwMjX6JSFEjHreXcL4QnhTO3hiv04ByY379aLpL4OrOzX2bPgJgpR94+f6jZ0+iDIyTMwtFQ==
	Content-Length: 134

	{ "response": {
		"availablePercent": 99.9993696969697,
		"unavailablePercent": 0,
		"utilizedPercent": 0.00063030303030303,
		"maintenancePercent": 0
	}}

.. [#tenancy] Users will only be able to see capacity details for the :term:`Delivery Services` their :term:`Tenant` is allowed to see.
