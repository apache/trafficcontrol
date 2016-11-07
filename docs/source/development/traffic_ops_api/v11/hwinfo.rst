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


.. _to-api-v11-hwinfo:

Hardware Info
=============

.. _to-api-v11-hwinfo-route:

/api/1.1/hwinfo
+++++++++++++++

**GET /api/1.1/hwinfo.json**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +--------------------+--------+----------------------------------------------------------------------+
  | Parameter          | Type   | Description                                                          |
  +====================+========+======================================================================+
  | ``serverId``       | string | Local unique identifier for this specific server's hardware info     |
  +--------------------+--------+----------------------------------------------------------------------+
  | ``serverHostName`` | string | Hostname for this specific server's hardware info                    |
  +--------------------+--------+----------------------------------------------------------------------+
  | ``lastUpdated``    | string | The Time and Date for the last update for this server.               |
  +--------------------+--------+----------------------------------------------------------------------+
  | ``val``            | string | Freeform value used to track anything about a server's hardware info |
  +--------------------+--------+----------------------------------------------------------------------+
  | ``description``    | string | Freeform description for this specific server's hardware info        |
  +--------------------+--------+----------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "serverId": "odol-atsmid-cen-09",
           "lastUpdated": "2014-05-27 09:06:02",
           "val": "D1S4",
           "description": "Physical Disk 0:1:0"
        },
        {
           "serverId": "odol-atsmid-cen-09",
           "lastUpdated": "2014-05-27 09:06:02",
           "val": "D1S4",
           "description": "Physical Disk 0:1:1"
        }
     ],
    }

