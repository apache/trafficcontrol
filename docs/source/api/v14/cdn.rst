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

**POST api/1.4/cdns/:name/dnsseckeys/ksk/generate**

  Authentication Required: Yes

  Role(s) Required: admin

  **Request Properties**

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
     "response": "Successfully generated ksk dnssec keys for my-cdn-name"
    }

|

**GET /api/1.4/cdns/dnsseckeys/refresh**

  Refresh the DNSSEC keys for all CDNs. This call initiates a background process to refresh outdated keys, and immediately returns a response that the process has started.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Route Parameters**

  None.

  **Response Properties**

  +--------------+--------+------------------+
  |  Parameter   |  Type  |   Description    |
  +==============+========+==================+
  | ``response`` | string | success response |
  +--------------+--------+------------------+

  **Response Example**
  ::

    {
      "response": "Checking DNSSEC keys for refresh in the background"
    }

|

.. _to-api-v12-cdn-dnsseckeys:

DNSSEC Keys
+++++++++++

**GET /api/1.2/cdns/name/:name/dnsseckeys**

  Gets a list of dnsseckeys for a CDN and all associated Delivery Services.

  Authentication Required: Yes

  Role(s) Required: Admin

  **Request Route Parameters**

	  +----------+----------+-------------+
  |   Name   | Required | Description |
  +==========+==========+=============+
  | ``name`` | yes      |             |
  +----------+----------+-------------+

  **Response Properties**

  +-------------------------------+--------+---------------------------------------------------------------+
  |           Parameter           |  Type  |                          Description                          |
  +===============================+========+===============================================================+
  | ``cdn name/ds xml_id``        | string | identifier for ds or cdn                                      |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>zsk/ksk``                  | array  | collection of zsk/ksk data                                    |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>>ttl``                     | string | time-to-live for dnssec requests                              |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>>inceptionDate``           | string | epoch timestamp for when the keys were created                |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>>expirationDate``          | string | epoch timestamp representing the expiration of the keys       |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>>private``                 | string | encoded private key                                           |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>>public``                  | string | encoded public key                                            |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``>>name``                    | string | domain name                                                   |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``version``                   | string | API version                                                   |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``ksk>>dsRecord>>algorithm``  | string | The algorithm of the referenced DNSKEY-recor.                 |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``ksk>>dsRecord>>digestType`` | string | Cryptographic hash algorithm used to create the Digest value. |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``ksk>>dsRecord>>digest``     | string | A cryptographic hash value of the referenced DNSKEY-record.   |
  +-------------------------------+--------+---------------------------------------------------------------+
  | ``ksk>>dsRecord>>text``       | string | The DS Record text, to be inserted in the parent resolver.    |
  +-------------------------------+--------+---------------------------------------------------------------+

  **Response Example** ::

    {
      "response": {
        "cdn1": {
          "zsk": {
            "ttl": "60",
            "inceptionDate": "1426196750",
            "private": "zsk private key",
            "public": "zsk public key",
            "expirationDate": "1428788750",
            "name": "foo.kabletown.com."
          },
          "ksk": {
            "name": "foo.kabletown.com.",
            "expirationDate": "1457732750",
            "public": "ksk public key",
            "private": "ksk private key",
            "inceptionDate": "1426196750",
            "ttl": "60",
            "dsRecord": {
              "algorithm": "5",
              "digestType": "2",
              "digest": "abc123def456",
              "text": "foo.kabletown.com.\t30\tIN\tDS\t12345 8 2 DEADBEEF123456789"
            }
          }
        },
        "ds-01": {
          "zsk": {
            "ttl": "60",
            "inceptionDate": "1426196750",
            "private": "zsk private key",
            "public": "zsk public key",
            "expirationDate": "1428788750",
            "name": "ds-01.foo.kabletown.com."
          },
          "ksk": {
            "name": "ds-01.foo.kabletown.com.",
            "expirationDate": "1457732750",
            "public": "ksk public key",
            "private": "ksk private key",
            "inceptionDate": "1426196750"
          }
        },
        ... repeated for each ds in the cdn
      }
    }

|
