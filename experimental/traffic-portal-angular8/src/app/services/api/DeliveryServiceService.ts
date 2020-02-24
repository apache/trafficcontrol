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

import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { merge } from 'rxjs/index';
import { map, mergeAll, reduce } from 'rxjs/operators';

import { APIService } from './apiservice';

import { DataPoint, DataSet, DataSetWithSummary, DeliveryService, TPSData, Type } from '../../models';

function constructDataSetFromResponse (r: any): DataSetWithSummary {
	if (!r.series || !r.series.name) {
		console.debug(r);
		throw new Error('No series data for response!');
	}

	const data = new Array<DataPoint>();
	if (r.series.values !== null && r.series.values !== undefined) {
		for (const v of r.series.values) {
			if (v[1] === null) {
				continue;
			}
			data.push({t: new Date(v[0]), y: v[1].toFixed(3)} as DataPoint);
		}
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
		min = null;
		max = null;
		fifth = null;
		nfifth = null;
		neight = null;
		mean = null;
	}

	return {
		dataSet: {label: r.series.name.split('.')[0], data: data} as DataSet,
		min: min,
		max: max,
		fifthPercentile: fifth,
		ninetyFifthPercentile: nfifth,
		ninetyEighthPercentile: neight,
		mean: mean
	} as DataSetWithSummary;
}


@Injectable({providedIn: 'root'})
export class DeliveryServiceService extends APIService {

	private deliveryServiceTypes: Array<Type>;

