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
 * @property {string} changeType
 * @property {string} [createdAt]
 * @property {number} [id]
 * @property {string} [lastUpdated]
 * @property {import("./DeliveryServiceService").DeliveryService} [original]
 * @property {import("./DeliveryServiceService").DeliveryService} [requested]
 */

/**
 * Represents a comment on a DSR.
 *
 * @typedef DSRComment
 * @property {number} [authorId]
 * @property {string} [author]
 * @property {number} deliveryServiceRequestId
 * @property {number} [id]
 * @property {string} [lastUpdated]
 * @property {string} value
 * @property {string} [xmlId]
 */

/**
 * DeliveryServiceRequestService provides methods for interacting with the parts
 * of the Traffic Ops API that relate to Delivery Service Requests.
 */
class DeliveryServiceRequestService {

	/**
	 * @type {string}
	 * @readonly
	 * @private
	 */
	apiVersion;

	/**
	 *
	 * @param {import("angular").IHttpService} $http Angular HTTP service.
	 * @param {import("../models/MessageModel")} messageModel Service for displaying messages/alerts.
	 * @param {{api:{next: string; unstable: string; stable: string}}} ENV Environment configuration.
	 */
	constructor($http, messageModel, ENV) {
		this.apiVersion = ENV.api.next;
		this.$http = $http;
		this.messageModel = messageModel;
	}

	/**
	 * Get Delivery Service Requests.
	 *
	 * @param {Record<PropertyKey, unknown>} params
	 * @returns {Promise<DeliveryServiceRequest[]>}
	 */
	async getDeliveryServiceRequests(params) {
		const result = await this.$http.get(`${this.apiVersion}deliveryservice_requests`, { params });
		return result.data.response;
	}

	/**
	 * Creates a new DSR.
	 *
	 * @param {DeliveryServiceRequest} dsRequest The DSR to be create.
	 * @returns {Promise<DeliveryServiceRequest & {id: number}>} The newly created DSR.
	 */
	async createDeliveryServiceRequest(dsRequest) {

		// strip out any falsy values or duplicates from consistentHashQueryParams
		if (dsRequest.requested) {
			dsRequest.requested.consistentHashQueryParams = Array.from(new Set(dsRequest.requested.consistentHashQueryParams)).filter(i => i);
		}

		try {
			const result = await this.$http.post(`${this.apiVersion}deliveryservice_requests`, dsRequest);
			return result.data.response;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}

	/**
	 * Replaces an existing DSR with the provided, new definition.
	 *
	 * @param {number} id The ID of the DSR being modified.
	 * @param {DeliveryServiceRequest} dsRequest The new, desired definition of the DSR.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DeliveryServiceRequest}>} The full API response.
	 */
	async updateDeliveryServiceRequest(id, dsRequest) {

		// strip out any falsy values or duplicates from consistentHashQueryParams
		if (dsRequest.requested) {
			dsRequest.requested.consistentHashQueryParams = Array.from(new Set(dsRequest.requested.consistentHashQueryParams)).filter(i => i);
		}

		try {
			const result = await this.$http.put(`${this.apiVersion}deliveryservice_requests`, dsRequest, { params: { id } });
			return result.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}

	/**
	 * Deletes the identified DSR.
	 *
	 * @param {number} id The ID of the DSR to be deleted.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	async deleteDeliveryServiceRequest(id) {
		try {
			const response = await this.$http.delete(`${this.apiVersion}deliveryservice_requests`, { params: { id } });
			return response.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}

	/**
	 * Assigns a DS to a user.
	 *
	 * @param {number} id The ID of the DSR being assigned.
	 * @param {string} assignee The username of the user to whom the DSR is being assigned.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DeliveryServiceRequest}>} The full API response.
	 */
	async assignDeliveryServiceRequest(id, assignee) {
		try {
			const result = await this.$http.put(`${this.apiVersion}deliveryservice_requests/${id}/assign`, { assignee });
			return result.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}

	/**
	 * Sets the status of a DSR.
	 *
	 * @param {number} id The ID of the DSR being modified.
	 * @param {string} status The new status of the DSR.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DeliveryServiceRequest}>} The full API response.
	 */
	async updateDeliveryServiceRequestStatus(id, status) {
		try {
			const result = await this.$http.put(`${this.apiVersion}deliveryservice_requests/${id}/status`, { status });
			return result.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}

	/**
	 * Gets the comments associated with DSRs.
	 *
	 * @param {Record<PropertyKey, unknown>} params Any and all query string parameters for the request.
	 * @returns {Promise<DSRComment[]>} The requested comments.
	 */
	async getDeliveryServiceRequestComments(params) {
		try {
			const result = await this.$http.get(`${this.apiVersion}deliveryservice_request_comments`, { params });
			return result.data.response;
		} catch (err) {
			throw err;
		}
	}

	/**
	 * Creates a new comment on a DSR.
	 *
	 * @param {DSRComment} comment The new comment to be created.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DSRComment}>} The full API response.
	 */
	async createDeliveryServiceRequestComment(comment) {
		const response = await this.$http.post(`${this.apiVersion}deliveryservice_request_comments`, comment);
		return response.data;
	}

	/**
	 * Updates an existing comment on a DSR.
	 *
	 * @param {DSRComment} comment The comment as desired.
	 * @returns {Promise<{alerts: {level: string; text: string}[]; response: DSRComment}>} The full API response.
	 */
	async updateDeliveryServiceRequestComment(comment) {
		try {
			const result = await this.$http.put(`${this.apiVersion}deliveryservice_request_comments`, comment, { params: { id: comment.id } });
			return result.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}

	/**
	 * Removes a comment from its DSR.
	 *
	 * @param {DSRComment} comment The comment to be deleted.
	 * @returns {Promise<{alerts: {level: string; text: string}[]}>} The full API response.
	 */
	async deleteDeliveryServiceRequestComment(comment) {
		try {
			const response = await this.$http.delete(`${this.apiVersion}deliveryservice_request_comments`, { params: { id: comment.id } });
			return response.data;
		} catch (err) {
			this.messageModel.setMessages(err.data.alerts, false);
			throw err;
		}
	}
}

DeliveryServiceRequestService.$inject = ["$http", "messageModel", "ENV"];
module.exports = DeliveryServiceRequestService;
