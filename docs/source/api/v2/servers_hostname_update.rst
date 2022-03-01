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

.. _to-api-v2-servers-hostname-update:

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

	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name                       | Required | Description                                                                                                  |
	+============================+==========+==============================================================================================================+
	| updated (Deprecated)       | no       | The value to set for the queue update flag on this server. May be 'true' or 'false'.                         |
	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| reval_updated (Deprecated) | no       | The value to set for the reval update flag on this server. May be 'true' or 'false'.                         |
	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| config_update_time         | no       | The value to set for when a queue update is requested for this server. Must be a valid RFC333Nano timestamp. |
	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| config_apply_time          | no       | The value to set for when a queue update is applied for this server. Must be a valid RFC333Nano timestamp.   |
	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| revalidate_update_time     | no       | The value to set for when a reval update is requested for this server. Must be a valid RFC333Nano timestamp. |
	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+
	| revalidate_apply_time      | no       | The value to set for when a reval update is applied for this server. Must be a valid RFC333Nano timestamp.   |
	+----------------------------+----------+--------------------------------------------------------------------------------------------------------------+

The ``updated`` and ``reval_updated`` boolean values should be considered deprecated and may be removed in a future release. Prefer to use the timestamp values associated with ``config_update_time``, ``config_apply_time``, ``revalidate_update_time``, and ``revalidate_apply_time``.

.. note:: You will not be able to send both a boolean value (``updated`` and ``reval_updated``) and the corresponding timestamp values (``config_update_time``, ``config_apply_time``, ``revalidate_update_time``, and ``revalidate_apply_time``) in the same request as they will conflict.

.. code-block:: http
	:caption: Request Example

	GET /api/2.0/servers/my-edge/update?config_update_time=2022-01-31T12%3A00%3A00.123456-07%3A00&revalidate_update_time=2022-01-31T12%3A00%3A00.123456-07%3A00 HTTP/1.1
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
					"text" : "successfully set server 'my-edge' config_update_time=2022-01-31T12:00:00.123456-07:00 revalidate_update_time=2022-01-31T12:00:00.123456-07:00",
					"level" : "success"
				}
			]
	}
