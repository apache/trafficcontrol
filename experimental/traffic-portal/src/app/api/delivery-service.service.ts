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
import type {
	Capacity,
	Health,
	DSStats,
	RequestDeliveryService,
	ResponseDeliveryService,
	SteeringConfiguration,
	TypeFromResponse,
	ResponseDeliveryServiceSSLKey
} from "trafficops-types";

import type {
	DataPoint,
	DataSetWithSummary,
	TPSData,
} from "src/app/models";

import { LoggingService } from "../shared/logging.service";

import { APIService } from "./base-api.service";

/**
 * Generates a default, blank data set with the given label.
 *
 * @param label The dataset's label.
 * @returns A dataset with no data points and all metrics set to `0`.
 */
function defaultDataSet(label: string): DataSetWithSummary {
	return {
		dataSet: {
			data: [],
			label,
		},
		fifthPercentile: 0,
		max: 0,
		mean: 0,
		min: 0,
		ninetyEighthPercentile: 0,
		ninetyFifthPercentile: 0,
	};
}

/**
 * Constructs a data set from the API response.
 *
 * @param r The parsed response body.
 * @returns A DataSetWithSummary that was massaged out of the raw response.
 */
export function constructDataSetFromResponse(r: DSStats): DataSetWithSummary {
	if (!r.series) {
		// logging service not accessible in this scope
		// eslint-disable-next-line no-console
		console.debug("raw DS stats response:", r);
		throw new Error("invalid data set response");
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
		mean = r.summary.average;
	} else {
		min = -1;
		max = -1;
		fifth = -1;
		nfifth = -1;
		neight = -1;
		mean = -1;
	}

	return {
		dataSet: {data, label: r.series.name.split(".")[0]},
		fifthPercentile: fifth,
		max,
		mean,
		min,
		ninetyEighthPercentile: neight,
		ninetyFifthPercentile: nfifth
	};
}

/**
 * DeliveryServiceService exposes API functionality related to Delivery Services.
 */
@Injectable()
export class DeliveryServiceService extends APIService {

	/** This is where DS Types are cached, as they are presumed to not change (often). */
	private deliveryServiceTypes: Array<TypeFromResponse>;

	constructor(http: HttpClient, private readonly log: LoggingService) {
		super(http);
		this.deliveryServiceTypes = new Array<TypeFromResponse>();
	}

	/**
	 * Gets a list of all Steering Configurations.
	 *
	 * @returns An array of Steering Configurations.
	 */
	public async getSteering(): Promise<Array<SteeringConfiguration>> {
		const path = "steering";
		return this.get<Array<SteeringConfiguration>>(path).toPromise();
	}

	/**
	 * Get a single Delivery Service.
	 *
	 * @param id A unique identifier for a Delivery Service - either its numeric
	 * ID or its "XML_ID".
	 * @returns the Delivery Service with the given identifier.
	 */
	public async getDeliveryServices(id: string | number): Promise<ResponseDeliveryService>;
	/**
	 * Gets a list of all visible Delivery Services
	 *
	 * @returns An array of {@link ResponseDeliveryService} objects.
	 */
	public async getDeliveryServices(): Promise<Array<ResponseDeliveryService>>;
	/**
	 * Gets a list of all visible Delivery Services
	 *
	 * @param id A unique identifier for a single Delivery Service to fetch.
	 * This may be either its numeric ID or its "XML_ID".
	 * @throws {Error} If no DS with a given `id` is found, or if more than one
	 * is found.
	 * @returns One or more Delivery Services.
	 */
	public async getDeliveryServices(id?: string | number): Promise<ResponseDeliveryService[] | ResponseDeliveryService> {
		const path = "deliveryservices";
		if (id) {
			let params;
			switch (typeof id) {
				case "string":
					params = {xmlId: id};
					break;
				case "number":
					params = { id };
			}
			const r = await this.get<[ResponseDeliveryService]>(path, undefined, params).toPromise();
			if (r.length !== 1) {
				throw new Error(`expected exactly one Delivery Service by identifier '${id}', got: ${r.length}`);
			}
			return r[0];
		}
		return this.get<Array<ResponseDeliveryService>>(path).toPromise();
	}

	/**
	 * Creates a new Delivery Service
	 *
	 * @param ds The new Delivery Service object
	 * @returns A boolean value indicating the success of the operation
	 */
	public async createDeliveryService(ds: RequestDeliveryService): Promise<ResponseDeliveryService> {
		const path = "deliveryservices";
		return this.post<ResponseDeliveryService>(path, ds).toPromise();
	}

