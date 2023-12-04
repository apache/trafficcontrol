/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * This is a minimal type definition for Delivery Services. Expand as necessary.
 * @typedef DeliveryService
 * @property {"ACTIVE" | "INACTIVE" | "PRIMED"} active
 * @property {number} cdnId
 * @property {string[]} consistentHashQueryParams
 * @property {string[]} [exampleURLs]
 * @property {number} [id]
 * @property {string} [lastUpdated]
 * @property {string[]} requiredCapabilities
 * @property {string} routingName
 * @property {boolean} signed
 * @property {null|number} [sslKeyVersion]
 * @property {null|[string, ...string[]]} tlsVersions
 * @property {null|string} topology
 * @property {string} [trResponseHeaders]
 * @property {string|null|undefined} [type]
 * @property {number} typeId
 * @property {string} xmlId
 */

/**
 * The type of elements in the array response to GET requests from `deliveryservices_required_capabilities`.
 * @typedef DSRequiredCapability
 * @property {number} deliveryServiceID
 * @property {string} [lastUpdated]
 * @property {string} requiredCapability
 * @property {string} xmlId
 */

/**
 * The type of elements in the array response from `steering`.
 * @typedef SteeringDefinition
 * @property {string} deliveryService
 * @property {boolean} clientSteering
 * @property {{order: number; weight: number; deliveryService: string}[]} targets
 * @property {{deliveryService: string; pattern: string}[]} filters
 */

/**
 * A single target of a Delivery Service.
 * @typedef SteeringTarget
 * @property {string} deliveryService
 * @property {number} deliveryServiceId
 * @property {string} target
 * @property {number} targetId
 * @property {string} type
 * @property {number} typeId
 * @property {number} value
 */

/**
 * The result of testing a consistent hashing regular expression against Traffic
 * Router.
 * @typedef ConsistentHashResponse
 * @property {string} resultingPathToConsistentHash
 * @property {string} consistentHashRegex
 * @property {string} requestPath
 */

class DeliveryServiceService {
	/**
	 * DeliveryServiceService handles API requests dealing with Delivery Services.
	 *
	 * @param {import("angular").IHttpService} $http Angular HTTP service.
	 * @param {import("../service/utils/LocationUtils")} locationUtils Utilities for manipulating Angular routing.
	 * @param {import("../models/MessageModel")} messageModel Service for displaying messages/alerts.
	 * @param {{api:{next: string; unstable: string; stable: string}}} ENV Environment configuration.
	 */
	constructor($http, locationUtils, messageModel, ENV) {
		this.$http = $http;
		this.locationUtils = locationUtils;
		this.messageModel = messageModel;
		this.ENV = ENV;
	}

	/**
	 * Get Delivery Services.
	 *
	 * @param {Record<string, unknown>} [params] Any and all query string parameters.
	 * @returns {Promise<DeliveryService[]>} The response property of the response.
	 */
	async getDeliveryServices(params) {
		const result = await this.$http.get(`${this.ENV.api.unstable}deliveryservices`, { params });
		return result.data.response;
	}

	/**
	 * Get the Delivery Service with the given ID.
	 *
	 * @param {number} id The ID of the desired Delivery Service.
	 * @returns {Promise<DeliveryService>} The requested Delivery Service.
	 */
	async getDeliveryService(id) {
		const result = await this.$http.get(`${this.ENV.api.unstable}deliveryservices`, { params: { id } });
		return result.data.response[0];
	}

	/**
	 * Creates the given Delivery Service.
	 *
	 * @param {DeliveryService} ds The Delivery Service being created.
	 * @returns {Promise<{alerts: {level: string; text: string}[], response: DeliveryService}>} The full API response.
	 */
	async createDeliveryService(ds) {
		// strip out any falsy values or duplicates from consistentHashQueryParams
		ds.consistentHashQueryParams = Array.from(new Set(ds.consistentHashQueryParams)).filter(i => i);

		const response = await this.$http.post(`${this.ENV.api.unstable}deliveryservices`, ds);
		return response.data;
	}

