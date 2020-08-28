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

.. _to-api-v1-cachegroup-parameterID-parameter:

*****************************************
``cachegroup/{{parameter ID}}/parameter``
*****************************************
.. deprecated:: ATCv4
.. danger:: This endpoint does not appear to work, and thus its use is strongly discouraged!

``GET``
=======
Extract identifying information about all :term:`Cache Groups` with a specific :term:`Parameter`

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------------+----------+-----------------------+
	| Name         | Required | Description           |
	+==============+==========+=======================+
	| parameter_ID | yes      | A :ref:`parameter-id` |
	+--------------+----------+-----------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/cachegroup/1/parameter HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:cachegroups: An array of all :term:`Cache Groups` with an associated :term:`Parameter` identifiable by the ``parameter_id`` request path parameter

	:id:   An integer that is the :term:`Cache Group`'s :ref:`cache-group-id`
	:name: A string that is the :ref:`cache-group-name` of the :term:`Cache Group`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Encoding: gzip
	Content-Length: 161
	Content-Type: application/json
	Date: Tue, 03 Dec 2019 15:15:26 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 03 Dec 2019 19:15:26 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: H03AKuJ2IjG3wb6SEplNtIjm8ka3JJdRxc2HyOkNzjHdsh8p7UcJ1teYvYUf8yMNDt8HHBaKzIDoHODLwhktjA==

	{
		"alerts": [
			{
				"level": "warning",
				"text": "This endpoint is deprecated, please use 'GET /cachegroupparameters & GET /cachegroups' instead"
			}
		],
		"response": {
			"cachegroups": [
				{
					"name": "CDN_in_a_Box_Edge",
					"id": 7
				},
				{
					"name": "CDN_in_a_Box_Mid",
					"id": 6
				},
				{
					"name": "TRAFFIC_ANALYTICS",
					"id": 1
				},
				{
					"name": "TRAFFIC_OPS",
					"id": 2
				},
				{
					"name": "TRAFFIC_OPS_DB",
					"id": 3
				},
				{
					"name": "TRAFFIC_PORTAL",
					"id": 4
				},
				{
					"name": "TRAFFIC_STATS",
					"id": 5
				}
			]
	}}
