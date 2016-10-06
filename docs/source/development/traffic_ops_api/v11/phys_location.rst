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

.. _to-api-v11-phys-loc:

Physical Location
=================

.. _to-api-v11-phys-loc-route:

/api/1.1/phys_locations
+++++++++++++++++++++++

**GET /api/1.1/phys_locations**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``address``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``comments``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phone``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``poc``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``region``            | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``regionId``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``shortName``         | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``state``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``zip``               | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "region": "Mile High",
           "region": "4",
           "poc": "Jane Doe",
           "lastUpdated": "2014-10-02 08:22:43",
           "name": "Albuquerque",
           "comments": "Albuquerque",
           "phone": "(123) 555-1111",
           "state": "NM",
           "email": "jane.doe@email.com",
           "city": "Albuquerque",
           "zip": "87107",
           "id": "2",
           "address": "123 East 3rd St",
           "shortName": "Albuquerque"
        },
        {
           "region": "Mile High",
           "region": "4",
           "poc": "Jane Doe",
           "lastUpdated": "2014-10-02 08:22:43",
           "name": "Albuquerque",
           "comments": "Albuquerque",
           "phone": "(123) 555-1111",
           "state": "NM",
           "email": "jane.doe@email.com",
           "city": "Albuquerque",
           "zip": "87107",
           "id": "2",
           "address": "123 East 3rd St",
           "shortName": "Albuquerque"
        }
     ]
    }

|

**GET /api/1.1/phys_locations/:id**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  | ``id``    | yes      | Physical location ID.                       |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``address``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``comments``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phone``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``poc``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``region``            | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``regionId``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``shortName``         | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``state``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``zip``               | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "region": "Mile High",
           "region": "4",
           "poc": "Jane Doe",
           "lastUpdated": "2014-10-02 08:22:43",
           "name": "Albuquerque",
           "comments": "Albuquerque",
           "phone": "(123) 555-1111",
           "state": "NM",
           "email": "jane.doe@email.com",
           "city": "Albuquerque",
           "zip": "87107",
           "id": "2",
           "address": "123 East 3rd St",
           "shortName": "Albuquerque"
        }
     ]
    }

|


**GET /api/1.1/phys_locations/trimmed**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+---------+------------------------------------------------+
  | Parameter            | Type    | Description                                    |
  +======================+=========+================================================+
  |``name``              | string  |                                                |
  +----------------------+---------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "name": "Albuquerque"
        },
        {
           "name": "Ashburn"
        }
     ],
    }


