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

.. _to-api-v13-staticdnsentries:

StaticDNSEntries
================

.. _to-api-v13-static-dns-entry-route:

/api/1.3/staticdnsentries
+++++++++++++++++++++++++

**GET /api/1.3/staticdnsentries**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +-----------------------+----------+-------------------------------------------------------------------------------------------+
  | Name                  | Required | Description                                                                               |
  +=======================+==========+===========================================================================================+
  | ``id``                | no       | Filter StaticDNSEntries by id.                                                            |
  +-----------------------+----------+-------------------------------------------------------------------------------------------+
  | ``address``           | no       | Filter StaticDNSENtries by address.                                                       |
  | ``cachegroup``        | no       | Filter StaticDNSENtries by cachegroup name.                                               |
  | ``cachegroupId``      | no       | Filter StaticDNSENtries by cachegroup id.                                                 |
  | ``deliveryservice``   | no       | Filter StaticDNSENtries by deliveryservice (xml_id).                                      |
  | ``deliveryserviceId`` | no       | Filter StaticDNSENtries by deliveryserviceId.                                             |
  | ``host``              | no       | Filter StaticDNSENtries by host.                                                          |
  | ``ttl``               | no       | Filter StaticDNSENtries by ttl.                                                           |
  | ``type``              | no       | Filter StaticDNSENtries by type (valid types are A_RECORD, AAAA_RECORD and CNAME_RECORD). |
  | ``typeId``            | no       | Filter StaticDNSENtries by typeId.                                                        |
  +-----------------------+----------+-------------------------------------------------------------------------------------------+

  **Response Properties**

  +-----------------------+--------+---------------------------------------------+
  | Parameter             | Type   | Description                                 |
  +=======================+========+=============================================+
  | ``address``           | string | The fully qualified domain name (FQDN)      |
  +-----------------------+--------+---------------------------------------------+
  | ``cachegroup``        | string | The Cachegroup Name                         |
  +-----------------------+--------+---------------------------------------------+
  | ``cachegroupId``      | int    | The Cachegroup id                           |
  +-----------------------+--------+---------------------------------------------+
  | ``deliveryservice``   | string | The DeliveryService Name                    |
  +-----------------------+--------+---------------------------------------------+
  | ``deliveryserviceId`` | int    | The DeliveryService id                      |
  +-----------------------+--------+---------------------------------------------+
  | ``host``              | string | The hostname                                |
  +-----------------------+--------+---------------------------------------------+
  | ``id``                | int    | Local unique identifier                     |
  +-----------------------+--------+---------------------------------------------+
  | ``lastUpdated``       | string | The Time / Date this entry was last updated |
  +-----------------------+--------+---------------------------------------------+
  | ``ttl``               | int    | The Total Time to Live                      |
  +-----------------------+--------+---------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "address": "this.one.is.a.hostname",
           "cachegroup": "cachegroup1",
           "cachegroupId": 18,
           "deliveryservice": "ds1",
           "deliveryserviceId": 28,
           "lastUpdated": "2012-09-25 20:27:28",
           "host": "host1",
           "id": 21,
           "ttl": 30,
           "type": "CNAME_RECORD",
           "typeId": 19
        },
        {
           "address": "this.two.is.a.hostname",
           "cachegroup": "cachegroup2",
           "cachegroupId": 18,
           "deliveryservice": "ds1",
           "deliveryserviceId": 30,
           "lastUpdated": "2012-09-25 20:27:28",
           "host": "host2",
           "id": 22,
           "ttl": 10,
           "type": "A_RECORD",
           "typeId": 18
        }
     ]
    }

|

