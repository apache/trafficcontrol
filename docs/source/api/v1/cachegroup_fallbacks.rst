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

.. _to-api-v1-cachegroup_fallbacks:

************************
``cachegroup_fallbacks``
************************

.. deprecated:: 1.4

	The :ref:`to-api-v1-cachegroups` and :ref:`to-api-v1-cachegroups-id` endpoints now contain a list of :ref:`cache-group-fallbacks` in the output, and support it in input, and so this endpoint is redundant.

``GET``
=======
Retrieve the :ref:`cache-group-fallbacks` of a :term:`Cache Group`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------------+----------+-----------------------------------------------------------------------------------------------------------+
	| Name         | Required | Description                                                                                               |
	+==============+==========+===========================================================================================================+
	| cacheGroupId | yes      | The :ref:`cache-group-id` of a :term:`Cache Group` whose :ref:`cache-group-fallbacks` shall be retrieved  |
	+--------------+----------+-----------------------------------------------------------------------------------------------------------+
	| fallbackId   | no       | The integral, unique identifier of a single :ref:`"fallback" <cache-group-fallbacks>` :term:`Cache Group` |
	+--------------+----------+-----------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/cachegroup_fallbacks?cacheGroupId=7 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:cacheGroupId:   An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` described by this entry
:cacheGroupName: The :ref:`cache-group-name` of the :term:`Cache Group` described by this entry as a string
:fallbackId:     An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will "fall back"
:fallbackName:   The :ref:`cache-group-name` of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will "fall back" as a string
:fallbackOrder:  The place in the list of :ref:`cache-group-fallbacks` of the :term:`Cache Group` identified by ``cacheGroupId`` and ``cacheGroupName`` where the :term:`Cache Group` identified by ``fallbackId`` and ``fallbackName`` starting from index 1.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Encoding: gzip
	Content-Length: 189
	Content-Type: application/json
	Date: Mon, 02 Dec 2019 22:26:27 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 03 Dec 2019 02:26:27 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: zSAeB8nxonyinsg1/at/l0/9FRRPw7N27DpkcZxRIwEzDOEY5XVfYcCHHFg1d/Q2JWtWZ9iRhs8mK5rLbKkccw==

	{ "alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'GET /cachegroups' instead"
		}
	],
	"response": [
		{
			"cacheGroupId": 7,
			"fallbackOrder": 2,
			"fallbackName": "test",
			"fallbackId": 8,
			"cacheGroupName": "CDN_in_a_Box_Edge"
		}
	]}


``POST``
========
Creates :ref:`"fallback" <cache-group-fallbacks>` configuration for a :term:`Cache Group`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
The request payload for this endpoint **must** be an array, even if only one "fallback" relationship is being created.

:cacheGroupId:  An integer that is the :ref:`cache-group-id` of a :term:`Cache Group` to which to assign a :ref:`fallback <cache-group-fallbacks>`
:fallbackId:    An integer that is the :ref:`cache-group-id` of a :term:`Cache Group` on which the :term:`Cache Group` identified by ``cacheGroupId`` will "fall back"
:fallbackOrder:  The place in the list of :ref:`cache-group-fallbacks` of the :term:`Cache Group` identified by ``cacheGroupId`` and ``cacheGroupName`` where the :term:`Cache Group` identified by ``fallbackId`` and ``fallbackName`` starting from index 1.

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/cachegroup_fallbacks HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 57

	[{"cacheGroupId": 7, "fallbackId": 8, "fallbackOrder": 2}]

Response Structure
------------------
:cacheGroupId:   An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` described by this entry
:cacheGroupName: The :ref:`cache-group-name` of the :term:`Cache Group` described by this entry as a string
:fallbackId:     An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will "fall back"
:fallbackName:   The :ref:`cache-group-name` of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will "fall back" as a string
:fallbackOrder:  The place in the list of :ref:`cache-group-fallbacks` of the :term:`Cache Group` identified by ``cacheGroupId`` and ``cacheGroupName`` where the :term:`Cache Group` identified by ``fallbackId`` and ``fallbackName`` starting from index 1.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Encoding: gzip
	Content-Length: 174
	Content-Type: application/json
	Date: Mon, 02 Dec 2019 22:23:22 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 03 Dec 2019 02:23:22 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: S8CMeR3P22itBNYOQaIjiQPMDoq2AzGt0/oBYpMPm1b8/iKeZfGSS4zyt4WYbVJrgrzFZYGUhBEJe6uimQYdCQ==

	{ "alerts": [
		{
			"level": "success",
			"text": "Backup configuration CREATE for cache group 7 successful."
		},
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'POST /cachegroups with a non-empty 'fallbacks' array' instead"
		}
	]}


