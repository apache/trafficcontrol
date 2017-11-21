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

.. _to-api-v12-federation-federationresolver:

Federation Federation Resolver
==============================

.. _to-api-v12-federation-federationresolver-route:

/api/1.2/federations/:id/federation_resolvers
+++++++++++++++++++++++++++++++++++++++++++++

**GET /api/1.2/federations/:id/federation_resolvers**

  Retrieves federation resolvers assigned to a federation.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              |   Type   |                 Description                    |
  +===================+==========+================================================+
  | ``federation``    | string   | Federation ID.                                 |
  +-------------------+----------+------------------------------------------------+

  **Response Properties**

  +---------------------+--------+----------------------------------------------------+
  |    Parameter        |  Type  |                   Description                      |
  +=====================+========+====================================================+
  | ``id``              |  int   |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``ipAddress``       | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``type``            | string |                                                    |
  +---------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
			"id": 41
			"ipAddress": "2.2.2.2/16",
			"type": "RESOLVE4"
		}
      ]
    }

|

**POST /api/1.2/federations/:id/federation_resolvers**

  Create one or more federation / federation resolver assignments.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``fedResolverIds``              | yes      | An array of federation resolver IDs.                              |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``replace``                     | no       | Replace existing fed/ds assignments? (true|false)                 |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "fedResolverIds": [ 2, 3, 4, 5, 6 ],
        "replace": true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``fedResolverIds``                 | array  | An array of federation resolver IDs.                              |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``replace``                        | array  | Existing fed/fed resolver assignments replaced? (true|false).     |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "5 resolvers(s) were assigned to the cname. federation"
                  }
          ],
        "response": {
            "fedResolverIds" : [ 2, 3, 4, 5, 6 ],
            "replace" : true
        }
    }

|
