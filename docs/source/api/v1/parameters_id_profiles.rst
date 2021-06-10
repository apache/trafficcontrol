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

.. _to-api-v1-parameters-id-profiles:

******************************
``parameters/{{ID}}/profiles``
******************************
.. deprecated:: ATCv4

``GET``
=======
Retrieves all :term:`Profiles` assigned to a specific :term:`Parameter`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------+
	| Name |                    Description                                                              |
	+======+=============================================================================================+
	|  ID  | The :ref:`parameter-id` of the :term:`Parameter` for which :term:`Profiles` shall be listed |
	+------+---------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Structure

	GET /api/1.4/parameters/4/profiles HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
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
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 05 Dec 2018 20:51:23 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: y5fA9q1VogDGxL66ka+ofTtLo3JiTj+Bdrvc4DnfrjFyzqll+537WySFj1nE0C29Twx5l/C8JEHy3Byaz/wbfA==
	Content-Length: 184

	{ "response": [
		{
			"routingDisabled": 0,
			"lastUpdated": "2018-12-05 17:50:49.007102+00",
			"name": "GLOBAL",
			"type": "UNK_PROFILE",
			"id": 1,
			"description": "Global Traffic Ops profile, DO NOT DELETE"
		}
	],
	"alerts": [
		{
			"level": "warning",
			"text": "This endpoint is deprecated, please use 'GET /profiles' instead"
		}
	]}
