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

.. _to-api-v1-parameters-profile-name:

*******************************
``parameters/profile/{{name}}``
*******************************

``GET``
=======
Gets details about a specific :term:`Profile`'s :term:`Parameters`

.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-profiles-name-name-parameters` instead.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------------------+
	| Name | Description                                                                                |
	+======+============================================================================================+
	| name | The :ref:`profile-name` of the :term:`Profile` for which :term:`Parameters` will be listed |
	+------+--------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/parameters/profile/GLOBAL HTTP/1.1
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
	Whole-Content-Sha512: NudgZXUNyKNpmSFf856KEjyy+Pin/bFhG9NoRBDAxYbRKt2T5fF5Ze7sUNZfFI5n/ZZsgbx6Tsgtfd7oM6j+eg==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 21:08:56 GMT
	Content-Length: 628

	{ "alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use /profiles/name/{name}/parameters instead"
		}],
			"response": [
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
