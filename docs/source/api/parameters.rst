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

.. _to-api-parameters:

**************
``parameters``
**************

``GET``
=======
Gets all parameters configured in Traffic Ops

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+---------------------------------------------------+
	| Name       | Required | Description                                       |
	+============+==========+===================================================+
	| id         | no       | Filter parameters by integral, unique identifier  |
	+------------+----------+---------------------------------------------------+
	| name       | no       | Filter parameters by name                         |
	+------------+----------+---------------------------------------------------+
	| configFile | no       | Filter parameters by configuration file           |
	+------------+----------+---------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/parameters?configFile=records.config&name=location HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:configFile:  The *base* filename to which this parameter belongs
:id:          An integral, unique identifier for this parameter
:lastUpdated: The date and time at which this parameter was last modified in ISO format
:name:        The parameter name
:profiles:    An array of profile names that use this parameter
:secure:      When ``true``, the parameter value is visible only to "admin"-role users
:value:       The parameter value - if ``secure`` is true and the user does not have the "admin" role this will be obfuscated (at the time of this writing the obfuscation value is defined to be ``"********"``) but **not** missing

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
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
Creates one or more new parameters.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
The request body may be in one of two formats, a single parameter object or an array of parameter objects. Each parameter object shall have the following keys:

.. caution:: At the time of this writing, there is a bug in the Go rewrite of this endpoint such that the "array format" will not be accepted by the server. Watch `GitHub issue #3093 <https://github.com/apache/trafficcontrol/issues/3093>`_ for further developments

:name:       Parameter name
:configFile: The *base* filename of the configuration file to which this parameter shall belong e.g. "foo" not "/path/to/foo"
:secure:     A boolean value which, when ``true`` will prohibit users who do not have the "admin" role from viewing the parameter's ``value`` (at the time of this writing the obfuscation value is defined to be ``"********"``)
:value:      Parameter value

.. code-block:: http
	:caption: Request Example - Single Object Format

	POST /api/1.4/parameters HTTP/1.1
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

	POST /api/1.4/parameters HTTP/1.1
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
:configFile:  The *base* filename to which this parameter belongs
:id:          An integral, unique identifier for this parameter
:lastUpdated: The date and time at which this parameter was last modified in ISO format
:name:        The parameter name
:profiles:    An array of profile names that use this parameter - should be ``null`` immediately after parameter creation
:secure:      When ``true``, the parameter value is visible only to "admin"-role users
:value:       The parameter value - if ``secure`` is true and the user does not have the "admin" role this will be obfuscated (at the time of this writing the obfuscation value is defined to be ``"********"``) but **not** missing

.. code-block:: http
	:caption: Response Example - Single Object Format

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
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
