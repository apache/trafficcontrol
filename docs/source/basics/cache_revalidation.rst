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

.. _cache-revalidation:

**************************************
Cache Control Headers and Revalidation
**************************************
The HTTP/1.1 specification in :rfc:`2616#section-14.9` allows for origin servers and clients to influence how caches treat their requests and responses. By default, the Traffic Control CDN will honor cache control headers. Most commonly, origin servers will tell the downstream caches how long a response can be cached.

.. note:: The terms "content revalidation" and "content invalidation" are often convoluted when referring to the same behavior. Within the context of Traffic Control, the two should be considered synonymous.

.. code-block:: http
	:caption: This Response may Only be Cached for 86400 Seconds

	HTTP/1.1 200 OK
	Date: Sun, 14 Dec 2014 23:22:44 GMT
	Server: Apache/2.2.15 (Red Hat)
	Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
	ETag: "1aa008f-2d-50a3559482cc0"
	Cache-Control: max-age=86400
	Content-Length: 45
	Connection: close
	Content-Type: text/html; charset=UTF-8

	<!DOCTYPE html><html><body>This is a fun file</body></html>

The ``max-age`` directive in the ``Cache-Control`` header tells downstream caching systems that the maximum time for which they are allowed to cache this response is the specified number of seconds. The origin can also add an ``Expires:`` header, explicitly telling the cache the time this response is to be expired. When a response is expired it usually doesn't get deleted from the cache, but, when a request comes in that would have hit on this response if it was not expired, the cache *revalidates* the response. In particular, this is the way :abbr:`ATS (Apache Traffic Server)` handles content revalidation. Instead of requesting the object again from the origin server, the cache will send a request to the origin indicating what version of the response it has, and asking if it has changed. If it changed, the server will send a ``200 OK`` response, with the new data. If it has not changed, the origin server will send back a ``304 Not Modified`` response indicating the response is still valid, and that the cache can reset the timer on the response expiration. To indicate what version the client (cache) has it will add an ``If-Not-Modified-Since:`` header, or an ``If-None-Match:`` header. For example, in the ``If-None-Match:`` case, the origin will have sent an ``ETag`` header that uniquely identifies the response. The client can then later use that in a revalidation request to check if the ``ETag`` of the requested content has changed.

.. code-block:: http
	:caption: The Cache Server Sends a Request with the Old ``ETag`` Value in the ``If-None-Match`` Header

	GET /foo/bar/fun.html HTTP/1.1
	If-None-Match: "1aa008f-2d-50a3559482cc0"
	Host: www.origin.com

If the content has changed (meaning, the new response would not have had the same ``ETag``) the server MUST respond with the up-to-date content, usually in the body of a ``200 OK`` response.

.. code-block:: http
	:caption: The Origin Responds with the Modified Content and a New ``ETag``

	HTTP/1.1 200 OK
	Date: Sun, 18 Dec 2014 3:22:44 GMT
	Server: Apache/2.2.15 (Red Hat)
	Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
	ETag: "1aa008f-2d-50aa00feadd"
	Cache-Control: max-age=604800
	Content-Length: 49
	Connection: close
	Content-Type: text/html; charset=UTF-8

	<!DOCTYPE html><html><body>This is NOT a fun file</body></html>


If the content did not change (meaning, the response would have had the same ``ETag``) the server SHOULD respond with a ``304 Not Modified``. In most cases, the server will also send back an ``ETag`` header, since the client is allowed to send multiple ``ETag`` values in its ``If-None-Match`` header to check against multiple cached versions of the content, and the ``ETag`` will tell it which specifically is the current version. This is a very rare use case, and :abbr:`ATS (Apache Traffic Server)` will not make use of this feature without a plugin to modify its behavior.

.. code-block:: http
	:caption: The Content has not been Modified so the Server Indicates the Cached Version is Up-To-Date

	HTTP/1.1 304 Not Modified
	Date: Sun, 18 Dec 2014 3:22:44 GMT
	Server: Apache/2.2.15 (Red Hat)
	Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
	ETag: "1aa008f-2d-50a3559482cc0"
	Cache-Control: max-age=604800
	Connection: close
