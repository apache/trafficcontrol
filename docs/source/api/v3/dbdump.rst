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

.. _to-api-v3-dbdump:

**********
``dbdump``
**********
.. caution:: This is an extremely dangerous thing to do, as it exposes the entirety of the database, including possibly sensitive information. Administrators and systems engineers are advised to instead use database-specific tools to make server transitions more securely.

Dumps the Traffic Ops database as an SQL script that should recreate its schema and contents exactly.

.. impl-detail:: The script is output using the :manpage:`pg_dump(1)` utility, and is thus compatible for use with the :manpage:`pg_restore(1)` utility.

``GET``
=======
Fetches the database dump.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  ``undefined`` - outputs an SQL script, not JSON

Request Structure
-----------------
No parameters available

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/dbdump HTTP/1.1
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
	Content-Type: application/sql
	Content-Disposition: attachment
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: YwvPB0ZToyzT8ilBnDlWWdwV+E3f2Xgus1OKrkNaipQqgrw5zGwq0rC1U9TZ8Zl6kAGcRZgCYnr1EWfHXpJRkg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 09 Sep 2019 21:08:28 GMT
	Transfer-Encoding: chunked

	-- Actual text omitted - it's huge
