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
 * This is the controller for (almost) all tables of Delivery Services.
 *
 * @param {string} tableName
 * @param {import("../../../api/DeliveryServiceService").DeliveryService[]} deliveryServices
 * @param {{deliveryService: string; targets: {deliveryService: string}[]}[]} steeringTargets
 * @param {import("angular").IAnchorScrollService} $anchorScroll
 * @param {*} $scope
 * @param {*} $state
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/DeliveryServiceUtils")} deliveryServiceUtils
 * @param {import("../../../models/MessageModel")} messageModel
 * @param {import("../../../models/PropertiesModel")} propertiesModel
 * @param {import("../../../models/UserModel")} userModel
 */
function TableDeliveryServicesController(tableName, deliveryServices, steeringTargets, $anchorScroll, $scope, $state, $location, $uibModal, deliveryServiceService, deliveryServiceRequestService, deliveryServiceUtils, locationUtils, messageModel, propertiesModel, userModel) {
	$scope.tableName = tableName;

	/** The columns of the ag-grid table */
	$scope.columns = [
		{
			headerName: "Active",
			field: "active",
			hide: false,
			valueGetter: ({data}) => {
				if (propertiesModel.properties.deliveryServices?.exposeInactive || data.active === "ACTIVE") {
					return data.active;
				}
				return "INACTIVE";
			}
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
			valueFormatter: params => params.data.consistentHashQueryParams.join(', '),
			tooltipValueGetter: params => params.data.consistentHashQueryParams.join(', ')
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
			valueGetter: params => deliveryServiceUtils.geoLimits[params.data.geoLimit],
			tooltipValueGetter: params => deliveryServiceUtils.geoLimits[params.data.geoLimit]
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
			valueGetter: params => deliveryServiceUtils.geoProviders[params.data.geoProvider],
			tooltipValueGetter: params => deliveryServiceUtils.geoProviders[params.data.geoProvider]
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
			valueGetter: params => deliveryServiceUtils.protocols[params.data.protocol],
			tooltipValueGetter: params => deliveryServiceUtils.protocols[params.data.protocol]
		},
		{
			headerName: "Qstring Handling",
			field: "qstringIgnore",
			hide: true,
			valueGetter: params => deliveryServiceUtils.qstrings[params.data.qstringIgnore],
			tooltipValueGetter: params => deliveryServiceUtils.qstrings[params.data.qstringIgnore]
		},
		{
			headerName: "Range Request Handling",
			field: "rangeRequestHandling",
			hide: true,
			valueGetter: params => deliveryServiceUtils.rrhs[params.data.rangeRequestHandling],
			tooltipValueGetter: params => deliveryServiceUtils.rrhs[params.data.rangeRequestHandling]
		},
		{
			headerName: "Regex Remap Expression",
			field: "regexRemap",
			hide: true
		},
		{
			headerName: "Regional Max Origin Connections",
			field: "regional",
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
			headerName: "Required Capability(ies)",
			field: "requiredCapabilities",
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
			headerName: "Target For",
			field: "isTargetFor",
			hide: true,
			valueGetter: params => params.data.isTargetsFor,
			tooltipValueGetter: params => `Steering targets for: ${params.data.isTargetsFor.join(", ")}`
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

	let dsRequestsEnabled = propertiesModel.properties?.dsRequests?.enabled;

	let showCustomCharts = propertiesModel.properties.deliveryServices?.charts.customLink.show;

	/**
	 * @param {string} typeName
	 */
	function createDeliveryService(typeName) {
		locationUtils.navigateToPath(`/delivery-services/new?dsType=${typeName}`);
	}

	/**
	 * Opens a dialog that the user uses to clone the given Delivery Service.
	 *
	 * @param {{readonly id: number; readonly xmlId: string;}} ds
	 */
	async function clone(ds) {
		const params = {
			title: `Clone Delivery Service: ${ds.xmlId}`,
			message: "Please select a <a href='https://traffic-control-cdn.readthedocs.io/en/latest/overview/delivery_services.html#ds-types' target='_blank'>content routing category</a> for the clone"
		};

		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/select/dialog.select.tpl.html",
			controller: "DialogSelectController",
			size: "md",
			resolve: {
				params,
				// the following represent the 4 categories of delivery services
				// the ids are arbitrary but the dialog.select dropdown needs them
				collection: () => [
					{ id: 1, name: "ANY_MAP" },
					{ id: 2, name: "DNS" },
					{ id: 3, name: "HTTP" },
					{ id: 4, name: "STEERING" }
				]
			}
		});

		const {name} = await modalInstance.result;
		locationUtils.navigateToPath(`/delivery-services/${ds.id}/clone?dsType=${name}`);
	}

	/**
	 * Opens a dialog asking the user for confirmation before deleting the given
	 * Delivery Service.
	 *
	 * @param {import("../../../api/DeliveryServiceService").DeliveryService} ds
	 */
	async function confirmDelete(ds) {
		const params = {
			title: `Delete Delivery Service: ${ds.xmlId}`,
			key: ds.xmlId
		};

		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/delete/dialog.delete.tpl.html",
			controller: "DialogDeleteController",
			size: "md",
			resolve: { params }
		});
		try {
			await modalInstance.result;
			if (dsRequestsEnabled) {
				return createDeliveryServiceDeleteRequest(ds);
			}
			try {
				await deliveryServiceService.deleteDeliveryService(ds);
				messageModel.setMessages([ { level: "success", text: `Delivery service [ ${ds.xmlId} ] deleted` } ], false);
				$scope.refresh();
			} catch (fault) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages(fault.data.alerts, false);
			}
		} catch {}
	}

	/**
	 * Creates a new DSR to delete the given Delivery Service.
	 *
	 * @param {import("../../../api/DeliveryServiceService").DeliveryService} ds
	 */
	async function createDeliveryServiceDeleteRequest(ds) {
		const params = {
			title: "Delivery Service Delete Request",
			message: "All delivery service deletions must be reviewed."
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/deliveryServiceRequest/dialog.deliveryServiceRequest.tpl.html",
			controller: "DialogDeliveryServiceRequestController",
			size: "md",
			resolve: {
				params,
				statuses: () => {
					const statuses = [
						{ id: $scope.DRAFT, name: "Save Request as Draft" },
						{ id: $scope.SUBMITTED, name: "Submit Request for Review and Deployment" }
					];
					if (userModel.user.role === propertiesModel.properties?.dsRequests?.overrideRole) {
						statuses.push({ id: $scope.COMPLETE, name: "Fulfill Request Immediately" });
					}
					return statuses;
				}
			}
		});
		const options = await modalInstance.result;
		let status = 'draft';
		if (options.status.id == $scope.SUBMITTED || options.status.id == $scope.COMPLETE) {
			status = 'submitted';
		}

		const dsRequest = {
			changeType: 'delete',
			status: status,
			original: ds
		};

		// if the user chooses to complete/fulfill the delete request immediately, the ds will be deleted and behind the
		// scenes a delivery service request will be created and marked as complete
		if (options.status.id == $scope.COMPLETE) {
			try {
				// first delete the ds
				await deliveryServiceService.deleteDeliveryService(ds);
				const response = await deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);
				const comment = {
					deliveryServiceRequestId: response.id,
					value: options.comment
				};
				// then create the ds request comment
				await deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
				const promises = [];
				// assign the ds request
				promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.username));
				// set the status to 'complete'
				promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, "complete"));
				// and finally refresh the delivery services table
				messageModel.setMessages([ { level: "success", text: `Delivery service [ ${ds.xmlId} ] deleted` } ], false);
				$scope.refresh();
			} catch (fault) {
				$anchorScroll(); // scrolls window to top
				messageModel.setMessages(fault.data.alerts, false);
			}
		} else {
			const response = await deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest);
			const comment = {
				deliveryServiceRequestId: response.id,
				value: options.comment
			};
			await deliveryServiceRequestService.createDeliveryServiceRequestComment(comment);
			const xmlId = (dsRequest.requested) ? dsRequest.requested.xmlId : dsRequest.original.xmlId;
			messageModel.setMessages([ { level: "success", text: `Created request to ${dsRequest.changeType} the ${xmlId} delivery service` } ], true);
			locationUtils.navigateToPath("/delivery-service-requests");
		}
	}


	$scope.DRAFT = 0;
	$scope.SUBMITTED = 1;
	$scope.REJECTED = 2;
	$scope.PENDING = 3;
	$scope.COMPLETE = 4;

	/**
	 * @deprecated This should instead just be an ng-href.
	 * @param {{readonly id: number; type: string; xmlId: string}} ds
	 */
	function viewCharts(ds) {
		if (showCustomCharts) {
			deliveryServiceUtils.openCharts(ds);
		} else {
			locationUtils.navigateToPath(`/delivery-services/${ds.id}/charts?dsType=${ds.type}`);
		}
	}

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	async function selectDSType() {
		const params = {
			title: "Create Delivery Service",
			message: "Please select a <a href='https://traffic-control-cdn.readthedocs.io/en/latest/overview/delivery_services.html#ds-types' target='_blank'>content routing category</a>"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
			controller: 'DialogSelectController',
			size: 'md',
			resolve: {
				params,
				// the following represent the 4 categories of delivery services
				// the ids are arbitrary but the dialog.select dropdown needs them
				collection: () => [
					{ id: 1, name: 'ANY_MAP' },
					{ id: 2, name: 'DNS' },
					{ id: 3, name: 'HTTP' },
					{ id: 4, name: 'STEERING' }
				]
			}
		});
		try {
			const type = await modalInstance.result;
			createDeliveryService(type.name);
		} catch {
			// do nothing
		}
	}

	this.$onInit = function() {
		const xmlIds = [];
		for(const ds of deliveryServices) {
			xmlIds.push(ds.xmlId);
		}
		const dsTargets = deliveryServiceUtils.getSteeringTargetsForDS(xmlIds, steeringTargets);
		/** All the delivery services - lastUpdated fields converted to actual Dates */
		$scope.deliveryServices = deliveryServices.map(
			x => ({...x, isTargetsFor: Array.from(dsTargets[x.xmlId]), lastUpdated: x.lastUpdated ? new Date(x.lastUpdated.replace("+00", "Z").replace(" ", "T")) : x.lastUpdated})
		);
	}

	async function compareDSs() {
		const params = {
			title: "Compare Delivery Services",
			message: "Please select 2 delivery services to compare",
			label: "xmlId"
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/compare/dialog.compare.tpl.html',
			controller: 'DialogCompareController',
			size: 'md',
			resolve: {
				params,
				collection: deliveryServiceService => deliveryServiceService.getDeliveryServices()
			}
		});
		try {
			const dss = await modalInstance.result;
			$location.path(`${$location.path()}/compare/${dss[0].id}/${dss[1].id}`);
		} catch {
			// do nothing
		}
	}

	$scope.options = {
		onRowClick: params => {
			const selection = window.getSelection()?.toString();
			if(!selection) {
				locationUtils.navigateToPath(`/delivery-services/${params.data.id}?dsType=${params.data.type}`);
				// Event is outside the digest cycle, so we need to trigger one.
				$scope.$apply();
			}
		}
	};

	$scope.dropDownOptions = [
		{
			onClick: selectDSType,
			text: "Create Delivery Service",
			type: 1
		},
		{
			onClick: compareDSs,
			text: "Compare Delivery Services",
			type: 1
		}
	];

	$scope.contextMenuOptions = [
		{
			getHref: ds => `#!/delivery-services/${ds.id}?dsType=${ds.type}`,
			getText: ds => `Open ${ds.xmlId} in a new tab`,
			newTab: true,
			type: 2
		},
		{type: 0},
		{
			getHref: ds => `#!/delivery-services/${ds.id}?dsType=${ds.type}`,
			text: "Edit",
			type: 2
		},
		{
			onClick: clone,
			text: "Clone",
			type: 1
		},
		{
			onClick: confirmDelete,
			text: "Delete",
			type: 1
		},
		{type: 0},
		{
			onClick: viewCharts,
			text: "View Charts",
			type: 1
		},
		{type: 0},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/ssl-keys?dsType=${ds.type}`,
			text: "Manage SSL Keys",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/url-sig-keys?dsType=${ds.type}`,
			text: "Manage URL Sig Keys",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/uri-signing-keys?dsType=${ds.type}`,
			text: "Manage URI Signing Keys",
			type: 2
		},
		{ type: 0 },
		{
			getHref: ds => `#!/delivery-services/${ds.id}/jobs?dsType=${ds.type}`,
			text: "Manage Invalidation Requests",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/origins?dsType=${ds.type}`,
			text: "Manage Origins",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/regexes?dsType=${ds.type}`,
			text: "Manage Regexes",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/servers?dsType=${ds.type}`,
			text: "Manage Servers",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/targets?dsType=${ds.type}`,
			text: "Manage Targets",
			type: 2
		},
		{
			getHref: ds => `#!/delivery-services/${ds.id}/static-dns-entries?dsType=${ds.type}`,
			text: "Manage Static DNS Entries",
			type: 2
		},

	];
}

TableDeliveryServicesController.$inject = ["tableName", "deliveryServices", "steeringTargets", "$anchorScroll", "$scope", "$state", "$location", "$uibModal", "deliveryServiceService", "deliveryServiceRequestService", "deliveryServiceUtils", "locationUtils", "messageModel", "propertiesModel", "userModel"];
module.exports = TableDeliveryServicesController;
