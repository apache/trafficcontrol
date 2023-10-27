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
 * @typedef CDN
 * @property {boolean} dnssecEnabled
 * @property {string} domainName
 * @property {number} id
 * @property {string} lastUpdated
 * @property {string} name
 */

/**
 * @param {CDN} cdn
 * @returns  {string}
 */
const getHref = cdn => `#!/cdns/${cdn.id}`;

/**
 * @param {CDN[]} cdns
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/CDNService")} cdnService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableCDNsController = function(cdns, $scope, $state, $uibModal, locationUtils, cdnService, messageModel) {

	/**** Constants, scope data, etc. ****/

	/** The columns of the ag-grid table */
	$scope.columns = [
		{
			headerName: "DNSSEC Enabled",
			field: "dnssecEnabled",
			hide: false
		},
		{
			headerName: "Domain",
			field: "domainName",
			hide: false
		},
		{
			headerName: "ID",
			field: "id",
			filter: "agNumberColumnFilter",
			hide: true
		},
		{
			headerName: "Last Updated",
			field: "lastUpdated",
			hide: true,
			filter: "agDateColumnFilter",
			relative: true
		},
		{
			headerName: "Name",
			field: "name",
			hide: false,
		}
	];

	/** @type {import("../agGrid/CommonGridController").CGC.DropDownOption[]} */
	$scope.dropDownOptions = [{
		name: "createCDNMenuItem",
		href: "#!/cdns/new",
		text: "Create New CDN",
		type: 2
	}, {
		type: 0
	}, {
		onClick: function (entry) {
			$scope.confirmQueueServerUpdates(entry);
		},
		text: "Queue CDN Server Updates",
		type: 1
	}, {
		onClick: function (entry) {
			$scope.confirmClearServerUpdates(entry);
		},
		text: "Clear CDN Server Updates",
		type: 1
	}];

	/** Reloads all resolved data for the view. */
	$scope.refresh = function() {
		$state.reload();
	};

	/**
	 * Deletes a CDN if confirmation is given.
	 * @param {CDN} cdn
	 */
	function confirmDelete(cdn) {
		const params = {
			title: `Delete CDN: ${cdn.name}`,
			key: cdn.name
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/delete/dialog.delete.tpl.html",
			controller: "DialogDeleteController",
			size: "md",
			resolve: {params}
		});
		modalInstance.result.then(() => {
			cdnService.deleteCDN(cdn.id).then(
				result => {
					messageModel.setMessages(result.alerts, false);
					$scope.refresh();
				}
			);
		}).catch(
			e => console.error("failed to delete CDN:", e)
		);
	};

	/**
	 * Queues servers updates on a CDN if confirmation is given.
	 * @param {CDN} cdn
	 */
	function confirmQueueServerUpdates(cdn) {
		const params = {
			title: `Queue Server Updates: ${cdn.name}`,
			message: `Are you sure you want to queue server updates for all ${cdn.name} servers?`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/confirm/dialog.confirm.tpl.html",
			controller: "DialogConfirmController",
			size: "md",
			resolve: {params}
		});
		modalInstance.result.then(() => cdnService.queueServerUpdates(cdn.id));
	};

	/**
	 * Clears servers updates on a CDN if confirmation is given.
	 * @param {CDN} cdn
	 */
	 function confirmClearServerUpdates(cdn) {
		const params = {
			title: `Clear Server Updates: ${cdn.name}`,
			message: `Are you sure you want to clear server updates for all ${cdn.name} servers?`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/confirm/dialog.confirm.tpl.html",
			controller: "DialogConfirmController",
			size: "md",
			resolve: {params}
		});
		modalInstance.result.then(function() {
			cdnService.clearServerUpdates(cdn.id);
		});
	};

	/** @type {import("../agGrid/CommonGridController").CGC.ContextMenuOption[]} */
	$scope.contextMenuOptions = [
		{
			getHref,
			getText: cdn => `Open ${cdn.name} in a new tab`,
			newTab: true,
			type: 2
		},
		{type: 0},
		{
			getHref,
			text: "Edit",
			type: 2
		},
		{
			onClick: cdn => confirmDelete(cdn),
			text: "Delete",
			type: 1
		},
		{type: 0},
		{
			getHref: cdn => `#!/cdns/${cdn.id}/config/changes`,
			text: "Diff Snapshot",
			type: 2
		},
		{type: 0},
		{
			onClick: cdn => confirmQueueServerUpdates(cdn),
			text: "Queue Server Updates",
			type: 1
		},
		{
			onClick: cdn => confirmClearServerUpdates(cdn),
			text: "Clear Server Updates",
			type: 1
		},
		{type: 0},
		{
			getHref: cdn => `#!/cdns/${cdn.id}/dnssec-keys`,
			text: "Manage DNSSEC Keys",
			type: 2
		},
		{
			getHref: cdn => `#!/cdns/${cdn.id}/federations`,
			text: "Manage Federations",
			type: 2
		},
		{
			getHref: cdn => `#!/cdns/${cdn.id}/delivery-services`,
			text: "Manage Delivery Services",
			type: 2
		},
		{
			getHref: cdn => `#!/profiles?cdnName=${encodeURIComponent(cdn.name)}`,
			text: "Manage Profiles",
			type: 2
		},
		{
			getHref: cdn => `#!/cdns/${cdn.id}/servers`,
			text: "Manage Servers",
			type: 2
		},
		{
			getHref: cdn => `#!/cdns/${cdn.id}/notifications`,
			text: "Manage Notifications",
			type: 2
		}
	];

	/** Options, configuration, data and callbacks for the ag-grid table. */
	/** @type {import("../agGrid/CommonGridController").CGC.GridSettings} */
	$scope.gridOptions = {
		onRowClick: function(row) {
			locationUtils.navigateToPath(`/cdns/${row.data.id}`);
		}
	};

	$scope.defaultData = {
		dnssecEnabled: false,
		domainName: "",
		name: ""
	};

	$scope.cdns = cdns.map(
		cdn => ({...cdn, lastUpdated: new Date(cdn.lastUpdated.replace(" ", "T").replace("+00", "Z"))})
	);

};

TableCDNsController.$inject = ["cdns", "$scope", '$state', "$uibModal", "locationUtils", "cdnService", "messageModel"];
module.exports = TableCDNsController;
