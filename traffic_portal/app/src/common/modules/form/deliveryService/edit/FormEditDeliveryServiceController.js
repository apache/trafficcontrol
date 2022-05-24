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

var FormEditDeliveryServiceController = function(deliveryService, origin, topologies, type, types, $scope, $state, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService, deliveryServiceRequestService, messageModel, propertiesModel, userModel) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, dsCurrent: deliveryService, origin: origin, topologies: topologies, type: type, types: types, $scope: $scope }));

	this.$onInit = function() {
		$scope.originalRoutingName = deliveryService.routingName;
		$scope.originalCDN = deliveryService.cdnId;
	};

	var createDeliveryServiceDeleteRequest = function(deliveryService) {
		var params = {
			title: "Delivery Service Delete Request",
			message: 'All delivery service deletions must be reviewed.'
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
					if (userModel.user.role == propertiesModel.properties.dsRequests.overrideRole) {
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
				changeType: 'delete',
				status: status,
				original: deliveryService
			};

			// if the user chooses to complete/fulfill the delete request immediately, a delivery service request will be made and marked as complete,
			// then if that is successful, the DS will be deleted
			if (options.status.id === $scope.COMPLETE) {
				// first create the ds request
				deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest)
					.then(function(response) {
						var comment = {
							deliveryServiceRequestId: response.id,
							value: options.comment
						};
						// then create the ds request comment
						deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
							then(
								function() {
									var promises = [];
									// assign the ds request
									promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username));
									// set the status to 'complete'
									promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, 'complete'));
								}
							// then, if all that works, delete the ds and navigate to the /delivery-services page
							).then(
								function() {
									deliveryServiceService.deleteDeliveryService(deliveryService).
										then(
											function() {
												messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + deliveryService.xmlId + ' ] deleted' } ], true);
												locationUtils.navigateToPath('/delivery-services');
											}
										);
								}
							);
					}
					// handle any failures just once
					).catch(function(fault) {
						$anchorScroll(); // scrolls window to top
						messageModel.setMessages(fault.data.alerts, false);
					}
				);
			} else {
				deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
					then(
						function(response) {
							var comment = {
								deliveryServiceRequestId: response.id,
								value: options.comment
							};
							deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
								then(
									function() {
										const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original.xmlId;
										messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + xmlId + ' delivery service' } ], true);
										locationUtils.navigateToPath('/delivery-service-requests');
									}
								);
						}
					);
			}
		}, function () {
			// do nothing
		});
	};

	var createDeliveryServiceUpdateRequest = function(dsRequest, dsRequestComment, autoFulfilled) {
		return deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
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

	$scope.deliveryServiceName = angular.copy(deliveryService.xmlId);

	$scope.settings = {
		isNew: false,
		isRequest: false,
		saveLabel: 'Update',
		deleteLabel: 'Delete'
	};

	$scope.restrictTLS = deliveryService.tlsVersions instanceof Array && deliveryService.tlsVersions.length > 0;

	$scope.save = function(deliveryService) {
		if (deliveryService.sslKeyVersion !== null && deliveryService.sslKeyVersion !== 0 &&
			($scope.originalRoutingName !== deliveryService.routingName || $scope.originalCDN !== deliveryService.cdnId) &&
			type.indexOf("HTTP") !== -1) {

			let params = {
				title: "Cannot update Delivery Service",
				message: "Delivery Service has SSL Keys that cannot be updated"
			};

			$uibModal.open({
				templateUrl: "common/modules/dialog/text/dialog.text.tpl.html",
				controller: "DialogTextController",
				size: "md",
				resolve: {
					params: function() {
						return params;
					},
					text: function() {
						return null;
					}
				}
			});
			return;
		}

		if (!$scope.restrictTLS) {
			deliveryService.tlsVersions = null;
		}

		// if ds requests are enabled in traffic_portal_properties.json, we'll create a ds request, else just update the ds
		if ($scope.dsRequestsEnabled) {
			var params = {
				title: "Delivery Service Update Request",
				message: 'All delivery service updates must be reviewed for completeness and accuracy before deployment.'
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
						if (userModel.user.role == propertiesModel.properties.dsRequests.overrideRole) {
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
					changeType: 'update',
					status: status,
					requested: deliveryService
				};
				// if the user chooses to complete/fulfill the update request immediately, a delivery service request will be made and marked as complete,
				// then if that is successful, the DS will be updated
				if (options.status.id == $scope.COMPLETE) {
					createDeliveryServiceUpdateRequest(dsRequest, options.comment, true).then(
						function() {
							deliveryServiceService.updateDeliveryService(deliveryService).
								then(
									function() {
										$state.reload(); // reloads all the resolves for the view
										messageModel.setMessages([ { level: 'success', text: 'Delivery Service [ ' + deliveryService.xmlId + ' ] updated' } ], false);
									}
								).catch(function(fault) {
									// if the ds update fails, send to dsr view w/ error message
									locationUtils.navigateToPath('/delivery-service-requests');
									messageModel.setMessages(fault.data.alerts, true);
							});
						}).catch(function(fault) {
							$anchorScroll(); // scrolls window to top
							messageModel.setMessages(fault.data.alerts, false);
					});
				} else {
					createDeliveryServiceUpdateRequest(dsRequest, options.comment, false);
				}

			}, function () {
				// do nothing
			});
		} else {
			deliveryServiceService.updateDeliveryService(deliveryService).
				then(
					function() {
						$state.reload(); // reloads all the resolves for the view
						messageModel.setMessages([ { level: 'success', text: 'Delivery Service [ ' + deliveryService.xmlId + ' ] updated' } ], false);
					},
					function(fault) {
						$anchorScroll(); // scrolls window to top
						messageModel.setMessages(fault.data.alerts, false);
					}
				);
		}
	};

	$scope.confirmDelete = function(deliveryService) {
		var params = {
			title: 'Delete Delivery Service: ' + deliveryService.xmlId,
			key: deliveryService.xmlId
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
			controller: 'DialogDeleteController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			if ($scope.dsRequestsEnabled) {
				createDeliveryServiceDeleteRequest(deliveryService);
			} else {
				deliveryServiceService.deleteDeliveryService(deliveryService)
					.then(
						function() {
							messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + deliveryService.xmlId + ' ] deleted. ' +
									'Perform a CDN snapshot then queue updates on all servers in the cdn for the changes to take affect.' } ], true);
							locationUtils.navigateToPath('/delivery-services');
						},
						function(fault) {
							$anchorScroll(); // scrolls window to top
							messageModel.setMessages(fault.data.alerts, false);
						}
					);
			}
		}, function () {
			// do nothing
		});
	};

};

FormEditDeliveryServiceController.$inject = ['deliveryService', 'origin', 'topologies', 'type', 'types', '$scope', '$state', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = FormEditDeliveryServiceController;
