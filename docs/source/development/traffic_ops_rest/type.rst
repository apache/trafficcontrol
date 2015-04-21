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

.. _to-api-type:

Types
=====

**GET /api/1.1/types.json**

  Authentication Required: Yes

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``lastUpdated``       | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``useInTable``        | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``description``       | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "lastUpdated": "2013-10-23 15:28:31",
           "useInTable": "staticdnsentry",
           "name": "AAAA_RECORD",
           "id": "22",
           "description": "Static DNS AAAA entry"
        }
     ],
     "version": "1.1"
    }


|


**GET /api/1.1/types/trimmed.json**

  Authentication Required: Yes

  Response Content Type: application/json

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "name": "AAAA_RECORD"
        },
        {
           "name": "ACTIVE_DIRECTORY"
        },
        {
           "name": "A_RECORD"
        },
        {
           "name": "CCR"
        }
     ],
     "version": "1.1"
    }

