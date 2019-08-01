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

.. _to-api-cachegroups:

***************
``cachegroups``
***************

``GET``
=======
Extract information about :term:`Cache Groups`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| type      | no       | Return only :term:`Cache Groups` that are of the :term:`type` identified by this integral, unique identifier  |
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

	GET /api/1.3/cachegroups?type=23 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...


Response Structure
------------------
:fallbacks: An array of :term:`Cache Group` names that are registered as "fallbacks" for use when this :term:`Cache Group` is unavailable.\ [#fallbacks]_

	.. versionadded:: ATCv4.0

		This field was added to all versions of this endpoint with Apache Traffic Control version 4.0

:fallbackToClosest:             If ``true``, Traffic Router will direct clients to peers of this :term:`Cache Group` in the event that it becomes unavailable.\ [#fallbacks]_
:id:                            A numeric, unique identifier for the :term:`Cache Group`
:lastUpdated:                   The time and date at which this entry was last updated in ISO format
:latitude:                      Latitude for the :term:`Cache Group`
:longitude:                     Longitude for the :term:`Cache Group`
:name:                          The name of the :term:`Cache Group` entry
:parentCachegroupId:            ID of this :term:`Cache Group`'s parent :term:`Cache Group` (if any)
:parentCachegroupName:          Name of this :term:`Cache Group`'s parent :term:`Cache Group` (if any)
:secondaryParentCachegroupId:   ID of this :term:`Cache Group`'s secondary parent :term:`Cache Group` (if any)
:secondaryParentCachegroupName: Name of this :term:`Cache Group`'s secondary parent :term:`Cache Group` (if any)
:shortName:                     Abbreviation of the :term:`Cache Group` name
:typeId:                        Unique identifier for the ':term:`Type`' of :term:`Cache Group` entry
:typeName:                      The name of the :term:`type` of :term:`Cache Group` entry

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: oV6ifEgoFy+v049tVjSsRdWQf4bxjrUvIYfDdgpUtlxiC7gzCv31m5bXQ8EUBW4eg2hfYM+BsGvJpnNDZB7pUg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 07 Nov 2018 19:46:36 GMT
	Content-Length: 379

	{ "response": [
		{
			"id": 7,
			"name": "CDN_in_a_Box_Edge",
			"shortName": "ciabEdge",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": "CDN_in_a_Box_Mid",
			"parentCachegroupId": 6,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": [],
			"localizationMethods": [],
			"typeName": "EDGE_LOC",
			"typeId": 23,
			"lastUpdated": "2018-11-07 14:45:43+00",
			"fallbacks": []
		}
	]}


``POST``
========
Creates a :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:fallbacks: An optional field which, when present, should contain an array of names of other :term:`Cache Groups` on which the Traffic Router will fall back in the event that this :term:`Cache Group` fails/becomes unavailable\ [#fallbacks]_

	.. versionadded:: ATCv4.0

		Support for this field was added to all versions of this endpoint with Apache Traffic Control version 4.0

:fallbackToClosest: If ``true``, the Traffic Router will fall back on the 'closest' :term:`Cache Group` to this one, when this one fails\ [#fallbacks]_

	.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

:latitude:                    An optional field which, if present, will define the latitude for the :term:`Cache Group` to ISO-standard double specification\ [#optional]_
:longitude:                   An optional field which, if present, will define the longitude for the :term:`Cache Group` to ISO-standard double specification\ [#optional]_
:localizationMethods:         Array of enabled localization methods (as strings)
:fallbacks:                   Array of fallback server hostnames.
:name:                        The name of the :term:`Cache Group`
:parentCachegroupId:          An optional field which, if present, should be an integral, unique identifier for this :term:`Cache Group`'s primary parent
:secondaryParentCachegroupId: An optional field which, if present, should be an integral, unique identifier for this :term:`Cache Group`'s secondary parent
:shortName:                   An abbreviation of the ``name``
:typeId:                      An integral, unique identifier for the :term:`type` of :term:`Cache Group`; one of:

	EDGE_LOC
		Indicates a group of Edge-tier caches
	MID_LOC
		Indicates a group of Mid-tier caches
	ORG_LOC
		Indicates a group of origin servers (though only one server will typically be in any given ORG_LOC)

	.. note:: The actual, integral, unique identifiers for these types must first be obtained, generally via :ref:`to-api-types`.

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/cachegroups HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 252
	Content-Type: application/x-www-form-urlencoded

	{
		"fallbackToClosest": false,
		"latitude": 0,
		"longitude": 0,
		"localizationMethods": [],
		"fallbacks": [],
		"name": "test",
		"parentCachegroupId": 7,
		"shortName": "test",
		"typeId": 23
	}

Response Structure
------------------
:fallbacks: An array of :term:`Cache Group` names that are registered as "fallbacks" for use when this :term:`Cache Group` is unavailable\ [#fallbacks]_

	.. versionadded:: ATCv4.0

		This field was added to all versions of this endpoint with Apache Traffic Control version 4.0

:fallbackToClosest:             If ``true``, Traffic Router will direct clients to peers of this :term:`Cache Group` in the event that it becomes unavailable\ [#fallbacks]_
:id:                            A numeric, unique identifier for the :term:`Cache Group`
:lastUpdated:                   The time and date at which this entry was last updated in ISO format
:latitude:                      Latitude for the :term:`Cache Group`
:longitude:                     Longitude for the :term:`Cache Group`
:localizationMethods:           Array of enabled localization methods (as strings)
:fallbacks:                     Array of fallback server hostnames
:name:                          The name of the :term:`Cache Group` entry
:parentCachegroupId:            ID of this :term:`Cache Group`'s parent :term:`Cache Group` (if any)
:parentCachegroupName:          Name of this :term:`Cache Group`'s parent :term:`Cache Group` (if any)
:secondaryParentCachegroupId:   ID of this :term:`Cache Group`'s secondary parent :term:`Cache Group` (if any)
:secondaryParentCachegroupName: Name of this :term:`Cache Group`'s secondary parent :term:`Cache Group` (if any)
:shortName:                     Abbreviation of the :term:`Cache Group` name
:typeId:                        Unique identifier for the ':term:`Type`' of :term:`Cache Group` entry
:typeName:                      The name of the :term:`type` of :term:`Cache Group` entry


.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: YvZlh3rpfl3nBq6SbNVhbkt3IvckbB9amqGW2JhLxWK9K3cxjBq5J2sIHBUhrLKUhE9afpxtvaYrLRxjt1/YMQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 07 Nov 2018 22:11:50 GMT
	Content-Length: 379

	{ "alerts": [
		{
			"text": "cachegroup was created.",
			"level": "success"
		}
	],
	"response": {
		"id": 10,
		"name": "test",
		"shortName": "test",
		"latitude": 0,
		"longitude": 0,
		"parentCachegroupName": "CDN_in_a_Box_Mid",
		"parentCachegroupId": 7,
		"secondaryParentCachegroupName": null,
		"secondaryParentCachegroupId": null,
		"fallbackToClosest": false,
		"localizationMethods": [],
		"fallbacks": [],
		"typeName": "EDGE_LOC",
		"typeId": 23,
		"lastUpdated": "2018-11-07 22:11:50+00"
	}}

.. [#fallbacks] Traffic Router will first check for a ``fallbacks`` array and, when that is empty/unset/all the :term:`Cache Groups` in it are also unavailable, will subsequently check for ``fallbackToClosest``. If that is ``true``, then it falls back to the geographically closest :term:`Cache Group` capable of serving the same content or, when it is ``false``/no such :term:`Cache Group` exists/said :term:`Cache Group` is also unavailable, will respond to clients with a failure response indicating the problem.
.. [#optional] While these fields are technically optional, note that if they are not specified many things may break. For this reason, Traffic Portal requires them when creating or editing :term:`Cache Groups`.

.. This doesn't appear to exist anymore - can't reproduce in CIAB nor production
.. ``/api/1.1/cachegroups/:parameter_id/parameter/available``
.. ==========================================================
