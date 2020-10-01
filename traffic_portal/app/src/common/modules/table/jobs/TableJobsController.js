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

var TableJobsController = function(jobs, $scope, $state, $uibModal, locationUtils, jobService, messageModel, dateUtils) {

	$scope.jobs = jobs;

	$scope.getHourOffsetDate = dateUtils.getHourOffsetDate;

	$scope.createJob = function() {
		locationUtils.navigateToPath('/jobs/new');
	};

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.confirmRemoveJob = function(job, $event) {
		if ($event) {
			$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else
		}
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

	angular.element(document).ready(function () {
		$('#jobsTable').dataTable({
			"lengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"columnDefs": [
				{ "width": "5%", "targets": 6 },
				{ 'orderable': false, 'targets': 6 }
			],
			"aaSorting": []
		});
	});

};

TableJobsController.$inject = ['jobs', '$scope', '$state', '$uibModal', 'locationUtils', 'jobService', 'messageModel', 'dateUtils'];
module.exports = TableJobsController;
