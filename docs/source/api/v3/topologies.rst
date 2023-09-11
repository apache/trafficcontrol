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

.. _to-api-v3-topologies:

**************
``topologies``
**************

``GET``
=======
Retrieves :term:`Topologies`.

:Auth. Required: Yes
:Roles Required: "read-only"
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+-----------------------------------------------------+
	| Name | Required | Description                                         |
	+======+==========+=====================================================+
	| name | no       | Return the :term:`Topology` with this name          |
	+------+----------+-----------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/topologies HTTP/1.1
	User-Agent: python-requests/2.23.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:description:           A short sentence that describes the :term:`Topology`.
:lastUpdated:           The date and time at which this :term:`Topology` was last updated, in ISO-like format
:name:                  The name of the :term:`Topology`. This can only be letters, numbers, and dashes.
:nodes:                 An array of nodes in the :term:`Topology`

	:cachegroup:            The name of a :term:`Cache Group`
	:parents:               The indices of the parents of this node in the nodes array, 0-indexed.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 13 Apr 2020 18:22:32 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: lF4MCJCinuQWz0flLAAZBrzbuPVsHrNn2BtTozRZojEjGpm76IsXBQK5QOwSwBoHac+D0C1Z3p7M8kdjcfgIIg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 13 Apr 2020 17:22:32 GMT
	Content-Length: 205

	{
		"response": [
			{
				"description": "This is my topology",
				"name": "my-topology",
				"nodes": [
					{
						"cachegroup": "edge1",
						"parents": [
							7
						]
					},
					{
						"cachegroup": "edge2",
						"parents": [
							7,
							8
						]
					},
					{
						"cachegroup": "edge3",
						"parents": [
							8,
							9
						]
					},
					{
						"cachegroup": "edge4",
						"parents": [
							9
						]
					},
					{
						"cachegroup": "mid1",
						"parents": []
					},
					{
						"cachegroup": "mid2",
						"parents": [
							4
						]
					},
					{
						"cachegroup": "mid3",
						"parents": [
							4
						]
					},
					{
						"cachegroup": "mid4",
						"parents": [
							5
						]
					},
					{
						"cachegroup": "mid5",
						"parents": [
							5,
							6
						]
					},
					{
						"cachegroup": "mid6",
						"parents": [
							6
						]
					}
				],
				"lastUpdated": "2020-04-13 17:12:34+00"
			}
		]
	}

``POST``
========
Create a new :term:`Topology`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:description:           A short sentence that describes the topology.
:name:                  The name of the topology. This can only be letters, numbers, and dashes.
:nodes:                 An array of nodes in the :term:`Topology`

	:cachegroup:            The name of a :term:`Cache Group` with at least 1 server in it
	:parents:               The indices of the parents of this node in the nodes array, 0-indexed.

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/topologies HTTP/1.1
	User-Agent: python-requests/2.23.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 924
	Content-Type: application/json

	{
		"name": "my-topology",
		"description": "This is my topology",
		"nodes": [
			{
				"cachegroup": "edge1",
				"parents": [
					7
				]
			},
			{
				"cachegroup": "edge2",
				"parents": [
					7,
					8
				]
			},
			{
				"cachegroup": "edge3",
				"parents": [
					8,
					9
				]
			},
			{
				"cachegroup": "edge4",
				"parents": [
					9
				]
			},
			{
				"cachegroup": "mid1",
				"parents": []
			},
			{
				"cachegroup": "mid2",
				"parents": [
					4
				]
			},
			{
				"cachegroup": "mid3",
				"parents": [
					4
				]
			},
			{
				"cachegroup": "mid4",
				"parents": [
					5
				]
			},
			{
				"cachegroup": "mid5",
				"parents": [
					5,
					6
				]
			},
			{
				"cachegroup": "mid6",
				"parents": [
					6
				]
			}
		]
	}

Response Structure
------------------
:description:           A short sentence that describes the topology.
:lastUpdated:           The date and time at which this :term:`Topology` was last updated, in ISO-like format
:name:                  The name of the topology. This can only be letters, numbers, and dashes.
:nodes:                 An array of nodes in the :term:`Topology`

	:cachegroup:            The name of a :term:`Cache Group`
	:parents:               The indices of the parents of this node in the nodes array, 0-indexed.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 13 Apr 2020 18:12:34 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: ftNcDRjYCDMkQM+o/szayKZriQZHGpcT0vNY0HpKgy88i0pXeEEeLGbUPh6LXtK7TvL76EgGECTzvCkcm+2LVA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 13 Apr 2020 17:12:34 GMT
	Content-Length: 239

	{
		"alerts": [
			{
				"text": "topology was created.",
				"level": "success"
			}
		],
		"response": {
			"description": "This is my topology",
			"name": "my-topology",
			"nodes": [
				{
					"cachegroup": "edge1",
					"parents": [
						7
					]
				},
				{
					"cachegroup": "edge2",
					"parents": [
						7,
						8
					]
				},
				{
					"cachegroup": "edge3",
					"parents": [
						8,
						9
					]
				},
				{
					"cachegroup": "edge4",
					"parents": [
						9
					]
				},
				{
					"cachegroup": "mid1",
					"parents": []
				},
				{
					"cachegroup": "mid2",
					"parents": [
						4
					]
				},
				{
					"cachegroup": "mid3",
					"parents": [
						4
					]
				},
				{
					"cachegroup": "mid4",
					"parents": [
						5
					]
				},
				{
					"cachegroup": "mid5",
					"parents": [
						5,
						6
					]
				},
				{
					"cachegroup": "mid6",
					"parents": [
						6
					]
				}
			],
			"lastUpdated": "2020-04-13 17:12:34+00"
		}
	}

