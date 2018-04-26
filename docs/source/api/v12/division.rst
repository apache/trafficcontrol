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

.. _to-api-v12-division:

Divisions
=========

.. _to-api-v12-division-route:

/api/1.2/divisions
++++++++++++++++++

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
  |``lastUpdated``       | string |                                                 |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Division name                                   |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "1"
           "name": "Central",
           "lastUpdated": "2014-10-02 08:22:43"
        },
        {
           "id": "2"
           "name": "West",
           "lastUpdated": "2014-10-02 08:22:43"
        }
     ]
    }

|


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
  |``lastUpdated``       | string |                                                 |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Division name                                   |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "1"
           "name": "Central",
           "lastUpdated": "2014-10-02 08:22:43"
        }
     ]
    }

|


**PUT /api/1.2/divisions/:id**
  Update a division

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              |   Type   |                 Description                    |
  +===================+==========+================================================+
  | ``id``            | int      | Division id.                                   |
  +-------------------+----------+------------------------------------------------+

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
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

	{
		"alerts": [
			{
				"level": "success",
				"text": "Division update was successful."
			}
		],
		"response": {
			"id": "1",
			"lastUpdated": "2014-03-18 08:57:39",
			"name": "mydivision1"
		}
	}
  
|


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
