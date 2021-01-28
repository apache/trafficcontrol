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

.. _to-api-v3-cdns-dnsseckeys-refresh:

***************************
``cdns/dnsseckeys/refresh``
***************************

``GET``
=======
Refresh the DNSSEC keys for all CDNs. This call initiates a background process to refresh outdated keys, and immediately returns a response that the process has started.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object (string)

Request Structure
-----------------
No parameters available

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 21:37:30 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Uwl+924m6Ye3NraFP+RBpldkhcNTTDyXHZbzRaYV95p9tP56Z61gckeKSr1oQIkNXjXcCsDN5Dmum7Zk1AR6Hw==
	Content-Length: 69

	{
		"response": "Checking DNSSEC keys for refresh in the background"
	}
