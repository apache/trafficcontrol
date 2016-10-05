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

.. _to-api-v11-region:

Regions
=======

.. _to-api-v11-regions-route:

/api/1.1/regions
++++++++++++++++

**GET /api/1.1/regions**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``id``                | string | Region ID.                                     |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string | Region name.                                   |
  +----------------------+--------+------------------------------------------------+
  |``division``          | string | Division ID.                                   |
  +----------------------+--------+------------------------------------------------+
  |``divisionName``      | string | Division name.                                 |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "6",
           "name": "Atlanta",
           "division": "2",
           "divisionName": "West"
        },
        {
           "id": "7",
           "name": "Denver",
           "division": "2",
           "divisionName": "West"
        },
     ]
    }


**GET /api/1.1/regions/:id**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  |   ``id``  |   yes    | Region id.                                  |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``id``                | string | Region ID.                                     |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string | Region name.                                   |
  +----------------------+--------+------------------------------------------------+
  |``division``          | string | Division ID.                                   |
  +----------------------+--------+------------------------------------------------+
  |``divisionName``      | string | Division name.                                 |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "6",
           "name": "Atlanta",
           "division": "2",
           "divisionName": "West"
        }
     ]
    }


