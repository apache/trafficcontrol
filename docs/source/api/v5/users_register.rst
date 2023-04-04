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

.. _to-api-users-register:

******************
``users/register``
******************

``POST``
========
Register a user and send registration email.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: USER:CREATE, USER:READ
:Response Type:  ``undefined``

Request Structure
-----------------
:email:    Email address of the new user The given email is validated (circuitously) by `GitHub user asaskevich's regular expression <https://github.com/asaskevich/govalidator/blob/9a090521c4893a35ca9a228628abf8ba93f63108/patterns.go#L7>`_ . Note that it can't actually distinguish a valid, deliverable, email address but merely ensure the email is in a commonly-found format.
:role:     The name of the :term:`Role` which will be afforded to the new user. It restricts the allowed values to identifiers for :term:`Roles` with at only Permissions the requesting user has.
:tenantId: A field containing the integral, unique identifier of the :term:`Tenant` to which the new user will belong. It restricts the allowed values to identifiers for :term:`Tenants` within the requesting user's :term:`Tenant`'s permissions.

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/users/register HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 59
	Content-Type: application/json

	{
		"email": "test@example.com",
		"role": "admin",
		"tenantId": 1
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
	Date: Thu, 13 Dec 2018 22:03:22 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: yvf++Oqxvu3uOIAYbWLUgJKxZ4T60Mi5H9eGTxrKLxnRsHw0PdDIrThbTnWtATBkak4vU/dPHLLXKW85LUTEWg==
	Content-Length: 160

	{ "alerts": [
		{
			"level": "success",
			"text": "Sent user registration to test@example.com with the following permissions [ role: read-only | tenant: root ]"
		}
	]}
