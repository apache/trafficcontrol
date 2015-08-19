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
  | ``typeName``           | string | The name of the type of Cache Group entry                                |
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


**GET /api/1.2/cachegroupparameters.json**

  Authentication Required: Yes
  
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

