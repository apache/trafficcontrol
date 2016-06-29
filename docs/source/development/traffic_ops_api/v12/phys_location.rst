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

.. _to-api-v12-phys-loc:

Physical Location
=================

.. _to-api-v12-phys-loc-route:

/api/1.2/phys_locations
+++++++++++++++++++++++

**GET /api/1.2/phys_locations.json**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``region``            | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``poc``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``comments``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``phone``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``state``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``email``             | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``city``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``zip``               | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``address``           | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``shortName``         | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "region": "Mile High",
           "poc": "Jane Doe",
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
           "region": "Chicago",
           "poc": "John Doe",
           "name": "Chicago",
           "comments": "",
           "phone": "(321) 555-1111",
           "state": "IL",
           "email": "john.doe@email.com",
           "city": "Chicago",
           "zip": "60636",
           "id": "3",
           "address": "123 East 4th Street",
           "shortName": "chicago"
        }
     ],
    }

|

**GET /api/1.2/phys_locations/trimmed.json**

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
     ]
    }

|

**POST /api/1.2/regions/:region_name/phys_locations**
  Create physical location.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  region_name: the name of the region to create physical location into.

  **Request Route Parameters**

  +-----------------+----------+-----------------------------------+
  | Name            | Required | Description                       |
  +=================+==========+===================================+
  | ``region_name`` | yes      | The name of the physical location |
  +-----------------+----------+-----------------------------------+

  **Request Properties**
  
  +-----------------+----------+---------------------------------------------------+
  | Parameter       | Required | Description                                       |
  +=================+==========+===================================================+
  | ``name``        | yes      | The name of the location                          |
  +-----------------+----------+---------------------------------------------------+
  | ``shortName``   | yes      | The short name of the location                    |
  +-----------------+----------+---------------------------------------------------+
  | ``address``     | yes      |                                                   |
  +-----------------+----------+---------------------------------------------------+
  | ``city``        | yes      |                                                   |
  +-----------------+----------+---------------------------------------------------+
  | ``state``       | yes      |                                                   |
  +-----------------+----------+---------------------------------------------------+
  | ``zip``         | yes      |                                                   |
  +-----------------+----------+---------------------------------------------------+
  | ``phone``       | no       |                                                   |
  +-----------------+----------+---------------------------------------------------+
  | ``poc``         | no       | Point of contact                                  |
  +-----------------+----------+---------------------------------------------------+
  | ``email``       | no       |                                                   |
  +-----------------+----------+---------------------------------------------------+
  | ``comments``    | no       |                                                   |
  +-----------------+----------+---------------------------------------------------+

  **Request Example** ::

    {
        "name" : "my physical location1",
        "shortName" : "myphylocation1",
        "address" : "",
        "city" : "Shanghai",
        "state": "SH",
        "zip": "200000",
        "comments": "this is physical location1"
    }
   
|

  **Response Properties**

  +-----------------+--------+---------------------------------------------------+
  | Parameter       | Type   | Description                                       |
  +=================+========+===================================================+
  | ``id``          | string | The id of the physical location created.          |
  +-----------------+--------+---------------------------------------------------+
  | ``name``        | string | The name of the location                          |
  +-----------------+--------+---------------------------------------------------+
  | ``shortName``   | string | The short name of the location                    |
  +-----------------+--------+---------------------------------------------------+
  | ``regionName``  | string | The region name the physical location belongs to. |
  +-----------------+--------+---------------------------------------------------+
  | ``regionId``    | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``address``     | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``city``        | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``state``       | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``zip``         | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``phone``       | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``poc``         | string | Point of contact                                  |
  +-----------------+--------+---------------------------------------------------+
  | ``email``       | string |                                                   |
  +-----------------+--------+---------------------------------------------------+
  | ``comments``    | string |                                                   |
  +-----------------+--------+---------------------------------------------------+

  **Response Example** ::

    {
      "response": {
        'shortName': 'myphylocati',
        'regionName': 'myregion1',
        'name': 'my physical location1',
        'poc': '',
        'phone': '',
        'comments': 'this is physical location1',
        'state': 'SH',
        'email': '',
        'zip': '20000',
        'region_id': '20',
        'city': 'Shanghai',
        'address': '',
        'id': '200'
     }
   }

|
