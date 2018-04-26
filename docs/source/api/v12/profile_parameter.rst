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

.. _to-api-v12-profileparameters:


Profile parameters
==================

.. _to-api-v12-profileparameters-route:

/api/1.2/profileparameters
++++++++++++++++++++++++++

**POST /api/1.2/profileparameters**

    Associate parameter to profile.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Properties**
	This accept two formats: single profile-parameter, profile-parameter array.

	Single profile-parameter format:

	+------------------+----------+----------------------------------------------------+
	| Parameter        | Required | Description                                        |
	+==================+==========+====================================================+
	| ``profileId``    | yes      | profile id.                                        |
	+------------------+----------+----------------------------------------------------+
	| ``parameterId``  | yes      | parameter id.                                      |
	+------------------+----------+----------------------------------------------------+

	Profile-parameter array format:

	+------------------+----------+----------------------------------------------------+
	| Parameter        | Required | Description                                        |
	+==================+==========+====================================================+
	|                  | yes      | profile-parameter array.                           |
	+------------------+----------+----------------------------------------------------+
	| ``>profileId``   | yes      | profile id.                                        |
	+------------------+----------+----------------------------------------------------+
	| ``>parameterId`` | yes      | parameter id.                                      |
	+------------------+----------+----------------------------------------------------+

  **Request Example** ::

    Single profile-parameter format:

    {
      "profileId": 2,
      "parameterId": 6
    }

    Profile-parameter array format:

    [
        {
          "profileId": 2,
          "parameterId": 6
        },
        {
          "profileId": 2,
          "parameterId": 7
        },
        {
          "profileId": 3,
          "parameterId": 6
        }
    ]

 	**Response Properties**

	+-------------------+---------+-----------------------------------------------------+
	|  Parameter        |  Type   |           Description                               |
	+===================+=========+=====================================================+
	| ``response``      | array   | Profile-parameter associations.                     |
	+-------------------+---------+-----------------------------------------------------+
	| ``>profileId``    | string  | Profile id.                                         |
	+-------------------+---------+-----------------------------------------------------+
	| ``>parameterId``  | string  | Parameter id.                                       |
	+-------------------+---------+-----------------------------------------------------+
	| ``alerts``        | array   | A collection of alert messages.                     |
	+-------------------+---------+-----------------------------------------------------+
	| ``>level``        | string  | success, info, warning or error.                    |
	+-------------------+---------+-----------------------------------------------------+
	| ``>text``         | string  | Alert message.                                      |
	+-------------------+---------+-----------------------------------------------------+
	| ``version``       | string  |                                                     |
	+-------------------+---------+-----------------------------------------------------+

  **Response Example** ::

    {
      "response":[
        {
          "profileId": "2",
          "parameterId": "6"
        },
        {
          "profileId": "2",
          "parameterId": "7"
        },
        {
          "profileId": "3",
          "parameterId": "6"
        }
      ]
      "alerts":[
        {
          "level": "success",
          "text": "Profile parameter associations were created."
        }
      ]
    }

|

**DELETE /api/1.2/profileparameters/{:profile_id}/{:parameter_id}**

    Delete a profile parameter association.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+------------------+----------+----------------------------------------------------+
	| Name             | Required | Description                                        |
	+==================+==========+====================================================+
	| ``profile_id``   | yes      | profile id.                                        |
	+------------------+----------+----------------------------------------------------+
	| ``parameter_id`` | yes      | parameter id.                                      |
	+------------------+----------+----------------------------------------------------+

 	**Response Properties**

	+-------------------+--------+-----------------------------------------------------+
	|  Parameter        |  Type  |           Description                               |
	+===================+========+=====================================================+
	| ``alerts``        | array  | A collection of alert messages.                     |
	+-------------------+--------+-----------------------------------------------------+
	| ``>level``        | string | success, info, warning or error.                    |
	+-------------------+--------+-----------------------------------------------------+
	| ``>text``         | string | Alert message.                                      |
	+-------------------+--------+-----------------------------------------------------+
	| ``version``       | string |                                                     |
	+-------------------+--------+-----------------------------------------------------+

  **Response Example** ::

    {
      "alerts":[
        {
          "level": "success",
          "text": "Profile parameter association was deleted."
        }
      ]
    }

|

