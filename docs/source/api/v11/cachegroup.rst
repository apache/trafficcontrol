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

.. _to-api-v11-cachegroup:

***********
Cache Group
***********

.. _to-api-v11-cachegroups-route:

``/api/1.1/cachegroups``
========================
Extract information about all Cache Groups.

``GET``
-------
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

.. table:: Response Properties

	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| Parameter                         | Type    | Description                                                              |
	+===================================+=========+==========================================================================+
	| ``id``                            | integer | Local unique identifier for the Cache Group                              |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``lastUpdated``                   | string  | The Time / Date this entry was last updated                              |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``latitude``                      | float   | Latitude for the Cache Group                                             |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``longitude``                     | float   | Longitude for the Cache Group                                            |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``name``                          | string  | The name of the Cache Group entry                                        |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``parentCachegroupId``            | integer | Parent cachegroup ID.                                                    |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``parentCachegroupName``          | string  | Parent cachegroup name.                                                  |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``secondaryParentCachegroupId``   | string  | Secondary parent cachegroup ID.                                          |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``secondaryParentCachegroupName`` | string  | Secondary parent cachegroup name.                                        |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``shortName``                     | string  | Abbreviation of the Cache Group Name                                     |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``typeId``                        | integer | Unique identifier for the 'Type' of Cache Group entry                    |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``typeName``                      | string  | The name of the type of Cache Group entry                                |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``fallbackToClosest``             | boolean | If 'true', Traffic Router will direct clients to peers of this Cache     |
	|                                   |         | Group in the event that it becomes unavailable.                          |
	+-----------------------------------+---------+--------------------------------------------------------------------------+

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: json
	:caption: Response Example

	{
		"response": [
			{
				"id": "21",
				"lastUpdated": "2012-09-25 20:27:28",
				"latitude": "0",
				"longitude": "0",
				"name": "dc-chicago",
				"parentCachegroupId": null,
				"parentCachegroupName": null,
				"secondaryParentCachegroupId": null,
				"secondaryParentCachegroupName": null,
				"shortName": "dcchi",
				"typeName": "MID_LOC",
				"typeId": "4",
				"fallbackToClosest": null,
				"localizationMethods": null
			},
			{
				"id": "22",
				"lastUpdated": "2012-09-25 20:27:28",
				"latitude": "0",
				"longitude": "0",
				"name": "dc-chicago-1",
				"parentCachegroupId": null,
				"parentCachegroupName": null,
				"secondaryParentCachegroupId": null,
				"secondaryParentCachegroupName": null,
				"shortName": "dcchi",
				"typeName": "MID_LOC",
				"typeId": "4",
				"fallbackToClosest": null,
				"localizationMethods": null
			}
		],
	}

``/api/1.1/cachegroups/trimmed``
================================
Extract just the names of all Cache Groups.

``GET``
-------
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

.. table:: Response Properties

	+----------------------+--------+------------------------------------------------+
	| Parameter            | Type   | Description                                    |
	+======================+========+================================================+
	|``name``              | string | The name of the cache group                    |
	+----------------------+--------+------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	{
		"response": [
			{
				"name": "dc-chicago"
			},
			{
				"name": "dc-cmc"
			}
		],
	}

``/api/1.1/cachegroups/:id``
============================
Extract information about a single Cache Group.

``GET``
-------

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

.. table:: Request Path Parameters

	+------------------+---------+---------------------------------------------------------------+
	| Parameter        | Type    | Description                                                   |
	+==================+=========+===============================================================+
	| ``id``           | integer | The ID of a Cache Group                                       |
	+------------------+---------+---------------------------------------------------------------+

.. table:: Response Properties

	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| Parameter                         | Type    | Description                                                              |
	+===================================+=========+==========================================================================+
	| ``id``                            | integer | Local unique identifier for the Cache Group                              |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``lastUpdated``                   | string  | The Time / Date this entry was last updated                              |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``latitude``                      | float   | Latitude for the Cache Group                                             |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``longitude``                     | float   | Longitude for the Cache Group                                            |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``name``                          | string  | The name of the Cache Group entry                                        |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``parentCachegroupId``            | integer | Parent cachegroup ID.                                                    |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``parentCachegroupName``          | string  | Parent cachegroup name.                                                  |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``secondaryParentCachegroupId``   | integer | Secondary parent cachegroup ID.                                          |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``secondaryParentCachegroupName`` | string  | Secondary parent cachegroup name.                                        |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``shortName``                     | string  | Abbreviation of the Cache Group Name                                     |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``typeId``                        | integer | Unique identifier for the 'Type' of Cache Group entry                    |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``typeName``                      | string  | The name of the type of Cache Group entry                                |
	+-----------------------------------+---------+--------------------------------------------------------------------------+
	| ``fallbackToClosest``             | boolean | If 'true', Traffic Router will direct clients to peers of this Cache     |
	|                                   |         | Group in the event that it becomes unavailable.                          |
	+-----------------------------------+---------+--------------------------------------------------------------------------+

