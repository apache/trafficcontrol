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

.. _to-api-v4-jobs:

********
``jobs``
********

``GET``
=======
Retrieve :term:`Content Invalidation Jobs`.

:Auth. Required:       Yes
:Roles Required:       None\ [#tenancy]_
:Permissions Required: JOB:READ, DELIVERY-SERVICE:READ\ [#tenancy]_
:Response Type:        Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| Name                 | Required | Description                                                                                                                          |
	+======================+==========+======================================================================================================================================+
	| assetUrl             | no       | Return only :term:`Content Invalidation Jobs` with this :ref:`job-asset-url`                                                         |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| cdn                  | no       | Return only :term:`Content Invalidation Jobs` for :term:`Delivery Services` within the CDN with this name                            |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| createdBy            | no       | Return only :term:`Content Invalidation Jobs` that were created by the user with this username                                       |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| deliveryService      | no       | Return only :term:`Content Invalidation Jobs` that operate on the :term:`Delivery Service` with this :ref:`ds-xmlid`                 |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| dsId                 | no       | Return only :term:`Content Invalidation Jobs` pending on the :term:`Delivery Service` identified by this integral, unique identifier |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| id                   | no       | Return only the single :term:`Content Invalidation Job` with this :ref:`job-id`                                                      |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| maxRevalDurationDays | no       | Return only :term:`Content Invalidation Jobs` with a :ref:`job-start-time` that is within the window defined by the                  |
	|                      |          | ``maxRevalDurationDays`` :term:`Parameter` in :ref:`the-global-profile`                                                              |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| userId               | no       | Return only :term:`Content Invalidation Jobs` created by the user identified by this integral, unique identifier                     |
	+----------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+


.. code-block:: http
	:caption: Request Example

	GET /api/4.0/jobs?id=1&dsId=1&userId=2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.20.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:id:               The :ref:`job-id`
:assetUrl:         The :ref:`job-asset-url`
:createdBy:        The :ref:`job-created-by`
:deliveryService:  The :ref:`job-ds`
:ttlHours:         The :ref:`job-ttl`
:invalidationType: The :ref:`job-invalidation-type`
:startTime:        The :ref:`job-start-time`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 12 Nov 2021 19:30:36 GMT
	Content-Length: 206

	{ "response": [{
		"id": 1,
		"assetUrl": "http://origin.infra.ciab.test/.+",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"ttlHours": 72,
		"invalidationType": "REFETCH",
		"startTime": "2021-11-09T01:02:03Z"
	}]}



``POST``
========
Creates a new :term:`Content Invalidation Jobs`.

.. caution:: Creating a :term:`Content Invalidation Job` immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that :term:`Content Invalidation Jobs` become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required:       Yes
:Roles Required:       "operations" or "admin"\ [#tenancy]_
:Permissions Required: JOB:CREATE, JOB:READ, DELIVERY-SERVICE:READ, DELIVERY-SERVICE:UPDATE\ [#tenancy]_
:Response Type:        Object

Request Structure
-----------------
:deliveryService:  The :ref:`job-ds`
:invalidationType: The :ref:`job-invalidation-type`
:regex:            The :ref:`job-regex`
:startTime:        The :ref:`job-start-time`
:ttlHours:         The :ref:`job-ttl`

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/jobs HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Transfer-Encoding: chunked
	Content-Type: application/json

	{
		"deliveryService": "demo1",
		"invalidationType": "REFRESH",
		"regex": "/.+",
		"startTime": "2021-11-09T01:02:03Z",
		"ttlHours": 72
	}


Response Structure
------------------
:assetUrl:         The :ref:`job-asset-url`
:createdBy:        The :ref:`job-created-by`
:deliveryService:  The :ref:`job-ds`
:id:               The :ref:`job-id`.
:invalidationType: The :ref:`job-invalidation-type`
:ttlHours:         The :ref:`job-ttl`
:startTime:        The :ref:`job-start-time`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Location: https://localhost:6443/api/4.0/jobs?id=1
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 08 Nov 2021 15:44:46 GMT
	Content-Length: 265

	{
		"alerts": [
			{
				"text": "Invalidation (REFRESH) request created for http://origin.infra.ciab.test/.+, start:2021-11-09 01:02:03 +0000 UTC end 2021-11-12 01:02:03 +0000 UTC",
				"level": "success"
			}
		],
		"response": {
			"id": 1,
			"assetUrl": "http://origin.infra.ciab.test/.+",
			"createdBy": "admin",
			"deliveryService": "demo1",
			"ttlHours": 72,
			"invalidationType": "REFRESH",
			"startTime": "2021-11-09T01:02:03Z"
		}
	}



``PUT``
=======
Replaces an existing :term:`Content Invalidation Job` with a new one provided in the request. This method of editing a :term:`Content Invalidation Job` does not prevent the requesting user from changing fields that normally only have one value. Use with care.

.. caution:: Modifying a :term:`Content Invalidation Job` immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that :term:`Content Invalidation Jobs` become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required:       Yes
:Roles Required:       "operations" or "admin"\ [#tenancy]_
:Permissions Required: JOB:UPDATE, DELIVERY-SERVICE:UPDATE, JOB:READ, DELIVERY-SERVICE:READ\ [#tenancy]_
:Response Type:        Object

Request Structure
-----------------
.. table:: Query Parameters

	+------+----------+----------------------------------------------------------------------------------------+
	| Name | Required | Description                                                                            |
	+======+==========+========================================================================================+
	| id   | yes      | The integral, unique identifier of the :term:`Content Invalidation Job` being modified |
	+------+----------+----------------------------------------------------------------------------------------+

:assetUrl:         The :ref:`job-asset-url` - the scheme and authority parts of the regular expression cannot be changed
:createdBy:        The :ref:`job-created-by`\ [#immutable]_
:deliveryService:  The :ref:`job-ds`\ [#immutable]_
:id:               The :ref:`job-id`\ [#immutable]_
:invalidationType: The :ref:`job-invalidation-type`
:ttlHours:         The :ref:`job-ttl`
:startTime:        The :ref:`job-start-time`

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/jobs?id=1 HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 191

	{
		"assetUrl": "http://origin.infra.ciab.test/.+",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 1,
		"invalidationType": "REFETCH",
		"startTime": "2021-11-09T01:02:03Z",
		"ttlHours": 72
	}


Response Structure
------------------
:assetUrl:         The :ref:`job-asset-url`
:createdBy:        The :ref:`job-created-by`
:deliveryService:  The :ref:`job-ds`
:id:               The :ref:`job-id`
:invalidationType: The :ref:`job-invalidation-type`
:ttlHours:         The :ref:`job-ttl`
:startTime:        The :ref:`job-start-time`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 08 Nov 2021 16:43:35 GMT
	Content-Length: 266

	{ "alerts": [
		{
			"text": "Invalidation request created for http://origin.infra.ciab.test/.+, start: 2021-11-09 01:02:03 +0000 UTC end: 2021-11-12 01:02:03 +0000 UTC invalidation type: REFETCH",
			"level": "success"
		}
	],
	"response": {
		"assetUrl": "http://origin.infra.ciab.test/.+",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 1,
		"invalidationType": "REFETCH",
		"startTime": "2021-11-09T01:02:03Z",
		"ttlHours": 72
	}}


``DELETE``
==========
Deletes a :term:`Content Invalidation Job`.

.. tip:: :term:`Content Invalidation Jobs` that have passed their :abbr:`TTL (Time To Live)` are not automatically deleted - for record-keeping purposes - so use this to clean up old jobs that are no longer useful.

.. caution:: Deleting a :term:`Content Invalidation Job` immediately triggers a CDN-wide revalidation update. In the case that the global :term:`Parameter` ``use_reval_pending`` has a value of exactly ``"0"``, this will instead trigger a CDN-wide "Queue Updates". This means that :term:`Content Invalidation Jobs` become active **immediately** at their ``startTime`` - unlike most other configuration changes they do not wait for a :term:`Snapshot` or a "Queue Updates". Furthermore, if the global :term:`Parameter` ``use_reval_pending`` *is* ``"0"``, this will cause all pending configuration changes to propagate to all :term:`cache servers` in the CDN. Take care when using this endpoint.

:Auth. Required:       Yes
:Roles Required:       "operations" or "admin"\ [#tenancy]_
:Permissions Required: JOB:DELETE, JOB:READ, DELIVERY-SERVICE:UPDATE, DELIVERY-SERVICE:READ\ [#tenancy]_
:Response Type:        Object

Request Structure
-----------------
.. table:: Query Parameters

	+------+----------+----------------------------------------------------------------------------------------+
	| Name | Required | Description                                                                            |
	+======+==========+========================================================================================+
	| id   | yes      | The integral, unique identifier of the :term:`Content Invalidation Job` being modified |
	+------+----------+----------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/jobs?id=1 HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0


Response Structure
------------------
:assetUrl:         The :ref:`job-asset-url` of the deleted :term:`Content Invalidation Job`
:createdBy:        The :ref:`job-created-by` of the deleted :term:`Content Invalidation Job`
:deliveryService:  The :ref:`job-ds` of the deleted :term:`Content Invalidation Job`
:id:               The :ref:`job-id`. of the deleted :term:`Content Invalidation Job`
:invalidationType: The :ref:`job-invalidation-type` of the deleted :term:`Content Invalidation Job`
:ttlHours:         The :ref:`job-ttl` of the deleted :term:`Content Invalidation Job`
:startTime:        The :ref:`job-start-time` of the deleted :term:`Content Invalidation Job`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 08 Nov 2021 16:54:32 GMT
	Content-Length: 230

	{ "alerts": [
		{
			"text": "Content invalidation job was deleted",
			"level": "success"
		}
	],
	"response": {
		"assetUrl": "http://origin.infra.ciab.test/.+",
		"createdBy": "admin",
		"deliveryService": "demo1",
		"id": 1,
		"invalidationType": "REFETCH",
		"startTime": "2021-11-09T01:02:03Z",
		"ttlHours": 72
	}}


.. [#tenancy] When viewing :term:`Content Invalidation Jobs`, only those jobs that operate on a :term:`Delivery Service` visible to the requesting user's :term:`Tenant` will be returned. Likewise, creating a new :term:`Content Invalidation Jobs` requires that the target :term:`Delivery Service` is modifiable by the requesting user's :term:`Tenant`. However, when modifying or deleting an existing :term:`Content Invalidation Jobs`, the operation can be completed if and only if the requesting user's :term:`Tenant` is the same as the job's :term:`Delivery Service`'s :term:`Tenant` or a descendant thereof, **and** if the requesting user's :term:`Tenant` is the same as the :term:`Tenant` of the *user who initially created the job* or a descendant thereof.
.. [#immutable] This field must exist, but it must *not* be different than the same field of the existing job (i.e. as seen in a GET_ response). That is, this cannot be changed.
