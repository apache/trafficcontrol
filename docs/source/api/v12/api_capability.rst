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

.. _to-api-v12-api_capability:

API-Capabilities
================

.. _to-api-v12-api-capability-route:

/api/1.2/api_capabilities
+++++++++++++++++++++++++

**GET /api/1.2/api_capabilities**

  Get all API-capability mappings.

  Authentication Required: Yes

  Role(s) Required: None

  **Query Parameters**

  +----------------+----------+--------+------------------------------------+
  |    Name        | Required | Type   |         Description                |
  +================+==========+========+====================================+
  | ``capability`` |   no     | string | Capability name.                   |
  +----------------+----------+--------+------------------------------------+

  **Response Properties**

  +-------------------+--------+--------------------------------------------------+
  |    Parameter      |  Type  |                   Description                    |
  +===================+========+==================================================+
  | ``id``            | int    | Mapping id.                                      |
  +-------------------+--------+--------------------------------------------------+
  | ``httpMethod``    | enum   | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
  +-------------------+--------+--------------------------------------------------+
  | ``httpRoute``     | string | API route.                                       |
  +-------------------+--------+--------------------------------------------------+
  | ``capability``    | string | Capability name.                                 |
  +-------------------+--------+--------------------------------------------------+
  | ``lastUpdated``   | string |                                                  |
  +-------------------+--------+--------------------------------------------------+

  **Response Example** ::

    {
     "response": [
           {
              "id": "6",
              "httpMethod": "GET",
              "httpRoute": "/api/*/asns",
              "capability": "asn-read",
              "lastUpdated": "2017-04-02 08:22:43"
           },
           {
              "id": "7",
              "httpMethod": "GET",
              "httpRoute": "/api/*/asns/*",
              "capability": "asn-read",
              "lastUpdated": "2017-04-02 08:22:43"
           }
        ]
    }

|

**GET /api/1.2/api_capabilities/:id**

  Get an API-capability mapping by id.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-------------+----------+-------+-------------------------------------+
  |    Name     | Required |  Type |         Description                 |
  +=============+==========+=======+=====================================+
  |   ``id``    |   yes    | int   | Mapping id.                         |
  +-------------+----------+-------+-------------------------------------+

  **Response Properties**

  +-------------------+--------+--------------------------------------------------+
  |    Parameter      |  Type  |                   Description                    |
  +===================+========+==================================================+
  | ``id``            | int    | Mapping id.                                      |
  +-------------------+--------+--------------------------------------------------+
  | ``httpMethod``    | enum   | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
  +-------------------+--------+--------------------------------------------------+
  | ``httpRoute``     | string | API route.                                       |
  +-------------------+--------+--------------------------------------------------+
  | ``capability``    | string | Capability name.                                 |
  +-------------------+--------+--------------------------------------------------+
  | ``lastUpdated``   | string |                                                  |
  +-------------------+--------+--------------------------------------------------+

  **Response Example** ::

    {
     "response": [
           {
              "id": "6",
              "httpMethod": "GET",
              "httpRoute": "/api/*/asns",
              "capability": "asn-read",
              "lastUpdated": "2017-04-02 08:22:43"
           }
        ]
    }

|

**POST /api/1.2/api_capabilities**

  Create an API-capability mapping.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Properties**

  +----------------+----------+--------+--------------------------------------------------+
  |    Name        | Required | Type   |                Description                       |
  +================+==========+========+==================================================+
  | ``httpMethod`` | yes      | enum   | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
  +----------------+----------+--------+--------------------------------------------------+
  | ``httpRoute``  | yes      | string | API route.                                       |
  +----------------+----------+--------+--------------------------------------------------+
  | ``capability`` | yes      | string | Capability name                                  |
  +----------------+----------+--------+--------------------------------------------------+

  **Request Example** ::

    {
        "httpMethod": "POST",
        "httpRoute": "/api/*/cdns",
        "capability": "cdn-write"
    }

  **Response Properties**

  +--------------------+--------+--------------------------------------------------+
  |    Parameter       |  Type  |                   Description                    |
  +====================+========+==================================================+
  | ``response``       |  hash  | The details of the creation, if success.         |
  +--------------------+--------+--------------------------------------------------+
  | ``>id``            | int    | Mapping id.                                      |
  +--------------------+--------+--------------------------------------------------+
  | ``>httpMethod``    | enum   | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
  +--------------------+--------+--------------------------------------------------+
  | ``>httpRoute``     | string | API route.                                       |
  +--------------------+--------+--------------------------------------------------+
  | ``>capability``    | string | Capability name                                  |
  +--------------------+--------+--------------------------------------------------+
  | ``>lastUpdated``   | string |                                                  |
  +--------------------+--------+--------------------------------------------------+
  | ``alerts``         | array  | A collection of alert messages.                  |
  +--------------------+--------+--------------------------------------------------+
  | ``>level``         | string | Success, info, warning or error.                 |
  +--------------------+--------+--------------------------------------------------+
  | ``>text``          | string | Alert message.                                   |
  +--------------------+--------+--------------------------------------------------+


  **Response Example** ::

    {
        "response":{
              "id": "6",
              "httpMethod": "POST",
              "httpRoute": "/api/*/cdns",
              "capability": "cdn-write",
              "lastUpdated": "2017-04-02 08:22:43"
        },
        "alerts":[
            {
                "level": "success",
                "text": "API-capability mapping was created."
            }
        ]
    }

