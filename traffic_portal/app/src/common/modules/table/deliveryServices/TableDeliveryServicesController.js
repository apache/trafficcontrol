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

var TableDeliveryServicesController = function(deliveryServices, $scope, $state, $location, $uibModal, dateUtils, deliveryServiceUtils, locationUtils, propertiesModel) {

    var protocols = deliveryServiceUtils.protocols;

    var qstrings = deliveryServiceUtils.qstrings;

    var createDeliveryService = function(typeName) {
        var path = '/delivery-services/new?type=' + typeName;
        locationUtils.navigateToPath(path);
    };

    $scope.deliveryServices = deliveryServices;

    $scope.showChartsButton = propertiesModel.properties.deliveryServices.charts.show;

    $scope.openCharts = deliveryServiceUtils.openCharts;

    $scope.getRelativeTime = dateUtils.getRelativeTime;

    $scope.editDeliveryService = function(ds) {
        var path = '/delivery-services/' + ds.id + '?type=' + ds.type;
        locationUtils.navigateToPath(path);
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.protocol = function(ds) {
        return protocols[ds.protocol];
    };

    $scope.qstring = function(ds) {
        return qstrings[ds.qstringIgnore];
    };

    $scope.selectDSType = function() {
        var params = {
            title: 'Create Delivery Service',
            message: "Please select a content routing category"
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
                    // the following represent the 4 categories of delivery services
                    // the ids are arbitrary but the dialog.select dropdown needs them
                    return [
                        { id: 1, name: 'ANY_MAP' },
                        { id: 2, name: 'DNS' },
                        { id: 3, name: 'HTTP' },
                        { id: 4, name: 'STEERING' }
                    ];
                }
            }
        });
        modalInstance.result.then(function(type) {
            createDeliveryService(type.name);
        }, function () {
            // do nothing
        });
    };

    $scope.compareDSs = function() {
        var params = {
            title: 'Compare Delivery Services',
            message: "Please select 2 delivery services to compare",
            label: "xmlId"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/compare/dialog.compare.tpl.html',
            controller: 'DialogCompareController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function(deliveryServiceService) {
                    return deliveryServiceService.getDeliveryServices();
                }
            }
        });
        modalInstance.result.then(function(dss) {
            $location.path($location.path() + '/compare/' + dss[0].id + '/' + dss[1].id);
        }, function () {
            // do nothing
        });
    };

    angular.element(document).ready(function () {
        $('#deliveryServicesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "columnDefs": [
                { 'orderable': false, 'targets': 12 }
            ],
            "aaSorting": []
        });
    });

};

TableDeliveryServicesController.$inject = ['deliveryServices', '$scope', '$state', '$location', '$uibModal', 'dateUtils', 'deliveryServiceUtils', 'locationUtils', 'propertiesModel'];
module.exports = TableDeliveryServicesController;
