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

.. _to-api-v4-topologies-name-queue_update:

************************************
``topologies/{{name}}/queue_update``
************************************

``POST``
========
:term:`Queue` or "dequeue" updates for all servers assigned to the :term:`Cache Groups` in a specific :term:`Topology`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:QUEUE, TOPOLOGY:READ, SERVER:READ, CACHE-GROUP:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------+
	| Name | Description                                                               |
	+======+===========================================================================+
	| name | The name of the :term:`Topology` on which to queue or dequeue updates.    |
	+------+---------------------------------------------------------------------------+

:action: One of "queue" or "dequeue" as appropriate
:cdnId:  The integral, unique identifier for the CDN on which to (de)queue updates

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/topologies/demo1-top/queue_update HTTP/1.1
	User-Agent: python-requests/2.24.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 28

	{
		"action": "queue",
		"cdnId": 1
	}

Response Structure
------------------
:action:   The action processed, either ``"queue"`` or ``"dequeue"``
:cdnId:    The CDN ID on which :term:`Queue Updates` was performed or cleared
:topology: The name of the :term:`Topology` on which :term:`Queue Updates` was performed or cleared

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 08 Sep 2020 17:35:42 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: nmu3TMVmllcHeEstLuiqPsEpypNV2jcs5PyrqsqJKkexxxb8N7qk84AWzTJWUpsfdVWrj/JzRiCPGJS4hw0phQ==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 08 Sep 2020 16:35:42 GMT
	Content-Length: 79

	{
		"response": {
			"action": "queue",
			"cdnId": 1,
			"topology": "demo1-top"
		}
	}
