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

var TableDeliveryServicesRequestsController = function (tableName, dsRequests, $scope, $state, $uibModal, $anchorScroll, $q, $location, $document, dateUtils, locationUtils, typeService, deliveryServiceService, deliveryServiceRequestService, messageModel, propertiesModel, userModel) {

	/**
	 * Gets value to display a default tooltip.
	 */
	function defaultTooltip(params) {
		return params.value;
	}

	/**
	 * Formats the contents of a 'createdAt' and 'lastUpdated' column cell as "relative to now".
	 */
	function dateCellFormatter(params) {
		return params.value ? dateUtils.getRelativeTime(params.value) : params.value;
	}

	columns = [
		{
			headerName: "Delivery Service",
			field: "deliveryService.xmlId",
			hide: false
		},
		{
			headerName: "Type",
			field: "changeType",
			hide: false
		},
		{
			headerName: "Status",
			field: "status",
			hide: false
		},
		{
			headerName: "Author",
			field: "author",
			hide: false
		},
		{
			headerName: "Assignee",
			field: "assignee",
			hide: false
		},
		{
			headerName: "Last Edited By",
			field: "lastEditedBy",
			hide: true
		},
		{
			headerName: "Last Updated",
			field: "lastUpdated",
			hide: true,
			filter: "agDateColumnFilter",
			tooltipValueGetter: () => dateCellFormatter,
			valueFormatter: dateCellFormatter
		},
		{
			headerName: "Created",
			field: "createdAt",
			hide: false,
			filter: "agDateColumnFilter",
			tooltipValueGetter: () => dateCellFormatter,
			valueFormatter: dateCellFormatter
		}
	];

	var createComment = function (request, placeholder) {
		var params = {
			title: 'Add Comment',
			placeholder: placeholder,
			text: null
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/textarea/dialog.textarea.tpl.html',
			controller: 'DialogTextareaController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function (commentValue) {
			var comment = {
				deliveryServiceRequestId: request.id,
				value: commentValue
			};
			deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
		}, function () {
			// do nothing
		});
	};

	$scope.DRAFT = 0;
	$scope.SUBMITTED = 1;
	$scope.REJECTED = 2;
	$scope.PENDING = 3;
	$scope.COMPLETE = 4;

	/** All of the ds requests - createdAt and lastUpdated fields converted to actual Date */
	$scope.dsRequests = dsRequests.map(
		function(x) {
			x.createdAt = x.createdAt ? new Date(x.createdAt.replace("+00", "Z")) : x.createdAt;
			x.lastUpdated = x.lastUpdated ? new Date(x.lastUpdated.replace("+00", "Z")) : x.lastUpdated;
		});

	$scope.quickSearch = '';

	$scope.pageSize = 100;

	$scope.mouseDownSelectionText = "";

	/** Options, configuration, data and callbacks for the ag-grid table. */
	$scope.gridOptions = {
		onCellMouseDown: function() {
			$scope.mouseDownSelectionText = window.getSelection().toString();
		},
		onRowClicked: function(params) {
			const selection = window.getSelection().toString();
			if(selection === "" || selection === $scope.mouseDownSelectionText) {
				let path = '/delivery-service-requests/' + params.data.id + '?type=';
				typeService.getType(params.data.deliveryService.typeId)
					.then(function (result) {
						path += result.name;
						locationUtils.navigateToPath(path);
					});
				$scope.$apply();
			}
			$scope.mouseDownSelectionText = "";
		},
		columnDefs: columns,
		enableCellTextSelection: true,
		suppressMenuHide: true,
		multiSortKey: 'ctrl',
		alwaysShowVerticalScroll: true,
		defaultColDef: {
			filter: true,
			sortable: true,
			resizable: true,
			tooltipValueGetter: () => defaultTooltip
		},
		rowClassRules: {
			'draft-request': function(params) {
				return params.data.status === 'draft';
			},
			'submitted-request': function(params) {
				return params.data.status === 'submitted';
			},
			'pending-request': function(params) {
				return params.data.status === 'pending';
			},
			'completed-request': function(params) {
				return params.data.status === 'complete';
			},
			'rejected-request': function(params) {
				return params.data.status === 'rejected';
			},
		},
		rowData: dsRequests,
		pagination: true,
		paginationPageSize: $scope.pageSize,
		rowBuffer: 0,
		onColumnResized: function(params) {
			localStorage.setItem(tableName + "_table_columns", JSON.stringify($scope.gridOptions.columnApi.getColumnState()));
		},
		tooltipShowDelay: 500,
		allowContextMenuWithControlKey: true,
		preventDefaultOnContextMenu: true,
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
			$scope.request = params.data;
			$scope.$apply();
		},
		onColumnVisible: function(params) {
			if (params.visible){
				return;
			}
			for (let column of params.columns) {
				if (column.filterActive) {
					const filterModel = $scope.gridOptions.api.getFilterModel();
					if (column.colId in filterModel) {
						delete filterModel[column.colId];
						$scope.gridOptions.api.setFilterModel(filterModel);
					}
				}
			}
		},
        onFirstDataRendered: function() {
			try {
				const filterState = JSON.parse(localStorage.getItem(tableName + "_table_filters")) || {};
				$scope.gridOptions.api.setFilterModel(filterState);
			} catch (e) {
				console.error("Failure to load stored filter state:", e);
			}

			$scope.gridOptions.api.addEventListener("filterChanged", function() {
				localStorage.setItem(tableName + "_table_filters", JSON.stringify($scope.gridOptions.api.getFilterModel()));
			});
		},
		onGridReady: function() {
			try {
				// need to create the show/hide column checkboxes and bind to the current visibility
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
		},
		colResizeDefault: "shift"
	};

	/** This is used to position the context menu under the cursor. */
	$scope.menuStyle = {
		left: 0,
		top: 0,
	};

	/** Toggles the visibility of a column that has the ID provided as 'col'. */
	$scope.toggleVisibility = function(col) {
		const visible = $scope.gridOptions.columnApi.getColumn(col).isVisible();
		$scope.gridOptions.columnApi.setColumnVisible(col, !visible);
	};

	/** Downloads the table as a CSV */
	$scope.exportCSV = function() {
		const params = {
			allColumns: true,
			fileName: "delivery_service_requests.csv",
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

	$scope.fulfillable = function (request) {
		return request && request.status == 'submitted';
	};

	$scope.rejectable = function (request) {
		return request && request.status == 'submitted';
	};

	$scope.completeable = function (request) {
		return request && request.status == 'pending';
	};

	$scope.open = function (request) {
		return request && (request.status == 'draft' || request.status == 'submitted');
	};

	$scope.assignRequest = function (request, assign, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: 'Assign Delivery Service Request',
			message: (assign) ? 'Are you sure you want to assign this delivery service request to yourself?' : 'Are you sure you want to unassign this delivery service request?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function () {
			var assigneeId = (assign) ? userModel.user.id : null;
			deliveryServiceRequestService.assignDeliveryServiceRequest(request.id, assigneeId).then(function () {
				$scope.refresh();
				if (assign) {
					messageModel.setMessages([ { level: 'success', text: 'Delivery service request was assigned' } ], false);
				} else {
					messageModel.setMessages([ { level: 'success', text: 'Delivery service request was unassigned' } ], false);
				}
			});
		}, function () {
			// do nothing
		});
	};

	$scope.editStatus = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: "Edit Delivery Service Request Status",
			message: 'Please select the appropriate status for this request.'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				},
				collection: function () {
					return [
						{id: $scope.DRAFT, name: 'Save as Draft'},
						{id: $scope.SUBMITTED, name: 'Submit for Review / Deployment'}
					];
				}
			}
		});
		modalInstance.result.then(function (action) {
			var status = (action.id == $scope.DRAFT) ? 'draft' : 'submitted';
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, status).then(function () {
				$scope.refresh();
				messageModel.setMessages([ { level: 'success', text: 'Delivery service request status was updated' } ], false);
			});
		}, function () {
			// do nothing
		});
	};

	$scope.rejectRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		// only the user assigned to the request can mark it as rejected (unless the user has override capabilities)
		if ((request.assigneeId != userModel.user.id) && (userModel.user.roleName != propertiesModel.properties.dsRequests.overrideRole)) {
			messageModel.setMessages([{
				level: 'error',
				text: 'Only the assignee can mark a delivery service request as rejected'
			}], false);
			$anchorScroll(); // scrolls window to top
			return;
		}

		var params = {
			title: 'Reject Delivery Service Request',
			message: 'Are you sure you want to reject this delivery service request?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function () {
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, 'rejected').then(
				function () {
					$scope.refresh();
					messageModel.setMessages([ { level: 'success', text: 'Delivery service request was rejected' } ], false);
					createComment(request, 'Enter rejection reason...');
				});
		}, function () {
			// do nothing
		});
	};

	$scope.completeRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		// only the user assigned to the request can mark it as complete (unless the user has override capabilities)
		if ((request.assigneeId != userModel.user.id) && (userModel.user.roleName != propertiesModel.properties.dsRequests.overrideRole)) {
			messageModel.setMessages([{
				level: 'error',
				text: 'Only the assignee can mark a delivery service request as complete'
			}], false);
			$anchorScroll(); // scrolls window to top
			return;
		}

		var params = {
			title: 'Complete Delivery Service Request',
			message: 'Are you sure you want to mark this delivery service request as complete?'
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function () {
			deliveryServiceRequestService.updateDeliveryServiceRequestStatus(request.id, 'complete').then(function () {
				$scope.refresh();
				messageModel.setMessages([ { level: 'success', text: 'Delivery service request marked as complete' } ], false);
				createComment(request, 'Enter comment...');
			});
		}, function () {
			// do nothing
		});
	};

	$scope.deleteRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var params = {
			title: 'Delete the ' + request.deliveryService.xmlId + ' ' + request.changeType + ' request?',
			key: request.deliveryService.xmlId + ' ' + request.changeType + ' request'
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
		modalInstance.result.then(function () {
			deliveryServiceRequestService.deleteDeliveryServiceRequest(request.id).then(function () {
				messageModel.setMessages([{level: 'success', text: 'Delivery service request was deleted'}], false);
				$scope.refresh();
			});
		}, function () {
			// do nothing
		});
	};

	$scope.fulfillRequest = function (request, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		var path = '/delivery-service-requests/' + request.id + '?type=';
		typeService.getType(request.deliveryService.typeId)
			.then(function (result) {
				path += result.name;
				locationUtils.navigateToPath(path);
			});
	};

	$scope.refresh = function () {
		$state.reload(); // reloads all the resolves for the view
	};

	/**** Initialization code, including loading user columns from localstorage ****/
	angular.element(document).ready(function () {

		// clicks outside the context menu will hide it
		$document.bind("click", function(e) {
			$scope.showMenu = false;
			e.stopPropagation();
			$scope.$apply();
		});
	});

};

TableDeliveryServicesRequestsController.$inject = ['tableName', 'dsRequests', '$scope', '$state', '$uibModal', '$anchorScroll', '$q', '$location', '$document', 'dateUtils', 'locationUtils', 'typeService', 'deliveryServiceService', 'deliveryServiceRequestService', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = TableDeliveryServicesRequestsController;
