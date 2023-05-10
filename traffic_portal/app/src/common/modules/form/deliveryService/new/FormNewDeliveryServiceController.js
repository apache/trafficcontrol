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
 * This is the controller for a form used to create a new Delivery Service.
 *
 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {unknown} origin
 * @param {unknown[]} topologies
 * @param {string} type
 * @param {unknown[]} types
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {{open: ({}) => {result: Promise<*>}}} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../../api/DeliveryServiceRequestService")} deliveryServiceRequestService
 * @param {import("../../../../models/MessageModel")} messageModel
 * @param {import("../../../../models/PropertiesModel")} propertiesModel
 * @param {import("../../../../models/UserModel")} userModel
 */
var FormNewDeliveryServiceController = function(deliveryService, origin, topologies, type, types, $scope, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService, deliveryServiceRequestService, messageModel, propertiesModel, userModel) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller("FormDeliveryServiceController", { deliveryService: deliveryService,
		dsCurrent: deliveryService, origin: origin, topologies: topologies, type: type, types: types, $scope: $scope }));

	$scope.deliveryServiceName = "New";

	$scope.settings = {
		isNew: true,
		isRequest: false,
		saveLabel: "Create"
	};

	$scope.restrictTLS = false;

	/**
	 * Creates a new request to create a Delivery Service.
	 *
	 * @param {import("../../../../api/DeliveryServiceRequestService").DeliveryServiceRequest} dsRequest The creation DSR being created.
	 * @param {string} dsRequestComment The initial comment to put on the DSR.
	 * @param {boolean} autoFulfilled If `true`, the request is immediately marked as complete.
	 * @returns {Promise<void>}
	 */
	async function createDeliveryServiceCreateRequest(dsRequest, dsRequestComment, autoFulfilled) {
		const response = await deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);
		const comment = {
			deliveryServiceRequestId: response.id,
			value: dsRequestComment
		};

		await deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
		if (!autoFulfilled) {
			const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original?.xmlId;
			messageModel.setMessages([ { level: "success", text: `Created request to ${dsRequest.changeType} the ${xmlId} delivery service` } ], true);
			locationUtils.navigateToPath("/delivery-service-requests");
			return;
		}

		await Promise.all([
			// assign the ds request
			deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username),
			// set the status to "complete"
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, "complete")
		]);
	};

	/**
	 * Creates a new Delivery Service using the DSR system.
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService The DS being created.
	 * @returns {Promise<void>}
	 */
	async function createDSUsingRequest(deliveryService) {
		const params = {
			title: "Delivery Service Create Request",
			message: "All new delivery services must be reviewed for completeness and accuracy before deployment."
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/deliveryServiceRequest/dialog.deliveryServiceRequest.tpl.html",
			controller: "DialogDeliveryServiceRequestController",
			size: "md",
			resolve: {
				params,
				statuses: () => {
					const statuses = [
						{ id: $scope.DRAFT, name: "Save Request as Draft" },
						{ id: $scope.SUBMITTED, name: "Submit Request for Review and Deployment" }
					];
					if (userModel.user.role === propertiesModel.properties.dsRequests.overrideRole) {
						statuses.push({ id: $scope.COMPLETE, name: "Fulfill Request Immediately" });
					}
					return statuses;
				}
			}
		});

		/** @type {{status: {id: number}; comment: string}} */
		let options;
		try {
			options = await modalInstance.result;
		} catch {
			// This means the user cancelled
			return;
		}

		let status = "draft";
		if (options.status.id == $scope.SUBMITTED || options.status.id == $scope.COMPLETE) {
			status = "submitted";
		};
		const dsRequest = {
			changeType: "create",
			status: status,
			requested: deliveryService
		};
		// if the user chooses to complete/fulfill the create request immediately, the ds will be created and behind the
		// scenes a delivery service request will be created and marked as complete
		if (options.status.id === $scope.COMPLETE) {
			try {
				const result = await deliveryServiceService.createDeliveryService(deliveryService);
				await createDeliveryServiceCreateRequest(dsRequest, options.comment, true);
				messageModel.setMessages(result.alerts, true);
				locationUtils.navigateToPath(`/delivery-services/${result.response.id}?dsType=${result.response.type}`);
			} catch(fault) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages(fault.data.alerts, false);
			}
			return;
		}
		return createDeliveryServiceCreateRequest(dsRequest, options.comment, false);
	}

	/**
	 * Handles the user clicking the "Create" button.
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
	 * @returns {Promise<void>}
	 */
	$scope.save = async function(deliveryService) {
		deliveryService.requiredCapabilities = Object.entries($scope.selectedCapabilities).filter(sc => (sc[1])).map(sc => sc[0])
		$scope.loadGeoLimitCountries(deliveryService);

		if (!$scope.restrictTLS) {
			deliveryService.tlsVersions = null;
		}
		// if ds requests are enabled in traffic_portal_properties.json, we'll create a ds request, else just create the ds
		if ($scope.dsRequestsEnabled) {
			await createDSUsingRequest(deliveryService);
			return;
		}
		try {
			const result = await deliveryServiceService.createDeliveryService(deliveryService);
			messageModel.setMessages(result.alerts, true);
			locationUtils.navigateToPath(`/delivery-services/${result.response.id}?dsType=${result.response.type}`);
		} catch(fault) {
			$anchorScroll(); // scrolls window to top
			messageModel.setMessages(fault.data.alerts, false);
		}
	};
};

FormNewDeliveryServiceController.$inject = ["deliveryService", "origin", "topologies", "type", "types", "$scope", "$controller", "$uibModal", "$anchorScroll", "locationUtils", "deliveryServiceService", "deliveryServiceRequestService", "messageModel", "propertiesModel", "userModel"];
module.exports = FormNewDeliveryServiceController;
