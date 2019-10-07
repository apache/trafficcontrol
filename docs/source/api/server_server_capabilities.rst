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

.. _to-api-server-server-capabilities:

**************
``server_server_capabilities``
**************

``GET``
=======
Gets all associations of :term:`Server Capabilities` to servers

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name             | Required | Description                                                                                                   |
	+==================+==========+===============================================================================================================+
	| serverId         | no       | Filter associated :term:`Server Capabilities` by server id                                                    |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| serverHostName   | no       | Filter associated :term:`Server Capabilities` by server host name                                             |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| serverCapability | no       | Filter associated :term:`Server Capabilities` by server capability                                            |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby          | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|                  |          | array                                                                                                         |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder        | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit            | no       | Choose the maximum number of results to return                                                                |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset           | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit.         |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+
	| page             | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|                  |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|                  |          | defined to make use of ``page``.                                                                              |
	+------------------+----------+---------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/server_server_capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:serverHostName:   The server's host name 
:serverId:         The server's id
:lastUpdated:      The date and time at which this association between the server and the :term:`Server Capability` was last updated, in an ISO-like format
:serverCapability: The :term:`Server Capability`'s name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: UFO3/jcBFmFZM7CsrsIwTfPc5v8gUiXqJm6BNp1boPb4EQBnWNXZh/DbBwhMAOJoeqDImoDlrLnrVjQGO4AooA==
	X-Server-Name: traffic_ops_golang/
	Mon, 07 Oct 2019 22:15:11 GMT
	Content-Length: 150

	{
		"response": [
			{
				"lastUpdated": "2019-10-07 22:05:31+00",
				"serverHostName": "atlanta-org-1",
				"serverId": 260,
				"serverCapability": "ram"
			},
			{
				"lastUpdated": "2019-10-07 22:05:31+00",
				"serverHostName": "atlanta-org-2",
				"serverId": 261,
				"serverCapability": "disk"
			}
		]
	}

``POST``
========
Associates a :term:`Server Capability` to a server.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:serverId:         The server's id to associate
:serverCapability: The :term:`Server Capability`'s name to associate

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/server_server_capabilities HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
		"serverId": 1,
		"serverCapability": "disk"
	}

Response Structure
------------------
:serverId:         The server's id
:lastUpdated:      The date and time at which this association between the server and the :term:`Server Capability` was last updated, in an ISO-like format
:serverCapability: The :term:`Server Capability`'s name

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Oct 2019 22:15:11 GMT
	Content-Length: 157

	{
		"alerts": [
			{
				"text": "server server_capability was created.",
				"level": "success"
			}
		],
		"response": {
			"lastUpdated": "2019-10-07 22:15:11+00",
			"serverId": 1,
			"serverCapability": "disk"
		}
	}

``DELETE``
=======
Disassociate a server from a :term:`Server Capability`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------------+----------+------------------------------------------------------------------+
	| Name             | Required | Description                                                      |
	+==================+==========+==================================================================+
	| serverId         | yes      | ID of the server to disassociate                                 |
	+------------------+----------+------------------------------------------------------------------+
	| serverCapability | yes      | :term:`Server Capability` name to disassociate from given server |
	+------------------+----------+------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/server_server_capabilities?serverId=1&serverCapability=disk HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: UFO3/jcBFmFZM7CsrsIwTfPc5v8gUiXqJm6BNp1boPb4EQBnWNXZh/DbBwhMAOJoeqDImoDlrLnrVjQGO4AooA==
	X-Server-Name: traffic_ops_golang/
	Mon, 07 Oct 2019 22:15:11 GMT
	Content-Length: 96

	{
		"alerts": [
			{
				"text": "server server_capability was deleted.",
				"level": "success"
			}
		]
	}