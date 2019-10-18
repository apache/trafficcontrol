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
 * This file contains definitons for objects of which Delivery Services are composed - including
 * some convenience functions for the magic numbers used to represent certain properties.
*/

/**
 * Represents the `geoLimit` field of a Delivery Service
*/
export enum GeoLimit {
	None = 0,
	CZFOnly = 1,
	CZFAndCountryCodes = 2
}

/**
 * This namespace merges with the `GeoLimit` enum to provide a seamless method on the object. The
 * mechanics of this are a mystery to me.
*/
export namespace GeoLimit {
	export function toString (g: GeoLimit): string {
		switch (g) {
			case GeoLimit.None:
				return 'None';
			case GeoLimit.CZFOnly:
				return "Serve content only if the client's IP is found in the Coverage Zone File";
			case GeoLimit.CZFAndCountryCodes:
				/* tslint:disable */
				return "Serve content only if the client's IP is found in the Coverage Zone File OR if the client can be determined to be within a country specified by the 'GeoLimit Countries' list";
				/* tslint:enable */
			default:
				return 'UNKNOWN';
		}
	}
}

/**
 * Defines the supported Geograhic IP mapping database providers and their respective magic number
 * identifiers.
*/
export enum GeoProvider {
	MaxMind = 0,
	Neustar = 1
}

/**
 * Represents a single entry in a Delivery Service's `matchList` field.
*/
export interface DeliveryServiceMatch {
	pattern: string;
	setNumber: number;
	type: string;
}

/**
 * Represents the allowed routing protocols and their respective magic number identifiers.
*/
export enum Protocol {
	HTTP = 0,
	HTTPS = 1,
	HTTP_AND_HTTPS = 2,
	HTTP_TO_HTTPS = 3
}

/**
 * This namespace merges with the `Protocol` enum to provide a seamless method to convert those
 * values to verbose explanations.
*/
export namespace Protocol {
	export function toString (p: Protocol): string {
		switch (p) {
			case Protocol.HTTP:
				return 'Serve only unsecured HTTP requests';
			case Protocol.HTTPS:
				return 'Serve only secured HTTPS requests';
			case Protocol.HTTP_AND_HTTPS:
				return 'Serve both unsecured HTTP requests and secured HTTPS requests';
			case Protocol.HTTP_TO_HTTPS:
				return 'Serve secured HTTPS requests normally, but redirect unsecured HTTP requests to use HTTPS';
			default:
				return 'UNKNOWN';
		}
	}
}

/**
 * Represents the allowed values of the `qstringIgnore` field of a `DeliveryService`
*/
export enum QStringHandling {
	USE = 0,
	IGNORE = 1,
	DROP = 2
}

/**
 * This namespace merges with the `QStringHandling` enum to provide a seamless method to convert
 * those values to verbose explanations.
*/
export namespace QStringHandling {
	export function toString (q: QStringHandling): string {
		switch (q) {
			case QStringHandling.USE:
				return 'Use the query parameter string when deciding if a URL is cached, and pass it in upstream requests to the Mid-tier/origin';
			case QStringHandling.IGNORE:
				/* tslint:disable */
				return 'Do not use the query parameter string when deciding if a URL is cached, but do pass it in upstream requests to the Mid-tier/origin';
				/* tslint:enable */
			case QStringHandling.DROP:
				return 'Immediately strip URLs of their query parameter strings before checking cached objects or making upstream requests';
			default:
				return 'UNKNOWN';
		}
	}
}

/**
 * Represents the allowed values of the `rangeRequestHandling` field of a `Delivery Service`
*/
export enum RangeRequestHandling {
	NONE = 0,
	BACKGROUND_FETCH = 1,
	CACHE_RANGE_REQUESTS = 2
}

/**
 * This namespace merges with the `RangeRequestHandling` enum to provide a seamless method to convert
 * those values to verbose explanations.
*/
export namespace RangeRequestHandling {
	export function toString (r: RangeRequestHandling): string {
		switch (r) {
			case RangeRequestHandling.NONE:
				return 'Do not cache Range requests';
			case RangeRequestHandling.BACKGROUND_FETCH:
				return 'Use the background_fetch plugin to serve Range requests while quietly caching the entire object';
			case RangeRequestHandling.CACHE_RANGE_REQUESTS:
				return 'Use the cache_range_requests plugin to directly cache object ranges';
		}
	}
}

/**
 * Represents a single Delivery Service of arbitrary type
*/
export interface DeliveryService {
	active:                     boolean;
	anonymousBlockingEnabled:   boolean;
	cacheurl?:                  string;
	ccrDnsTtl?:                 number;
	cdnId:                      number;
	cdnName?:                   string;
	checkPath?:                 string;
	consistentHashRegex?:       RegExp;
	consistentHashQueryParams?: Array<string>;
	deepCachingType?:           string;
	displayName:                string;
	dnsBypassCname?:            string;
	dnsBypassIp?:               string;
	dnsBypassIp6?:              string;
	dnsBypassTtl?:              number;
	dscp:                       number;
	edgeHeaderRewrite?:         string;
	exampleURLs?:               Array<string>;
	fqPacingRateL?:             number;
	geoLimit:                   GeoLimit;
	geoLimitCountries?:         string;
	geoLimitRedirectURL?:       string;
	geoProvider:                GeoProvider;
	globalMaxMbps?:             string;
	globalMaxTps?:              string;
	httpBypassFqdn?:            string;
	id?:                        number;
	infoUrl?:                   string;
	initialDispersion?:         number;
	ipv6RoutingEnabled:         boolean;
	lastUpdated?:               Date;
	logsEnabled:                boolean;
	longDesc:                   string;
	longDesc1?:                 string;
	longDesc2?:                 string;
	matchList?:                 DeliveryServiceMatch[];
	maxDnsAnswers?:             number;
	maxQriginConnections?:      number;
	midHeaderRewrite?:          string;
	missLat:                    number;
	missLong:                   number;
	multiSiteOrigin:            boolean;
	orgServerFqdn?:             string;
	originShield?:              string;
	profileDescription?:        string;
	profileId?:                 number;
	profileName?:               string;
	protocol?:                  Protocol;
	qstringIgnore?:             QStringHandling;
	rangeRequestHandling?:      RangeRequestHandling;
	regexRemap?:                string;
	regionalGeoBlocking:        boolean;
	remapText?:                 string;
	routingName:                string;
	signed?:                    boolean;
	signingAlgorithm?:          string;
	sslKeyVersion?:             number;
	tenant?:                    string;
	tenantId?:                  number;
	trRequestHeaders?:          string;
	trResponseHeaders?:         string;
	type?:                      string;
	typeId:                     number;
	xmlId:                      string;
}

/**
 * Determines if the Delivery Service is a candidate for bypassing
 * @returns `true` if it can have bypass settings, `false` otherwise.
*/
export function bypassable (ds: DeliveryService): boolean {
	if (!ds.type) {
		return false;
	}

	switch (ds.type) {
		case 'HTTP':
		case 'HTTP_LIVE':
		case 'HTTP_LIVE_NATNL':
		case 'DNS':
		case 'DNS_LIVE':
		case 'DNS_LIVE_NATNL':
			return true;
		default:
			return false;
	}
}
