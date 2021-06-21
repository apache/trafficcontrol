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

.. _to-api-v3-cdns-dnsseckeys-generate:

****************************
``cdns/dnsseckeys/generate``
****************************

``POST``
========
Generates :abbr:`ZSK (Zone-Signing Key)` and :abbr:`KSK (Key-Signing Key)` keypairs for a CDN and all associated :term:`Delivery Services`.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
:effectiveDate:         An optional string containing the date and time at which the newly-generated :abbr:`ZSK (Zone-Signing Key)` and :abbr:`KSK (Key-Signing Key)` become effective, in :RFC:`3339` format. Defaults to the current time if not specified.
:key:                   Name of the CDN
:kskExpirationDays:     Expiration (in days) for the :abbr:`KSKs (Key-Signing Keys)`
:ttl:                   Time, in seconds, for which the keypairs shall remain valid
:zskExpirationDays:     Expiration (in days) for the :abbr:`ZSKs (Zone-Signing Keys)`

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/cdns/dnsseckeys/generate HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 130

	{
		"key": "CDN-in-a-Box",
		"kskExpirationDays": 1095,
		"ttl": 3600,
		"zskExpirationDays": 1095
	}


Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 19:42:15 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: O9SPWzeMNFgg6I/PPeXittBIhdh3/zUKK1NwNlYIM9SszSrk0h/Dfz7tnwgnA7h/s6M4eYBJxykDpCfVC7xpeg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 18:42:15 GMT
	Content-Length: 89

	{
		"response": "Successfully created dnssec keys for CDN-in-a-Box"
	}
