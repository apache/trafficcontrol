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

.. _to-api-regions:

***********
``regions``
***********

``GET``
=======
Retrieves information about :term:`Regions`

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: REGION:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| division  | no       | Filter :term:`Regions` by the integral, unique identifier of the :term:`Division` which contains them         |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| id        | no       | Filter :term:`Regions` by integral, unique identifier                                                         |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| name      | no       | Filter :term:`Regions` by name                                                                                |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - either one of them "division", "id", "name"                              |
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

	GET /api/5.0/regions?division=1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
-------------------
:divisionName: The name of the division which contains this :term:Region
:division:     The integral, unique identifier of the division which contains this :term:Region
:id:           An integral, unique identifier for this :term:Region
:lastUpdated:  The date and time at which this :term:Region was last updated in :rfc:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name: The :term:Region name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: nSYbR+fRXaxhYl7dWgf0Udo2AsiXEnwvED1CPbk7ZNWK03I3TOhtmCQx9ABnJJ6xKYnlt6EKMeopVTK0nKU+SQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 06 Dec 2018 01:58:38 GMT
	Content-Length: 117

	{ "response": [
		{
			"divisionName": "Quebec",
			"division": 1,
			"id": 2,
			"lastUpdated": "2023-05-25T15:59:33.7096-06:00",
			"name": "Montreal"
		}
	]}

.. _to-api-regions-post:

``POST``
========
Creates a new :term:Region

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: REGION:CREATE, REGION:READ
:Response Type:  Object

Request Structure
-----------------
:division: The integral, unique identifier of the division which shall contain the new :term:Region
:name:     The name of the :term:Region

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/regions HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 65
	Content-Type: application/json

	{
		"division": 4,
		"name": "Manchester"
	}

Response Structure
------------------
:divisionName: The name of the division which contains this :term:Region
:division:     The integral, unique identifier of the division which contains this :term:Region
:id:           An integral, unique identifier for this :term:Region
:lastUpdated:  The date and time at which this :term:Region was last updated in :rfc:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name: The :term:Region name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: ezxk+iP7o7KE7zpWmGc0j8nz5k+1wAzY0HiNiA2xswTQrt+N+6CgQqUV2r9G1HAsPNr0HF2PhYs/Xr7DrYOw0A==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 06 Dec 2018 02:14:45 GMT
	Content-Length: 178

	{ "alerts": [
		{
			"text": "region was created.",
			"level": "success"
		}
	],
	"response": {
		"divisionName": "England",
		"division": 3,
		"id": 5,
		"lastUpdated": "2023-05-25T15:59:33.7096-06:00",
		"name": "Manchester"
	}}

``DELETE``
==========
Deletes a :term:`Region`\ . If no query parameter is specified, a ``400 Bad Request`` response is returned.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: REGION:DELETE
:Response Type:  Object

Request Structure
-----------------

.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------+
	| Name      | Required | Description                                          |
	+===========+==========+======================================================+
	| id        | no       | Delete :term:`Region` by integral, unique identifier |
	+-----------+----------+------------------------------------------------------+
	| name      | no       | Delete :term:`Region` by name                        |
	+-----------+----------+------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/regions?name=Manchester HTTP/1.1
	User-Agent: curl/7.29.0
	Host: trafficops.infra.ciab.test
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 07 Feb 2020 13:56:24 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: yNqXKcoiohEbJrEkH8LD1zifh87dIusuvUqgQnYueyKqCXkfd5bQvQ0OhQ2AAdSZa/oe2SAqMjojGsUlxHtIQw==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 07 Feb 2020 12:56:24 GMT
	Content-Length: 62

	{
		"alerts": [
			{
				"text": "region was deleted.",
				"level": "success"
			}
		]
	}
