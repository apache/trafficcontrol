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

.. _to-api-v3-vault-bucket-bucket-key-key-values:

**********************************************
``vault/bucket/{{bucket}}/key/{{key}}/values``
**********************************************
Retrieves the `object <https://docs.riak.com/riak/kv/latest/learn/concepts/keys-and-objects/index.html#objects>`_ stored under a given `key <https://docs.riak.com/riak/kv/latest/learn/concepts/keys-and-objects/index.html#keys>`_ from a given `bucket <https://docs.riak.com/riak/kv/latest/learn/concepts/buckets/index.html>`_ in :ref:`tv-overview`.

.. deprecated:: ATCv6

``GET``
=======
:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------+------------------------------------------+
	| Name   | Description                              |
	+========+==========================================+
	| bucket | The bucket that the key is stored under  |
	+--------+------------------------------------------+
	| key    | The key that the object is stored under  |
	+--------+------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/vault/bucket/ssl/key/demo1-latest/values HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
The response structure varies according to what is stored. Top-level keys will always be ``String`` type, but the values can be any type.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Sun, 23 Feb 2020 23:29:27 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: LUq4+MUCgIqxBWqbuA/Pbsdq5Qfs7vPUcZ0Cu2FcOUyP0X8xXit3BJrdOLA+dSSJ3tSQ7Qij1+0ahzuouuLT7Q==
	X-Server-Name: traffic_ops_golang/
	Date: Sun, 23 Feb 2020 22:29:27 GMT
	Transfer-Encoding: chunked

	{
		"alerts": [
		{
			"text": "This endpoint is deprecated, and will be removed in the future",
			"level": "warning"
		}],
		"response": {
			"cdn": "CDN-in-a-Box",
			"certificate": {
				"crt": "...",
				"csr": "...",
				"key": "..."
			},
			"deliveryservice": "demo1",
			"hostname": "*.demo1.mycdn.ciab.test",
			"key": "demo1",
			"version": 1
		}
	}
