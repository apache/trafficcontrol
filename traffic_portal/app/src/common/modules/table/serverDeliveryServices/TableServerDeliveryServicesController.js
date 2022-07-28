/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

function TableServerDeliveryServicesController(server, deliveryServices, filter, $controller, $scope, $uibModal, locationUtils, serverUtils, deliveryServiceService, serverService) {

	// extends the TableDeliveryServicesController to inherit common methods
	angular.extend(this, $controller("TableDeliveryServicesController", { tableName: "serverDS", deliveryServices, filter, $scope }));

	server = Array.isArray(server) ? server[0] : server;

	$scope.breadCrumbs = [
		{
			href: "#!/servers",
			text: "Servers"
		},
		{
			href: `#!/servers/${server.id}`,
			text: server.hostName
		},
		{
			text: "Delivery Services"
		}
	];

	/**
	 * Removes the assignment of a Delivery Service to the table's server.
	 *
	 * @param {number} dsId The ID of the Delivery Service being removed.
	 */
	async function removeDeliveryService(dsId) {
		await deliveryServiceService.deleteDeliveryServiceServer(dsId, $scope.server.id);
		$scope.refresh();
	};

	$scope.dropDownOptions = [
		{
			onClick: cloneDsAssignments,
			text: "Clone Delivery Service Assignments",
			type: 1
		}
	];

	if (serverUtils.isEdge(server) || serverUtils.isOrigin(server)) {
		$scope.dropDownOptions.unshift({
			onClick: selectDeliveryServices,
			text: "Assign Delivery Services",
			type: 1
		});
	}

	/**
	 * Asks a user for confirmation before removing a Delivery Service assignment
	 * from the table's server.
	 *
	 * @param {{id: number; xmlId: string}} ds The Delivery Service being removed.
	 */
	async function confirmRemoveDS(ds) {
		const params = {
			title: "Remove Delivery Service from Server?",
			message: `Are you sure you want to remove ${ds.xmlId} from this server?`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/confirm/dialog.confirm.tpl.html",
			controller: "DialogConfirmController",
			size: "md",
			resolve: { params }
		});
		try {
			await modalInstance.result;
			removeDeliveryService(ds.id);
		} catch {
			// do nothing
		}
	};

	$scope.contextMenuOptions.splice(1, 0, {
		onClick: confirmRemoveDS,
		getText: ds => `Remove ${ds.xmlId}`,
		type: 1
	});


	/**
	 * Removes a Delivery Service from this server after obtaining user
	 * confirmation.
	 *
	 * @param {{readonly id: number; readonly xmlId: string}} ds
	 */
	async function confirmRemoveDS(ds) {
		const params = {
			title: "Remove Delivery Service from Server?",
			message: `Are you sure you want to remove ${ds.xmlId} from this server?`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/confirm/dialog.confirm.tpl.html",
			controller: "DialogConfirmController",
			size: "md",
			resolve: { params }
		});
		try {
			await modalInstance.result;
			await deliveryServiceService.deleteDeliveryServiceServer(ds.id, server.id)
			$scope.refresh();
		} catch {
			// do nothing
		}
	}

	/**
	 * Opens a dialog that allows the user to select another cache server to
	 * which to assign the same Delivery Services assigned to this server (
	 * overriding any existing assignments for the selected cache server).
	 */
	async function cloneDsAssignments() {
		const params = {
			title: "Clone Delivery Service Assignments",
			message: `Please select another ${server.type} cache server to assign these ${deliveryServices.length} Delivery Services to.` +
				"<br/>" +
				"<br/>" +
				`<strong style='text-transform: uppercase'>Warning: this cannot be undone</strong> - Any Delivery Services currently assigned to the selected cache server will be lost and replaced with these ${deliveryServices.length} Delivery Service assignments.`,
			labelFunction: item => `${item.hostName}.${item.domainName}`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/select/dialog.select.tpl.html",
			controller: "DialogSelectController",
			size: "md",
			resolve: {
				params,
				collection: async () => {
					const opts = {
						type: server.type,
						orderby: "hostName",
						cdn: server.cdnId
					};
					const ss = await serverService.getServers(opts);
					return ss.filter(s => s.id !== server.id);
				}
			}
		});

		let selectedServer;
		try {
			selectedServer = await modalInstance.result;
		} catch {
			return;
		}
		const dsIds = deliveryServices.map(ds=>ds.id);
		await serverService.assignDeliveryServices(selectedServer, dsIds, true, true);
		locationUtils.navigateToPath(`/servers/${selectedServer.id}/delivery-services`);
	};

	/**
	 * Opens a dialog that allows the user to modify the server's Delivery
	 * Service assignments.
	 */
	async function selectDeliveryServices() {
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/table/serverDeliveryServices/table.assignDeliveryServices.tpl.html",
			controller: "TableAssignDeliveryServicesController",
			size: "lg",
			resolve: {
				server: () => server,
				deliveryServices: deliveryServiceService => deliveryServiceService.getDeliveryServices({ cdn: server.cdnId }),
				assignedDeliveryServices: () => deliveryServices
			}
		});

		let selectedDSIDs;
		try {
			selectedDSIDs = await modalInstance.result;
		} catch {
			return;
		}
		await serverService.assignDeliveryServices(server, selectedDSIDs, true, false);
		$scope.refresh();
	};

	$scope.contextMenuOptions.splice(1, 0, {
		getText: ds => `Remove ${ds.xmlId}`,
		onClick: confirmRemoveDS,
		type: 1
	});

	$scope.dropDownOptions = [
		{
			onClick: selectDeliveryServices,
			text: "Assign Delivery Services",
			type: 1
		},
		{
			onClick: cloneDsAssignments,
			text: "Clone Delivery Service Assignments",
			type: 1
		}
	];
};

TableServerDeliveryServicesController.$inject = ["server", "deliveryServices", "filter", "$controller", "$scope", "$uibModal", "locationUtils", "serverUtils", "deliveryServiceService", "serverService"];
module.exports = TableServerDeliveryServicesController;
