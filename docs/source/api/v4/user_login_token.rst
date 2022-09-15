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

.. _to-api-v4-user-login-token:

********************
``user/login/token``
********************

``POST``
========
Authentication of a user using a token. Normally, the token is obtained via a call to either :ref:`to-api-v4-user-reset_password` or :ref:`to-api-v4-users-register`.

:Auth. Required: No
:Roles Required: None
:Permissions Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
:t: A :abbr:`UUID (Universal Unique Identifier)` generated for the user.

	.. impl-detail:: Though not strictly necessary for authentication provided direct database access, the tokens generated for use with this endpoint are compliant with :RFC:`4122`.

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/user/login/token HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 44
	Content-Type: application/json

	{
		"t": "18EE200C-FF24-11E8-BF01-870C776752A3"
	}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: FuS3TkVosxHtpxRGMJ2on+WnFdYTNSPjxz/Gh1iT4UCJ2/P0twUbAGQ3tTx9EfGiAzg9CNQiVUFGnYjJZ6NCpg==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 20 Sep 2019 15:02:43 GMT
	Content-Length: 66

	{ "alerts": [
		{
			"text": "Successfully logged in.",
			"level": "success"
		}
	]}