.. note:: The default value of ``fallbackToClosest`` is 'true', and if it is 'null' Traffic Control components will still interpret it as 'true'.

.. code-block:: json
	:caption: Response Example

	{
		"response": [
			{
				"id": 21,
				"lastUpdated": "2012-09-25 20:27:28",
				"latitude": 0,
				"longitude": 0,
				"name": "dc-chicago",
				"parentCachegroupId": null,
				"parentCachegroupName": null,
				"secondaryParentCachegroupId": null,
				"secondaryParentCachegroupName": null,
				"shortName": "dcchi",
				"typeName": "MID_LOC",
				"typeId": 4,
				"fallbackToClosest": null,
				"localizationMethods": null
			}
		],
	}

``/api/1.1/cachegroup/:parameter_id/parameter``
===============================================
Extract identifying information about all cachegroups with a specific parameter

``GET``
-------
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

.. table:: Request Path Parameters

	+------------------+----------+-----------------------+
	|       Name       | Required | Description           |
	+==================+==========+=======================+
	| ``parameter_id`` | yes      | the ID of a parameter |
	+------------------+----------+-----------------------+

.. table:: Response Properties

	+-----------------+---------------------------------+--------------------------------------------------------------------------------------------------------------------+
	|    Parameter    |  Type                           | Description                                                                                                        |
	+=================+=================================+====================================================================================================================+
	| ``cachegroups`` | array of ``cachegroup`` objects | A list of all cachegroups with an associated parameter identifiable by the ``parameter_id`` request path parameter |
	+-----------------+---------------------------------+--------------------------------------------------------------------------------------------------------------------+

.. table:: ``cachegroup`` Properties

	+----------------+---------+-----------------------------+
	|    Parameter   |  Type   | Description                 |
	+================+=========+=============================+
	| ``name``       | string  | The name of the Cache Group |
	+----------------+---------+-----------------------------+
	| ``id``         | integer | The ID of the Cache Group   |
	+----------------+---------+-----------------------------+

.. code-block:: json
	:caption: Response Example

	{
		"response": {
			"cachegroups": [
				{
						"name": "dc-chicago",
						"id": 21
				},
				{
						"name": "dc-cmc",
						"id": 22
				}
			]
		},
	}


``/api/1.1/cachegroupparameters``
=================================
Extract information about parameters associated with cachegroups

``GET``
-------
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

.. table:: Response Properties

	+--------------------------+-------------------------------------------+-----------------------------------------+
	|        Parameter         |  Type                                     |               Description               |
	+==========================+===========================================+=========================================+
	| ``cachegroupParameters`` | array of ``cachegroupParameter`` objects  | A collection of cache group parameters. |
	+--------------------------+-------------------------------------------+-----------------------------------------+

.. table:: ``cachegroupParameter`` Properties

	+-------------------------+---------+-----------------------------------------+
	| Parameter               | Type    | Description                             |
	+=========================+=========+=========================================+
	| ``parameter``           | integer | ID of the parameter                     |
	+-------------------------+---------+-----------------------------------------+
	| ``last_updated``        | string  | Date and time of last modification      |
	+-------------------------+---------+-----------------------------------------+
	| ``cachegroup``          | string  | Name of the Cache Group                 |
	+-------------------------+---------+-----------------------------------------+

.. code-block:: json
	:caption: Response Example

	{
		"response": {
			"cachegroupParameters": [
				{
						"parameter": "379",
						"last_updated": "2013-08-05 18:49:37",
						"cachegroup": "us-ca-sanjose"
				},
				{
						"parameter": "380",
						"last_updated": "2013-08-05 18:49:37",
						"cachegroup": "us-ca-sanjose"
				},
				{
						"parameter": "379",
						"last_updated": "2013-08-05 18:49:37",
						"cachegroup": "us-ma-woburn"
				}
			]
		},
	}


``/api/1.1/cachegroups/:parameter_id/parameter/available``
==========================================================
Extract information about what Cache Groups are available to be assigned a specific parameter.

``GET``
-------
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

.. table:: Request Route Parameters

	+------------------+----------+-----------------------------------------------------------------+
	|       Name       | Required | Description                                                     |
	+==================+==========+=================================================================+
	| ``parameter_id`` | yes      | The ID of the parameter for which availability is being checked |
	+------------------+----------+-----------------------------------------------------------------+

.. table:: Response Properties

	+----------------------+---------+------------------------------------------------+
	| Parameter            | Type    | Description                                    |
	+======================+=========+================================================+
	|``name``              | string  | The name of the Cache Group                    |
	+----------------------+---------+------------------------------------------------+
	|``id``                | integer | The ID of the Cache Group                      |
	+----------------------+---------+------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	{
		"response": [
			{
				"name": "dc-chicago",
				"id": 21
			},
			{
				"name": "dc-cmc",
				"id": 22
			}
		],
	}
