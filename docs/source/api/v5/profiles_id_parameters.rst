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

.. _to-api-profiles-id-parameters:

******************************
``profiles/{{ID}}/parameters``
******************************

``GET``
=======

Retrieves all :term:`Parameters` assigned to the :term:`Profile`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: PROFILE:READ, PARAMETER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------------------+
	| Name | Description                                                                              |
	+======+==========================================================================================+
	|  ID  | The :ref:`profile-id` of the :term:`Profile` for which :term:`Parameters` will be listed |
	+------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/parameters/profile/GLOBAL HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:id:          The :term:`Parameter`'s :ref:`parameter-id`
:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in :rfc:`3339` format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

:name:     :ref:`parameter-name` of the :term:`Parameter`
:profiles: An array of :term:`Profile` :ref:`Names <profile-name>` that use this :term:`Parameter`
:secure:   A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:    The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: NudgZXUNyKNpmSFf856KEjyy+Pin/bFhG9NoRBDAxYbRKt2T5fF5Ze7sUNZfFI5n/ZZsgbx6Tsgtfd7oM6j+eg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 21:08:56 GMT
	Content-Length: 542

	{ "response": [
		{
			"configFile": "global",
			"id": 4,
			"lastUpdated": "2018-12-05T23:52:59.696337+05:30",
			"name": "tm.instance_name",
			"secure": false,
			"value": "Traffic Ops CDN"
		},
		{
			"configFile": "global",
			"id": 5,
			"lastUpdated": "2018-12-05T23:52:59.696337+05:30",
			"name": "tm.toolname",
			"secure": false,
			"value": "Traffic Ops"
		},
		{
			"configFile": "regex_revalidate.config",
			"id": 7,
			"lastUpdated": "2018-12-05T23:52:59.696337+05:30",
			"name": "maxRevalDurationDays",
			"secure": false,
			"value": "90"
		}
	]}

``POST``
========

Associates :term:`Parameters` to a :term:`Profile`. If the :term:`Parameter` does not exist, creates it and associates it to the :term:`Profile`. If the :term:`Parameter` already exists, associates it to the :term:`Profile`. If the :term:`Parameter` is already associated with the :term:`Profile`, keep the association.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PROFILE:UPDATE, PROFILE:READ, PARAMETER:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------------+
	| Name | Description                                                                               |
	+======+===========================================================================================+
	|  ID  | The :ref:`profile-id` of the :term:`Profile` to which :term:`Parameters` will be assigned |
	+------+-------------------------------------------------------------------------------------------+

This endpoint accepts two formats for the request payload:

Single Object Format
	For assigning a single :term:`Parameter` to a single :term:`Profile`
Parameter Array Format
	For making multiple assignments of :term:`Parameters` to :term:`Profiles` simultaneously

.. warning:: Most API endpoints dealing with :term:`Parameters` treat :ref:`parameter-secure` as a boolean value, whereas this endpoint takes the legacy approach of treating it as an integer. Be careful when passing data back and forth, as boolean values will **not** be accepted by this endpoint!

Single Parameter Format
"""""""""""""""""""""""
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Response Example - Single Parameter Format

	POST /api/5.0/profiles/18/parameters HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 99
	Content-Type: application/json

	{
		"name": "test",
		"configFile": "quest",
		"value": "A test parameter for API examples",
		"secure": 0
	}


Parameter Array Format
""""""""""""""""""""""
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Request Example - Parameter Array Format

	POST /api/5.0/profiles/18/parameters HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 212
	Content-Type: application/json

	[{
		"name": "test",
		"configFile": "quest",
		"value": "A test parameter for API examples",
		"secure": 0
	},
	{
		"name": "foo",
		"configFile": "bar",
		"value": "Another test parameter for API examples",
		"secure": 0
	}]

Response Structure
------------------
:parameters: An array of objects representing the :term:`Parameters` which have been assigned

	:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
	:name:        :ref:`parameter-name` of the :term:`Parameter`
	:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
	:value:       The :term:`Parameter`'s :ref:`parameter-value`

:profileId:   The :ref:`profile-id` of the :term:`Profile` to which the :term:`Parameter`\ (s) have been assigned
:profileName: :ref:`profile-name` of the :term:`Profile` to which the :term:`Parameter`\ (s) have been assigned

.. code-block:: http
	:caption: Response Example - Single Parameter Format

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: R2QUyCaNvKvVv/PNVNmEd/ma5h/iP1fMJlqhv+x2jE/zNpHJ1KVXt6s3btB8nnHv6IDF/gt5kIzQ0mbW5U8bpg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 14:45:28 GMT
	Content-Length: 253

	{ "alerts": [
		{
			"text": "Assign parameters successfully to profile test",
			"level": "success"
		}
	],
	"response": {
		"parameters": [
			{
				"configFile": "quest",
				"name": "test",
				"secure": 0,
				"value": "A test parameter for API examples",
				"id": 126
			}
		],
		"profileId": 18,
		"profileName": "test"
	}}

.. note:: The format of the request does not affect the format of the response. ``parameters`` will be an array either way.
