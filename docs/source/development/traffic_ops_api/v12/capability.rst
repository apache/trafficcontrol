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

.. _to-api-v12-capability:

Capabilities
============

.. _to-api-v12-capability-route:

/api/1.2/capabilities
+++++++++++++++++++++

**GET /api/1.2/capabilities**

  Get all capabilities.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +-------------------+--------+-------------------------------------------------+
  |    Parameter      |  Type  |                   Description                   |
  +===================+========+=================================================+
  | ``name``          | string | Capability name.                                |
  +-------------------+--------+-------------------------------------------------+
  | ``description``   | string | Describing the APIs covered by the capability.  |
  +-------------------+--------+-------------------------------------------------+
  | ``lastUpdated``   | string |                                                 |
  +-------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
           {
              "name": "cdn-read",
              "description": "View CDN configuration",
              "lastUpdated": "2017-04-02 08:22:43"
           },
           {
              "name": "cdn-write",
              "description": "Create, edit or delete CDN configuration",
              "lastUpdated": "2017-04-02 08:22:43"
           }
        ]
    }

|

**GET /api/1.2/capabilities/:name**

  Get a capability by name.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-------------+----------+--------+------------------------------------+
  |    Name     | Required | Type   |          Description               |
  +=============+==========+========+====================================+
  |   ``name``  |   yes    | string | Capability name.                   |
  +-------------+----------+--------+------------------------------------+

  **Response Properties**

  +-------------------+--------+-------------------------------------------------+
  |    Parameter      |  Type  |                   Description                   |
  +===================+========+=================================================+
  | ``name``          | string | Capability name.                                |
  +-------------------+--------+-------------------------------------------------+
  | ``description``   | string | Describing the APIs covered by the capability.  |
  +-------------------+--------+-------------------------------------------------+
  | ``lastUpdated``   | string |                                                 |
  +-------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
           {
              "name": "cdn-read",
              "description": "View CDN configuration",
              "lastUpdated": "2017-04-02 08:22:43"
           }
        ]
    }

|

**POST /api/1.2/capabilities**

  Create a capability.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Parameters**

  +-----------------+----------+--------+-------------------------------------------------+
  |      Name       | Required | Type   |          Description                            |
  +=================+==========+========+=================================================+
  |   ``name``      | yes      | string | Capability name.                                |
  +-----------------+----------+--------+-------------------------------------------------+
  | ``description`` | yes      | string | Describing the APIs covered by the capability.  |
  +-----------------+----------+--------+-------------------------------------------------+

  **Request Example** ::

    {
        "name": "cdn-write",
        "description": "Create, edit or delete CDN configuration"
    }

  **Response Properties**

  +--------------------+--------+-------------------------------------------------+
  |    Parameter       |  Type  |                   Description                   |
  +====================+========+=================================================+
  | ``response``       |  hash  | The details of the creation, if success.        |
  +--------------------+--------+-------------------------------------------------+
  | ``>name``          | string | Capability name.                                |
  +--------------------+--------+-------------------------------------------------+
  | ``>description``   | string | Describing the APIs covered by the capability.  |
  +--------------------+--------+-------------------------------------------------+
  | ``alerts``         | array  | A collection of alert messages.                 |
  +--------------------+--------+-------------------------------------------------+
  | ``>level``         | string | Success, info, warning or error.                |
  +--------------------+--------+-------------------------------------------------+
  | ``>text``          | string | Alert message.                                  |
  +--------------------+--------+-------------------------------------------------+


  **Response Example** ::

    {
        "response":{
            "name": "cdn-write",
            "description": "Create, edit or delete CDN configuration"
        },
        "alerts":[
            {
                "level": "success",
                "text": "Capability was created."
            }
        ]
    }

|

**PUT /api/1.2/capabilities/{:name}**

  Edit a capability.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              |   Type   |                 Description                    |
  +===================+==========+================================================+
  | ``name``          | int      | Capability name.                               |
  +-------------------+----------+------------------------------------------------+

  **Request Properties**

  +-------------------+--------+-------------------------------------------------+
  |    Parameter      |  Type  |                   Description                   |
  +===================+========+=================================================+
  | ``description``   | string | Describing the APIs covered by the capability.  |
  +-------------------+--------+-------------------------------------------------+


  **Request Example** ::

    {
        "description": "View CDN configuration"
    }

  **Response Properties**

  +--------------------+--------+-------------------------------------------------+
  |    Parameter       |  Type  |                   Description                   |
  +====================+========+=================================================+
  | ``response``       |  hash  | The details of the update, if success.          |
  +--------------------+--------+-------------------------------------------------+
  | ``>name``          | string | Capability name.                                |
  +--------------------+--------+-------------------------------------------------+
  | ``>description``   |  int   | Describing the APIs covered by the capability.  |
  +--------------------+--------+-------------------------------------------------+
  | ``alerts``         | array  | A collection of alert messages.                 |
  +--------------------+--------+-------------------------------------------------+
  | ``>level``         | string | Success, info, warning or error.                |
  +--------------------+--------+-------------------------------------------------+
  | ``>text``          | string | Alert message.                                  |
  +--------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
        "response":{
            "name": "cdn-read",
            "description": "View CDN configuration"
        },
        "alerts":[
            {
                "level": "success",
                "text": "Capability was updated."
            }
        ]
    }

|

**DELETE /api/1.2/capabilities/{:name}**

  Delete a capability.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-----------------+----------+------------------------------------------------+
  | Name            | Required | Description                                    |
  +=================+==========+================================================+
  | ``name``        | yes      | Capability name.                               |
  +-----------------+----------+------------------------------------------------+

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
                            "text": "Capability deleted."
                    }
            ],
    }
