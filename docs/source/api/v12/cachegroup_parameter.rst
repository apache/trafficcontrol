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

.. _to-api-v12-cachegroupparameters:


Cache Group parameters
======================

.. _to-api-v12-cachegroupparameters-route:

/api/1.2/cachegroupparameters
+++++++++++++++++++++++++++++

**POST /api/1.2/cachegroupparameters**

	Assign parameter(s) to cache group(s).

	Authentication Required: Yes

	Role(s) Required:  Admin or Operations

	**Request Properties**

	Two formats are acceptable.

	Single cachegroup-parameter format:

	+------------------+----------+----------------------------------------------------+
	| Parameter        | Required | Description                                        |
	+==================+==========+====================================================+
	| ``cacheGroupId`` | yes      | cache group id.                                    |
	+------------------+----------+----------------------------------------------------+
	| ``parameterId``  | yes      | parameter id.                                      |
	+------------------+----------+----------------------------------------------------+

	Profile-parameter array format:

	+-----------------------+----------+----------------------------------------------------+
	| Parameter             | Required | Description                                        |
	+=======================+==========+====================================================+
	|                       | yes      | cachegroup-parameter array.                        |
	+-----------------------+----------+----------------------------------------------------+
	| ``>cacheGroupId``     | yes      | cache group id.                                    |
	+-----------------------+----------+----------------------------------------------------+
	| ``>parameterId``      | yes      | parameter id.                                      |
	+-----------------------+----------+----------------------------------------------------+

  **Request Example** ::

    Single cachegroup-parameter format:

    {
      "cacheGroupId": 2,
      "parameterId": 6
    }

    Cachegroup-parameter array format:

    [
        {
          "cacheGroupId": 2,
          "parameterId": 6
        },
        {
          "cacheGroupId": 2,
          "parameterId": 7
        },
        {
          "cacheGroupId": 3,
          "parameterId": 6
        }
    ]

 	**Response Properties**

	+-------------------+---------+-----------------------------------------------------+
	|  Parameter        |  Type   |           Description                               |
	+===================+=========+=====================================================+
	| ``response``      | array   | Cache group-parameter associations.                 |
	+-------------------+---------+-----------------------------------------------------+
	| ``>cacheGroupId`` | string  | Cache Group id.                                     |
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
          "cacheGroupId": "2",
          "parameterId": "6"
        },
        {
          "cacheGroupId": "2",
          "parameterId": "7"
        },
        {
          "cacheGroupId": "3",
          "parameterId": "6"
        }
      ]
      "alerts":[
        {
          "level": "success",
          "text": "Cache group parameter associations were created."
        }
      ]
    }

|

**DELETE /api/1.2/cachegroupparameters/{:cachegroup_id}/{:parameter_id}**

    Delete a cache group parameter association.

	Authentication Required: Yes

	Role(s) Required:  Admin or Operations

	**Request Route Parameters**

	+------------------+----------+----------------------------------------------------+
	| Name             | Required | Description                                        |
	+==================+==========+====================================================+
	| ``cachegroup_id``| yes      | cache group id.                                    |
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
          "text": "Cache group parameter association was deleted."
        }
      ]
    }

|
