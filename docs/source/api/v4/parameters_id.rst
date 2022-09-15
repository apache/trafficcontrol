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

.. _to-api-v4-parameters-id:

*********************
``parameters/{{ID}}``
*********************

``PUT``
=======
Replaces a :term:`Parameter`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PARAMETER:UPDATE, PARAMETER:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	|  ID  | The :ref:`parameter-id` of the :term:`Parameter` which will be deleted |
	+------+------------------------------------------------------------------------+

:configFile:  The :term:`Parameter`'s :ref:`parameter-config-file`
:name:        :ref:`parameter-name` of the :term:`Parameter`
:secure:      A boolean value that describes whether or not the :term:`Parameter` is :ref:`parameter-secure`
:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/parameters/124 HTTP/1.1
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
Deletes the specified :term:`Parameter`. If, however, the :term:`Parameter` is associated with one or more :term:`Profiles`, deletion will fail.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PARAMETER:DELETE, PARAMETER:READ
:Response TYpe:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	|  ID  | The :ref:`parameter-id` of the :term:`Parameter` which will be deleted |
	+------+------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/parameters/124 HTTP/1.1
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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
