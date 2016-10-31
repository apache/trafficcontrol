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


.. |arrow| image:: fwda.png

.. _rl-caching_proxy:

Caching Proxies
===============
The main function of a CDN is to proxy requests from clients to origin servers 
and cache the results. 
To proxy, in the CDN context, is to obtain content using HTTP from an origin 
server on behalf of a client. To cache is to store the results so they can be 
reused when other clients are requesting the same content. There are three 
types of proxies in use on the Internet today which are described below.

.. index::
  Reverse Proxy

.. _rl-rev-proxy:

|arrow| Reverse Proxy
---------------------
  A reverse proxy acts on behalf of the origin server. The client is mostly unaware it is communicating with a proxy and not the actual origin. 
  All EDGE caches in a Traffic Control CDN are reverse proxies. 
  To the end user a Traffic Control based CDN appears as a reverse proxy since 
  it retrieves content from the origin server, acting on behalf of that origin server. The client requests a URL that has 
  a hostname which resolves to the reverse proxy's IP address and, in compliance 
  with the HTTP 1.1 specification, the client sends a ``Host:`` header to the reverse 
  proxy that matches the hostname in the URL. 
  The proxy looks up this hostname in a 
  list of mappings to find the origin hostname; if the hostname of the Host header is not found in the list, 
  the proxy will send an error (``404 Not Found``) to the client. 
  If the supplied hostname is found in this list of mappings, the proxy checks the cache, and when the content is not already present, connects to the 
  origin the requested ``Host:`` maps to and requests the path of the original URL, providing the origin hostname in the ``Host`` header.  The proxy then stores the URL in cache and serves the contents to the client. When there are subsequent requests for 
  the same URL, a caching proxy serves the content out of cache thereby reducing 
  latency and network traffic.

.. seealso:: `ATS documentation on reverse proxy <https://docs.trafficserver.apache.org/en/latest/admin/reverse-proxy-http-redirects.en.html#http-reverse-proxy>`_.

To insert a reverse proxy into the previous HTTP 1.1 example, the reverse proxy requires provisioning 
for ``www.origin.com``. By adding a remap rule to the cache, the reverse proxy then maps requests to 
this origin. The content owner must inform the clients, by updating the URL, to receive the content 
from the cache and not from the origin server directly. For this example, the remap rule on the 
cache is: ``http://www-origin-cache.cdn.com http://www.origin.com``.

..  Note:: In the previous example minimal headers were shown on both the request and response. In the examples that follow, the origin server response is more realistic. 

::

  HTTP/1.1 200 OK
  Date: Sun, 14 Dec 2014 23:22:44 GMT
  Server: Apache/2.2.15 (Red Hat)
  Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
  ETag: "1aa008f-2d-50a3559482cc0"
  Content-Length: 45
  Connection: close
  Content-Type: text/html; charset=UTF-8

  <html><body>This is a fun file</body></html>

The client is given the URL ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` (note the different hostname) and when attempting to obtain that URL, the following occurs:

1. The client sends a request to the LDNS server to resolve the name ``www-origin-cache.cdn.com`` to an IPv4 address.

2. Similar to the previous case, the LDNS server resolves the name ``www-origin-cache.cdn.com`` to an IPv4 address, in this example, this address is 55.44.33.22.

3. The client opens a TCP connection from a random port locally, to port 80 (the HTTP default) on 55.44.33.22, and sends the following: ::

    GET /foo/bar/fun.html HTTP/1.1
    Host: www-origin-cache.cdn.com

4. The reverse proxy looks up ``www-origin-cache.cdn.com`` in its remap rules, and finds the origin is ``www.origin.com``.

5. The proxy checks its cache to see if the response for ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` is already in the cache.

6a. If the response is not in the cache:

  1. The proxy uses DNS to get the IPv4 address for ``www.origin.com``, connect to it on port 80, and sends: ::

   	GET /foo/bar/fun.html HTTP/1.1
   	Host: www.origin.com

  2. The origin server responds with the headers and content as shown: ::

      HTTP/1.1 200 OK
      Date: Sun, 14 Dec 2014 23:22:44 GMT
      Server: Apache/2.2.15 (Red Hat)
      Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
      ETag: "1aa008f-2d-50a3559482cc0"
      Content-Length: 45
      Connection: close
      Content-Type: text/html; charset=UTF-8

      <html><body>This is a fun file</body></html>

  3. The proxy sends the origin response on to the client adding a ``Via:`` header (and maybe others): ::

      HTTP/1.1 200 OK
      Date: Sun, 14 Dec 2014 23:22:44 GMT
      Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
      ETag: "1aa008f-2d-50a3559482cc0"
      Content-Length: 45
      Connection: close
      Content-Type: text/html; charset=UTF-8
      Age: 0
      Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
      Server: ATS/4.2.1

    	<html><body>This is a fun file</body></html>

