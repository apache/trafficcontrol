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
import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";

import { Observable, of } from "rxjs";
import { merge } from "rxjs/index";
import { map, mergeAll, reduce } from "rxjs/operators";

import {
	CDN,
	DataPoint,
	DataSet,
	DataSetWithSummary,
	DeliveryService,
	DSCapacity,
	DSHealth,
	InvalidationJob,
	TPSData,
	Type
} from "../../models";
import { APIService } from "./apiservice";

/**
 * The type of a raw response returned from the API that has to be massaged
 * into a DataSet.
 */
interface DataResponse {
	series: {
		name: string;
		values: Array<[number, number | null]>;
	};
	summary?: {
		min: number;
		max: number;
		fifthPercentile: number;
		ninetyFifthPercentile: number;
		ninetyEightPercentile: number;
		mean: number;
	};
}

/**
 * Checks that a given object represents a proper data set.
 *
 * @param r The 'response' object from the API response.
 * @returns Always 'true' - if the assertion fails, an error is thrown rather than returning 'false'.
 * @throws When 'r' is not a 'DataResponse'.
 */
function isDataResponse(r: object): r is DataResponse {
	if (!Object.prototype.hasOwnProperty.call(r, "series")) {
		throw new Error("no series data");
	}
	if (!Object.prototype.hasOwnProperty.call((r as {series: unknown}), "series")) {
		throw new Error("series data has no name");
	}
	const nameType = typeof((r as {series: {name: unknown}}).series.name);
	if (nameType !== "string") {
		throw new Error(`invalid series name, expected a string, got ${nameType}`);
	}
	if (!Object.prototype.hasOwnProperty.call((r as {series: object}).series, "values") ||
		(r as {series: {values: unknown}}).series === null) {
		// just fix this silently.
		(r as {series: Record<symbol | string, unknown>}).series.values = new Array<[number, number]>();
	} else if (!((r as {series: {values: unknown}}).series.values instanceof Array)) {
		throw new Error(`series values are not an array or missing/null, got: ${typeof(r as {series: {values: unknown}}).series.values}`);
	}

	// At this point we assume the summary data either isn't present or
	// is fully compliant with the expected format. That's because the
	// common problem is old API versions not returning the 'series'
	// property - there is no known issue that would cause it to not
	// return a proper 'summary' (if one is returned at all).
	return true;
}

/**
 * Constructs a data set from the API response.
 *
 * @param r The parsed response body.
 * @returns A DataSetWithSummary that was massaged out of the raw response.
 */
function constructDataSetFromResponse(r: object): DataSetWithSummary {
	try {
		if (!isDataResponse(r)) {
			throw new Error("response is not a data series");
		}
	} catch (e) {
		console.log("response:", r);
		throw new Error(`invalid data set response: ${e}`);
	}

	const data = new Array<DataPoint>();
	for (const v of r.series.values) {
		if (v[1] === null) {
			continue;
		}
		data.push({t: new Date(v[0]), y: v[1].toFixed(3)});
	}

	let min: number;
	let max: number;
	let fifth: number;
	let nfifth: number;
	let neight: number;
	let mean: number;
	if (r.summary) {
		min = r.summary.min;
		max = r.summary.max;
		fifth = r.summary.fifthPercentile;
		nfifth = r.summary.ninetyFifthPercentile;
		neight = r.summary.ninetyEightPercentile;
		mean = r.summary.mean;
	} else {
		min = -1;
		max = -1;
		fifth = -1;
		nfifth = -1;
		neight = -1;
		mean = -1;
	}

	return {
		dataSet: {data, label: r.series.name.split(".")[0]} as DataSet,
		fifthPercentile: fifth,
		max,
		mean,
		min,
		ninetyEighthPercentile: neight,
		ninetyFifthPercentile: nfifth
	} as DataSetWithSummary;
}

/**
 * DeliveryServiceService exposes API functionality related to Delivery Services.
 */
@Injectable({providedIn: "root"})
export class DeliveryServiceService extends APIService {

