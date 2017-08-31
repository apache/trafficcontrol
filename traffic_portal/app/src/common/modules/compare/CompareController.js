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

var CompareController = function(item1Name, item2Name, item1, item2, $scope, $window) {

	$scope.item1Name = item1Name;
	$scope.item2Name = item2Name;

	var performDiff = function(item1JSON, item2JSON, destination) {
		var div = null,
			source = '';

		var display = document.getElementById(destination),
			fragment = document.createDocumentFragment();

		if (item1JSON && item2JSON) {
			var diff = JsDiff.diffJson(item1JSON, item2JSON);
			diff.forEach(function(part){
				source = part.added ? $scope.item2Name : part.removed ? $scope.item1Name : '';
				div = document.createElement('div');
				div.className = part.added ? 'item2' : part.removed ? 'item1' : 'same';

				var sourceDiv = document.createElement('div');
				var sourceName = document.createTextNode(source);
				sourceDiv.className = 'source';
				sourceDiv.appendChild(sourceName);

				var partDiv = document.createElement('div');
				var partValue = document.createTextNode(part.value);
				partDiv.appendChild(partValue);

				div.appendChild(sourceDiv);
				div.appendChild(partDiv);
				fragment.appendChild(div);
			});

			display.innerHTML = '';
			display.appendChild(fragment);
		} else {
			display.innerHTML = 'Diff failed.';
		}
	};

	var compare = function() {
		$('#diff').html('<i class="fa fa-refresh fa-spin fa-1x fa-fw"></i> Comparing...');
		performDiff(item1, item2, 'diff');
	};

	$scope.back = function() {
		$window.history.back();
	};

	angular.element(document).ready(function () {
		compare();
	});

};

CompareController.$inject = ['item1Name', 'item2Name', 'item1', 'item2', '$scope', '$window'];
module.exports = CompareController;
