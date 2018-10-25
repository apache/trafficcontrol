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

************************
``/api/1.x/coordinates``
************************
.. versionadded:: 1.3

``GET``
=======
Gets a list of all coordinates in the Traffic Ops database

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------------+----------+---------------------------------------------------------------------+
	| Name            | Required | Description                                                         |
	+=================+==========+=====================================================================+
	| ``id``          | no       | Return only coordinates that have this integral, unique identifier  |
	+-----------------+----------+---------------------------------------------------------------------+
	| ``name``        | no       | Return only coordinates with this name                              |
	+-----------------+----------+---------------------------------------------------------------------+

Response Structure
------------------
:id:          Integral, unique, identifier for this coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in a ``ctime``-like format
:latitude:    Latitude of the coordinate
:longitude:   Longitude of the coordinate
:name:        The name of the coordinate - typically this just reflects the name of the Cache Group for which the coordinate was created

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 1,
			"name": "from_cachegroup_TRAFFIC_ANALYTICS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:04+00"
		},
		{
			"id": 2,
			"name": "from_cachegroup_TRAFFIC_OPS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:04+00"
		},
		{
			"id": 3,
			"name": "from_cachegroup_TRAFFIC_OPS_DB",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:04+00"
		},
		{
			"id": 4,
			"name": "from_cachegroup_TRAFFIC_PORTAL",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:04+00"
		},
		{
			"id": 5,
			"name": "from_cachegroup_TRAFFIC_STATS",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:04+00"
		},
		{
			"id": 6,
			"name": "from_cachegroup_CDN_in_a_Box_Mid",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:04+00"
		},
		{
			"id": 7,
			"name": "from_cachegroup_CDN_in_a_Box_Edge",
			"latitude": 38.897663,
			"longitude": -77.036574,
			"lastUpdated": "2018-10-24 16:07:05+00"
		}
	]}

``POST``
========
Creates a new coordinate pair

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Data Parameters

	+---------------+-----------------------+----------+-------------------------------------------------------------------+
	| Name          | Type                  | Required | Description                                                       |
	+===============+=======================+==========+===================================================================+
	| ``name``      | string                | yes      | The name of the new coordinate                                    |
	+---------------+-----------------------+----------+-------------------------------------------------------------------+
	| ``latitude``  | float (<=180, >=-180) | no       | The latitude of the new coordinate                                |
	+---------------+-----------------------+----------+-------------------------------------------------------------------+
	| ``longitude`` | float (<=90, >=-90)   | no       | The longitude of the new coordinate                               |
	+---------------+-----------------------+----------+-------------------------------------------------------------------+

Response Structure
------------------
:id:          Integral, unique, identifier for the newly created coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in a ``ctime``-like format
:latitude:    Latitude of the newly created coordinate
:longitude:   Longitude of the newly created coordinate
:name:        The name of the coordinate

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "coordinate was created.",
			"level": "success"
		}
	],
	"response": {
		"id": 10,
		"name": "testquest",
		"latitude": 0,
		"longitude": 0,
		"lastUpdated": "2018-10-25 16:40:33+00"
	}}

``PUT``
=======
Updates a coordinate

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+------------------------------------------------------------+
	| Name | Required | Description                                                |
	+======+==========+============================================================+
	| id   | yes      | The integral, unique identifier of the coordinate to edit  |
	+------+----------+------------------------------------------------------------+

.. table:: Request Data Parameters

	+---------------+-----------------------+----------+-------------------------------------------------------------------+
	| Name          | Type                  | Required | Description                                                       |
	+===============+=======================+==========+===================================================================+
	| ``name``      | string                | yes      | The new name of the coordinate entry                              |
	+---------------+-----------------------+----------+-------------------------------------------------------------------+
	| ``latitude``  | float (<=180, >=-180) | no       | The new latitude of the coordinate                                |
	+---------------+-----------------------+----------+-------------------------------------------------------------------+
	| ``longitude`` | float (<=90, >=-90)   | no       | The new longitude of the coordinate                               |
	+---------------+-----------------------+----------+-------------------------------------------------------------------+

Response Structure
------------------
:id:          Integral, unique, identifier for the coordinate pair
:lastUpdated: The time and date at which this entry was last updated, in a ``ctime``-like format
:latitude:    Latitude of the coordinate
:longitude:   Longitude of the coordinate
:name:        The name of the coordinate

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "coordinate was updated.",
			"level": "success"
		}
	],
	"response": {
		"id": 10,
		"name": "testquest",
		"latitude": -90,
		"longitude": 180,
		"lastUpdated": "2018-10-25 17:08:55+00"
	}}

``DELETE``
==========
Deletes a coordinate

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+-------------------------------------------------------------+
	| Name | Required | Description                                                 |
	+======+==========+=============================================================+
	| id   | yes      | The integral, unique identifier of the coordinate to delete |
	+------+----------+-------------------------------------------------------------+

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "coordinate was deleted.",
			"level": "success"
		}
	]}