	/** This is where DS Types are cached, as they are presumed to not change (often). */
	private deliveryServiceTypes: Array<Type>;

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
		this.deliveryServiceTypes = new Array<Type>();
	}

	public getDeliveryServices(id: string | number): Observable<DeliveryService>;
	public getDeliveryServices(): Observable<Array<DeliveryService>>;
	/**
	 * Gets a list of all visible Delivery Services
	 *
	 * @param id A unique identifier for a Delivery Service - either a numeric id or an "xml_id"
	 * @throws TypeError if ``id`` is not a proper type
	 * @returns An observable that will emit an array of `DeliveryService` objects.
	 */
	public getDeliveryServices(id?: string | number): Observable<DeliveryService[] | DeliveryService> {
		const path = "deliveryservices";
		if (id) {
			let params;
			switch (typeof id) {
				case "string":
					// Part of the API spec, unfortunately
					// eslint-disable-next-line @typescript-eslint/naming-convention
					params = {xml_id: id};
					break;
				case "number":
					params = {id: String(id)};
			}
			return this.get<[DeliveryService]>(path, undefined, params).pipe(map(
				r => {
					const ds = r[0];
					ds.lastUpdated = new Date((ds.lastUpdated as unknown as string).replace("+00", "Z"));
					return ds;
				}
			));
		}
		return this.get<Array<DeliveryService>>(path).pipe(map(r => r.map(
			ds => {
				ds.lastUpdated = new Date((ds.lastUpdated as unknown as string).replace("+00", "Z"));
				return ds;
			}
		)));
	}

	/**
	 * Creates a new Delivery Service
	 *
	 * @param ds The new Delivery Service object
	 * @returns An Observable that will emit a boolean value indicating the success of the operation
	 */
	public createDeliveryService(ds: DeliveryService): Observable<boolean> {
		const path = "deliveryservices";
		return this.post<DeliveryService>(path, ds).pipe(map(
			() => true,
			() => false
		));
	}

	/**
	 * Retrieves capacity statistics for the Delivery Service identified by a given, unique,
	 * integral value.
	 *
	 * @param d Either a {@link DeliveryService} or an integral, unique identifier of a Delivery Service
	 * @returns An Observable that emits an object that hopefully has the right keys to represent capacity.
	 * @throws If `d` is a {@link DeliveryService} that has no (valid) id
	 */
	public getDSCapacity(d: number | DeliveryService): Observable<DSCapacity> {
		let id: number;
		if (typeof d === "number") {
			id = d;
		} else {
			d = d;
			if (!d.id || d.id < 0) {
				throw new Error("Delivery Service id must be defined!");
			}
			id = d.id;
		}

		const path = `deliveryservices/${id}/capacity`;
		return this.get<DSCapacity>(path).pipe();
	}

	/**
	 * Retrieves the Cache Group health of a Delivery Service identified by a given, unique,
	 * integral value.
	 *
	 * @param d The integral, unique identifier of a Delivery Service
	 * @returns An Observable that emits a response from the health endpoint
	 */
	public getDSHealth(d: number): Observable<DSHealth> {
		const path = `deliveryservices/${d}/health`;
		return this.get<DSHealth>(path).pipe();
	}

	public getDSKBPS(d: string, start: Date, end: Date, interval: string, useMids: boolean, dataOnly: true): Observable<Array<DataPoint>>;
	public getDSKBPS(d: string, start: Date, end: Date, interval: string, useMids: boolean, dataOnly?: false): Observable<DataResponse>;
	/**
	 * Retrieves Delivery Service throughput statistics for a given time period, averaged over a given
	 * interval.
	 *
	 * @param d The `xml_id` of a Delivery Service
	 * @param start A date/time from which to start data collection
	 * @param end A date/time at which to end data collection
	 * @param interval A unit-suffixed interval over which data will be "binned"
	 * @param useMids Collect data regarding Mid-tier cache servers rather than Edge-tier cache servers
	 * @param dataOnly Only returns the data series, not any supplementing meta info found in the API response
	 * @returns An Observable that will emit an Array of datapoint Arrays (length 2 containing a date string and data value)
	 */
	public getDSKBPS(
		d: string,
		start: Date,
		end: Date,
		interval: string,
		useMids: boolean,
		dataOnly?: boolean
	): Observable<Array<DataPoint> | DataResponse> {
		const path = "deliveryservice_stats";
		const params = {
			deliveryServiceName: d,
			endDate: end.toISOString(),
			interval,
			metricType: "kbps",
			serverType: useMids ? "mid" : "edge",
			startDate: start.toISOString()
		};
		return this.get<object>(path, undefined, params).pipe(map(
			r => {
				try {
					if (!isDataResponse(r)) {
						throw new Error("invalid data from getDSKBPS");
					}
				} catch (e) {
					throw new Error(`invalid data set returned from ${path}: ${e}`);
				}
				if (dataOnly) {
					if (r.hasOwnProperty("series") && (r.series.hasOwnProperty("values"))) {
						return r.series.values.filter(ds => ds[1] !== null).map(
							ds => ({
								t: new Date(ds[0]),
								y: (ds[1] as number).toFixed(3)
							})
						);
					}
					throw new Error(`no data series found (path was "${path}")`);
				}
				return r;
			}
		));
	}

	/**
	 * Gets total TPS data for a Delivery Service. To get TPS data broken down by HTTP status, use {@link getAllDSTPSData}.
	 *
	 * @param d The name (xmlid) of the Delivery Service for which TPS stats will be fetched
	 * @param start The desired start date/time of the data range (must not have nonzero milliseconds!)
	 * @param end The desired end date/time of the data range (must not have nonzero milliseconds!)
	 * @param interval A string that describes the interval across which to 'bucket' data e.g. '60s'
	 * @param useMids If given (and true) will get stats for the Mid-tier instead of the Edge-tier (which is the default behavior).
	 * @returns An Observable that will emit the requested DataResponse.
	 */
	public getDSTPS(
		d: string,
		start: Date,
		end: Date,
		interval: string,
		useMids?: boolean
	): Observable<DataResponse> {
		const path = "deliveryservice_stats";
		const params = {
			deliveryServiceName: d,
			endDate: end.toISOString(),
			interval,
			metricType: "tps_total",
			serverType: useMids ? "mid" : "edge",
			startDate: start.toISOString()
		};
		return this.get<DataResponse>(path, undefined, params).pipe();
	}

	/**
	 * Gets total TPS data for a Delivery Service, as well as TPS data by HTTP response type.
	 *
	 * @param d The name (xmlid) of the Delivery Service for which TPS stats will be fetched
	 * @param start The desired start date/time of the data range (must not have nonzero milliseconds!)
	 * @param end The desired end date/time of the data range (must not have nonzero milliseconds!)
	 * @param interval A string that describes the interval across which to 'bucket' data e.g. '60s'
	 * @param useMids If given (and true) will get stats for the Mid-tier instead of the Edge-tier (which is the default behavior)
	 * @returns An Observable that will emit the requested TPSData.
	 */
	public getAllDSTPSData(
		d: string,
		start: Date,
		end: Date,
		interval: string,
		useMids?: boolean
	): Observable<TPSData> {
		const path = "deliveryservice_stats";
		const params: Record<string, string> = {
			deliveryServiceName: d,
			endDate: end.toISOString(),
			interval,
			serverType: useMids ? "mid" : "edge",
			startDate: start.toISOString()
		};
		const metricTypes = [
			"tps_total",
			"tps_2xx",
			"tps_3xx",
			"tps_4xx",
			"tps_5xx",
		];

		const observables = metricTypes.map(
			x => this.get<object>(path, undefined, {metricType: x, ...params}).pipe(map(constructDataSetFromResponse))
		);

		// TODO: Figure out why this double-merge is necessary
		const tasks = merge(observables).pipe(mergeAll());
		return tasks.pipe(reduce(
			(output: TPSData, result: DataSetWithSummary): TPSData => {
				switch (result.dataSet.label) {
					case "tps_total":
						output.total = result;
						break;
					case "tps_1xx":
						output.informational = result;
						break;
					case "tps_2xx":
						output.success = result;
						break;
					case "tps_3xx":
						output.redirection = result;
						break;
					case "tps_4xx":
						output.clientError = result;
						break;
					case "tps_5xx":
						output.serverError = result;
						break;
					default:
						throw new Error(`Unknown data set type: "${result.dataSet.label}"`);
				}
				return output;
			},
			({
				clientError: null,
				informational: null,
				redirection: null,
				serverError: null,
				success: null,
				total: null
			} as unknown) as TPSData
		));
	}

	/**
	 * This method is handled seperately from :js:method:`APIService.getTypes` because this information
	 * (should) never change, and therefore can be cached. This method makes an HTTP request iff the values are not already
	 * cached.
	 *
	 * @returns An Observable that will emit an array of all of the Type objects in Traffic Ops that refer specifically to Delivery Service
	 * 	types.
	 */
	public getDSTypes(): Observable<Array<Type>> {
		if (this.deliveryServiceTypes) {
			return of(this.deliveryServiceTypes);
		}
		const path = "types";
		return this.get<Array<Type>>(path, undefined, {useInTable: "deliveryservice"}).pipe(map(
			r => {
				this.deliveryServiceTypes = r;
				return r;
			}
		));
	}

	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns An Observable that will emit either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 */
	public getCDNs(id?: number): Observable<Map<string, CDN> | CDN> {
		const path = "cdns";
		if (id) {
			return this.get<[CDN]>(path, undefined, {id: String(id)}).pipe(map(
				r => r[0]
			));
		}
		return this.get<Array<CDN>>(path).pipe(map(
			r => {
				const ret = new Map<string, CDN>();
				for (const c of r) {
					ret.set(c.name, c);
				}
				return ret;
			}
		));
	}

	/**
	 * Creates a new content invalidation job.
	 *
	 * @param job The content invalidation job to be created.
	 * @returns An Observable that emits once: whether or not creation succeeded.
	 */
	public createInvalidationJob(job: InvalidationJob): Observable<boolean> {
		const path = "user/current/jobs";
		return this.post<InvalidationJob>(path, job).pipe(map(
			() => true,
			() => false
		));
	}
}
