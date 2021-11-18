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

.. _to-api-v3-deliveryservices-request:

****************************
``deliveryservices/request``
****************************

.. note:: This route does NOT do the same thing as :ref:`POST deliveryservice_requests<to-api-v3-deliveryservice-requests-post>`.

.. deprecated:: ATCv6
	This endpoint does not appear in Traffic Ops API version 4.0 - released with Apache Traffic Control version 6.0 - or later.

``POST``
========
Submits an emailed requesting that a :term:`Delivery Service` be created.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  ``undefined``

Request Structure
-----------------
:details: An object describing the actual parameters for the Delivery Service request

	:customer: Name of the customer associated with the :term:`Delivery Service` - must only contain alphanumeric characters and the characters :kbd:`@`, :kbd:`!`, :kbd:`#`, :kbd:`$`, :kbd:`%`, :kbd:`^`, :kbd:`&`, :kbd:`*`, :kbd:`(`, :kbd:`)`, :kbd:`[`, :kbd:`]`, :kbd:`.`, :kbd:`\ `, and :kbd:`-`

	.. versionchanged:: ATCv6
		Prior to ATC version 6, this field had no restrictions.

	:deepCachingType: An optional string describing when to do Deep Caching for this :term:`Delivery Service` - one of:

		NEVER
			Never use deep caching (default)
		ALWAYS
			Always use deep caching

	:deliveryProtocol: The protocol used to retrieve content from the CDN - one of:

		* http
		* https
		* http/https

	:hasNegativeCachingCustomization:  ``true`` if any customization is required for negative caching, ``false`` otherwise
	:hasOriginACLWhitelist:            ``true`` if access to the origin is restricted using an Access Control List (ACL or "whitelist") of IP addresses
	:hasOriginDynamicRemap:            If ``true``, this :term:`Delivery Service` can dynamically map to multiple origin URLs
	:hasSignedURLs:                    If ``true``, this :term:`Delivery Service`'s URLs are signed
	:headerRewriteEdge:                An optional string containing a header re-write rule to be used at the Edge tier
	:headerRewriteMid:                 An optional string containing a header re-write rule to be used at the Mid tier
	:headerRewriteRedirectRouter:      An optional string containing a header re-write rule to be used by the Traffic Router
	:maxLibrarySizeEstimate:           A special string that describes the estimated size of the sum total of content available through this :term:`Delivery Service`
	:negativeCachingCustomizationNote: A note remarking on the use, customization, or complications associated with negative caching for this :term:`Delivery Service`
	:notes:                            An optional string containing additional instructions or notes regarding the Request
	:originHeaders:                    An optional, comma-separated string of header values that must be passed to requests to the :term:`Delivery Service`'s origin
	:originTestFile:                   A URL path to a test file available on the :term:`Delivery Service`'s origin server
	:originURL:                        The URL of the :term:`Delivery Service`'s origin server
	:otherOriginSecurity:              An optional string describing any and all other origin security measures that need to be considered for access to the :term:`Delivery Service`'s origin
	:overflowService:                  An optional string containing the IP address or URL of an overflow point (used if rate limits are met or exceeded
	:peakBPSEstimate:                  A special string describing the estimated peak data transfer rate of the :term:`Delivery Service` in Bytes Per Second (BPS)
	:peakTPSEstimate:                  A special string describing the estimated peak transaction rate of the :term:`Delivery Service` in Transactions Per Second (TPS)
	:queryStringHandling:              A special string describing how the :term:`Delivery Service` should treat URLs containing query parameters
	:rangeRequestHandling:             A special string describing how the :term:`Delivery Service` should handle range requests
	:rateLimitingGBPS:                 An optional field which, if defined, should contain the maximum allowed data transfer rate for the :term:`Delivery Service` in GigaBytes Per Second (GBPS)
	:rateLimitingTPS:                  An optional field which, if defined, should contain the maximum allowed transaction rate for the :term:`Delivery Service` in Transactions Per Second (TPS)
	:routingName:                      An optional field which, if defined, should contain the routing name for the :term:`Delivery Service`, e.g. ``SomeRoutingName.DeliveryService_xml_id.CDNName.com``
	:routingType:                      The :term:`Delivery Service`'s routing type, should be one of:

		HTTP
			The Traffic Router re-directs clients to :term:`cache servers` using the HTTP ``302 REDIRECT`` response code
		DNS
			The Traffic Router responds to requests for name resolution of the :term:`Delivery Service`'s routing name with IP addresses of :term:`cache servers`
		STEERING
			This :term:`Delivery Service` routes clients to other :term:`Delivery Services` - which will in turn (generally) route them to clients
		ANY_MAP
			Some kind of undocumented black magic is used to get clients to... content, probably?

	:serviceAliases: An optional array of aliases for this :term:`Delivery Service`
	:serviceDesc:    A description of the :term:`Delivery Service`

:emailTo: The email to which the Delivery Service request will be sent

.. code-block:: json
	:caption: Request Example

	{ "emailTo": "foo@bar.com",
	"details": {
		"customer": "XYZ Corporation",
		"contentType": "static",
		"deepCachingType": "NEVER",
		"deliveryProtocol": "http",
		"routingType": "http",
		"routingName": "demo1",
		"serviceDesc": "service description goes here",
		"peakBPSEstimate": "less-than-5-Gbps",
		"peakTPSEstimate": "less-than-1000-TPS",
		"maxLibrarySizeEstimate": "less-than-200-GB",
		"originURL": "http://myorigin.com",
		"hasOriginDynamicRemap": false,
		"originTestFile": "http://origin.infra.ciab.test",
		"hasOriginACLWhitelist": false,
		"originHeaders": "",
		"otherOriginSecurity": "",
		"queryStringHandling": "ignore-in-cache-key-and-pass-up",
		"rangeRequestHandling": "range-requests-not-used",
		"hasSignedURLs": false,
		"hasNegativeCachingCustomization": false,
		"negativeCachingCustomizationNote": "",
		"serviceAliases": [],
		"rateLimitingGBPS": 50,
		"rateLimitingTPS": 5000,
		"overflowService": null,
		"headerRewriteEdge": "",
		"headerRewriteMid": "",
		"headerRewriteRedirectRouter": "",
		"notes": ""
	}}

Response Structure
------------------
.. code-block:: json
	:caption: Response Example

	{ "alerts": [{
		"level": "success",
		"text": "Delivery Service request sent to foo@bar.com."
	}]}
