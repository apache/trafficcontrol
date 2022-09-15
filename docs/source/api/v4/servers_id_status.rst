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

.. _to-api-v4-servers-id-status:

*************************
``servers/{{ID}}/status``
*************************

``PUT``
=======
Updates server status and queues updates on all descendant :term:`Topology` nodes or child caches if server type is EDGE or MID. Also, captures offline reason if status is set to ADMIN_DOWN or OFFLINE and prepends offline reason with the user that initiated the status change.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: SERVER:UPDATE, SERVER:READ, STATUS:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------+
	| Name | Description                                                                 |
	+======+=============================================================================+
	|  ID  | The integral, unique identifier of the server whose status is being changed |
	+------+-----------------------------------------------------------------------------+

:offlineReason: A string containing the reason for the status change
:status:        The name or integral, unique identifier of the server's new status

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/servers/13/status HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 56
	Content-Type: application/json

	{
		"status": "ADMIN_DOWN",
		"offlineReason": "Bad drives"
	}

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
	Date: Mon, 10 Dec 2018 18:08:44 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: LS1jCo5eMVKxmeYDol0I2LgLYazocSggR5hynNoLcPmMov9u2s3ulksPdQtG1N3aS+VM9tdMsCrahFPraLJVwg==
	Content-Length: 158

	{ "alerts": [
		{
			"level": "success",
			"text": "Updated status [ ADMIN_DOWN ] for quest.infra.ciab.test [ admin: Bad drives ] and queued updates on all child caches"
		}
	]}
