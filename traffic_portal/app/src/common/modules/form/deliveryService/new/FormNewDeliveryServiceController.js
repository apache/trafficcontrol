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

var FormNewDeliveryServiceController = function(deliveryService, origin, topologies, type, types, $scope, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService, deliveryServiceRequestService, messageModel, propertiesModel, userModel) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService,
		dsCurrent: deliveryService, origin: origin, topologies: topologies, type: type, types: types, $scope: $scope }));

	$scope.deliveryServiceName = 'New';

	$scope.settings = {
		isNew: true,
		isRequest: false,
		saveLabel: 'Create'
	};

	$scope.restrictTLS = false;

	var createDeliveryServiceCreateRequest = function(dsRequest, dsRequestComment, autoFulfilled) {
		deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
			then(
				function(response) {
					var comment = {
						deliveryServiceRequestId: response.id,
						value: dsRequestComment
					};
					var promises = [];

					deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
						then(
							function() {
								if (!autoFulfilled) {
									const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original.xmlId;
									messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + xmlId + ' delivery service' } ], true);
									locationUtils.navigateToPath('/delivery-service-requests');
								}
							}
						);

					if (autoFulfilled) {
						// assign the ds request
						promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username));
						// set the status to 'complete'
						promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, 'complete'));
					}
				}
			);
	};


	$scope.save = function(deliveryService) {
		if (!$scope.restrictTLS) {
			deliveryService.tlsVersions = null;
		}
		// if ds requests are enabled in traffic_portal_properties.json, we'll create a ds request, else just create the ds
		if ($scope.dsRequestsEnabled) {
			var params = {
				title: "Delivery Service Create Request",
				message: 'All new delivery services must be reviewed for completeness and accuracy before deployment.'
			};
			var modalInstance = $uibModal.open({
				templateUrl: 'common/modules/dialog/deliveryServiceRequest/dialog.deliveryServiceRequest.tpl.html',
				controller: 'DialogDeliveryServiceRequestController',
				size: 'md',
				resolve: {
					params: function () {
						return params;
					},
					statuses: function() {
						var statuses = [
							{ id: $scope.DRAFT, name: 'Save Request as Draft' },
							{ id: $scope.SUBMITTED, name: 'Submit Request for Review and Deployment' }
						];
						if (userModel.user.roleName == propertiesModel.properties.dsRequests.overrideRole) {
							statuses.push({ id: $scope.COMPLETE, name: 'Fulfill Request Immediately' });
						}
						return statuses;
					}
				}
			});
			modalInstance.result.then(function(options) {
				var status = 'draft';
				if (options.status.id == $scope.SUBMITTED || options.status.id == $scope.COMPLETE) {
					status = 'submitted';
				};
				var dsRequest = {
					changeType: 'create',
					status: status,
					requested: deliveryService
				};
				// if the user chooses to complete/fulfill the create request immediately, the ds will be created and behind the
				// scenes a delivery service request will be created and marked as complete
				if (options.status.id == $scope.COMPLETE) {
					deliveryServiceService.createDeliveryService(deliveryService).
						then(
							function(result) {
								createDeliveryServiceCreateRequest(dsRequest, options.comment, true);
								messageModel.setMessages([ { level: 'success', text: 'Delivery Service [ ' + deliveryService.xmlId + ' ] created' } ], true);
								locationUtils.navigateToPath('/delivery-services/' + result.data.response[0].id + '?type=' + result.data.response[0].type);
							},
							function(fault) {
								$anchorScroll(); // scrolls window to top
								messageModel.setMessages(fault.data.alerts, false);
							}
						);

				} else {
					createDeliveryServiceCreateRequest(dsRequest, options.comment, false);
				}

			}, function () {
				// do nothing
			});
		} else {
			deliveryServiceService.createDeliveryService(deliveryService).
				then(
					function(result) {
						messageModel.setMessages([ { level: 'success', text: 'Delivery Service [ ' + deliveryService.xmlId + ' ] created' } ], true);
						locationUtils.navigateToPath('/delivery-services/' + result.data.response[0].id + '?type=' + result.data.response[0].type);
					},
					function(fault) {
						$anchorScroll(); // scrolls window to top
						messageModel.setMessages(fault.data.alerts, false);
					}
			);
		}
	};

};

FormNewDeliveryServiceController.$inject = ['deliveryService', 'origin', 'topologies', 'type', 'types', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = FormNewDeliveryServiceController;
