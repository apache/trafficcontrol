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

.. _to-api-v14-cdns:

CDN
==========

.. _to-api-v14-cdns-route:

/api/1.4/cdns
++++++++++++++++++++

**GET api/1.4/cdns/:name/dnsseckeys/ksk/generate**

  Authentication Required: Yes

  Role(s) Required: None

  **Request Query Parameters**

  +--------------------+----------+----------------------------------------------------------------------------------------------------+
  | Name               | Required | Description                                                                                        |
  +====================+==========+====================================================================================================+
  | ``expirationDays`` | yes      | The number of days until the new generated ksk expires.                                            |
  +--------------------+----------+----------------------------------------------------------------------------------------------------+
  | ``effectiveDate``  | no       | The time the new generated ksk becomes effective, in RFC3339 format. Defaults to the current time. |
  +--------------------+----------+----------------------------------------------------------------------------------------------------+

  **Request Example** ::

    {
    	"expirationDays": 100,
    	"effectiveDate": "2021-01-01T00:00:00.0000000-04:00"
    }

|

  **Response Properties**

  +-----------------------------------+--------+--------------------------------------------------------------------------+
  | Parameter                         | Type   | Description                                                              |
  +===================================+========+==========================================================================+
  | ``response``                      | string | response string                                                          |
  +-----------------------------------+--------+--------------------------------------------------------------------------+

  **Response Example** ::

    {
     "response": Successfully generated ksk dnssec keys for my-cdn-name"
    }

|

