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

.. _to-api-tenants-id:

******************
``tenants/{{ID}}``
******************

``PUT``
=======
Updates a specific tenant.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: TENANT:UPDATE, TENANT:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------+
	| Name |                 Description                                   |
	+======+===============================================================+
	|  ID  | The integral, unique identifier for the tenant being modified |
	+------+---------------------------------------------------------------+

:active:   An optional boolean - default: ``false`` - which indicates whether or not the tenant shall be immediately active
:name:     The name of the tenant
:parentId: The integral, unique identifier of the parent of this tenant

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/tenants/9 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 59
	Content-Type: application/json

	{
		"active": true,
		"name": "quest",
		"parentId": 3
	}

Response Structure
------------------
:active:      A boolean which indicates whether or not the tenant is active
:id:          The integral, unique identifier of this tenant
:lastUpdated: The date and time at which the :term:Tenant was last updated, in :RFC:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name:     This tenant's name
:parentId: The integral, unique identifier of this tenant's parent

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 5soYQFrG2x5ZJ1e5UZIOLUv/928qyd2Bfgw21Wv85rqjLpyeT3djkfRVD1/xpKConulNrZs2czJKrrwZA7X61w==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 20:30:54 GMT
	Content-Length: 163

	{ "alerts": [
		{
			"text": "tenant was updated.",
			"level": "success"
		}
	],
	"response": {
		"id": 9,
		"name": "quest",
		"active": true,
		"lastUpdated": "2023-05-30T19:52:58.183642+00:00",
		"parentId": 3
	}}

``DELETE``
==========
Deletes a specific tenant.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: TENANT:DELETE, TENANT:READ
:Response Type:  ``undefined``


Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------+
	| Name |                 Description                                  |
	+======+==============================================================+
	|  ID  | The integral, unique identifier for the tenant being deleted |
	+------+--------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/tenants/9 HTTP/1.1
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
	Whole-Content-Sha512: KU0XIbFoD0Cy06kzH2Gl59pBqie/TEFJgh33mssGNwXJZlRkTLaSTHT8Df4X+pOs7UauZH10akGvaA0UTiN/vg==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 20:40:31 GMT
	Content-Length: 61

	{ "alerts": [
		{
			"text": "tenant was deleted.",
			"level": "success"
		}
	]}
