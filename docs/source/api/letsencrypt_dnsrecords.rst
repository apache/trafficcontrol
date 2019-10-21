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

.. _to-letsencrypt-dnsrecords:

**************************
``letsencrypt/dnsrecords``
**************************

.. versionadded:: 1.4

``GET``
========
Gets all DNS challenge records for Let's Encrypt DNS challenges.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
No parameters available


Response Structure
------------------
:fqdn:      The Fully Qualified Domain Name (FQDN) for the TXT record location to complete the DNS challenge
:record:    The record provided by Let's Encrypt to complete the DNS challenge

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": [
		{
			"fqdn":"_acme-challenge.demo1.example.com.",
			"record":"testRecord"
		}
	]}
