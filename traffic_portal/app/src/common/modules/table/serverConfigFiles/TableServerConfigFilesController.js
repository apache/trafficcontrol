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

var TableServerConfigFilesController = function(server, serverConfigFiles, $scope, $state, $uibModal, locationUtils) {

	$scope.server = server;

	$scope.configFiles = serverConfigFiles.configFiles;

	$scope.refresh = function() {
		$state.reload(); // reloads all the resolves for the view
	};

	$scope.view = function(name, url) {
		var params = {
			title: name
		};
		var modalInstance = $uibModal.open({
			templateUrl: 'common/modules/dialog/text/dialog.text.tpl.html',
			controller: 'DialogTextController',
			windowClass: 'dialog-90',
			resolve: {
				params: function () {
					return params;
				},
				text: function(serverService) {
					return serverService.getServerConfigFile(url);
				}
			}
		});
		modalInstance.result.then(function() {}, function() {}); // do nothing on modal close

	};

	$scope.download = function(name, url, $event) {
		$event.stopPropagation(); // this kills the click event so it doesn't trigger anything else

		// we're going to trick the browser into opening a download dialog
		// generate a temp <a> tag
		var link = document.createElement("a");
		link.href = url;

		// keep it hidden
		link.style = "visibility:hidden";
		link.download = name;

		// briefly append the <a> tag and remove it after auto click
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	};

	$scope.navigateToPath = locationUtils.navigateToPath;

	angular.element(document).ready(function () {
		$('#configFilesTable').dataTable({
			"aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
			"iDisplayLength": 25,
			"aaSorting": [],
			"columnDefs": [
				{ 'orderable': false, 'targets': 3 },
				{ "width": "5%", "targets": 3 }
			]
		});
	});

};

TableServerConfigFilesController.$inject = ['server', 'serverConfigFiles', '$scope', '$state', '$uibModal', 'locationUtils'];
module.exports = TableServerConfigFilesController;
