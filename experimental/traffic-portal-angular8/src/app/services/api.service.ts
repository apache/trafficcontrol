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
import { HttpClient, HttpHeaders, HttpResponse, HttpParams } from '@angular/common/http';
import { BehaviorSubject, Observable, throwError } from 'rxjs';
import { merge } from 'rxjs/index';
import { map, mergeAll, first, catchError, reduce } from 'rxjs/operators';

import { CDN } from '../models/cdn';
import { DataPoint, DataSet, DataSetWithSummary, TPSData } from '../models/data';
import { DeliveryService } from '../models/deliveryservice';
import { InvalidationJob } from '../models/invalidation';
import { Type } from '../models/type';
import { Role, User, Capability } from '../models/user';

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

/**
 * The APIService provides access to the Traffic Ops API. Its methods should be kept API-version
 * agnostic (from the caller's perspective), and always return `Observable`s.
*/
@Injectable({ providedIn: 'root' })
export class APIService {

	/**
	 * The current API version to use
	 * @todo Get this from the environment
	 */
	public API_VERSION = '1.4';

	private deliveryServiceTypes: Array<Type>;

	// private cookies: string;

	constructor (private readonly http: HttpClient) {

	}

	private delete (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('delete', path, data);
	}
	private get (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('get', path, data);
	}
	private head (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('head', path, data);
	}
	private options (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('options', path, data);
	}
	private patch (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('patch', path, data);
	}
	private post (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('post', path, data);
	}
	private push (path: string, data?: any): Observable<HttpResponse<any>> {
		return this.do('push', path, data);
	}

	private do (method: string, path: string, data?: Object): Observable<HttpResponse<any>> {

		/* tslint:disable */
		const options = {headers: new HttpHeaders({'Content-Type': 'application/json'}),
		                 observe: 'response' as 'response',
		                 responseType: 'json' as 'json',
		                 body: data};
		/* tslint:enable */
		return this.http.request(method, path, options).pipe(map((response) => {
			// TODO pass alerts to the alert service
			// (TODO create the alert service)
			return response as HttpResponse<any>;
		}));
	}

	/**
	 * Performs authentication with the Traffic Ops server.
	 * @param u The username to be used for authentication
	 * @param p The password of user `u`
	 * @returns An observable that will emit the entire HTTP response
	*/
	public login (u: string, p: string): Observable<HttpResponse<any>> {
		const path = '/api/' + this.API_VERSION + '/user/login';
		return this.post(path, {u, p});
	}

