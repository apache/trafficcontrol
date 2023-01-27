/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
/**
 * This file contains definitons for objects of which Delivery Services are
 * composed - including some convenience functions for the magic numbers used to
 * represent certain properties.
 */

/**
 * Represents the `geoLimit` field of a Delivery Service
 */
export const enum GeoLimit {
	/**
	 * No geographic limiting is to be done.
	 */
	NONE = 0,
	/**
	 * Only clients found in a Coverage Zone File may be permitted access.
	 */
	CZF_ONLY = 1,
	/**
	 * Only clients found in a Coverage Zone File OR can be geo-located within a
	 * set of country codes may be permitted access.
	 */
	CZF_AND_COUNTRY_CODES = 2
}

/**
 * Defines the supported Geograhic IP mapping database providers and their
 * respective magic number identifiers.
 */
export const enum GeoProvider {
	/** The standard database used for geo-location. */
	MAX_MIND = 0,
	/** An alternative database with dubious support. */
	NEUSTAR = 1
}

/**
 * Represents a single entry in a Delivery Service's `matchList` field.
 */
export interface DeliveryServiceMatch {
	/** A regular expression matching on something depending on the 'type'. */
	pattern: string;
	/**
	 * The number in the set of the expression, which has vague but incredibly
	 * important meaning.
	 */
	setNumber: number;
	/**
	 * The type of match which determines how it's used.
	 */
	type: string;
}

/**
 * Represents the allowed routing protocols and their respective magic number
 * identifiers.
 */
export const enum Protocol {
	/** Serve HTTP traffic only. */
	HTTP = 0,
	/** Serve HTTPS traffic only. */
	HTTPS = 1,
	/** Serve both HTTPS and HTTP traffic. */
	HTTP_AND_HTTPS = 2,
	/** Redirect HTTP requests to HTTPS URLs and serve HTTPS traffic normally. */
	HTTP_TO_HTTPS = 3
}

/**
 * Represents the allowed values of the `qstringIgnore` field of a
 * `DeliveryService`.
 */
export const enum QStringHandling {
	/** Use the query string in the cache key and pass in upstream requests. */
	USE = 0,
	/**
	 * Don't use the query string in the cache key but do pass it in upstream
	 * requests.
	 */
	IGNORE = 1,
	/**
	 * Neither use the query string in the cache key nor pass it in upstream
	 * requests.
	 */
	DROP = 2
}

/**
 * Represents the allowed values of the `rangeRequestHandling` field of a
 * `Delivery Service`.
 */
export const enum RangeRequestHandling {
	/** Range requests will not be cached. */
	NONE = 0,
	/**
	 * The entire object will be fetched in the background to be cached, with
	 * the requested range served when it becomes available.
	 */
	BACKGROUND_FETCH = 1,
	/**
	 * Cache range requests like any other request.
	 */
	CACHE_RANGE_REQUESTS = 2
}

/**
 * Represents a single Delivery Service of arbitrary type
 */
