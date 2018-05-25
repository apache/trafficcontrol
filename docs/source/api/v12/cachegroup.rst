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

.. _to-api-v12-cachegroup:

Cache Group
===========

.. _to-api-v12-cachegroups-route:

/api/1.2/cachegroups
++++++++++++++++++++

**GET /api/1.2/cachegroups**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``type``        | no       | Filter cache groups by Type ID.                   |
  +-----------------+----------+---------------------------------------------------+

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  | ``id``                            | string | Local unique identifier for the Cache Group                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``lastUpdated``                   | string | The Time / Date this entry was last updated                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``latitude``                      | string | Latitude for the Cache Group                                             |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``longitude``                     | string | Longitude for the Cache Group                                            |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``name``                          | string | The name of the Cache Group entry                                        |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``parentCachegroupId``            | string | Parent cachegroup ID.                                                    |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``parentCachegroupName``          | string | Parent cachegroup name.                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``secondaryParentCachegroupId``   | string | Secondary parent cachegroup ID.                                          |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``secondaryParentCachegroupName`` | string | Secondary parent cachegroup name.                                        |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``shortName``                     | string | Abbreviation of the Cache Group Name                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``typeId``                        | string | Unique identifier for the 'Type' of Cache Group entry                    |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``typeName``                      | string | The name of the type of Cache Group entry                                |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``fallbackToClosest``             | bool   | Behaviour during non-availability/ failure of configured fallbacks       |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "21",
           "lastUpdated": "2012-09-25 20:27:28",
           "latitude": "0",
           "longitude": "0",
           "name": "dc-chicago",
           "parentCachegroupId": null,
           "parentCachegroupName": null,
           "secondaryParentCachegroupId": null,
           "secondaryParentCachegroupName": null,
           "shortName": "dcchi",
           "typeName": "MID_LOC",
           "typeId": "4",
           "fallbackToClosest":true
        },
        {
           "id": "22",
           "lastUpdated": "2012-09-25 20:27:28",
           "latitude": "0",
           "longitude": "0",
           "name": "dc-chicago-1",
           "parentCachegroupId": null,
           "parentCachegroupName": null,
           "secondaryParentCachegroupId": null,
           "secondaryParentCachegroupName": null,
           "shortName": "dcchi",
           "typeName": "MID_LOC",
           "typeId": "4",
           "fallbackToClosest":false
        }
     ],
    }

|

**GET /api/1.2/cachegroups/trimmed**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

      {
       "response": [
          {
             "name": "dc-chicago"
          },
          {
             "name": "dc-cmc"
          }
       ],
     }

|

**GET /api/1.2/cachegroups/:id**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  | ``id``                            | string | Local unique identifier for the Cache Group                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``lastUpdated``                   | string | The Time / Date this entry was last updated                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``latitude``                      | string | Latitude for the Cache Group                                             |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``longitude``                     | string | Longitude for the Cache Group                                            |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``name``                          | string | The name of the Cache Group entry                                        |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``parentCachegroupId``            | string | Parent cachegroup ID.                                                    |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``parentCachegroupName``          | string | Parent cachegroup name.                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``secondaryParentCachegroupId``   | string | Secondary parent cachegroup ID.                                          |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``secondaryParentCachegroupName`` | string | Secondary parent cachegroup name.                                        |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``shortName``                     | string | Abbreviation of the Cache Group Name                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``typeId``                        | string | Unique identifier for the 'Type' of Cache Group entry                    |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``typeName``                      | string | The name of the type of Cache Group entry                                |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``fallbackToClosest``             | bool   | Behaviour during non-availability/ failure of configured fallbacks       |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": "21",
           "lastUpdated": "2012-09-25 20:27:28",
           "latitude": "0",
           "longitude": "0",
           "name": "dc-chicago",
           "parentCachegroupId": null,
           "parentCachegroupName": null,
           "secondaryParentCachegroupId": null,
           "secondaryParentCachegroupName": null,
           "shortName": "dcchi",
           "typeName": "MID_LOC",
           "typeId": "4",
           "fallbackToClosest":true
        }
     ],
    }

|

**GET /api/1.2/cachegroups/:id/parameters**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  | ``id``                            |   int  | Local unique identifier for the parameter                                |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``name``                          | string | Name of the parameter                                                    |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``value``                         | string | Value of the parameter                                                   |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``configFile``                    | string | Config file associated with the parameter                                |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``secure``                        |  bool  | Is the parameter value only visible to admin users                       |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``lastUpdated``                   | string | The Time / Date this entry was last updated                              |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
            "id": "1100",
            "name": "cgw.originUrl",
            "value": "http://to-short.g.foo.net/data/",
            "configFile": "foo.config",
            "secure": false,
            "lastUpdated": "2015-08-27 15:11:49"
        },
        { ... }
     ]
    }

