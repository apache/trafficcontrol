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
