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

.. _to-api-v1-hwinfo:

**********
``hwinfo``
**********
.. deprecated:: 1.1
	This endpoint still works, but it is unused and serves no purpose. It will always return an empty ``response`` array unless the database is manually altered.

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name           | Required | Description                                                                                                   |
	+================+==========+===============================================================================================================+
	| id             | no       | An integral, unique identifier of a specific hwinfo object which will be retrieved                            |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| serverHostName | no       | The name of the server for which hwinfo objects will be retrieved                                             |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| serverId       | no       | The integral, unique identifier of a server for which hwinfo objects will be retrieved                        |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| description    | no       | The description of a hwinfo object; only hwinfo objects with descriptions matching this will be retrieved     |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| val            | no       | The value of a hwinfo object; only hwinfo objects with values matching this will be retrieved                 |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| lastUpdated    | no       | Only hwinfo objects that were last updated at this ISO-format date and time will be retrieved                 |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby        | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|                |          | array                                                                                                         |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder      | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit          | no       | Choose the maximum number of results to return. Default if not specified is 1000.                             |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset         | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+
	| page           | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|                |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect.                       |
	+----------------+----------+---------------------------------------------------------------------------------------------------------------+

.. caution:: The ``lastUpdated`` query parameter doesn't seem to work properly, and its use is therefore discouraged.

.. code-block:: http
	:caption: Request Example

	GET /api/1.3/hwinfo HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:description:    Freeform description for this specific server's hardware info
:lastUpdated:    The Time and Date for the last update for this server
:serverHostName: Hostname for this specific server's hardware info
:serverId:       Local unique identifier for this specific server's hardware info
:val:            Freeform value used to track anything about a server's hardware info

Also, in addition to the regular ``response`` field and any and all ``alerts``, this endpoint returns an extra top-level JSON key: ``limit``.

:limit: The number of results to which the result was limited. Should be exactly as specified in the `Request Structure`_.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 22 Nov 2019 20:28:07 GMT; Max-Age=3600; HttpOnly
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 22 Nov 2019 19:28:07 GMT
	Content-Length: 138

	{ "alerts": [
		{
			"text": "This endpoint is deprecated, and will be removed in the future",
			"level": "warning"
		}
	],
	"response": [
		{
			"description": "quest",
			"lastUpdated": "2019-11-22 19:31:26+00",
			"serverHostName": "dns",
			"serverId": 1,
			"val": "test"
		}
	],
	"limit": 1000
	}
