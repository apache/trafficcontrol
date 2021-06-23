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

.. _to-api-v3-api_capabilities:

********************
``api_capabilities``
********************
Deals with the capabilities that may be associated with API endpoints and methods. These capabilities are assigned to :term:`Roles`, of which a user may have one or more. Capabilities support "wildcarding" or "globbing" using asterisks to group multiple routes into a single capability

``GET``
=======
Get all API-capability mappings.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| Name           | Required | Description                                                                                            |
	+================+==========+========================================================================================================+
	| capability     | no       | Return only the Capability that has this name                                                          |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| id             | no       | Return only the Capability that has this integral, unique identifier                                   |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| httpMethod     | no       | Return only Capabilities which have this ``httpMethod``                                                |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| route          | no       | Return only Capabilities which have this ``route``                                                     |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| lastUpdated    | no       | Return only Capabilites which were last updated at this **exact** date and time\ [#lastUpdatedFormat]_ |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| sortOrder      | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")               |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| limit          | no       | Choose the maximum number of results to return                                                         |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| offset         | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit   |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| page           | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are           |
	|                |          | ``limit`` long and the first page is 1. If ``offset`` was defined, this query parameter has no         |
	|                |          | effect. ``limit`` must be defined to make use of ``page``.                                             |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| newerThan      | no       | Return only Capabilities that were most recently updated no earlier than this date/time, which may be  |
	|                |          | given as an :rfc:`3339`-formatted string or as number of nanoseconds since the Unix Epoch (midnight    |
	|                |          | on January 1\ :sup:`st` 1970 UTC).                                                                     |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+
	| olderThan      | no       | Return only Capabilities that were most recently updated no later than this date/time, which may be    |
	|                |          | given as an :rfc:`3339`-formatted string or as number of nanoseconds since the Unix Epoch (midnight    |
	|                |          | on January 1\ :sup:`st` 1970 UTC).                                                                     |
	+----------------+----------+--------------------------------------------------------------------------------------------------------+

.. versionadded:: ATCv6
	The ``newerThan`` and ``olderThan`` query string parameters were added to all API versions as of :abbr:`ATC (Apache Traffic Control)` version 6.0.

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/api_capabilities?capability=types-write HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:capability:  Capability name
:httpMethod:  An HTTP request method, practically one of:

	- ``GET``
	- ``POST``
	- ``PUT``
	- ``PATCH``
	- ``DELETE``

:httpRoute:   The request route for which this capability applies - relative to the Traffic Ops server's URL
:id:          An integer which uniquely identifies this capability
:lastUpdated: The time at which this capability was last updated, in :ref:`non-rfc-datetime`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 01 Nov 2018 14:45:24 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: wptErtIop/AfTTQ+1MZdA2YpPXEOuLFfrPQvvaHqO/uX5fRruOVYW+7p8JTrtH1xg1WN+x6FnjQnSHuWwcpyJg==
	Content-Length: 393

	{ "response": [
		{
			"httpMethod": "POST",
			"lastUpdated": "2018-11-01 14:10:22.794114+00",
			"httpRoute": "types",
			"id": 261,
			"capability": "types-write"
		},
		{
			"httpMethod": "PUT",
			"lastUpdated": "2018-11-01 14:10:22.795917+00",
			"httpRoute": "types/*",
			"id": 262,
			"capability": "types-write"
		},
		{
			"httpMethod": "DELETE",
			"lastUpdated": "2018-11-01 14:10:22.799748+00",
			"httpRoute": "types/*",
			"id": 263,
			"capability": "types-write"
		}
	]}

.. [#lastUpdatedFormat] Unlike the ``newerThan`` and ``olderThan`` query string parameters which can accept either RFC3339 strings or nanoseconds, this **must** be RFC3339 and **must not** have sub-second precision. This also means that the format of the returned ``lastUpdated`` fields on the actual response objects is unnacceptable as input for this query string parameter.
