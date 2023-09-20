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

How to Read this Documentation
==============================
Each endpoint for each version is on its own page, titled with the request path. The request paths shown on each endpoint's page are - unless otherwise noted - only usable by being appended to the request path prefix ``/api/<version>/`` where ``<version>`` is the API version being requested. The API versions officially supported as of the time of this writing are 3.0, 3.1, 4.0, 4.1, and 5.0. All endpoints are documented as though they were being used in version 3.1 in the version 3 documentation, version 4.1 in the version 4 documentation, and version 5.0 in the version 5 documentation. If an endpoint or request method of an endpoint is only available after a specific version, that will be noted next to the method or endpoint name. If changes were made to the structure of an endpoint's input or output, the version number and nature of the change will be noted.

Every endpoint is documented with a section for each method, containing the subsections "Request Structure" and "Response Structure" which identify all properties and structure of the Request to and Response from the endpoint. Before these subsections, three key pieces of information will be provided:

Auth. Required
	This will either be 'Yes' to indicate that a user must be authenticated (or "logged-in") via e.g. :ref:`to-api-user-login` to use this method of the endpoint, or 'No' to indicate that this is not required.
Roles Required
	.. deprecated:: ATCv7.0
		Roles for the use of authentication/permissions have been deprecated in favor of role based permissions, see :pr:`5848`.

	Any permissions roles that are allowed to use this method of the endpoint will be listed here. Users with roles not listed here will be unable to properly use these endpoints.
Permissions Required
	Any permissions that are needed to use this endpoint. Users with roles that don't have the permissions will be unable to properly use these endpoints.
Response Type
	Unless otherwise noted, all responses are JSON objects. See `Response Structure`_ for more information.

The methods of endpoints that require/accept data payloads - unless otherwise noted - always interpret the content of the payload as a JSON object, regardless of the request's ``Content-Type`` header. Because of this, all payloads are - unless otherwise noted - JSON objects. The Request Structure and Response Structure subsections will contain explanations of the fields before any examples like e.g.

:foo: A constant field that always contains "foo"
:bar: An array of objects that each represent a "bar" object

	:name:  The bar's name
	:value: The bar's value (an integer)

All fields are mandatory in a request payload, or always present in a response payload unless otherwise noted in the field description.

