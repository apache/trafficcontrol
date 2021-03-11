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

.. _regionalgeo-qht:

*************************************
Configure Regional Geo-blocking (RGB)
*************************************
.. Note:: :abbr:`RGB (Regional Geographic-based Blocking)` is only supported for HTTP :term:`Delivery Services`.

#. Prepare an :abbr:`RGB (Regional Geographic-based Blocking)` configuration file. :abbr:`RGB (Regional Geographic-based Blocking)` uses a configuration file in JSON format to define regional geographic blocking rules for :term:`Delivery Services`. The file needs to be put on an HTTP server accessible to Traffic Router.

	.. code-block:: json
		:caption: Example Configuration File

		{
		"deliveryServices":
			[
				{
					"deliveryServiceId": "hls-live",
					"urlRegex": ".*live4\\.m3u8",
					"geoLocation": {"includePostalCode":["N0H", "L9V", "L9W"],
									"coordinateRange": [{"minLat" : -12, "maxLat": 13, "minLon" : 55, "maxLon": 56}, {"minLat" : -13, "maxLat": 14, "minLon" : 55, "maxLon": 56}]},
					"redirectUrl": "http://third-party.com/blacked_out.html"
				},
				{
					"deliveryServiceId": "hls-live",
					"urlRegex": ".*live5\\.m3u8",
					"ipWhiteList": ["185.68.71.9/22","142.232.0.79/24"],
					"geoLocation": {"excludePostalCode":["N0H", "L9V"]},
					"redirectUrl": "/live5_low_bitrate.m3u8",
					"isSteeringDS": "false"
				},
				{
					"deliveryServiceId": "linear-steering",
					"urlRegex": ".*live3\\.m3u8",
					"ipWhiteList": ["185.68.71.9/22","142.232.0.79/24"],
					"geoLocation": {"excludePostalCode":["N0H", "L9V"]},
					"redirectUrl": "http://ip-slate.cdn.example.com/slate.m3u8",
					"isSteeringDS": "true"
				}
			]
		}

	``deliveryServiceId``
		Should be equal to the ``ID`` or ``xml_id`` field of the intended :term:`Delivery Service` as configured in Traffic Portal
	``urlRegex``
		A regular expression to be used to determine to what URLs the rule shall apply; a URL that matches it is subject to the rule
	``geoLocation``
		An object that currently supports only the keys ``includePostalCode``, ``excludePostalCode`` (mutually exclusive) and ``coordinateRange``. When the ``includePostalCode`` key is used, only the clients whose :abbr:`FSA (Forward Sortation Areas)`\ s - the first three postal characters of Canadian postal codes - are in the ``includePostalCode`` list are able to view the content at URLs matched by the ``urlRegex``. When ``excludePostalCode`` is used, any client whose :abbr:`FSA (Forward Sortation Areas)` is not in the ``excludePostalCode`` list will be allowed to view the content. The ``coordinateRange`` key is used to specify a list of latitude and longitude ranges. This is used in regional geo blocking, in case the client does not have a postal code associated with it.
	``redirectUrl``
		The URL that will be returned to the blocked clients. Without a domain name in the URL, the URL will still be served in the same :term:`Delivery Service`. Thus Traffic Router will redirect the client to a chosen :term:`cache server` assigned to the :term:`Delivery Service`. If the URL includes a domain name, Traffic Router simply redirects the client to the defined URL. In the latter case, the redirect URL must not match the ``urlRegex`` value, or an infinite loop of  HTTP ``302 Found`` responses will occur at the Traffic Router.  Steering-:ref:`ds-types` :term:`Delivery Services` must contain an :abbr:`FQDN (Fully Qualified Domain Name)` as the re-direct or Traffic Router will return a DENIED to the client.  This is because steering services do not have caches associated to them, so a relative ``redirectURL`` can not be turned into a :abbr:`FQDN (Fully Qualified Domain Name)`.
	``ipWhiteList``
		An optional element that is an array of :abbr:`CIDR (Classless Inter-Domain Routing)` blocks indicating the IPv4 and/or IPv6 subnets that are allowed by the rule. If this list exists and the value is not empty, client IP will be matched against the :abbr:`CIDR (Classless Inter-Domain Routing)` list, bypassing the value of ``geoLocation``. If there is no match in the white list, Traffic Router defers to the value of ``geoLocation`` to determine if content ought to be blocked.


