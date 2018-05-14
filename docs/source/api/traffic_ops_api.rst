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

API Overview
************
The Traffic Ops API provides programmatic access to read and write CDN data providing authorized API consumers with the ability to monitor CDN performance and configure CDN settings and parameters.

Response Structure
------------------
All successful responses have the following structure: ::

    {
      "response": <JSON object with main response>,
    }

To make the documentation easier to read, only the ``<JSON object with main response>`` is documented, even though the response and version fields are always present. 

Using API Endpoints
-------------------
1. Authenticate with your Traffic Portal or Traffic Ops user account credentials.
2. Upon successful user authentication, note the mojolicious cookie value in the response headers. 
3. Pass the mojolicious cookie value, along with any subsequent calls to an authenticated API endpoint.

Example: ::
  
    [jvd@laika ~]$ curl -H "Accept: application/json" http://localhost:3000/api/1.1/usage/asns.json
    {"alerts":[{"level":"error","text":"Unauthorized, please log in."}]}
    [jvd@laika ~]$
    [jvd@laika ~]$ curl -v -H "Accept: application/json" -v -X POST --data '{ "u":"admin", "p":"secret_passwd" }' http://localhost:3000/api/1.1/user/login
    * Hostname was NOT found in DNS cache
    *   Trying ::1...
    * connect to ::1 port 3000 failed: Connection refused
    *   Trying 127.0.0.1...
    * Connected to localhost (127.0.0.1) port 3000 (#0)
    > POST /api/1.1/user/login HTTP/1.1
    > User-Agent: curl/7.37.1
    > Host: localhost:3000
    > Accept: application/json
    > Content-Length: 32
    > Content-Type: application/x-www-form-urlencoded
    >
    * upload completely sent off: 32 out of 32 bytes
    < HTTP/1.1 200 OK
    < Connection: keep-alive
    < Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
    < Access-Control-Allow-Origin: http://localhost:8080
    < Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
    < Set-Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862; expires=Sun, 19 Apr 2015 00:10:01 GMT; path=/; HttpOnly
    < Content-Type: application/json
    < Date: Sat, 18 Apr 2015 20:10:01 GMT
    < Access-Control-Allow-Credentials: true
    < Content-Length: 81
    < Cache-Control: no-cache, no-store, max-age=0, must-revalidate
    * Server Mojolicious (Perl) is not blacklisted
    < Server: Mojolicious (Perl)
    <
    * Connection #0 to host localhost left intact
    {"alerts":[{"level":"success","text":"Successfully logged in."}]}
    [jvd@laika ~]$

    [jvd@laika ~]$ curl -H'Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862;' -H "Accept: application/json" http://localhost:3000/api/1.1/asns.json
    {"response":{"asns":[{"lastUpdated":"2012-09-17 15:41:22", .. asn data deleted ..   ,}
    [jvd@laika ~]$

API Errors
----------

**Response Properties**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
|``alerts``            | array  | A collection of alert messages.                |
+----------------------+--------+------------------------------------------------+
| ``>level``           | string | Success, info, warning or error.               |
+----------------------+--------+------------------------------------------------+
| ``>text``            | string | Alert message.                                 |
+----------------------+--------+------------------------------------------------+

The 3 most common errors returned by Traffic Ops are:

401 Unauthorized
  When you don't supply the right cookie, this is the response. :: 

    [jvd@laika ~]$ curl -v -H "Accept: application/json" http://localhost:3000/api/1.1/usage/asns.json
    * Hostname was NOT found in DNS cache
    *   Trying ::1...
    * connect to ::1 port 3000 failed: Connection refused
    *   Trying 127.0.0.1...
    * Connected to localhost (127.0.0.1) port 3000 (#0)
    > GET /api/1.1/usage/asns.json HTTP/1.1
    > User-Agent: curl/7.37.1
    > Host: localhost:3000
    > Accept: application/json
    >
    < HTTP/1.1 401 Unauthorized
    < Cache-Control: no-cache, no-store, max-age=0, must-revalidate
    < Content-Length: 84
    * Server Mojolicious (Perl) is not blacklisted
    < Server: Mojolicious (Perl)
    < Connection: keep-alive
    < Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
    < Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
    < Access-Control-Allow-Origin: http://localhost:8080
    < Date: Sat, 18 Apr 2015 20:36:12 GMT
    < Content-Type: application/json
    < Access-Control-Allow-Credentials: true
    <
    * Connection #0 to host localhost left intact
    {"alerts":[{"level":"error","text":"Unauthorized, please log in."}]}
    [jvd@laika ~]$

404 Not Found
  When the resource (path) is non existent Traffic Ops returns a 404::

    [jvd@laika ~]$ curl -v -H'Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862;' -H "Accept: application/json" http://localhost:3000/api/1.1/asnsjj.json
    * Hostname was NOT found in DNS cache
    *   Trying ::1...
    * connect to ::1 port 3000 failed: Connection refused
    *   Trying 127.0.0.1...
    * Connected to localhost (127.0.0.1) port 3000 (#0)
    > GET /api/1.1/asnsjj.json HTTP/1.1
    > User-Agent: curl/7.37.1
    > Host: localhost:3000
    > Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862;
    > Accept: application/json
    >
    < HTTP/1.1 404 Not Found
    * Server Mojolicious (Perl) is not blacklisted
    < Server: Mojolicious (Perl)
    < Content-Length: 75
    < Cache-Control: no-cache, no-store, max-age=0, must-revalidate
    < Content-Type: application/json
    < Date: Sat, 18 Apr 2015 20:37:43 GMT
    < Access-Control-Allow-Credentials: true
    < Set-Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAzODYzLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--8a5a61b91473bc785d4073fe711de8d2c63f02dd; expires=Sun, 19 Apr 2015 00:37:43 GMT; path=/; HttpOnly
    < Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
    < Connection: keep-alive
    < Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
    < Access-Control-Allow-Origin: http://localhost:8080
    <
    * Connection #0 to host localhost left intact
    {"alerts":[{"text":"Resource not found.","level":"error"}]}
    [jvd@laika ~]$

500 Internal Server Error
  When you are asking for a correct path, but the database doesn't match, it returns a 500:: 

    [jvd@laika ~]$ curl -v -H'Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862;' -H "Accept: application/json" http://localhost:3000/api/1.1/servers/hostname/jj/details.json
    * Hostname was NOT found in DNS cache
    *   Trying ::1...
    * connect to ::1 port 3000 failed: Connection refused
    *   Trying 127.0.0.1...
    * Connected to localhost (127.0.0.1) port 3000 (#0)
    > GET /api/1.1/servers/hostname/jj/details.json HTTP/1.1
    > User-Agent: curl/7.37.1
    > Host: localhost:3000
    > Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862;
    > Accept: application/json
    >
    < HTTP/1.1 500 Internal Server Error
    * Server Mojolicious (Perl) is not blacklisted
    < Server: Mojolicious (Perl)
    < Cache-Control: no-cache, no-store, max-age=0, must-revalidate
    < Content-Length: 93
    < Set-Cookie: mojolicious=eyJhdXRoX2RhdGEiOiJhZG1pbiIsImV4cGlyZXMiOjE0Mjk0MDQzMDZ9--1b08977e91f8f68b0ff5d5e5f6481c76ddfd0853; expires=Sun, 19 Apr 2015 00:45:06 GMT; path=/; HttpOnly
    < Content-Type: application/json
    < Date: Sat, 18 Apr 2015 20:45:06 GMT
    < Access-Control-Allow-Credentials: true
    < Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
    < Connection: keep-alive
    < Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
    < Access-Control-Allow-Origin: http://localhost:8080
    <
    * Connection #0 to host localhost left intact
    {"alerts":[{"level":"error","text":"An error occurred. Please contact your administrator."}]}
    [jvd@laika ~]$

  The rest of the API documentation will only document the ``200 OK`` case, where no errors have occured.

TrafficOps Native Client Libraries
----------------------------------

TrafficOps client libraries are available in both Golang and Python.  You can read more about them at https://github.com/apache/incubator-trafficcontrol/tree/master/traffic_control/clients
