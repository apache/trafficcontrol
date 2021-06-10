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

.. _to-api-v3-cachegroups:

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

	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                              |
	+===========+==========+==========================================================================================================================+
	| id        | no       | Return the only :term:`Cache Group` that has this id                                                                     |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| name      | no       | Return only the :term:`Cache Group` identified by this :ref:`cache-group-name`                                           |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| type      | no       | Return only :term:`Cache Groups` that are of the :ref:`cache-group-type` identified by this integral, unique identifier  |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| topology  | no       | Return only :term:`Cache Groups` that are used in the :term:`Topology` identified by this unique identifier              |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` array      |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                 |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                           |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                     |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long  and the     |
	|           |          | first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of |
	|           |          | ``page``.                                                                                                                |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/cachegroups?type=23 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...


Response Structure
------------------
:fallbacks:                     An array of strings that are :ref:`Cache Group names <cache-group-name>` that are registered as :ref:`cache-group-fallbacks` for this :term:`Cache Group`\ [#fallbacks]_
:fallbackToClosest:             A boolean value that defines the :ref:`cache-group-fallback-to-closest` behavior of this :term:`Cache Group`\ [#fallbacks]_
:id:                            An integer that is the :ref:`cache-group-id` of the :term:`Cache Group`
:lastUpdated:                   The time and date at which this entry was last updated in :ref:`non-rfc-datetime`
:latitude:                      A floating-point :ref:`cache-group-latitude` for the :term:`Cache Group`
:localizationMethods:           An array of :ref:`cache-group-localization-methods` as strings
:longitude:                     A floating-point :ref:`cache-group-longitude` for the :term:`Cache Group`
:name:                          A string containing the :ref:`cache-group-name` of the :term:`Cache Group`
:parentCachegroupId:            An integer that is the :ref:`cache-group-id` of this :term:`Cache Group`'s :ref:`cache-group-parent` - or ``null`` if it doesn't have a :ref:`cache-group-parent`
:parentCachegroupName:          A string containing the :ref:`cache-group-name` of this :term:`Cache Group`'s :ref:`cache-group-parent` - or ``null`` if it doesn't have a :ref:`cache-group-parent`
:secondaryParentCachegroupId:   An integer that is the :ref:`cache-group-id` of this :term:`Cache Group`'s :ref:`cache-group-secondary-parent` - or ``null`` if it doesn't have a :ref:`cache-group-secondary-parent`
:secondaryParentCachegroupName: A string containing the :ref:`cache-group-name` of this :term:`Cache Group`'s :ref:`cache-group-secondary-parent` :term:`Cache Group` - or ``null`` if it doesn't have a :ref:`cache-group-secondary-parent`
:shortName:                     A string containing the :ref:`cache-group-short-name` of the :term:`Cache Group`
:typeId:                        An integral, unique identifier for the ':term:`Type`' of the :term:`Cache Group`
:typeName:                      A string that names the :ref:`cache-group-type` of this :term:`Cache Group`

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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
:fallbacks:         An optional field which, when present, should contain an array of strings that are the :ref:`Names <cache-group-name>` of other :term:`Cache Groups` which will be the :ref:`cache-group-fallbacks`\ [#fallbacks]_
:fallbackToClosest: A boolean that sets the :ref:`cache-group-fallback-to-closest` behavior of the :term:`Cache Group`\ [#fallbacks]_

	.. note:: The default value of ``fallbackToClosest`` is ``true``, and if it is ``null`` Traffic Control components will still interpret it as though it were ``true``.

:latitude:                    An optional field which, if present, should be a floating-point number that will define the :ref:`cache-group-latitude` for the :term:`Cache Group`\ [#optional]_
:localizationMethods:         Array of :ref:`cache-group-localization-methods` (as strings)

	.. tip:: This field has no defined meaning if the :ref:`cache-group-type` identified by ``typeId`` is not "EDGE_LOC".

:longitude:                   An optional field which, if present, should be a floating-point number that will define the :ref:`cache-group-longitude` for the :term:`Cache Group`\ [#optional]_
:name:                        The :ref:`cache-group-name` of the :term:`Cache Group`
:parentCachegroupId:          An optional field which, if present, should be an integer that is the :ref:`cache-group-id` of a :ref:`cache-group-parent` for this :term:`Cache Group`.
:secondaryParentCachegroupId: An optional field which, if present, should be an integral, unique identifier for this :term:`Cache Group`'s secondary parent
:shortName:                   An abbreviation of the ``name``
:typeId:                      An integral, unique identifier for the :ref:`Cache Group's Type <cache-group-type>`

	.. note:: The actual, integral, unique identifiers for these :term:`Types` must first be obtained, generally via :ref:`to-api-v3-types`.

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/cachegroups HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 252
	Content-Type: application/json

	{
		"name": "test",
		"shortName": "test",
		"latitude": 0,
		"longitude": 0,
		"fallbackToClosest": true,
		"localizationMethods": [
			"DEEP_CZ",
			"CZ",
			"GEO"
		],
		"typeId": 23,
	}

