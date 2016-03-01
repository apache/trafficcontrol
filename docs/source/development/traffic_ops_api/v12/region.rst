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

.. _to-api-v12-region:

Regions
=======

.. _to-api-v12-regions-route:

/api/1.2/regions
++++++++++++++++

**GET /api/1.2/regions.json**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "name": "Atlanta",
           "id": "6"
        },
        {
           "name": "Beltway",
           "id": "1"
        }
     ],
    }

|

**POST /api/1.2/divisions/:division_name/regions**
  Create Region

  Authentication Required: Yes

  Role(s) Required: admin or oper

  division_name - The name of division to create new region into.

  ** Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              | Required | Description                                    |
  +===================+==========+================================================+
  | ``division_name`` | yes      | The name of division will create new region in |
  +-------------------+----------+------------------------------------------------+

  **Request Properties**

  +-------------------+----------+------------------------------------------+
  | Parameter         | Required | Description                              |
  +===================+==========+==========================================+
  | ``name``          | yes      | The name of the region                   |
  +-------------------+----------+------------------------------------------+

  **Request Example** ::

    {
        "name": "myregion1",
    }

|

  **Response Properties**

  +-------------------+--------+-------------------------------------------+
  | Parameter         | Type   | Description                               |
  +===================+========+===========================================+
  | ``name``          | string | name of region created                    |
  +-------------------+--------+-------------------------------------------+
  | ``id``            | string | id of region created                      |
  +-------------------+--------+-------------------------------------------+
  | ``division_name`` | string | the division name the region belongs to.  |
  +-------------------+--------+-------------------------------------------+
  | ``division_id``   | string | the id of division the region belongs to. |
  +-------------------+--------+-------------------------------------------+

  **Response Example** ::

    {
      "response": {
        'division_name': 'mydivision1',
        'divsion_id': '4',
        'name': 'myregion1',
        'id': '19'
       }
    }

|