|

**GET /api/1.2/cachegroups/:id/unassigned_parameters**

  Retrieves all parameters NOT assigned to the cache group.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +------------------+----------+-----------------------+
  |       Name       | Required | Description           |
  +==================+==========+=======================+
  | ``id``           | yes      | Cache group id        |
  +------------------+----------+-----------------------+

  **Response Properties**

  +------------------+---------+--------------------------------------------------------------------------------+
  |    Parameter     |  Type   |                    Description                                                 |
  +==================+=========+================================================================================+
  | ``last_updated`` | string  | The Time / Date this server entry was last updated                             |
  +------------------+---------+--------------------------------------------------------------------------------+
  | ``secure``       | boolean | When true, the parameter is accessible only by admin users. Defaults to false. |
  +------------------+---------+--------------------------------------------------------------------------------+
  | ``value``        | string  | The parameter value, only visible to admin if secure is true                   |
  +------------------+---------+--------------------------------------------------------------------------------+
  | ``name``         | string  | The parameter name                                                             |
  +------------------+---------+--------------------------------------------------------------------------------+
  | ``config_file``  | string  | The parameter config_file                                                      |
  +------------------+---------+--------------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "last_updated": "2012-09-17 21:41:22",
           "secure": false,
           "value": "0,1,2,3,4,5,6",
           "name": "Drive_Letters",
           "config_file": "storage.config"
        },
        {
           "last_updated": "2012-09-17 21:41:22",
           "secure": true,
           "value": "STRING __HOSTNAME__",
           "name": "CONFIG proxy.config.proxy_name",
           "config_file": "records.config"
        }
     ],
    }

|

**GET /api/1.2/cachegroup/:parameter_id/parameter**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +------------------+----------+-------------+
  |       Name       | Required | Description |
  +==================+==========+=============+
  | ``parameter_id`` | yes      |             |
  +------------------+----------+-------------+

  **Response Properties**

  +-----------------+--------+-------------+
  |    Parameter    |  Type  | Description |
  +=================+========+=============+
  | ``cachegroups`` | array  |             |
  +-----------------+--------+-------------+
  | ``>name``       | string |             |
  +-----------------+--------+-------------+
  | ``>id``         | string |             |
  +-----------------+--------+-------------+

  **Response Example** ::

    {
     "response": {
        "cachegroups": [
           {
              "name": "dc-chicago",
              "id": "21"
           },
           {
              "name": "dc-cmc",
              "id": "22"
           }
        ]
     },
    }

|

**GET /api/1.2/cachegroupparameters**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +--------------------------+--------+-----------------------------------------+
  |        Parameter         |  Type  |               Description               |
  +==========================+========+=========================================+
  | ``cachegroupParameters`` | array  | A collection of cache group parameters. |
  +--------------------------+--------+-----------------------------------------+
  | ``>parameter``           | string |                                         |
  +--------------------------+--------+-----------------------------------------+
  | ``>lastUpdated``         | string |                                         |
  +--------------------------+--------+-----------------------------------------+
  | ``>cachegroup``          | string |                                         |
  +--------------------------+--------+-----------------------------------------+

  **Response Example** ::

    {
     "response": {
        "cachegroupParameters": [
           {
              "parameter": "379",
              "lastUpdated": "2013-08-05 18:49:37",
              "cachegroup": "us-ca-sanjose"
           },
           {
              "parameter": "380",
              "lastUpdated": "2013-08-05 18:49:37",
              "cachegroup": "us-ca-sanjose"
           },
           {
              "parameter": "379",
              "lastUpdated": "2013-08-05 18:49:37",
              "cachegroup": "us-ma-woburn"
           }
        ]
     },
    }

|

**GET /api/1.2/cachegroups/:parameter_id/parameter/available**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +------------------+----------+-------------+
  |       Name       | Required | Description |
  +==================+==========+=============+
  | ``parameter_id`` | yes      |             |
  +------------------+----------+-------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``name``              |        |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                |        |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "name": "dc-chicago",
           "id": "21"
        },
        {
           "name": "dc-cmc",
           "id": "22"
        }
     ],
    }

|

