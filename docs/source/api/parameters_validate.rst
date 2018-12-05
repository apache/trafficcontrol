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

.. _to-api-parameters-validate:

***********************
``parameters/validate``
***********************
.. deprecated:: 1.1
	To check for the existence of a parameter with specific name, value etc., use the query parameters of the :ref:`to-api-parameters` endpoint instead.

``POST``
========
Returns a successful response and message if a parameter matching the one in the payload exists, and an error response and message if no such parameter is found.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Object or ``undefined`` - no ``response`` key is returned if the provided parameter could not be matched

Request Structure
-----------------
:name:       Parameter name
:configFile: The *base* filename of the configuration file to which this parameter belongs e.g. "foo" not "/path/to/foo"
:secure:     When ``true``, the parameter value is visible only to "admin"-role users
:value:       The parameter value - if ``secure`` is true and the user does not have the "admin" role this will be obfuscated (at the time of this writing the obfuscation value is defined to be ``"********"``) but **not** missing

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/parameters/validate HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 80
	Content-Type: application/json

	{
		"name": "foo",
		"value": "bar",
		"configFile": "records.config",
		"secure": true
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
	:caption: Response Example - Parameter Found

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 05 Dec 2018 20:35:42 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 06 Dec 2018 00:35:42 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: CcsN9WhMPnvlPtBAcTnecILm1eM1ZxEySwmk3rdCclydPu0cMgefRVI/aRYe+IDAKWFmpeZHg+g1Ed11R7dfWg==
	Content-Length: 149

	{ "alerts": [
		{
			"level": "success",
			"text": "Parameter exists."
		}
	],
	"response": {
		"value": "bar",
		"name": "foo",
		"secure": 0,
		"id": 125,
		"configFile": "records.config"
	}}

.. code-block:: http
	:caption: Response Example - Parameter Not Found

	HTTP/1.1 400 Bad Request
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Wed, 05 Dec 2018 20:42:10 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 06 Dec 2018 00:42:10 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: kUNe70iQz1eEjsSZK3hk3WaJ3eTpBsepdDRUYeXTgEII3lBD5NiXobShT6zGhWJTsalHbNegjWbfAWsly/XEQQ==
	Content-Length: 116

	{ "alerts": [
		{
			"level": "error",
			"text": "parameter [name:fooa, config_file:records.config, value:bar] does not exist."
		}
	]}

.. note:: This endpoint returns a client-side error response when the parameter was not found - as such any API tools that wish to use this endpoint should be aware that a client-side error response code may not actually mean that an error occurred. However, neither can it be said that a ``400`` response code means that the parameter wasn't found; that response code is also returned in the event of _true_ client-side errors e.g. a malformed JSON payload in the request.