Response Structure
------------------
:fallbacks:                     An array of strings that are :ref:`Cache Group names <cache-group-name>` that are registered as :ref:`cache-group-fallbacks` for this :term:`Cache Group`\ [#fallbacks]_
:fallbackToClosest:             A boolean value that defines the :ref:`cache-group-fallback-to-closest` behavior of this :term:`Cache Group`\ [#fallbacks]_
:id:                            An integer that is the :ref:`cache-group-id` of the :term:`Cache Group`
:lastUpdated:                   The time and date at which this entry was last updated in :ref:`non-rfc-datetime`
:latitude:                      A floating-point :ref:`cache-group-latitude` for the :term:`Cache Group`
:localizationMethods:           An array of :ref:`cache-group-localization-methods` as strings
:longitude:                     A floating-point :ref:`cache-group-longitude` for the :term:`Cache Group`
:name:                          A string containing the :ref:`cache-group-name` of the :term:`Cache Group`
:parentCachegroupId:            An integer that is the :ref:`cache-group-id` of this :term:`Cache Group`'s :ref:`cache-group-parent` - or ``null`` if it doesn't have a :ref:`cache-group-parent`
:parentCachegroupName:          A string containing the :ref:`cache-group-name` of this :term:`Cache Group`'s :ref:`cache-group-parent` - or ``null`` if it doesn't have a :ref:`cache-group-parent`
:secondaryParentCachegroupId:   An integer that is the :ref:`cache-group-id` of this :term:`Cache Group`'s :ref:`cache-group-secondary-parent` - or ``null`` if it doesn't have a :ref:`cache-group-secondary-parent`
:secondaryParentCachegroupName: A string containing the :ref:`cache-group-name` of this :term:`Cache Group`'s :ref:`cache-group-secondary-parent` :term:`Cache Group` - or ``null`` if it doesn't have a :ref:`cache-group-secondary-parent`
:shortName:                     A string containing the :ref:`cache-group-short-name` of the :term:`Cache Group`
:typeId:                        An integral, unique identifier for the ':term:`Type`' of the :term:`Cache Group`
:typeName:                      A string that names the :ref:`cache-group-type` of this :term:`Cache Group`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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
		"id": 8,
		"name": "test",
		"shortName": "test",
		"latitude": 0,
		"longitude": 0,
		"parentCachegroupName": null,
		"parentCachegroupId": null,
		"secondaryParentCachegroupName": null,
		"secondaryParentCachegroupId": null,
		"fallbackToClosest": true,
		"localizationMethods": [
			"DEEP_CZ",
			"CZ",
			"GEO"
		],
		"typeName": "EDGE_LOC",
		"typeId": 23,
		"lastUpdated": "2019-12-02 22:21:08+00",
		"fallbacks": []
	}}

.. [#fallbacks] Traffic Router will first check for a ``fallbacks`` array and, when that is empty/unset/all the :term:`Cache Groups` in it are also unavailable, will subsequently check for ``fallbackToClosest``. If that is ``true``, then it falls back to the geographically closest :term:`Cache Group` capable of serving the same content or, when it is ``false``/no such :term:`Cache Group` exists/said :term:`Cache Group` is also unavailable, will respond to clients with a failure response indicating the problem.
.. [#optional] While these fields are technically optional, note that if they are not specified many things may break. For this reason, Traffic Portal requires them when creating or editing :term:`Cache Groups`.
