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

.. _to-api-v3-consistenthash:

******************
``consistenthash``
******************
Test Pattern-Based Consistent Hashing for a :term:`Delivery Service` using a regular expression and a request path

``POST``
========
Queries database for an active Traffic Router on a given CDN and sends GET request to get the resulting path to consistent hash with a given regex and request path.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
:regex:       The regular expression to apply to the request path to get a resulting path that will be used for consistent hashing
:requestPath: The request path to use to test the regular expression against
:cdnId:       The unique identifier of a CDN that will be used to query for an active Traffic Router

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/consistenthash HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.54.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 80
	Content-Type: application/json

	{"regex":"/.*?(/.*?/).*?(m3u8)","requestPath":"/test/path/asset.m3u8","cdnId":2}

Response Structure
------------------
:resultingPathToConsistentHash: The resulting path that Traffic Router will use for consistent hashing
:consistentHashRegex:           The regex used by Traffic Router derived from POST 'regex' parameter
:requestPath:                   The request path used by Traffic Router to test regex against

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: QMDFOnUfqH4TcZ4YnUQyqnXDier0YiUMIfwBGDcT7ySjw9uASBGsLQW35lpnKFi4as0vYlHuSSGpe4hHGsladQ==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 12 Feb 2019 21:32:05 GMT
	Content-Length: 142

	{ "response": {
		"resultingPathToConsistentHash":"/path/m3u8",
		"consistentHashRegex":"/.*?(/.*?/).*?(m3u8)",
		"requestPath":"/test/path/asset.m3u8"
	}}
