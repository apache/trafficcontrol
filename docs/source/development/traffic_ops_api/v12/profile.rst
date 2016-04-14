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

.. _to-api-v12-profile:


Profiles
========

.. _to-api-v12-profiles-route:

/api/1.2/profiles
+++++++++++++++++

**GET /api/1.2/profiles.json**

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

**GET /api/1.2/profiles/trimmed.json**

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

**POST /api/1.2/profiles**
    Create a new empty  profile. 

	Authentication Required: Yes

	Role(s) Required: admin or oper

	**Request Properties**

	+-----------------------+--------+----------------------------------------------------+
	|    Parameter          |  Type  |                    Description                     |
	+=======================+========+====================================================+
	| ``name``              | string | The name of the new profile                        |
	+-----------------------+--------+----------------------------------------------------+
	| ``description``       | string | new profile description                            |
	+-----------------------+--------+----------------------------------------------------+

  **Request Example** ::

    {
      "name": "CCR_COPY",
      "description": "CCR_COPY description",
    }

|

	**Response Properties**

	+-----------------------+--------+----------------------------------------------------+
	|    Parameter          |  Type  |                    Description                     |
	+=======================+========+====================================================+
	| ``id``                | string | Id of the new profile                              |
	+-----------------------+--------+----------------------------------------------------+
	| ``name``              | string | The name of the new profile                        |
	+-----------------------+--------+----------------------------------------------------+
	| ``description``       | string | new profile description                            |
	+-----------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "66",
            "name": "CCR_COPY",
            "description": "CCR_COPY description",
        }
     ]
    }

|

**POST /api/1.2/profiles/name/:profile_name/copy/:profile_copy_from**
    Copy profile to a new profile. The new profile name must not exist. 

	Authentication Required: Yes

	Role(s) Required: admin or oper

	**Request Route Parameters**
   
	+-----------------------+----------+-------------------------------+
	| Name                  | Required | Description                   |
	+=======================+==========+===============================+
	| ``profile_name``      | yes      | The name of profile to copy   |
	+-----------------------+----------+-------------------------------+
	| ``profile_copy_from`` | yes      | The name of profile copy from |
	+-----------------------+----------+-------------------------------+


	**Response Properties**

	+-----------------------+--------+----------------------------------------------------+
	|    Parameter          |  Type  |                    Description                     |
	+=======================+========+====================================================+
	| ``id``                | string | Id of the new profile                              |
	+-----------------------+--------+----------------------------------------------------+
	| ``name``              | string | The name of the new profile                        |
	+-----------------------+--------+----------------------------------------------------+
	| ``profile_copy_from`` | string | The name of profile to copy                        |
	+-----------------------+--------+----------------------------------------------------+
	| ``id_copy_from``      | string | The id of profile to copy                          |
	+-----------------------+--------+----------------------------------------------------+
	| ``description``       | string | new profile's description (copied)                 |
	+-----------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "66",
            "name": "CCR_COPY",
            "profile_copy_from": "CCR1",
            "description": "CCR_COPY description",
            "id_copy_from": "3"
        }
     ]
    }

|
