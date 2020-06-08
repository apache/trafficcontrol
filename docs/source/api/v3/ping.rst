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

.. _to-api-v3-ping:

********
``ping``
********
Checks whether Traffic Ops is online.

``GET``
=======
:Auth. Required: No
:Response Type:  ``undefined``

Request Structure
-----------------
No parameters available.

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/ping HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:ping:          Returns an object containing only the ``"ping"`` property, which always has the value ``"pong"``.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Whole-Content-Sha512: 0HqcLcYHCB4AFYGFzcAsP2h+PCMlYxk/TqMajcy3fWCzY730Tv32k5trUaJLeSBbRx2FUi7z/sTAuzikdg0E4g==
	X-Server-Name: traffic_ops_golang/
	Date: Sun, 23 Feb 2020 20:52:01 GMT
	Content-Length: 40
	Content-Type: application/x-gzip

	{
		"ping": "pong"
	}
