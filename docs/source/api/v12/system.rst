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

.. _to-api-v12-sys:

System
======

.. _to-api-v12-sys-route:

/api/1.1/system
+++++++++++++++

**GET /api/1.2/system/info.json**

  Authentication Required: Yes

  Role(s) Required: None

  **Response Properties**

  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  |            Key             |  Type  |                                                             Description                                                              |
  +============================+========+======================================================================================================================================+
  | ``parameters``             | hash   | This is a hash with the parameter names that describe the Traffic Ops installation as keys.                                          |
  |                            |        | These are all the parameters in the ``GLOBAL`` profile.                                                                              |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>tm.toolname``           | string | The name of the Traffic Ops tool. Usually "Traffic Ops". Used in the About screen and in the comments headers of the files generated |
  |                            |        | (``# DO NOT EDIT - Generated for atsec-lax-04 by Traffic Ops (https://traffops.kabletown.net/) on Fri Mar  6 05:15:15 UTC 2015``).   |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>tm.instance_name``      | string | The name of the Traffic Ops instance. Can be used when multiple instances are active. Visible in the About page.                     |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>traffic_rtr_fwd_proxy`` | string | When collecting stats from Traffic Router, Traffic Ops uses this forward proxy to pull the stats through.                            |
  |                            |        | This can be any of the MID tier caches, or a forward cache specifically deployed for this purpose. Setting                           |
  |                            |        | this variable can significantly lighten the load on the Traffic Router stats system and it is recommended to                         |
  |                            |        | set this parameter on a production system.                                                                                           |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>tm.url``                | string | The URL for this Traffic Ops instance. Used in the About screen and in the comments headers of the files generated                   |
  |                            |        | (``# DO NOT EDIT - Generated for atsec-lax-04 by Traffic Ops (https://traffops.kabletown.net/) on Fri Mar  6 05:15:15 UTC 2015``).   |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>traffic_mon_fwd_proxy`` | string | When collecting stats from Traffic Monitor, Traffic Ops uses this forward proxy to pull the stats through.                           |
  |                            |        | This can be any of the MID tier caches, or a forward cache specifically deployed for this purpose. Setting                           |
  |                            |        | this variable can significantly lighten the load on the Traffic Monitor system and it is recommended to                              |
  |                            |        | set this parameter on a production system.                                                                                           |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>tm.logourl``            | string | This is the URL of the logo for Traffic Ops and can be relative if the logo is under traffic_ops/app/public.                         |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
  | ``>tm.infourl``            | string | This is the "for more information go here" URL, which is visible in the About page.                                                  |
  +----------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+

  **Response Example** ::

    {
      "response": {
        "parameters": {
          "tm.toolname": "Traffic Ops",
          "tm.infourl": "http:\/\/staging-03.cdnlab.kabletown.net\/tm\/info",
          "traffic_mon_fwd_proxy": "http:\/\/proxy.kabletown.net:81",
          "traffic_rtr_fwd_proxy": "http:\/\/proxy.kabletown.net:81",
          "tm.logourl": "\/images\/tc_logo.png",
          "tm.url": "https:\/\/tm.kabletown.net\/",
          "tm.instance_name": "Kabletown CDN"
        }
      },
    }

|

