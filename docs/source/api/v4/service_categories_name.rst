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

.. _to-api-v4-service-categories-name:

*******************************
``service_categories/{{name}}``
*******************************

``PUT``
========
Update a service category.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVICE-CATEGORY:UPDATE, SERVICE-CATEGORY:READ
:Response Type:  Object

Request Structure
-----------------
:name:        The :term:`Service Category`'s new name

.. table:: Request Path Parameters

	+------------+------------------------------------------------------------------------+
	| Name       | Description                                                            |
	+============+========================================================================+
	| name       | The current name of the :term:`Service Category`                       |
	+------------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/service_categories/sc-name HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 48
	Content-Type: application/json

	{
		"name": "New Name",
	}

Response Structure
------------------
:name:        This :term:`Service Category`'s name
:lastUpdated: The date and time at which this :term:`Service Category` was last modified, in :ref:`non-rfc-datetime`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: +pJm4c3O+JTaSXNt+LP+u240Ba/SsvSSDOQ4rDc6hcyZ0FIL+iY/WWrMHhpLulRGKGY88bM4YPCMaxGn3FZ9yQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 11 Mar 2020 20:12:20 GMT
	Content-Length: 189

	{
		"alerts": [
			{
				"text": "Service Category was updated.",
				"level": "success"
			}
		],
		"response": {
			"lastUpdated": "2020-03-11 14:12:20-06",
			"name": "New Name"
		}
	}

``DELETE``
==========
Deletes a specific :term:`Service Category`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVICE-CATEGORY:DELETE, SERVICE-CATEGORY:READ
:Response Type:  ``undefined``


Request Structure
-----------------
.. table:: Request Path Parameters

	+------------+------------------------------------------------------------------------+
	| Name       | Description                                                            |
	+============+========================================================================+
	| name       | The current name of the :term:`Service Category` to be deleted         |
	+------------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/service_categories/my-service-category HTTP/1.1
	User-Agent: python-requests/2.23.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0


Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 17 Aug 2020 16:13:31 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: yErJobzG9IA0khvqZQK+Yi7X4pFVvOqxn6PjrdzN5DnKVm/K8Kka3REul1XmKJnMXVRY8RayoEVGDm16mBFe4Q==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 17 Aug 2020 15:13:31 GMT
	Content-Length: 103

	{
		"alerts": [
			{
				"text": "my-service-category was deleted.",
				"level": "success"
			}
		]
	}
