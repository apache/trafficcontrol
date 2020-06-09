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

.. _to-api-v1-servers-server-configfiles-ats-filename:

***************************************************
``servers/{{server}}/configfiles/ats/{{filename}}``
***************************************************

.. seealso:: The :ref:`to-api-v1-servers-server-configfiles-ats` endpoint

.. deprecated:: 1.4
	Using the API to retrieve generated configuration files for servers is deprecated, and unavailable in more recent API versions. Also, in ATC version 4.x and higher, it is not guaranteed that configuration files will be output correctly, or even successfully. Instead, configuration file generation is now the responsibility of :ref:`atstccfg`.

``GET``
=======
Returns the requested configuration file for download.

:Auth. Required: Yes
:Roles Required: "operations"
:Response Type:  **NOT PRESENT** - endpoint returns custom :mimetype:`text/plain` response (represents the contents of the requested configuration file)

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+-------------------+--------------------------------------------------------------+
	| Parameter | Type              | Description                                                  |
	+===========+===================+==============================================================+
	| server    | string or integer | Either the name or integral, unique, identifier of a server  |
	+-----------+-------------------+--------------------------------------------------------------+
	| filename  | string            | The name of a configuration file used by ``server``          |
	+-----------+-------------------+--------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.2/servers/edge/configfiles/ats/hosting.config HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
.. note:: If the file identified by ``filename`` does exist, but is configured at a higher level than "server", a JSON response will be returned and the ``alerts`` array will contain a ``"level": "error"`` node which identifies the correct scope of the configuration file.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: text/plain;charset=UTF-8
	Date: Thu, 15 Nov 2018 15:32:25 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: EmhHogPfcxQq2zHmFFJtjwzZiUHNgOZvE572Se/H/54gwarkkKjm89+xJr7fQbfytc7xWYApzwfjNl6LfbM0hg==
	Content-Length: 107

	# DO NOT EDIT - Generated for edge by Traffic Ops (trafficops.infra.ciab.test:443) on Thu Nov 15 15:32:25 UTC 2018
	hostname=*   volume=1
