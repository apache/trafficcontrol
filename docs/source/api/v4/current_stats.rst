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


.. _to-api-v4-current-stats:

*****************
``current_stats``
*****************
An API endpoint that returns current statistics for each CDN and an aggregate across them all.

``GET``
=======
Retrieves current stats for each CDN. Also includes aggregate stats across them.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: CDN:READ
:Response Type:  Array

Request Structure
-----------------
No parameters available.

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/current_stats HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cdn:         The name of the CDN
:connections: Current number of TCP connections maintained
:capacity:    85 percent capacity of the CDN in Gps
:bandwidth:   The total amount of bandwidth in Gbs

.. note:: If ``cdn`` name is total and capacity is omitted it represents the aggregate stats across CDNs

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=; Path=/; HttpOnly
	Whole-Content-Sha512: Rs3wgd7v5dP0bOQs4I3J1q6mnWIMSM2AKSAWirK1kymvDYOoFISArF7Kkypgy10I34yn7FtFdMh6U7ABaS1Tjw==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 14 Nov 2019 15:35:31 GMT
	Content-Length: 138

	{"response": {
		"currentStats": [
			{
				"bandwidth": null,
				"capacity": null,
				"cdn": "ALL",
				"connections": null
			},
			{
				"bandwidth": 0.000104,
				"capacity": 17,
				"cdn": "CDN-in-a-Box",
				"connections": 4
			},
			{
				"bandwidth": 0.000104,
				"cdn": "total",
				"connections": 4
			}
		]
	}}
