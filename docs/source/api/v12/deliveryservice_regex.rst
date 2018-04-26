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


.. _to-api-v12-ds-regexes:

Delivery Service Regexes
========================

.. _to-api-v12-ds-regexes-route:


**GET /api/1.2/deliveryservices_regexes**

  Retrieves regexes for all delivery services.

  Authentication Required: Yes

  Role(s) Required: Admin or Oper

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``dsName``       | array  | Delivery service name.                                                  |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``regexes``      | array  | An array of regexes for the delivery service.                           |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``>type``        | string | The regex type.                                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``>pattern``     | string | The regex pattern.                                                      |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``>setNumber``   | string | The order in which the regex is evaluated.                              |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "dsName": "foo-bar",
          "regexes": [
            {
              "type": "HOST_REGEXP",
              "pattern": ".*\.foo-bar\..*",
              "setNumber": 0
            },
            {
              "type": "HOST_REGEXP",
              "pattern": "foo.bar.com",
              "setNumber": 1
            }
			    ]
		    },
		    { ... }
      ]
    }

|

**GET /api/1.2/deliveryservices/{:dsId}/regexes**

  Retrieves regexes for a specific delivery service.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  |  ``dsId`` |   yes    | Delivery service id.                        |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``id``           | string | Delivery service regex ID.                                              |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``type``         | string | Delivery service regex type ID.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``typeName``     | string | Delivery service regex type name.                                       |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``pattern``      | string | Delivery service regex pattern.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``setNumber``    | string | The order in which the regex is evaluated for the delivery service.     |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": [
        {
          "id": 852,
          "type": 18,
          "typeName": "HOST_REGEXP",
          "pattern": ".*\.foo-bar\..*",
          "setNumber": 0
        },
        {
          "id": 853,
          "type": 18,
          "typeName": "HOST_REGEXP",
          "pattern": "foo.bar.com",
          "setNumber": 1
        }
      ]
    }

|

**GET /api/1.2/deliveryservices/{:dsId}/regexes/{:id}**

  Retrieves a regex for a specific delivery service.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  | ``dsId``  |   yes    | Delivery service id.                        |
  +-----------+----------+---------------------------------------------+
  | ``id``    |   yes    | Delivery service regex id.                  |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``id``           | string | Delivery service regex ID.                                              |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``type``         | string | Delivery service regex type ID.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``typeName``     | string | Delivery service regex type name.                                       |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``pattern``      | string | Delivery service regex pattern.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``setNumber``    | string | The order in which the regex is evaluated for the delivery service.     |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "id": 852,
          "type": 18,
          "typeName": "HOST_REGEXP",
          "pattern": ".*\.foo-bar\..*",
          "setNumber": 0
        }
      ]
    }

|

**POST /api/1.2/deliveryservices/{:dsId}/regexes**

  Create a regex for a delivery service.

  Authentication Required: Yes

  Role(s) Required: Admin or Oper

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  | ``dsId``  |   yes    | Delivery service id.                        |
  +-----------+----------+---------------------------------------------+

  **Request Properties**

  +---------------+----------+---------------------------------------------+
  |  Parameter    | Required |                Description                  |
  +===============+==========+=============================================+
  | ``pattern``   |   yes    | Regex pattern.                              |
  +---------------+----------+---------------------------------------------+
  | ``type``      |   yes    | Regex type ID.                              |
  +---------------+----------+---------------------------------------------+
  | ``setNumber`` |   yes    | Regex type ID.                              |
  +---------------+----------+---------------------------------------------+

  **Request Example** ::

    {
        "pattern": ".*\.foo-bar\..*"
        "type": 18
        "setNumber": 0
    }

|

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``id``           | string | Delivery service regex ID.                                              |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``type``         | string | Delivery service regex type ID.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``typeName``     | string | Delivery service regex type name.                                       |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``pattern``      | string | Delivery service regex pattern.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``setNumber``    | string | The order in which the regex is evaluated for the delivery service.     |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response":{
        "id": 852,
        "type": 18,
        "typeName": "HOST_REGEXP",
        "pattern": ".*\.foo-bar\..*",
        "setNumber": 0
      },
      "alerts":[
        {
          "level": "success",
          "text": "Delivery service regex creation was successful."
        }
      ]
    }

|

**PUT /api/1.2/deliveryservices/{:dsId}/regexes/{:id}**

  Update a regex for a delivery service.

  Authentication Required: Yes

  Role(s) Required: Admin or Oper

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  | ``dsId``  |   yes    | Delivery service id.                        |
  +-----------+----------+---------------------------------------------+
  | ``id``    |   yes    | Delivery service regex id.                  |
  +-----------+----------+---------------------------------------------+

  **Request Properties**

  +---------------+----------+---------------------------------------------+
  |  Parameter    | Required |                Description                  |
  +===============+==========+=============================================+
  | ``pattern``   |   yes    | Regex pattern.                              |
  +---------------+----------+---------------------------------------------+
  | ``type``      |   yes    | Regex type ID.                              |
  +---------------+----------+---------------------------------------------+
  | ``setNumber`` |   yes    | Regex type ID.                              |
  +---------------+----------+---------------------------------------------+

  **Request Example** ::

    {
        "pattern": ".*\.foo-bar\..*"
        "type": 18
        "setNumber": 0
    }

|

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``id``           | string | Delivery service regex ID.                                              |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``type``         | string | Delivery service regex type ID.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``typeName``     | string | Delivery service regex type name.                                       |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``pattern``      | string | Delivery service regex pattern.                                         |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``setNumber``    | string | The order in which the regex is evaluated for the delivery service.     |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response":{
        "id": 852,
        "type": 18,
        "typeName": "HOST_REGEXP",
        "pattern": ".*\.foo-bar\..*",
        "setNumber": 0
      },
      "alerts":[
        {
          "level": "success",
          "text": "Delivery service regex update was successful."
        }
      ]
    }

|

**DELETE /api/1.2/deliveryservices/{:dsId}/regexes/{:id}**

  Delete delivery service regex.

  Authentication Required: Yes

  Role(s) Required: Admin or Oper

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  | ``dsId``  |   yes    | Delivery service id.                        |
  +-----------+----------+---------------------------------------------+
  | ``id``    |   yes    | Delivery service regex id.                  |
  +-----------+----------+---------------------------------------------+

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
                            "text": "Delivery service regex delete was successful."
                    }
            ],
    }

|