export interface DeliveryService {
	/** Whether or not the Delivery Service is actively routed. */
	active:                     boolean;
	/** Whether or not anonymization services are blocked. */
	anonymousBlockingEnabled:   boolean;
	/** The TTL of DNS responses from the Traffic Router, in seconds. */
	ccrDnsTtl?:                 number;
	/** The ID of the CDN to which the Delivery Service belongs. */
	cdnId:                      number;
	/** The Name of the CDN to which the Delivery Service belongs. */
	cdnName?:                   string;
	/** A sample path which may be requested to ensure the origin is working. */
	checkPath?:                 string;
	/**
	 * A regular expression used to extract request path fragments for use as
	 * keys in "consistent hashing" for routing purposes.
	 */
	consistentHashRegex?:       string;
	/**
	 * A set of the query parameters that are important for Traffic Router to
	 * consider when performing "consistent hashing".
	 */
	consistentHashQueryParams?: Array<string>;
	/**
	 * Whether or not to use "deep caching".
	 */
	deepCachingType?:           string;
	/** A human-friendly name for the Delivery Service. */
	displayName:                string;
	/** An FQDN to use for DNS-routed bypass scenarios. */
	dnsBypassCname?:            string;
	/** An IPv4 address to use for DNS-routed bypass scenarios. */
	dnsBypassIp?:               string;
	/** An IPv6 address to use for DNS-routed bypass scenarios. */
	dnsBypassIp6?:              string;
	/** The TTL of DNS responses served in bypass scenarios. */
	dnsBypassTtl?:              number;
	/** The Delivery Service's DSCP. */
	dscp:                       number;
	/** Extra header rewrite text used at the Edge tier. */
	edgeHeaderRewrite?:         string;
	/**
	 * A list of the URLs which may be used to request Delivery Service content.
	 */
	exampleURLs?:               Array<string>;
	/**
	 * Describes limitation of content availability based on geographic
	 * location.
	 */
	geoLimit:                   GeoLimit;
	/**
	 * The countries from which content access is allowed in the event that
	 * geographic limiting is taking place with a setting that allows for
	 * specific country codes to access content.
	 */
	geoLimitCountries?:         string;
	/**
	 * A URL to which to re-direct users who are blocked because of
	 * geographic-based access limiting
	 */
	geoLimitRedirectURL?:       string;
	/**
	 * The provider of the IP-address-to-geographic-location database.
	 */
	geoProvider:                GeoProvider;
	/**
	 * The globally allowed maximum megabits per second to be served for the
	 * Delivery Service.
	 */
	globalMaxMbps?:             string;
	/**
	 * The globally allowed maximum transactions per second to be served for the
	 * Delivery Service.
	 */
	globalMaxTps?:              string;
	/**
	 * A URL to be used in HTTP-routed bypass scenarios.
	 */
	httpBypassFqdn?:            string;
	/**
	 * An integral, unique identifier for the Delivery Service.
	 */
	id?:                        number;
	/**
	 * A URL from which information about a Delivery Service may be obtained.
	 * Historically, this has been used to link to the support ticket that
	 * requested the Delivery Service's creation.
	 */
	infoUrl?:                   string;
	/**
	 * The number of caches across which to spread content.
	 */
	initialDispersion?:         number;
	/**
	 * whether or not routing of IPv6 clients should be supported.
	 */
	ipv6RoutingEnabled:         boolean;
	/** When the Delivery Service was last updated via the API. */
	lastUpdated?:               Date;
	/** Whether or not logging should be enabled for the Delivery Service. */
	logsEnabled:                boolean;
	/** A textual description of arbitrary length. */
	longDesc:                   string;
	/** A textual description of arbitrary length. */
	longDesc1?:                 string;
	/** A textual description of arbitrary length. */
	longDesc2?:                 string;
	/**
	 * A list of regular expressions for routing purposes which should not ever
	 * be modified by hand.
	 */
	matchList?:                 DeliveryServiceMatch[];
	/**
	 * Sets the maximum number of answers Traffic Router may provide in a single
	 * DNS response for this Delivery Service.
	 */
	maxDnsAnswers?:             number;
	/**
	 * The maximum number of connections that cache servers are allowed to open
	 * to the Origin(s).
	 */
	maxOriginConnections?:      number;
	/** Extra header rewrite text to be used at the Mid-tier. */
	midHeaderRewrite?:          string;
	/** The latitude that should be used when geo-location of a client fails. */
	missLat:                    number;
	/** The longitude that should be used when geo-location of a client fails. */
	missLong:                   number;
	/** Whether or not Multi-Site Origin is in use. */
	multiSiteOrigin:            boolean;
	/** The URL of the Origin server, which I think means nothing for MSO. */
	orgServerFqdn?:             string;
	/** A string used to shield the Origin, somehow. */
	originShield?:              string;
	/** A description of the Profile used by the Delivery Service (read-only) */
	profileDescription?:        string;
	/** An integral, unique identifer for the Profile used by the Delivery Service. */
	profileId?:                 number;
	/** The name of the Profile used by the Delivery Service. */
	profileName?:               string;
	/** The protocols served by the Delivery Service. */
	protocol?:                  Protocol;
	/**
	 * How query strings ought to be handled by cache servers serving content
	 * for this Delivery Service.
	 */
	qstringIgnore?:             QStringHandling;
	/**
	 * How HTTP Range requests ought to be handled by cache servers serving
	 * content for this Delivery Service.
	 */
	rangeRequestHandling?:      RangeRequestHandling;
	/**
	 * some raw text to be inserted into regex_remap.config.
	 */
	regexRemap?:                string;
	/**
	 * Whether or not regional geo-blocking should be used.
	 */
	regionalGeoBlocking:        boolean;
	/** some raw text to be inserted into remap.config. */
	remapText?:                 string;
	/** The lowest-level DNS label used in URLs for the Delivery Service. */
	routingName:                string;
	/**
	 * Whether or not responses from the cache servers for this Delivery
	 * Service's content will be signed.
	 */
	signed?:                    boolean;
	/**
	 * The algorithm used to sign responses from the cache servers for this
	 * Delivery Service's content.
	 */
	signingAlgorithm?:          string;
	/**
	 * The generation of SSL key used by this Delivery Service.
	 */
	sslKeyVersion?:             number;
	/** The name of the Tenant to whom this Delivery Service belongs. */
	tenant?:                    string;
	/**
	 * An integral, unique identifier for the Tenant to whom this Delivery
	 * Service belongs.
	 */
	tenantId:                   number;
	/**
	 * HTTP headers that should be logged from client requests by Traffic
	 * Router.
	 */
	trRequestHeaders?:          string;
	/**
	 * Extra HTTP headers that Traffic Router should provide in responses.
	 */
	trResponseHeaders?:         string;
	/** The type of the Delivery Service. */
	type?:                      string;
	/** The integral, unique identifier of the type of the Delivery Service. */
	typeId:                     number;
	/** The second-lowest-level DNS label used in the Delivery Service's URLs. */
	xmlId:                      string;
}

/**
 * DSCapacity represents a response from the API to a request for the capacity
 * of a Delivery Service.
 */
export interface DSCapacity {
	availablePercent: number;
	maintenancePercent: number;
	utilizedPercent: number;
}

/**
 * DSHealth represents a response from the API to a request for the health of a
 * Delivery Service.
 */
export interface DSHealth {
	totalOnline: number;
	totalOffline: number;
}
