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
 * This is an extension of the general DS form that is meant to handle editing
 * an existing DS - as apposed to creating a new one.
 *
 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {unknown} origin
 * @param {unknown[]} topologies
 * @param {string} type
 * @param {unknown[]} types
 * @param {import("angular").IScope & Record<PropertyKey, any>} $scope
 * @param {{reload: ()=>void}} $state
 * @param {import("angular").IControllerService} $controller
 * @param {{open: ({})=>{result: Promise<*>}}} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../../api/DeliveryServiceRequestService")} deliveryServiceRequestService
 * @param {import("../../../../models/MessageModel")} messageModel
 * @param {import("../../../../models/PropertiesModel")} propertiesModel
 * @param {import("../../../../models/UserModel")} userModel
 */
var FormEditDeliveryServiceController = function(deliveryService, origin, topologies, type, types, $scope, $state, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService, deliveryServiceRequestService, messageModel, propertiesModel, userModel) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller("FormDeliveryServiceController", { deliveryService: deliveryService, dsCurrent: deliveryService, origin: origin, topologies: topologies, type: type, types: types, $scope: $scope }));

	this.$onInit = function() {
		$scope.originalRoutingName = deliveryService.routingName;
		$scope.originalCDN = deliveryService.cdnId;
		$scope.loadGeoLimitCountriesRaw(deliveryService);
	};

	/**
	 * Creates a deletion DSR and the immediately fulfills it.
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
	 * @param {import("../../../../api/DeliveryServiceRequestService").DeliveryServiceRequest} dsRequest
	 * @param {string} commentValue
	 * @returns {Promise<void>}
	 */
	async function createAndCompleteDeletionRequest(deliveryService, dsRequest, commentValue) {
		// first create the ds request
		try {
			const response = await deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);

			const comment = {
				deliveryServiceRequestId: response.id,
				value: commentValue
			};

			// then create the ds request comment
			await deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);

			await Promise.all([
				// assign the ds request
				deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username),
				// set the status to 'complete'
				deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, "complete")
			]);
			// then, if all that works, delete the ds and navigate to the /delivery-services page
			await deliveryServiceService.deleteDeliveryService(deliveryService);
			messageModel.setMessages([ { level: "success", text: `Delivery service [ ${deliveryService.xmlId} ] deleted` } ], true);
			locationUtils.navigateToPath("/delivery-services");
		} catch (fault) {
			$anchorScroll(); // scrolls window to top
			messageModel.setMessages(fault.data.alerts, false);
		}
	}

	/**
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
	 */
	async function createDeliveryServiceDeleteRequest(deliveryService) {
		const params = {
			title: "Delivery Service Delete Request",
			message: "All delivery service deletions must be reviewed."
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
					if (userModel.user.role == propertiesModel.properties.dsRequests.overrideRole) {
						statuses.push({ id: $scope.COMPLETE, name: "Fulfill Request Immediately" });
					}
					return statuses;
				}
			}
		});
		/** @type {{comment: string; status: {id: number}}} */
		let options;
		try {
			options = await modalInstance.result;
		} catch {
			// This means the user cancelled.
			return;
		}
		let status = "draft";
		if (options.status.id === $scope.SUBMITTED || options.status.id === $scope.COMPLETE) {
			status = "submitted";
		};

		const dsRequest = {
			changeType: "delete",
			status: status,
			original: deliveryService
		};

		// if the user chooses to complete/fulfill the delete request immediately, a delivery service request will be made and marked as complete,
		// then if that is successful, the DS will be deleted
		if (options.status.id === $scope.COMPLETE) {
			await createAndCompleteDeletionRequest(deliveryService, dsRequest, options.comment);
			return
		}
		const response = await deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);
		const comment = {
			deliveryServiceRequestId: response.id,
			value: options.comment
		};
		await deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
		const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original.xmlId;
		messageModel.setMessages([ { level: "success", text: `Created request to ${dsRequest.changeType} the ${xmlId} delivery service` } ], true);
		locationUtils.navigateToPath("/delivery-service-requests");
	}

	/**
	 * Creates a new DSR for updating a Delivery Service.
	 *
	 * @param {*} dsRequest
	 * @param {string} dsRequestComment
	 * @param {boolean} autoFulfilled
	 * @returns {Promise<void>}
	 */
	async function createDeliveryServiceUpdateRequest(dsRequest, dsRequestComment, autoFulfilled) {
		const response = await deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);
		const comment = {
			deliveryServiceRequestId: response.id,
			value: dsRequestComment
		};
		await deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
		if (!autoFulfilled) {
			const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original.xmlId;
			messageModel.setMessages([ { level: "success", text: `Created request to ${dsRequest.changeType} the ${xmlId} delivery service` } ], true);
			locationUtils.navigateToPath("/delivery-service-requests");
			return;
		}
		await Promise.all([
			// assign the ds request
			deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username),
			// set the status to 'complete'
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, 'complete')
		]);
	};

	$scope.deliveryServiceName = angular.copy(deliveryService.xmlId);

	$scope.settings = {
		isNew: false,
		isRequest: false,
		saveLabel: "Update",
		deleteLabel: "Delete"
	};

	$scope.restrictTLS = Array.isArray(deliveryService.tlsVersions) && deliveryService.tlsVersions.length > 0;

	/**
	 * Saves the edit changes to the Delivery Service by creating a DSR to
	 * effect the change.
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService & {id: number}} deliveryService
	 * @returns
	 */
	async function saveWithRequest(deliveryService) {
		const params = {
			title: "Delivery Service Update Request",
			message: "All delivery service updates must be reviewed for completeness and accuracy before deployment."
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
					if (userModel.user.role == propertiesModel.properties.dsRequests.overrideRole) {
						statuses.push({ id: $scope.COMPLETE, name: "Fulfill Request Immediately" });
					}
					return statuses;
				}
			}
		});
		let options;
		try {
			options = await modalInstance.result;
		} catch {
			// this means they cancelled
			return;
		}
		let status = "draft";
		if (options.status.id == $scope.SUBMITTED || options.status.id == $scope.COMPLETE) {
			status = "submitted";
		};
		const dsRequest = {
			changeType: "update",
			status: status,
			requested: deliveryService
		};
		// if the user chooses to complete/fulfill the update request immediately, a delivery service request will be made and marked as complete,
		// then if that is successful, the DS will be updated
		if (options.status.id == $scope.COMPLETE) {
			try {
				await createDeliveryServiceUpdateRequest(dsRequest, options.comment, true);
			} catch(fault) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages(fault.data.alerts, false);
				return;
			}
			try {
				const response = await deliveryServiceService.updateDeliveryService(deliveryService);
				$state.reload(); // reloads all the resolves for the view
				messageModel.setMessages(response.alerts, false);
			} catch(fault) {
				// if the ds update fails, send to dsr view w/ error message
				locationUtils.navigateToPath('/delivery-service-requests');
				messageModel.setMessages(fault.data.alerts, true);
			}
		} else {
			createDeliveryServiceUpdateRequest(dsRequest, options.comment, false);
		}
	}

	/**
	 * Saves the edit changes to the Delivery Service.
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService & {id: number}} deliveryService
	 * @returns
	 */
	$scope.save = async function(deliveryService) {
		deliveryService.requiredCapabilities = Object.entries($scope.selectedCapabilities).filter(sc => (sc[1])).map(sc => sc[0])
		$scope.loadGeoLimitCountries(deliveryService);

		if (
			deliveryService.sslKeyVersion !== null &&
			deliveryService.sslKeyVersion !== 0 &&
			(
				$scope.originalRoutingName !== deliveryService.routingName ||
				$scope.originalCDN !== deliveryService.cdnId
			) &&
			type.includes("HTTP")
		) {

			const params = {
				title: "Cannot update Delivery Service",
				message: "Delivery Service has SSL Keys that cannot be updated"
			};

			$uibModal.open({
				templateUrl: "common/modules/dialog/text/dialog.text.tpl.html",
				controller: "DialogTextController",
				size: "md",
				resolve: {
					params,
					text: () => null
				}
			});
			return;
		}

		if (!$scope.restrictTLS) {
			deliveryService.tlsVersions = null;
		}

		// if ds requests are enabled in traffic_portal_properties.json, we'll create a ds request, else just update the ds
		if ($scope.dsRequestsEnabled) {
			saveWithRequest(deliveryService);
		} else {
			try {
				const response = await deliveryServiceService.updateDeliveryService(deliveryService);
				$state.reload(); // reloads all the resolves for the view
				messageModel.setMessages(response.alerts, false);
			} catch(fault) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages(fault.data.alerts, false);
			}
		}
	};

	/**
	 * Opens a dialog to confirm deletion of the current Delivery Service.
	 *
	 * @param {import("../../../../api/DeliveryServiceService").DeliveryService} deliveryService
	 */
	$scope.confirmDelete = async function(deliveryService) {
		const params = {
			title: `Delete Delivery Service: ${deliveryService.xmlId}`,
			key: deliveryService.xmlId
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/delete/dialog.delete.tpl.html",
			controller: "DialogDeleteController",
			size: "md",
			resolve: { params }
		});
		try {
			await modalInstance.result;
		} catch {
			// This means the user cancelled.
			return;
		}
		if ($scope.dsRequestsEnabled) {
			createDeliveryServiceDeleteRequest(deliveryService);
			return;
		}
		try {
			const response = await deliveryServiceService.deleteDeliveryService(deliveryService);
			messageModel.setMessages(response.alerts, true);
			locationUtils.navigateToPath("/delivery-services");
		} catch(fault) {
			$anchorScroll(); // scrolls window to top
			messageModel.setMessages(fault.data.alerts, false);
		}
	};

};

FormEditDeliveryServiceController.$inject = ["deliveryService", "origin", "topologies", "type", "types", "$scope", "$state", "$controller", "$uibModal", "$anchorScroll", "locationUtils", "deliveryServiceService", "deliveryServiceRequestService", "messageModel", "propertiesModel", "userModel"];
module.exports = FormEditDeliveryServiceController;
