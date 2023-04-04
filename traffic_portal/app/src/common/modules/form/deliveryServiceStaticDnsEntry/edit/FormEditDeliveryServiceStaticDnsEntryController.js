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
 * @param {*} deliveryService
 * @param {*} staticDnsEntry
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/StaticDnsEntryService")} staticDnsEntryService
 */
var FormEditDeliveryServiceStaticDnsEntryController = function(deliveryService, staticDnsEntry, $scope, $controller, $uibModal, $anchorScroll, locationUtils, staticDnsEntryService) {

    // extends the FormDeliveryServiceController to inherit common methods
    angular.extend(this, $controller('FormDeliveryServiceStaticDnsEntryController', { deliveryService: deliveryService, staticDnsEntry: staticDnsEntry, $scope: $scope }));

    var deleteDeliveryServiceStaticDnsEntry = function(dsId, staticDnsEntryId) {
        staticDnsEntryService.deleteDeliveryServiceStaticDnsEntry(staticDnsEntryId)
            .then(function() {
                locationUtils.navigateToPath('/delivery-services/' + dsId + '/static-dns-entries');
            });
    };

    $scope.staticDnsEntry = staticDnsEntry;
    $scope.host = angular.copy($scope.staticDnsEntry.host);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(dsId, staticDnsEntry) {
        staticDnsEntryService.updateDeliveryServiceStaticDnsEntry(staticDnsEntry.id, staticDnsEntry).
        then(function() {
            $scope.host = angular.copy(staticDnsEntry.host);
            $anchorScroll(); // scrolls window to top
        });
    };

    $scope.confirmDelete = function(staticDnsEntry) {
        var params = {
            title: 'Delete Delivery Service Static DNS Entry on host: ' + staticDnsEntry.host + ' with address: ' + staticDnsEntry.address,
            key: staticDnsEntry.host
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
            deleteDeliveryServiceStaticDnsEntry(deliveryService.id, staticDnsEntry.id);
        }, function () {
            // do nothing
        });
    };

};

FormEditDeliveryServiceStaticDnsEntryController.$inject = ['deliveryService', 'staticDnsEntry', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'staticDnsEntryService'];
module.exports = FormEditDeliveryServiceStaticDnsEntryController;
