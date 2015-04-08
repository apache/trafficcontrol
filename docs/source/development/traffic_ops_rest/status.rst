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

.. _to-api-status:

Status
======

**GET /api/1.1/statuses.json**

Authentication Required: Yes

Response Content Type: application/json

**Response Messages**

::


  HTTP Status Code: 200
  Reason: Success

**Response Properties**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
|``lastUpdated``       | string |                                                |
+----------------------+--------+------------------------------------------------+
|``name``              | string |                                                |
+----------------------+--------+------------------------------------------------+
|``id``                | string |                                                |
+----------------------+--------+------------------------------------------------+
|``description``       | string |                                                |
+----------------------+--------+------------------------------------------------+

**Response Example**


::

  {
   "response": [
      {
         "lastUpdated": "2013-02-13 23:34:29",
         "name": "ADMIN_DOWN",
         "id": "4",
         "description": "Temporary down. Edge: XMPP client will send status OFFLINE to CCR."
      },
      {
         "lastUpdated": "2013-02-13 23:34:29",
         "name": "CCR_IGNORE",
         "id": "5",
         "description": "Edge: 12M will not include caches in this state in CCR config files. Mid: N\/A for now"
      }
   ],
   "version": "1.1"
  }

For error messages, see :ref:`reference-label-401`.
