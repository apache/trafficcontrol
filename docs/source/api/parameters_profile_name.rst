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

.. _to-api-parameters-profile-name:

*******************************
``parameters/profile/{{name}}``
*******************************
.. deprecated:: 1.1
	Use :ref:`to-api-profiles-name-name-parameters` instead

``GET``
=======
Gets details about a specific profile's parameters

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

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

	GET /api/1.4/parameters/profile/GLOBAL HTTP/1.1
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
	Whole-Content-Sha512: NudgZXUNyKNpmSFf856KEjyy+Pin/bFhG9NoRBDAxYbRKt2T5fF5Ze7sUNZfFI5n/ZZsgbx6Tsgtfd7oM6j+eg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 21:08:56 GMT
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
