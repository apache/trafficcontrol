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

.. _to-api-v13-coordinate:

Coordinate
==========

.. _to-api-v13-coordinates-route:

/api/1.3/coordinates
++++++++++++++++++++

**GET /api/1.3/coordinates**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``id``          | no       | Filter Coordinates by ID.                         |
  +-----------------+----------+---------------------------------------------------+
  | ``name``        | no       | Filter Coordinates by name.                       |
  +-----------------+----------+---------------------------------------------------+

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  | ``id``                            | int    | Local unique identifier for the Coordinate                               |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``lastUpdated``                   | string | The Time / Date this entry was last updated                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``latitude``                      | float  | Latitude of the Coordinate                                               |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``longitude``                     | float  | Longitude of the Coordinate                                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``name``                          | string | The name of the Coordinate                                               |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": 21,
           "lastUpdated": "2012-09-25 20:27:28",
           "latitude": 0,
           "longitude": 0,
           "name": "dc-chicago"
        },
        {
           "id": 22,
           "lastUpdated": "2012-09-25 20:27:28",
           "latitude": 0,
           "longitude": 0,
           "name": "dc-chicago-1"
        }
     ]
    }

|

**POST /api/1.3/coordinates**

  Create Coordinate.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``name``                        | yes      | The name of the Coordinate entry                                  |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``latitude``                    | no       | Latitude of the Coordinate                                        |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``longitude``                   | no       | Longitude of the Coordinate                                       |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "name": "my_coordinate",
        "latitude": 1.2,
        "longitude": 4.5
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``id``                             | int    | The id of the Coordinate                                          |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``name``                           | string | The name of the Coordinate                                        |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``latitude``                       | float  | Latitude of the Coordinate                                        |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``longitude``                      | float  | Longitude of the Coordinate                                       |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``lastUpdated``                    | string | The Time / Date this entry was last updated                       |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``alerts``                         | array  | A collection of alert messages.                                   |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``>level``                         | string | Success, info, warning or error.                                  |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``>text``                          | string | Alert message.                                                    |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "coordinate was created"
                  }
          ],
        "response": {
            'longitude' : 4.5,
            'lastUpdated' : '2016-01-25 13:55:30',
            'name' : 'my_coordinate',
            'latitude' : 1.2,
            'id' : 1
        }
    }
   
|

**PUT /api/1.3/coordinates**

  Update coordinate.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Query Parameters**

  +------+----------+------------------------------------+
  | Name | Required | Description                        |
  +======+==========+====================================+
  | id   | yes      | The id of the coordinate to edit.  |
  +------+----------+------------------------------------+

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``id``                          | yes      | The id of the Coordinate                                          |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``name``                        | yes      | The name of the Coordinate entry                                  |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``latitude``                    | no       | Latitude of the Coordinate                                        |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``longitude``                   | no       | Longitude of the Coordinate                                       |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "id": 1,
        "name": "my_coordinate",
        "latitude": 12,
        "longitude": 45
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``id``                             | int    | The id of the Coordinate                                          |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``name``                           | string | The name of the Coordinate                                        |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``latitude``                       | float  | Latitude of the Coordinate                                        |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``longitude``                      | float  | Longitude of the Coordinate                                       |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``lastUpdated``                    | string | The Time / Date this entry was last updated                       |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``alerts``                         | array  | A collection of alert messages.                                   |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``>level``                         | string | Success, info, warning or error.                                  |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``>text``                          | string | Alert message.                                                    |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "coordinate was updated"
                  }
          ],
        "response": {
            'longitude' : 45,
            'lastUpdated' : '2016-01-25 13:55:30',
            'name' : 'my_coordinate',
            'latitude' : 12,
            'id' : 1
        }
    }

|

**DELETE /api/1.3/coordinates**

  Delete coordinate.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Query Parameters**

  +------+----------+--------------------------------------+
  | Name | Required | Description                          |
  +======+==========+======================================+
  | id   | yes      | The id of the coordinate to delete.  |
  +------+----------+--------------------------------------+
  
  **Response Properties**

  +-------------+--------+----------------------------------+
  |  Parameter  |  Type  |           Description            |
  +=============+========+==================================+
  | ``alerts``  | array  | A collection of alert messages.  |
  +-------------+--------+----------------------------------+
  | ``>level``  | string | Success, info, warning or error. |
  +-------------+--------+----------------------------------+
  | ``>text``   | string | Alert message.                   |
  +-------------+--------+----------------------------------+

  **Response Example** ::

    {
          "alerts": [
                    {
                            "level": "success",
                            "text": "coordinate was deleted"
                    }
            ]
    }

|

