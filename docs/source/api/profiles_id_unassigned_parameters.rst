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

.. _to-api-profiles-id-unassigned_parameters:

*****************************************
``profiles/{{ID}}/unassigned_parameters``
*****************************************
.. warning:: There are **very** few good reasons to use this endpoint - be sure not limit said use.

``GET``
=======
Retrieves all parameters NOT assigned to the specified profile.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------------------+
	| Name | Description                                                                                   |
	+======+===============================================================================================+
	|  ID  | The integral, unique identifier of the profile for which unassigned parameters will be listed |
	+------+-----------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/profiles/9/unassigned_parameters HTTP/1.1
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
	Whole-Content-Sha512: iO7YHU+0spCPSaR6oDrVIQwxSS1GoSyi8K6ng4eemuxqOxB9FdfPgBpXN8w+xmxf2ZwRMLXHv5S6cfIoNNDnqw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 05 Dec 2018 21:37:50 GMT
	Transfer-Encoding: chunked

	{ "response": [
		{
			"configFile": "parent.config",
			"id": 1,
			"lastUpdated": "2018-12-05 17:50:47+00",
			"name": "mso.parent_retry",
			"secure": false,
			"value": "simple_retry"
		},
		{
			"configFile": "parent.config",
			"id": 2,
			"lastUpdated": "2018-12-05 17:50:47+00",
			"name": "mso.parent_retry",
			"secure": false,
			"value": "unavailable_server_retry"
		}
	]

.. note:: The response example for this endpoint has been truncated to only the first two elements of the resulting array, as the output was hundreds of lines long.
