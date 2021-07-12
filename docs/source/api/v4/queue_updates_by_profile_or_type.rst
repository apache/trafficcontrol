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

.. _to-api-cdn-locks:

*************************
``queue_updates/{{cdn}}``
*************************

.. versionadded:: 4.0

``PUT``
=======
Allows a user to queue updates on the servers of a CDN filtered by ``type`` and/or ``profile``.
:term:`Queue` or dequeue updates for a list of servers.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-------+---------------------------------------------------------------------------------------------+
	| Name  | Description                                                                                 |
	+=======+=============================================================================================+
	|  cdn  | The name of the cdn, for which the servers are being queued or dequeued                     |
	+-------+---------------------------------------------------------------------------------------------+

.. table:: Request Query Parameters

	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                   |
	+===========+==========+===============================================================================================================+
	| type      | no       | The name of the ``type`` of servers, for which the updates need to be queued or dequeued.                     |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+
	| profile   | no       | The name of the ``profile`` of servers, for which the updates need to be queued or dequeued.                  |
	+-----------+----------+---------------------------------------------------------------------------------------------------------------+

:action: A string describing what action to take regarding server updates; one of:

	queue
		:term:`Queue Updates` for the server, propagating configuration changes to the actual server
	dequeue
		Cancels any pending updates on the server

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/queue_updates/cdn1?type=EDGE HTTP/2
	Host: localhost:8443
	User-Agent: curl/7.64.2
	Accept: */*
	Cookie: mojolicious=...
	Content-Type: application/json
	Content-Length: 21

	{
		"action": "queue"
	}

Response Structure
------------------
:action: The action processed, one of:

	queue
		:term:`Queue Updates` was performed on the list of servers, propagating configuration changes to the actual servers
	dequeue
		Canceled any pending updates on the list of servers

:cdnID: The integral, unique identifier of the cdn on which ``action`` was taken
:typeID: The integral, unique identifier of the type on which ``action`` was taken

.. code-block:: http
	:caption: Response Example

	HTTP/2 200
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=eyJhdXRoX2Rhd...; Path=/; Expires=Fri, 09 Jul 2021 20:20:24 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 5ifVe8ihG3kwrRSfu14N7IF7ldLgCtGDoNS/aQEHag3lVAce6vLLALrD4YdiDl7NLwOzifq1MC7SY8YcyHEipQ==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 09 Jul 2021 19:20:25 GMT
	Content-Length: 54

	{
		"response": {
			"action": "queue",
			"cdnID": 5,
			"typeID": 11
		}
	}
