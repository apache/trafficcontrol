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

.. _to-api-v3-deliveryservices-sslkeys-generate:

*************************************
``deliveryservices/sslkeys/generate``
*************************************

``POST``
========
Generates an SSL certificate, csr, and private key for a :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object (string)

Request Structure
-----------------
:authType:        The certificate provider correlating to an :abbr:`ACME (Automatic Certificate Management Environment)` account in :ref:`cdn.conf` or Let's Encrypt.
:businessUnit:    A required field which will represent the business unit for which the SSL certificate was generated
:cdn:             A required field which will represent the CDN of the :term:`Delivery Service` for which keys will be generated
:city:            A required field which will represent the resident city of the generated SSL certificate
:country:         A required field which will represent the resident country of the generated SSL certificate
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` for which keys will be generated
:hostname:        The desired hostname of the :term:`Delivery Service`

	.. note:: In most cases, this must be the same as the :term:`Delivery Service` URL'

:key:             The :ref:`ds-xmlid` of the :term:`Delivery Service` for which keys will be generated
:organization:    A required field which will represent the organization for which the SSL certificate was generated
:state:           A required field which will represent the resident state or province of the generated SSL certificate
:version:         An integer that defines the "version" of the key - which may be thought of as the sequential generation; that is, the higher the number the more recent the key

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/deliveryservices/sslkeys/generate HTTP/1.1
	Content-Type: application/json

	{
		"cdn": "test-cdn",
		"key": "ds-01",
		"businessUnit": "CDN Engineering",
		"version": "3",
		"hostname": "tr.ds-01.mycdn.ciab.test",
		"country": "US",
		"organization": "CDN",
		"city": "Denver",
		"state": "Colorado"
	}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": "Successfully created ssl keys for ds-01" }