**POST /api/1.2/profiles/name/{:name}/parameters**

    Associate parameters to a profile. If the parameter does not exist, create it and associate to the profile. If the parameter already exists, associate it to the profile. If the parameter already associate the profile, keep the association.
    If the profile does not exist, the API returns fail.

    Authentication Required: Yes

    Role(s) Required:  admin or oper. If there is parameter's secure equals 1 in the request properties, need admin role. 

	**Request Route Parameters**

	+------------+----------+----------------------------------------------------+
	| Name       | Required | Description                                        |
	+============+==========+====================================================+
	| ``name``   | yes      | profile name.                                      |
	+------------+----------+----------------------------------------------------+

    **Request Properties**
    The request properties accept 2 formats, both single paramter and parameters array formats are acceptable.

    single parameter format:

    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | Name            | Required | Type    | Description                                                                          |
    +=================+==========+=========+======================================================================================+
    | ``name``        | yes      | string  | parameter name                                                                       |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``configFile``  | yes      | string  | parameter config_file                                                                |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``value``       | yes      | string  | parameter value                                                                      |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``secure``      | yes      | integer | secure flag, when 1, the parameter is accessible only by admin users. Defaults to 0. |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+

    array parameters format:

    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | Name            | Required | Type    | Description                                                                          |
    +=================+==========+=========+======================================================================================+
    |                 | yes      | array   | parameters array                                                                     |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>name``       | yes      | string  | parameter name                                                                       |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>configFile`` | yes      | string  | parameter config_file                                                                |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>value``      | yes      | string  | parameter value                                                                      |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>secure``     | yes      | integer | secure flag, when 1, the parameter is accessible only by admin users. Defaults to 0. |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+

  **Request Example** ::

    1. single parameter format exampe:  
    {
        "name":"param1", 
        "configFile":"configFile1",  
        "value":"value1",   
        "secure":0,  
    }

    2. array format example:  
    [
      {
          "name":"param1",
          "configFile":"configFile1",
          "value":"value1",
          "secure":0,
      },
      {
          "name":"param2",
          "configFile":"configFile2",
          "value":"value2",
          "secure":1,
      }
    ]


  **Response Properties** ::

    +------------------+---------+--------------------------------------------------------------------------------------+
    | Name             | Type    | Description                                                                          |
    +==================+=========+======================================================================================+
    | ``response``     |         | Parameters associated with the profile.                                              |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>profileName`` | string  | profile name                                                                         |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>profileId``   | integer | profile index                                                                        |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>parameters``  | array   | parameters array                                                                     |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>id``         | integer | parameter index                                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>name``       | string  | parameter name                                                                       |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>configFile`` | string  | parameter config_file                                                                |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>value``      | string  | parameter value                                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>secure``     | integer | secure flag, when 1, the parameter is accessible only by admin users. Defaults to 0. |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``alerts``       | array   | A collection of alert messages.                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>level``       | string  | success, info, warning or error.                                                     |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>text``        | string  | Alert message.                                                                       |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``version``      | string  |                                                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response":{
        "profileName": "CCR1",
        "profileId" : "12",
        "parameters":[
            {
                "name":"param1",
                "configFile":"configFile1"
                "value":"value1",
                "secure":"0",
            },
            {
                "name":"param2",
                "configFile":"configFile2"
                "value":"value2",
                "secure":"1",
            }
        ]
      }
      "alerts":[
        {
          "level": "success",
          "text": ""Assign parameters successfully to profile CCR1"
        }
      ]
    }

|

