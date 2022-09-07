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

.. _to-api-v4-acme-accounts-providers:

***************************
``acme_accounts/providers``
***************************

.. versionadded:: 4.0

``GET``
=======
Gets a list of all :abbr:`ACME (Automatic Certificate Management Environment)` providers set up in :ref:`cdn.conf` and Let's Encrypt.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ACME:READ
:Response Type:  Array

Request Structure
-----------------
No parameters available


Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": [
		"CertAuth1",
		"CertAuth2",
		"CertAuth3"
	]}
