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

.. _to-api-v3-federations-id-users:

****************************
``federations/{{ID}}/users``
****************************

``GET``
=======
Retrieves users assigned to a federation.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------+
	| Name |                 Description                                                         |
	+======+=====================================================================================+
	|  ID  | The integral, unique identifier of the federation for which users will be retrieved |
	+------+-------------------------------------------------------------------------------------+

.. table:: Request Query Parameters

	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                                          |
	+===========+==========+======================================================================================================================================+
	| userID    | no       | Show only the user that has this integral, unique identifier                                                                         |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| role      | no       | Show only the users that have this role                                                                                              |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response``                        |
	|           |          | array                                                                                                                                |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                             |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                                       |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                                 |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first page is 1. |
	|           |          | If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``.                    |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:company:  The company to which the user belongs
:email:    The user's email address
:fullName: The user's full name
:id:       An integral, unique identifier for the user
:role:     The user's highest role
:username: The user's short "username"

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:31:34 GMT
	X-Server-Name: traffic_ops_golang/
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:31:34 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: eQQoF2xlbK2I2oTja7zrt/FlkLzCgwpU2zb2+rmIjHbHJ3MnmsSczSamIAAyTzs5gDaqcuUX1G35ZB8d7Bj82g==
	content-length: 101

	{ "response": [
		{
			"fullName": null,
			"email": null,
			"id": 2,
			"role": "admin",
			"company": null,
			"username": "admin"
		}
	]}


``POST``
========
Assigns one or more users to a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
:userIds: An array of integral, unique identifiers for users which will be assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, will cause any conflicting assignments already in place to be overridden by this request

	.. note:: If ``replace`` is not given (and/or not ``true``), then any conflicts with existing assignments will cause the entire operation to fail.

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/federations/1/users HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 34
	Content-Type: application/json

	{
		"userIds": [2],
		"replace": true
	}

Response Structure
------------------
:userIds: An array of integral, unique identifiers for users which have been assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, caused any conflicting assignments already in place to be overridden by this request

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:29:19 GMT
	X-Server-Name: traffic_ops_golang/
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:29:19 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: MvPmgOAs58aSOGvh+iEilflgOexbaexg+qE2IPrQZX0H4iSX4JvEys9adbGE9a9yaLj9uUMxg77N6ZyDhVqsbQ==
	content-length: 137

	{ "alerts": [
		{
			"level": "success",
			"text": "1 user(s) were assigned to the test.quest. federation"
		}
	],
	"response": {
		"userIds": [
			2
		],
		"replace": true
	}}