	/**
	 * Retrieves capacity statistics for the Delivery Service identified by a given, unique,
	 * integral value.
	 *
	 * @param d Either a {@link DeliveryService} or an integral, unique identifier of a Delivery Service
	 * @returns An object that hopefully has the right keys to represent capacity.
	 * @throws If `d` is a {@link DeliveryService} that has no (valid) id
	 */
	public async getDSCapacity(d: number | ResponseDeliveryService): Promise<Capacity> {
		const id = typeof(d) === "number" ? d : d.id;
		return this.get<Capacity>(`deliveryservices/${id}/capacity`).toPromise();
	}

	/**
	 * Retrieves the Cache Group health of a Delivery Service identified by a given, unique,
	 * integral value.
	 *
	 * @param d The integral, unique identifier of a Delivery Service
	 * @returns A response from the health endpoint
	 */
	public async getDSHealth(d: number | ResponseDeliveryService): Promise<Health> {
		const id = typeof(d) === "number" ? d : d.id;
		return this.get<Health>(`deliveryservices/${id}/health`).toPromise();
	}

	/**
	 * Retrieves Delivery Service throughput statistics for a given time period,
	 * averaged over a given interval.
	 *
	 * @param d The `xml_id` of a Delivery Service, or the Delivery Service
	 * itself.
	 * @param start A date/time from which to start data collection.
	 * @param end A date/time at which to end data collection.
	 * @param interval A unit-suffixed interval over which data will be
	 * "binned".
	 * @param useMids Collect data regarding Mid-tier cache servers rather than
	 * Edge-tier cache servers
	 * @param dataOnly Only returns the data series, not any supplementing meta
	 * info found in the API response.
	 * @returns An Array of {@link DataPoint}s.
	 */
	public async getDSKBPS(
		d: string | ResponseDeliveryService,
		start: Date,
		end: Date,
		interval: string,
		useMids: boolean,
		dataOnly: true
	): Promise<Array<DataPoint>>;
	/**
	 * Retrieves Delivery Service throughput statistics for a given time period,
	 * averaged over a given interval.
	 *
	 * @param d The `xml_id` of a Delivery Service, or the Delivery Service
	 * itself.
	 * @param start A date/time from which to start data collection.
	 * @param end A date/time at which to end data collection.
	 * @param interval A unit-suffixed interval over which data will be
	 * "binned".
	 * @param useMids Collect data regarding Mid-tier cache servers rather than
	 * Edge-tier cache servers
	 * @returns The API response.
	 */
	public async getDSKBPS(
		d: string | ResponseDeliveryService,
		start: Date,
		end: Date,
		interval: string,
		useMids: boolean
	): Promise<DSStats>;
	/**
	 * Retrieves Delivery Service throughput statistics for a given time period,
	 * averaged over a given interval.
	 *
	 * @param d The `xml_id` of a Delivery Service, or the Delivery Service
	 * itself.
	 * @param start A date/time from which to start data collection.
	 * @param end A date/time at which to end data collection.
	 * @param interval A unit-suffixed interval over which data will be
	 * "binned".
	 * @param useMids Collect data regarding Mid-tier cache servers rather than
	 * Edge-tier cache servers
	 * @param dataOnly Only returns the data series, not any supplementing meta
	 * info found in the API response.
	 * @returns An Array of {@link DataPoint}s if only data was requested, or
	 * the entire API response otherwise.
	 */
	public async getDSKBPS(
		d: string | ResponseDeliveryService,
		start: Date,
		end: Date,
		interval: string,
		useMids: boolean,
		dataOnly?: boolean
	): Promise<Array<DataPoint> | DSStats> {
		const path = "deliveryservice_stats";
		const params = {
			deliveryServiceName: typeof(d) === "string" ? d : d.xmlId,
			endDate: end.toISOString(),
			interval,
			metricType: "kbps",
			serverType: useMids ? "mid" : "edge",
			startDate: start.toISOString()
		};

		if (dataOnly) {
			const r = await this.get<DSStats>(path, undefined, {exclude: "summary", ...params}).toPromise();
			if (r.series && r.series.values) {
				const series = [];
				for (const [t, y] of r.series.values) {
					if (y !== null) {
						series.push({
							t,
							y: y.toFixed(3)
						});
					}
				}
				return series;
			}
			this.log.debug("data response:", r);
			throw new Error("no data series found");
		}
		return this.get<DSStats>(path, undefined, params).toPromise();
	}