|

**PUT /api/1.2/api_capabilities/{:id}**

  Edit an API-capability mapping.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-------------------+----------+--------+---------------------------------------+
  | Name              | Required | Type   |           Description                 |
  +===================+==========+========+=======================================+
  |   ``id``          |   yes    | string | Mapping id.                           |
  +-------------------+----------+--------+---------------------------------------+

  **Request Properties**

  +-------------------+--------+--------------------------------------------------+
  |    Parameter      |  Type  |                   Description                    |
  +===================+========+==================================================+
  | ``httpMethod``    | enum   | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
  +-------------------+--------+--------------------------------------------------+
  | ``httpRoute``     | string | API route.                                       |
  +-------------------+--------+--------------------------------------------------+
  | ``capability``    | string | Capability name                                  |
  +-------------------+--------+--------------------------------------------------+


  **Request Example** ::

    {
        "httpMethod": "GET",
        "httpRoute": "/api/*/cdns",
        "capability": "cdn-read"
    }

  **Response Properties**

  +--------------------+--------+--------------------------------------------------+
  |    Parameter       |  Type  |                   Description                    |
  +====================+========+==================================================+
  | ``response``       |  hash  | The details of the creation, if success.         |
  +--------------------+--------+--------------------------------------------------+
  | ``>id``            | int    | Mapping id.                                      |
  +--------------------+--------+--------------------------------------------------+
  | ``>httpMethod``    | enum   | One of: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'. |
  +--------------------+--------+--------------------------------------------------+
  | ``>httpRoute``     | string | API route.                                       |
  +--------------------+--------+--------------------------------------------------+
  | ``>capability``    | string | Capability name                                  |
  +--------------------+--------+--------------------------------------------------+
  | ``>lastUpdated``   | string |                                                  |
  +--------------------+--------+--------------------------------------------------+
  | ``alerts``         | array  | A collection of alert messages.                  |
  +--------------------+--------+--------------------------------------------------+
  | ``>level``         | string | Success, info, warning or error.                 |
  +--------------------+--------+--------------------------------------------------+
  | ``>text``          | string | Alert message.                                   |
  +--------------------+--------+--------------------------------------------------+

  **Response Example** ::

    {
        "response":{
              "id": "6",
              "httpMethod": "GET",
              "httpRoute": "/api/*/cdns",
              "capability": "cdn-read",
              "lastUpdated": "2017-04-02 08:22:43"
        },
        "alerts":[
            {
                "level": "success",
                "text": "API-capability mapping was updated."
            }
        ]
    }

|

**DELETE /api/1.2/api_capabilities/{:id}**

  Delete a capability.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-------------------+----------+--------+---------------------------------------+
  | Name              | Required | Type   |           Description                 |
  +===================+==========+========+=======================================+
  |   ``id``          |   yes    | string | Mapping id.                           |
  +-------------------+----------+--------+---------------------------------------+

  **Response Properties**

  +-----------------+----------+------------------------------------------------+
  |  Parameter      |  Type    |           Description                          |
  +=================+==========+================================================+
  |  ``alerts``     |  array   |  A collection of alert messages.               |
  +-----------------+----------+------------------------------------------------+
  |  ``>level``     |  string  |  success, info, warning or error.              |
  +-----------------+----------+------------------------------------------------+
  |  ``>text``      |  string  |  Alert message.                                |
  +-----------------+----------+------------------------------------------------+

  **Response Example** ::

    {
          "alerts": [
                    {
                            "level": "success",
                            "text": "API-capability mapping deleted."
                    }
            ],
    }
