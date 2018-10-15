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

.. _to-api:

***************
Traffic Ops API
***************
The Traffic Ops API provides programmatic access to read and write Traffic Control data which allows for the monitoring of CDN performance and configuration of Traffic Control settings and parameters.

API Routes
==========
.. table:: All Implemented API routes

	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| 1.0                                |   1.1                                              |   1.2                                              |   1.3                                              |
	+====================================+====================================================+====================================================+====================================================+
	| /asns                              |   :ref:`to-api-asns`                               |   :ref:`to-api-asns`                               |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /availableds                       |   :ref:`to-api-v11-ds-route`                       |   :ref:`to-api-v12-ds-route`                       |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| *Not Implemented*                  |   *Not Implemented*                                |   *Not Implemented*                                |   :ref:`to-api-v13-coordinates-route`              |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datacrans                         |   /api/1.1/crans.json                              |   /api/1.2/crans.json                              |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datacrans/orderby/:field          |   /api/1.1/crans.json                              |   /api/1.2/crans.json                              |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datadeliveryservice               |   :ref:`to-api-v11-ds-route`                       |   :ref:`to-api-v12-ds-route`                       |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datadeliveryserviceserver         |   /api/1.1/deliveryserviceserver.json              |   /api/1.2/deliveryserviceserver.json              |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datadomains                       |   /api/1.1/cdns/domains.json                       |   /api/1.2/cdns/domains.json                       |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| *Not Implemented*                  |  *Not Implemented*                                 |   :ref:`to-api-v12-ds-stats-route`                 |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datahwinfo                        |   :ref:`to-api-v11-hwinfo-route`                   |   :ref:`to-api-v12-hwinfo-route`                   |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datalinks                         |   /api/1.1/deliveryserviceserver.json              |   /api/1.2/deliveryserviceserver.json              |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datalinks/orderby/:field          |   /api/1.1/deliveryserviceserver.json              |   /api/1.2/deliveryserviceserver.json              |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datalogs                          |   :ref:`to-api-v11-change-logs-route`              |   :ref:`to-api-v12-change-logs-route`              |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataparameter                     |   :ref:`to-api-v11-parameters-route`               |   :ref:`to-api-v12-parameters-route`               |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataparameter/:parameter          |   /api/1.1/parameters/profile/:parameter.json      |   /api/1.2/parameters/profile/:parameter.json      |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataphys_location                 |   :ref:`to-api-v11-phys-loc-route`                 |   :ref:`to-api-v12-phys-loc-route`                 |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataprofile                       |   :ref:`to-api-v11-profiles-route`                 |   :ref:`to-api-v12-profiles-route`                 |   *Not Implemented*                                |
	|                                    |                                                    |                                                    |                                                    |
	| /dataprofile/orderby/name          |                                                    |                                                    |                                                    |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataregion                        |   :ref:`to-api-v11-regions-route`                  |   :ref:`to-api-v12-regions-route`                  |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datarole                          |   :ref:`to-api-v11-roles-route`                    |   :ref:`to-api-v12-roles-route`                    |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datarole/orderby/:field           |   :ref:`to-api-v11-roles-route`                    |   :ref:`to-api-v12-roles-route`                    |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataserver                        |   :ref:`to-api-v11-servers-route`                  |   :ref:`to-api-v12-servers-route`                  |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataserver/orderby/:field         |   :ref:`to-api-v11-servers-route`                  |   :ref:`to-api-v12-servers-route`                  |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /dataserverdetail/select/:hostname |   /api/1.1/servers/hostname/:hostname/details.json |   /api/1.2/servers/hostname/:hostname/details.json |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datastaticdnsentry                |   :ref:`to-api-v11-static-dns-route`               |   :ref:`to-api-v12-static-dns-route`               |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datastatus                        |   :ref:`to-api-v11-statuses-route`                 |   :ref:`to-api-v12-statuses-route`                 |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datastatus/orderby/name           |   :ref:`to-api-v11-statuses-route`                 |   :ref:`to-api-v12-statuses-route`                 |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datatype                          |   :ref:`to-api-v11-types-route`                    |   :ref:`to-api-v12-types-route`                    |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datatype/orderby/:field           |   :ref:`to-api-v11-types-route`                    |   :ref:`to-api-v12-types-route`                    |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datauser                          |   :ref:`to-api-v11-users-route`                    |   :ref:`to-api-v12-users-route`                    |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| /datauser/orderby/:field           |   :ref:`to-api-v11-users-route`                    |   :ref:`to-api-v12-users-route`                    |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| *Not Implemented*                  |   *Not Implemented*                                |   :ref:`to-api-v12-configfiles_ats-route`          |   *Not Implemented*                                |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| *Not Implemented*                  |   *Not Implemented*                                |   *Not Implemented*                                |   :ref:`to-api-v13-static-dns-entry-route`         |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+
	| *Not Implemented*                  |   *Not Implemented*                                |   *Not Implemented*                                |   :ref:`to-api-v13-origin-route`                   |
	+------------------------------------+----------------------------------------------------+----------------------------------------------------+----------------------------------------------------+


