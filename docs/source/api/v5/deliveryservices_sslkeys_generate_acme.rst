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

.. _to-api-deliveryservices-sslkeys-generate-acme:

******************************************
``deliveryservices/sslkeys/generate/acme``
******************************************

``POST``
========
Generates an SSL certificate and private key using :abbr:`ACME (Automatic Certificate Management Environment)` protocol for a :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: DS-SECURITY-KEY:UPDATE, ACME:READ, DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE
:Response Type:  Object (string)

Request Structure
-----------------
:authType:        The certificate provider correlating to an :abbr:`ACME (Automatic Certificate Management Environment)` account in :ref:`cdn.conf` or Let's Encrypt.
:key:             The :ref:`ds-xmlid` of the :term:`Delivery Service` for which keys will be generated [#needOne]_
:deliveryservice: The :ref:`ds-xmlid` of the :term:`Delivery Service` for which keys will be generated [#needOne]_
:version:         An integer that defines the "version" of the key - which may be thought of as the sequential generation; that is, the higher the number the more recent the key
:hostname:        The desired hostname of the :term:`Delivery Service`

	.. note:: In most cases, this must be the same as the :ref:`ds-example-urls`.

:cdn:             The name of the CDN of the :term:`Delivery Service` for which the certs will be generated

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/deliveryservices/sslkeys/generate/acme HTTP/1.1
	Content-Type: application/json

	{
		"authType": "Lets Encrypt",
		"key": "ds-01",
		"deliveryservice": "ds-01",
		"version": "3",
		"hostname": "tr.ds-01.mycdn.ciab.test",
		"cdn":"test-cdn"
	}


Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [{
		"level": "success",
		"text": "Beginning async ACME call for demo1 using Lets Encrypt. This may take a few minutes. Status updates can be found here: /api/5.0/async_status/1"
	}]}

.. [#needOne] Either the ``key`` or the ``deliveryservice`` field must be provided. If both are provided, then they must match.
