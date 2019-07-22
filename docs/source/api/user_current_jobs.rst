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

.. _to-api-user-current-jobs:

*********************
``user/current/jobs``
*********************

``GET``
=======
.. deprecated:: 1.1
	Use the ``userId`` query parameter of a ``GET`` request to the :ref:`to-api-jobs` endpoint instead.

Retrieves the user's list of running and pending content invalidation jobs.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+---------+----------+-------------------------------------------------------------------------------------------------------------------------------------------+
	|  Name   | Required | Description                                                                                                                               |
	+=========+==========+===========================================================================================================================================+
	| keyword | no       | Return only jobs that have this keyword - keyword corresponds to the operation or type of job (currently only "PURGE" is a valid keyword) |
	+---------+----------+-------------------------------------------------------------------------------------------------------------------------------------------+

.. deprecated:: 1.1
	This query parameter has been deprecated because the only supported keyword is "PURGE". Jobs used to be much more versatile, but such versatility is no longer required of them. This still "works", but never has any effect on the output, except to make it an empty array if it is anything other than "PURGE".

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
:assetType: The type of asset being revalidated e.g. "file"

	.. deprecated:: 1.1
		This field still exists, but has no purpose as all assets are now treated as remote files; i.e. it will always be ``"file"``.

:createdBy:       The username of the user who initiated the job
:deliveryService: The 'xml_id' that uniquely identifies the :term:`Delivery Service` on which this job operates
:enteredTime:     The date and time at which the job was created, in ISO format
:id:              An integral, unique identifier for this job
:keyword:         A keyword that represents the operation being performed by the job:

	PURGE
		This job will prevent caching of URLs matching the ``assetURL`` until it is removed (or its Time to Live expires)

:objectName: A deprecated field of unknown use - it only still exists for legacy compatibility reasons, and will always be ``null``
:objectType: A deprecated field of unknown use - it only still exists for legacy compatibility reasons, and will always be ``null``
:parameters: A string containing key/value pairs representing parameters associated with the job - currently only uses Time to Live e.g. ``"TTL:48h"``
:startTime:  The date and time at which the job began or will begin, in ISO format
:status:     A deprecated field of unknown use - it only still exists for legacy compatibility reasons, and appears to always be ``"PENDING"``
:username:   The username of the user who created this revalidation job

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 13 Dec 2018 14:23:54 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 13 Dec 2018 18:23:54 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Ijr9pDZ4XwPIBX0Qnl+yihTYa8bK7TjJdrpDiV9VNg9k7OC9FSNQV4HSmX35KUAKMFpIHe/azutbvr0xZzQucg==
	Content-Length: 301

	{ "response": [
		{
			"keyword": "PURGE",
			"objectName": null,
			"assetUrl": "http://origin.infra.ciab.test/.*\\.jpg",
			"assetType": "file",
			"status": "PENDING",
			"username": "admin",
			"parameters": "TTL:1h",
			"enteredTime": "2018-12-13 13:56:35+00",
			"objectType": null,
			"agent": "dummy",
			"id": 1,
			"startTime": "2018-12-13 13:56:09+00"
		}
	]}

``POST``
========
Creates a new content revalidation job.

.. Note:: This method forces a HTTP *revalidation* of the content, and not a new ``GET`` - the origin needs to support revalidation according to the HTTP/1.1 specification, and send a ``200 OK`` or ``304 Not Modified`` HTTP response as appropriate.

:Auth. Required: Yes
:Roles Required: "portal"

	.. versionchanged:: ATCv3.0.2
		For security reasons, the endpoint was reworked so that regardless of tenancy, the "portal" :term:`Role` or higher is required.

:Response Type:  ``undefined``

Request Structure
-----------------
:dsId: The integral, unique identifier of the :term:`Delivery Service` on which the revalidation job shall operate

:regex: This should be a `PCRE <http://www.pcre.org/>`_-compatible regular expression for the path to match for forcing the revalidation

	.. warning:: This is concatenated directly to the origin URL of the :term:`Delivery Service` identified by ``dsId`` to make the full regular expression. Thus it is not necessary to restate the URL but it should be noted that if the origin URL does not end with a backslash (``/``) then this should begin with an escaped backslash to ensure proper behavior (otherwise it will match against FQDNs, which leads to undefined behavior in Traffic Control).

	.. note:: Be careful to only match on the content that must be removed - revalidation is an expensive operation for many origins, and a simple ``/.*`` can cause an overload in requests to the origin.

:startTime: The time and date at which the revalidation rule will be made active, in ISO format
:ttl:       Specifies the Time To Live (TTL) - in hours - for which the revalidation rule will remain active after ``startTime``

	.. note:: It usually makes sense to make this the same as the ``Cache-Control`` header from the origin which sets the object time to live in cache (by ``max-age`` or ``Expires``). Entering a longer TTL here will make the caches do unnecessary work.

:urgent: An optional boolean which, if present and ``true``, marks the job as "urgent", which has no meaning to machines but is visible to humans for their consideration

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/user/current/jobs HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 79
	Content-Type: application/json

	{
		"dsId": 1,
		"regex": "\\/.*\\.jpg",
		"startTime": "2018-12-13 13:55:09",
		"ttl": 1
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
	Date: Thu, 13 Dec 2018 13:56:35 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 13 Dec 2018 17:56:35 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Uyz2P6gkzsSu8xESEHSKQCG6+6Xw0o+wgjx30+UTBFNIZzFYlkjDK1wZdIUYUPdSbPcTRy5ZaxT1qFpl8+4aGQ==
	Content-Length: 141

	{ "alerts": [
		{
			"level": "success",
			"text": "Invalidate content request submitted for demo1 [ http://origin.infra.ciab.test.*\\.jpg - TTL:1h ]"
		}
	]}

