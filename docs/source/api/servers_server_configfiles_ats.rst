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

.. _to-api-servers-server-configfiles-ats:

**************************************
``servers/{{server}}/configfiles/ats``
**************************************

.. seealso:: The :ref:`to-api-servers-server-configfiles-ats-filename`, :ref:`to-api-cdns-cdn-configfiles-ats-filename`, and :ref:`to-api-profiles-profile-configfiles-ats-filename` endpoints.

``GET``
=======
Gets a list of the configuration files used by ``server``

:Auth. Required: Yes
:Roles Required: "operations"
:Response Type:  **NOT PRESENT** - endpoint returns custom application/json response

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+-------------------+--------------------------------------------------------------+
	| Parameter | Type              | Description                                                  |
	+===========+===================+==============================================================+
	| server    | string or integer | Either the name or integral, unique, identifier of a server  |
	+-----------+-------------------+--------------------------------------------------------------+

Response Structure
------------------
:info: An object that provides information about ``server`` as it is understood by Traffic Ops

	:cdnId:         The integral, unique, identifier of the CDN to which ``server`` is assigned
	:cdnName:       The name of the CDN to which ``server`` is assigned
	:profileName:   The :ref:`profile-name` of the :term:`Profile` used by this server
	:profileId:     The :ref:`profile-id` the :term:`Profile` used by this server
	:serverId:      An integral, unique, identifier for ``server``
	:serverIpv4:    IPv4 address of the server
	:serverTcpPort: The port number on which ``server`` listens for incoming TCP connections
	:toRevProxyUrl: An optional field which, if present, gives a URL that resolves to a proxy for Traffic Ops which ``server`` ought to use rather than directly contacting ``toUrl``
	:toUrl:         A full URL that resolves to the Traffic Ops instance

:configFiles: An array of objects which each represent a configuration file used by the server

	:apiUri:      An optional field which, if present, gives a path relative to the Traffic Ops instance (or reverse proxy when applicable) URL where the actual file's contents may be retrieved\ [1]_
	:fnameOnDisk: The filename of the configuration file as stored on the server
	:location:    The directory location of the configuration file as stored on the server
	:scope:       The "scope" of the configuration file, which will be one of:

		"cdns"
			The file is used by all caches in the CDN
		"profiles"
			The file is used by all servers with the same :term:`Profile`
		"servers"
			The most specific grouping of servers which use this file is simply a collection of distinct servers

	:url:         An optional field which, if present, gives the full URL used to retrieve the actual file's contents\ [1]_

.. versionchanged:: Traffic Control 2.0
	Elements of the ``"configFile"`` array may no longer have the ``"contents"`` key - all file contents are now retrieved via a network request

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: text/plain;charset=UTF-8
	Date: Thu, 15 Nov 2018 15:28:10 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Thu, 15 Nov 2018 19:28:10 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: K6pRI4MkN8O9+wKW8MG3w6nTnmLHtCZKqzXCjw4JfoMYIVJC6fVTN9ysGML71VF2T7ZAIP1TveWhjaH/fNr7sQ==
	Transfer-Encoding: chunked

	{ "info": {
		"profileId": 9,
		"toUrl": null,
		"serverIpv4": "172.16.239.100",
		"serverTcpPort": 80,
		"serverName": "edge",
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"serverId": 10,
		"profileName": "ATS_EDGE_TIER_CACHE"
	},
	"configFiles": [
		{
			"fnameOnDisk": "astats.config",
			"location": "/etc/trafficserver",
			"apiUri": "/api/1.2/profiles/ATS_EDGE_TIER_CACHE/configfiles/ats/astats.config",
			"scope": "profiles"
		},
		{
			"fnameOnDisk": "cache.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/profiles/ATS_EDGE_TIER_CACHE/configfiles/ats/cache.config",
			"scope": "profiles"
		},
		{
			"fnameOnDisk": "cacheurl_foo.config",
			"location": "/etc/trafficserver",
			"apiUri": "/api/1.2/cdns/CDN-in-a-Box/configfiles/ats/cacheurl_foo.config",
			"scope": "cdns"
		},
		{
			"fnameOnDisk": "hdr_rw_foo.config",
			"location": "/etc/trafficserver",
			"apiUri": "/api/1.2/cdns/CDN-in-a-Box/configfiles/ats/hdr_rw_foo.config",
			"scope": "cdns"
		},
		{
			"fnameOnDisk": "hosting.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/servers/edge/configfiles/ats/hosting.config",
			"scope": "servers"
		},
		{
			"fnameOnDisk": "ip_allow.config",
			"location": "/etc/trafficserver",
			"apiUri": "/api/1.2/servers/edge/configfiles/ats/ip_allow.config",
			"scope": "servers"
		},
		{
			"fnameOnDisk": "parent.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/servers/edge/configfiles/ats/parent.config",
			"scope": "servers"
		},
		{
			"fnameOnDisk": "plugin.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/profiles/ATS_EDGE_TIER_CACHE/configfiles/ats/plugin.config",
			"scope": "profiles"
		},
		{
			"fnameOnDisk": "records.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/profiles/ATS_EDGE_TIER_CACHE/configfiles/ats/records.config",
			"scope": "profiles"
		},
		{
			"fnameOnDisk": "regex_remap_foo.config",
			"location": "/etc/trafficserver",
			"apiUri": "/api/1.2/cdns/CDN-in-a-Box/configfiles/ats/regex_remap_foo.config",
			"scope": "cdns"
		},
		{
			"fnameOnDisk": "regex_revalidate.config",
			"location": "/etc/trafficserver",
			"apiUri": "/api/1.2/cdns/CDN-in-a-Box/configfiles/ats/regex_revalidate.config",
			"scope": "cdns"
		},
		{
			"fnameOnDisk": "remap.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/servers/edge/configfiles/ats/remap.config",
			"scope": "servers"
		},
		{
			"fnameOnDisk": "storage.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/profiles/ATS_EDGE_TIER_CACHE/configfiles/ats/storage.config",
			"scope": "profiles"
		},
		{
			"fnameOnDisk": "volume.config",
			"location": "/etc/trafficserver/",
			"apiUri": "/api/1.2/profiles/ATS_EDGE_TIER_CACHE/configfiles/ats/volume.config",
			"scope": "profiles"
		}
	]}

.. note:: Some DSCP-related files like e.g. ``set_dscp_0.config`` have been removed from this response, which otherwise reflects a stock CDN-in-a-Box configuration. This was done both for brevity's sake, and due to the expectation that these will disappear from the default configuration in the (hopefully near) future.

.. [1] Exactly one of these fields is guaranteed to exist for any given configuration file - although "apiUrl" is far more common.
