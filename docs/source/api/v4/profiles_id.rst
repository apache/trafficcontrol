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

.. _to-api-v4-profiles-id:

*******************
``profiles/{{ID}}``
*******************

``PUT``
=======
Replaces the specified :term:`Profile` with the one in the request payload

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PROFILE:UPDATE, PROFILE:READ
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

	PUT /api/4.0/profiles/16 HTTP/1.1
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
:Permissions Required: PROFILE:DELETE, PROFILE:READ
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

	DELETE /api/4.0/profiles/16 HTTP/1.1
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
