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

.. _to-api-v1-deliveryservice_user:

************************
``deliveryservice_user``
************************
.. deprecated:: ATCv4

``POST``
========
Assigns one or more :term:`Delivery Services` to a user.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:userId:           An integral, unique identifier for the user to whom the :term:`Delivery Service(s) <Delivery Service>` identified in ``deliveryServices`` will be assigned
:deliveryServices: An array of integral, unique identifiers for the :term:`Delivery Service(s) <Delivery Service>` being assigned to the user identified by ``userId``
:replace:          An optional field which, when present and ``true`` will replace existing user/ds assignments? (true|false)

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/deliveryservice_user HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 38
	Content-Type: application/json

	{"userId": 5, "deliveryServices": [1]}

Response Structure
------------------
:userId:           The integral, unique identifier of the user to whom the :term:`Delivery Service(s) <Delivery Service>` identified in ``deliveryServices`` are assigned
:deliveryServices: An array of integral, unique identifiers of :term:`Delivery Services` assigned to the user identified by ``userId``
:replace:          If ``true``, any and all existing, conflicting :term:`Delivery Service` assignments were overwritten by this assignment operation

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 21:37:30 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Uwl+924m6Ye3NraFP+RBpldkhcNTTDyXHZbzRaYV95p9tP56Z61gckeKSr1oQIkNXjXcCsDN5Dmum7Zk1AR6Hw==
	Content-Length: 127

	{ "alerts": [
		{
			"level": "success",
			"text": "Delivery service assignments complete."
		},
		{
			"level": "warning",
			"text": "This endpoint and its functionality is deprecated, and will be removed in the future"
		}
	],
	"response": {
		"userId": 5,
		"deliveryServices": [
			1
		]
	}}
