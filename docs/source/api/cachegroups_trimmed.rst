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

.. _to-api-cachegroups-trimmed:

***********************
``cachegroups/trimmed``
***********************
Extract just the names of all :term:`Cache Groups`.

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available

Response Structure
------------------
:name: The name of the :term:`Cache Group`                    |

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: OyOKqpB24AMlrENIEoA4la/3rclnuKMayvzskmPNPXrDMQksGt0UjVwORYmMdmIS5dQHuIlglBlksvLtqjziHQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 20:23:23 GMT
	Content-Length: 216

	{ "response": [
		{
			"name": "TRAFFIC_ANALYTICS"
		},
		{
			"name": "TRAFFIC_OPS"
		},
		{
			"name": "TRAFFIC_OPS_DB"
		},
		{
			"name": "TRAFFIC_PORTAL"
		},
		{
			"name": "TRAFFIC_STATS"
		},
		{
			"name": "CDN_in_a_Box_Mid"
		},
		{
			"name": "CDN_in_a_Box_Edge"
		},
		{
			"name": "test"
		}
	]}

