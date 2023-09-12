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

.. _to-api-cachegroups-id:

**********************
``cachegroups/{{ID}}``
**********************

``PUT``
=======
Update :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: CACHE-GROUP:UPDATE, CACHE-GROUP:READ, TYPE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------------------------------------------------+
	| Parameter | Description                                        |
	+===========+====================================================+
	| ID        | The :ref:`cache-group-id` of a :term:`Cache Group` |
	+-----------+----------------------------------------------------+

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

	.. note:: The actual, integral, unique identifiers for these :term:`Types` must first be obtained, generally via :ref:`to-api-types`.

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/cachegroups/8 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 118
	Content-Type: application/json

	{
		"latitude": 0.0,
		"longitude": 0.0,
		"name": "test",
		"fallbacks": [],
		"fallbackToClosest": true,
		"shortName": "test",
		"typeId": 23,
		"localizationMethods": ["GEO"]
	}

Response Structure
------------------
:fallbacks:         An array of strings that are :ref:`Cache Group names <cache-group-name>` that are registered as :ref:`cache-group-fallbacks` for this :term:`Cache Group`\ [#fallbacks]_
:fallbackToClosest: A boolean value that defines the :ref:`cache-group-fallback-to-closest` behavior of this :term:`Cache Group`\ [#fallbacks]_
:id:                An integer that is the :ref:`cache-group-id` of the :term:`Cache Group`
:lastUpdated:       The time and date at which this entry was last updated in :rfc:`3339`

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

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
	Whole-Content-Sha512: t1W65/2kj25QyHt0Ib0xpBaAR2sXu2kOsRZ49WjKZp/AK5S1YWhX7VNWCuUGiN1VNM4QRNqODC/7ewhYDFUncA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 19:14:28 GMT
	Content-Length: 385

	{ "alerts": [
		{
			"text": "cachegroup was updated.",
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
		"fallbacks": [],
		"fallbackToClosest": true,
		"localizationMethods": [
			"GEO"
		],
		"typeName": "EDGE_LOC",
		"typeId": 23,
		"lastUpdated": "2023-05-30T19:52:58.183642+00:00"
	}}


``DELETE``
==========
Delete a :term:`Cache Group`. A :term:`Cache Group` which has assigned servers or is the :ref:`cache-group-parent` of one or more other :term:`Cache Groups` cannot be deleted.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: CACHE-GROUP:DELETE, CACHE-GROUP:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+------------------------------------------------------------------+
	| Parameter | Description                                                      |
	+===========+==================================================================+
	| ID        | The :ref:`cache-group-id` of a :term:`Cache Group` to be deleted |
	+-----------+------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/cachegroups/42 HTTP/1.1
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
	Whole-Content-Sha512: 5jZBgO7h1eNF70J/cmlbi3Hf9KJPx+WLMblH/pSKF3FWb/10GUHIN35ZOB+lN5LZYCkmk3izGbTFkiruG8I41Q==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 19:14:28 GMT
	Content-Length: 57

	{ "alerts": [
		{
			"text": "cachegroup was deleted.",
			"level": "success"
		}
	]}

.. [#fallbacks] Traffic Router will first check for a ``fallbacks`` array and, when that is empty/unset/all the :term:`Cache Groups` in it are also unavailable, will subsequently check for ``fallbackToClosest``. If that is ``true``, then it falls back to the geographically closest :term:`Cache Group` capable of serving the same content or, when it is ``false``/no such :term:`Cache Group` exists/said :term:`Cache Group` is also unavailable, will respond to clients with a failure response indicating the problem.
.. [#optional] While these fields are technically optional, note that if they are not specified many things may break. For this reason, Traffic Portal requires them when creating or editing :term:`Cache Groups`.