``PUT``
=======
Updates a specific :term:`Topology`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+---------------------------------------------------------+
	| Name | Required | Description                                             |
	+======+==========+=========================================================+
	| name | yes      | The name of the :term:`Topology` to be updated          |
	+------+----------+---------------------------------------------------------+

:description:           A short sentence that describes the :term:`Topology`.
:name:                  The name of the :term:`Topology`. This can only be letters, numbers, and dashes.
:nodes:                 An array of nodes in the :term:`Topology`

	:cachegroup:            The name of a :term:`Cache Group` with at least 1 server in it
	:parents:               The indices of the parents of this node in the nodes array, 0-indexed.

.. code-block:: http
	:caption: Request Example

	PUT /api/3.0/topologies?name=my-topology HTTP/1.1
	User-Agent: python-requests/2.23.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 853
	Content-Type: application/json

	{
		"name": "my-topology",
		"description": "The description is updated, too",
		"nodes": [
			{
				"cachegroup": "edge1",
				"parents": [
					6
				]
			},
			{
				"cachegroup": "edge2",
				"parents": [
					6,
					7
				]
			},
			{
				"cachegroup": "edge3",
				"parents": [
					7,
					8
				]
			},
			{
				"cachegroup": "edge4",
				"parents": [
					8
				]
			},
			{
				"cachegroup": "mid2",
				"parents": []
			},
			{
				"cachegroup": "mid3",
				"parents": []
			},
			{
				"cachegroup": "mid4",
				"parents": [
					4
				]
			},
			{
				"cachegroup": "mid5",
				"parents": [
					4,
					5
				]
			},
			{
				"cachegroup": "mid6",
				"parents": [
					5
				]
			}
		]
	}

Response Structure
------------------
:description:           A short sentence that describes the :term:`Topology`.
:lastUpdated:           The date and time at which this :term:`Topology` was last updated, in ISO-like format
:name:                  The name of the :term:`Topology`. This can only be letters, numbers, and dashes.
:nodes:                 An array of nodes in the :term:`Topology`

	:cachegroup:            The name of a :term:`Cache Group`
	:parents:               The indices of the parents of this node in the nodes array, 0-indexed.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 13 Apr 2020 18:33:13 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: WVOtsoOVrEWcVjWM2TmT5DXy/a5Q0ygTZEQRhbkHHUmz9dgVLK1F5Joc9jtKA8gZu8/eM5+Tqqguh3mzrhAy/Q==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 13 Apr 2020 17:33:13 GMT
	Content-Length: 237

	{
		"alerts": [
			{
				"text": "topology was updated.",
				"level": "success"
			}
		],
		"response": {
			"description": "The description is updated, too",
			"name": "my-topology",
			"nodes": [
				{
					"cachegroup": "edge1",
					"parents": [
						6
					]
				},
				{
					"cachegroup": "edge2",
					"parents": [
						6,
						7
					]
				},
				{
					"cachegroup": "edge3",
					"parents": [
						7,
						8
					]
				},
				{
					"cachegroup": "edge4",
					"parents": [
						8
					]
				},
				{
					"cachegroup": "mid2",
					"parents": []
				},
				{
					"cachegroup": "mid3",
					"parents": []
				},
				{
					"cachegroup": "mid4",
					"parents": [
						4
					]
				},
				{
					"cachegroup": "mid5",
					"parents": [
						4,
						5
					]
				},
				{
					"cachegroup": "mid6",
					"parents": [
						5
					]
				}
			],
			"lastUpdated": "2020-04-13 17:33:13+00"
		}
	}

``DELETE``
==========
Deletes a specific :term:`Topology`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``


Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------+---------------------------------------------------------+
	| Name | Required | Description                                             |
	+======+==========+=========================================================+
	| name | yes      | The name of the :term:`Topology` to be deleted          |
	+------+----------+---------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/topologies?name=my-topology HTTP/1.1
	User-Agent: python-requests/2.23.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 13 Apr 2020 18:35:32 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: yErJobzG9IA0khvqZQK+Yi7X4pFVvOqxn6PjrdzN5DnKVm/K8Kka3REul1XmKJnMXVRY8RayoEVGDm16mBFe4Q==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 13 Apr 2020 17:35:32 GMT
	Content-Length: 87

	{
		"alerts": [
			{
				"text": "topology was deleted.",
				"level": "success"
			}
		]
	}
