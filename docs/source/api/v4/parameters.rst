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

.. _to-api-v4-parameters:

**************
``parameters``
**************

``GET``
=======
Gets all :term:`Parameters` configured in Traffic Ops

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: PARAMETER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                                   |
	+=============+==========+===============================================================================================================+
	| configFile  | no       | Filter :term:`Parameters` by :ref:`parameter-config-file`                                                     |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| id          | no       | Filters :term:`Parameters` by :ref:`parameter-id`                                                             |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| name        | no       | Filter :term:`Parameters` by :ref:`parameter-name`                                                            |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| value       | no       | Filter :term:`Parameters` by :ref:`parameter-value`                                                           |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| orderby     | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|             |          | array                                                                                                         |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder   | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| limit       | no       | Choose the maximum number of results to return                                                                |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| offset      | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit.         |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+
	| page        | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|             |          | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|             |          | defined to make use of ``page``.                                                                              |
	+-------------+----------+---------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/parameters?configFile=records.config&name=location HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:id:          The :term:`Parameter`'s :ref:`parameter-id`
:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in :ref:`non-rfc-datetime`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:profiles:    An array of :term:`Profile` :ref:`Names <profile-name>` that use this :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: UFO3/jcBFmFZM7CsrsIwTfPc5v8gUiXqJm6BNp1boPb4EQBnWNXZh/DbBwhMAOJoeqDImoDlrLnrVjQGO4AooA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 18:23:39 GMT
	Content-Length: 212

	{ "response": [
		{
			"configFile": "records.config",
			"id": 29,
			"lastUpdated": "2018-12-05 17:51:02+00",
			"name": "location",
			"profiles": [
				"ATS_EDGE_TIER_CACHE",
				"ATS_MID_TIER_CACHE"
			],
			"secure": false,
			"value": "/etc/trafficserver/"
		}
	]}

``POST``
========
Creates one or more new :term:`Parameters`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PARAMETER:CREATE, PARAMETER:READ
:Response Type:  Array

Request Structure
-----------------
The request body may be in one of two formats, a single :term:`Parameter` object or an array of :term:`Parameter` objects. Each :term:`Parameter` object shall have the following keys:

:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Request Example - Single Object Format

	POST /api/4.0/parameters HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 84
	Content-Type: application/json

	{
		"name": "test",
		"value": "quest",
		"configFile": "records.config",
		"secure": false
	}

.. code-block:: http
	:caption: Request Example - Array Format

	POST /api/4.0/parameters HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 180
	Content-Type: application/json

	[{
		"name": "test",
		"value": "quest",
		"configFile": "records.config",
		"secure": false
	},
	{
		"name": "foo",
		"value": "bar",
		"configFile": "records.config",
		"secure": false
	}]

Response Structure
------------------
:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:id:          The :term:`Parameter`'s :ref:`parameter-id`
:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in :ref:`non-rfc-datetime`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:profiles:    An array of :term:`Profile` :ref:`Names <profile-name>` that use this :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Response Example - Single Object Format

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: eQrl48zWids0kDpfCYmmtYMpegjnFxfOVvlBYxxLSfp7P7p6oWX4uiC+/Cfh2X9i3G+MQ36eH95gukJqOBOGbQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 19:18:21 GMT
	Content-Length: 212

	{ "alerts": [
		{
			"text": "param was created.",
			"level": "success"
		}
	],
	"response": {
		"configFile": "records.config",
		"id": 124,
		"lastUpdated": "2018-12-05 19:18:21+00",
		"name": "test",
		"profiles": null,
		"secure": false,
		"value": "quest"
	}}
