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

.. _to-api-cachegroups_id:

*******************************
``/api/1.x/cachegroups/{{id}}``
*******************************
Extracts information about a single Cache Group

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+---------+---------------------------------------------------------------+
	| Parameter        | Type    | Description                                                   |
	+==================+=========+===============================================================+
	| ``id``           | integer | The ID of a Cache Group                                       |
	+------------------+---------+---------------------------------------------------------------+

Response Structure
------------------
:fallbackToClosest:             If 'true', Traffic Router will direct clients to peers of this Cache Group in the event that it becomes unavailable.
:id:                            Local unique identifier for the Cache Group
:lastUpdated:                   The Time / Date this entry was last updated
:latitude:                      Latitude for the Cache Group
:longitude:                     Longitude for the Cache Group
:name:                          The name of the Cache Group entry
:parentCachegroupId:            Parent Cache Group ID
:parentCachegroupName:          Parent Cache Group name
:secondaryParentCachegroupId:   Secondary parent Cache Group ID
:secondaryParentCachegroupName: Secondary parent Cache Group name
:shortName:                     Abbreviation of the Cache Group Name
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
		}
	]}
