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

.. _to-api-parameter:

Parameter
=========
**GET /api/1.1/parameters.json**


.. Description.

Authentication Required: Yes

Response Content Type: application/json

**Response Messages**

::


  HTTP Status Code: 200
  Reason: Success


**Return Values**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
|``last_updated``      | string |                                                |
+----------------------+--------+------------------------------------------------+
|``value``             | string |                                                |
+----------------------+--------+------------------------------------------------+
|``name``              | string |                                                |
+----------------------+--------+------------------------------------------------+
|``config_file``       | string |                                                |
+----------------------+--------+------------------------------------------------+

**Response Example**


::


  {
   "response": [
      {
         "last_updated": "2012-09-17 21:41:22",
         "value": "foo.bar.net",
         "name": "domain_name",
         "config_file": "FooConfig.xml"
      },
      {
         "last_updated": "2012-09-17 21:41:22",
         "value": "0,1,2,3,4,5,6",
         "name": "Drive_Letters",
         "config_file": "storage.config"
      },
      {
         "last_updated": "2012-09-17 21:41:22",
         "value": "STRING __HOSTNAME__",
         "name": "CONFIG proxy.config.proxy_name",
         "config_file": "records.config"
      }
   ],
   "version": "1.1"
  }

For error messages, see :ref:`reference-label-401`.

|

**GET /api/1.1/parameters/profile/:profile_name.json**


.. Description.


Authentication Required: Yes

**Request Route Parameters**

+-----------------+----------+---------------------------------------------------+
| Name            | Required | Description                                       |
+=================+==========+===================================================+
|profile_name     | yes      |                                                   |
+-----------------+----------+---------------------------------------------------+

Response Content Type: application/json

**Response Messages**

::


  HTTP Status Code: 200
  Reason: Success


**Return Values**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
|``last_updated``      | string |                                                |
+----------------------+--------+------------------------------------------------+
|``value``             | string |                                                |
+----------------------+--------+------------------------------------------------+
|``name``              | string |                                                |
+----------------------+--------+------------------------------------------------+
|``config_file``       | string |                                                |
+----------------------+--------+------------------------------------------------+
|``version``           | string |                                                |
+----------------------+--------+------------------------------------------------+

**Response Example**


::


  {
   "response": [
      {
         "last_updated": "2012-09-17 21:41:22",
         "value": "foo.bar.net",
         "name": "domain_name",
         "config_file": "FooConfig.xml"
      },
      {
         "last_updated": "2012-09-17 21:41:22",
         "value": "0,1,2,3,4,5,6",
         "name": "Drive_Letters",
         "config_file": "storage.config"
      },
      {
         "last_updated": "2012-09-17 21:41:22",
         "value": "STRING __HOSTNAME__",
         "name": "CONFIG proxy.config.proxy_name",
         "config_file": "records.config"
      }
   ],
   "version": "1.1"
  }

For error messages, see :ref:`reference-label-401`.