	/**
	 * Fetches the current user from Traffic Ops
	 * @returns An observable that will emit a `User` object representing the current user.
	*/
	public getCurrentUser (): Observable<User> {
		const path = '/api/' + this.API_VERSION + '/user/current';
		return this.get(path).pipe(map(
			r => {
				return r.body.response as User;
			}
		));
	}

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
	 * Gets an array of all users in Traffic Ops
	 * @returns An Observable that will emit an Array of User objects.
	*/
	public getUsers (nameOrID?: string | number): Observable<Array<User> | User> {
		const path = '/api/' + this.API_VERSION + '/users';
		if (nameOrID) {
			switch (typeof nameOrID) {
				case 'string':
					return this.get(path + '?username=' + encodeURIComponent(nameOrID)).pipe(map(
						r => {
							for (const u of (r.body.response as Array<User>)) {
								if (u.username === nameOrID) {
									return u;
								}
							}
							return null;
						}
					));
				case 'number':
					return this.get(path + '?id=' + nameOrID.toString()).pipe(map(
						r => {
							for (const u of (r.body.response as Array<User>)) {
								if (u.id === nameOrID) {
									return u;
								}
							}
							return null;
						}
					));
				default:
					throw new TypeError("expected a username or ID, got '" + typeof(nameOrID) + "'");
					return null;
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<User>;
			}
		));
	}

	public getRoles (id: number): Observable<Role>;
	public getRoles (name: string): Observable<Role>;
	public getRoles (): Observable<Array<Role>>;
	/**
	 * Fetches one or all Roles from Traffic Ops
	 * @param name Optionally, the name of a single Role which will be fetched
	 * @param id Optionally, the integral, unique identifier of a single Role which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns an Observable that will emit either an Array of Roles, or a single Role, depending on whether
	 *	`name`/`id` was passed
	 * (In the event that `name`/`id` is given but does not match any Role, `null` will be emitted)
	*/
	public getRoles (nameOrID?: string | number) {
		const path = '/api/' + this.API_VERSION + '/roles';
		if (nameOrID) {
			switch (typeof nameOrID) {
				case 'string':
					return this.get(path + '?name=' + nameOrID).pipe(map(
						r => {
							for (const role of (r.body.response as Array<Role>)) {
								if (role.name === nameOrID) {
									return role;
								}
							}
							return null;
						}
					));
					break;
				case 'number':
					return this.get(path + '?id=' + nameOrID.toString()).pipe(map(
						r => {
							for (const role of (r.body.response as Array<Role>)) {
								if (role.id === nameOrID) {
									return role;
								}
							}
						}
					));
					break;
				default:
					throw new TypeError("expected a name or ID, got '" + typeof(nameOrID) + "'");
					break;
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<Role>;
			}
		));
	}

	public getCapabilities (name: string): Observable<Capability>;
	public getCapabilities (): Observable<Array<Capability>>;
	/**
	 * Fetches one or all Capabilities from Traffic Ops
	 * @param name Optionally, the name of a single Capability which will be fetched
	 * @throws {TypeError} When called with an improper argument.
	 * @returns an Observable that will emit either an Array of Capabilities, or a single Capability,
	 * depending on whether `name`/`id` was passed
	 * (In the event that `name`/`id` is given but does not match any Capability, `null` will be emitted)
	*/
	public getCapabilities (name?: string) {
		const path = '/api/' + this.API_VERSION + '/capabilities';
		if (name) {
			return this.get(path + '?name=' + encodeURIComponent(name)).pipe(map(
				r => {
					for (const cap of (r.body.response as Array<Capability>)) {
						if (cap.name === name) {
							return cap;
						}
					}
					return null;
				}
			));
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<Capability>;
			}
		));
	}

	/**
	 * Gets one or all Types from Traffic Ops
	 * @param name Optionally, the name of a single Type which will be returned
	 * @returns An Observable that will emit either a Map of Type names to full Type objects, or a single Type, depending on whether
	 * 	`name` was passed
	 * (In the event that `name` is given but does not match any Type, `null` will be emitted)
	*/
	public getTypes (name?: string): Observable<Map<string, Type> | Type> {
		const path = '/api/' + this.API_VERSION + '/types';
		if (name) {
			return this.get(path + '?name=' + encodeURIComponent(name)).pipe(map(
				r => {
					for (const t of (r.body.response as Array<Type>)) {
						if (t.name === name) {
							return t;
						}
					}
					return null;
				}
			));
		}
		return this.get(path).pipe(map(
			r => {
				const ret = new Map<string, Type>();
				for (const t of (r.body.response as Array<Type>)) {
					ret.set(t.name, t);
				}
				return ret;
			}
		));
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

	/**
	 * Gets one or all CDNs from Traffic Ops
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns An Observable that will emit either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any CDN, `null` will be emitted)
	*/
	public getCDNs (id?: number): Observable<Map<string, CDN> | CDN> {
		const path = '/api/' + this.API_VERSION + '/cdns';
		if (id) {
			return this.get(path + '?id=' + String(id)).pipe(map(
				r => {
					for (const c of (r.body.response as Array<CDN>)) {
						if (c.id === id) {
							return c;
						}
					}
				}
			));
		}
		return this.get(path).pipe(map(
			r => {
				const ret = new Map<string, CDN>();
				for (const c of (r.body.response as Array<CDN>)) {
					ret.set(c.name, c);
				}
				return ret;
			}
		));
	}

	public getInvalidationJobs (opts?: {id: number} |
	                                   {userId: number} |
	                                   {user: User} |
	                                   {dsId: number} |
	                                   {deliveryService: DeliveryService}): Observable<Array<InvalidationJob>> {
		let path = '/api/' + this.API_VERSION + '/jobs';
		if (opts) {
			path += '?';
			if (opts.hasOwnProperty('id')) {
				path += 'id=' + String((opts as {id: number}).id);
			} else if (opts.hasOwnProperty('dsId')) {
				path += 'dsId=' + String((opts as {dsId: number}).dsId);
			} else if (opts.hasOwnProperty('userId')) {
				path += 'userId=' + String((opts as {userId: number}).userId);
			} else if (opts.hasOwnProperty('deliveryService')) {
				path += 'dsId=' + String((opts as {deliveryService: DeliveryService}).deliveryService.id);
			} else {
				path += 'userId=' + String((opts as {user: User}).user.id);
			}
		}
		return this.get(path).pipe(map(
			r => {
				return r.body.response as Array<InvalidationJob>;
			}
		));
	}

	public createInvalidationJob (job: InvalidationJob): Observable<boolean> {
		const path = '/api/' + this.API_VERSION + '/user/current/jobs';
		return this.post(path, job).pipe(map(
			r => true,
			e => false
		));
	}
}
