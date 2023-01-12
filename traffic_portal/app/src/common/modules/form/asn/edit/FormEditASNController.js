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
 * @param {*} asn
 * @param {*} $scope
 * @param {import("angular").IControllerService} $controller
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/ASNService")} asnService
 */
var FormEditASNController = function(asn, $scope, $controller, $uibModal, $anchorScroll, locationUtils, asnService) {

    // extends the FormASNController to inherit common methods
    angular.extend(this, $controller('FormASNController', { asn: asn, $scope: $scope }));

    var deleteASN = function(asn) {
        asnService.deleteASN(asn.id)
            .then(function() {
                locationUtils.navigateToPath('/asns');
            });
    };

    $scope.asnName = angular.copy(asn.asn);

    $scope.settings = {
        isNew: false,
        saveLabel: 'Update'
    };

    $scope.save = function(asn) {
        asnService.updateASN(asn).
            then(function() {
                $scope.asnName = angular.copy(asn.asn);
                $anchorScroll(); // scrolls window to top
            });
    };

    $scope.confirmDelete = function(asn) {
        var params = {
            title: 'Delete ASN: ' + asn.asn,
            key: asn.asn.toString()
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
            deleteASN(asn);
        }, function () {
            // do nothing
        });
    };

};

FormEditASNController.$inject = ['asn', '$scope', '$controller', '$uibModal', '$anchorScroll', 'locationUtils', 'asnService'];
module.exports = FormEditASNController;
