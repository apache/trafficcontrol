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

.. _to-api-v12-steering-targets:

Steering Targets
================

.. _to-api-v12-steering-target-route:

**GET /api/1.2/steering/:dsId/targets**

  Get all targets for a steering delivery service.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-------------+----------+---------------------------------------------+
  |   Name      | Required |                Description                  |
  +=============+==========+=============================================+
  | ``dsId``    |   yes    | DS ID.                                      |
  +-------------+----------+---------------------------------------------+

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``deliveryServiceId`` |  int   | DS ID                                           |
  +----------------------+--------+-------------------------------------------------+
  |``deliveryService``   | string | DS XML ID                                       |
  +----------------------+--------+-------------------------------------------------+
  |``targetId``          |  int   | Target DS ID                                    |
  +----------------------+--------+-------------------------------------------------+
  |``target``            | string | Target DS XML ID                                |
  +----------------------+--------+-------------------------------------------------+
  |``value``             |  int   | Value is weight or order depending on type      |
  +----------------------+--------+-------------------------------------------------+
  |``typeId``            |  int   | Steering target type ID                         |
  +----------------------+--------+-------------------------------------------------+
  |``type``              | string | Steering target type name                       |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "deliveryServiceId": 1
           "deliveryService": "steering-ds-one",
           "targetId": 2,
           "target": "steering-target-one",
           "value": 1,
           "typeId": 35,
           "type": "STEERING_ORDER"
        },
        {
           "deliveryServiceId": 1
           "deliveryService": "steering-ds-one",
           "targetId": 3,
           "target": "steering-target-two",
           "value": 2,
           "typeId": 35,
           "type": "STEERING_ORDER"
        },
     ]
    }

|

**GET /api/1.2/steering/:dsId/targets/:targetId**

  Get a steering target.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-------------+----------+---------------------------------------------+
  |   Name      | Required |                Description                  |
  +=============+==========+=============================================+
  | ``dsId``    |   yes    | DS ID.                                      |
  +-------------+----------+---------------------------------------------+
  | ``targetId``|   yes    | DS Target ID.                               |
  +-------------+----------+---------------------------------------------+

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``deliveryServiceId`` |  int   | DS ID                                           |
  +----------------------+--------+-------------------------------------------------+
  |``deliveryService``   | string | DS XML ID                                       |
  +----------------------+--------+-------------------------------------------------+
  |``targetId``          |  int   | Target DS ID                                    |
  +----------------------+--------+-------------------------------------------------+
  |``target``            | string | Target DS XML ID                                |
  +----------------------+--------+-------------------------------------------------+
  |``value``             |  int   | Value is weight or order depending on type      |
  +----------------------+--------+-------------------------------------------------+
  |``typeId``            |  int   | Steering target type ID                         |
  +----------------------+--------+-------------------------------------------------+
  |``type``              | string | Steering target type name                       |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "deliveryServiceId": 1
           "deliveryService": "steering-ds-one",
           "targetId": 2,
           "target": "steering-target-one",
           "value": 1,
           "typeId": 35,
           "type": "STEERING_ORDER"
        }
     ]
    }

|


**PUT /api/1.2/steering/:dsId/targets/:targetId**

  Update a steering target.

  Authentication Required: Yes

  Role(s) Required: Portal

  **Request Route Parameters**

  +-------------+----------+---------------------------------------------+
  |   Name      | Required |                Description                  |
  +=============+==========+=============================================+
  | ``dsId``    |   yes    | DS ID.                                      |
  +-------------+----------+---------------------------------------------+
  | ``targetId``|   yes    | DS Target ID.                               |
  +-------------+----------+---------------------------------------------+

  **Request Properties**

  +------------------------+----------+--------------------------+
  | Parameter              | Required | Description              |
  +========================+==========+==========================+
  | ``value``              | yes      | Target value             |
  +------------------------+----------+--------------------------+
  | ``typeId``             | yes      | Target type ID           |
  +------------------------+----------+--------------------------+

  **Request Example** ::

    {
        "value": 34,
        "typeId": 46,
    }

