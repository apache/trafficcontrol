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

.. _to-api-v3-servers-hostname-update:

*************************************
``servers/{{HostName-Or-ID}}/update``
*************************************

``POST``
========
:term:`Queue` or dequeue updates and revalidation updates for a specific server.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  undefined

Request Structure
-----------------
.. table:: Request Path Parameters

	+------------------+---------------------------------------------------------------------------------------------------------+
	| Name             | Description                                                                                             |
	+==================+=========================================================================================================+
	|  HostName-OR-ID  | The hostName or integral, unique identifier of the server on which updates are being queued or dequeued |
	+------------------+---------------------------------------------------------------------------------------------------------+

.. table:: Request Query Parameters

	+----------------------------+----------+--------------------------------------------------------------------------------------+
	| Name                       | Required | Description                                                                          |
	+============================+==========+======================================================================================+
	| updated (Deprecated)       | no       | The value to set for the queue update flag on this server. May be 'true' or 'false'. |
	+----------------------------+----------+--------------------------------------------------------------------------------------+
	| reval_updated (Deprecated) | no       | The value to set for the queue update flag on this server. May be 'true' or 'false'. |
	+----------------------------+----------+--------------------------------------------------------------------------------------+

.. deprecated:: 3.1

	The boolean values above have resulted in an unintended race condition. These fields are removed in APIv4 and replaced.

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/servers/my-edge/update?updated=true&reval_updated=false HTTP/1.1
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 10 Dec 2018 18:20:04 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 9Mmo9hIFZyF5gAvfdJD//VH9eNgiHVLinXt88H0GlJSHhwND8gMxaFyC+f9XZfiNAoGd1MKi1934ZJGmaIR6qQ==
	Content-Length: 49

	{
		"alerts" :
			[
				{
					"text" : "successfully set server 'my-edge' updated=true reval_updated=false",
					"level" : "success"
				}
			]
	}
