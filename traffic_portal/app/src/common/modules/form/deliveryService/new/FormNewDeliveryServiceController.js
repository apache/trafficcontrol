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

var FormNewDeliveryServiceController = function(deliveryService, type, types, $scope, $controller, $uibModal, deliveryServiceRequestService) {

	// extends the FormDeliveryServiceController to inherit common methods
	angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, type: type, types: types, $scope: $scope }));

	$scope.deliveryServiceName = 'New';

	$scope.settings = {
		isNew: true,
		isRequest: false,
		saveLabel: 'Create'
	};

	$scope.save = function(deliveryService) {
		var params = {
			title: "Create Delivery Service",
			message: 'All new delivery services must be reviewed for completeness and accuracy before deployment. A request will be created for you. Please select the status of your request.'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function() {
					return [
						{ id: $scope.DRAFT, name: 'Save Request as Draft' },
						{ id: $scope.SUBMITTED, name: 'Submit Request for Review / Deployment' }
					];
				}
			}
		});
		modalInstance.result.then(function(action) {
			var dsRequest = {
				changeType: 'create',
				status: (action.id == $scope.SUBMITTED) ? 'submitted' : 'draft',
				request: deliveryService
			};
			deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);
		}, function () {
			// do nothing
		});

	};

};

FormNewDeliveryServiceController.$inject = ['deliveryService', 'type', 'types', '$scope', '$controller', '$uibModal', 'deliveryServiceRequestService'];
module.exports = FormNewDeliveryServiceController;
