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

.. _to-api-jobs:

********
``jobs``
********

``GET``
=======
Get all jobs (currently limited to invalidate content (PURGE) jobs) sorted by start time (descending).

:Auth. Required: Yes
:Roles Required: "operations" or "admin"\ [1]_

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name   | Required | Description                                                                                                  |
	+========+==========+==============================================================================================================+
	| dsId   | no       | Return only invalidation jobs pending on the Delivery Service identified by this integral, unique identifier |
	+--------+----------+--------------------------------------------------------------------------------------------------------------+
	| userId | no       | Return only invalidation jobs created by the user identified by this integral, unique identifier             |
	+--------+----------+--------------------------------------------------------------------------------------------------------------+

.. note:: If the ``dsId`` parameter is given, an error will be returned if the thereby identified Delivery Service is not visible to the logged-in user's Tenant

Response Structure
------------------
+----------------------+--------+-------------------------------------------------+
| Parameter            | Type   | Description                                     |
+======================+========+=================================================+
|``id``                |  int   | Job id                                          |
+----------------------+--------+-------------------------------------------------+
|``assetUrl``          | string | URL of the asset to invalidate.                 |
+----------------------+--------+-------------------------------------------------+
|``deliveryService``   | string | Unique identifier of the job's DS.              |
+----------------------+--------+-------------------------------------------------+
|``keyword``           | string | Job keyword (PURGE)                             |
+----------------------+--------+-------------------------------------------------+
|``parameters``        | string | Parameters associated with the job.             |
+----------------------+--------+-------------------------------------------------+
|``startTime``         | string | Start time of the job.                          |
+----------------------+--------+-------------------------------------------------+
|``createdBy``         | string | Username that initiated the job.                |
+----------------------+--------+-------------------------------------------------+

**Response Example** ::

	{
	 "response": [
			{
				 "id": 1
				 "assetUrl": "http:\/\/foo-bar.domain.net\/taco.html",
				 "deliveryService": "foo-bar",
				 "keyword": "PURGE",
				 "parameters": "TTL:48h",
				 "startTime": "2015-05-14 08:56:36-06",
				 "createdBy": "jdog24"
			},
			{
				 "id": 2
				 "assetUrl": "http:\/\/foo-bar.domain.net\/bell.html",
				 "deliveryService": "foo-bar",
				 "keyword": "PURGE",
				 "parameters": "TTL:72h",
				 "startTime": "2015-05-16 08:56:36-06",
				 "createdBy": "jdog24"
			}
	 ]
	}

.. [1] Only jobs running on Delivery Services that the logged-in user's Tenant has permissions to see will be returned - regardless of role.

``POST``
========
Creates a new content invalidation job

:Auth. Required: Yes
:Roles Required:

**GET /api/1.2/jobs/:id**

Get a job by ID (currently limited to invalidate content (PURGE) jobs).

Authentication Required: Yes

Role(s) Required: Operations or Admin

**Response Properties**

+----------------------+--------+-------------------------------------------------+
| Parameter            | Type   | Description                                     |
+======================+========+=================================================+
|``id``                |  int   | Job id                                          |
+----------------------+--------+-------------------------------------------------+
|``assetUrl``          | string | URL of the asset to invalidate.                 |
+----------------------+--------+-------------------------------------------------+
|``deliveryService``   | string | Unique identifier of the job's DS.              |
+----------------------+--------+-------------------------------------------------+
|``keyword``           | string | Job keyword (PURGE)                             |
+----------------------+--------+-------------------------------------------------+
|``parameters``        | string | Parameters associated with the job.             |
+----------------------+--------+-------------------------------------------------+
|``startTime``         | string | Start time of the job.                          |
+----------------------+--------+-------------------------------------------------+
|``createdBy``         | string | Username that initiated the job.                |
+----------------------+--------+-------------------------------------------------+

**Response Example** ::

	{
	 "response": [
			{
				 "id": 1
				 "assetUrl": "http:\/\/foo-bar.domain.net\/taco.html",
				 "deliveryService": "foo-bar",
				 "keyword": "PURGE",
				 "parameters": "TTL:48h",
				 "startTime": "2015-05-14 08:56:36-06",
				 "createdBy": "jdog24"
			}
	 ]
	}

|
