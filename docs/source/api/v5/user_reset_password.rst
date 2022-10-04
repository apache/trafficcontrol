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

.. _to-api-user-reset_password:

***********************
``user/reset_password``
***********************

``POST``
========
Sends an email to reset a user's password.

:Auth. Required: No
:Roles Required: None
:Permissions Required: None
:Response Type:  ``undefined``

Request Structure
-----------------
:email: The email address of the user to initiate password reset

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/user/reset_password HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 35
	Content-Type: application/json

	{
		"email": "test@example.com"
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
	Date: Thu, 13 Dec 2018 22:11:53 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: lKWwVYbgKklk7bljnQJZwWV5bljIk+GkooD6SAc3CSexVKvfbL9dgL5iBc/BNNRk2pIU5I/1GgldcDLrXsF1ZA==
	Content-Length: 109

	{ "alerts": [
		{
			"level": "success",
			"text": "Successfully sent password reset to email 'test@example.com'"
		}
	]}
