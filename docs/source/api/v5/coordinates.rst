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

.. _to-api-coordinates:

***************
``coordinates``
***************

``GET``
=======
Gets a list of all coordinates in the Traffic Ops database

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: COORDINATE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| id        | no       | Return only coordinates that have this integral, unique identifier                                            |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| name      | no       | Return only coordinates with this name                                                                        |
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

Response Structure
------------------
:id:          Integral, unique, identifier for this coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in :rfc:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:latitude:  Latitude of the coordinate
:longitude: Longitude of the coordinate
:name:      The name of the coordinate - typically this just reflects the name of the :term:`Cache Group` for which the coordinate was created

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: Y2vxC3hpxIg6aRNBBT7i2hbAViIJp+dJoqHIzu3acFM+vGay/I5E+eZYOC9RY8hcJPrKNXysZOD8DOb9KsFgaw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 21:32:28 GMT
	Content-Length: 942

	{ "response": [
		{
			"id": 1,
			"name": "from_cachegroup_TRAFFIC_ANALYTICS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:04.596321Z"
		},
		{
			"id": 2,
			"name": "from_cachegroup_TRAFFIC_OPS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:04.596321Z"
		},
		{
			"id": 3,
			"name": "from_cachegroup_TRAFFIC_OPS_DB",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:04.596321Z"
		},
		{
			"id": 4,
			"name": "from_cachegroup_TRAFFIC_PORTAL",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:04.596321Z"
		},
		{
			"id": 5,
			"name": "from_cachegroup_TRAFFIC_STATS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:04.596321Z"
		},
		{
			"id": 6,
			"name": "from_cachegroup_CDN_in_a_Box_Mid",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:04.596321Z"
		},
		{
			"id": 7,
			"name": "from_cachegroup_CDN_in_a_Box_Edge",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24T16:07:05.596321Z"
		}
	]}

``POST``
========
Creates a new coordinate pair.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: COORDINATE:CREATE, COORDINATE:READ
:Response Type:  Object

Request Structure
-----------------
:name:      The name of the new coordinate
:latitude:  The desired latitude of the new coordinate (must be on the interval [-180, 180])
:longitude: The desired longitude of the new coordinate (must be on the interval [-90, 90])

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/coordinates HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 47
	Content-Type: application/json

	{"name": "test", "latitude": 0, "longitude": 0}

Response Structure
------------------
:id:          Integral, unique, identifier for the newly created coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in :RFC:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:latitude:    Latitude of the newly created coordinate
:longitude:   Longitude of the newly created coordinate
:name:        The name of the coordinate

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 7pWdeZyIIXE1P7o/JVon+5eSCbDw+FGamAzdXzWHXJ8IhF+Vh+/tWFCkzHYw3rP2kBVwZu+gqLffjQpBCMjt7A==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Nov 2018 17:48:55 GMT
	Content-Length: 165

	{ "alerts": [
		{
			"text": "Coordinate 'test' (#9) created",
			"level": "success"
		}
	],
	"response": {
		"id": 9,
		"name": "test",
		"latitude": 0,
		"longitude": 0,
		"lastUpdated": "2018-11-15T17:48:55.596321Z"
	}}


``PUT``
=======
Updates a coordinate

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: COORDINATE:UPDATE, COORDINATE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+------------------------------------------------------------+
	| Name | Required | Description                                                |
	+======+==========+============================================================+
	| id   | yes      | The integral, unique identifier of the coordinate to edit  |
	+------+----------+------------------------------------------------------------+

:name:      The name of the new coordinate
:latitude:  The desired new latitude of the coordinate (must be on the interval [-180, 180])
:longitude: The desired new longitude of the coordinate (must be on the interval [-90, 90])

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/coordinates?id=9 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 48
	Content-Type: application/json

	{"name": "quest", "latitude": 0, "longitude": 0}

Response Structure
------------------
:id:          Integral, unique, identifier for the coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in :RFC:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:latitude:  Latitude of the coordinate
:longitude: Longitude of the coordinate
:name:      The name of the coordinate

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: zd03Uvbnv8EbSZZ75Xp5tnnYStZsZTdyPxXnoqK4QZ5WhELLPL8iHlRfOaiLTbrUWUeJ8ue2HRz6aBS/iXCCGA==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Nov 2018 17:54:30 GMT
	Content-Length: 166

	{ "alerts": [
		{
			"text": "Coordinate 'quest' (#9) updated",
			"level": "success"
		}
	],
	"response": {
		"id": 9,
		"name": "quest",
		"latitude": 0,
		"longitude": 0,
		"lastUpdated": "2018-11-15T17:54:30.596321Z"
	}}

``DELETE``
==========
Deletes a coordinate

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: COORDINATE:DELETE, COORDINATE:READ
:Response Type:  Object

Request Structure
-----------------
:id:          Integral, unique, identifier for the coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in :RFC:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:latitude:  Latitude of the coordinate
:longitude: Longitude of the coordinate
:name:      The name of the coordinate

.. table:: Request Query Parameters

	+------+----------+-------------------------------------------------------------+
	| Name | Required | Description                                                 |
	+======+==========+=============================================================+
	| id   | yes      | The integral, unique identifier of the coordinate to delete |
	+------+----------+-------------------------------------------------------------+

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
	Whole-Content-Sha512: 82x/Wdckqgk4LN5LIlZfBJ26xkDrUVUGDjs5QFa/Lzap7dU3OZkjv8XW41xeFYj8PDmxHIpb7hiVObvLaxnEDA==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Nov 2018 17:57:42 GMT
	Content-Length: 65

	{ "alerts": [
		{
			"text": "Coordinate 'quest' (#9) deleted",
			"level": "success"
		}
	],
		"response": {
		"id": 9,
		"name": "quest",
		"latitude": 0,
		"longitude": 0,
		"lastUpdated": "2018-11-15T17:54:30.596321Z"
	}}
