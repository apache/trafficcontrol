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

.. _to-api-parameters-id:

*********************
``parameters/{{ID}}``
*********************

``GET``
=======
Gets details about a specific parameter

.. deprecated:: 1.1
	Use the ``id`` query parameter of the :ref:`to-api-parameters` endpoint instead

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	|  ID  | The integral, unique identifier of the parameter which will be deleted |
	+------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/parameters/29 HTTP/1.1
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
	Date: Wed, 05 Dec 2018 19:01:54 GMT
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

``PUT``
=======
Replaces a parameter.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	|  ID  | The integral, unique identifier of the parameter which will be deleted |
	+------+------------------------------------------------------------------------+

:name:       Parameter name
:configFile: The *base* filename of the configuration file to which this parameter shall belong e.g. "foo" not "/path/to/foo"
:secure:     A boolean value which, when ``true`` will prohibit users who do not have the "admin" role from viewing the parameter's ``value`` (at the time of this writing the obfuscation value is defined to be ``"********"``)
:value:      Parameter value

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/parameters/124 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 81
	Content-Type: application/json

	{
		"name": "foo",
		"value": "bar",
		"configFile": "records.config",
		"secure": false
	}

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
	Whole-Content-Sha512: DMxS2gKceFVKRtezON/vsnrC+zI8onASSHaGv5i3wwvUvyt9KEe72gxQd6ZgVcSq3K8ZpkH6g3UI/WtEfdp5vA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 20:21:07 GMT
	Content-Length: 209

	{ "alerts": [
		{
			"text": "param was updated.",
			"level": "success"
		}
	],
	"response": {
		"configFile": "records.config",
		"id": 125,
		"lastUpdated": "2018-12-05 20:21:07+00",
		"name": "foo",
		"profiles": null,
		"secure": false,
		"value": "bar"
	}}

``DELETE``
==========
Deletes the specified parameter. If, however, the parameter is associated with one or more profiles, deletion will fail.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response TYpe:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	|  ID  | The integral, unique identifier of the parameter which will be deleted |
	+------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/parameters/124 HTTP/1.1
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: hJjQq2Seg7sqWt+jKgp6gwRxUtoVU34PFoc9wEaweXdaIBTn/BscoUuyw2/n+V8GZPqpeQcihZE50/0oQhdtHw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 19:20:30 GMT
	Content-Length: 60

	{ "alerts": [
		{
			"text": "param was deleted.",
			"level": "success"
		}
	]}
