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

var TableDeliveryServicesController = function(deliveryServices, $anchorScroll, $scope, $state, $location, $uibModal, $window, deliveryServiceService, deliveryServiceRequestService, dateUtils, deliveryServiceUtils, locationUtils, messageModel, propertiesModel, userModel) {

    let deliveryServicesTable;

    var protocols = deliveryServiceUtils.protocols;

    var qstrings = deliveryServiceUtils.qstrings;

    var geoProviders = deliveryServiceUtils.geoProviders;

    var geoLimits = deliveryServiceUtils.geoLimits;

    var rrhs = deliveryServiceUtils.rrhs;

    var dsRequestsEnabled = propertiesModel.properties.dsRequests.enabled;

    var showCustomCharts = propertiesModel.properties.deliveryServices.charts.customLink.show;

    var createDeliveryService = function(typeName) {
        var path = '/delivery-services/new?type=' + typeName;
        locationUtils.navigateToPath(path);
    };

    var clone = function(ds) {
        var params = {
            title: 'Clone Delivery Service: ' + ds.xmlId,
            message: "Please select a content routing category for the clone"
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
            locationUtils.navigateToPath('/delivery-services/' + ds.id + '/clone?type=' + type.name);
        });
    };

    var confirmDelete = function(deliveryService) {
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
            if (dsRequestsEnabled) {
                createDeliveryServiceDeleteRequest(deliveryService);
            } else {
                deliveryServiceService.deleteDeliveryService(deliveryService)
                    .then(
                        function() {
                            messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + deliveryService.xmlId + ' ] deleted' } ], false);
                            $scope.refresh();                        },
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
                changeType: 'delete',
                status: status,
                deliveryService: deliveryService
            };

            // if the user chooses to complete/fulfill the delete request immediately, the ds will be deleted and behind the
            // scenes a delivery service request will be created and marked as complete
            if (options.status.id == $scope.COMPLETE) {
                // first delete the ds
                deliveryServiceService.deleteDeliveryService(deliveryService)
                    .then(
                        function() {
                            // then create the ds request
                            deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
                            then(
                                function(response) {
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
                                            promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.id));
                                            // set the status to 'complete'
                                            promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, 'complete'));
                                            // and finally refresh the delivery services table
                                            messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + deliveryService.xmlId + ' ] deleted' } ], false);
                                            $scope.refresh();
                                        }
                                    );
                                }
                            );
                        },
                        function(fault) {
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
                                        messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + dsRequest.deliveryService.xmlId + ' delivery service' } ], true);
                                        locationUtils.navigateToPath('/delivery-service-requests');
                                    }
                                );
                        }
                    );
            }
        });
    };

    $scope.deliveryServices = deliveryServices;

    $scope.getRelativeTime = dateUtils.getRelativeTime;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.DRAFT = 0;
    $scope.SUBMITTED = 1;
    $scope.REJECTED = 2;
    $scope.PENDING = 3;
    $scope.COMPLETE = 4;

    $scope.columns = [
        { "name": "Active", "visible": true, "searchable": true },
        { "name": "Anonymous Blocking", "visible": false, "searchable": false },
        { "name": "CDN", "visible": true, "searchable": true },
        { "name": "Check Path", "visible": false, "searchable": false },
        { "name": "Consistent Hash Query Params", "visible": false, "searchable": false },
        { "name": "Consistent Hash Regex", "visible": false, "searchable": false },
        { "name": "Deep Caching Type", "visible": false, "searchable": false },
        { "name": "Display Name", "visible": true, "searchable": true },
        { "name": "DNS Bypass CNAME", "visible": false, "searchable": false },
        { "name": "DNS Bypass IP", "visible": false, "searchable": false },
        { "name": "DNS Bypass IPv6", "visible": false, "searchable": false },
        { "name": "DNS Bypass TTL", "visible": false, "searchable": false },
        { "name": "DNS TTL", "visible": false, "searchable": false },
        { "name": "DSCP", "visible": true, "searchable": true },
        { "name": "ECS Enabled", "visible": false, "searchable": false },
        { "name": "Edge Header Rewrite Rules", "visible": false, "searchable": false },
        { "name": "First Header Rewrite Rules", "visible": false, "searchable": false },
        { "name": "FQ Pacing Rate", "visible": false, "searchable": false },
        { "name": "Geo Limit", "visible": false, "searchable": false },
        { "name": "Geo Limit Countries", "visible": false, "searchable": false },
        { "name": "Geo Limit Redirect URL", "visible": false, "searchable": false },
        { "name": "Geolocation Provider", "visible": false, "searchable": false },
        { "name": "Geo Miss Latitude", "visible": false, "searchable": false },
        { "name": "Geo Miss Longitude", "visible": false, "searchable": false },
        { "name": "Global Max Mbps", "visible": false, "searchable": false },
        { "name": "Global Max TPS", "visible": false, "searchable": false },
        { "name": "HTTP Bypass FQDN", "visible": false, "searchable": false },
        { "name": "ID", "visible": false, "searchable": false },
        { "name": "Info URL", "visible": true, "searchable": true },
        { "name": "Initial Dispersion", "visible": false, "searchable": false },
        { "name": "Inner Header Rewrite Rules", "visible": false, "searchable": false },
        { "name": "IPv6 Routing", "visible": false, "searchable": false },
        { "name": "Last Header Rewrite Rules", "visible": false, "searchable": false },
        { "name": "Last Updated", "visible": false, "searchable": false },
        { "name": "Long Desc 1", "visible": false, "searchable": false },
        { "name": "Long Desc 2", "visible": false, "searchable": false },
        { "name": "Long Desc 3", "visible": false, "searchable": false },
        { "name": "Max DNS Answers", "visible": false, "searchable": false },
        { "name": "Max Origin Connections", "visible": false, "searchable": false },
        { "name": "Mid Header Rewrite Rules", "visible": false, "searchable": false },
        { "name": "Multi-Site Origin", "visible": false, "searchable": false },
        { "name": "Origin Shield", "visible": false, "searchable": false },
        { "name": "Origin FQDN", "visible": true, "searchable": true },
        { "name": "Profile", "visible": true, "searchable": true },
        { "name": "Protocol", "visible": true, "searchable": true },
        { "name": "Qstring Handling", "visible": false, "searchable": false },
        { "name": "Range Request Handling", "visible": false, "searchable": false },
        { "name": "Regex Remap Expression", "visible": false, "searchable": false },
        { "name": "Regional Geoblocking", "visible": false, "searchable": false },
        { "name": "Raw Remap Text", "visible": false, "searchable": false },
        { "name": "Routing Name", "visible": false, "searchable": false },
        { "name": "Service Category", "visible": false, "searchable": false },
        { "name": "Signed", "visible": false, "searchable": false },
        { "name": "Signing Algorithm", "visible": false, "searchable": false },
        { "name": "Range Slice Block Size", "visible": false, "searchable": false },
        { "name": "Tenant", "visible": true, "searchable": true },
        { "name": "Topology", "visible": true, "searchable": true },
        { "name": "TR Request Headers", "visible": false, "searchable": false },
        { "name": "TR Response Headers", "visible": false, "searchable": false },
        { "name": "Type", "visible": true, "searchable": true },
        { "name": "XML ID (Key)", "visible": true, "searchable": true }
    ];

    $scope.contextMenuItems = [
        {
            text: 'Open in New Tab',
            click: function ($itemScope) {
                $window.open('/#!/delivery-services/' + $itemScope.ds.id + '?type=' + $itemScope.ds.type, '_blank');
            }
        },
        null, // Divider
        {
            text: 'Edit',
            click: function ($itemScope) {
                $scope.editDeliveryService($itemScope.ds);
            }
        },
        {
            text: 'Clone',
            click: function ($itemScope) {
                clone($itemScope.ds);
            }
        },
        {
            text: 'Delete',
            click: function ($itemScope) {
                confirmDelete($itemScope.ds);
            }
        },
        null, // Divider
        {
            text: 'View Charts',
            click: function ($itemScope, evt) {
                $scope.viewCharts($itemScope.ds, evt);
            }
        },
        null, // Divider
        {
            text: 'Manage SSL Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/ssl-keys?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage URL Sig Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/url-sig-keys?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage URI Signing Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/uri-signing-keys?type=' + $itemScope.ds.type);
            }
        },
        null, // Divider
        {
            text: 'Manage Invalidation Requests',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/jobs?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Origins',
            displayed: function ($itemScope) {
                // only show for non-steering* delivery services
                return $itemScope.ds.type.indexOf('STEERING') == -1;
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/origins?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Regexes',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/regexes?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Required Server Capabilities',
            displayed: function ($itemScope) {
                // only show for DNS* or HTTP* delivery services
                return ($itemScope.ds.type.indexOf('DNS') != -1 || $itemScope.ds.type.indexOf('HTTP') != -1);
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/required-server-capabilities?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Servers',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/servers?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Targets',
            displayed: function ($itemScope) {
                // only show for steering* delivery services
                return $itemScope.ds.type.indexOf('STEERING') != -1;
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/targets?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Static DNS Entries',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/static-dns-entries?type=' + $itemScope.ds.type);
            }
        }
    ];

    $scope.editDeliveryService = function(ds) {
        var path = '/delivery-services/' + ds.id + '?type=' + ds.type;
        locationUtils.navigateToPath(path);
    };

    $scope.viewCharts = function(ds, $event) {
        if ($event) {
            $event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
        }

        if (showCustomCharts) {
            deliveryServiceUtils.openCharts(ds);
        } else {
            locationUtils.navigateToPath('/delivery-services/' + ds.id + '/charts?type=' + ds.type);
        }
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

    $scope.geoProvider = function(ds) {
        return geoProviders[ds.geoProvider];
    };

    $scope.geoLimit = function(ds) {
        return geoLimits[ds.geoLimit];
    };

    $scope.rrh = function(ds) {
        return rrhs[ds.rangeRequestHandling];
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

    $scope.toggleVisibility = function(colName) {
        const col = deliveryServicesTable.column(colName + ':name');
        col.visible(!col.visible());
        deliveryServicesTable.rows().invalidate().draw();
    };

    angular.element(document).ready(function () {
        deliveryServicesTable = $('#deliveryServicesTable').DataTable({
            "lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "aaSorting": [],
            "columns": $scope.columns,
            "initComplete": function(settings, json) {
                try {
                    // need to create the show/hide column checkboxes and bind to the current visibility
                    $scope.columns = JSON.parse(localStorage.getItem('DataTables_deliveryServicesTable_/')).columns;
                } catch (e) {
                    console.error("Failure to retrieve required column info from localStorage (key=DataTables_deliveryServicesTable_/):", e);
                }
            }
        });
    });

};

TableDeliveryServicesController.$inject = ['deliveryServices', '$anchorScroll', '$scope', '$state', '$location', '$uibModal', '$window', 'deliveryServiceService', 'deliveryServiceRequestService', 'dateUtils', 'deliveryServiceUtils', 'locationUtils', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = TableDeliveryServicesController;
