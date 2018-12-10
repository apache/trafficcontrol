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

.. _to-api-profiles-name-name-parameters:

*************************************
``profiles/name/{{name}}/parameters``
*************************************

``GET``
=======
Retrieves all parameters associated with a given profile

:Auth. Required: Yes
:Roles Required: None
:Response Type:  None

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------+
	| Name | Description                                                 |
	+======+=============================================================+
	| name | The name of the profile for which parameters will be listed |
	+------+-------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/profiles/name/GLOBAL/parameters HTTP/1.1
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

**Response Example** ::

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: NudgZXUNyKNpmSFf856KEjyy+Pin/bFhG9NoRBDAxYbRKt2T5fF5Ze7sUNZfFI5n/ZZsgbx6Tsgtfd7oM6j+eg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 21:52:08 GMT
	Content-Length: 542

	{ "response": [
		{
			"configFile": "global",
			"id": 4,
			"lastUpdated": "2018-12-05 17:50:49+00",
			"name": "tm.instance_name",
			"secure": false,
			"value": "Traffic Ops CDN"
		},
		{
			"configFile": "global",
			"id": 5,
			"lastUpdated": "2018-12-05 17:50:49+00",
			"name": "tm.toolname",
			"secure": false,
			"value": "Traffic Ops"
		},
		{
			"configFile": "global",
			"id": 6,
			"lastUpdated": "2018-12-05 17:50:51+00",
			"name": "use_tenancy",
			"secure": false,
			"value": "1"
		},
		{
			"configFile": "regex_revalidate.config",
			"id": 7,
			"lastUpdated": "2018-12-05 17:50:49+00",
			"name": "maxRevalDurationDays",
			"secure": false,
			"value": "90"
		}
	]}

``POST``
========
Associate parameters to a profile. If the parameter does not exist, create it and associate to the profile. If the parameter already exists, associate it to the profile. If the parameter is already associated with the profile, keep the association.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------+
	| Name | Description                                                  |
	+======+==============================================================+
	| name | The name of the profile to which parameters will be assigned |
	+------+--------------------------------------------------------------+

This endpoint accepts two formats for the request payload:

Single Parameter Format
	Specify a single parameter to assign to the specified profile
Parameter Array Format
	Specify multiple parameters to assign to the specified profile

.. warning:: Most API endpoints dealing with parameters treat ``secure`` as a boolean value, whereas this endpoint takes the legacy approach of treating it as an integer. Be careful when passing data back and forth, as boolean values will **not** be accepted by this endpoint!

Single Parameter Format
"""""""""""""""""""""""
:configFile: The *base* filename of the configuration file to which this parameter shall belong e.g. "foo" not "/path/to/foo"
:name:       Parameter name
:secure:     An integer which, when any number other than ``0``, will prohibit users who do not have the "admin" role from viewing the parameter's ``value`` (at the time of this writing the obfuscation value is defined to be ``"********"``)
:value:      Parameter value

.. code-block:: http
	:caption: Request Example - Single Parameter Format

	POST /api/1.4/profiles/name/test/parameters HTTP/1.1
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
:configFile: The *base* filename of the configuration file to which this parameter shall belong e.g. "foo" not "/path/to/foo"
:name:       Parameter name
:secure:     An integer which, when any number other than ``0``, will prohibit users who do not have the "admin" role from viewing the parameter's ``value`` (at the time of this writing the obfuscation value is defined to be ``"********"``)
:value:      Parameter value

.. code-block:: http
	:caption: Request Example - Parameter Array Format

	POST /api/1.4/profiles/name/test/parameters HTTP/1.1
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
:parameters: An array of objects representing the parameters which have been assigned

	:configFile: The *base* filename of the configuration file to which this parameter shall belong e.g. "foo" not "/path/to/foo"
	:name:       Parameter name
	:secure:     An integer which, when any number other than ``0``, will prohibit users who do not have the "admin" role from viewing the parameter's ``value`` (at the time of this writing the obfuscation value is defined to be ``"********"``)
	:value:      Parameter value

:profileId:   The integral, unique identifier for the profile to which the parameter(s) have been assigned
:profileName: Name of the profile to which the parameter(s) have been assigned

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: HQWqOkoYHjvcxheWPrHOb0oZnUC+qLG1LO4OjtsLLnZYVUIu/qgJrzvziPnKq3FEHUWaZrnDCZM/iZD8AXOKBw==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 14:20:22 GMT
	Content-Length: 357

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
			},
			{
				"configFile": "bar",
				"name": "foo",
				"secure": 0,
				"value": "Another test parameter for API examples",
				"id": 129
			}
		],
		"profileId": 18,
		"profileName": "test"
	}}

.. note:: The format of the request does not affect the format of the response. ``parameters`` will be an array either way.