**POST /api/1.2/cachegroups**

  Create cache group.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``name``                        | yes      | The name of the Cache Group entry                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``shortName``                   | yes      | Abbreviation of the Cache Group Name                              |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``latitude``                    | no       | Latitude for the Cache Group                                      |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``longitude``                   | no       | Longitude for the Cache Group                                     |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``parentCachegroup``            | no       | Name of Parent Cache Group entry.                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``secondaryParentCachegroup``   | no       | Name of Secondary Parent Cache Group entry.                       |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``typeId``                      | yes      | The type of Cache Group entry, "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``fallbackToClosest``           | no       | Behaviour on configured fallbacks failure, true / false           |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "name": "cache_group_edge",
        "shortName": "cg_edge",
        "latitude": 12,
        "longitude": 45,
        "parentCachegroup": "cache_group_mid",
        "typeId": 6,
        "fallbackToClosest":true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``id``                             | string | The id of cache group                                             |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``name``                           | string | The name of the Cache Group entry                                 |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``shortName``                      | string | Abbreviation of the Cache Group Name                              |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``latitude``                       | string | Latitude for the Cache Group                                      |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``longitude``                      | string | Longitude for the Cache Group                                     |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``parentCachegroup``               | string | Name of Parent Cache Group entry.                                 |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``parentCachegroupId``             | string | id of Parent Cache Group entry.                                   |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``secondaryParentCachegroup``      | string | Name of Secondary Parent Cache Group entry.                       |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``secondaryParentCachegroupId``    | string | id of Secondary Parent Cache Group entry.                         |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``typeName``                       | string | The type of Cache Group entry, "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``fallbackToClosest``              | bool   | Behaviour during non-availability/failure of configured fallbacks |
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
                          "text": "Cachegroup successfully created: cache_group_edge"
                  }
          ],
        "response": {
            'longitude' : '45',
            'lastUpdated' : '2016-01-25 13:55:30',
            'shortName' : 'cg_edge',
            'name' : 'cache_group_edge',
            'parentCachegroup' : 'cache_group_mid',
            'secondaryParentCachegroup' : null,
            'latitude' : '12',
            'typeName' : 'EDGE_LOC',
            'id' : '104',
            'parentCachegroupId' : '103',
            'secondaryParentCachegroupId' : null,
            'fallbackToClosest':true
        }
    }
   
|

**PUT /api/1.2/cachegroups/{:id}**

  Update cache group.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +------+----------+------------------------------------+
  | Name | Required | Description                        |
  +======+==========+====================================+
  | id   | yes      | The id of the cache group to edit. |
  +------+----------+------------------------------------+

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``name``                        | yes      | The name of the Cache Group entry                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``shortName``                   | yes      | Abbreviation of the Cache Group Name                              |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``latitude``                    | no       | Latitude for the Cache Group                                      |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``longitude``                   | no       | Longitude for the Cache Group                                     |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``parentCachegroup``            | no       | Name of Parent Cache Group entry.                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``secondaryParentCachegroup``   | no       | Name of Secondary Parent Cache Group entry.                       |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``typeName``                    | yes      | The type of Cache Group entry, "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``fallbackToClosest``           | no       | Behaviour on configured fallbacks failure, true / false           |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "name": "cache_group_edge",
        "shortName": "cg_edge",
        "latitude": 12,
        "longitude": 45,
        "parentCachegroup": "cache_group_mid",
        "typeName": "EDGE_LOC",
        "fallbackToClosest":true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``id``                             | string | The id of cache group                                             |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``name``                           | string | The name of the Cache Group entry                                 |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``shortName``                      | string | Abbreviation of the Cache Group Name                              |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``latitude``                       | string | Latitude for the Cache Group                                      |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``longitude``                      | string | Longitude for the Cache Group                                     |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``parentCachegroup``               | string | Name of Parent Cache Group entry.                                 |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``parentCachegroupId``             | string | id of Parent Cache Group entry.                                   |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``secondaryParentCachegroup``      | string | Name of Secondary Parent Cache Group entry.                       |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``secondaryParentCachegroupId``    | string | id of Secondary Parent Cache Group entry.                         |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``typeName``                       | string | The type of Cache Group entry, "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``fallbackToClosest``              | bool   | Behaviour during non-availability/failure of configured fallbacks |
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
                          "text": "Cachegroup was updated: cache_group_edge"
                  }
          ],
        "response": {
            'longitude' : '45',
            'lastUpdated' : '2016-01-25 13:55:30',
            'shortName' : 'cg_edge',
            'name' : 'cache_group_edge',
            'parentCachegroup' : 'cache_group_mid',
            'secondaryParentCachegroup' : null,
            'latitude' : '12',
            'typeName' : 'EDGE_LOC',
            'id' : '104',
            'parentCachegroupId' : '103',
            'secondaryParentCachegroupId' : null,
            'fallbackToClosest':true
        }
    }

