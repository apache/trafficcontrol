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

var FormEditDeliveryServiceStaticDnsEntryController = function(deliveryService, staticDnsEntry, $scope, $controller, $uibModal, $anchorScroll, locationUtils, staticDnsEntryService) {

    // extends the FormDeliveryServiceController to inherit common methods
    angular.extend(this, $controller('FormDeliveryServiceStaticDnsEntryController', { deliveryService: deliveryService, staticDnsEntry: staticDnsEntry, $scope: $scope }));

    // var deleteDeliveryServiceRegex = function(dsId, regexId) {
    //     deliveryServiceRegexService.deleteDeliveryServiceRegex(dsId, regexId)
    //         .then(function() {
    //             locationUtils.navigateToPath('/delivery-services/' + dsId + '/regexes');
    //         });
    // };

    $scope.staticDnsEntry = staticDnsEntry[0];
    $scope.host = angular.copy($scope.staticDnsEntry.host);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    // $scope.save = function(dsId, regex) {
    //     deliveryServiceRegexService.updateDeliveryServiceRegex(regex).
    //     then(function() {
    //         $scope.regexPattern = angular.copy(regex.pattern);
    //         $anchorScroll(); // scrolls window to top
    //     });
    // };

    // $scope.confirmDelete = function(regex) {
    //     var params = {
    //         title: 'Delete Delivery Service Regex: ' + regex.pattern,
    //         key: regex.pattern
    //     };
    //     var modalInstance = $uibModal.open({
    //         templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
    //         controller: 'DialogDeleteController',
    //         size: 'md',
    //         resolve: {
    //             params: function () {
    //                 return params;
    //             }
    //         }
    //     });
    //     modalInstance.result.then(function() {
    //         deleteDeliveryServiceRegex(deliveryService.id, regex.id);
    //     }, function () {
    //         // do nothing
    //     });
    // };

};

FormEditDeliveryServiceStaticDnsEntryController.$inject = ['deliveryService', 'staticDnsEntry', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'staticDnsEntryService'];
module.exports = FormEditDeliveryServiceStaticDnsEntryController;
