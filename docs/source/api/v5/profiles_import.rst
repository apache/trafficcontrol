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

.. _to-api-profiles-import:

*******************
``profiles/import``
*******************

``POST``
========

Imports a :term:`Profile` that was exported via :ref:`to-api-profiles-id-export`

.. note:: On import of the :term:`Profile` :term:`Parameters` if a :term:`Parameter` already exists with the same :ref:`parameter-name`, :ref:`parameter-config-file` and :ref:`parameter-value` it will link that to the :term:`Profile` instead of creating it.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: PROFILE:CREATE, PARAMETER:CREATE, PROFILE:READ, PARAMETER:READ
:Response Type:  Object

Request Structure
-----------------

:profile:     The exported :term:`Profile`

	:cdn:         The name of the :ref:`profile-cdn` to which this :term:`Profile` belongs
	:description: The :term:`Profile`'s :ref:`profile-description`
	:name:        The :term:`Profile`'s :ref:`profile-name`
	:type:        The :term:`Profile`'s :ref:`profile-type`

:parameters:  An array of :term:`Parameters` in use by this :term:`Profile`

	:config_file: The :term:`Parameter`'s :ref:`parameter-config-file`
	:name:        :ref:`parameter-name` of the :term:`Parameter`
	:value:       The :term:`Parameter`'s :ref:`parameter-value`

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/profiles/import HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Type: application/json

	{ "profile": {
		"name": "GLOBAL",
		"description": "Global Traffic Ops profile",
		"cdn": "ALL",
		"type": "UNK_PROFILE"
	},
	"parameters": [
		{
			"config_file": "global",
			"name": "tm.instance_name",
			"value": "Traffic Ops CDN"
		},
		{
			"config_file": "global",
			"name": "tm.toolname",
			"value": "Traffic Ops"
		}
	]}

Response Structure
------------------
:cdn:         The name of the :ref:`profile-cdn` to which this :term:`Profile` belongs
:description: The :term:`Profile`'s :ref:`profile-description`
:name:        The :term:`Profile`'s :ref:`profile-name`
:type:        The :term:`Profile`'s :ref:`profile-type`
:id:          The :term:`Profile`'s :ref:`profile-id`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: mzP7DVxFAGhICxqagwDyBDRea7oBZPMAx7NCDeOBVCRqlcCFFe7XL3JP58b80aaVOW/2ZGfg/jpYF70cdDfzQA==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 13 Sep 2019 20:14:42 GMT
	Transfer-Encoding: gzip


	{ "alerts": [
		{
			"level": "success",
			"text": "Profile imported [ Global ] with 2 new and 0 existing parameters"
		}
	],
	"response": {
		"cdn": "ALL",
		"name": "Global",
		"id": 18,
		"type": "UNK_PROFILE",
		"description": "Global Traffic Ops profile"
	}}
