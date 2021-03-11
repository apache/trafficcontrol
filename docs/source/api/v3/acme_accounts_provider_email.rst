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

.. _to-api-v3-acme-accounts-provider-email:

****************************************
``acme_accounts/{{provider}}/{{email}}``
****************************************

.. versionadded:: 3.1

``DELETE``
==========
Delete :term:`ACME Account` information.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+----------+-----------------------------------------------------------------------------------------------------------------+
	| Name     |                       Description                                                                               |
	+==========+=================================================================================================================+
	| provider | The :abbr:`ACME (Automatic Certificate Management Environment)` provider for the account to be deleted          |
	+----------+-----------------------------------------------------------------------------------------------------------------+
	| email    | The email used in the :term:`ACME Account` to be deleted                                                        |
	+----------+-----------------------------------------------------------------------------------------------------------------+

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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 10 Dec 2020 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: rGD2sOMHYF0sga1zuTytyLHCUkkc3ZwQRKvZ/HuPzObOP4WztKTOVXB4uhs3iJqBg9zRB2TucMxONHN+3/yShQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 10 Dec 2020 14:24:34 GMT
	Content-Length: 60

	{"alerts": [
		{
			"text": "Acme account deleted",
			"level":"success"
		}
	]}
