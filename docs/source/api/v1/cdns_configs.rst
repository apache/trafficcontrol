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
.. _to-api-v1-cdns-config:

****************
``cdns/configs``
****************
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-v1-cdns` instead.
.. danger:: This endpoint does not appear to work, and thus its use is strongly discouraged!

``GET``
=======
Retrieves CDN configuration information.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Properties
-------------------
:id:          The integral, unique identifier for this CDN
:name:        The CDN's name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: z9P1NkxGebPncUhaChDHtYKYI+XVZfhE6Y84TuwoASZFIMfISELwADLpvpPTN+wwnzBfREksLYn+0313QoBWhA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 20:46:57 GMT
	Content-Length: 237

	{ "response": [
		{
			"id": 1,
			"name": "ALL"
		},
		{
			"id": 2,
			"name": "CDN-in-a-Box"
		}
	],
	"alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /cdns instead",
			"level": "warning"
		}
	]}
