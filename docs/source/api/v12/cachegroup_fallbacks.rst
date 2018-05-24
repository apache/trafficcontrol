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

.. _to-api-v12-cachegroupfallbacks:

Cache Group Fallback
====================

.. _to-api-v12-cachegroupfallbacks-route:

/api/1.2/cachegroup_fallbacks
++++++++++++++++++++++++++++++

**GET /api/1.2/cachegroup_fallbacks?cacheGroupId={id}**
**GET /api/1.2/cachegroup_fallbacks?fallbackId={id}**
**GET /api/1.2/cachegroup_fallbacks?cacheGroupId={id}&fallbackId={id}**

  Retrieve fallback related configurations for a cache group.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  Query parameter is mandatory. Either one of the parameters must be used. Both can also be used simultaneously.

  +-----------------+---------------------------------------------------------------------------+
  | Name            | Description                                                               |
  +=================+===========================================================================+
  | cacheGroupId    | The id of the cache group whose backup configurations has to be retrieved |
  +-----------------+---------------------------------------------------------------------------+
  | fallbackId      | The id of the fallback cache group associated with a cache group          |
  +-----------------+---------------------------------------------------------------------------+

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  |                                   | array  | parameters array                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupId``                 | int    | Cache group id                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackId``                   | int    | fallback cache group id                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupName``               | string | Cache group name                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackName``                 | string | Fallback cache group  name                                               |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackOrder``                | int    | Ordering list in the list of backups                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
       "response": [
          {
             "cacheGroupId":1,
             "cacheGroupName":"GROUP1",
             "fallbackId":2,
             "fallbackOrder":10,
             "fallbackName":"GROUP2"
          }
       ]
    }

|

**POST /api/1.2/cachegroup_fallbacks**

  Creates fallback configuration for the cache group. New fallbacks can be added only via POST.   

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Parameters**
  The request parameters should be in array format.

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  |                                   | array  | parameters array                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupId``                 | int    | Cache group id                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackId``                   | int    | Fallback cache group id                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackOrder``                | int    | Ordering list in the list of backups                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Request Example** ::

    [
       {
          "cacheGroupId": 1, 
          "fallbackId": 3, 
          "fallbackOrder": 10
       }
    ]

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  |                                   | array  | parameters array                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupId``                 | int    | Cache group id                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackId``                   | int    | fallback cache group id                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupName``               | string | Cache group name                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackName``                 | string | Fallback cache group  name                                               |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackOrder``                | int    | Ordering list in the list of backups                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``alerts``                        | array  | A collection of alert messages.                                          |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>level``                        | string | Success, info, warning or error.                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>text``                         | string | Alert message.                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+


  **Response Example** ::

    {
       "alerts": [
          {
             "level":"success",
             "text":"Backup configuration CREATE for cache group 1 successful."
          }
       ],
       "response": [
          {
             "cacheGroupId":1,
             "cacheGroupName":"GROUP1",
             "fallbackId":3,
             "fallbackName":"GROUP2",
             "fallbackorder":10,
          }
       ]
    }
   
|

**PUT /api/1.2/cachegroup_fallbacks**

  Updates an existing fallback configuration for the cache group. 

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Parameters**
  The request parameters should be in array format.

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  |                                   | array  | parameters array                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupId``                 | int    | Cache group id                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackId``                   | int    | Fallback cache group id                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackOrder``                | int    | Ordering list in the list of backups                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Request Example** ::

    [
       {
          "cacheGroupId": 1, 
          "fallbackId": 3, 
          "fallbackOrder": 10
       }
    ]

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  |                                   | array  | parameters array                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupId``                 | int    | Cache group id                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackId``                   | int    | fallback cache group id                                                  |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>cacheGroupName``               | string | Cache group name                                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackName``                 | string | Fallback cache group  name                                               |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>fallbackOrder``                | int    | Ordering list in the list of backups                                     |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``alerts``                        | array  | A collection of alert messages.                                          |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>level``                        | string | Success, info, warning or error.                                         |
  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | ``>text``                         | string | Alert message.                                                           |
  +-----------------------------------+--------+--------------------------------------------------------------------------+


  **Response Example** ::

    {
       "alerts": [
          {
             "level":"success",
             "text":"Backup configuration UPDATE for cache group 1 successful."
          }
       ],
       "response": [
          {
             "cacheGroupId":1,
             "cacheGroupName":"GROUP1",
             "fallbackId":3,
             "fallbackName":"GROUP2",
             "fallbackorder":10,
          }
       ]
    }
   
|

**DELETE /api/1.2/cachegroup_fallbacks?cacheGroupId={id}**
**DELETE /api/1.2/cachegroup_fallbacks?fallbackId={id}**
**DELETE /api/1.2/cachegroup_fallbacks?fallbackId={id}&cacheGroupId={id}**

  Delete fallback list assigned to the cache group.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Query Parameters**

  Query parameter is mandatory. Either one of the parameters must be used. Both can also be used simultaneously.

  +-----------------+----------+--------------------------------------------------------------------------------------+
  | Name            | Required |     Description                                                                      |
  +=================+==========+======================================================================================+
  | cacheGroupId    | Yes      | The id of the cache group whose backup configurations has to be deleted              |
  +-----------------+----------+--------------------------------------------------------------------------------------+
  | fallbackId      | Yes      | The id of the fallback cachegroup which has to be deleted from the list of fallbacks |
  +-----------------+----------+--------------------------------------------------------------------------------------+
  
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
                            "text": "Backup configuration DELETED"
                    }
            ],
    }

|

