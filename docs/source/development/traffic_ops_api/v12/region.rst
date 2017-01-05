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

.. _to-api-v12-region:

Regions
=======

.. _to-api-v12-regions-route:

/api/1.2/regions
++++++++++++++++

**GET /api/1.1/regions**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +---------------+----------+---------------------------------------------+
  |   Name        | Required |                Description                  |
  +===============+==========+=============================================+
  |  ``division`` |    no    | Filter regions by Division ID.              |
  +---------------+----------+---------------------------------------------+

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


**PUT /api/1.2/regions/:id**
  Update a region

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              |   Type   |                 Description                    |
  +===================+==========+================================================+
  | ``id``            | int      | Region id.                                     |
  +-------------------+----------+------------------------------------------------+

  **Request Properties**

  +----------------+----------+--------------------------+
  | Parameter      | Required | Description              |
  +================+==========+==========================+
  | ``name``       | yes      | The name of the region   |
  +----------------+----------+--------------------------+
  | ``division``   | yes      | The division Id          |
  +----------------+----------+--------------------------+

  **Request Example** ::

    {
        "name": "myregion1",
        "division": "4"
    }

|

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``division``          | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``divisionName``      | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

	{
		"alerts": [
			{
				"level": "success",
				"text": "Region update was successful."
			}
		],
		"response": {
			"id": "1",
			"lastUpdated": "2014-03-18 08:57:39",
			"name": "myregion1",
			"division": "4",
			"divisionName": "mydivision1"
		}
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
  | ``divisionName``  | string | the division name the region belongs to.  |
  +-------------------+--------+-------------------------------------------+
  | ``divisionId``    | string | the id of division the region belongs to. |
  +-------------------+--------+-------------------------------------------+

  **Response Example** ::

    {
      "response": {
        'divisionName': 'mydivision1',
        'divsionId': '4',
        'name': 'myregion1',
        'id': '19'
       }
    }

|