	/**
	 * Replaces an existing Delivery Service with the new provided definition.
	 *
	 * @param {DeliveryService & {id: number}} ds The Delivery Service being updated (ID MUST be specified).
	 * @returns {Promise<{alerts: {level: string; text: string}[], response: DeliveryService}>} The full API response.
	 */
	async updateDeliveryService(ds) {
		// strip out any falsy values or duplicates from consistentHashQueryParams
		ds.consistentHashQueryParams = Array.from(new Set(ds.consistentHashQueryParams)).filter(i => i);

		const response = await this.$http.put(`${this.ENV.api.unstable}deliveryservices/${ds.id}`, ds);
		return response.data;
	}

	/**
	 * Deletes an existing Delivery Service.
	 *
	 * @param {DeliveryService} ds The Delivery Service to be deleted.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response
	 */
	async deleteDeliveryService(ds) {
		const response = await this.$http.delete(`${this.ENV.api.unstable}deliveryservices/${ds.id}`);
		return response.data;
	}

	/**
	 * Gets the server capabilities required by the identified Delivery Service.
	 *
	 * @param {number} deliveryServiceID The ID of the Delivery Service in question.
	 * @returns {Promise<DeliveryService[]>} The Server Capabilities required by the DS with the given ID.
	 */
	async getServerCapabilities(deliveryServiceID) {
		const result = await this.$http.get(`${this.ENV.api.unstable}deliveryservices`, { params: { deliveryServiceID } });
		return result.data.response;
	};

	/**
	 * Gets steering information.
	 * @returns {Promise<SteeringDefinition[]>}
	 */
	async getSteering() {
		const r = await this.$http.get(`${this.ENV.api.unstable}steering/`)
		return r.data.response;
	}

	/**
	 * Adds a Capability requirement to a Delivery Service.
	 * @param {number} deliveryServiceID The ID of the Delivery Service to which to add a Capability requirement.
	 * @param {string} requiredCapability The name of the Capability being added as a requirement.
	 * @returns {Promise<{alerts: {text: string; level: string}[]; response: {deliveryServiceID: number; lastUpdated: string; requiredCapability: string}}>} The full API response.
	 */
	async addServerCapability(deliveryServiceID, requiredCapability) {
		try {
			const result = await this.$http.post(`${this.ENV.api.unstable}deliveryservices_required_capabilities`, { deliveryServiceID, requiredCapability });
			return result.data;
		} catch (err) {
			if (err.data && err.data.alerts) {
				this.messageModel.setMessages(err.data.alerts, false);
			}
			throw err;
		}
	}

	/**
	 * Get the Delivery Service for which the identified server is responsible for serving content.
	 *
	 * This includes assignments through direct assignment as well as Topology ancestry.
	 *
	 * @param {number} serverID The ID of the server in question.
	 * @returns {Promise<DeliveryService[]>} The Delivery Services for which the identified server is responsible for serving content.
	 */
	async getServerDeliveryServices(serverID) {
		const result = await this.$http.get(`${this.ENV.api.unstable}servers/${serverID}/deliveryservices`);
		return result.data.response;
	}

	/**
	 * Gets the targets of the given steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the steering Delivery Service in question.
	 * @returns {Promise<SteeringTarget[]>} The targets of the identified Delivery Service.
	 */
	async getDeliveryServiceTargets(dsID) {
		const result = await this.$http.get(`${this.ENV.api.unstable}steering/${dsID}/targets`);
		return result.data.response;
	}

	/**
	 * Gets a particular target definition for a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the Steering Delivery Service in question.
	 * @param {number} target The ID of the target Delivery Service.
	 * @returns {Promise<SteeringTarget|undefined>} The definition of the requested
	 * target - or `undefined` if the target DS is not actually a target of the
	 * steering DS.
	 */
	async getDeliveryServiceTarget(dsID, target) {
		const result = await this.$http.get(`${this.ENV.api.unstable}steering/${dsID}/targets`, { params: { target } });
		return result.data.response[0];
	}

