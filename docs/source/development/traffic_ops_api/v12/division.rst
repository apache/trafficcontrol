.. 
.. Copyright 2016 Cisco
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

.. _to-api-v12-division:

Divisions
=========

.. _to-api-v12-division-route:

/api/1.2/divisions
++++++++++++++++

**GET /api/1.2/divisions**
  Get all divisions.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``id``                | string | Division id                                     |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Division name                                   |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "1"
           "name": "Central"
        },
        {
           "id": "2"
           "name": "West"
        }
     ],
    }

|

/api/1.2/divisions/:id
++++++++++++++++

**GET /api/1.2/divisions/:id**
  Get division by Id.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``id``                | string | Division id                                     |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Division name                                   |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "1"
           "name": "Central"
        }
     ],
    }

|

/api/1.2/divisions
++++++++++++++++++

**POST /api/1.2/divisions**
  Create division

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Properties**

  +-----------+----------+--------------------------+
  | Parameter | Required | Description              |
  +===========+==========+==========================+
  | ``name``  | yes      | The name of the division |
  +-----------+----------+--------------------------+
 
  **Request Example** ::

    {
        "name": "mydivision1"
    }

|

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
      "response": {
        'name': 'mydivision1',
        'id': '4'
      }
    }

|
