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

.. _to-api-profiles:

************
``profiles``
************

``GET``
=======
:Auth. Required: Yes
:Roles Required: None
:Permissions Required: PROFILE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------+----------+--------------------------------------------------------------------------------------------------------+
	|  Name | Required | Description                                                                                            |
	+=======+==========+========================================================================================================+
	|  cdn  |   no     | Used to filter :term:`Profiles` by the integral, unique identifier of the CDN to which they belong     |
	+-------+----------+--------------------------------------------------------------------------------------------------------+
	|  id   |   no     | Filters :term:`Profiles` by :ref:`profile-id`                                                          |
	+-------+----------+--------------------------------------------------------------------------------------------------------+
	| name  |   no     | Filters :term:`Profiles` by :ref:`profile-name`                                                        |
	+-------+----------+--------------------------------------------------------------------------------------------------------+
	| param |   no     | Used to filter :term:`Profiles` by the :ref:`parameter-id` of a :term:`Parameter` associated with them |
	+-------+----------+--------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/profiles?name=ATS_EDGE_TIER_CACHE HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cdn:         The integral, unique identifier of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:cdnName:     The name of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:description: The :term:`Profile`'s :ref:`profile-description`
:id:          The :term:`Profile`'s :ref:`profile-id`
:lastUpdated: The date and time at which this :term:`Profile` was last updated, in :rfc:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name:            The :term:`Profile`'s :ref:`profile-name`
:routingDisabled: The :term:`Profile`'s :ref:`profile-routing-disabled` setting
:type:            The :term:`Profile`'s :ref:`profile-type`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: QEpKM/DwHBRvue9K7XKrpwKFKhW6yCMQ2vSQgxE7dWFGJaqC4KOUO92bsJU/5fjI9qlB+1uMT2kz6mFb1Wzp/w==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 07 Dec 2018 20:40:31 GMT
	Content-Length: 220

	{ "response": [
		{
			"id": 9,
			"lastUpdated": "2023-08-17T09:06:18-06:00",
			"name": "ATS_EDGE_TIER_CACHE",
			"description": "Edge Cache - Apache Traffic Server",
			"cdnName": "CDN-in-a-Box",
			"cdn": 2,
			"routingDisabled": false,
			"type": "ATS_PROFILE"
		}
	]}

``POST``
========
Creates a new :term:`Profile`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PROFILE:CREATE, PROFILE:READ
:Response Type:  Object

Request Structure
-----------------
:cdn:             The integral, unique identifier of the :ref:`profile-cdn` to which this :term:`Profile` shall belong
:description:     The :term:`Profile`'s :ref:`profile-description`
:name:            The :term:`Profile`'s :ref:`profile-name`
:routingDisabled: The :term:`Profile`'s :ref:`profile-routing-disabled` setting
:type:            The :term:`Profile`'s :ref:`profile-type`

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/profiles HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 125
	Content-Type: application/json

	{
		"name": "test",
		"description": "A test profile for API examples",
		"cdn": 2,
		"type": "UNK_PROFILE",
		"routingDisabled": true
	}

Response Structure
------------------
:cdn:         The integral, unique identifier of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:cdnName:     The name of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:description: The :term:`Profile`'s :ref:`profile-description`
:id:          The :term:`Profile`'s :ref:`profile-id`
:lastUpdated: The date and time at which this :term:`Profile` was last updated, in :rfc:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name:            The :term:`Profile`'s :ref:`profile-name`
:routingDisabled: The :term:`Profile`'s :ref:`profile-routing-disabled` setting
:type:            The :term:`Profile`'s :ref:`profile-type`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: UGV3PCnYBY0J3siICR0f9VVRNdUK1+9zsDDP6T9yt6t+AoHckHe6bvzOli9to/fGhC2zz5l9Nc1ro4taJUDD8g==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 07 Dec 2018 21:24:49 GMT
	Content-Length: 251

	{ "alerts": [
		{
			"text": "profile was created.",
			"level": "success"
		}
	],
	"response": {
		"id": 16,
		"lastUpdated": "2023-08-17T09:06:18-06:00",
		"name": "test",
		"description": "A test profile for API examples",
		"cdnName": null,
		"cdn": 2,
		"routingDisabled": true,
		"type": "UNK_PROFILE"
	}}
