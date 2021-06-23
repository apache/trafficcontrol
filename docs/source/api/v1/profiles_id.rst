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

.. _to-api-v1-profiles-id:

*******************
``profiles/{{ID}}``
*******************

``GET``
=======
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-v1-profiles` with the query parameter ``id`` instead.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+--------------------------------------------------------------+
	| Parameter | Description                                                  |
	+===========+==============================================================+
	|    id     | The :ref:`profile-id` of the :term:`Profile` to be retrieved |
	+-----------+--------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.1/profiles/9 HTTP/1.1
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
:lastUpdated: The date and time at which this :term:`Profile` was last updated, in :ref:`non-rfc-datetime`
:name:        The :term:`Profile`'s :ref:`profile-name`
:params:      An array of :term:`Parameters` in use by this :term:`Profile`

	:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
	:id:          The :term:`Parameter`'s :ref:`parameter-id`
	:lastUpdated: The date and time at which this :term:`Parameter` was last updated, in :ref:`non-rfc-datetime`
	:name:        :ref:`parameter-name` of the :term:`Parameter`
	:profiles:    An array of :term:`Profile` :ref:`Names <profile-name>` that use this :term:`Parameter`
	:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
	:value:       The :term:`Parameter`'s :ref:`parameter-value`

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
	}],
	"alerts": [
		{
			"text": "This endpoint is deprecated, please use GET /profiles with query parameter id instead",
			"level": "warning"
		}
	]}

.. note:: The response example for this endpoint has been truncated to only the first two elements of the resulting ``params`` array, as the output was hundreds of lines long.

``PUT``
=======
Replaces the specified :term:`Profile` with the one in the request payload

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------+
	| Name | Description                                                 |
	+======+=============================================================+
	|  ID  | The :ref:`profile-id` of the :term:`Profile` being modified |
	+------+-------------------------------------------------------------+

:cdn:             The integral, unique identifier of the :ref:`profile-cdn` to which this :term:`Profile` will belong
:description:     The :term:`Profile`'s new :ref:`profile-description`
:name:            The :term:`Profile`'s new :ref:`profile-name`
:routingDisabled: The :term:`Profile`'s new :ref:`profile-routing-disabled` setting
:type:            The :term:`Profile`'s new :ref:`profile-type`

	.. warning:: Changing this will likely break something, be **VERY** careful when modifying this value

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
:cdn:             The integral, unique identifier of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:cdnName:         The name of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:description:     The :term:`Profile`'s :ref:`profile-description`
:id:              The :term:`Profile`'s :ref:`profile-id`
:lastUpdated:     The date and time at which this :term:`Profile` was last updated, in :ref:`non-rfc-datetime`
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
Allows user to delete a :term:`Profile`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------+
	| Name | Description                                                |
	+======+============================================================+
	|  ID  | The :ref:`profile-id` of the :term:`Profile` being deleted |
	+------+------------------------------------------------------------+

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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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
