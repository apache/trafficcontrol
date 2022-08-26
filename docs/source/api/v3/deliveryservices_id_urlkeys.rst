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

.. _to-api-v3-deliveryservices-id-urlkeys:

***********************************
``deliveryservices/{{ID}}/urlkeys``
***********************************

``GET``
=======
.. seealso:: :ref:`to-api-v3-deliveryservices-xmlid-xmlid-urlkeys`

Retrieves URL signing keys for a :term:`Delivery Service`.

.. caution:: This method will return the :term:`Delivery Service`'s **PRIVATE** URL signing keys! Be wary of using this endpoint and **NEVER** share the output with anyone who would be unable to see it on their own.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------------------------------------------+
	| Name | Description                                                                            |
	+======+========================================================================================+
	| id   | Filter for the :term:`Delivery Service` identified by this integral, unique identifier |
	+------+----------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/deliveryservices/1/urlkeys HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:key<N>: The private URL signing key for this :term:`Delivery Service` as a base-64-encoded string, where ``<N>`` is the "generation" of the key e.g. the first key will always be named ``"key0"``. Up to 16 concurrent generations are retained at any time (``<N>`` is always on the interval [0,15])

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Sun, 23 Feb 2020 16:34:56 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: cTc5OPE3hM+CiyQPCy36zD2tsQcfkvIqQ7/D82WGMWHm+ACW3YbcKhgPnSQU6+Tuj4jya52Kx9+nw5+OonFvPQ==
	X-Server-Name: traffic_ops_golang/
	Date: Sun, 23 Feb 2020 15:34:56 GMT
	Content-Length: 533

	{
		"response": {
			"key0": "...",
			"key1": "...",
			"key2": "...",
			"key3": "...",
			"key4": "...",
			"key5": "...",
			"key6": "...",
			"key7": "...",
			"key8": "...",
			"key9": "...",
			"key10": "...",
			"key11": "...",
			"key12": "...",
			"key13": "...",
			"key14": "...",
			"key15": "..."
		}
	}
