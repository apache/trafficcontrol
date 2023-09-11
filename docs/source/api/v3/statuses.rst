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

.. _to-api-v3-statuses:

************
``statuses``
************

``GET``
=======
Retrieves a list of all server :term:`Statuses`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                          |
	+=============+==========+======================================================================================================+
	| description | no       | Return only :term:`Statuses` with this *exact* description                                           |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| id          | no       | Return only the :term:`Status` with this integral, unique identifier                                 |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| name        | no       | Return only :term:`Statuses` with this name                                                          |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| orderby     | no       | Choose the ordering of the results - must be the name of one                                         |
	|             |          | of the fields of the objects in the ``response`` array                                               |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| sortOrder   | no       | Changes the order of sorting. Either ascending (default or "asc") or                                 |
	|             |          | descending ("desc")                                                                                  |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| limit       | no       | Choose the maximum number of results to return                                                       |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| offset      | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit |
	+-------------+----------+------------------------------------------------------------------------------------------------------+
	| page        | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are         |
	|             |          | ``limit`` long and the first page is 1. If ``offset`` was defined, this query parameter has no       |
	|             |          | effect. ``limit`` must be defined to make use of ``page``.                                           |
	+-------------+----------+------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/statuses?name=REPORTED HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:description: A short description of the status
:id:          The integral, unique identifier of this status
:lastUpdated: The date and time at which this status was last modified, in :ref:`non-rfc-datetime`
:name:        The name of the status

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: dHNip9kpTGGS1w39/fWcFehNktgmXZus8XaufnmDpv0PyG/3fK/KfoCO3ZOj9V74/CCffps7doEygWeL/xRtKA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 20:56:59 GMT
	Content-Length: 150

	{ "response": [
		{
			"description": "Server is online and reported in the health protocol.",
			"id": 3,
			"lastUpdated": "2018-12-10 19:11:17+00",
			"name": "REPORTED"
		}
	]}

``POST``
========
Creates a Server :term:`Status`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: STATUS:CREATE, STATUS:READ
:Response Type:  Array

Request Structure
-----------------
:description: Create a :term:`Status` with this description
:name:        Create a :term:`Status` with this name

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/statuses HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

	{ "description": "test", "name": "example" }

Response Structure
------------------
:description: A short description of the status
:id:          The integral, unique identifier of this status
:lastUpdated: The date and time at which this status was last modified, in :ref:`non-rfc-datetime`
:name:        The name of the status

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Wed, 21 Jun 2023 19:25:41 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 49FsWlQoEywA+MwYHFXcYmnLokUI4CWeDJLh8BGRB8V4ju9DckzvUUkFNGa7oXvDgEBpsxI4HoPuk8TCluvLTw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 21 Jun 2023 18:25:41 GMT
	Content-Length: 78

	{ "alerts": [
		{
			"text": "status was created.",
			"level": "success"
		}
	],"response": [
		{
			"description": "test",
			"id": 31,
			"lastUpdated": "2023-06-21 12:21:52-06",
			"name": "example"
		}
	]}