	/**
	 * Updates a particular target definition of a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the steering Delivery Service in question.
	 * @param {number} targetID The ID of the Delivery Service for which the target definition will be updated.
	 * @param {SteeringTarget} target The new, desired definition of the target.
	 * @returns {Promise<SteeringTarget>} The steering target definition after update.
	 */
	async updateDeliveryServiceTarget(dsID, targetID, target) {
		let result;
		try {
			result = await this.$http.put(`${this.ENV.api.unstable}steering/${dsID}/targets/${targetID}`, target);
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		this.messageModel.setMessages(result.data.alerts, true);
		this.locationUtils.navigateToPath(`/delivery-services/${dsID}/targets`);
		return result.data.response;
	}

	/**
	 * Creates a new target definition for a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the Steering Delivery Service in question.
	 * @param {number} target The definition of the new target.
	 * @returns {Promise<SteeringTarget>} The newly created target definition.
	 */
	async createDeliveryServiceTarget(dsID, target) {
		let result;
		try {
			result = await this.$http.post(`${this.ENV.api.unstable}steering/${dsID}/targets`, target);
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		this.messageModel.setMessages(result.data.alerts, true);
		this.locationUtils.navigateToPath(`/delivery-services/${dsID}/targets`);
		return result.data.response;
	}

	/**
	 * Removes a target from a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the Steering Delivery Service in question.
	 * @param {number} targetID The ID of the Delivery Service being removed as a target.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	async deleteDeliveryServiceTarget(dsID, targetID) {
		let result;
		try {
			result = await this.$http.delete(`${this.ENV.api.unstable}steering/${dsID}/targets/${targetID}`);
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, true);
			throw err;
		}
		this.messageModel.setMessages(result.data.alerts, true);
		this.locationUtils.navigateToPath('/delivery-services/' + dsID + '/targets');
		return result.data;
	}

	/**
	 * Removes a server from a Delivery Service's directly assigned servers.
	 *
	 * Cannot be used on a Delivery Service that uses a Topology.
	 *
	 * @param {number} dsID The ID of the Delivery Service in question.
	 * @param {number} serverID The ID of the server being removed.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	async deleteDeliveryServiceServer(dsID, serverID) {
		let result;
		try {
			result = await this.$http.delete(`${this.ENV.api.unstable}deliveryserviceserver/${dsID}/${serverID}`);
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		this.messageModel.setMessages(result.data.alerts, false);
		return result.data;
	}

	/**
	 * Assigns a set of servers to a Delivery Service *overriding any existing
	 * direct assignments*.
	 *
	 * This cannot be used with a Delivery Service that uses a Topology.
	 *
	 * @param {number} dsId The ID of the Delivery Service in question.
	 * @param {number[]} servers The IDs of the servers being directly assigned to the Delivery Service.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: {dsId: number; replace: true; servers: number[]}}>} The full API response.
	 */
	async assignDeliveryServiceServers(dsId, servers) {
		let result;
		try {
			result = await this.$http.post(`${this.ENV.api.unstable}deliveryserviceserver`, { dsId, servers, replace: true });
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		this.messageModel.setMessages(result.data.alerts, false);
		return result.data;
	}

	/**
	 * Tests a consistent hashing regular expression against Traffic Router, to
	 * see how it would be routed.
	 *
	 * @param {string|RegExp} regex The regular expression being tested.
	 * @param {string} requestPath The sample client request path.
	 * @param {number} cdnId The ID of the CDN within which the request is to be made.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: ConsistentHashResponse}>} The full API response.
	 */
	async getConsistentHashResult(regex, requestPath, cdnId) {
		const url = `${this.ENV.api.unstable}consistenthash`;
		const params = { regex, requestPath, cdnId };

		try {
			const result = await this.$http.post(url, params);
			return result.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}
}

DeliveryServiceService.$inject = ["$http", "locationUtils", "messageModel", "ENV"];
module.exports = DeliveryServiceService;