|

  **Response Properties**

  +------------------------+----------+--------------------------+
  | Parameter              | Type     | Description              |
  +========================+==========+==========================+
  | ``deliveryServiceId``  | int      | Steering DS ID           |
  +------------------------+----------+--------------------------+
  | ``deliveryService``    | string   | DS XML ID                |
  +------------------------+----------+--------------------------+
  | ``targetId``           | int      | Target DS ID             |
  +------------------------+----------+--------------------------+
  | ``target``             | string   | Target DS XML ID         |
  +------------------------+----------+--------------------------+
  | ``value``              | string   | Target value             |
  +------------------------+----------+--------------------------+
  | ``typeId``             | int      | Target type ID           |
  +------------------------+----------+--------------------------+
  | ``type``               | string   | Steering target type name|
  +------------------------+----------+--------------------------+

  **Response Example** ::

	{
		"response": {
			"deliveryServiceId": 1,
			"deliveryService": "steering-ds-one",
			"targetId": 2,
			"target": "steering-target-two",
			"value": "34",
			"typeId": 45,
			"type": "STEERING_ORDER"
		},
		"alerts": [
			{
				"level": "success",
				"text": "Delivery service steering target update was successful."
			}
		]
	}

|


**POST /api/1.2/steering/:dsId/targets**

  Create a steering target.

  Authentication Required: Yes

  Role(s) Required: Portal

  **Request Route Parameters**

  +-------------+----------+---------------------------------------------+
  |   Name      | Required |                Description                  |
  +=============+==========+=============================================+
  | ``dsId``    |   yes    | DS ID.                                      |
  +-------------+----------+---------------------------------------------+

  **Request Properties**

  +------------------------+----------+--------------------------+
  | Parameter              | Required | Description              |
  +========================+==========+==========================+
  | ``targetId``           | yes      | Target DS ID             |
  +------------------------+----------+--------------------------+
  | ``value``              | yes      | Target value             |
  +------------------------+----------+--------------------------+
  | ``typeId``             | yes      | Target type ID           |
  +------------------------+----------+--------------------------+

  **Request Example** ::

    {
        "targetId": 6,
        "value": 22,
        "typeId": 47,
    }

|

  **Response Properties**

  +------------------------+----------+--------------------------+
  | Parameter              | Type     | Description              |
  +========================+==========+==========================+
  | ``deliveryServiceId``  | int      | Steering DS ID           |
  +------------------------+----------+--------------------------+
  | ``deliveryService``    | string   | DS XML ID                |
  +------------------------+----------+--------------------------+
  | ``targetId``           | int      | Target DS ID             |
  +------------------------+----------+--------------------------+
  | ``target``             | string   | Target DS XML ID         |
  +------------------------+----------+--------------------------+
  | ``value``              | string   | Target value             |
  +------------------------+----------+--------------------------+
  | ``typeId``             | int      | Target type ID           |
  +------------------------+----------+--------------------------+
  | ``type``               | string   | Steering target type name|
  +------------------------+----------+--------------------------+

  **Response Example** ::

	{
		"response": {
			"deliveryServiceId": 1,
			"deliveryService": "steering-ds-one",
			"targetId": 6,
			"target": "steering-target-six",
			"value": "22",
			"typeId": 47,
			"type": "STEERING_ORDER"
		},
		"alerts": [
			{
				"level": "success",
				"text": "Delivery service target creation was successful."
			}
		]
	}

|

**DELETE /api/1.2/steering/:dsId/targets/:targetId**

  Delete a steering target.

  Authentication Required: Yes

  Role(s) Required: Portal

  **Request Route Parameters**

  +-------------+----------+---------------------------------------------+
  |   Name      | Required |                Description                  |
  +=============+==========+=============================================+
  | ``dsId``    |   yes    | DS ID.                                      |
  +-------------+----------+---------------------------------------------+
  | ``targetId``|   yes    | DS Target ID.                               |
  +-------------+----------+---------------------------------------------+

  **Response Example** ::

    {
          "alerts": [
                    {
                            "level": "success",
                            "text": "Delivery service target delete was successful."
                    }
            ],
    }

|

