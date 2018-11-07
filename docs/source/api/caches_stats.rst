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


.. _to-api-caches_stats:

****************
``caches/stats``
****************
An API endpoint that returns cache statistics using the :ref:`tm-api`.

.. seealso:: This gives a set of basic statistics for *all caches* at the current time. For statistics from time ranges and/or aggregated over a specific CDN, use :ref:`to-api-cache_stats`.

``GET``
=======
Retrieves cache stats from Traffic Monitor. Also includes rows for aggregates.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:cachegroup:  The name of the Cache Group to which this cache belongs
:connections: Current number of TCP connections maintained by the cache
:healthy:     ``true`` if Traffic Monitor has marked the cache as "healthy", ``false`` otherwise

	.. seealso:: :ref:`health-proto`

:hostname:    The (short) hostname of the cache
:ip:          The IP address of the cache
:kbps:        Cache upload speed (to clients) in Kilobits per second
:profile:     The name of the profile in use by this cache
:status:      The status of the cache

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"profile": "ALL",
			"connections": 0,
			"ip": null,
			"status": "ALL",
			"healthy": true,
			"kbps": 0,
			"hostname": "ALL",
			"cachegroup": "ALL"
		}
	]}