.. toctree::
	:maxdepth: 4
	:hidden:

	api_capabilities
	api_capabilities_id
	asns
	asns_id
	cache_stats
	caches_stats
	cachegroup_parameterID_parameter
	cachegroupparameters
	cachegroups
	cachegroups_id
	cachegroups_trimmed
	v11/index
	v12/index
	v13/index
	v14/index



Response Structure
==================
All successful responses have the following structure:

.. code-block:: json

	{
		"response": "<JSON object with main response>",
	}

To make the documentation easier to read, only the ``<JSON object with main response>`` is documented, even though the response endpoints may return other top-level objects (most commonly the ``"alerts"`` object.

Using API Endpoints
===================
#. Authenticate with valid Traffic Control user account credentials (the same used by Traffic Portal).
#. Upon successful user authentication, note the Mojolicious cookie value in the response headers\ [1]_.

	.. note:: Many tools have methods for doing this without manual intervention - a web browser for instance will automatically remember and properly handle cookies. Another common tool, cURL, has command line switches that will also accomplish this. Most high-level programming language libraries will implement a cookie-handling method as well.

#. Pass the Mojolicious cookie value, along with any subsequent calls to an authenticated API endpoint.

Example Session
---------------
A user makes a request to the ``/api/1.1/asns`` endpoint.

.. code-block:: http

	GET /api/1.1/asns HTTP/1.1
	Accept: application/json
	Host: trafficops.infra.ciab.test
	User-Agent: example

The response JSON indicates an authentication error.

.. code-block:: http

	HTTP/1.1 401 UNAUTHORIZED
	Content-Length: 68
	Content-Type: application/json
	Date: Tue, 02 Oct 2018 13:12:30 GMT

	{ "alerts": [
		{
			"level":"error",
			"text":"Unauthorized, please log in."
		}
	]}

To authenticate, the user sends a POST request containing their login information to the ``/api/1.3/user/login`` endpoint.

.. code-block:: http

	POST /api/1.1/user/login HTTP/1.1
	User-Agent: example
	Host: trafficops.infra.ciab.test
	Accept: application/json
	Content-Length: 32
	Content-Type: application/x-www-form-urlencoded

Traffic Ops responds with a Mojolicious cookie to be used for future requests, and a message indicating the success or failure (in this case success) of the login operation.

.. code-block:: http

	HTTP/1.1 200 OK
	Connection: keep-alive
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Set-Cookie: mojolicious=eyJhdXRoX2RhdGEiOiJhZG1pbiIsImV4cGlyZXMiOjE1Mzg1MDY5OTgsImJ5IjoidHJhZmZpY2NvbnRyb2wtZ28tdG9jb29raWUifQ--bcc9aade79b6de436cb4962ef5cec397f7ac5bd2; Path=/; Expires=Tue, 02 Oct 2018 19:03:18 GMT; HttpOnly
	Content-Type: application/json
	Date: Tue, 02 Oct 2018 12:53:32 GMT
	Access-Control-Allow-Credentials: true
	Content-Length: 81
	X-Server-Name: traffic_ops_golang/

	{ "alerts": [
		{
			"level": "success",
			"text": "Successfully logged in."
		}
	]}

Using this cookie, the user can now access their original target - the ``/api/1.1/asns`` endpoint...

.. code-block:: http

	GET /api/1.1/asns HTTP/1.1
	Accept: application/json
	Cookie: mojolicious=eyJleHBpcmVzIjoxNDI5NDAyMjAxLCJhdXRoX2RhdGEiOiJhZG1pbiJ9--f990d03b7180b1ece97c3bb5ca69803cd6a79862;
	Host: trafficops.infra.ciab.test
	User-Agent: Example

\... and the Traffic Ops server will now happily service this request.

.. code-block:: http

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Connection: keep-alive
	Content-Encoding: gzip
	Content-Length: 48
	Content-Type: application/json
	Date: Tue, 02 Oct 2018 12:55:57 GMT
	Set-Cookie: mojolicious=eyJhdXRoX2RhdGEiOi…ccd4eae46c6; Path=/; HttpOnly
	Whole-Content-SHA512: u+Q5X7z/DMTc/VzRGaFlJBA8btA8EC…dnA85HCYTm8vVwsQCvle+uVc1nA==
	X-Server-Name: traffic_ops_golang/

	{ "response": {
		"asns": [
			{
				"lastUpdated": "2012-09-17 21:41:22",
				"id": 27,
				"asn": 7015,
				"cachegroup": "us-ma-woburn",
				"cachegroupId": 2
			},
			{
				"lastUpdated": "2012-09-17 21:41:22",
				"id": 28,
				"asn": 7016,
				"cachegroup": "us-pa-pittsburgh",
				"cachegroupID": 3
			}
		]
	}}

API Errors
==========
If an API endpoint has something to say besides the actual response (usually an error message), it will add a top-level object to the response JSON with the key ``"alerts"``. This will be an array of objects that represent messages from the server, each with the following string fields:

:``level``: ``"success"``, ``"info"``, ``"warning"`` or ``"error"`` as appropriate
:``text``: The alert's actual message

The most common errors returned by Traffic Ops are:

401 Unauthorized
	When a Mojolicious cookie is supplied that is invalid or expired, or the login credentials are incorrect the server responds with a ``401 UNAUTHORIZED`` response code.

	.. code-block:: http
		:caption: Example of a Response to a Login Request with Bad Credentials

		HTTP/1.1 401 Unauthorized
		Access-Control-Allow-Credentials: true
		Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
		Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
		Access-Control-Allow-Origin: *
		Content-Type: application/json
		Whole-Content-Sha512: xRKu2Q7Yj07UA6A6SyxMNmcBpuBcW2/bzuKO5eTZ2y4V27rXfP/5bSkNPesomJbiOO+xSmiybDsHlcL3P+pzpg==
		X-Server-Name: traffic_ops_golang/
		Date: Tue, 02 Oct 2018 13:28:30 GMT
		Content-Length: 69

		{ "alerts": [
			{
				"text": "Invalid username or password.",
				"level": "error"
			}
		]}

404 Not Found
	When the requested resource (path) doesn't exist, Traffic Ops returns a ``404 NOT FOUND`` response code.

	.. code-block:: http
		:caption: Example Response to ``GET /not/an/api/path HTTP/1.1`` with Proper Cookies

		HTTP/1.1 404 Not Found
		Access-Control-Allow-Credentials: true
		Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
		Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
		Access-Control-Allow-Origin: *
		Cache-Control: no-cache, no-store, max-age=0, must-revalidate
		Content-Type: text/html;charset=UTF-8
		Date: Tue, 02 Oct 2018 13:58:56 GMT
		Server: Mojolicious (Perl)
		Set-Cookie: mojolicious=eyJhdXRoX2RhdGEiOiJhZG1pbiIsImV4cGlyZXMiOjE1Mzg1MDMxMzYsImJ5IjoidHJhZmZpY2NvbnRyb2wtZ28tdG9jb29raWUifQ----9b144306f8bb6020eadb950647b3dc0eebeb7eae; expires=Tue, 02 Oct 2018 17:58:56 GMT; path=/; HttpOnly
		Vary: Accept-Encoding
		Whole-Content-Sha512: Ff5hO8ZUNUMbwCW0mBuUlsvrSmm/Giijpq7O3uLivLZ6VOu6eGom4Jag6UqlBbbDBnP6AG7l1Szdt74TT6NidA==
		Transfer-Encoding: chunked

		The content of this response will be the Legacy UI login page (which is omitted because it's huge)


500 Internal Server Error
	When a server-side error occurs, the Perl API will return a ``500 INTERNAL SERVER ERROR`` response (the below example request will result in a ``400 BAD REQUEST`` response if using the v1.3 API instead - as this will use the Go server's API)

	.. code-block:: http
		:caption: Example Response to ``GET /api/1.1/servers/hostname/jj/details.json`` ('jj' doesn't exist)

		HTTP/1.1 500 Internal Server Error
		Access-Control-Allow-Credentials: true
		Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
		Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
		Cache-Control: no-cache, no-store, max-age=0, must-revalidate
		Content-Length: 93
		Content-Type: application/json
		Date: Tue, 02 Oct 2018 17:29:42 GMT
		Server: Mojolicious (Perl)
		Set-Cookie: mojolicious=eyJhdXRoX2RhdGEiOiJhZG1pbiIsImV4cGlyZXMiOjE0Mjk0MDQzMDZ9--1b08977e91f8f68b0ff5d5e5f6481c76ddfd0853; expires=Sun, 19 Apr 2015 00:45:06 GMT; path=/; HttpOnly
		Vary: Accept-Encoding
		Whole-Content-Sha512: gFa4NYFmofCbV7YqgwyFRzKk90+KNgoZu6p2Nx98J4Gy7/2j55tYknvk53WXuMdMKKrgYMop4uiYOla1k1ozQQ==

		{ "alerts": [
			{
				"level": "error",
				"text": "An error occurred. Please contact your administrator."
			}
		]}

The rest of the API documentation will only document the ``200 OK`` case, where no errors have occurred.

TrafficOps Native Client Libraries
==================================
TrafficOps client libraries are available in Java, Go and Python. You can read (very little) more about them in `the client README <https://github.com/apache/trafficcontrol/tree/master/traffic_control/clients>`_.

.. [1] A cookie obtained by logging in through Traffic Portal can be used to access API endpoints under the Traffic Portal domain name - since it will proxy such requests back to Traffic Ops. This is not recommended in actual deployments, however, because it will involve an extra network connection which could be avoided by simply using the Traffic Ops domain itself.
