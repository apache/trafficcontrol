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
.. _to-api-v4-vault-ping:

**************
``vault/ping``
**************

``GET``
=======
Pings Traffic Vault to retrieve status.

:Auth. Required: Yes
:Roles Required: "read-only"
:Permissions Required: TRAFFIC-VAULT:READ
:Response Type:  Object

Request Structure
-----------------
No parameters available.

Response Properties
-------------------
:status:        The status returned from the ping request to the Traffic Vault server
:server:        The Traffic Vault server that was pinged

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 25 Feb 2020 15:37:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: z9P1NkxGebPncUhaChDHtYKYI+XVZfhE6Y84TuwoASZFIMfISELwADLpvpPTN+wwnzBfREksLYn+0313QoBWhA==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 25 Feb 2020 14:37:55 GMT
	Content-Length: 90

	{ "response":
		{
			"status": "OK",
			"server": "trafficvault.infra.ciab.test:8087"
		}
	}
