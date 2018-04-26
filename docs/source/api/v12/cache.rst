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


.. _to-api-v12-cache:

Cache
=====

.. _to-api-v12-cache-route:

/api/1.2/caches/stats
+++++++++++++++++++++

**GET /api/1.2/caches/stats**

  Retrieves cache stats from Traffic Monitor. Also includes rows for aggregates.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  | Parameter                  | Type          | Description                                                                             |
  +============================+===============+=========================================================================================+
  |``profile``                 | string        | The profile of the cache.                                                               |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``cachegroup``              | string        | The cache group of the cache.                                                           |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``hostname``                | string        | The hostname of the cache.                                                              |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``ip``                      | string        | The IP address of the cache.                                                            |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``status``                  | string        | The status of the cache.                                                                |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``healthy``                 | string        | Has Traffic Monitor marked the cache as healthy or unhealthy?                           |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``connections``             | string        | Cache connections                                                                       |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+
  |``kbps``                    | string        | Cache kbps out                                                                          |
  +----------------------------+---------------+-----------------------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "profile": "ALL",
          "cachegroup": "ALL",
          "hostname": "ALL",
          "ip": null,
          "status": "ALL",
          "healthy": true,
          "connections": 934424,
          "kbps": 618631875
        },
        {
          "profile": "EDGE1_FOO_721-ATS621-45",
          "cachegroup": "us-nm-albuquerque",
          "hostname": "foo-bar-alb-01",
          "ip": "2.2.2.2",
          "status": "REPORTED",
          "healthy": true,
          "connections": 373,
          "kbps": 390136
        },
      ]
    }

|
