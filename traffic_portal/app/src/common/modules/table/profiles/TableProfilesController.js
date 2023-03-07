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
 * @typedef ProfileType @type {"ATS_PROFILE" | "TR_PROFILE" | "TM_PROFILE" | "TS_PROFILE" | "TP_PROFILE" | "INFLUXDB_PROFILE" | "RIAK_PROFILE" | "SPLUNK_PROFILE" | "DS_PROFILE" | "ORG_PROFILE" | "KAFKA_PROFILE" | "LOGSTASH_PROFILE" | "ES_PROFILE" | "UNK_PROFILE" | "GROVE_PROFILE"}
 */

/**
 * @typedef Profile
 * @property {number} cdn
 * @property {string} cdnName
 * @property {string} description
 * @property {number} id
 * @property {string} lastUpdated
 * @property {string} name
 * @property {boolean} routingDisabled
 * @property {ProfileType} type
 */

/**
 * @typedef ProfileGridParams
 * @property {Profile} data
 */

/**
 * The controller for the Profiles table view
 *
 * @param {Profile[]} profiles
 * @param {*} $scope
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/ProfileService")} profileService
 * @param {import("../../../models/MessageModel")} messageModel
 * @param {import("../../../service/utils/FileUtils")} fileUtils
 */
var TableProfilesController = function(profiles, $scope, $location, $uibModal, locationUtils, profileService, messageModel, fileUtils) {

	/**** Constants, scope data, etc. ****/

	/** The columns of the ag-grid table */
	$scope.columns = [
		{
			headerName: "CDN",
			field: "cdnName",
			hide: false,
			/**
			 * @param {ProfileGridParams} params
			 * @returns {string}
			 */
			tooltipValueGetter: params => `#${params.data.cdn}`
		},
		{
			headerName: "Description",
			field: "description",
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
		},
		{
			headerName: "Routing Disabled",
			field: "routingDisabled",
			hide: false
		},
		{
			headerName: "Type",
			field: "type",
			hide: false,
			/**
			 * @param {ProfileGridParams} params
			 * @returns {string}
			 */
			tooltipValueGetter: params => {
				switch(params.data.type) {
					case "ATS_PROFILE":
						return "Trafficserver Cache Server Profile";
					case "DS_PROFILE":
						return "Delivery Service Profile";
					case "ES_PROFILE":
						return "Elasticsearch Server Profile";
					case "GROVE_PROFILE":
						return "Grove Cache Server Profile";
					case "INFLUXDB_PROFILE":
						return "InfluxDB Server Profile";
					case "KAFKA_PROFILE":
						return "Kafka Server Profile";
					case "LOGSTASH_PROFILE":
						return "Logstash Server Profile";
					case "ORG_PROFILE":
						return "Origin Profile";
					case "RIAK_PROFILE":
						return "Traffic Vault Server Profile";
					case "SPLUNK_PROFILE":
						return "Splunk Server Profile";
					case "TM_PROFILE":
						return "Traffic Monitor Server Profile";
					case "TP_PROFILE":
						return "Traffic Portal Server Profile";
					case "TR_PROFILE":
						return "Traffic Router Server Profile";
					case "TS_PROFILE":
						return "Traffic Stats Server Profile";
					case "UNK_PROFILE":
						return "Other Profile";
				}
			}
		}
	];

	/**
	 * Opens a dialog that prompts the user to upload a Profile to be imported.
	 */
	function importProfile() {
		const params = {
			title: "Import Profile",
			message: "Drop Profile Here"
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/import/dialog.import.tpl.html",
			controller: "DialogImportController",
			size: "lg",
			resolve: {params}
		});
		modalInstance.result.then(
			importJSON => {
				profileService.importProfile(importJSON);
			}
		);
	};

	/**
	 * Opens a dialog that prompts the user to select two Profiles to compare,
	 * then navigates to the Profile comparison page for them (assuming it was
	 * not instead cancelled).
	 */
	function compareProfiles() {
		const params = {
			title: "Compare Profiles",
			message: "Please select 2 profiles to compare",
			labelFunction: item => `${item.name} (${item.type})`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/compare/dialog.compare.tpl.html",
			controller: "DialogCompareController",
			size: "md",
			resolve: {
				collection: profileService => profileService.getProfiles({ orderby: "name" }),
				params
			}
		});
		modalInstance.result.then(
			([a, b]) => {
				$location.path(`${$location.path()}/${a.id}/${b.id}/compare/diff`);
			}
		);
	};

	/** @type {import("../agGrid/CommonGridController").CGC.DropDownOption[]} */
	$scope.dropDownOptions = [
		{
			name: "createProfileMenuItem",
			href: "#!/profiles/new",
			text: "Create New Profile",
			type: 2
		},
		{type: 0},
		{
			onClick: importProfile,
			text: "Import Profile",
			type: 1
		},
		{
			onClick: compareProfiles,
			text: "Compare Profiles",
			type: 1
		}
	];

	/**
	 * Deletes a Profile after getting confirmation from the user.
	 *
	 * @param {Profile} profile
	 */
	function confirmDelete(profile) {
		const params = {
			title: `Delete Profile: ${profile.name}`,
			key: profile.name
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/delete/dialog.delete.tpl.html",
			controller: "DialogDeleteController",
			size: "md",
			resolve: {params}
		});
		modalInstance.result.then(
			async () => {
				const result = await profileService.deleteProfile(profile.id)
				messageModel.setMessages(result.alerts, false);
				$scope.refresh();
			}
		);
	};

	/**
	 * Clones the given Profile.
	 *
	 * @param {Profile} profile
	 */
	function cloneProfile(profile) {
		const params = {
			title: "Clone Profile",
			message: `You are about to clone the ${profile.name} profile. Your clone will have the same attributes and parameter assignments as the ${profile.name} profile.<br><br>Please enter a name for your cloned profile.`
		};
		const modalInstance = $uibModal.open({
			templateUrl: "common/modules/dialog/input/dialog.input.tpl.html",
			controller: "DialogInputController",
			size: "md",
			resolve: {params}
		});
		modalInstance.result.then(
			clonedProfileName => {
				profileService.cloneProfile(profile.name, clonedProfileName);
			}
		);
	};

	/**
	 * Downloads the given Profile.
	 *
	 * @param {Profile} profile
	 */
	async function exportProfile(profile) {
		const result = await profileService.exportProfile(profile.id)
		fileUtils.exportJSON(result, profile.name, "json");
	};

	/** @type {import("../agGrid/CommonGridController").CGC.ContextMenuOption[]} */
	$scope.contextMenuOptions = [
		{
			getHref: profile => `#!/profiles/${profile.id}`,
			getText: profile => `Open ${profile.name} in a new tab`,
			newTab: true,
			type: 2
		},
		{type: 0},
		{
			getHref: profile => `#!/profiles/${profile.id}`,
			text: "Edit",
			type: 2
		},
		{
			onClick: profile => confirmDelete(profile),
			text: "Delete",
			type: 1
		},
		{type: 0},
		{
			onClick: cloneProfile,
			text: "Clone Profile",
			type: 1
		},
		{
			onClick: exportProfile,
			text: "Export Profile",
			type: 1
		},
		{type: 0},
		{
			getHref: profile => `#!/profiles/${profile.id}/parameters`,
			text: "Manage Parameters",
			type: 2
		},
		{
			getHref: profile => `#!/servers?profileName=${profile.name}`,
			isDisabled: profile => profile.type === "DS_PROFILE",
			text: "View Servers",
			type: 2
		},
		{
			getHref: profile => `#!/delivery-services?profileName=${profile.name}`,
			isDisabled: profile => profile.type !== "DS_PROFILE",
			text: "View Delivery Services",
			type: 2
		}
	];

	/** Options, configuration, data and callbacks for the ag-grid table. */
	/** @type {import("../agGrid/CommonGridController").CGC.GridSettings} */
	$scope.gridOptions = {
		onRowClick: row => {
			locationUtils.navigateToPath(`/profiles/${row.data.id}`);
		}
	};

	$scope.defaultData = {
		cdn: -1,
		cdnName: "",
		description: "",
		id: -1,
		lastUpdated: new Date(),
		name: "",
		routingDisabled: true,
		type: "ATS_PROFILE"
	};

	$scope.profiles = profiles.map(
		profile =>
			({...profile, lastUpdated: new Date(profile.lastUpdated.replace(" ", "T").replace("+00", "Z"))})
	);
};

TableProfilesController.$inject = ["profiles", "$scope", "$location", "$uibModal", "locationUtils", "profileService", "messageModel", "fileUtils"];
module.exports = TableProfilesController;
