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

.. _to-api-v3-cdns-name-name-sslkeys:

******************************
``cdns/name/{{name}}/sslkeys``
******************************

``GET``
=======
Returns SSL certificates for all :term:`Delivery Services` that are a part of the CDN.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------+
	| Name | Description                                        |
	+======+====================================================+
	| name | The name of the CDN for which keys will be fetched |
	+------+----------------------------------------------------+

Response Structure
------------------
:certificate: An object representing The SSL keys used for the :term:`Delivery Service` identified by ``deliveryservice``

	:key: Base 64-encoded private key for SSL certificate
	:crt: Base 64-encoded SSL certificate

:deliveryservice: A string that is the :ref:`ds-xmlid` of the :term:`Delivery Service` using the SSL key within ``certificate``

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"deliveryservice": "ds1",
			"certificate": {
				"crt": "base64encodedcrt1",
				"key": "base64encodedkey1"
			}
		},
		{
			"deliveryservice": "ds2",
			"certificate": {
				"crt": "base64encodedcrt2",
				"key": "base64encodedkey2"
			}
		}
	]}
