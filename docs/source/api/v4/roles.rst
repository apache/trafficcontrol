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

.. _to-api-v4-roles:

*********
``roles``
*********

``GET``
=======
Retrieves all user :term:`Roles`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: ROLE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| id        | no       | Return only the :term:`Role` identified by this integral, unique identifier                                   |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| name      | no       | Return only the :term:`Role` with this name                                                                   |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|           |          | array                                                                                                         |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|           |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|           |          | defined to make use of ``page``.                                                                              |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/roles?name=read-only HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:permissions:  An array of the names of the Permissions given to this :term:`Role`
:description:  A description of the :term:`Role`
:name:         The name of the :term:`Role`
:lastUpdated: The date and time at which this :term:`Role` was last updated, in :rfc:`3339` format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: TEDXlQqWMSnJbL10JtFdbw0nqciNpjc4bd6m7iAB8aymakWeF+ghs1k5LayjdzHcjeDE8UNF/HXSxOFvoLFEuA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 25 Aug 2021 20:10:34 GMT
	Content-Length: 888

	{ "response": [
		{
			"name": "read-only",
			"description": "Has access to all read capabilities",
			"permissions": [
				"auth",
				"api-endpoints-read",
				"asns-read",
				"cache-config-files-read",
				"cache-groups-read",
				"capabilities-read",
				"cdns-read",
				"cdn-security-keys-read",
				"change-logs-read",
				"consistenthash-read",
				"coordinates-read",
				"delivery-services-read",
				"delivery-service-security-keys-read",
				"delivery-service-requests-read",
				"delivery-service-servers-read",
				"divisions-read",
				"to-extensions-read",
				"federations-read",
				"hwinfo-read",
				"jobs-read",
				"origins-read",
				"parameters-read",
				"phys-locations-read",
				"profiles-read",
				"regions-read",
				"roles-read",
				"server-capabilities-read",
				"servers-read",
				"service-categories-read",
				"stats-read",
				"statuses-read",
				"static-dns-entries-read",
				"steering-read",
				"steering-targets-read",
				"system-info-read",
				"tenants-read",
				"types-read",
				"users-read"
			],
			"lastUpdated": "2021-05-03T14:50:18.93513-06:00",
		}
	]}

``POST``
========
Creates a new :term:`Role`.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: ROLE:CREATE, ROLE:READ
:Response Type: Object

Request Structure
-----------------
:permissions:  An optional array of permission names that will be granted to the new :term:`Role`\ [#permissions]_
:description:  A helpful description of the :term:`Role`'s purpose.
:name:         The name of the new :term:`Role`

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/roles HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 56
	Content-Type: application/json

	{
		"name": "test",
		"description": "quest"
	}


Response Structure
------------------
:permissions: An array of the names of the Permissions given to this :term:`Role`

	.. tip:: This can be ``null`` *or* empty, depending on whether it was present in the request body, or merely empty. Obviously, it can also be a populated array.

:description: A description of the :term:`Role`
:name:        The name of the :term:`Role`
:lastUpdated: The date and time at which this :term:`Role` was last updated, in :rfc:`3339` format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: gzfc7m/in5vVsVP+Y9h6JJfDhgpXKn9VAzoiPENhKbQfP8Q6jug08Rt2AK/3Nz1cx5zZ8P9IjVxDdIg7mlC8bw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 04 Sep 2019 17:44:42 GMT
	Content-Length: 128

	{ "alerts": [{
		"text": "role was created.",
		"level": "success"
	}],
	"response": {
		"name": "test",
		"description": "quest",
		"permissions": null,
		"lastUpdated": "2021-05-03T14:50:18.93513-06:00"
	}}

``PUT``
=======
Replaces an existing :term:`Role` with one provided by the request\ [#admin]_.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: ROLE:UPDATE, ROLE:READ
:Response Type:

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+--------------------------------------------------------------------+
	| Name | Required | Description                                                        |
	+======+==========+====================================================================+
	| name | yes      | The name of the :term:`Role` to be updated                         |
	+------+----------+--------------------------------------------------------------------+

:permissions: An optional array of permission names that will be granted to the new :term:`Role`

	.. warning:: When not present, the affected :term:`Role`'s Permissions will be unchanged - *not* removed, unlike when the array is empty.

:description: A helpful description of the :term:`Role`'s purpose.
:name:        The new name of the :term:`Role`

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/roles?name=test HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 56
	Content-Type: application/json

	{
		"name":"test",
		"description": "quest_updated"
	}

Response Structure
------------------
:permissions: An array of the names of the Permissions given to this :term:`Role`

	.. tip:: This can be ``null`` *or* empty, depending on whether it was present in the request body, or merely empty. Obviously, it can also be a populated array.

	.. warning:: If no ``permissions`` array was given in the request, this will *always* be ``null``, even if the :term:`Role` has Permissions that would have gone unchanged.

:description: A description of the :term:`Role`
:name:        The name of the :term:`Role`
:lastUpdated: The date and time at which this :term:`Role` was last updated, in :rfc:`3339` format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: mlHQenE1Q3gjrIK2lC2hfueQOaTCpdYEfboN0A9vYPUIwTiaF5ZaAMPQBdfGyiAhgHRxowITs3bR7s1L++oFTQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 05 Sep 2019 12:56:46 GMT
	Content-Length: 136

	{
		"alerts": [
			{
				"text": "role was updated.",
				"level": "success"
			}
		],
		"response": {
			"name": "test",
			"description": "quest_updated",
			"permissions": null,
			"lastUpdated": "2021-05-03T14:50:18.93513-06:00"
		}
	}


``DELETE``
==========
Deletes a :term:`Role`\ [#admin]_.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: ROLE:DELETE, ROLE:READ
:Response Type: ``undefined``

Request Structure
-----------------
.. table:: Request  Query Parameters

	+------+----------+--------------------------------------------------------------------+
	| Name | Required | Description                                                        |
	+======+==========+====================================================================+
	| name | yes      | The name of the :term:`Role` to be deleted                         |
	+------+----------+--------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/roles?name=test HTTP/1.1
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
	Whole-Content-Sha512: 10jeFZihtbvAus/XyHAW8rhgS9JBD+X/ezCp1iExYkEcHxN4gjr1L6x8zDFXORueBSlFldgtbWKT7QsmwCHUWA==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 05 Sep 2019 13:02:06 GMT
	Content-Length: 60

	{ "alerts": [{
		"text": "role was deleted.",
		"level": "success"
	}]}

.. [#permissions] ``permissions`` cannot include permissions that are not included in the permissions of the requesting user. In POST requests, if ``permissions`` is omitted or explicitly ``null``, it is treated as an empty set/array.
.. [#admin] The special :term:`Role` with the name "admin" cannot be modified or deleted - regardless of user Permissions.
