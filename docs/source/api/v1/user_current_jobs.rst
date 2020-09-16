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

.. _to-api-v1-user-current-jobs:

*********************
``user/current/jobs``
*********************
.. deprecated:: ATCv4

	Both request methods supported by this endpoint are implemented (better) by :ref:`to-api-v1-jobs`, and in the future that will be the only way to interact with jobs. Developers and administrators are encouraged to switch at their earliest convenience.

``GET``
=======

Retrieves the user's list of running and pending content invalidation jobs.

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Response Type:  Array

Request Structure
-----------------
.. versionchanged:: ATCv4
	Prior to version 4 of Traffic Control, the deprecated ``keyword`` query parameter was available to filter jobs in the response. As only one keyword is meaningful, this was never used and so has been removed.

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/user/current/jobs?keyword=PURGE HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:agent: The name of the agent or process responsible for running the job

	.. deprecated:: 1.1
		This field is no longer used, but does still exist for legacy compatibility reasons. It will always be ``"dummy"``.

:assetUrl:  A regular expression - matching URLs will be operated upon according to ``keyword``
:assetType: The type of asset being invalidated

	.. deprecated:: 1.1
		This field still exists, but has no purpose as all assets are now treated as remote files; i.e. it will always be ``"file"``.

:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` on which this job operates
:enteredTime:     The date and time at which the job was created, in the same format as the ``last_updated`` fields seen throughout other API responses

	.. versionchanged:: ATCv4
		This used to be in the legacy ``YYYY-MM-DD HH:MM:SS`` format, but as of Traffic Control version 4 they are standardized to match the format of other date strings in API responses

:id:      An integral, unique identifier for this job
:keyword: A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetURL`` until it is removed (or its Time to Live expires)

:objectName: A deprecated field of unknown use - it only still exists for legacy compatibility reasons, and will always be ``null``
:objectType: A deprecated field of unknown use - it only still exists for legacy compatibility reasons, and will always be ``null``
:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began or will begin, in the same format as the ``last_updated`` fields seen throughout other API responses

	.. versionchanged:: ATCv4
		This used to be in the legacy ``YYYY-MM-DD HH:MM:SS`` format, but as of Traffic Control version 4 they are standardized to match the format of other date strings in API responses

:status:   A deprecated field of unknown use - it only still exists for legacy compatibility reasons, and appears to always be ``"PENDING"``
:username: The username of the user who created this revalidation job

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: RxFZN2+OvP3HEyp+KlCPDFT74PwPFNjxBjibGIMPhbRjVEb8PhdaF7Gq61wklNRfda4PgTP2tzOheiM0oUzUTQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Jun 2019 13:23:18 GMT
	Content-Length: 747

	{ "alerts": [
		{
			"text": "This endpoint is deprecated, please use the 'userId' or 'createdBy' query parameters of a GET request to /jobs instead",
			"level": "warning"
		}
	],
	"response": [{
		"agent": 1,
		"assetType": "file",
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"deliveryService": "demo1",
		"enteredTime": "2019-06-19 13:19:51+00",
		"id": 3,
		"keyword": "PURGE",
		"objectName": null,
		"objectType": null,
		"parameters": "TTL:3h",
		"username": "admin"
	}]}

``POST``
========

Creates a new content revalidation job.

.. caution:: Creating a content invalidation job immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that content invalidation jobs become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required: Yes
:Roles Required: "portal"\ [#tenancy]_

	.. versionchanged:: ATCv3.1.0
		For security reasons, the endpoint was reworked so that regardless of tenancy, the "portal" :term:`Role` or higher is required.

:Response Type:  ``undefined``

Request Structure
-----------------
:dsId:  The integral, unique identifier of the :term:`Delivery Service` on which the revalidation job shall operate
:regex: This should be a `PCRE <http://www.pcre.org/>`_-compatible regular expression for the path to match for forcing the revalidation

	.. warning:: This is concatenated directly to the origin URL of the :term:`Delivery Service` identified by ``dsId`` to make the full regular expression. Thus it is not necessary to restate the URL but it should be noted that if the origin URL does not end with a backslash (``/``) then this should begin with an escaped backslash to ensure proper behavior (otherwise it will match against FQDNs, which leads to undefined behavior in Traffic Control).

	.. note:: Be careful to only match on the content that must be removed - revalidation is an expensive operation for many origins, and a simple ``/.*`` can cause an overload in requests to the origin.

:startTime: This can be a string in the legacy ``YYYY-MM-DD HH:MM:SS`` format, or a string in :rfc:`3339` format, or a string representing a date in the same non-standard format as the ``last_updated`` fields common in other API responses, or finally it can be a number indicating the number of milliseconds since the Unix Epoch (January 1, 1970 UTC). This date must be in the future, and unlike a ``POST`` request to :ref:`to-api-v1-jobs`, it must be *within two days from the time of creation*.

	.. versionchanged:: ATCv4
		Prior to Traffic Control version 4, this used to **only** accept the legacy ``YYYY-MM-DD HH:MM:SS`` date string format, but this constraint has been relaxed. Developers are encouraged to submit date/time strings in either :rfc:`3339` format or as a numerical Unix timestamp (in milliseconds).

:ttl: Specifies the :abbr:`TTL (Time To Live)` - in hours - for which the revalidation rule will remain active after ``startTime``
:urgent: An optional boolean which, if present and ``true``, marks the job as "urgent", which has no meaning whatsoever, and in fact is not even stored by Traffic Control. So don't use it.

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/user/current/jobs HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 67
	Content-Type: application/json

	{
		"dsId": 1,
		"startTime": "2019-06-21T00:00:00Z",
		"regex": "/.*",
		"ttl": 3
	}


Response Structure
------------------
.. versionchanged:: ATCv4
	This method of this endpoint used to only return a successful ``alert`` (presuming success), but in ATCv4 a representation of the newly-created content invalidation job was added to the response.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Location: https://trafficops.infra.ciab.test/api/1.4/jobs?id=3
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: zQrzB3SLXTbpxLaVWq4WHeONUfEirXDaLRlCi/4+fekgtbjnDgGnA+Sq6MGaxRyQ92/96IsYjAP3Re6ZoN7rzg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 19 Jun 2019 13:19:51 GMT
	Content-Length: 235

	{ "alerts": [
		{
			"text": "This endpoint is deprecated, please use the POST method /jobs instead",
			"level": "warning"
		},
		{
			"text": "Invalidation Job creation was successful",
			"level": "success"
		}
	],
	"response": {
		"assetUrl": "http://origin.infra.ciab.test/.*",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 1,
		"keyword": "PURGE",
		"parameters": "TTL:3h",
		"startTime": "2019-06-21 00:00:00+00"
	}}

.. [#tenancy] When viewing content invalidation jobs, only those jobs that operate on a :term:`Delivery Service` visible to the requesting user's :term:`Tenant` will be returned. Likewise, creating a new content invalidation job requires that the target :term:`Delivery Service` is modifiable by the requesting user's :term:`Tenant`.