**POST /api/1.2/profiles/id/{:id}/parameters**

    Associate parameters to a profile. If the parameter does not exist, create it and associate to the profile. If the parameter already exists, associate it to the profile. If the parameter already associate the profile, keep the association.
    If the profile does not exist, the API returns fail.

    Authentication Required: Yes

    Role(s) Required:  admin or oper. If there is parameter's secure equals 1 in the request properties, need admin role. 

	**Request Route Parameters**

	+------------+----------+----------------------------------------------------+
	| Name       | Required | Description                                        |
	+============+==========+====================================================+
	| ``id``     | yes      | profile name.                                      |
	+------------+----------+----------------------------------------------------+

    **Request Properties**
    The request properties accept 2 formats, both single paramter and parameters array formats are acceptable.

    single parameter format:

    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | Name            | Required | Type    | Description                                                                          |
    +=================+==========+=========+======================================================================================+
    | ``name``        | yes      | string  | parameter name                                                                       |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``configFile``  | yes      | string  | parameter config_file                                                                |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``value``       | yes      | string  | parameter value                                                                      |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``secure``      | yes      | integer | secure flag, when 1, the parameter is accessible only by admin users. Defaults to 0. |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+

    array parameters format:

    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | Name            | Required | Type    | Description                                                                          |
    +=================+==========+=========+======================================================================================+
    |                 | yes      | array   | parameters array                                                                     |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>name``       | yes      | string  | parameter name                                                                       |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>configFile`` | yes      | string  | parameter config_file                                                                |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>value``      | yes      | string  | parameter value                                                                      |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+
    | ``>secure``     | yes      | integer | secure flag, when 1, the parameter is accessible only by admin users. Defaults to 0. |
    +-----------------+----------+---------+--------------------------------------------------------------------------------------+

  **Request Example** ::

    1. single parameter format exampe:  
    {
        "name":"param1", 
        "configFile":"configFile1",  
        "value":"value1",   
        "secure":0,  
    }

    2. array format example:  
    [
      {
          "name":"param1",
          "configFile":"configFile1",
          "value":"value1",
          "secure":0,
      },
      {
          "name":"param2",
          "configFile":"configFile2",
          "value":"value2",
          "secure":1,
      }
    ]


  **Response Properties** ::

    +------------------+---------+--------------------------------------------------------------------------------------+
    | Name             | Type    | Description                                                                          |
    +==================+=========+======================================================================================+
    | ``response``     |         | Parameters associated with the profile.                                              |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>profileName`` | string  | profile name                                                                         |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>profileId``   | integer | profile index                                                                        |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>parameters``  | array   | parameters array                                                                     |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>id``         | integer | parameter index                                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>name``       | string  | parameter name                                                                       |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>configFile`` | string  | parameter config_file                                                                |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>value``      | string  | parameter value                                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>>secure``     | integer | secure flag, when 1, the parameter is accessible only by admin users. Defaults to 0. |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``alerts``       | array   | A collection of alert messages.                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>level``       | string  | success, info, warning or error.                                                     |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``>text``        | string  | Alert message.                                                                       |
    +------------------+---------+--------------------------------------------------------------------------------------+
    | ``version``      | string  |                                                                                      |
    +------------------+---------+--------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response":{
        "profileName": "CCR1",
        "profileId" : "12",
        "parameters":[
            {
                "name":"param1",
                "configFile":"configFile1"
                "value":"value1",
                "secure":"0",
            },
            {
                "name":"param2",
                "configFile":"configFile2"
                "value":"value2",
                "secure":"1",
            }
        ]
      }
      "alerts":[
        {
          "level": "success",
          "text": ""Assign parameters successfully to profile CCR1"
        }
      ]
    }

|

**POST /api/1.2/profileparameter**

  Create one or more profile / parameter assignments.

  Authentication Required: Yes

  Role(s) Required: Admin or Operations

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``profileId``                   | yes      | The ID of the profile.                                            |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``paramIds``                    | yes      | An array of parameter IDs.                                        |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``replace``                     | no       | Replace existing profile/param assignments? (true|false)          |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "profileId": 22,
        "paramIds": [ 2, 3, 4, 5, 6 ],
        "replace": true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``profileId``                      | int    | The ID of the profile.                                            |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``paramIds``                       | array  | An array of parameter IDs.                                        |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``replace``                        |  bool  | Existing profile/param assignments replaced? (true|false).        |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "14 parameters where assigned to the foo profile."
                  }
          ],
        "response": {
            "profileId" : 22,
            "paramIds" : [ 2, 3, 4, 5, 6 ],
            "replace" : true
        }
    }

|

**POST /api/1.2/parameterprofile**

  Create one or more parameter / profile assignments.

  Authentication Required: Yes

  Role(s) Required: Admin or Operations

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``paramId``                     | yes      | The ID of the parameter.                                          |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``profileIds``                  | yes      | An array of profile IDs.                                          |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``replace``                     | no       | Replace existing param/profile assignments? (true|false)          |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "paramId": 22,
        "profileIds": [ 2, 3, 4, 5, 6 ],
        "replace": true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``paramId``                        | int    | The ID of the parameter.                                          |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``profileIds``                     | array  | An array of profile IDs.                                          |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``replace``                        |  bool  | Existing param/profile assignments replaced? (true|false).        |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "14 profiles where assigned to the bar parameter."
                  }
          ],
        "response": {
            "paramId" : 22,
            "profileIds" : [ 2, 3, 4, 5, 6 ],
            "replace" : true
        }
    }

|