6b. If it *is* in the cache:
 
  The proxy responds to the client with the previously retrieved result: ::

      HTTP/1.1 200 OK
      Date: Sun, 14 Dec 2014 23:22:44 GMT
      Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
      ETag: "1aa008f-2d-50a3559482cc0"
      Content-Length: 45
      Connection: close
      Content-Type: text/html; charset=UTF-8
      Age: 39711
      Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
      Server: ATS/4.2.1

      <html><body>This is a fun file</body></html>


.. index::
  Forward Proxy

.. _rl-fwd-proxy:

|arrow| Forward Proxy
---------------------
  A forward proxy acts on behalf of the client. The origin server is mostly 
  unaware of the proxy, the client requests the proxy to retrieve content from a 
  particular origin server. All MID caches in a Traffic Control based CDN are 
  forward proxies. In a forward proxy scenario, the client is explicitely configured  to use the
  the proxy's IP address and port as a forward proxy. The client always connects to the forward 
  proxy for content. The content provider does not have to change the URL the 
  client obtains, and is unaware of the proxy in the middle. 

..  seealso:: `ATS documentation on forward proxy <https://docs.trafficserver.apache.org/en/latest/admin/forward-proxy.en.html>`_.

Below is an example of the client retrieving the URL ``http://www.origin.com/foo/bar/fun.html`` through a forward proxy:

1. The client requires configuration to use the proxy, as opposed to the reverse proxy example. Assume the client configuration is through preferences entries or other to use the proxy IP address 99.88.77.66 and proxy port 8080.

2. To retrieve ``http://www.origin.com/foo/bar/fun.html`` URL, the client connects to 99.88.77.66 on port 8080 and sends: 
 
 ::

  GET http://www.origin.com/foo/bar/fun.html HTTP/1.1


 ..  Note:: In this case, the client places the entire URL after GET, including protocol and hostname (``http://www.origin.com``),  but in the reverse proxy and direct-to-origin case it  puts only the path portion of the URL (``/foo/bar/fun.html``) after the GET. 

3. The proxy verifies whether the response for ``http://www-origin-cache.cdn.com/foo/bar/fun.html`` is already in the cache.

4a. If it is not in the cache:

  1. The proxy uses DNS to obtain the IPv4 address for ``www.origin.com``, connects to it on port 80, and sends: ::


      GET /foo/bar/fun.html HTTP/1.1
      Host: www.origin.com


  2. The origin server responds with the headers and content as shown below: ::


      HTTP/1.1 200 OK
      Date: Sun, 14 Dec 2014 23:22:44 GMT
      Server: Apache/2.2.15 (Red Hat)
      Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
      ETag: "1aa008f-2d-50a3559482cc0"
      Content-Length: 45
      Connection: close
      Content-Type: text/html; charset=UTF-8

      <html><body>This is a fun file</body></html>


  3. The proxy sends this on to the client adding a ``Via:`` header (and maybe others): ::

      HTTP/1.1 200 OK
      Date: Sun, 14 Dec 2014 23:22:44 GMT
      Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
      ETag: "1aa008f-2d-50a3559482cc0"
      Content-Length: 45
      Connection: close
      Content-Type: text/html; charset=UTF-8
      Age: 0
      Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
      Server: ATS/4.2.1
          
      <html><body>This is a fun file</body></html>


4b. If it *is* in the cache:
 
  The proxy responds to the client with the previously retrieved result: ::

    HTTP/1.1 200 OK
    Date: Sun, 14 Dec 2014 23:22:44 GMT
    Last-Modified: Sun, 14 Dec 2014 23:18:51 GMT
    ETag: "1aa008f-2d-50a3559482cc0"
    Content-Length: 45
    Connection: close
    Content-Type: text/html; charset=UTF-8
    Age: 99711
    Via: http/1.1 cache01.cdn.kabletown.net (ApacheTrafficServer/4.2.1 [uScSsSfUpSeN:t cCSi p sS])
    Server: ATS/4.2.1
          
    <html><body>This is a fun file</body></html>

.. index::
  Transparent Proxy
  
|arrow| Transparent Proxy 
-------------------------
  Neither the origin nor the client are aware of the actions performed by the transparent proxies. A Traffic Control based CDN does not use transparent proxies.   If you are interested you can learn more about transparent proxies on `wikipedia <http://en.wikipedia.org/wiki/Proxy_server#Transparent_proxy>`_.

