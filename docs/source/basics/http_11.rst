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

********
HTTP/1.1
********
For a comprehensive look at Traffic Control, it is important to understand basic HTTP/1.1 protocol operations and how :term:`cache servers` function.

.. seealso:: For complete details on HTTP/1.1 see :rfc:`2616`.

What follows is a sequence of events that take place when a client requests content from an HTTP/1.1-compliant server.

#. The client sends a request to the :abbr:`LDNS (Local DNS)` server to resolve the name ``www.origin.com`` to an IP address, then sends an HTTP request to that IP.

	.. Note:: A DNS response is accompanied by a :abbr:`TTL (Time To Live)` which indicates for how long a name resolution should be considered valid. While longer DNS :abbr:`TTL (Time To Live)`\ s of a day (86400 seconds) or more are quite common in other use cases, in CDN use-cases DNS :abbr:`TTL (Time To Live)`\ s are often below a minute.

	.. code-block:: http
		:caption: A Client Request for ``/foo/bar/fun.html`` from ``www.origin.com``

		GET /foo/bar/fun.html HTTP/1.1
		Host: www.origin.com

#. The server at ``www.origin.com`` looks up the content of the path ``/foo/bar/fun.html`` and sends it in a response to the client.

	.. code-block:: http
		:caption: Server Response

		HTTP/1.1 200 OK
		Content-Type: text/html; charset=UTF-8
		Content-Length: 45

		<!DOCTYPE html><html><body>This is a fun file</body></html>