In most cases, JSON objects have been "pretty-printed" by inserting line breaks and indentation. This means that the ``Content-Length`` HTTP header does not, in general, accurately portray the length of the content displayed in Request Examples and Response Examples. Also, the Traffic Ops endpoints will ignore any content negotiation, meaning that the ``Content-Type`` header of a request is totally meaningless. A utility may choose to pass the data as e.g. :mimetype:`application/x-www-form-urlencoded` (cURL's default ``Content-Type``) when constructing a Request Example, but the example itself will most often show :mimetype:`application/json` in order for syntax highlighting to properly work.

.. _to-api-response-structure:

Response Structure
------------------
Unless otherwise noted, all response payloads come as JSON objects.

.. code-block:: json
	:caption: Response Structure

	{
		"response": "<JSON object with main response>",
	}

To make the documentation easier to read, only the ``<JSON object with main response>`` is documented, even though the response endpoints may return other top-level objects (most commonly the ``"alerts"`` object). The field definitions listed in the Response Structure subsection of an endpoint method are the elements of this object. Sometimes the ``response`` object is a string, sometimes it's an object that maps keys to values, sometimes it's an array that contains many arbitrary objects, and sometimes it isn't present at all. For ease of reading, the field lists delegate the distinction to be made by the ``Response Type`` field directly under the request method heading.

Response Type Meanings
""""""""""""""""""""""
Array
	The fields in the field list refer to the keys of the objects in the ``response`` array.
Object
	The fields in the field list refer to the keys of the ``response`` object.
``undefined``
	No ``response`` object is present in the response payload. Unless the format is otherwise noted, this means that there should be no field list in the "Response Structure" subsection.

Summary
-------
The top-level ``summary`` object is used to provide summary statistics about object collections. In general the use of ``summary`` is left to be defined by API endpoints (subject to some restrictions). When an endpoint uses the ``summary`` object, its fields will be explained in a subsection of the "Response Structure" named "Summary Fields".

The following, reserved properties of ``summary`` are guaranteed to always have their herein-described meaning.

.. _reserved-summary-fields:

``count``
	``count`` contains an unsigned integer that defines the total number of results that could possibly be returned given the non-pagination query parameters supplied by the client.

.. _non-rfc-datetime:

Traffic Ops's Custom Date/Time Format
-------------------------------------
Traffic Ops will often return responses from its API that include dates. As of the time of this writing, the vast majority of those dates are written in a non-:RFC:`3339` format (this is tracked by :issue:`5911`). This is most commonly the case in the ``last_updated`` properties of objects returned as JSON-encoded documents. The format used is :samp:`{YYYY}-{MM}-{DD} {hh}:{mm}:{ss}±{ZZ}`, where ``YYYY`` is the 4-digit year, ``MM`` is the two-digit (zero padded) month, ``DD`` is the two-digit (zero padded) day of the month, ``hh`` is the two-digit (zero padded) hour of the day, ``mm`` is the two-digit (zero padded) minute of the hour, ``ss`` is the two-digit (zero padded) second of the minute, and ``ZZ`` is the two-digit (zero padded) timezone offset in hours of the date/time's local timezone from UTC (the offset can be positive or negative as indicated by a ``+`` or a ``-`` directly before it, where the sample has a ``±``).

.. note:: In practice, all Traffic Ops API responses use the UTC timezone (offset ``+00``), but do note that this custom format is not capable of representing all timezones.

.. code-block:: text
	:caption: Example Date/Timestamp

	2021-06-07 08:01:02+00

Using API Endpoints
===================
#. Authenticate with valid Traffic Control user account credentials (the same used by Traffic Portal).
#. Upon successful user authentication, note the Mojolicious cookie value in the response headers\ [1]_.

	.. note:: Many tools have methods for doing this without manual intervention - a web browser for instance will automatically remember and properly handle cookies. Another common tool, cURL, has command line switches that will also accomplish this. Most high-level programming language libraries will implement a cookie-handling method as well.

#. Pass the Mojolicious cookie value, along with any subsequent calls to an authenticated API endpoint.

.. note:: Although many endpoints in API version 1.x supported a ``.json`` suffix, API version 2.x does not support it at all. Even when using API version 1.x using the ``.json`` suffix should be avoided at all costs, because there's no real consistency regarding when it may be used, and the output of API endpoints, in general, are not capable of representing POSIX-compliant files (as a 'file extension' might imply).

Example Session
---------------
A user makes a request to the ``/api/4.0/asns`` endpoint.

.. code-block:: http

	GET /api/4.0/asns HTTP/1.1
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

To authenticate, the user sends a POST request containing their login information to the ``/api/4.0/user/login`` endpoint.

.. code-block:: http

	POST /api/4.0/user/login HTTP/1.1
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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

Using this cookie, the user can now access their original target - the ``/api/4.0/asns`` endpoint...

.. code-block:: http

	GET /api/4.0/asns HTTP/1.1
	Accept: application/json
	Cookie: mojolicious=...;
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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

:level: ``"success"``, ``"info"``, ``"warning"`` or ``"error"`` as appropriate
:text: The alert's actual message

The most common errors returned by Traffic Ops are:

401 Unauthorized
	When a "mojolicious" cookie is supplied that is invalid or expired, or the login credentials are incorrect the server responds with a ``401 UNAUTHORIZED`` response code.

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
		X-Server-Name: traffic_ops_golang/
		Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
		Vary: Accept-Encoding
		Whole-Content-Sha512: Ff5hO8ZUNUMbwCW0mBuUlsvrSmm/Giijpq7O3uLivLZ6VOu6eGom4Jag6UqlBbbDBnP6AG7l1Szdt74TT6NidA==
		Transfer-Encoding: chunked

		{ "alerts": [
			{
				"level": "error",
				"text": "Resource not found."
			}
		]}


500 Internal Server Error
	When a server-side error occurs, the API will return a ``500 INTERNAL SERVER ERROR`` response.

	.. code-block:: http
		:caption: Example Response to ``GET /api/4.0/servers``. (when a server error such as a postgres failure occured)

		HTTP/1.1 500 Internal Server Error
		Access-Control-Allow-Credentials: true
		Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
		Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
		Cache-Control: no-cache, no-store, max-age=0, must-revalidate
		Content-Length: 93
		Content-Type: application/json
		Date: Tue, 02 Oct 2018 17:29:42 GMT
		X-Server-Name: traffic_ops_golang/
		Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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
TrafficOps client libraries are available in Java, Go and Python. You can read (very little) more about them in the client README at :atc-file:`traffic_control/clients`.

API V3 Routes
=============
API routes available in version 3.

.. deprecated:: ATCv7
	Traffic Ops API version 3 is deprecated in favor of version 4.

.. toctree::
	:maxdepth: 4
	:glob:

	v3/*

API V4 Routes
=============
API routes available in version 4.

.. toctree::
	:maxdepth: 4
	:glob:

	v4/*

API V5 Routes
=============
API routes available in version 5.

.. toctree::
	:maxdepth: 4
	:glob:

	v5/*


.. [1] A cookie obtained by logging in through Traffic Portal can be used to access API endpoints under the Traffic Portal domain name - since it will proxy such requests back to Traffic Ops. This is not recommended in actual deployments, however, because it will involve an extra network connection which could be avoided by simply using the Traffic Ops domain itself.
