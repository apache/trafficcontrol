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
.. _to-api-logs:

********
``logs``
********

.. note:: This endpoint's responses will contain a cookie (``last_seen_log``) that is used by :ref:`to-api-logs-newcount` to determine the time of last access. Be sure your client uses cookies properly if you intend to use :ref:`to-api-logs-newcount` in concert with this endpoint!

``GET``
=======
Fetches a list of changes that have been made to the Traffic Control system.

:Auth. Required:       Yes
:Roles Required:       None
:Permissions Required: LOG:READ
:Response Type:        Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------+----------+-------------------------------------------------------------------------------------------------------------------------------------+
	| Name   | Required | Description                                                                                                                         |
	+========+==========+=====================================================================================================================================+
	| days   | no       | An integer number of days of change logs to return - 0 means "no limit" (which could be quite a lot!)                               |
	+--------+----------+-------------------------------------------------------------------------------------------------------------------------------------+
	| limit  | no       | The number of records to which to limit the response, by default there is no limit applied                                          |
	+--------+----------+-------------------------------------------------------------------------------------------------------------------------------------+
	| offset | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                                |
	+--------+----------+-------------------------------------------------------------------------------------------------------------------------------------+
	| page   | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first page is 1.|
	|        |          | If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``.                   |
	+--------+----------+-------------------------------------------------------------------------------------------------------------------------------------+
	| user   | no       | A name to which to limit the response too                                                                                           |
	+--------+----------+-------------------------------------------------------------------------------------------------------------------------------------+

.. versionadded:: ATCv6
	The ``username``, ``page``, ``offset`` query parameters were added to this in endpoint across across all API versions in :abbr:`ATC (Apache Traffic Control)` version 6.0.0.

.. versionchanged:: 4.0
	The ``username`` query string parameter was renamed to ``user`` so that it has the same name as the response property by which it filters.

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/logs?days=1&limit=2&username=admin HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:lastUpdated: Date and time at which the change was made, in :rfc:`3339` format (the name was chosen for consistency for the rest of the API; changelog entries are never "updated")
:message:     Log detail about what occurred
:user:        username of the user who made the change

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; HttpOnly, last_seen_log=2021-11-22T02:34:06.583699419Z;
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 22 Nov 2021 02:34:06 GMT
	Content-Length: 220

	{ "response": [
		{
			"lastUpdated": "2021-11-22T01:59:32.692767Z",
			"message": "CDN: CDN-in-a-Box, ID: 2, ACTION: server updates queued on 6 servers",
			"user": "admin"
		},
		{
			"lastUpdated": "2021-11-22T01:59:30.624049Z",
			"message": "CDN: CDN-in-a-Box, ID: -1, ACTION: Snapshot of CRConfig and Monitor",
			"user": "admin"
		}
	],
	"summary": {
		"count": 467
	}}

Summary Fields
""""""""""""""
The ``summary`` object returned by this method of this endpoint uses only the ``count`` :ref:`standard property <reserved-summary-fields>`.