	/**
	 * Gets total TPS data for a Delivery Service. To get TPS data broken down
	 * by HTTP status, use {@link getAllDSTPSData}.
	 *
	 * @param d The Delivery Service or its "XML_ID" for which TPS stats will be
	 * fetched.
	 * @param start The desired start date/time of the data range (must not have
	 * nonzero milliseconds!).
	 * @param end The desired end date/time of the data range (must not have
	 * nonzero milliseconds!).
	 * @param interval A string that describes the interval across which to
	 * 'bucket' data e.g. '60s'.
	 * @param useMids If given (and true) will get stats for the Mid-tier
	 * instead of the Edge-tier (which is the default behavior).
	 * @returns The requested stats data.
	 */
	public async getDSTPS(
		d: string | ResponseDeliveryService,
		start: Date,
		end: Date,
		interval: string,
		useMids?: boolean
	): Promise<DSStats> {
		const path = "deliveryservice_stats";
		const params = {
			deliveryServiceName: typeof(d) === "string" ? d : d.xmlId,
			endDate: end.toISOString(),
			interval,
			metricType: "tps_total",
			serverType: useMids ? "mid" : "edge",
			startDate: start.toISOString()
		};
		return this.get<DSStats>(path, undefined, params).toPromise();
	}

	/**
	 * Gets total TPS data for a Delivery Service, as well as TPS data by HTTP
	 * response type.
	 *
	 * @param d The Delivery Service or its "XML_ID" for which TPS stats will be
	 * fetched.
	 * @param start The desired start date/time of the data range (must not have
	 * nonzero milliseconds!).
	 * @param end The desired end date/time of the data range (must not have
	 * nonzero milliseconds!).
	 * @param interval A string that describes the interval across which to
	 * 'bucket' data e.g. '60s'.
	 * @param useMids If given (and true) will get stats for the Mid-tier
	 * instead of the Edge-tier (which is the default behavior).
	 * @returns The requested TPSData.
	 */
	public async getAllDSTPSData(
		d: string | ResponseDeliveryService,
		start: Date,
		end: Date,
		interval: string,
		useMids?: boolean
	): Promise<TPSData> {
		const path = "deliveryservice_stats";
		const params: Record<string, string> = {
			deliveryServiceName: typeof(d) === "string" ? d : d.xmlId,
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
			async x => this.get<DSStats>(path, undefined, {metricType: x, ...params}).toPromise().then(constructDataSetFromResponse)
		);

		const data = await Promise.all(observables);
		const output: TPSData = {
			clientError: defaultDataSet("tps_4xx"),
			redirection: defaultDataSet("tps_3xx"),
			serverError: defaultDataSet("tps_5xx"),
			success: defaultDataSet("tps_2xx"),
			total: defaultDataSet("tps_total"),
		};
		for (const dataSet of data) {
			switch (dataSet.dataSet.label) {
				case "tps_total":
					output.total = dataSet;
					break;
				case "tps_2xx":
					output.success = dataSet;
					break;
				case "tps_3xx":
					output.redirection = dataSet;
					break;
				case "tps_4xx":
					output.clientError = dataSet;
					break;
				case "tps_5xx":
					output.serverError = dataSet;
					break;
				default:
					throw new Error(`Unknown data set type: "${dataSet.dataSet.label}"`);
			}
		}
		return output;
	}

	/**
	 * This method is handled separately from `TypeService.getTypes`
	 * because this information (should) never change, and therefore can be
	 * cached. This method makes an HTTP request iff the values are not already
	 * cached.
	 *
	 * @returns An array of all of the Type objects in Traffic Ops that refer
	 * specifically to Delivery Service types.
	 */
	public async getDSTypes(): Promise<Array<TypeFromResponse>> {
		if (this.deliveryServiceTypes.length > 0) {
			return this.deliveryServiceTypes;
		}
		const path = "types";
		const r = await this.get<Array<TypeFromResponse>>(path, undefined, {useInTable: "deliveryservice"}).toPromise();
		this.deliveryServiceTypes = r;
		return r;
	}

	/**
	 * Gets a Delivery Service's SSL Keys
	 *
	 * @param ds The delivery service xmlid or object
	 * @returns The DS ssl keys
	 */
	public async getSSLKeys(ds: string | ResponseDeliveryService): Promise<ResponseDeliveryServiceSSLKey> {
		const xmlId = typeof ds === "string" ? ds : ds.xmlId;
		const path = `deliveryservices/xmlId/${xmlId}/sslkeys`;
		return this.get<ResponseDeliveryServiceSSLKey>(path, undefined, {decode: true}).toPromise();
	}
}
