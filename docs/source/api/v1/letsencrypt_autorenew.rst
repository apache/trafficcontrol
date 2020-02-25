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

.. _to-api-v1-letsencrypt-autorenew:

*************************
``letsencrypt/autorenew``
*************************

.. versionadded:: 1.5

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
:LetsEncryptExpirations: A list of objects with information regarding certificate expiration for all Let's Encrypt certificates

	:XmlId:       The :term:`Delivery Service`'s uniquely identifying :ref:`ds-xmlid`
	:Version:     An integer that defines the "version" of the key - which may be thought of as the sequential generation; that is, the higher the number the more recent the key
	:Expiration:  The expiration date of the certificate for the :term:`Delivery Service` in :rfc:`3339` format
	:AuthType:    The authority type of the certificate for the :term:`Delivery Service`
	:Error:       Any errors received in the renewal process

:SelfSignedExpirations:  A list of objects with information regarding certificate expiration for all self signed certificates

	:XmlId:       The :term:`Delivery Service`'s uniquely identifying :ref:`ds-xmlid`
	:Version:     An integer that defines the "version" of the key - which may be thought of as the sequential generation; that is, the higher the number the more recent the key
	:Expiration:  The expiration date of the certificate for the :term:`Delivery Service` in :rfc:`3339` format
	:AuthType:    The authority type of the certificate for the :term:`Delivery Service`
	:Error:       Any errors received in the renewal process

:OtherExpirations:       A list of objects with information regarding certificate expiration for all other certificates

	:XmlId:       The :term:`Delivery Service`'s uniquely identifying :ref:`ds-xmlid`
	:Version:     An integer that defines the "version" of the key - which may be thought of as the sequential generation; that is, the higher the number the more recent the key
	:Expiration:  The expiration date of the certificate for the :term:`Delivery Service` in :rfc:`3339` format
	:AuthType:    The authority type of the certificate for the :term:`Delivery Service`
	:Error:       Any errors received in the renewal process

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": {
		"LetsEncryptExpirations": [
			{
				"XmlId":"demo2",
				"Version":1,
				"Expiration":"2020-08-18T13:53:06Z",
				"AuthType":"Lets Encrypt",
				"Error":null
			}
		],
		"SelfSignedExpirations": [
			{
				"XmlId":"demo1",
				"Version":3,
				"Expiration":"2020-08-18T13:53:06Z",
				"AuthType":"Self Signed",
				"Error":null
			}
		],
		"OtherExpirations":null
	}}
