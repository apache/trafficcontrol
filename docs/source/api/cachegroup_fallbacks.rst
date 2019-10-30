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

.. _to-api-cachegroup_fallbacks:

************************
``cachegroup_fallbacks``
************************

.. deprecated:: ATCv4.0

	The :ref:`to-api-cachegroups` and :ref:`to-api-cachegroups-id` endpoints now contain a list of "fallbacks" in the output, and support it in input, and so this endpoint is redundant.

``GET``
=======
Retrieve fallback-related configurations for a :term:`Cache Group`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------------+----------+-----------------------------------------------------------------------------------------------------------+
	| Name         | Required | Description                                                                                               |
	+==============+==========+===========================================================================================================+
	| cacheGroupId |yes\ [1]_ | The integral, unique identifier of a :term:`Cache Group` whose fallback configurations shall be retrieved |
	+--------------+----------+-----------------------------------------------------------------------------------------------------------+
	| fallbackId   |yes\ [1]_ | The integral, unique identifier of a fallback :term:`Cache Group`                                         |
	+--------------+----------+-----------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.3/cachegroup_fallbacks?cacheGroupId=11&fallbackId=7 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cacheGroupId:   The integral, unique identifier of the :term:`Cache Group` described by this entry
:cacheGroupName: The name of the :term:`Cache Group` described by this entry
:fallbackId:     The integral, unique identifier of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will fall back
:fallbackName:   The name of the :term:`Cache Group` on which the :term:`Cache Group` described by this entry will fall back
:fallbackOrder:  The order of the fallback described by "fallbackId" and "fallbackName" in the list of fallbacks for the :term:`Cache Group` described by this entry

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 14 Nov 2018 15:40:34 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Wed, 14 Nov 2018 19:40:34 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 9kauJ9tA4Ca5ElMHZk0fIJpQr+Wcx6NHiqWrnZJvyupRIOBQiUec3UW/fI9HdtE98xkrthz1daXKmdUkDhon8Q==
	Content-Length: 125

	{ "response": [
		{
			"cacheGroupId": 11,
			"fallbackOrder": 1,
			"fallbackName": "CDN_in_a_Box_Edge",
			"fallbackId": 7,
			"cacheGroupName": "test"
		}
	]}


.. [1] At least one of these must be provided, not necessarily both (though both is perfectly valid).

``POST``
========
Creates fallback configuration for a :term:`Cache Group`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
The request payload for this endpoint **must** be an array, even if only one fallback relationship is being created.

:cacheGroupId:  Integral, unique identifier of a :term:`Cache Group` to which to assign a fallback
:fallbackId:    Integral, unique identifier of a :term:`Cache Group` on which the :term:`Cache Group` identified by ``cacheGroupId`` will fall back
:fallbackOrder: The order of this fallback for the :term:`Cache Group` identified by ``cacheGroupId``

.. code-block:: http
	:caption: Request Example

	POST /api/1.3/cachegroup_fallbacks HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 59
	Content-Type: application/x-www-form-urlencoded

	[{"cacheGroupId": 11, "fallbackId": 7, "fallbackOrder": 1}]

Response Structure
------------------
:cacheGroupId:   The integral, unique identifier of the :term:`Cache Group` to which this fallback was assigned
:cacheGroupName: The name of the :term:`Cache Group` to which this fallback was assigned
:fallbackId:     The integral, unique identifier of the :term:`Cache Group` on which this entries :term:`Cache Group` will fall back
:fallbackName:   The name of the :term:`Cache Group` on which this entries :term:`Cache Group` will fall back
:fallbackOrder:  The order of the fallback described by "fallbackId" and "fallbackName" in the list of fallbacks for the :term:`Cache Group` described by this entry


.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 08 Nov 2018 14:59:46 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 08 Nov 2018 18:59:46 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 0twD50R5e7V2DtVrALQxzr2DmeHPPu8rTY8aGU4dFkx4XnOzjeRK5z+SYCrZEZ9Mh8QnWha3yZ2PtlxVTZt1YA==
	Content-Length: 225

	{ "alerts": [
		{
			"level": "success",
			"text": "Backup configuration CREATE for cache group 11 successful."
		}
	],
	"response": [
		{
			"cacheGroupId": 11,
			"fallbackName": "CDN_in_a_Box_Edge",
			"fallbackOrder": 1,
			"fallbackId": 7,
			"cacheGroupName": "test"
		}
	]}


