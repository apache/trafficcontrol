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

.. _to-api-v12-profileparameters:


Profile parameters
==================

.. _to-api-v12-profileparameters-route:

/api/1.2/profileparameters
++++++++++++++++++++++++++

**POST /api/1.2/profileparameters/{:id}**

    Associate parameters to a profile.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+------------------+----------+----------------------------------------------------+
	| Name             | Required | Description                                        |
	+==================+==========+====================================================+
	| ``id``           | yes      | profile id.                                        |
	+------------------+----------+----------------------------------------------------+

	**Request Properties**

	+------------------+----------+----------------------------------------------------+
	| Parameter        | Required | Description                                        |
	+==================+==========+====================================================+
	| ``parametersId`` | yes      | id array of parameters to associate to the profile,|
	|                  |          | for example: [ 1, 2, 32 ]                          |
	+------------------+----------+----------------------------------------------------+

  **Request Example** ::

    {
      "parametersId": [ 4, 5 ]
    }

 	**Response Properties**

	+-------------------+--------+-----------------------------------------------------+
	|  Parameter        |  Type  |           Description                               |
	+===================+========+=====================================================+
	| ``response``      |        | Parameters associated with the profile.             |
	+-------------------+--------+-----------------------------------------------------+
	| ``>id``           | string | Profile id.                                         |
	+-------------------+--------+-----------------------------------------------------+
	| ``>parametersId`` | array  | id array of parameters associated with the profile, |
	|                   |        | for example [ 1, 2, 32, 100 ]                       |
	+-------------------+--------+-----------------------------------------------------+
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
      "response":{
        "id": "3",
        "parametersId": [ 3, 4, 5 ]
      }
      "alerts":[
        {
          "level": "success",
          "text": "Parameters were associated to profile: 3"
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
