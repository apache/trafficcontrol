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


.. _to-api-caches-stats:

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
:cachegroup:  The name of the :term:`Cache Group` to which this cache belongs
:connections: Current number of TCP connections maintained by the cache
:healthy:     ``true`` if Traffic Monitor has marked the cache as "healthy", ``false`` otherwise

	.. seealso:: :ref:`health-proto`

:hostname:    The (short) hostname of the cache
:ip:          The IP address of the cache
:kbps:        Cache upload speed (to clients) in Kilobits per second
:profile:     The :ref:`profile-name` of the :term:`Profile` in use by this :term:`cache server`
:status:      The status of the cache

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 20:25:01 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 15 Nov 2018 00:25:01 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: DqbLgitanS8q81/qKC1i+ImMiEMF+SW4G9rb79FWdeWcgwFjL810tlTRp1nNNfHV+tajgjyK+wMHobqVyaNEfA==
	Content-Length: 133

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
