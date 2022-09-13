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

.. _to-api-user-logout:

***************
``user/logout``
***************

``POST``
========
User logout. Invalidates the session cookie of the currently logged-in user.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
No parameters available

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 13 Dec 2018 21:25:36 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 6KEdr1ZC512zkOl03KwvQE0L7qrJ/+ek6ztymkYy9p8gdPUyYyzGEAJ/Ldb8GY0UBFYmgqeZ3yWHvTcEsOQMiw==
	Content-Length: 61

	{ "alerts": [
		{
			"level": "success",
			"text": "You are logged out."
		}
	]}
