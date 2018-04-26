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

.. _to-api-v12-topology:

Snapshot CRConfig
=================

.. _to-api-v12-topology-route:

/api/1.2/snapshot/{:cdn_name}
+++++++++++++++++++++++++++++

**GET /api/1.2/cdns/{:cdn_name}/snapshot**

  Retrieves the CURRENT snapshot for a CDN which doesn't necessarily represent the current state of the CDN. The contents of this snapshot are currently used by Traffic Monitor and Traffic Router.

  Authentication Required: Yes

  Role(s) Required: Admin or Operations

  **Request Route Parameters**

  +----------------+----------+---------------------------------------------+
  |   Name         | Required |                Description                  |
  +================+==========+=============================================+
  |  ``cdn_name``  |   yes    | CDN name.                                   |
  +----------------+----------+---------------------------------------------+

  **Response Properties**

  +-----------------------+--------+------------------------------------------------------------------------------+
  |    Parameter          |  Type  |                               Description                                    |
  +=======================+========+==============================================================================+
  | ``config``            |  hash  | General CDN configuration settings.                                          |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``contentRouters``    |  hash  | A list of Traffic Routers.                                                   |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``contentServers``    |  hash  | A list of Traffic Servers and the delivery services associated with each.    |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``deliveryServices``  |  hash  | A list of delivery services.                                                 |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``edgeLocations``     |  hash  | A list of cache groups.                                                      |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``stats``             |  hash  | Snapshot properties.                                                         |
  +-----------------------+--------+------------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": {
			"config": { ... },
			"contentRouters": { ... },
			"contentServers": { ... },
			"deliveryServices": { ... },
			"edgeLocations": { ... },
			"stats": { ... },
		},
    }

|

**GET /api/1.2/cdns/{:cdn_name}/snapshot/new**

  Retrieves a PENDING snapshot for a CDN which represents the current state of the CDN. The contents of this snapshot are NOT currently used by Traffic Monitor and Traffic Router. Once a snapshot is performed, this snapshot will become the CURRENT snapshot and will be used by Traffic Monitor and Traffic Router.

  Authentication Required: Yes

  Role(s) Required: Admin or Operations

  **Request Route Parameters**

  +----------------+----------+---------------------------------------------+
  |   Name         | Required |                Description                  |
  +================+==========+=============================================+
  |  ``cdn_name``  |   yes    | CDN name.                                   |
  +----------------+----------+---------------------------------------------+

  **Response Properties**

  +-----------------------+--------+------------------------------------------------------------------------------+
  |    Parameter          |  Type  |                               Description                                    |
  +=======================+========+==============================================================================+
  | ``config``            |  hash  | General CDN configuration settings.                                          |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``contentRouters``    |  hash  | A list of Traffic Routers.                                                   |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``contentServers``    |  hash  | A list of Traffic Servers and the delivery services associated with each.    |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``deliveryServices``  |  hash  | A list of delivery services.                                                 |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``edgeLocations``     |  hash  | A list of cache groups.                                                      |
  +-----------------------+--------+------------------------------------------------------------------------------+
  | ``stats``             |  hash  | Snapshot properties.                                                         |
  +-----------------------+--------+------------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": {
			"config": { ... },
			"contentRouters": { ... },
			"contentServers": { ... },
			"deliveryServices": { ... },
			"edgeLocations": { ... },
			"stats": { ... },
		},
    }

|

**PUT /api/1.2/snapshot/{:cdn_name}**

  Authentication Required: Yes

  Role(s) Required: Admin or Operations

  **Request Route Parameters**

  +----------+----------+-------------------------------------------+
  | Name     | Required | Description                               |
  +==========+==========+===========================================+
  | cdn_name | yes      | The name of the cdn to snapshot configure |
  +----------+----------+-------------------------------------------+

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |response              | string |  "SUCCESS"                                     |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": "SUCCESS"
    }

|
