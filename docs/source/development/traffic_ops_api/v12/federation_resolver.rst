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

.. _to-api-v12-federation-resolver:

Federation Resolver
===================

.. _to-api-v12-federation-resolver-route:

/api/1.2/federation_resolvers
+++++++++++++++++++++++++++++

**POST /api/1.2/federation_resolvers**

  Create a federation resolver.

  Authentication Required: Yes

  Role(s) Required: ADMIN

  **Request Properties**

  +-------------------------+----------+------------------------------------------+
  | Parameter               | Required | Description                              |
  +=========================+==========+==========================================+
  | ``ipAddress``           | yes      | IP or CIDR range                         |
  +-------------------------+----------+------------------------------------------+
  | ``typeId``              | yes      | Type Id where useintable=federation      |
  +-------------------------+----------+------------------------------------------+

  **Request Example** ::

    {
        "ipAddress": "2.2.2.2/32",
        "typeId": 245
    }

|

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``id``                | int    |                                                |
  +----------------------+--------+------------------------------------------------+
  |``ipAddress``         | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``type``              | int    |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "Federation resolver created [ IP = 2.2.2.2/32 ] with id: 27"
                  }
          ],
        "response": {
            "id" : 27,
            "ipAddress" : "2.2.2.2/32",
            "typeId" : 245,
        }
    }

|

**DELETE /api/1.2/federation_resolvers/:id**

  Deletes a federation resolver.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``resolver``    | yes      | Federation resolver ID.                           |
  +-----------------+----------+---------------------------------------------------+

   **Response Example** ::

    {
           "alerts": [
                     {
                             "level": "success",
                             "text": "Federation resolver deleted [ IP = 2.2.2.2/32 ] with id: 27"
                     }
             ],
    }

|
