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

.. _to-api-v12-federation:

Federation 
==========

.. _to-api-v12-federation-route:

/api/1.2/federations
++++++++++++++++++++

**GET /api/1.2/federations.json**

  Retrieves a list of federations for a user's delivery service(s).

  Authentication Required: Yes

  Role(s) Required: Federation

  **Response Properties**

  +---------------------+--------+----------------------------------------------------+
  |    Parameter        |  Type  |                   Description                      |
  +=====================+========+====================================================+
  | ``cname``           | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``ttl``             |  int   | Time to live for the cname.                        |
  +---------------------+--------+----------------------------------------------------+
  | ``deliveryService`` | string | Unique string that describes the deliveryservice.  |
  +---------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "mappings": [
            "cname": "cname-01.",
            "ttl": 8865,
          ]
          "deliveryService": "ds-01",
        }
      ]
    }

|

**POST /api/1.2/federations.json**

  Allows a user to add federations for their delivery service(s).

  Authentication Required: Yes

  Role(s) Required: Federation

  **Request Properties**

  +---------------------+--------+----------------------------------------------------+
  |    Parameter        |  Type  |                   Description                      |
  +=====================+========+====================================================+
  | ``deliveryService`` | string | Unique string that describes the deliveryservice.  |
  +---------------------+--------+----------------------------------------------------+
  | ``resolve4``        | array  | Array of IPv4 Addresses.                           |
  +---------------------+--------+----------------------------------------------------+
  | ``resolve6``        | array  | Array of IPv6 Addresses.                           |
  +---------------------+--------+----------------------------------------------------+

  **Request Example** ::

    {
      "federations": [
        {
          "deliveryService": "ccp-omg-01",
          "mappings": {
            "resolve4": [
              "255.255.255.255"
            ],
            "resolve6": [
              "FE80::0202:B3FF:FE1E:8329",
            ]
          }
        }
      ]
    }

|

**DELETE /api/1.2/federations.json**

  Deletes **all** federations associated with a user's delivery service(s).

  Authentication Required: Yes

  Role(s) Required: Federation

|


**PUT /api/1.2/federations.json**

  Deletes **all** federations associated with a user's delivery service(s) then adds the new federations.

  Authentication Required: Yes

  Role(s) Required: Federation

  **Request Properties**

  +---------------------+--------+----------------------------------------------------+
  |    Parameter        |  Type  |                   Description                      |
  +=====================+========+====================================================+
  | ``deliveryService`` | string | Unique string that describes the deliveryservice.  |
  +---------------------+--------+----------------------------------------------------+
  | ``resolve4``        | array  | Array of IPv4 Addresses.                           |
  +---------------------+--------+----------------------------------------------------+
  | ``resolve6``        | array  | Array of IPv6 Addresses.                           |
  +---------------------+--------+----------------------------------------------------+

  **Request Example** ::

    {
      "federations": [
        {
          "deliveryService": "ccp-omg-01",
          "mappings": {
            "resolve4": [
              "255.255.255.255"
            ],
            "resolve6": [
              "FE80::0202:B3FF:FE1E:8329",
            ]
          }
        }
      ]
    }

|
