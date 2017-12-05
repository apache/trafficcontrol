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

.. _to-api-v12-profile:


Profiles
========

.. _to-api-v12-profiles-route:

/api/1.2/profiles
+++++++++++++++++

**GET /api/1.2/profiles**

	Authentication Required: Yes

	Role(s) Required: None

	**Request Query Parameters**

	+---------------+----------+----------------------------------------------------+
	|    Name       | Required |                    Description                     |
	+===============+==========+====================================================+
	| ``param``     |   no     | Used to filter profiles by parameter ID.           |
	+---------------+----------+----------------------------------------------------+
	| ``cdn``       |   no     | Used to filter profiles by CDN ID.                 |
	+---------------+----------+----------------------------------------------------+

	**Response Properties**

	+---------------------+--------+------------------------------------------------------+
	|      Parameter      |  Type  |                    Description                       |
	+=====================+========+======================================================+
	| ``id``              | string | Primary key                                          |
	+---------------------+--------+------------------------------------------------------+
	| ``name``            | string | The name for the profile                             |
	+---------------------+--------+------------------------------------------------------+
	| ``description``     | string | The description for the profile                      |
	+---------------------+--------+------------------------------------------------------+
	| ``cdn``             |  int   | The CDN ID                                           |
	+---------------------+--------+------------------------------------------------------+
	| ``cdnName``         | string | The CDN name                                         |
	+---------------------+--------+------------------------------------------------------+
	| ``type``            | string | Profile type                                         |
	+---------------------+--------+------------------------------------------------------+
	| ``routingDisabled`` |  bool  | Traffic router routing disabled - defaults to false. |
	+---------------------+--------+------------------------------------------------------+
	| ``lastUpdated``     | array  | The Time / Date this server entry was last updated   |
	+---------------------+--------+------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "8",
            "name": "EDGE_27_PROFILE",
            "description": "A profile with all the Foo parameters"
            "cdn": 1
            "cdnName": "cdn1"
            "type": "ATS_PROFILE"
            "routingDisabled": false
            "lastUpdated": "2012-10-08 19:34:45",
        }
     ]
    }

|

**GET /api/1.2/profiles/trimmed**

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
            "name": "EDGE_27_PROFILE"
        }
     ]
    }

|

**GET /api/1.2/profiles/:id**

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+-----------------+------------+------------------------------------------------+
	|    Parameter    |  Required  |                    Description                 |
	+=================+============+================================================+
	| ``id``          |    yes     | The ID of the profile.                         |
	+-----------------+------------+------------------------------------------------+

	**Response Properties**

	+---------------------+--------+----------------------------------------------------+
	|      Parameter      |  Type  |                    Description                     |
	+=====================+========+====================================================+
	| ``id``              | string | Primary key                                        |
	+---------------------+--------+----------------------------------------------------+
	| ``name``            | string | The name for the profile                           |
	+---------------------+--------+----------------------------------------------------+
	| ``description``     | string | The description for the profile                    |
	+---------------------+--------+----------------------------------------------------+
	| ``cdn``             |  int   | The CDN ID                                         |
	+---------------------+--------+----------------------------------------------------+
	| ``cdnName``         | string | The CDN name                                       |
	+---------------------+--------+----------------------------------------------------+
	| ``type``            | string | Profile type                                       |
	+---------------------+--------+----------------------------------------------------+
	| ``routingDisabled`` |  bool  | Traffic router routing disabled                    |
	+---------------------+--------+----------------------------------------------------+
	| ``lastUpdated``     | array  | The Time / Date this server entry was last updated |
	+---------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "8",
            "name": "EDGE_27_PROFILE",
            "description": "A profile with all the Foo parameters"
            "cdn": 1
            "cdnName": "cdn1"
            "type": "ATS_PROFILE"
            "routingDisabled": true
            "lastUpdated": "2012-10-08 19:34:45",
        }
     ]
    }

|


**POST /api/1.2/profiles**
    Create a new empty profile.

	Authentication Required: Yes

	Role(s) Required: admin or oper

	**Request Properties**

	+-----------------------+--------+----------+-----------------------------------------+
	|  Parameter            |  Type  | Required |           Description                   |
	+=======================+========+==========+=========================================+
	| ``name``              | string | yes      | Profile name                            |
	+-----------------------+--------+----------+-----------------------------------------+
	| ``description``       | string | yes      | Profile description                     |
	+-----------------------+--------+----------+-----------------------------------------+
	| ``cdn``               |  int   | no       | CDN ID                                  |
	+-----------------------+--------+----------+-----------------------------------------+
	| ``type``              | string | yes      | Profile type                            |
	+-----------------------+--------+----------+-----------------------------------------+
	| ``routingDisabled``   |  bool  | no       | Traffic router routing disabled.        |
	|                       |        |          | Defaults to false.                      |
	+-----------------------+--------+----------+-----------------------------------------+


  **Request Example** ::

    {
      "name": "EDGE_28_PROFILE",
      "description": "EDGE_28_PROFILE description",
      "cdn": 1,
      "type": "ATS_PROFILE",
      "routingDisabled": false
    }

