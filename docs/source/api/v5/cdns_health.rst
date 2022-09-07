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

.. _to-api-cdns-health:

***************
``cdns/health``
***************
Extract health information from all :term:`Cache Groups` across all CDNs

.. seealso:: :ref:`health-proto`

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Permissions Required: CACHE-GROUP:READ
:Response Type:  Object

Request Structure
-----------------
No parameters available

Response Structure
------------------
:cachegroups:  An array of objects describing the health of each Cache Group

	:name:    The name of the Cache Group
	:offline: The number of OFFLINE caches in the Cache Group
	:online:  The number of ONLINE caches in the Cache Group

:totalOffline: Total number of OFFLINE caches across all Cache Groups which are assigned to any CDN
:totalOnline:  Total number of ONLINE caches across all Cache Groups which are assigned to any CDN

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"totalOffline": 0,
		"totalOnline": 1,
		"cachegroups": [
			{
					"offline": 0,
					"name": "CDN_in_a_Box_Edge",
					"online": 1
				}
		]
	}}