**POST /api/1.3/staticdnsentries**

  Create a StaticDNSEntry.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Parameters**

  +-----------------------+----------+----------------------------------------+
  | Name                  | Required | Description                            |
  +=======================+==========+========================================+
  | ``address``           | yes      | The fully qualified domain name (FQDN) |
  +-----------------------+----------+----------------------------------------+
  | ``cachegroupId``      | yes      | The Cachegroup id                      |
  +-----------------------+----------+----------------------------------------+
  | ``deliveryserviceId`` | yes      | The DeliveryService id                 |
  +-----------------------+----------+----------------------------------------+
  | ``host``              | yes      | The hostname                           |
  +-----------------------+----------+----------------------------------------+
  | ``ttl``               | yes      | The Total Time to Live                 |
  +-----------------------+----------+----------------------------------------+

  **Request Example** ::

    {
        "address": "this.one.is.a.hostname",
        "cachegroupId": 18,
        "deliveryserviceId": 20,
        "host": 20,
        "ttl": 10
    }

  **Response Properties**

  +-----------------------+--------+----------------------------------------+
  | Parameter             | Type   | Description                            |
  +=======================+========+========================================+
  | ``id``                | int    | The id of the StaticDNSEntry           |
  +-----------------------+--------+----------------------------------------+
  | ``address``           | string | The fully qualified domain name (FQDN) |
  +-----------------------+--------+----------------------------------------+
  | ``cachegroupId``      | int    | The Cachegroup id                      |
  +-----------------------+--------+----------------------------------------+
  | ``deliveryserviceId`` | int    | The DeliveryService id                 |
  +-----------------------+--------+----------------------------------------+
  | ``host``              | string | The hostname                           |
  +-----------------------+--------+----------------------------------------+
  | ``ttl``               | int    | The Total Time to Live                 |
  +-----------------------+--------+----------------------------------------+
  | ``alerts``            | array  | A collection of alert messages.        |
  +-----------------------+--------+----------------------------------------+
  | ``>level``            | string | Success, info, warning or error.       |
  +-----------------------+--------+----------------------------------------+
  | ``>text``             | string | Alert message.                         |
  +-----------------------+--------+----------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "staticdnsentry was created"
                  }
          ],
          "response": {
            "address": "this.one.is.a.hostname",
            "cachegroupId": 18,
            "deliveryserviceId": 20,
            "lastUpdated" : "2016-01-25 13:55:30",
            "host": 20,
            "id" : 1,
            "ttl": 10
        }
    }
   
|

**PUT /api/1.3/staticdnsentries**

  Update staticdnsentries.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Query Parameters**

  +------+----------+---------------------------------------+
  | Name | Required | Description                           |
  +======+==========+=======================================+
  | id   | yes      | The id of the staticdnsentry to edit. |
  +------+----------+---------------------------------------+

  **Request Parameters**

  +-----------------------+----------+----------------------------------------+
  | Name                  | Required | Description                            |
  +=======================+==========+========================================+
  | ``address``           | yes      | The fully qualified domain name (FQDN) |
  +-----------------------+----------+----------------------------------------+
  | ``cachegroupId``      | yes      | The Cachegroup id                      |
  +-----------------------+----------+----------------------------------------+
  | ``deliveryserviceId`` | yes      | The DeliveryService id                 |
  +-----------------------+----------+----------------------------------------+
  | ``host``              | yes      | The hostname                           |
  +-----------------------+----------+----------------------------------------+
  | ``ttl``               | yes      | The Total Time to Live                 |
  +-----------------------+----------+----------------------------------------+

  **Request Example** ::

    {
        "address": "this.one.is.a.hostname",
        "cachegroupId": 18,
        "deliveryserviceId": 20,
        "host": 20,
        "ttl": 10
    }

  **Response Properties**

  +-----------------------+--------+----------------------------------------+
  | Parameter             | Type   | Description                            |
  +=======================+========+========================================+
  | ``id``                | int    | The id of the StaticDNSEntry           |
  +-----------------------+--------+----------------------------------------+
  | ``address``           | string | The fully qualified domain name (FQDN) |
  +-----------------------+--------+----------------------------------------+
  | ``cachegroupId``      | int    | The Cachegroup id                      |
  +-----------------------+--------+----------------------------------------+
  | ``deliveryserviceId`` | int    | The DeliveryService id                 |
  +-----------------------+--------+----------------------------------------+
  | ``host``              | string | The hostname                           |
  +-----------------------+--------+----------------------------------------+
  | ``ttl``               | int    | The Total Time to Live                 |
  +-----------------------+--------+----------------------------------------+
  | ``alerts``            | array  | A collection of alert messages.        |
  +-----------------------+--------+----------------------------------------+
  | ``>level``            | string | Success, info, warning or error.       |
  +-----------------------+--------+----------------------------------------+
  | ``>text``             | string | Alert message.                         |
  +-----------------------+--------+----------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "staticdnsentry was updated"
                  }
          ],
        "response": {
            "address": "this.one.is.a.hostname",
            "cachegroupId": 18,
            "deliveryserviceId": 20,
            "lastUpdated" : "2016-01-25 13:55:30",
            "host": 20,
            "id" : 1,
            "ttl": 10
        }
    }

|

**DELETE /api/1.3/staticdnsentries**

  Delete staticdnsentries.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Query Parameters**

  +------+----------+-----------------------------------------+
  | Name | Required | Description                             |
  +======+==========+=========================================+
  | id   | yes      | The id of the staticdnsentry to delete. |
  +------+----------+-----------------------------------------+
  
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
                            "text": "staticdnsentry was deleted"
                    }
            ]
    }

|

