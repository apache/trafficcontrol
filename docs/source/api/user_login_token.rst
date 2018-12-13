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

.. _to-api-user-login-token:

********************
``user/login/token``
********************
.. caution:: This page is a stub! Much of it may be missing or just downright wrong - it needs a lot of love from people with the domain knowledge required to update it.

``POST``
========
Authentication of a user using a token. Normally, the token is obtained via a call to either :ref:`to-api-user-reset_password` or :ref:`to-api-users-register`.

:Auth. Required: No
:Roles Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
:t: The login token

.. code-block:: http
	:caption: Request Example

	POST /api/1.3/user/login/token HTTP/1.1
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 13 Dec 2018 22:16:25 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Fri, 14 Dec 2018 02:16:25 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: uDowfYsW7ADmZyfahD21A+KuDdycQ3a4ma5kbPO/9RXsvgL9bqNC0Ocpi4QLxJN1Ffe1jroYoiqcnjlK9KX/5Q==
	Content-Length: 65

	{ "alerts": [
		{
			"level": "success",
			"text": "Successfully logged in."
		}
	]}
