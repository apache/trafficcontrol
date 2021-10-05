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

var TableDeliveryServicesController = function(tableName, deliveryServices, filter, $anchorScroll, $document, $scope, $state, $location, $uibModal, deliveryServiceService, deliveryServiceRequestService, dateUtils, deliveryServiceUtils, locationUtils, messageModel, propertiesModel, userModel) {

    /**
     * Gets value to display a default tooltip.
     */
    function defaultTooltip(params) {
        return params.value;
    }

    /**
     * Formats the contents of a 'lastUpdated' column cell as "relative to now".
     */
    function dateCellFormatter(params) {
        return params.value ? dateUtils.getRelativeTime(params.value) : params.value;
    }

    /** The columns of the ag-grid table */
    const columns = [
        {
            headerName: "Active",
            field: "active",
            hide: false
        },
        {
            headerName: "Anonymous Blocking",
            field: "anonymousBlockingEnabled",
            hide: true
        },
        {
            headerName: "CDN",
            field: "cdnName",
            hide: false
        },
        {
            headerName: "Check Path",
            field: "checkPath",
            hide: true
        },
        {
            headerName: "Consistent Hash Query Params",
            field: "consistentHashQueryParams",
            hide: true,
            valueFormatter: function(params) {
                return params.data.consistentHashQueryParams.join(', ');
            },
            tooltipValueGetter: function(params) {
                return params.data.consistentHashQueryParams.join(', ');
            }
        },
        {
            headerName: "Consistent Hash Regex",
            field: "consistentHashRegex",
            hide: true
        },
        {
            headerName: "Deep Caching Type",
            field: "deepCachingType",
            hide: true
        },
        {
            headerName: "Display Name",
            field: "displayName",
            hide: false
        },
        {
            headerName: "DNS Bypass CNAME",
            field: "dnsBypassCname",
            hide: true
        },
        {
            headerName: "DNS Bypass IP",
            field: "dnsBypassIp",
            hide: true
        },
        {
            headerName: "DNS Bypass IPv6",
            field: "dnsBypassIp6",
            hide: true
        },
        {
            headerName: "DNS Bypass TTL",
            field: "dnsBypassTtl",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "DNS TTL",
            field: "ccrDnsTtl",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "DSCP",
            field: "dscp",
            hide: false,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "ECS Enabled",
            field: "ecsEnabled",
            hide: true
        },
        {
            headerName: "Edge Header Rewrite Rules",
            field: "edgeHeaderRewrite",
            hide: true
        },
        {
            headerName: "First Header Rewrite Rules",
            field: "firstHeaderRewrite",
            hide: true
        },
        {
            headerName: "FQ Pacing Rate",
            field: "fqPacingRate",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Geo Limit",
            field: "geoLimit",
            hide: true,
            valueFormatter: function(params) {
                return deliveryServiceUtils.geoLimits[params.data.geoLimit];
            },
            tooltipValueGetter: function(params) {
                return deliveryServiceUtils.geoLimits[params.data.geoLimit];
            }
        },
        {
            headerName: "Geo Limit Countries",
            field: "geoLimitCountries",
            hide: true
        },
        {
            headerName: "Geo Limit Redirect URL",
            field: "geoLimitRedirectURL",
            hide: true
        },
        {
            headerName: "Geolocation Provider",
            field: "geoProvider",
            hide: true,
            valueFormatter: function(params) {
                return deliveryServiceUtils.geoProviders[params.data.geoProvider];
            },
            tooltipValueGetter: function(params) {
                return deliveryServiceUtils.geoProviders[params.data.geoProvider];
            }
        },
        {
            headerName: "Geo Miss Latitude",
            field: "missLat",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Geo Miss Longitude",
            field: "missLong",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Global Max Mbps",
            field: "globalMaxMbps",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Global Max TPS",
            field: "globalMaxTps",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "HTTP Bypass FQDN",
            field: "httpBypassFqdn",
            hide: true
        },
        {
            headerName: "ID",
            field: "id",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Info URL",
            field: "infoUrl",
            hide: true
        },
        {
            headerName: "Initial Dispersion",
            field: "initialDispersion",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Inner Header Rewrite Rules",
            field: "innerHeaderRewrite",
            hide: true
        },
        {
            headerName: "IPv6 Routing",
            field: "ipv6RoutingEnabled",
            hide: true
        },
        {
            headerName: "Last Header Rewrite Rules",
            field: "lastHeaderRewrite",
            hide: true
        },
        {
            headerName: "Last Updated",
            field: "lastUpdated",
            hide: true,
            filter: "agDateColumnFilter",
            valueFormatter: dateCellFormatter
        },
        {
            headerName: "Long Desc",
            field: "longDesc",
            hide: true
        },
        {
            headerName: "Max DNS Answers",
            field: "maxDnsAnswers",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Max Origin Connections",
            field: "maxOriginConnections",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Max Request Header Bytes",
            field: "maxRequestHeaderBytes",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Mid Header Rewrite Rules",
            field: "midHeaderRewrite",
            hide: true
        },
        {
            headerName: "Multi-Site Origin",
            field: "multiSiteOrigin",
            hide: true
        },
        {
            headerName: "Origin Shield",
            field: "originShield",
            hide: true
        },
        {
            headerName: "Origin FQDN",
            field: "orgServerFqdn",
            hide: false
        },
        {
            headerName: "Profile",
            field: "profileName",
            hide: true
        },
        {
            headerName: "Protocol",
            field: "protocol",
            hide: false,
            valueFormatter: function(params) {
                return deliveryServiceUtils.protocols[params.data.protocol];
            },
            tooltipValueGetter: function(params) {
                return deliveryServiceUtils.protocols[params.data.protocol];
            }
        },
        {
            headerName: "Qstring Handling",
            field: "qstringIgnore",
            hide: true,
            valueFormatter: function(params) {
                return deliveryServiceUtils.qstrings[params.data.qstringIgnore];
            },
            tooltipValueGetter: function(params) {
                return deliveryServiceUtils.qstrings[params.data.qstringIgnore];
            }
        },
        {
            headerName: "Range Request Handling",
            field: "rangeRequestHandling",
            hide: true,
            valueFormatter: function(params) {
                return deliveryServiceUtils.rrhs[params.data.rangeRequestHandling];
            },
            tooltipValueGetter: function(params) {
                return deliveryServiceUtils.rrhs[params.data.rangeRequestHandling];
            }
        },
        {
            headerName: "Regex Remap Expression",
            field: "regexRemap",
            hide: true
        },
        {
            headerName: "Regional Geoblocking",
            field: "regionalGeoBlocking",
            hide: true
        },
        {
            headerName: "Raw Remap Text",
            field: "remapText",
            hide: true
        },
        {
            headerName: "Routing Name",
            field: "routingName",
            hide: true
        },
        {
            headerName: "Service Category",
            field: "serviceCategory",
            hide: true
        },
        {
            headerName: "Signed",
            field: "signed",
            hide: true
        },
        {
            headerName: "Signing Algorithm",
            field: "signingAlgorithm",
            hide: true
        },
        {
            headerName: "Range Slice Block Size",
            field: "rangeSliceBlockSize",
            hide: true,
            filter: "agNumberColumnFilter"
        },
        {
            headerName: "Tenant",
            field: "tenant",
            hide: false
        },
        {
            headerName: "Topology",
            field: "topology",
            hide: false
        },
        {
            headerName: "TR Request Headers",
            field: "trRequestHeaders",
            hide: true
        },
        {
            headerName: "TR Response Headers",
            field: "trResponseHeaders",
            hide: true
        },
        {
            headerName: "Type",
            field: "type",
            hide: false
        },
        {
            headerName: "XML ID (Key)",
            field: "xmlId",
            hide: false
        }
    ];

    let dsRequestsEnabled = propertiesModel.properties.dsRequests.enabled;

    let showCustomCharts = propertiesModel.properties.deliveryServices.charts.customLink.show;

    var createDeliveryService = function(typeName) {
        var path = '/delivery-services/new?type=' + typeName;
        locationUtils.navigateToPath(path);
    };

    $scope.clone = function(ds, event) {
        event.stopPropagation();
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

    $scope.confirmDelete = function(deliveryService, event) {
        event.stopPropagation();
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
                                            promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username));
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
                                        const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original.xmlId;
                                        messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + xmlId + ' delivery service' } ], true);
                                        locationUtils.navigateToPath('/delivery-service-requests');
                                    }
                                );
                        }
                    );
            }
        });
    };

    /** All of the delivery services - lastUpdated fields converted to actual Dates */
    $scope.deliveryServices = deliveryServices.map(
        function(x) {
            x.lastUpdated = x.lastUpdated ? new Date(x.lastUpdated.replace("+00", "Z")) : x.lastUpdated;
        });

    /** The currently selected server - at the moment only used by the context menu */
    $scope.deliveryService = {
        xmlId: "",
        id: -1
    };

    $scope.quickSearch = '';

    $scope.pageSize = 100;

    $scope.mouseDownSelectionText = "";

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.DRAFT = 0;
    $scope.SUBMITTED = 1;
    $scope.REJECTED = 2;
    $scope.PENDING = 3;
    $scope.COMPLETE = 4;

    $scope.viewCharts = function(ds, $event) {
        $event.stopPropagation();
        if (showCustomCharts) {
            deliveryServiceUtils.openCharts(ds);
        } else {
            locationUtils.navigateToPath('/delivery-services/' + ds.id + '/charts?type=' + ds.type);
        }
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
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

    /** Toggles the visibility of a column that has the ID provided as 'col'. */
    $scope.toggleVisibility = function(col) {
        const visible = $scope.gridOptions.columnApi.getColumn(col).isVisible();
        $scope.gridOptions.columnApi.setColumnVisible(col, !visible);
    };

    /** Options, configuration, data and callbacks for the ag-grid table. */
    $scope.gridOptions = {
        columnDefs: columns,
        enableCellTextSelection:true,
        suppressMenuHide: true,
        multiSortKey: 'ctrl',
        alwaysShowVerticalScroll: true,
        defaultColDef: {
            filter: true,
            sortable: true,
            resizable: true,
            tooltipValueGetter: defaultTooltip
        },
        rowData: deliveryServices,
        pagination: true,
        paginationPageSize: $scope.pageSize,
        rowBuffer: 0,
        tooltipShowDelay: 500,
        allowContextMenuWithControlKey: true,
        preventDefaultOnContextMenu: true,
        colResizeDefault: "shift",
        onColumnVisible: function(params) {
            if (params.visible){
                return;
            }
            const filterModel = $scope.gridOptions.api.getFilterModel();
            for (let column of params.columns) {
                if (column.filterActive) {
                    if (column.colId in filterModel) {
                        delete filterModel[column.colId];
                        $scope.gridOptions.api.setFilterModel(filterModel);
                    }
                }
            }
        },
        onCellContextMenu: function(params) {
            $scope.showMenu = true;
            $scope.menuStyle.left = String(params.event.clientX) + "px";
            $scope.menuStyle.top = String(params.event.clientY) + "px";
            $scope.menuStyle.bottom = "unset";
            $scope.menuStyle.right = "unset";
            $scope.$apply();
            const boundingRect = document.getElementById("context-menu").getBoundingClientRect();

            if (boundingRect.bottom > window.innerHeight){
                $scope.menuStyle.bottom = String(window.innerHeight - params.event.clientY) + "px";
                $scope.menuStyle.top = "unset";
            }
            if (boundingRect.right > window.innerWidth) {
                $scope.menuStyle.right = String(window.innerWidth - params.event.clientX) + "px";
                $scope.menuStyle.left = "unset";
            }
            $scope.deliveryService = params.data;
            $scope.$apply();
        },
        onCellMouseDown: function() {
            $scope.mouseDownSelectionText = window.getSelection().toString();
        },
        onRowClicked: function(params) {
            const selection = window.getSelection().toString();
            if(selection === "" || selection === $scope.mouseDownSelectionText) {
                locationUtils.navigateToPath('/delivery-services/' + params.data.id + '?type=' + params.data.type);
                // Event is outside the digest cycle, so we need to trigger one.
                $scope.$apply();
            }
            $scope.mouseDownSelectionText = "";
        },
        onColumnResized: function(params) {
            localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
        },
        onFirstDataRendered: function(event) {
            try {
                const filterState = JSON.parse(localStorage.getItem(tableName + "_table_filters")) || {};
                // apply any filter provided to the controller
                Object.assign(filterState, filter);
                $scope.gridOptions.api.setFilterModel(filterState);
            } catch (e) {
                console.error("Failure to load stored filter state:", e);
            }

            $scope.gridOptions.api.addEventListener("filterChanged", function() {
                localStorage.setItem(tableName + "_table_filters", JSON.stringify($scope.gridOptions.api.getFilterModel()));
            });
        },
        onGridReady: function() {
            try { // need to create the show/hide column checkboxes and bind to the current visibility
                const colstates = JSON.parse(localStorage.getItem(tableName + "_table_columns"));
                if (colstates) {
                    if (!$scope.gridOptions.columnApi.setColumnState(colstates)) {
                        console.error("Failed to load stored column state: one or more columns not found");
                    }
                } else {
                    $scope.gridOptions.api.sizeColumnsToFit();
                }
            } catch (e) {
                console.error("Failure to retrieve required column info from localStorage (key=" + tableName + "_table_columns):", e);
            }

            try {
                const sortState = JSON.parse(localStorage.getItem(tableName + "_table_sort"));
                $scope.gridOptions.api.setSortModel(sortState);
            } catch (e) {
                console.error("Failure to load stored sort state:", e);
            }

            try {
                $scope.quickSearch = localStorage.getItem(tableName + "_quick_search");
                $scope.gridOptions.api.setQuickFilter($scope.quickSearch);
            } catch (e) {
                console.error("Failure to load stored quick search:", e);
            }

            try {
                const ps = localStorage.getItem(tableName + "_page_size");
                if (ps && ps > 0) {
                    $scope.pageSize = Number(ps);
                    $scope.gridOptions.api.paginationSetPageSize($scope.pageSize);
                }
            } catch (e) {
                console.error("Failure to load stored page size:", e);
            }
            
            try {
                const page = parseInt(localStorage.getItem(tableName + "_table_page"));
                const totalPages = $scope.gridOptions.api.paginationGetTotalPages();
                if (page !== undefined && page > 0 && page <= totalPages-1) {
                    $scope.gridOptions.api.paginationGoToPage(page);
                }
            } catch (e) {
                console.error("Failed to load stored page number:", e);
            }

            $scope.gridOptions.api.addEventListener("paginationChanged", function() {
                localStorage.setItem(tableName + "_table_page", $scope.gridOptions.api.paginationGetCurrentPage());
            });

            $scope.gridOptions.api.addEventListener("sortChanged", function() {
                localStorage.setItem(tableName + "_table_sort", JSON.stringify($scope.gridOptions.api.getSortModel()));
            });

            $scope.gridOptions.api.addEventListener("columnMoved", function() {
                localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
            });

            $scope.gridOptions.api.addEventListener("columnVisible", function() {
                $scope.gridOptions.api.sizeColumnsToFit();
                try {
                    const colStates = $scope.gridOptions.columnApi.getColumnState();
                    localStorage.setItem(tableName + "_table_columns", JSON.stringify(colStates));
                } catch (e) {
                    console.error("Failed to store column defs to local storage:", e);
                }
            });
        }
    };

    /** This is used to position the context menu under the cursor. */
    $scope.menuStyle = {
        left: 0,
        top: 0,
    };

    /** Controls whether or not the context menu is visible. */
    $scope.showMenu = false;

    /** Downloads the table as a CSV */
    $scope.exportCSV = function() {
        const params = {
            allColumns: true,
            fileName: "delivery-services.csv",
        };
        $scope.gridOptions.api.exportDataAsCsv(params);
    };

    $scope.onQuickSearchChanged = function() {
        $scope.gridOptions.api.setQuickFilter($scope.quickSearch);
        localStorage.setItem(tableName + "_quick_search", $scope.quickSearch);
    };

    $scope.onPageSizeChanged = function() {
        const value = Number($scope.pageSize);
        $scope.gridOptions.api.paginationSetPageSize(value);
        localStorage.setItem(tableName + "_page_size", value);
    };

    $scope.clearTableFilters = function() {
        // clear the quick search
        $scope.quickSearch = '';
        $scope.onQuickSearchChanged();
        // clear any column filters
        $scope.gridOptions.api.setFilterModel(null);
    };

    angular.element(document).ready(function () {
        // clicks outside the context menu will hide it
        $document.bind("click", function(e) {
            $scope.showMenu = false;
            e.stopPropagation();
            $scope.$apply();
        });
    });

};

TableDeliveryServicesController.$inject = ['tableName', 'deliveryServices', 'filter', '$anchorScroll', '$document', '$scope', '$state', '$location', '$uibModal', 'deliveryServiceService', 'deliveryServiceRequestService', 'dateUtils', 'deliveryServiceUtils', 'locationUtils', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = TableDeliveryServicesController;