``PUT``
=======
Updates an existing :ref:`fallback <cache-group-fallbacks>` configuration for one or more :term:`Cache Groups`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
The request payload for this endpoint **must** be an array, even if only one fallback relationship is being updated.

:cacheGroupId:  An integer that is the :ref:`cache-group-id` of a :term:`Cache Group` to which to assign a :ref:`fallback <cache-group-fallbacks>`
:fallbackId:    An integer that is the :ref:`cache-group-id` of a :term:`Cache Group` on which the :term:`Cache Group` identified by ``cacheGroupId`` will "fall back"
:fallbackOrder:  The place in the list of :ref:`cache-group-fallbacks` of the :term:`Cache Group` identified by ``cacheGroupId`` and ``cacheGroupName`` where the :term:`Cache Group` identified by ``fallbackId`` and ``fallbackName`` starting from index 1.

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/cachegroup_fallbacks HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 58

	[{"cacheGroupId": 7, "fallbackId": 8, "fallbackOrder": 2}]

Response Structure
------------------
:cacheGroupId:   An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` described by this entry
:cacheGroupName: The :ref:`cache-group-name` of the :term:`Cache Group` described by this entry as a string
:fallbackId:     An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will "fall back"
:fallbackName:   The :ref:`cache-group-name` of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will "fall back" as a string
:fallbackOrder:  The place in the list of :ref:`cache-group-fallbacks` of the :term:`Cache Group` identified by ``cacheGroupId`` and ``cacheGroupName`` where the :term:`Cache Group` identified by ``fallbackId`` and ``fallbackName`` starting from index 1.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Encoding: gzip
	Content-Length: 237
	Content-Type: application/json
	Date: Mon, 02 Dec 2019 22:28:55 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 03 Dec 2019 02:28:55 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: /rGLP3gbnqFUjDhC/4mSYr2a2HoVsGTukxHX8CbURnwDS5LV7U6gwvlOcgtMfEUyX1FEa4+1Xa94tiL/dRFj6w==

	{ "alerts": [
		{
			"level": "success",
			"text": "Backup configuration UPDATE for cache group 7 successful."
		},
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'PUT /cachegroups' instead"
		}
	],
	"response": [
		{
			"cacheGroupId": 7,
			"fallbackOrder": 2,
			"fallbackName": "test",
			"fallbackId": 8,
			"cacheGroupName": "CDN_in_a_Box_Edge"
		}
	]}

``DELETE``
==========
Remove one or more :ref:`cache-group-fallbacks` from one or more :term:`Cache Groups`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------------+----------+--------------------------------------------------------------------------------------------------------------+
	| Name         | Required | Description                                                                                                  |
	+==============+==========+==============================================================================================================+
	| cacheGroupId |yes\ [2]_ | The :ref:`cache-group-id` of a :term:`Cache Group` from which :ref:`cache-group-fallbacks` are being removed |
	+--------------+----------+--------------------------------------------------------------------------------------------------------------+
	| fallbackId   |yes\ [2]_ | The :ref:`cache-group-id` of a :ref:`"fallback" <cache-group-fallbacks>` :term:`Cache Group`                 |
	+--------------+----------+--------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/cachegroup_fallbacks?fallbackId=8 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

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
	Content-Encoding: gzip
	Content-Length: 186
	Content-Type: application/json
	Date: Mon, 02 Dec 2019 22:30:58 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Tue, 03 Dec 2019 02:30:58 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: iag1k8Ym4K6nrpahJwzyA45m2RO6159gSRg4ozUvg69/TKrTLyggMeAIVdzbwn8+ayOFq01lTK1Ho9jQFJ5j2w==

	{ "alerts": [
		{
			"level": "success",
			"text": "Cachegroup 8 DELETED from all the configured fallback lists"
		},
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'PUT /cachegroups with an empty 'fallbacks' array' instead"
		}
	]}

.. [2] At least one of "cacheGroupId" or "fallbackId" must be sent with the request. If both are sent, a single fallback relationship is deleted, whereas using only "cacheGroupId" will result in all fallbacks being removed from the :term:`Cache Group` identified by that integral, unique identifier, and using only "fallbackId" will remove the :term:`Cache Group` identified by *that* integral, unique identifier from all other :term:`Cache Groups`' fallback lists.
