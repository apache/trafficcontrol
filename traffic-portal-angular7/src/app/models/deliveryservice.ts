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
export enum GeoLimit {
	None = 0,
	CZFOnly = 1,
	CZFAndCountryCodes = 2
};

export enum GeoProvider {
	MaxMind = 0,
	Neustar = 1
};

export class DeliveryServiceMatch {
	pattern: string;
	setNumber: number;
	type: string;
};

export enum Protocol {
	HTTP = 0,
	HTTPS = 1,
	HTTP_AND_HTTPS = 2,
	HTTP_TO_HTTPS = 3
};

export enum QStringHandling {
	USE = 0,
	IGNORE = 1,
	DROP = 2
};

export enum RangeRequestHandling {
	NONE = 0,
	BACKGROUND_FETCH = 1,
	CACHE_RANGE_REQUESTS = 2
};

export class DeliveryService {
	active: boolean;
	anonymousBlockingEnabled: boolean;
	cacheurl?: string;
	ccrDnsTtl: number;
	cdnId: number;
	cdnName?: string;
	checkPath?: string;
	deepCachingType?: string;
	displayName: string;
	dnsBypassCname?: string;
	dnsBypassIp?: string;
	dnsBypassIp6?: string;
	dnsBypassTtl?: number;
	dscp: number;
	edgeHeaderRewrite?: string;
	exampleURLs?: string;
	fqPacingRateL?: number;
	geoLimit: GeoLimit;
	geoLimitCountries?: string;
	geoLimitRedirectURL?: string;
	geoProvider: GeoProvider;
	globalMaxMbps?: string;
	globalMaxTps?: string;
	httpBypassFqdn?: string;
	id: number;
	infoUrl?: string;
	initialDispersion?: number;
	ipv6RoutingEnabled: boolean;
	lastUpdated: Date;
	logsEnabled: boolean;
	longDesc: string;
	longDesc1?: string;
	longDesc2?: string;
	matchList?: DeliveryServiceMatch[];
	maxDnsAnswers?: number;
	midHeaderRewrite?: string;
	missLat: number;
	missLong: number;
	multiSiteOrigin: boolean;
	orgServerFqdn?: string;
	originShield?: string;
	profileDescription?: string;
	profileId?: number;
	profileName?: string;
	protocol?: Protocol;
	qstringIgnore?: QStringHandling;
	regexRemap?: string;
	regionalGeoBlocking: boolean;
	remapText?: string;
	routingName: string;
	signed?: boolean;
	signingAlgorithm?: string;
	sslKeyVersion?: number;
	tenant?: string;
	tenantId?: number;
	trRequestHeaders?: string;
	trResponseHeaders?: string;
	type?: string;
	typeId: number; // TODO: Deterministic? Use enum if so
	xmlId: string;
};
