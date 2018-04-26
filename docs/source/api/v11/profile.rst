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

.. _to-api-v11-profile:

Profiles
========

.. _to-api-v11-profiles-route:

/api/1.1/profiles
+++++++++++++++++

**GET /api/1.1/profiles**

	Authentication Required: Yes

	Role(s) Required: None

	**Response Properties**

	+-----------------+--------+----------------------------------------------------+
	|    Parameter    |  Type  |                    Description                     |
	+=================+========+====================================================+
	| ``lastUpdated`` | array  | The Time / Date this server entry was last updated |
	+-----------------+--------+----------------------------------------------------+
	| ``name``        | string | The name for the profile                           |
	+-----------------+--------+----------------------------------------------------+
	| ``id``          | string | Primary key                                        |
	+-----------------+--------+----------------------------------------------------+
	| ``description`` | string | The description for the profile                    |
	+-----------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "lastUpdated": "2012-10-08 19:34:45",
            "name": "CCR_TOP",
            "id": "8",
            "description": "Content Router for top.foobar.net"
        }
     ]
    }

|

**GET /api/1.1/profiles/trimmed**

	Authentication Required: Yes

	Role(s) Required: None

	**Response Properties**

	+-----------------+--------+----------------------------------------------------+
	|    Parameter    |  Type  |                    Description                     |
	+=================+========+====================================================+
	| ``name``        | string | The name for the profile                           |
	+-----------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "name": "CCR_TOP"
        }
     ]
    }

|

**GET /api/1.1/profiles/:id**

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+-----------------+------------+------------------------------------------------+
	|    Parameter    |  Required  |                    Description                 |
	+=================+============+================================================+
	| ``id``          |    yes     | The ID of the profile.                         |
	+-----------------+------------+------------------------------------------------+

	**Response Properties**

	+-----------------+--------+----------------------------------------------------+
	|    Parameter    |  Type  |                    Description                     |
	+=================+========+====================================================+
	| ``lastUpdated`` | array  | The Time / Date this server entry was last updated |
	+-----------------+--------+----------------------------------------------------+
	| ``name``        | string | The name for the profile                           |
	+-----------------+--------+----------------------------------------------------+
	| ``id``          | string | Primary key                                        |
	+-----------------+--------+----------------------------------------------------+
	| ``description`` | string | The description for the profile                    |
	+-----------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "lastUpdated": "2012-10-08 19:34:45",
            "name": "CCR_TOP",
            "id": "8",
            "description": "Content Router for top.foobar.net"
        }
     ]
    }

|

