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

.. _to-api-profiles-profile-configfiles-ats-filename:

*****************************************************
``profiles/{{profile}}/configfiles/ats/{{filename}}``
*****************************************************

.. seealso:: The :ref:`to-api-servers-server-configfiles-ats` endpoint

``GET``
=======
Returns the requested configuration file for download.

:Auth. Required: Yes
:Roles Required: "operations"
:Response Type:  **NOT PRESENT** - endpoint returns custom text/plain response (represents the contents of the requested configuration file)

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+-------------------+--------------------------------------------------------------------------+
	| Parameter | Type              | Description                                                              |
	+===========+===================+==========================================================================+
	| profile   | string or integer | Either the :ref:`profile-name` or :ref:`profile-id` of a :term:`Profile` |
	+-----------+-------------------+--------------------------------------------------------------------------+
	| filename  | string            | The name of a configuration file used by ``profile``                     |
	+-----------+-------------------+--------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/profiles/ATS_MID_TIER_CACHE/configfiles/ats/volume.config HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. note:: If the file identified by ``filename`` doesn't exist at the :term:`Profile`, a JSON response will be returned and the ``alerts`` array will contain a ``"level": "error"`` node which suggests other scopes to check for the configuration file.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: text/plain;charset=UTF-8
	Date: Thu, 15 Nov 2018 15:23:44 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 15 Nov 2018 19:23:44 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: C1Hrs4y3qSThOeZJo5aDJu1QjD/r/7vO6c7E7TaFXx67kWat91uk9BSvieXN5yrOE4HkGsiGBkNZjjZ3hb5mYw==
	Content-Length: 211

	# DO NOT EDIT - Generated for ATS_MID_TIER_CACHE by Traffic Ops (trafficops.infra.ciab.test:443) on Thu Nov 15 15:23:44 UTC 2018
	# TRAFFIC OPS NOTE: This is running with forced volumes - the size is irrelevant
	volume=1 scheme=http size=100%
