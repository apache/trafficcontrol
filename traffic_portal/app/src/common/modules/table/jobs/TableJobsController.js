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
 * @param {*} jobs
 * @param {*} $scope
 * @param {*} $state
 * @param {import("../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../api/JobService")} jobService
 * @param {import("../../../models/MessageModel")} messageModel
 */
var TableJobsController = function(jobs, $scope, $state, $uibModal, locationUtils, jobService, messageModel) {

	/** @type {import("../agGrid/CommonGridController").CGC.ColumnDefinition} */
	$scope.columns = [
		{
			headerName: "Delivery Service",
			field: "deliveryService",
			hide: false
		},
		{
			headerName: "Asset URL",
			field: "assetUrl",
			hide: false
		},
		{
			headerName: "TTL (Hours)",
			field: "ttlHours",
			hide: false
		},
		{
			headerName: "Start (UTC)",
			field: "startTime",
			hide: false,
			filter: "agDateColumnFilter",
		},
		{
			headerName: "Expires (UTC)",
			field: "expires",
			hide: false,
			filter: "agDateColumnFilter",
		},
		{
			headerName: "Created By",
			field: "createdBy",
			hide: false
		},
		{
			headerName: "Invalidation Type",
			field: "invalidationType",
			hide: false
		}
	];

	/** All of the jobs - startTime fields converted to actual Dates and derived expires field from TTL */
	$scope.jobs = jobs.map(
		function(x) {
			// need to convert this to a date object for ag-grid filter to work properly
			x.startTime = new Date(x.startTime.replace("+00", "Z"));

			x.expires = new Date(x.startTime.getTime() + x.ttlHours*3600*1000);
			return x;
		});

	/** @type {import("../agGrid/CommonGridController").CGC.DropDownOption[]} */
	$scope.dropDownOptions = [{
		name: "createJobMenuItem",
		onClick: function (){
			$scope.createJob();
		},
		text: "Create Invalidation Request",
		type: 1
	}];

	/** @type {import("../agGrid/CommonGridController").CGC.ContextMenuOption[]} */
	$scope.contextMenuOptions = [{
		onClick: function (job, $event) {
			$scope.confirmRemoveJob(job,  $event);
		},
		text: "Delete Invalidation Request",
		type: 1
	}];

	/** @type {import("../agGrid/CommonGridController").CGC.GridSettings} */
	$scope.gridOptions = {
		rowClassRules: {
			'active-job': function(params) {
				return params.data.expires > new Date();
			},
			'expired-job': function(params) {
				return params.data.expires <= new Date();
			}
		},
	};

	$scope.createJob = function() {
		locationUtils.navigateToPath('/jobs/new');
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.confirmRemoveJob = function(job) {
		const params = {
			title: 'Remove Invalidation Request?',
			message: 'Are you sure you want to remove the ' + job.assetUrl + ' invalidation request?<br><br>' +
				'NOTE: The invalidation request may have already been applied.'
		};
		const modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
			controller: 'DialogConfirmController',
			size: 'md',
			resolve: {
				params: function () {
					return params;
				}
			}
		});
		modalInstance.result.then(function() {
			jobService.deleteJob(job.id)
				.then(
					function(result) {
						messageModel.setMessages(result.data.alerts, false);
						$scope.refresh(); // refresh the jobs table
					}
				);
		});
	};

};

TableJobsController.$inject = ['jobs', '$scope', '$state', '$uibModal', 'locationUtils', 'jobService', 'messageModel'];
module.exports = TableJobsController;
