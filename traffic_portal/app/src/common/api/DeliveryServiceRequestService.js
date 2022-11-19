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
 * This is a minimal definition of a DSR. Add to it as necessary.
 *
 * @typedef DeliveryServiceRequest
 * @property {number} id
 * @property {?import("./DeliveryServiceService").DeliveryService} requested
 */

/**
 * The allowed values for a DSR's status.
 * @typedef {"draft" | "submitted" | "rejected" | "pending" | "complete"} DSRStatus
 */

/**
 * Represents a comment on a DSR.
 *
 * @typedef DSRComment
 * @property {number} authorId
 * @property {string} author
 * @property {number} deliveryServiceRequestId
 * @property {number} id
 * @property {string} lastUpdated
 * @property {string} value
 * @property {string} xmlId
 */

/**
 * DeliveryServiceRequestService provides methods for interacting with the parts
 * of the Traffic Ops API that relate to Delivery Service Requests.
 *
 * @param {import("angular").IHttpService} $http Angular HTTP service.
 * @param {import("../models/MessageModel")} messageModel Service for displaying messages/alerts.
 * @param {{api:{next: string; unstable: string; stable: string}}} ENV Environment configuration.
 */
var DeliveryServiceRequestService = function($http, messageModel, ENV) {

	const apiVersion = ENV.api.next;

	/**
	 * Get Delivery Service Requests.
	 *
	 * @param {Record<PropertyKey, unknown>} params
	 * @returns {Promise<DeliveryServiceRequest[]>}
	 */
	this.getDeliveryServiceRequests = async function(params) {
		const result = await $http.get(`${apiVersion}/deliveryservice_requests`, {params});
		return result.data.response;
	};

	/**
	 * Creates a new DSR.
	 *
	 * @param {DeliveryServiceRequest} dsRequest The DSR to be create.
	 * @returns {Promise<DeliveryServiceRequest>} The newly created DSR.
	 */
	this.createDeliveryServiceRequest = async function(dsRequest) {

		// strip out any falsy values or duplicates from consistentHashQueryParams
		if (dsRequest.requested) {
			dsRequest.requested.consistentHashQueryParams = Array.from(new Set(dsRequest.requested.consistentHashQueryParams)).filter(i => i);
		}

		try {
			const result = await $http.post(`${apiVersion}/deliveryservice_requests`, dsRequest);
			return result.data.response;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};

	/**
	 * Replaces an existing DSR with the provided, new definition.
	 *
	 * @param {number} id The ID of the DSR being modified.
	 * @param {DeliveryServiceRequest} dsRequest The new, desired definition of the DSR.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DeliveryServiceRequest}>} The full API response.
	 */
	this.updateDeliveryServiceRequest = async function(id, dsRequest) {

		// strip out any falsy values or duplicates from consistentHashQueryParams
		if (dsRequest.requested) {
			dsRequest.requested.consistentHashQueryParams = Array.from(new Set(dsRequest.requested.consistentHashQueryParams)).filter(i => i);
		}

		try {
			const result = await $http.put(`${apiVersion}deliveryservice_requests`, dsRequest, {params: {id}});
			return result.data;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};

	/**
	 * Deletes the identified DSR.
	 *
	 * @param {number} id The ID of the DSR to be deleted.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	this.deleteDeliveryServiceRequest = async function(id) {
		try {
			const response = await $http.delete(`${apiVersion}deliveryservice_requests`, {params: {id}});
			return response.data;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};

	/**
	 * Assigns a DS to a user.
	 *
	 * @param {number} id The ID of the DSR being assigned.
	 * @param {string} assignee The username of the user to whom the DSR is being assigned.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DeliveryServiceRequest}>} The full API response.
	 */
	this.assignDeliveryServiceRequest = async function(id, assignee) {
		try {
			const result = await $http.put(`${apiVersion}deliveryservice_requests/${id}/assign`, { assignee });
			return result.data;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};

	/**
	 * Sets the status of a DSR.
	 *
	 * @param {number} id The ID of the DSR being modified.
	 * @param {DSRStatus} status The new status of the DSR.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DeliveryServiceRequest}>} The full API response.
	 */
	this.updateDeliveryServiceRequestStatus = async function(id, status) {
		try {
			const result = await $http.put(`${apiVersion}deliveryservice_requests/${id}/status`, { status });
			return result.data;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};

	/**
	 * Gets the comments associated with DSRs.
	 *
	 * @param {Record<PropertyKey, unknown>} params Any and all query string parameters for the request.
	 * @returns {Promise<DSRComment[]>} The requested comments.
	 */
	this.getDeliveryServiceRequestComments = async function(params) {
		try {
			const result = await $http.get(`${apiVersion}deliveryservice_request_comments`, {params});
			return result.data.response;
		} catch(err) {
			throw err;
		}
	};

	/**
	 * Creates a new comment on a DSR.
	 *
	 * @param {DSRComment} comment The new comment to be created.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DSRComment}>} The full API response.
	 */
	this.createDeliveryServiceRequestComment = async function(comment) {
		const response = await  $http.post(`${apiVersion}deliveryservice_request_comments`, comment);
		return response.data;
	};

	/**
	 * Updates an existing comment on a DSR.
	 *
	 * @param {DSRComment} comment The comment as desired.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DSRComment}>} The full API response.
	 */
	this.updateDeliveryServiceRequestComment = async function(comment) {
		try {
			const result = await $http.put(`${apiVersion}deliveryservice_request_comments`, comment, {params: {id: comment.id}});
			return result.data;
		} catch(err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};

	/**
	 * Removes a comment from its DSR.
	 *
	 * @param {DSRComment} comment The comment to be deleted.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	this.deleteDeliveryServiceRequestComment = async function(comment) {
		try {
			const response = await $http.delete(`${apiVersion}deliveryservice_request_comments`, {params: {id: comment.id}});
			return response.data;
		} catch (err) {
			messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	};
};

DeliveryServiceRequestService.$inject = ["$http", "messageModel", "ENV"];
module.exports = DeliveryServiceRequestService;
