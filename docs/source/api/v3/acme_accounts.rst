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

.. _to-api-v3-acme-accounts:

*****************
``acme_accounts``
*****************

.. versionadded:: 3.1

``GET``
=======
Gets information for all :term:`ACME Account` s.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Array

Request Structure
-----------------
No parameters available


Response Structure
------------------
:email:       The email connected to the :term:`ACME Account`.
:privateKey:  The private key connected to the :term:`ACME Account`.
:uri:         The URI for the :term:`ACME Account`. Differs per provider.
:provider:    The :abbr:`ACME (Automatic Certificate Management Environment)` provider.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "response": [
		{
			"email": "sample@example.com",
			"privateKey": "-----BEGIN RSA PRIVATE KEY-----\nSampleKey\n-----END RSA PRIVATE KEY-----\n",
			"uri": "https://acme.example.com/acct/1",
			"provider": "Lets Encrypt"
		}
	]}


``POST``
========
Creates a new :term:`ACME Account`.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
The request body must be a single :term:`ACME Account` object with the following keys:

:email:       The email connected to the :term:`ACME Account`.
:privateKey:  The private key connected to the :term:`ACME Account`.
:uri:         The URI for the :term:`ACME Account`. Differs per provider.
:provider:    The :abbr:`ACME (Automatic Certificate Management Environment)` provider.

.. code-block:: http
	:caption: Request Example

	POST /api/3.1/acme_accounts HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 181
	Content-Type: application/json

	{
		"email": "sample@example.com",
		"privateKey": "-----BEGIN RSA PRIVATE KEY-----\nSampleKey\n-----END RSA PRIVATE KEY-----\n",
		"uri": "https://acme.example.com/acct/1",
		"provider": "Lets Encrypt"
	}

Response Structure
------------------
:email:       The email connected to the :term:`ACME Account`.
:privateKey:  The private key connected to the :term:`ACME Account`.
:uri:         The URI for the :term:`ACME Account`. Differs per provider.
:provider:    The :abbr:`ACME (Automatic Certificate Management Environment)` provider.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 201 Created
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 10 Dec 2020 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 19:18:21 GMT
	Content-Length: 253

	{ "alerts": [
		{
			"text": "Acme account created",
			"level":"success"
		}
	],
	"response": {
		"email": "sample@example.com",
		"privateKey": "-----BEGIN RSA PRIVATE KEY-----\nSampleKey\n-----END RSA PRIVATE KEY-----\n",
		"uri": "https://acme.example.com/acct/1",
		"provider": "Lets Encrypt"
	}}


``PUT``
=======
Updates an existing :term:`ACME Account`.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
The request body must be a single :term:`ACME Account` object with the following keys:

:email:       The email connected to the :term:`ACME Account`.
:privateKey:  The private key connected to the :term:`ACME Account`.
:uri:         The URI for the :term:`ACME Account`. Differs per provider.
:provider:    The :abbr:`ACME (Automatic Certificate Management Environment)` provider.

.. code-block:: http
	:caption: Request Example

	PUT /api/3.1/acme_accounts HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 181
	Content-Type: application/json

	{
		"email": "sample@example.com",
		"privateKey": "-----BEGIN RSA PRIVATE KEY-----\nSampleKey\n-----END RSA PRIVATE KEY-----\n",
		"uri": "https://acme.example.com/acct/1",
		"provider": "Lets Encrypt"
	}

Response Structure
------------------
:email:       The email connected to the :term:`ACME Account`.
:privateKey:  The private key connected to the :term:`ACME Account`.
:uri:         The URI for the :term:`ACME Account`. Differs per provider.
:provider:    The :abbr:`ACME (Automatic Certificate Management Environment)` provider.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 10 Dec 2020 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 19:18:21 GMT
	Content-Length: 253

	{ "alerts": [
		{
			"text": "Acme account updated",
			"level":"success"
		}
	],
	"response": {
		"email": "sample@example.com",
		"privateKey": "-----BEGIN RSA PRIVATE KEY-----\nSampleKey\n-----END RSA PRIVATE KEY-----\n",
		"uri": "https://acme.example.com/acct/1",
		"provider": "Lets Encrypt"
	}}
