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

var DialogImportController = function(params, $scope, $uibModalInstance, messageModel) {

	$scope.params = params;

	$scope.import = {
		value: null
	};

	$scope.import = function() {
		var importJSON;

		try {
			importJSON = JSON.parse($scope.import.value); // json need to be valid (ie. not malformed)
		}
		catch(e) {
			messageModel.setMessages([ { level: 'error', text: 'Invalid JSON. Please validate your import JSON.' } ], false);
			$scope.cancel();
			return;
		}

		$uibModalInstance.close(importJSON);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	var handleFileDragOver = function(event) {
		event.stopPropagation();
		event.preventDefault();
		event.dataTransfer.dropEffect = 'copy';
	};

	var handleFileDrop = function(event) {
		event.stopPropagation();
		event.preventDefault();

		var file = event.dataTransfer.files[0];

		if (!file) {
			return;
		}

		// show the meta details of the file
		var output = [];
		output.push('<li><strong>', encodeURIComponent(file.name), '</strong> - ',
			file.size, ' bytes, last modified: ',
			file.lastModifiedDate ? file.lastModifiedDate.toLocaleDateString() : 'n/a',
			'</li>');
		$('#fileDetails').html('<ul>' + output.join('') + '</ul>');

		// show the contents of the file
		var reader = new FileReader();
		reader.onload = function(e) {
			$scope.import.value = e.target.result;
			$('#fileContent').val($scope.import.value);
		};
		reader.readAsText(file);
	};

	angular.element(document).ready(function () {
		var dropZone = document.getElementById('importDropZone');
		dropZone.addEventListener('dragover', handleFileDragOver, false);
		dropZone.addEventListener('drop', handleFileDrop, false);
	});

};

DialogImportController.$inject = ['params', '$scope', '$uibModalInstance', 'messageModel'];
module.exports = DialogImportController;