|

**DELETE /api/1.2/cachegroups/{:id}**

  Delete cache group. The request to delete a cache group, which has servers or child cache group, will be rejected.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +------+----------+--------------------------------------+
  | Name | Required | Description                          |
  +======+==========+======================================+
  | id   | yes      | The id of the cache group to delete. |
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
                            "text": "Cachegroup was deleted: cache_group_edge"
                    }
            ],
    }

|

**POST /api/1.2/cachegroups/{:id}/queue_update**

  Queue or dequeue updates for all servers assigned to a cache group limited to a specific CDN.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +-----------------+----------+----------------------+
  | Name            | Required | Description          |
  +=================+==========+======================+
  | id              | yes      | the cachegroup id.   |
  +-----------------+----------+----------------------+

  **Request Properties**

  +--------------+---------+-----------------------------------------------+
  | Name         | Type    | Description                                   |
  +==============+=========+===============================================+
  | action       | string  | queue or dequeue                              |
  +--------------+---------+-----------------------------------------------+
  | cdn          | string  | cdn name or cdn ID is required                |
  +--------------+---------+-----------------------------------------------+
  | cdnId        | string  | cdn ID or cdn name is required                |
  +--------------+---------+-----------------------------------------------+

  **Response Properties**

  +-----------------+---------+----------------------------------------------------+
  | Name            | Type    | Description                                        |
  +=================+=========+====================================================+
  | action          | string  | The action processed, queue or dequeue.            |
  +-----------------+---------+----------------------------------------------------+
  | cachegroupId    | integer | cachegroup id                                      |
  +-----------------+---------+----------------------------------------------------+
  | cachegroupName  | string  | cachegroup name                                    |
  +-----------------+---------+----------------------------------------------------+
  | cdn             | string  | cdn name                                           |
  +-----------------+---------+----------------------------------------------------+
  | serverNames     | array   | servers name array in the cachegroup in cdn        |
  +-----------------+---------+----------------------------------------------------+

  **Response Example** ::

    {
      "response": {
            "cachegroupName": "us-il-chicago",
            "action": "queue",
            "serverNames":   [
                "atsec-chi-00",
                "atsec-chi-01",
                "atsec-chi-02",
                "atsec-chi-03",
            ],
            "cachegroupId": "93",
            "cdn": "cdn_number_1",
        }
    }

|

**POST /api/1.2/cachegroups/{:id}/deliveryservices**

  Assign deliveryservices for servers in cachegroup

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +------------------+----------+------------------------------------------------------------------------------+
  |      Name        | Required |           Description                                                        |
  +==================+==========+==============================================================================+
  |      id          |   yes    | The cachegroup id.                                                           |
  +------------------+----------+------------------------------------------------------------------------------+
  
  **Request Properties**

  +------------------+----------+------------------------------------------------------------------------------+
  |    Parameter     |   Type   |           Description                                                        |
  +==================+==========+==============================================================================+
  | deliveryServices |  array   | The Ids of the delivery services to assign to each server in the cachegroup. |
  +------------------+----------+------------------------------------------------------------------------------+

  **Request Example** ::

    {
        "deliveryServices": [ 234, 235 ]
    }

  **Response Properties**

  +--------------------+----------+--------------------------------------------------------+
  |    Parameter       |   Type   |           Description                                  |
  +====================+==========+========================================================+
  | response           |   hash   | The details of the assignment, if success.             |
  +--------------------+----------+--------------------------------------------------------+
  |  >id               |   int    | The cachegroup id.                                     |
  +--------------------+----------+--------------------------------------------------------+
  |  >serverNames      |  array   | The server name array in the cachegroup.               |
  +--------------------+----------+--------------------------------------------------------+
  |  >deliveryServices |  array   | The deliveryservice id array.                          |
  +--------------------+----------+--------------------------------------------------------+
  | alerts             |  array   | A collection of alert messages.                        |
  +--------------------+----------+--------------------------------------------------------+
  |  >level            |  string  | Success, info, warning or error.                       |
  +--------------------+----------+--------------------------------------------------------+
  |  >text             |  string  | Alert message.                                         |
  +--------------------+----------+--------------------------------------------------------+

  **Response Example** ::

    {
      "response": {
          "id": 3,
          "serverNames": [ "atlanta-edge-01", "atlanta-edge-07" ],
          "deliveryServices": [ 234, 235 ]
      }
      "alerts":
      [
          {
              "level": "success",
              "text": "Delivery services successfully assigned to all the servers of cache group 3."
          }
      ],
    }

|

