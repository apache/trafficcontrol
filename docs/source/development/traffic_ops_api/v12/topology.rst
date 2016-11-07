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

**PUT /api/1.2/snapshot/{:cdn_name}**

  Authentication Required: Yes

  Role(s) Required: admin or oper

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
