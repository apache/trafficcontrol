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

.. _to-api-v12-federation-deliveryservice:

Federation Delivery Service
===========================

.. _to-api-v12-federation-deliveryservice-route:

/api/1.2/federations/:id/deliveryservices
+++++++++++++++++++++++++++++++++++++++++

**GET /api/1.2/federations/:id/deliveryservices**

  Retrieves delivery services assigned to a federation.

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
  | ``cdn``             | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``type``            | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``xmlId``           | string |                                                    |
  +---------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
			"id": 41
			"cdn": "cdn1",
			"type": "DNS",
			"xmlId": "booya-12"
		}
      ]
    }

|

**POST /api/1.2/federations/:id/deliveryservices**

  Create one or more federation / delivery service assignments.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``dsIds``                       | yes      | An array of delivery service IDs.                                 |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``replace``                     | no       | Replace existing fed/ds assignments? (true|false)                 |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "dsIds": [ 2, 3, 4, 5, 6 ],
        "replace": true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``dsIds``                          | array  | An array of delivery service IDs.                                 |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``replace``                        | array  | Existing fed/ds assignments replaced? (true|false).               |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "5 delivery service(s) were assigned to the cname. federation"
                  }
          ],
        "response": {
            "dsIds" : [ 2, 3, 4, 5, 6 ],
            "replace" : true
        }
    }

|

**DELETE /api/1.2/federations/:id/deliveryservices/:id**

  Removes a delivery service from a federation.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``federation``  | yes      | Federation ID.                                    |
  +-----------------+----------+---------------------------------------------------+
  | ``ds``          | yes      | Delivery Service ID.                              |
  +-----------------+----------+---------------------------------------------------+

   **Response Example** ::

    {
           "alerts": [
                     {
                             "level": "success",
                             "text": "Removed delivery service [ booya-12 ] from federation [ cname1. ]"
                     }
             ],
    }

|



