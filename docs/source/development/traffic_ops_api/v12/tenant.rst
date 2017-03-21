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

.. _to-api-v12-tenant:

Tenants
=======

.. _to-api-v12-tenant-route:

/api/1.2/tenants
++++++++++++++++++

**GET /api/1.2/tenants**

  Get all tenants.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``id``                |  int   | Tenant id                                       |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Tenant name                                     |
  +----------------------+--------+-------------------------------------------------+
  |``active``            |  bool  | Active or inactive                              |
  +----------------------+--------+-------------------------------------------------+
  |``parentId``          |  int   | Parent tenant ID                                |
  +----------------------+--------+-------------------------------------------------+
  |``parentName``        | string | Parent tenant name                              |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": 1
           "name": "root",
           "active": true,
           "parentId": null,
           "parentName": null,
        },
        {
           "id": 2
           "name": "tenant-a",
           "active": true,
           "parentId": 1
           "parentName": "root"
        }
     ]
    }

|


**GET /api/1.2/tenants/:id**

  Get a tenant by ID.

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``id``                |  int   | Tenant id                                       |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Tenant name                                     |
  +----------------------+--------+-------------------------------------------------+
  |``active``            |  bool  | Active or inactive                              |
  +----------------------+--------+-------------------------------------------------+
  |``parentId``          |  int   | Parent tenant ID                                |
  +----------------------+--------+-------------------------------------------------+
  |``parentName``        | string | Parent tenant name                              |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "id": 2
           "name": "tenant-a",
           "active": true,
           "parentId": 1,
           "parentName": "root"
        }
     ]
    }

|


**PUT /api/1.2/tenants/:id**

  Update a tenant.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Route Parameters**

  +-------------------+----------+------------------------------------------------+
  | Name              |   Type   |                 Description                    |
  +===================+==========+================================================+
  | ``id``            |   int    | Tenant id                                      |
  +-------------------+----------+------------------------------------------------+

  **Request Properties**

  +-------------------+----------+--------------------------+
  | Parameter         | Required | Description              |
  +===================+==========+==========================+
  | ``name``          | yes      | The name of the tenant   |
  +-------------------+----------+--------------------------+
  | ``active``        | yes      | True or false            |
  +-------------------+----------+--------------------------+
  | ``parentId``      | yes      | Parent tenant            |
  +-------------------+----------+--------------------------+

  **Request Example** ::

    {
        "name": "my-tenant"
        "active": true
        "parentId": 1
    }

|

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``id``                |  int   | Tenant id                                       |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Tenant name                                     |
  +----------------------+--------+-------------------------------------------------+
  |``active``            |  bool  | Active or inactive                              |
  +----------------------+--------+-------------------------------------------------+
  |``parentId``          |  int   | Parent tenant ID                                |
  +----------------------+--------+-------------------------------------------------+
  |``parentName``        | string | Parent tenant name                              |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

	{
		"response": {
			"id": 2,
			"name": "my-tenant",
			"active": true,
			"parentId": 1,
			"parentName": "root",
			"lastUpdated": "2014-03-18 08:57:39"
		},
		"alerts": [
			{
				"level": "success",
				"text": "Tenant update was successful."
			}
		]
	}

|


**POST /api/1.2/tenants**

  Create a tenant.

  Authentication Required: Yes

  Role(s) Required: admin or oper

  **Request Properties**

  +-------------------+----------+--------------------------+
  | Parameter         | Required | Description              |
  +===================+==========+==========================+
  | ``name``          | yes      | The name of the tenant   |
  +-------------------+----------+--------------------------+
  | ``active``        | no       | Defaults to false        |
  +-------------------+----------+--------------------------+
  | ``parentId``      | yes      | Parent tenant            |
  +-------------------+----------+--------------------------+

  **Request Example** ::

    {
        "name": "your-tenant"
        "parentId": 2
    }

|

  **Response Properties**

  +----------------------+--------+-------------------------------------------------+
  | Parameter            | Type   | Description                                     |
  +======================+========+=================================================+
  |``id``                |  int   | Tenant id                                       |
  +----------------------+--------+-------------------------------------------------+
  |``name``              | string | Tenant name                                     |
  +----------------------+--------+-------------------------------------------------+
  |``active``            |  bool  | Active or inactive                              |
  +----------------------+--------+-------------------------------------------------+
  |``parentId``          |  int   | Parent tenant ID                                |
  +----------------------+--------+-------------------------------------------------+
  |``parentName``        | string | Parent tenant name                              |
  +----------------------+--------+-------------------------------------------------+

  **Response Example** ::

	{
		"response": {
			"id": 2,
			"name": "your-tenant",
			"active": false,
			"parentId": 2,
			"parentName": "my-tenant",
			"lastUpdated": "2014-03-18 08:57:39"
		},
		"alerts": [
			{
				"level": "success",
				"text": "Tenant create was successful."
			}
		]
	}

|
