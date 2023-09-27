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

.. _to-api-tenants:

***********
``tenants``
***********

``GET``
=======
Get all requested :term:`Tenants`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: TENANT:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+------------------------------------------------------------------------------------+
	| Name      | Description                                                                        |
	+===========+====================================================================================+
	| active    | If ``true``, return only active :term:`Tenants`; if ``false`` return only inactive |
	|           | :term:`Tenants`                                                                    |
	+-----------+------------------------------------------------------------------------------------+
	| id        | Return only :term:`Tenants` with this integral, unique identifier                  |
	+-----------+------------------------------------------------------------------------------------+
	| name      | Return only :term:`Tenants` with this name                                         |
	+-----------+------------------------------------------------------------------------------------+
	| orderby   | Choose the ordering of the results - must be the name of one of the fields of the  |
	|           | objects in the ``response`` array                                                  |
	+-----------+------------------------------------------------------------------------------------+
	| sortOrder | Changes the order of sorting. Either ascending (default or "asc") or descending    |
	|           | ("desc")                                                                           |
	+-----------+------------------------------------------------------------------------------------+
	| limit     | Choose the maximum number of results to return                                     |
	+-----------+------------------------------------------------------------------------------------+
	| offset    | The number of results to skip before beginning to return results. Must use in      |
	|           | conjunction with limit                                                             |
	+-----------+------------------------------------------------------------------------------------+
	| page      | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, |
	|           | pages are ``limit`` long and the first page is 1. If ``offset`` was defined, this  |
	|           | query parameter has no effect. ``limit`` must be defined to make use of ``page``.  |
	+-----------+------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/tenants?name=root HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:active:      A boolean which indicates whether or not the :term:`Tenant` is active
:id:          The integral, unique identifier of this :term:`Tenant`
:lastUpdated: The date and time at which the :term:Tenant was last updated, in :RFC:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name:       This :term:`Tenant`'s name
:parentId:   The integral, unique identifier of this :term:`Tenant`'s parent
:parentName: The name of the parent of this :term:`Tenant`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: Yzr6TfhxgpZ3pbbrr4TRG4wC3PlnHDDzgs2igtz/1ppLSy2MzugqaGW4y5yzwzl5T3+7q6HWej7GQZt1XIVeZQ==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 19:57:58 GMT
	Content-Length: 106

	{ "response": [
		{
			"id": 1,
			"name": "root",
			"active": true,
			"lastUpdated": "2023-05-30T19:52:58.183642+00:00",
			"parentId": null
		}
	]}

``POST``
========
Create a new tenant.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: TENANT:CREATE, TENANT:READ
:Response Type:  Object

Request Structure
-----------------
:active:   An optional boolean - default: ``false`` - which indicates whether or not the tenant shall be immediately active
:name:     The name of the tenant
:parentId: The integral, unique identifier of the parent of this tenant

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/tenants HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 48
	Content-Type: application/json

	{
		"active": true,
		"name": "test",
		"parentId": 1
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
	Whole-Content-Sha512: ysdopC//JQI79BRUa61s6M2HzHxYHpo5RdcuauOoqCYxiVOoUhNZfOVydVkv8zDN2qA374XKnym4kWj3VzQIXg==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 11 Dec 2018 19:37:16 GMT
	Content-Length: 162

	{ "alerts": [
		{
			"text": "tenant was created.",
			"level": "success"
		}
	],
	"response": {
		"id": 9,
		"name": "test",
		"active": true,
		"lastUpdated": "2023-05-30T19:52:58.183642+00:00",
		"parentId": 1
	}}
