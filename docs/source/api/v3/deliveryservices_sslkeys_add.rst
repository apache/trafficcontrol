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

.. _to-api-v3-deliveryservices-sslkeys-add:

********************************
``deliveryservices/sslkeys/add``
********************************

.. seealso:: In most cases it is preferable to allow Traffic Ops to generate the keys via :ref:`to-api-v3-deliveryservices-sslkeys-generate`, rather than uploading them manually using this endpoint.

``POST``
========
Allows user to upload an SSL certificate, csr, and private key for a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object (string)

Request Structure
-----------------
:cdn:         The name of the CDN to which the :term:`Delivery Service` belongs
:certificate: An object that contains the actual components of the SSL key

	:crt: The certificate for the :term:`Delivery Service` identified by ``key``
	:csr: The csr file for the :term:`Delivery Service` identified by ``key``
	:key: The private key for the :term:`Delivery Service` identified by ``key``

:hostname:        The desired hostname of the :term:`Delivery Service`

	.. note:: In most cases, this must be the same as the :term:`Delivery Service` URL'

:key:     The :ref:`ds-xmlid` of the :term:`Delivery Service` to which these keys will be assigned
:version: An integer that defines the "version" of the key - which may be thought of as the sequential generation; that is, the higher the number the more recent the key

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/deliveryservices/sslkeys/add HTTP/1.1
	Host: trafficops.infra.ciab.test
	Content-Type: application/json

	{
		"cdn": "test-cdn",
		"hostname": "tr.ds-01.mycdn.ciab.test",
		"key": "ds-01",
		"version": "1",
		"certificate": {
			"key": "some_key",
			"csr": "some_csr",
			"crt": "some_crt"
		}
	}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{
		"response": "Successfully added ssl keys for ds-01"
	}
