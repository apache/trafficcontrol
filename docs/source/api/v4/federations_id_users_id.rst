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

.. _to-api-v4-federations-id-users-id:

***************************************
``federations/{{ID}}/users/{{userID}}``
***************************************

``DELETE``
==========
Removes a user from a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: FEDERATION:UPDATE, FEDERATION:READ, USER:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+--------+----------------------------------------------------------------------------------------------------------------+
	|  Name  | Description                                                                                                    |
	+========+================================================================================================================+
	|   ID   | An integral, unique identifier for the federation from which the user identified by ``userID`` will be removed |
	+--------+----------------------------------------------------------------------------------------------------------------+
	| userID | An integral, unique identifier for the user who will be removed from the federation identified by ``ID``       |
	+--------+----------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Structure

	DELETE /api/4.0/federations/1/users/2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 01:14:04 GMT
	X-Server-Name: traffic_ops_golang/
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 05:14:04 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: xdF6l7jdd2t8au6lh4pFtDqYxTfehzke2aDBuytL7I74hK9KCT7ssLuYbfvD8ejdqqF3+jiBiFk7neQ8c4vVUQ==
	content-length: 93

	{ "alerts": [
		{
			"level": "success",
			"text": "Removed user [ admin ] from federation [ foo.bar. ]"
		}
	]}
