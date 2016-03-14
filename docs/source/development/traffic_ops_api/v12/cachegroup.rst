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

.. _to-api-v12-cachegroup:

Cache Group
===========

.. _to-api-v12-cachegroups-route:

/api/1.2/cachegroups
++++++++++++++++++++

**GET /api/1.2/cachegroups.json**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +------------------------+--------+--------------------------------------------------------------------------+
  | Parameter              | Type   | Description                                                              |
  +========================+========+==========================================================================+
  | ``longitude``          | string | Longitude for the Cache Group                                            |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``parentCachegroupId`` | string | Identifier that refers to the 'id' field of different Cache Group entry. |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``lastUpdated``        | string | The Time / Date this entry was last updated                              |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``typeName``           | string | The type name of Cache Group entry                                       |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``name``               | string | The name of the Cache Group entry                                        |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``typeId``             | string | Unique identifier for the 'Type' of Cache Group entry                    |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``latitude``           | string | Latitude for the Cache Group                                             |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``id``                 | string | Local unique identifier for the Cache Group                              |
  +------------------------+--------+--------------------------------------------------------------------------+
  | ``shortName``          | string | Abbreviation of the Cache Group Name                                     |
  +------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "longitude": "0",
           "parentCachegroupId": null,
           "lastUpdated": "2012-09-25 20:27:28",
           "typeName": "MID_LOC",
           "name": "dc-chicago",
           "parentCachegroupName": null,
           "typeId": "4",
           "latitude": "0",
           "id": "21",
           "shortName": "dcchi"
        },
        {
           "longitude": "0",
           "parentCachegroupId": null,
           "lastUpdated": "2012-09-25 20:32:03",
           "typeName": "MID_LOC",
           "name": "dc-cmc",
           "parentCachegroupName": null,
           "typeId": "4",
           "latitude": "0",
           "id": "22",
           "shortName": "dccmc"
        }
     ],
    }

|

**GET /api/1.2/cachegroups/trimmed.json**

  Authentication Required: Yes

  Role(s) Required Required: None
  
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

**GET /api/1.2/cachegroup/:parameter_id/parameter.json**

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

**GET /api/1.2/cachegroupparameters.json**

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
  | ``>last_updated``        | string |                                         |
  +--------------------------+--------+-----------------------------------------+
  | ``>cachegroup``          | string |                                         |
  +--------------------------+--------+-----------------------------------------+

  **Response Example** ::

    {
     "response": {
        "cachegroupParameters": [
           {
              "parameter": "379",
              "last_updated": "2013-08-05 18:49:37",
              "cachegroup": "us-ca-sanjose"
           },
           {
              "parameter": "380",
              "last_updated": "2013-08-05 18:49:37",
              "cachegroup": "us-ca-sanjose"
           },
           {
              "parameter": "379",
              "last_updated": "2013-08-05 18:49:37",
              "cachegroup": "us-ma-woburn"
           }
        ]
     },
    }

|

**GET /api/1.2/cachegroups/:parameter_id/parameter/available.json**

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

  **Request Route Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``name``                        | yes      | The name of the Cache Group entry                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``short_name``                  | yes      | Abbreviation of the Cache Group Name                              |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``latitude``                    | yes      | Latitude for the Cache Group                                      |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``longitude``                   | yes      | Longitude for the Cache Group                                     |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``parent_cachegroup``           | yes      | Name of Parent Cache Group entry.                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``secondary_parent_cachegroup`` | yes      | Name of Secondary Parent Cache Group entry.                       |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``type_name``                   | yes      | The type of Cache Group entry, "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "name": "cache_group_edge",
        "short_name": "cg_edge",
        "latitude": "123",
        "longitude": "456",
        "parent_cachegroup": "cache_group_mid",
        "type_name": "EDGE_LOC"
    }

  **Response Properties**

  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                                                  |
  +====================================+========+==============================================================================================+
  | ``id``                             | string | The id of cache group                                                                        |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``name``                           | string | The name of the Cache Group entry                                                            |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``short_name``                     | string | Abbreviation of the Cache Group Name                                                         |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``latitude``                       | string | Latitude for the Cache Group                                                                 |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``longitude``                      | string | Longitude for the Cache Group                                                                |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``parent_cachegroup``              | string | Name of Parent Cache Group entry.                                                            |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``parent_cachegroup_id``           | string | id of Parent Cache Group entry.                                                              |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``secondary_parent_cachegroup``    | string | Name of Secondary Parent Cache Group entry.                                                  |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``secondary_parent_cachegroup_id`` | string | id of Secondary Parent Cache Group entry.                                                    |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+
  | ``type``                           | string | The id of the type of Cache Group entry, its name must be "EDGE_LOC", "MID_LOC" or "ORG_LOC" |
  +------------------------------------+--------+----------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
        "response": {
            'longitude' : '456',
            'last_updated' : '2016-01-25 13:55:30',
            'short_name' : 'cg_edge',
            'name' : 'cache_group_edge',
            'parent_cachegroup' : 'cache_group_mid',
            'secondary_parent_cachegroup' : null,
            'latitude' : '123',
            'type' : '6',
            'id' : '104',
            'parent_cachegroup_id' : '103',
            'secondary_parent_cachegroup_id' : null
        }
    }
   
|
