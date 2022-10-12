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

.. _to-api-v3-letsencrypt-autorenew:

*************************
``letsencrypt/autorenew``
*************************

``POST``
========
Generates an SSL certificate and private key using Let's Encrypt for a :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
No parameters available


Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "alerts": [
		{
			"text": "Beginning async call to renew certificates. This may take a few minutes.",
			"level": "success"
		}
	]}
