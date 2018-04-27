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


.. _to-api-v12-asn:

ASN
===

.. _to-api-v12-asns-route:

/api/1.2/asns
+++++++++++++

**GET /api/1.2/asns**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +---------------------+----------+---------------------------------------------+
  |   Name              | Required |                Description                  |
  +=====================+==========+=============================================+
  |   ``cachegroup``    |    no    | Filter ASNs by cache group ID               |
  +---------------------+----------+---------------------------------------------+

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``lastUpdated``  | string | The Time / Date this server entry was last updated                      |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``id``           | string | Local unique identifier for the ASN                                     |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``asn``          | string | Autonomous System Numbers per APNIC for identifying a service provider. |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``cachegroup``   | string | Related cachegroup name                                                 |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``cachegroupId`` | string | Related cachegroup id                                                   |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "lastUpdated": "2012-09-17 21:41:22",
          "id": "27",
          "asn": "7015",
          "cachegroup": "us-ma-woburn",
          "cachegroupId": "27",
        },
        {
          "lastUpdated": "2012-09-17 21:41:22",
          "id": "28",
          "asn": "7016",
          "cachegroup": "us-pa-pittsburgh",
          "cachegroupId": "13"
        }
      ]
    }

|

**GET /api/1.2/asns/:id**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------+----------+---------------------------------------------+
  |   Name    | Required |                Description                  |
  +===========+==========+=============================================+
  |   ``id``  |   yes    | ASN id.                                     |
  +-----------+----------+---------------------------------------------+

  **Response Properties**

  +------------------+--------+-------------------------------------------------------------------------+
  |    Parameter     |  Type  |                               Description                               |
  +==================+========+=========================================================================+
  | ``lastUpdated``  | string | The Time / Date this server entry was last updated                      |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``id``           | string | Local unique identifier for the ASN                                     |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``asn``          | string | Autonomous System Numbers per APNIC for identifying a service provider. |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``cachegroup``   | string | Related cachegroup name                                                 |
  +------------------+--------+-------------------------------------------------------------------------+
  | ``cachegroupId`` | string | Related cachegroup id                                                   |
  +------------------+--------+-------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "lastUpdated": "2012-09-17 21:41:22",
          "id": "28",
          "asn": "7016",
          "cachegroup": "us-pa-pittsburgh",
          "cachegroupId": "13"
        }
      ]
    }

|

**PUT /api/1.2/asns/{:id}**

  Allows user to edit an ASN.

  Authentication Required: Yes

  Role(s) Required:  admin or oper

  **Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              |   Type   |                 Description                    |
  +===================+==========+================================================+
  | ``id``            | int      | ASN id.                                        |
  +-------------------+----------+------------------------------------------------+

  **Request Properties**

  +-------------------+--------+-------------------------------------------------+
  |    Parameter      |  Type  |                   Description                   |
  +===================+========+=================================================+
  | ``asn``           | string | ASN                                             |
  +-------------------+--------+-------------------------------------------------+
  | ``cachegroupId``  | string | The cachegroup the ASN belongs to               |
  +-------------------+--------+-------------------------------------------------+


  **Request Example** ::

    {
        "asn": "99",
        "cachegroupId": "177"
    }

  **Response Properties**

  +--------------------+--------+-------------------------------------------------+
  |    Parameter       |  Type  |                   Description                   |
  +====================+========+=================================================+
  | ``response``       |  hash  | The details of the update, if success.          |
  +--------------------+--------+-------------------------------------------------+
  | ``>name``          | string | CDN name.                                       |
  +--------------------+--------+-------------------------------------------------+
  | ``>id``            |  int   | CDN id.                                         |
  +--------------------+--------+-------------------------------------------------+
  | ``>dnssecEnabled`` | string | Whether dnssec is enabled.                      |
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
          "lastUpdated": "2012-09-17 21:41:22",
          "id": "28",
          "asn": "99",
          "cachegroup": "us-pa-pittsburgh",
          "cachegroupId": "177"
        },
        "alerts":[
            {
                "level": "success",
                "text": "cdn was updated."
            }
        ]
    }

|