``PUT``
=======
Updates an existing fallback configuration for one or more :term:`Cache Groups`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
The request payload for this endpoint **must** be an array, even if only one fallback relationship is being updated.
:cacheGroupId:  Integral, unique identifier of a :term:`Cache Group` to which to assign a fallback
:fallbackId:    Integral, unique identifier of a :term:`Cache Group` on which the :term:`Cache Group` identified by ``cacheGroupId`` will fall back
:fallbackOrder: The order of this fallback for the :term:`Cache Group` identified by ``cacheGroupId``

.. note:: The request data should be an array of these objects (and any number can be submitted per request), see the example

.. code-block:: http
	:caption: Request Example

	PUT /api/1.1/cachegroup_fallbacks HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 59
	Content-Type: application/x-www-form-urlencoded

	[{"cacheGroupId": 11, "fallbackId": 7, "fallbackOrder": 2}]

Response Structure
------------------
:cacheGroupId:   The integral, unique identifier of the :term:`Cache Group` to which this fallback was assigned
:cacheGroupName: The name of the :term:`Cache Group` to which this fallback was assigned
:fallbackId:     The integral, unique identifier of the :term:`Cache Group` on which this entries :term:`Cache Group` will fall back
:fallbackName:   The name of the :term:`Cache Group` on which this entries :term:`Cache Group` will fall back
:fallbackOrder:  The order of the fallback described by "fallbackId" and "fallbackName" in the list of fallbacks for the :term:`Cache Group` described by this entry

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 08 Nov 2018 15:07:06 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 08 Nov 2018 19:07:06 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 7QQpwDEmSpSPn6E3FAjxNw3E7xKP3TOBdnvZiBHQwOLmOH6Eiaq58f3eMPYAuK4qMSAKBj9Y2R//Fpa59YCMRw==
	Content-Length: 225

	{ "alerts": [
		{
			"level": "success",
			"text": "Backup configuration UPDATE for cache group 11 successful."
		}
	],
	"response": [
		{
			"cacheGroupId": 11,
			"fallbackName": "CDN_in_a_Box_Edge",
			"fallbackOrder": 2,
			"fallbackId": 7,
			"cacheGroupName": "test"
		}
	]}

``DELETE``
==========
Delete fallback list assigned to a :term:`Cache Group`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Query Parameters

	+--------------+----------+-----------------------------------------------------------------------------------------------------------+
	| Name         | Required | Description                                                                                               |
	+==============+==========+===========================================================================================================+
	| cacheGroupId |yes\ [2]_ | The integral, unique identifier of a :term:`Cache Group` whose fallback configurations shall be retrieved |
	+--------------+----------+-----------------------------------------------------------------------------------------------------------+
	| fallbackId   |yes\ [2]_ | The integral, unique identifier of a fallback :term:`Cache Group`                                         |
	+--------------+----------+-----------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.2/cachegroup_fallbacks?cacheGroupId=11&fallbackId=7 HTTP/1.1
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Thu, 08 Nov 2018 15:48:56 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 08 Nov 2018 19:48:56 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: MG2FNZ18EyAvy/IgdUPX4XRjJXYclXtp0e/kCMfimx9C427LNwjvL1seXkvu9crT2o68i0H2q1efshDJHO81IQ==
	Content-Length: 76

	{
		"response": "Backup Cachegroup 7  DELETED from cachegroup 11 fallback list"
	}


.. [2] At least one of "cacheGroupId" or "fallbackId" must be sent with the request. If both are sent, a single fallback relationship is deleted, whereas using only "cacheGroupId" will result in all fallbacks being removed from the :term:`Cache Group` identified by that integral, unique identifier, and using only "fallbackId" will remove the :term:`Cache Group` identified by *that* integral, unique identifier from all other :term:`Cache Groups`' fallback lists.
