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
 * @property {string[]} consistentHashQueryParams
 * @property {?number|null|undefined} id
 * @property {null|[string, ...string[]]} tlsVersions
 */

/**
 * The type of elements in the array response to GET requests from `deliveryservices_required_capabilities`.
 * @typedef DSRequiredCapability
 * @property {number} deliveryServiceID
 * @property {string} lastUpdated
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

/**
 * DeliveryServiceService handles API requests dealing with Delivery Services.
 *
 * @param {import("angular").IHttpService} $http Angular HTTP service.
 * @param {import("../service/utils/LocationUtils")} locationUtils Utilities for manipulating Angular routing.
 * @param {import("../models/MessageModel")} messageModel Service for displaying messages/alerts.
 * @param {{api:{next: string; unstable: string; stable: string}}} ENV Environment configuration.
 */
var DeliveryServiceService = function($http, locationUtils, messageModel, ENV) {

	/**
	 * Get Delivery Services.
	 *
	 * @param {Record<string, unknown>} params Any and all query string parameters.
	 * @returns {Promise<DeliveryService[]>} The response property of the response.
	 */
	this.getDeliveryServices = async function(params) {
		const result = await $http.get(`${ENV.api.next}deliveryservices`, {params});
		return result.data.response;
	};

	/**
	 * Get the Delivery Service with the given ID.
	 *
	 * @param {number} id The ID of the desired Delivery Service.
	 * @returns {Promise<DeliveryService>} The requested Delivery Service.
	 */
	this.getDeliveryService = async function(id) {
		const result = await $http.get(`${ENV.api.next}deliveryservices`, {params: {id}});
		return result.data.response[0];
	};

	/**
	 * Creates the given Delivery Service.
	 *
	 * @param {DeliveryService} ds The Delivery Service being created.
	 * @returns {Promise<{alerts: {level: string; text: string}[], response: DeliveryService}>} The full API response.
	 */
	this.createDeliveryService = async function(ds) {
		// strip out any falsy values or duplicates from consistentHashQueryParams
		ds.consistentHashQueryParams = Array.from(new Set(ds.consistentHashQueryParams)).filter(i => i);

		const response = await $http.post(`${ENV.api.next}deliveryservices`, ds);
		return response.data;
	};

	/**
	 * Replaces an existing Delivery Service with the new provided definition.
	 *
	 * @param {DeliveryService & {id: number}} ds The Delivery Service being updated (ID MUST be specified).
	 * @returns {Promise<{alerts: {level: string; text: string}[], response: DeliveryService}>} The full API response.
	 */
	this.updateDeliveryService = async function(ds) {
		// strip out any falsy values or duplicates from consistentHashQueryParams
		ds.consistentHashQueryParams = Array.from(new Set(ds.consistentHashQueryParams)).filter(i => i);

		const response = await $http.put(`${ENV.api.next}deliveryservices/${ds.id}`, ds);
		return response.data;
	};

	/**
	 * Deletes an existing Delivery Service.
	 *
	 * @param {DeliveryService} ds The Delivery Service to be deleted.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response
	 */
	this.deleteDeliveryService = async function(ds) {
		const response = await $http.delete(`${ENV.api.next}deliveryservices/${ds.id}`);
		return response.data;
	};

	/**
	 * Gets the server capabilities required by the identified Delivery Service.
	 *
	 * @param {number} deliveryServiceID The ID of the Delivery Service in question.
	 * @returns {Promise<DSRequiredCapability[]>} The Server Capabilities required by the DS with the given ID.
	 */
	this.getServerCapabilities = async function(deliveryServiceID) {
		const result = await $http.get(`${ENV.api.unstable}deliveryservices_required_capabilities`, { params: { deliveryServiceID } });
		return result.data.response;
	};

	/**
	 * Gets steering information.
	 * @returns {Promise<SteeringDefinition[]>}
	 */
	this.getSteering = async () => $http.get(`${ENV.api.unstable}steering/`).then(r => r.data.response);

	/**
	 * Adds a Capability requirement to a Delivery Service.
	 * @param {number} deliveryServiceID The ID of the Delivery Service to which to add a Capability requirement.
	 * @param {string} requiredCapability The name of the Capability being added as a requirement.
	 * @returns {Promise<{alerts: {text: string; level: string}[]; response: {deliveryServiceID: number; lastUpdated: string; requiredCapability: string}}>} The full API response.
	 */
	this.addServerCapability = async function(deliveryServiceID, requiredCapability) {
		try {
			const result = await $http.post(`${ENV.api.unstable}deliveryservices_required_capabilities`, { deliveryServiceID, requiredCapability});
			return result.data;
		} catch (err) {
			if (err.data && err.data.alerts) {
				messageModel.setMessages(err.data.alerts, false);
			}
			throw err;
		}
	};

	/**
	 * Removes the requirement of a particular Capability from the identified Delivery Service.
	 *
	 * @param {number} deliveryServiceID The ID of the Delivery Service from which a Capability requirement will be removed.
	 * @param {string} requiredCapability The name of the Capability being removed as a requirement.
	 * @returns {Promise<{alerts: {text: string; level: string}[]}>} The full API response.
	 */
	this.removeServerCapability = async function(deliveryServiceID, requiredCapability) {
		try {
			const result = await $http.delete(`${ENV.api.unstable}deliveryservices_required_capabilities`, { params: { deliveryServiceID, requiredCapability} });
			return result.data;
		} catch(err) {
			if (err.data && err.data.alerts) {
				messageModel.setMessages(err.data.alerts, false);
			}
			throw err;
		}
	};

	/**
	 * Get the Delivery Service for which the identified server is responsible for serving content.
	 *
	 * This includes assignments through direct assignment as well as Topology ancestry.
	 *
	 * @param {number} serverID The ID of the server in question.
	 * @returns {Promise<DeliveryService[]>} The Delivery Services for which the identified server is responsible for serving content.
	 */
	this.getServerDeliveryServices = async function(serverID) {
		const result = await $http.get(`${ENV.api.unstable}servers/${serverID}/deliveryservices`);
		return result.data.response;
	};

	/**
	 * Gets the targets of the given steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the steering Delivery Service in question.
	 * @returns {Promise<SteeringTarget[]>} The targets of the identified Delivery Service.
	 */
	this.getDeliveryServiceTargets = async function(dsID) {
		const result = await $http.get(`${ENV.api.unstable}steering/${dsID}/targets`);
		return result.data.response;
	};

	/**
	 * Gets a particular target definition for a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the Steering Delivery Service in question.
	 * @param {number} target The ID of the target Delivery Service.
	 * @returns {Promise<SteeringTarget|undefined>} The definition of the requested
	 * target - or `undefined` if the target DS is not actually a target of the
	 * steering DS.
	 */
	this.getDeliveryServiceTarget = async function(dsID, target) {
		const result = await $http.get(`${ENV.api.unstable}steering/${dsID}/targets`, {params: {target}});
		return result.data.response[0];
	};

	/**
	 * Updates a particular target definition of a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the steering Delivery Service in question.
	 * @param {number} targetID The ID of the Delivery Service for which the target definition will be updated.
	 * @param {SteeringTarget} target The new, desired definition of the target.
	 * @returns {Promise<SteeringTarget>} The steering target definition after update.
	 */
	this.updateDeliveryServiceTarget = async function(dsID, targetID, target) {
		let result;
		try {
			result = await $http.put(`${ENV.api.unstable}steering/${dsID}/targets/${targetID}`, target);
		} catch (err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		messageModel.setMessages(result.data.alerts, true);
		locationUtils.navigateToPath(`/delivery-services/${dsID}/targets`);
		return result.data.response;
	};

	/**
	 * Creates a new target definition for a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the Steering Delivery Service in question.
	 * @param {number} target The definition of the new target.
	 * @returns {Promise<SteeringTarget>} The newly created target definition.
	 */
	this.createDeliveryServiceTarget = async function(dsID, target) {
		let result;
		try {
			result = await $http.post(`${ENV.api.unstable}steering/${dsID}/targets`, target);
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		messageModel.setMessages(result.data.alerts, true);
		locationUtils.navigateToPath(`/delivery-services/${dsID}/targets`);
		return result.data.response;
	};

	/**
	 * Removes a target from a Steering Delivery Service.
	 *
	 * @param {number} dsID The ID of the Steering Delivery Service in question.
	 * @param {number} targetID The ID of the Delivery Service being removed as a target.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	this.deleteDeliveryServiceTarget = async function(dsID, targetID) {
		let result;
		try {
			result = await $http.delete(`${ENV.api.unstable}steering/${dsID}/targets/${targetID}`);
		} catch (err) {
			messageModel.setMessages(err.data.alerts, true);
			throw err;
		}
		messageModel.setMessages(result.data.alerts, true);
		locationUtils.navigateToPath('/delivery-services/' + dsID + '/targets');
		return result.data;
	};

	/**
	 * Removes a server from a Delivery Service's directly assigned servers.
	 *
	 * Cannot be used on a Delivery Service that uses a Topology.
	 *
	 * @param {number} dsID The ID of the Delivery Service in question.
	 * @param {number} serverID The ID of the server being removed.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	this.deleteDeliveryServiceServer = async function(dsID, serverID) {
		let result;
		try {
			result = await $http.delete(`${ENV.api.unstable}deliveryserviceserver/${dsID}/${serverID}`);
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		messageModel.setMessages(result.data.alerts, false);
		return result.data;
	};

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
	this.assignDeliveryServiceServers = async function(dsId, servers) {
		let result;
		try {
			result = await $http.post(`${ENV.api.unstable}deliveryserviceserver`, { dsId, servers, replace: true } );
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
		messageModel.setMessages(result.data.alerts, false);
		return result.data;
	};

	/**
	 * Tests a consistent hashing regular expression against Traffic Router, to
	 * see how it would be routed.
	 *
	 * @param {string|RegExp} regex The regular expression being tested.
	 * @param {string} requestPath The sample client request path.
	 * @param {number} cdnId The ID of the CDN within which the request is to be made.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: ConsistentHashResponse}>} The full API response.
	 */
	this.getConsistentHashResult = async function (regex, requestPath, cdnId) {
		const url = `${ENV.api.unstable}consistenthash`;
		const params = {regex, requestPath, cdnId};

		try {
			const result = await $http.post(url, params);
			return result.data;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};
};

DeliveryServiceService.$inject = ["$http", "locationUtils", "messageModel", "ENV"];
module.exports = DeliveryServiceService;
