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

.. _to-api-v3-user-login-oauth:

********************
``user/login/oauth``
********************

``POST``
========
Authentication of a user by exchanging a code for an encrypted JSON Web Token from an OAuth service. Traffic Ops will ``POST`` to the ``authCodeTokenUrl`` to exchange the code for an encrypted JSON Web Token.  It will then decode and validate the token, validate the key set domain, and send back a session cookie.

:Auth. Required: No
:Roles Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
:authCodeTokenUrl: URL for code-to-token conversion
:code: Code
:clientId: Client Id
:redirectUri: Redirect URI

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/user/login/oauth HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 26
	Content-Type: application/json

	{
		"authCodeTokenUrl": "https://url-to-convert-code-to-token.example.com",
		"code": "AbCd123",
		"clientId": "oauthClientId",
		"redirectUri": "https://traffic-portal.example.com/sso"
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
	Whole-Content-Sha512: UdO6T3tMNctnVusDXzRjVwwYOnD7jmnBzPEB9PvOt2bHajTv3SKTPiIZjDzvhU6EX4p+JoG4fA5wlhgxpsejIw==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 13 Dec 2018 15:21:33 GMT
	Content-Length: 65

	{ "alerts": [
		{
			"text": "Successfully logged in.",
			"level": "success"
		}
	]}
