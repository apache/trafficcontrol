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

.. _to-api-sslkey_expirations:

**********************
``sslkey_expirations``
**********************

``GET``
=======
Retrieves SSL certificate expiration information.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: ACME:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                            |
	+===================+==========+========================================================================================================+
	| days              | no       | Return only the expiration information for SSL certificates expiring in the next given number of days. |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/sslkey_expirations?days=30 HTTP/1.1
	Host: trafficops.infra.ciab.test
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:deliveryservice:   The :ref:`ds-xmlid` for the :term:`Delivery Service` corresponding to this SSL certificate.
:cdn:               The ID for the :abbr:`CDN (Content Delivery Network)` corresponding to this SSL certificate.
:provider:          The provider of this SSL certificate, generally the name of the Certificate Authority or the :abbr:`ACME (Automatic Certificate Management Environment)` account.
:expiration:        The expiration date of this SSL certificate.
:federated:         A boolean indicating if this SSL certificate is use in a federated :term:`Delivery Service`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 07 Jun 2021 22:52:20 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Jun 2021 21:52:20 GMT
	Content-Length: 384

	{ "response": [
		{
			"deliveryservice": "foo1",
			"cdn": "cdn1",
			"provider": "Self Signed",
			"expiration": "2022-08-02T15:38:06-06:00",
			"federated": false
		},
		{
			"deliveryservice": "foo2",
			"cdn": "cdn2",
			"provider": "Lets Encrypt",
			"expiration": "2022-07-12T12:14:00-06:00",
			"federated": true
		}
	]}
