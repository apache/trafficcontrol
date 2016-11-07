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

.. _to-api-v11-cachegroup:

Cache Group
===========

.. _to-api-v11-cachegroups-route:

/api/1.1/cachegroups
++++++++++++++++++++

**GET /api/1.1/cachegroups**

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
           "typeId": "4"
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
           "typeId": "4"
        }
     ],
    }

|

**GET /api/1.1/cachegroups/trimmed**

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

**GET /api/1.1/cachegroups/:id**

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
           "typeId": "4"
        }
     ],
    }

|


**GET /api/1.1/cachegroup/:parameter_id/parameter**

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


**GET /api/1.1/cachegroupparameters**

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

**GET /api/1.1/cachegroups/:parameter_id/parameter/available**

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

