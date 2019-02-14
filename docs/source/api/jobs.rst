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

.. _to-api-jobs:

********
``jobs``
********

``GET``
=======
Get all jobs (currently limited to invalidate content (PURGE) jobs) sorted by start time (descending).

:Auth. Required: Yes
:Roles Required: "operations" or "admin"
:Response Type:  Array

.. warning:: This endpoint will respect tenancy rules *if and only if*  the ``dsId`` query parameter is used.

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------+----------+----------------------------------------------------------------------------------------------------------------------+
	|  Name  | Required | Description                                                                                                          |
	+========+==========+======================================================================================================================+
	|  dsId  | no       | Return only invalidation jobs pending on the :term:`Delivery Service` identified by this integral, unique identifier |
	+--------+----------+----------------------------------------------------------------------------------------------------------------------+
	| userId | no       | Return only invalidation jobs created by the user identified by this integral, unique identifier                     |
	+--------+----------+----------------------------------------------------------------------------------------------------------------------+

.. note:: If the ``dsId`` parameter is given, an error will be returned if the thereby identified :term:`Delivery Service` is not visible to the logged-in user's Tenant

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/jobs?dsId=1&userId=2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The 'xml_id' that uniquely identifies the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetURL`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in ISO format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 05 Dec 2018 15:44:07 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Wed, 05 Dec 2018 19:44:07 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: PfKJyahkUTjK2iNAuY3NZiljuHHVNJyNkdKRdtzHZN9fg4+HidejGIC19tcyRCDATQyZQ49/BLEIJDAAaqTwzA==
	Content-Length: 202

	{ "response": [
		{
			"parameters": "TTL:48h",
			"keyword": "PURGE",
			"assetUrl": "http://origin.infra.ciab.test/.*\\.jpg",
			"createdBy": "admin",
			"startTime": "2018-12-05 15:43:42+00",
			"id": 1,
			"deliveryService": "demo1"
		}
	]}

.. TODO: figure out why POST at this endpoint is giving 'unauthenticated' instead of 'resource not found'
