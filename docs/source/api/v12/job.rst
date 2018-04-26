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

.. _to-api-v12-job:

Jobs
====

.. _to-api-v12-job-route:

/api/1.2/jobs
++++++++++++++++++

**GET /api/1.2/jobs**

  Get all jobs (currently limited to invalidate content (PURGE) jobs) sorted by start time (descending).

  Authentication Required: Yes

  Role(s) Required: Operations or Admin

  **Request Query Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``dsId``        | no       | Filter jobs by Delivery Service ID.               |
  +-----------------+----------+---------------------------------------------------+
  | ``userId``      | no       | Filter jobs by User ID.                           |
  +-----------------+----------+---------------------------------------------------+

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

|


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
