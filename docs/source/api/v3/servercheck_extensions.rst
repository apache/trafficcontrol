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

.. _to-api-v3-servercheck_extensions:

**************************
``servercheck/extensions``
**************************
.. seealso:: :ref:`admin-to-ext-script`

``GET``
=======
Retrieves the list of Traffic Ops extensions.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| Name             | Required | Description                                                                                                                  |
	+==================+==========+==============================================================================================================================+
	| id               | no       | Filter TO Extensions by the integral, unique identifier of an Extension                                                      |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| name             | no       | Filter TO Extensions by the name of an Extension                                                                             |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| script_file      | no       | Filter TO Extensions by the base filename of the script that runs for the Extension                                          |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| isactive         | no       | Boolean used to return either only active (1) or inactive(0) TO Extensions                                                   |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| type             | no       | Filter TO Extensions by the type of Extension                                                                                |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| sortOrder        | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                     |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| limit            | no       | Choose the maximum number of results to return                                                                               |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| offset           | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit.                        |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+
	| page             | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first    |
	|                  |          | page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``. |
	+------------------+----------+------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/servercheck/extensions HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:additional_config_json: A string containing a JSON-encoded object with extra configuration options... inside a JSON object...
:description:            A short description of the extension

	.. note:: This is, unfortunately, ``null`` for all default extensions

:id:       An integral, unique identifier for this extension definition
:info_url: A URL where info about this extension may be found
:isactive: An integer describing the boolean notion of whether or not the extension is active; one of:

	0
		disabled
	1
		enabled

:name:                  The name of the extension
:script_file:           The base filename of the script that runs for the extension
:servercheck_shortname: The name of the column in the table at 'Monitor' -> 'Cache Checks' in Traffic Portal, where "Check Extension" output is displayed
:type:                  The Check :term:`Type` of the extension. This will always be a CHECK_EXTENSION type with the naming convention of ``CHECK_EXTENSION_*``.
:version:               A (hopefully) semantic version number describing the version of the plugin

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Tue, 11 Dec 2018 20:51:48 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: n73jg9XR4V5Cwqq56Rf3wuIi99k3mM5u2NAjcZ/gQBu8jvAFymDlnZqKeJ+wTll1vjIsHpXCOVXV7+5UGakLgA==
	Transfer-Encoding: chunked

	{ "response": [
		{
			"script_file": "ToPingCheck.pl",
			"version": "1.0.0",
			"name": "ILO_PING",
			"description": null,
			"info_url": "-",
			"additional_config_json": "{ check_name: \"ILO\", \"base_url\": \"https://localhost\", \"select\": \"ilo_ip_address\", \"cron\": \"9 * * * *\" }",
			"isactive": 1,
			"type": "CHECK_EXTENSION_BOOL",
			"id": 1,
			"servercheck_short_name": "ILO"
		},
		{
			"script_file": "ToPingCheck.pl",
			"version": "1.0.0",
			"name": "10G_PING",
			"description": null,
			"info_url": "-",
			"additional_config_json": "{ check_name: \"10G\", \"base_url\": \"https://localhost\", \"select\": \"ip_address\", \"cron\": \"18 * * * *\" }",
			"isactive": 1,
			"type": "CHECK_EXTENSION_BOOL",
			"id": 2,
			"servercheck_short_name": "10G"
		}
	]}

``POST``
========
Creates a new Traffic Ops check extension.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  ``undefined``

Request Structure
-----------------
:additional_config_json: An optional string containing a JSON-encoded object with extra configuration options... inside a JSON object...
:description:            A short description of the extension
:info_url:               A URL where info about this extension may be found
:isactive:               An integer describing the boolean notion of whether or not the extension is active; one of:

	0
		disabled
	1
		enabled

:name:        The name of the extension
:script_file: The base filename of the script that runs for the extension

	.. seealso:: :ref:`admin-to-ext-script` for details on where the script should be located on the Traffic Ops server

:servercheck_shortname: The name of the column in the table at 'Monitor' -> 'Cache Checks' in Traffic Portal, where "Check Extension" output is displayed
:type:                  The :term:`Type` of extension. Must be CHECK_EXTENSION type with the naming convention of ``CHECK_EXTENSION_*``.
:version:               A (hopefully) semantic version number describing the version of the plugin

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/servercheck/extensions HTTP/1.1
	Host: cache1.example.com:6443
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 208
	Content-Type: application/json

	{
		"name": "test",
		"version": "0.0.1-1",
		"info_url": "",
		"script_file": "",
		"isactive": 0,
		"description": "A test extension for API examples",
		"servercheck_short_name": "test",
		"type": "CHECK_EXTENSION_NUM"
	}


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
	Date: Wed, 12 Dec 2018 16:37:44 GMT
	X-Server-Name: traffic_ops_golang/
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: 7M67PYnli6WzGQFS3g8Gh1SOyq6VENZMqm/kUffOTLLFfuWSEuSLA65R5R+VyJiNjdqOG5Bp78mk+JYcqhtVGw==
	Content-Length: 89

	{ "supplemental":
		{
			"id": 5
		},
	"alerts": [{
		"level": "success",
		"text": "Check Extension Loaded."
	}]}

.. [1] No roles are required to use this endpoint, however access is controlled by username. Only the reserved user ``extension`` is permitted the use of this endpoint.
