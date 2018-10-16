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

************************
``/api/1.x/cachegroups``
************************

``GET``
=======
Extract information about all Cache Groups.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	| ``type``        | no       | Filter cache groups by Type ID.                   |
	+-----------------+----------+---------------------------------------------------+

Response Structure
------------------
:fallbackToClosest:             If ``true``, Traffic Router will direct clients to peers of this Cache Group in the event that it becomes unavailable.
:id:                            A numeric, unique identifier for the Cache Group
:lastUpdated:                   The time and date at which this entry was last updated in ISO format
:latitude:                      Latitude for the Cache Group
:longitude:                     Longitude for the Cache Group
:name:                          The name of the Cache Group entry
:parentCachegroupId:            ID of this Cache Group's parent Cache Group (if any)
:parentCachegroupName:          Name of this Cache Group's parent Cache Group (if any)
:secondaryParentCachegroupId:   ID of this Cache Group's secondary parent Cache Group (if any)
:secondaryParentCachegroupName: Name of this Cache Group's secondary parent Cache Group (if any)
:shortName:                     Abbreviation of the Cache Group name
:typeId:                        Unique identifier for the 'Type' of Cache Group entry
:typeName:                      The name of the type of Cache Group entry

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 1,
			"name": "TRAFFIC_ANALYTICS",
			"shortName": "TRAFFIC_ANALYTICS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": null,
			"parentCachegroupId": null,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "TC_LOC",
			"typeId": 47,
			"lastUpdated": "2018-10-15 13:35:35+00"
		},
		{
			"id": 2,
			"name": "TRAFFIC_OPS",
			"shortName": "TRAFFIC_OPS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": null,
			"parentCachegroupId": null,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "TC_LOC",
			"typeId": 47,
			"lastUpdated": "2018-10-15 13:35:35+00"
		},
		{
			"id": 3,
			"name": "TRAFFIC_OPS_DB",
			"shortName": "TRAFFIC_OPS_DB",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": null,
			"parentCachegroupId": null,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "TC_LOC",
			"typeId": 47,
			"lastUpdated": "2018-10-15 13:35:36+00"
		},
		{
			"id": 4,
			"name": "TRAFFIC_PORTAL",
			"shortName": "TRAFFIC_PORTAL",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": null,
			"parentCachegroupId": null,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "TC_LOC",
			"typeId": 47,
			"lastUpdated": "2018-10-15 13:35:36+00"
		},
		{
			"id": 5,
			"name": "TRAFFIC_STATS",
			"shortName": "TRAFFIC_STATS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": null,
			"parentCachegroupId": null,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "TC_LOC",
			"typeId": 47,
			"lastUpdated": "2018-10-15 13:35:36+00"
		},
		{
			"id": 6,
			"name": "CDN_in_a_Box_Mid",
			"shortName": "ciabMid",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"parentCachegroupName": null,
			"parentCachegroupId": null,
			"secondaryParentCachegroupName": null,
			"secondaryParentCachegroupId": null,
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "MID_LOC",
			"typeId": 24,
			"lastUpdated": "2018-10-15 13:35:36+00"
		},
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
			"fallbackToClosest": null,
			"localizationMethods": null,
			"typeName": "EDGE_LOC",
			"typeId": 23,
			"lastUpdated": "2018-10-15 13:35:36+00"
		}
	]}

``POST``
========
Creates a cache group

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Data Parameters

	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| Name                            | Required | Type    |  Description                                                      |
	+=================================+==========+=========+===================================================================+
	| ``name``                        | yes      | string  | The name of the Cache Group entry                                 |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``shortName``                   | yes      | string  | Abbreviation of the Cache Group Name                              |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``latitude``                    | no       | float   | Latitude for the Cache Group                                      |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``longitude``                   | no       | float   | Longitude for the Cache Group                                     |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``parentCachegroup``            | no       | string  | Name of Parent Cache Group entry.                                 |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``secondaryParentCachegroup``   | no       | string  | Name of Secondary Parent Cache Group entry.                       |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``localizationMethods``         | no       | array   | Array of enabled localization methods (as strings)                |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``typeId``                      | yes      | integer | The type of Cache Group entry, "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+
	| ``fallbackToClosest``           | no       | boolean | Behaviour on configured fallbacks failure                         |
	+---------------------------------+----------+---------+-------------------------------------------------------------------+

Response Structure
------------------
:fallbackToClosest:             If ``true``, Traffic Router will direct clients to peers of this Cache Group in the event that it becomes unavailable.
:id:                            A numeric, unique identifier for the Cache Group
:lastUpdated:                   The time and date at which this entry was last updated in ISO format
:latitude:                      Latitude for the Cache Group
:longitude:                     Longitude for the Cache Group
:name:                          The name of the Cache Group entry
:parentCachegroupId:            ID of this Cache Group's parent Cache Group (if any)
:parentCachegroupName:          Name of this Cache Group's parent Cache Group (if any)
:secondaryParentCachegroupId:   ID of this Cache Group's secondary parent Cache Group (if any)
:secondaryParentCachegroupName: Name of this Cache Group's secondary parent Cache Group (if any)
:shortName:                     Abbreviation of the Cache Group name
:typeId:                        Unique identifier for the 'Type' of Cache Group entry
:typeName:                      The name of the type of Cache Group entry

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Cachegroup successfully created: cache_group_edge"
		}
	],
	"response": {
		"longitude" : 45,
		"lastUpdated" : "2016-01-25 13:55:30",
		"shortName" : "cg_edge",
		"name" : "cache_group_edge",
		"parentCachegroup" : "cache_group_mid",
		"secondaryParentCachegroup" : null,
		"latitude" : 12,
		"localizationMethods": [
			"CZ",
			"GEO"
		],
		"typeName" : "EDGE_LOC",
		"id" : 104,
		"parentCachegroupId" : 103,
		"secondaryParentCachegroupId" : null,
		"fallbackToClosest":true
	}}


.. This doesn't appear to exist anymore - can't reproduce in CIAB nor production
.. ``/api/1.1/cachegroups/:parameter_id/parameter/available``
.. ==========================================================
