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


.. _to-api-v12-config-diffs:

Configuration Differences
=========================

.. _to-api-v12-config-diffs-route:

/api/1.2/servers/:host-name/config_diffs.json
+++++++++++++++++++++++++++++++++++++++++++++

**GET servers/:host-name/config_diffs.json**

  Retrieves a all of the configuration file differences between a server and Traffic Ops

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-----------------+----------+----------------------------------------------------------+
  | Name            | Required | Description                                              |
  +=================+==========+==========================================================+
  | ``host-name``   | yes      | The host name of the server to get the differences for   |
  +-----------------+----------+----------------------------------------------------------+

  **Response Properties**

  +-----------------------+----------+--------------------------------------------------------------------------+
  | Parameter             | Type     | Description                                                              |
  +=======================+==========+==========================================================================+
  | ``fileName``          | string   | The name of the configuration file                                       |
  +-----------------------+----------+--------------------------------------------------------------------------+
  | ``dbLinesMissing``    | string[] | The lines not in the Traffic Ops database, but on the server             |
  +-----------------------+----------+--------------------------------------------------------------------------+
  | ``diskLinesMissing``  | string[] | The lines in the Traffic Ops database, but not on the server             |
  +-----------------------+----------+--------------------------------------------------------------------------+
  | ``timestamp``         | string   | The last time the server updated this entry                              |
  +-----------------------+----------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": [
        {
          "fileName": "configName1",
          "dbLinesMissing": [ 
            "LocalOnlyLine One",
            "Local Only Line Two"
          ],
          "diskLinesMissing": [
            "DBOnlyLine One",
            "DB Only Line Two"
          ],
          "timestamp": "2015-02-03 17:04:20"
        },
        {
          "fileName": "otherConfigName",
          "dbLinesMissing": [ 
            "Config Line Local",
            "Another Config Line Local"
          ],
          "diskLinesMissing": [
            "DB Only Line",
            "Another Config DB Line"
          ],
          "timestamp": "2015-02-03 17:04:20"
        },
     ],
    }

|

**PUT /api/1.2/servers/:host-name/:cfg-file-name**

  Updates the configuration file differences between the server and Traffic Ops.

  Authentication Required: Yes

  Role(s) Required: None

  **Request Route Parameters**

  +-------------------+----------+---------------------------------------------------------------+
  | Name              | Required | Description                                                   |
  +===================+==========+===============================================================+
  | ``host-name``     | yes      | The host name of the server to set the differences for        |
  +-------------------+----------+---------------------------------------------------------------+
  | ``cfg-file-name`` | yes      | The name of the configuration file to update differences for  |
  +-------------------+----------+---------------------------------------------------------------+

  **Request Properties**

  +-----------------------+----------+--------------------------------------------------------------------------+
  | Parameter             | Type     | Description                                                              |
  +=======================+==========+==========================================================================+
  | ``fileName``          | string   | The name of the configuration file (optional)                            |
  +-----------------------+----------+--------------------------------------------------------------------------+
  | ``dbLinesMissing``    | string[] | The lines not in the Traffic Ops database, but on the server             |
  +-----------------------+----------+--------------------------------------------------------------------------+
  | ``diskLinesMissing``  | string[] | The lines in the Traffic Ops database, but not on the server             |
  +-----------------------+----------+--------------------------------------------------------------------------+
  | ``timestamp``         | string   | The last time the server updated this entry                              |
  +-----------------------+----------+--------------------------------------------------------------------------+

  **Request Example** ::

    {
      "dbLinesMissing": [ 
        "LocalOnlyLine One",
        "Local Only Line Two"
      ],
      "diskLinesMissing": [
        "DBOnlyLine One",
        "DB Only Line Two"
      ],
      "timestamp": "2015-02-03 17:04:20"
    }

  **Response Properties**
  
  No Response body
