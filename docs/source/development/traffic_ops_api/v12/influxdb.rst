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

.. _to-api-v12-influxdb:

InfluxDB
==========

.. Note:: The documentation needs a thorough review!

**GET /api/1.2/traffic_monitor/stats.json**

Authentication Required: Yes

Role(s) Required: None

**Response Properties**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
| ``aaData``           | array  |                                                |
+----------------------+--------+------------------------------------------------+

**Response Example**
::

  {
   "aaData": [
      [
         "0",
         "ALL",
         "ALL",
         "ALL",
         "true",
         "ALL",
         "142035",
         "172365661.85"
      ],
      [
         1,
         "EDGE1_TOP_421_PSPP",
         "odol-atsec-atl-03",
         "us-ga-atlanta",
         "1",
         "REPORTED",
         "596",
         "923510.04",
         "69.241.82.126"
      ]
   ],
  }
  