	/**
	 * Gets a list of all visible Delivery Services
	 * @param A unique identifier for a Delivery Service - either a numeric id or an "xml_id"
	 * @throws TypeError if ``id`` is not a proper type
	 * @returns An observable that will emit an array of `DeliveryService` objects.
	*/
	public getDeliveryServices (id?: string | number): Observable<DeliveryService[] | DeliveryService> {
		let path = '/api/' + this.API_VERSION + '/deliveryservices';
		if (id) {
			if (typeof(id) === 'string') {
				path += '?xml_id=' + encodeURIComponent(id);
			} else if (typeof(id) === 'number') {
				path += '?id=' + String(id);
			} else {
				throw new TypeError("'id' must be a string or a number! (got: '" + typeof(id) + "')");
			}
			return this.get(path).pipe(map(
				r => {
					return r.body.response[0] as DeliveryService;
				}
			));
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as DeliveryService[];
			}
		));
	}

	/**
	 * Creates a new Delivery Service
	 * @param ds The new Delivery Service object
	 * @returns An Observable that will emit a boolean value indicating the success of the operation
	*/
	public createDeliveryService (ds: DeliveryService): Observable<boolean> {
		const path = '/api/' + this.API_VERSION + '/deliveryservices';
		return this.post(path, ds).pipe(map(
			r => {
				return true;
			},
			e => {
				return false;
			}
		));
	}

	/**
	 * Retrieves capacity statistics for the Delivery Service identified by a given, unique,
	 * integral value.
	 * @param d Either a {@link DeliveryService} or an integral, unique identifier of a Delivery Service
	 * @returrns An Observable that emits an object that hopefully has the right keys to represent capacity.
	 * @throws If `d` is a {@link DeliveryService} that has no (valid) id
	*/
	public getDSCapacity (d: number | DeliveryService): Observable<any> {
		let id: number;
		if (typeof d === 'number'){
			id = d;
		} else {
			d = d as DeliveryService;
			if (!d.id || d.id < 0) {
				throw new Error("Delivery Service id must be defined!");
			}
			id = d.id;
		}

		const path = '/api/' + this.API_VERSION + '/deliveryservices/' + String(id) + '/capacity';
		return this.get(path).pipe(map(
			r => {
				return r.body.response;
			}
		));
	}

	/**
	 * Retrieves the Cache Group health of a Delivery Service identified by a given, unique,
	 * integral value.
	 * @param d The integral, unique identifier of a Delivery Service
	 * @returns An Observable that emits a response from the health endpoint
	*/
	public getDSHealth (d: number): Observable<any> {
		const path = '/api/' + this.API_VERSION + '/deliveryservices/' + String(d) + '/health';
		return this.get(path).pipe(map(
			r => {
				return r.body.response;
			}
		));
	}

	/**
	 * Retrieves Delivery Service throughput statistics for a given time period, averaged over a given
	 * interval.
	 * @param d The `xml_id` of a Delivery Service
	 * @param start A date/time from which to start data collection
	 * @param end A date/time at which to end data collection
	 * @param interval A unit-suffixed interval over which data will be "binned"
	 * @param useMids Collect data regarding Mid-tier cache servers rather than Edge-tier cache servers
	 * @param dataOnly Only returns the data series, not any supplementing meta info found in the API response
	 * @returns An Observable that will emit an Array of datapoint Arrays (length 2 containing a date string and data value)
	*//* tslint:disable */
	public getDSKBPS (d: string,
	                  start: Date,
	                  end: Date,
	                  interval: string,
	                  useMids: boolean,
	                  dataOnly?: boolean): Observable<any | Array<DataPoint>> {
		/* tslint:enable */
		let path = '/api/' + this.API_VERSION + '/deliveryservice_stats?metricType=kbps';
		path += '&interval=' + interval;
		path += '&deliveryServiceName=' + d;
		path += '&startDate=' + start.toISOString();
		path += '&endDate=' + end.toISOString();
		path += '&serverType=' + (useMids ? 'mid' : 'edge');
		return this.get(path).pipe(map(
			r => {
				if (r && r.body && r.body.response) {
					const resp = r.body.response;
					if (dataOnly) {
						if (resp.hasOwnProperty('series') && (resp.series.hasOwnProperty('values'))) {
							return resp.series.values.filter(d => d[1] !== null).map(
								d => ({
									t: new Date(d[0]),
									y: d[1].toFixed(3)
								} as DataPoint)) as Array<DataPoint>;
						}
						throw new Error("No data series found! Path was '" + path + "'");
					}
					return r.body.response;
				}
				return null;
			}
		));
	}

	/**
	 * Gets total TPS data for a Delivery Service. To get TPS data broken down by HTTP status, use {@link getAllDSTPSData}.
	 * @param d The name (xmlid) of the Delivery Service for which TPS stats will be fetched
	 * @param start The desired start date/time of the data range (must not have nonzero milliseconds!)
	 * @param end The desired end date/time of the data range (must not have nonzero milliseconds!)
	 * @param interval A string that describes the interval across which to 'bucket' data e.g. '60s'
	 * @param useMids If given (and true) will get stats for the Mid-tier instead of the Edge-tier (which is the default behavior)
	 */
	public getDSTPS (d: string,
	                 start: Date,
	                 end: Date,
	                 interval: string,
	                 useMids?: boolean): Observable<any> {
		let path = '/api/' + this.API_VERSION + '/deliveryservice_stats?metricType=tps_total';
		path += '&interval=' + interval;
		path += '&deliveryServiceName=' + d;
		path += '&startDate=' + start.toISOString();
		path += '&endDate=' + end.toISOString();
		path += '&serverType=' + (useMids ? 'mid' : 'edge');
		return this.get(path).pipe(map(
			r => {
				if (r && r.body && r.body.response) {
					return r.body.response;
				}
				return null;
			}
		));
	}

	/**
	 * Gets total TPS data for a Delivery Service, as well as TPS data by HTTP response type.
	 * @param d The name (xmlid) of the Delivery Service for which TPS stats will be fetched
	 * @param start The desired start date/time of the data range (must not have nonzero milliseconds!)
	 * @param end The desired end date/time of the data range (must not have nonzero milliseconds!)
	 * @param interval A string that describes the interval across which to 'bucket' data e.g. '60s'
	 * @param useMids If given (and true) will get stats for the Mid-tier instead of the Edge-tier (which is the default behavior)
	 */
	public getAllDSTPSData (d: string,
	                        start: Date,
	                        end: Date,
	                        interval: string,
	                        useMids?: boolean): Observable<TPSData> {
		let path = '/api/' + this.API_VERSION + '/deliveryservice_stats?';
		path += 'interval=' + interval;
		path += '&deliveryServiceName=' + d;
		path += '&startDate=' + start.toISOString();
		path += '&endDate=' + end.toISOString();
		path += '&serverType=' + (useMids ? 'mid' : 'edge');
		path += '&metricType=';
		const paths = [
			path + 'tps_total',
			path + 'tps_2xx',
			path + 'tps_3xx',
			path + 'tps_4xx',
			path + 'tps_5xx',
		];

		const observables = paths.map(x => this.get(x).pipe(map(r => constructDataSetFromResponse(r.body.response))));

		const tasks = merge(observables).pipe(mergeAll());
		return tasks.pipe(reduce(
			(output: TPSData, result: DataSetWithSummary): TPSData => {
				switch (result.dataSet.label) {
					case 'tps_total':
						output.total = result;
						break;
					case 'tps_1xx':
						output.informational = result;
						break;
					case 'tps_2xx':
						output.success = result;
						break;
					case 'tps_3xx':
						output.redirection = result;
						break;
					case 'tps_4xx':
						output.clientError = result;
						break;
					case 'tps_5xx':
						output.serverError = result;
						break;
					default:
						console.debug(result);
						throw new Error("Unknown data set type: '" + result.dataSet.label + "'");
				}
				return output;
			},
			{
				total: null,
				informational: null,
				success: null,
				redirection: null,
				clientError: null,
				serverError: null
			} as TPSData
		)) as Observable<TPSData>;
	}

	/**
	 * This method is handled seperately from :js:method:`APIService.getTypes` because this information
	 * (should) never change, and therefore can be cached. This method makes an HTTP request iff the values are not already
	 * cached.
	 * @returns An Observable that will emit an array of all of the Type objects in Traffic Ops that refer specifically to Delivery Service
	 * 	types.
	*/
	public getDSTypes (): Observable<Array<Type>> {
		if (this.deliveryServiceTypes) {
			return new Observable(
				o => {
					o.next(this.deliveryServiceTypes);
					o.complete();
					return {unsubscribe () {}};
				}
			);
		}
		const path = '/api/' + this.API_VERSION + '/types?useInTable=deliveryservice';
		return this.get(path).pipe(map(
			r => {
				this.deliveryServiceTypes = r.body.response as Array<Type>;
				return r.body.response as Array<Type>;
			}
		));
	}
}