|

	**Response Properties**

	+-----------------------+--------+----------------------------------------------------+
	|    Parameter          |  Type  |                    Description                     |
	+=======================+========+====================================================+
	| ``id``                | string | Profile ID                                         |
	+-----------------------+--------+----------------------------------------------------+
	| ``name``              | string | Profile name                                       |
	+-----------------------+--------+----------------------------------------------------+
	| ``description``       | string | Profile description                                |
	+-----------------------+--------+----------------------------------------------------+
	| ``cdn``               |  int   | CDN ID                                             |
	+-----------------------+--------+----------------------------------------------------+
	| ``type``              | string | Profile type                                       |
	+-----------------------+--------+----------------------------------------------------+
	| ``routingDisabled``   |  bool  | Traffic router routing disabled                    |
	+-----------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "66",
            "name": "EDGE_28_PROFILE",
            "description": "EDGE_28_PROFILE description",
            "cdn": 1,
            "type": "ATS_PROFILE",
            "routingDisabled": false
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
	| ``profileCopyFrom``   | string | The name of profile to copy                        |
	+-----------------------+--------+----------------------------------------------------+
	| ``idCopyFrom``        | string | The id of profile to copy                          |
	+-----------------------+--------+----------------------------------------------------+
	| ``description``       | string | new profile's description (copied)                 |
	+-----------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "66",
            "name": "CCR_COPY",
            "profileCopyFrom": "CCR1",
            "description": "CCR_COPY description",
            "idCopyFrom": "3"
        }
     ]
    }

|

**PUT /api/1.2/profiles/{:id}**

    Allows user to edit a profile.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	| ``id``          | yes      | profile id.                                       |
	+-----------------+----------+---------------------------------------------------+

	**Request Properties**

	+-----------------------+--------+----------+--------------------------------------------+
	|  Parameter            |  Type  | Required |           Description                      |
	+=======================+========+==========+============================================+
	| ``name``              | string | yes      | Profile name                               |
	+-----------------------+--------+----------+--------------------------------------------+
	| ``description``       | string | yes      | Profile description                        |
	+-----------------------+--------+----------+--------------------------------------------+
	| ``cdn``               |  int   | no       | CDN ID - must use the same ID as any       |
	|                       |        |          | servers assigned to the profile.           |
	+-----------------------+--------+----------+--------------------------------------------+
	| ``type``              | string | yes      | Profile type                               |
	+-----------------------+--------+----------+--------------------------------------------+
	| ``routingDisabled``   |  bool  | no       | Traffic router routing disabled.           |
	|                       |        |          | When not present, value defaults to false. |
	+-----------------------+--------+----------+--------------------------------------------+

  **Request Example** ::

    {
      "name": "EDGE_28_PROFILE",
      "description": "EDGE_28_PROFILE description",
      "cdn": 1,
      "type": "ATS_PROFILE",
      "routingDisabled": false
    }

|

 	**Response Properties**

	+-----------------------+--------+----------------------------------------------------+
	|    Parameter          |  Type  |                    Description                     |
	+=======================+========+====================================================+
	| ``id``                | string | Profile ID                                         |
	+-----------------------+--------+----------------------------------------------------+
	| ``name``              | string | Profile name                                       |
	+-----------------------+--------+----------------------------------------------------+
	| ``description``       | string | Profile description                                |
	+-----------------------+--------+----------------------------------------------------+
	| ``cdn``               |  int   | CDN ID                                             |
	+-----------------------+--------+----------------------------------------------------+
	| ``type``              | string | Profile type                                       |
	+-----------------------+--------+----------------------------------------------------+
	| ``routingDisabled``   |  bool  | Traffic router routing disabled                    |
	+-----------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
      "response":{
        "id": "219",
        "name": "EDGE_28_PROFILE",
        "description": "EDGE_28_PROFILE description"
        "cdn": 1
        "type": "ATS_PROFILE",
        "routingDisabled": false
      }
      "alerts":[
        {
          "level": "success",
          "text": "Profile was updated: 219"
        }
      ]
    }

|

**DELETE /api/1.2/profiles/{:id}**

  Allows user to delete a profile.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+-----------------+----------+----------------------------+
	| Name            | Required | Description                |
	+=================+==========+============================+
	| ``id``          | yes      | profile id.                |
	+-----------------+----------+----------------------------+

 	**Response Properties**

	+-------------+--------+----------------------------------+
	|  Parameter  |  Type  |           Description            |
	+=============+========+==================================+
	| ``alerts``  | array  | A collection of alert messages.  |
	+-------------+--------+----------------------------------+
	| ``>level``  | string | success, info, warning or error. |
	+-------------+--------+----------------------------------+
	| ``>text``   | string | Alert message.                   |
	+-------------+--------+----------------------------------+
	| ``version`` | string |                                  |
	+-------------+--------+----------------------------------+

  **Response Example** ::

    {
      "alerts": [
        {
          "level": "success",
          "text": "Profile was deleted."
        }
      ]
    }

|

