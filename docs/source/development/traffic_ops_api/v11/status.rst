.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
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

.. _to-api-v11-status:

Status
======

.. _to-api-v11-statuses-route:

/api/1.1/statuses
+++++++++++++++++

**GET /api/1.1/statuses.json**

  Retrieves a list of the server status codes available. May be useful when the status is retrieved from other APIs as a number and not a string.

  Authentication Required: Yes

  **Response Properties**

  +-----------------+--------+--------------------------------------------------------------------------------+
  |    Parameter    |  Type  |                                  Description                                   |
  +=================+========+================================================================================+
  | ``lastUpdated`` | string | The Time / Date this server entry was last updated                             |
  +-----------------+--------+--------------------------------------------------------------------------------+
  | ``name``        | string | The string equivalent of the status                                            |
  +-----------------+--------+--------------------------------------------------------------------------------+
  | ``id``          | string | The id with which Traffic Ops stores this status, and references it internally |
  +-----------------+--------+--------------------------------------------------------------------------------+
  | ``description`` | string | A short description of the status                                              |
  +-----------------+--------+--------------------------------------------------------------------------------+

  **Response Example** ::

       {
        "response": [
          {
            "description": "Temporary down. Edge: XMPP client will send status OFFLINE to CCR, otherwise similar to REPORTED. Mid: Server will not be included in parent.config files for its edge caches",
            "id": "4",
            "name": "ADMIN_DOWN",
            "lastUpdated": "2013-02-13 16:34:29"
          },
          {
            "lastUpdated": "2013-02-13 16:34:29",
            "name": "CCR_IGNORE",
            "id": "5",
            "description": "Edge: 12M will not include caches in this state in CCR config files. Mid: N\/A for now"
          },
          {
            "description": "Edge: Puts server in CCR config file in this state, but CCR will never route traffic to it. Mid: Server will not be included in parent.config files for its edge caches",
            "id": "1",
            "lastUpdated": "2013-02-13 16:34:29",
            "name": "OFFLINE"
          },
          {
            "id": "2",
            "description": "Edge: Puts server in CCR config file in this state, and CCR will always route traffic to it. Mid: Server will be included in parent.config files for its edges",
            "lastUpdated": "2013-02-13 16:34:29",
            "name": "ONLINE"
          },
          {
            "id": "3",
            "description": "Edge: Puts server in CCR config file in this state, and CCR will adhere to the health protocol. Mid: N\/A for now",
            "name": "REPORTED",
            "lastUpdated": "2013-02-13 16:34:29"
          }
        ],
      }

