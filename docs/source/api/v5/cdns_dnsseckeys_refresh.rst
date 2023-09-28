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

.. _to-api-cdns-dnsseckeys-refresh:

***************************
``cdns/dnsseckeys/refresh``
***************************

``PUT``
=======
Refresh the DNSSEC keys for all CDNs. This call initiates a background process to refresh outdated keys, and immediately returns a response that the process has started.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: DNS-SEC:UPDATE, CDN:UPDATE, CDN:READ
:Response Type: ``undefined``

Request Structure
-----------------
No parameters available

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 202 Accepted
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Location: /api/5.0/async_status/3
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 20 Jul 2021 23:55:11 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: yJUGNCYygBYvHft4z0nxJ0/p230s3PdPT5Tld+8hIWfxmpmKDciY4D7+1Bf8S69ckmZR/yxY95kIZEbg9/jFgw==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 20 Jul 2021 22:55:11 GMT
	Content-Length: 176

	{
		"alerts": [
			{
				"text": "Starting DNSSEC key refresh in the background. This may take a few minutes. Status updates can be found here: /api/5.0/async_status/3",
				"level": "success"
			}
		]
	}