#. Add :abbr:`RGB (Regional Geographic-based Blocking)` :term:`Parameters` in Traffic Portal to the :term:`Delivery Service`'s Traffic Router(s)'s :term:`Profile`\ (s). The :ref:`parameter-config-file` value should be set to ``CRConfig.json``, and the following two :term:`Parameter` :ref:`parameter-name`/:ref:`parameter-value` pairs need to be specified:

	``regional_geoblock.polling.url``
		The URL of the RGB configuration file. Traffic Router will fetch the file from this URL using an HTTP ``GET`` request.
	``regional_geoblock.polling.interval``
		The interval on which Traffic Router polls the :abbr:`RGB (Regional Geographic-based Blocking)` configuration file.

	.. figure:: regionalgeo/01.png
		:width: 40%
		:align: center

#. Enable :abbr:`RGB (Regional Geographic-based Blocking)` for a :term:`Delivery Service` using the :ref:`Delivery Services view in Traffic Portal <tp-services-delivery-service>` (don't forget to save changes!)

	.. figure:: regionalgeo/02.png
		:width: 40%
		:align: center

#. Go to :ref:`the Traffic Portal CDNs view <tp-cdns>`, click on :guilabel:`Diff CDN Config Snapshot`, and click :guilabel:`Perform Snapshot`.

	.. figure:: regionalgeo/03.png
		:width: 40%
		:align: center

Traffic Router Access Log
=========================
.. seealso:: :ref:`tr-logs`

RGB extends the ``rtype`` field and adds a new field ``rgb`` in Traffic Router access.log to help to monitor this feature. A value of ``RGALT`` in the ``rtype`` field indicates that a request is redirected to an alternate URL by :abbr:`RGB (Regional Geographic-based Blocking)`; a value of ``RGDENY`` indicates that a request is denied by :abbr:`RGB (Regional Geographic-based Blocking)` because there is no matching rule in the :abbr:`RGB (Regional Geographic-based Blocking)` configuration file for this request. When :abbr:`RGB (Regional Geographic-based Blocking)` is enabled, the ``RGB`` field will be non-empty with following format:

``{FSA}:{allowed/disallowed}:{include/exclude postal}:{fallback config}:{allowed by whitelist}``


FSA
	:dfn:`FSA` part of the client’s postal code, which is retrieved from a geographic location database. If this field is empty, a dash (“-“) is filled in.
allowed/disallowed
	This flag shows if a request was allowed or disallowed by :abbr:`RGB (Regional Geographic-based Blocking)` (1 for yes, and 0 for no).
include/exclude postal
	This shows that when a rule in JSON is matched for a request, it's value is "I" if the rule matched because of an ``includePostalCode`` rule, "X" if the rule matched because of an ``excludePostalCode`` rule, or "-" if no rule matched.
fallback config
	When Traffic Router fails to parse an :abbr:`RGB (Regional Geographic-based Blocking)` configuration file as JSON, Traffic Router will handle requests with latest valid configuration that it had, but will set the ``fallback config`` flag to 1. If no fall-back occurred, then the flag is set to 0.
allowed by whitelist
	If a request is allowed by a ``whitelist`` field in the configuration, this flag is set to 1; for all other cases, it is 0.


.. code-block:: squid
	:caption: Example

	1446442214.685 qtype=HTTP chi=129.100.254.79 url="http://foo.geo2.cdn.com/live5.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=GEO rloc="-" rdtl=- rerr="-" rgb="N6G:1:X:0:0" pssc=302 ttms=3 rurl=http://cent6-44.geo2.cdn.com/live5.m3u8 rh="-"

	1446442219.181 qtype=HTTP chi=184.68.71.9 url="http://foo.geo2.cdn.com/live5.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=RGALT rloc="-" rdtl=- rerr="-" rgb="-:0:X:0:0" pssc=302 ttms=3 rurl=http://cent6-44.geo2.cdn.com/low_bitrate.m3u8 rh="-"

	1446445521.677 qtype=HTTP chi=24.114.29.79 url="http://foo.geo2.cdn.com/live51.m3u8" cqhm=GET cqhv=HTTP/1.1 rtype=RGDENY rloc="-" rdtl=- rerr="-" rgb="L4S:0:-:0:0" pssc=520 ttms=3 rurl="-" rh="-"


