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

.. _to-api-v1-servers-checks:

******************
``servers/checks``
******************
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-servercheck` instead.

Deals with the various checks associated with certain types of servers.

.. seealso:: :ref:`to-check-ext`

``GET``
=======
Fetches identifying and meta information as well as "check" values regarding all servers that have a :term:`Type` with a name beginning with "EDGE" or "MID" (ostensibly this is equivalent to all :term:`cache servers`).

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:adminState:   The name of the server's :term:`Status` - called "adminState" for legacy reasons
:cacheGroup:   The name of the :term:`Cache Group` to which the server belongs
:checks:       An optionally present map of the names of "checks" to their values. Only numeric and boolean checks are represented, and boolean checks are represented as integers with ``0`` meaning "false" and ``1`` meaning "true". Will not appear if the server in question has no valued "checks".
:hostName:     The (short) hostname of the server
:id:           The server's integral, unique identifier
:profile:      The name of the :term:`Profile` used by the server
:revalPending: A boolean that indicates whether or not the server has pending revalidations
:type:         The name of the server's :term:`Type`
:updPending:   A boolean that indicates whether or not the server has pending updates

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Thu, 23 Jan 2020 20:00:19 GMT; Max-Age=3600; HttpOnly
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 23 Jan 2020 19:00:19 GMT
	Content-Length: 449

	{ "alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /servercheck instead",
			"level": "warning"
		}
	],
	"response": [
		{
			"adminState": "REPORTED",
			"cacheGroup": "CDN_in_a_Box_Edge",
			"id": 12,
			"hostName": "edge",
			"revalPending": false,
			"profile": "ATS_EDGE_TIER_CACHE",
			"type": "EDGE",
			"updPending": false
		},
		{
			"adminState": "REPORTED",
			"cacheGroup": "CDN_in_a_Box_Mid",
			"id": 11,
			"hostName": "mid",
			"revalPending": false,
			"profile": "ATS_MID_TIER_CACHE",
			"type": "MID",
			"updPending": false
		}
	]}
