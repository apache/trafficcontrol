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

***************************************
Traffic Ops - Migrating from 2.0 to 2.2
***************************************

Per-DeliveryService Routing Names
---------------------------------
Before this release, DNS Delivery Services were hard-coded to use the name ``edge``, so that URLs would appear as e.g. ``edge.myds.mycdn.com``, and HTTP Delivery Services use the name ``tr`` [1]_, e.g. ``tr.myds.mycdn.com``. As of Traffic Control version 2.2, DNS routing names will default to ``cdn`` if left unspecified and can be set to any valid hostname, provided it does not contain the ``.`` character.

Prior to Traffic Control 2.2, the HTTP routing name was configurable via the ``http.routing.name`` option in in the Traffic Router ``http.properties`` configuration file. If your CDN uses that option to change the name from ``tr`` to something else, then you will need to perform the following steps for each CDN affected:

#. In Traffic Portal (if possible, else use the legacy Traffic Ops UI), create the following profile parameter ('Configure' -> 'Parameters' -> '+'). Be sure to double-check for typos, trailing spaces, etc.

	:name: upgrade_http_routing_name
	:config_file: temp
	:value: Whatever value is used for the affected CDN's ``http.routing.name``

#. Add this parameter to a single profile in the affected CDN

With those profile parameters in place, Traffic Ops can be safely upgraded to 2.2. Before taking a post-upgrade snapshot, make sure to check your Delivery Service example URLs for unexpected routing name changes. Once Traffic Ops has been upgraded to 2.2 and a post-upgrade snapshot has been taken, your Traffic Routers can be upgraded to 2.2 (Traffic Routers must be upgraded after Traffic Ops so that they can work with custom per-Delivery Service routing names).

Apache Traffic Server 7.x (Cachekey Plugin)
-------------------------------------------
In Traffic Ops 2.2 we have added support for Apache Traffic Server (ATS) 7.x. With 7.x comes support for the new Cachekey Plugin which replaces the now-deprecated Cacheurl Plugin.
While not needed immediately it is recommended to start replacing Cacheurl usages with Cachekey as soon as possible, because ATS 6.x already supports the new Cachekey Plugin.

It is also recommended to thoroughly vet your Cachekey replacement by comparing with an existing key value. There are inconsistencies in the 6.x version of Cachekey which have been
fixed in 7.x (or require this patch(`cachekeypatch`_) on 6.x to match 7.x). So to ensure you have a matching key value you should use the XDebug Plugin before fully implementing your Cachekey replacement.

.. _cachekeypatch: https://github.com/apache/trafficserver/commit/244288fab01bdad823f9de19dcece62a7e2a0c11

First, if you are currently using a regular expression for your Delivery Service you will have to remove that existing value. Then you will need to make a new Delivery Service profile and assign parameters to it that use the ``cachekey.config`` file.

Some common parameters are

static-prefix
	This is used for a simple domain replacement.

separator
	Used by Cachekey and in general is always a single space.

remove-path
	Removes path information from the URL.

remove-all-params
	Removes query parameters from the URL.

capture-prefix-uri
	This is usually used in concert with remove-path and remove-all-params parameters. Capture-prefix-uri will let you use your own, full regular expression for non-trivial cases.

Examples of Cacheurl to Cachekey Replacements
---------------------------------------------

Static Prefix Example
"""""""""""""""""""""
Original regex value: ::

	http://test.net/(.*) http://test-cdn.net/$1

.. table:: Cachekey Parameters

	+---------------+-----------------+---------------------------------+
	|Parameter      |File             |Value                            |
	+===============+=================+=================================+
	| static-prefix | cachekey.config | ``http://test-cdn.net/``        |
	+---------------+-----------------+---------------------------------+
	| separator     | cachekey.config | (empty space)                   |
	+---------------+-----------------+---------------------------------+


Removing Query Parameters and Path Information
""""""""""""""""""""""""""""""""""""""""""""""
Original regex value: ::

	http://([^?]+)(?:?|$) http://test-cdn.net/$1

.. table:: Cachekey Parameters

	+-----------------------+-----------------+-----------------------------------------------------+
	|Parameter              |File             |Value                                                |
	+=======================+=================+=====================================================+
	| remove-path           | cachekey.config | true                                                |
	+-----------------------+-----------------+-----------------------------------------------------+
	| remove-all-params     | cachekey.config | true                                                |
	+-----------------------+-----------------+-----------------------------------------------------+
	| separator             | cachekey.config | (empty space)                                       |
	+-----------------------+-----------------+-----------------------------------------------------+
	| capture-prefix-uri    | cachekey.config | ``/https?:\/\/([^?]*)/http:\/\/test-cdn.net\/$1/``  |
	+-----------------------+-----------------+-----------------------------------------------------+

Also note the ``s?`` used here so that both HTTP and HTTPS requests will end up with the same key value


Removing Query Parameters with a Static Prefix
""""""""""""""""""""""""""""""""""""""""""""""
Original regex value: ::

	http://test.net/([^?]+)(?:\?|$) http://test-cdn.net/$1

.. table:: Cachekey Parameters

	+-------------------+-----------------+---------------------------------+
	|Parameter          |File             |Value                            |
	+===================+=================+=================================+
	| static-prefix     | cachekey.config | ``http://test-cdn.net/``        |
	+-------------------+-----------------+---------------------------------+
	| separator         | cachekey.config | (empty space)                   |
	+-------------------+-----------------+---------------------------------+
	| remove-all-params | cachekey.config | true                            |
	+-------------------+-----------------+---------------------------------+

.. note:: Further documentation on the Cachekey Plugin can be found at `the Apache Traffic Server Documentation <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/cachekey.en.html>`_.

Apache Traffic Server 7.x (Logging)
-----------------------------------
Traffic Server has changed the logging format as of version 7.0. Previously, this was ``logs_xml.config`` - an XML file - and now it is ``logging.config`` - a Lua file. Traffic Control compensates for this
automatically depending upon the filename used for the logging parameters. The same parameters will work this new file, ``LogFormat.Format``, ``LogFormat.Name``, ``LogObject.Format`` etc.


Traffic Ops Profile Modifications
---------------------------------
When upgrading to ATS 7.x, the Traffic Ops EDGE and MID cache profiles must be modified to provide new configuration values. Traffic Server's recommended parameter changes can be found `on their wiki <https://cwiki.apache.org/confluence/display/TS/Upgrading+to+v7.0>`_.

Most users of Traffic Control have enough profiles to make the task of making these modifications manually a tedious and time-consuming process. A new utility ``traffic_ops/install/bin/convert_profile/convert_profile`` is provided to automatically convert an ATS 6.x profile into an ATS 7.x profile. This utility can be reused in the future for converting ATS 7.x profiles into ATS 8.x profiles.

Usage Example
"""""""""""""
#. Use Traffic Portal GUI to export profile to JSON ('Configure' -> 'Profiles' -> Desired profile -> 'More' -> 'Export Profile')
#. Modify the Traffic Server version numbers to match your current Traffic Server 6.x RPM version and planned Traffic Server 7.x RPM version
#. Run ``convert_profile -input_profile <exported_file> -rules convert622to713.json -out <new_profile_name>``
#. Review output messages and make manual updates as needed. If you have modified a default value which the script also wants to change, it will prompt you to make the update manually. You may either do this directly in the JSON file or through the Traffic Portal GUI after import.
#. Use Traffic Portal GUI to import the newly created profile ('Configure' -> 'Profiles' -> 'More' -> 'Import Profile')

.. [1] Another name previously used for HTTP Delivery Services was ``ccr``
