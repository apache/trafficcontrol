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

.. index::
	http/1.1
	HTTP

HTTP 1.1
========
For a comprehensive look at Traffic Control, it is important to understand basic HTTP 1.1 protocol operations and how caches function. The example below illustrates the fulfillment of an HTTP 1.1 request in a situation without CDN or proxy, followed by viewing the changes after inserting different types of (caching) proxies. Several of the examples below are simplified for clarification of the essentials.

For complete details on HTTP 1.1 see `RFC 2616 - Hypertext Transfer Protocol -- HTTP/1.1 <https://www.ietf.org/rfc/rfc2616.txt>`_.

Below are the steps of a client retrieving the URL ``http://www.origin.com/foo/bar/fun.html`` using HTTP/1.1 without proxies:

1. The client sends a request to the Local DNS (LDNS) server to resolve the name ``www.origin.com`` to an IPv4 address.

2. If the LDNS does not have this name (IPv4 mapping cached), it sends DNS requests to the ., .com, and .origin.com authoritative servers until it receives a response with the address for ``www.origin.com``. Per the DNS SPEC, this response has a Time To Live (TTL), which indicates how long this mapping can be cached at the LDNS server. In the example, the IP address found by the LDNS server for www.origin.com is 44.33.22.11.

  .. Note:: While longer DNS TTLs of a day (86400 seconds) or more are quite common in other use cases, in CDN use cases DNS TTLs are often below a minute.

3. The client opens a TCP connection from a random port locally to port 80 (the HTTP default) on 44.33.22.11, and sends this (showing the minimum HTTP 1.1 request, typically there are additional headers): ::

    GET /foo/bar/fun.html HTTP/1.1
    Host: www.origin.com

4. The server at ``www.origin.com`` looks up the Host: header to match that to a configuration section, usually referred to as a virtual host section. If the Host: header and configuration section match, the search continues for the content of the path ``/foo/bar/fun.html``, in the example, this is a file that contains ``<html><body>This is a fun file</body></html>``, so the server responds with the following: ::


      HTTP/1.1 200 OK
      Content-Type: text/html; charset=UTF-8
      Content-Length: 45

      <html><body>This is a fun file</body></html>

 At this point, HTTP transaction is complete.
