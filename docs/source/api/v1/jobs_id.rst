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

.. _to-api-v1-jobs-id:

***************
``jobs/{{ID}}``
***************
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-v1-jobs` with the ``id`` query parameter instead.

``GET``
=======
Get details about a specific content invalidation job.

:Auth. Required: Yes
:Roles Required: "operations" or "admin"\ [#tenancy]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name | Description                                                |
	+======+============================================================+
	|  ID  | An integral, unique identifier for the job to be inspected |
	+------+------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/jobs/3 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:assetUrl:        A regular expression - matching URLs will be operated upon according to ``keyword``
:createdBy:       The username of the user who initiated the job
:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetUrl`` until it is removed (or its Time to Live expires)

:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began, in a non-standard format

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: l7qvgOShdIFukHyOhi8es2BG6zJZ6RXTT7OKABtI8b1y+cE4nxFq11T5OG5yXjKo69eTYOD7xUUdLqneT2E/VA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Jun 2019 13:29:21 GMT
	Content-Length: 192

	{ "response": [{
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 3,
		"keyword": "PURGE",
		"parameters": "TTL:3h",
		"startTime": "2019-06-21 00:00:00+00"
	}],
	"alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /jobs with the 'id' parameter instead",
			"level": "warning"
		}
	]}


.. [#tenancy] When viewing content invalidation jobs, only those jobs that operate on a :term:`Delivery Service` visible to the requesting user's :term:`Tenant` will be returned.
