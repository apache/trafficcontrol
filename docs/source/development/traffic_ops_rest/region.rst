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

.. _to-api-region:

Regions
=======

**GET /api/1.1/regions.json**

  Authentication Required: 

  Response Content Type: application/json

  **Response Properties**

  +----------------------+--------+------------------------------------------------+
  | Parameter            | Type   | Description                                    |
  +======================+========+================================================+
  |``name``              | string |                                                |
  +----------------------+--------+------------------------------------------------+
  |``id``                | string |                                                |
  +----------------------+--------+------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
           "name": "Atlanta",
           "id": "6"
        },
        {
           "name": "Beltway",
           "id": "1"
        }
     ],
     "version": "1.1"
    }

