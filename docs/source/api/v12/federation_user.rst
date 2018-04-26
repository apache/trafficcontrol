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

.. _to-api-v12-federation-user:

Federation User
===============

.. _to-api-v12-federation-user-route:

/api/1.2/federations/:id/users
++++++++++++++++++++++++++++++

**GET /api/1.2/federations/:id/users**

  Retrieves users assigned to a federation.

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
  | ``company``         | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``id``              |  int   |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``username``        | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``role``            | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``email``           | string |                                                    |
  +---------------------+--------+----------------------------------------------------+
  | ``fullName``        | string |                                                    |
  +---------------------+--------+----------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
			"id": 41
			"username": "booya",
			"company": "XYZ Corporation",
			"role": "federation",
			"email": "booya@fooya.com",
			"fullName": "Booya Fooya"
		}
      ]
    }

|

**POST /api/1.2/federations/:id/users**

  Create one or more federation / user assignments.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Parameters**

  +---------------------------------+----------+-------------------------------------------------------------------+
  | Name                            | Required | Description                                                       |
  +=================================+==========+===================================================================+
  | ``userIds``                     | yes      | An array of user IDs.                                             |
  +---------------------------------+----------+-------------------------------------------------------------------+
  | ``replace``                     | no       | Replace existing fed/user assignments? (true|false)               |
  +---------------------------------+----------+-------------------------------------------------------------------+

  **Request Example** ::

    {
        "userIds": [ 2, 3, 4, 5, 6 ],
        "replace": true
    }

  **Response Properties**

  +------------------------------------+--------+-------------------------------------------------------------------+
  | Parameter                          | Type   | Description                                                       |
  +====================================+========+===================================================================+
  | ``userIds``                        | array  | An array of user IDs.                                             |
  +------------------------------------+--------+-------------------------------------------------------------------+
  | ``replace``                        | array  | Existing fed/user assignments replaced? (true|false).             |
  +------------------------------------+--------+-------------------------------------------------------------------+

  **Response Example** ::

    {
        "alerts": [
                  {
                          "level": "success",
                          "text": "5 user(s) were assigned to the cname. federation"
                  }
          ],
        "response": {
            "userIds" : [ 2, 3, 4, 5, 6 ],
            "replace" : true
        }
    }

|

**DELETE /api/1.2/federations/:id/users/:id**

  Removes a user from a federation.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Route Parameters**

  +-----------------+----------+---------------------------------------------------+
  | Name            | Required | Description                                       |
  +=================+==========+===================================================+
  | ``federation``  | yes      | Federation ID.                                    |
  +-----------------+----------+---------------------------------------------------+
  | ``user``        | yes      | User ID.                                          |
  +-----------------+----------+---------------------------------------------------+

   **Response Example** ::

    {
           "alerts": [
                     {
                             "level": "success",
                             "text": "Removed user [ bobmack ] from federation [ cname1. ]"
                     }
             ],
    }

|



