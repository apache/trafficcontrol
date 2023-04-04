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
 * @param {*} cdn
 * @param {*} dnssecKeys
 * @param {*} $scope
 * @param {import("angular").ILocationService} $location
 * @param {import("../../../service/utils/DateUtils")} dateUtils
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/LocationUtils")} locationUtils
 */
var FormCdnDnssecKeysController = function(cdn, dnssecKeys, $scope, $location, dateUtils, formUtils, locationUtils) {

	$scope.generate = function() {
		$location.path($location.path() + '/generate');
	};
	$scope.regenerateKSK = function() {
		$location.path($location.path() + '/regenerateKsk');
	};

	$scope.cdn = cdn;

	$scope.ksk_new = null;

	$scope.falseTrue = [
		{ value: true, label: 'true' },
		{ value: false, label: 'false' }
	];

	$scope.generateLabel = function() {
		var label = 'Generate DNSSEC Keys';
		if ($scope.ksk_new) {
			label = 'Regenerate DNSSEC Keys';
		}
		return label;
	};

	$scope.dateFormat = dateUtils.dateFormat;

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.hasError = (input) => formUtils.hasError(input);

	$scope.hasPropertyError = (input, property) => formUtils.hasPropertyError(input, property);

	var init = function() {
		if (dnssecKeys && dnssecKeys[cdn.name]) {
			$scope.ksk_new = dnssecKeys[cdn.name].ksk.find(ksk => ksk.status === 'new' );
		}
	};
	init();

};

FormCdnDnssecKeysController.$inject = ['cdn', 'dnssecKeys', '$scope', '$location', 'dateUtils', 'formUtils', 'locationUtils'];
module.exports = FormCdnDnssecKeysController;
