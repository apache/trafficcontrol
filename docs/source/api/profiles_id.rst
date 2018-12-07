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

.. _to-api-profiles-id:

*******************
``profiles/{{ID}}``
*******************

``GET``
=======
.. deprecated:: 1.1
	Use the ``id`` query parameter of a ``GET`` request to :ref:`to-api-profiles-id` instead.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+----------------------------------------------------------------+
	| Parameter |                           Description                          |
	+===========+================================================================+
	|    id     | The integral, unique identifier of the profile to be retrieved |
	+-----------+----------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.1/profiles/9 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cdn:         The integral, unique identifier of the CDN to which this profile belongs
:cdnName:     The CDN name
:description: A description of the profile
:id:          The integral, unique identifier of this profile
:lastUpdated: The date and time at which this profile was last updated
:name:        The name of the profile
:params:      An array of parameters in use by this profile

	:configFile:  The *base* filename to which this parameter belongs
	:id:          An integral, unique identifier for this parameter
	:lastUpdated: The date and time at which this parameter was last modified in ISO format
	:name:        The parameter name
	:profiles:    An array of profile names that use this parameter
	:secure:      When ``true``, the parameter value is visible only to "admin"-role users
	:value:       The parameter value - if ``secure`` is true and the user does not have the "admin" role this will be obfuscated (at the time of this writing the obfuscation value is defined to be ``"********"``) but **not** missing

:routingDisabled: A boolean which, if ``true`` will disable Traffic Router's routing to servers using this profile
:type:            The name of the 'type' of the profile

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: LCdG7AgeHqL4wpGraaoN8ks+/gYW//h1Q2OVBECk+T9/IC6tbJ3DWOgWX4u4dpudIDJ5mhRwBzicYvyyXWj3qA==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 07 Dec 2018 21:06:30 GMT
	Transfer-Encoding: chunked


	{ "response": [{
		"id": 9,
		"lastUpdated": "2018-12-05 17:51:00+00",
		"name": "ATS_EDGE_TIER_CACHE",
		"description": "Edge Cache - Apache Traffic Server",
		"cdnName": "CDN-in-a-Box",
		"cdn": 2,
		"routingDisabled": false,
		"type": "ATS_PROFILE",
		"params": [
			{
				"configFile": "records.config",
				"id": 9,
				"lastUpdated": null,
				"name": "CONFIG proxy.config.config_dir",
				"profiles": null,
				"secure": false,
				"value": "STRING /etc/trafficserver"
			},
			{
				"configFile": "records.config",
				"id": 10,
				"lastUpdated": null,
				"name": "CONFIG proxy.config.admin.user_id",
				"profiles": null,
				"secure": false,
				"value": "STRING ats"
			}
		]
	}]}

.. note:: The response example for this endpoint has been truncated to only the first two elements of the resulting ``params`` array, as the output was hundreds of lines long.

``PUT``
=======
Replaces the specified profile with the one in the response payload

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------+
	| Name | Description                                                   |
	+======+===============================================================+
	|  ID  | The integral, unique identifier of the profile being modified |
	+------+---------------------------------------------------------------+

:name:            New of the name profile
:description:     A new description of the new profile
:cdn:             The integral, unique identifier of the CDN to which the profile shall be assigned
:type:            The type of the profile

	.. warning:: Changing this will likely break something, be **VERY** careful when modifying this value

:routingDisabled: A boolean which, if ``true``, will prevent the Traffic Router from directing traffic to any servers assigned this profile

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/profiles/16 HTTP/1.1
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
:cdn:             The integral, unique identifier of the CDN to which this profile belongs
:cdnName:         The CDN name
:description:     A description of the profile
:id:              The integral, unique identifier of this profile
:lastUpdated:     The date and time at which this profile was last updated
:name:            The name of the profile
:routingDisabled: A boolean which, if ``true`` will disable Traffic Router's routing to servers using this profile
:type:            The name of the 'type' of the profile

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: Pnf+G9G3/+edt4b8PVsyGZHsNzaFEgphaGSminjRlRmMpWtuLAA20WZDUo3nX0QO81c2GCuFuEh9uMF2Vjeppg==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 07 Dec 2018 21:45:06 GMT
	Content-Length: 251

	{ "alerts": [
		{
			"text": "profile was updated.",
			"level": "success"
		}
	],
	"response": {
		"id": 16,
		"lastUpdated": "2018-12-07 21:45:06+00",
		"name": "test",
		"description": "A test profile for API examples",
		"cdnName": null,
		"cdn": 2,
		"routingDisabled": true,
		"type": "UNK_PROFILE"
	}}


``DELETE``
==========
Allows user to delete a profile.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------+
	| Name | Description                                                  |
	+======+==============================================================+
	|  ID  | The integral, unique identifier of the profile being deleted |
	+------+--------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/profiles/16 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
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
	Whole-Content-Sha512: HNmJkZaNW9yil08/3TnqZ5FllH6Rp+jgp3KI46FZdojLYcu+8jEhDLl1okoirdrHyU4R1c3hjCI0urN7PVvWDA==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 07 Dec 2018 21:55:33 GMT
	Content-Length: 62

	{ "alerts": [
		{
			"text": "profile was deleted.",
			"level": "success"
		}
	]}
