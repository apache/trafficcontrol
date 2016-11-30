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

var FormEditDeliveryServiceController = function(deliveryService, $scope, $controller, $uibModal, $anchorScroll, locationUtils, deliveryServiceService) {

    // extends the FormDeliveryServiceController to inherit common methods
    angular.extend(this, $controller('FormDeliveryServiceController', { deliveryService: deliveryService, $scope: $scope }));

    var deleteDeliveryService = function(deliveryService) {
        deliveryServiceService.deleteDeliveryService(deliveryService.id)
            .then(function() {
                locationUtils.navigateToPath('/configure/delivery-services');
            });
    };

    $scope.deliveryServiceName = angular.copy(deliveryService.displayName);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(deliveryService) {
        deliveryServiceService.updateDeliveryService(deliveryService).
            then(function() {
                $scope.deliveryServiceName = angular.copy(deliveryService.displayName);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(deliveryService) {
        var params = {
            title: 'Delete Delivery Service: ' + deliveryService.displayName,
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
            deleteDeliveryService(deliveryService);
        }, function () {
            // do nothing
        });
    };

};

FormEditDeliveryServiceController.$inject = ['deliveryService', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'deliveryServiceService'];
module.exports = FormEditDeliveryServiceController;
