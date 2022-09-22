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

.. _to-api-acme-autorenew:

******************
``acme_autorenew``
******************

``POST``
========
Generates SSL certificates and private keys for all :term:`Delivery Services` that have certificates expiring within the configured time. This uses:abbr:`ACME (Automatic Certificate Management Environment)` or Let's Encrypt depending on the certificate.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ACME:READ, DS-SECURITY-KEY:UPDATE, DELIVERY-SERVICE:UPDATE
:Response Type:  ``undefined``

Request Structure
-----------------
No parameters available


Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 202 Accepted
	Content-Type: application/json

	{ "alerts": [
		{
			"text": "Beginning async call to renew certificates. This may take a few minutes. Status updates can be found here: /api/5.0/async_status/1",
			"level": "success"
		}
	]}
