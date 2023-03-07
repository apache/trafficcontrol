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
 * @param {*} federation
 * @param {*} resolver
 * @param {*} types
 * @param {*} $scope
 * @param {*} $uibModalInstance
 * @param {import("../../../service/utils/FormUtils")} formUtils
 */
var DialogFederationResolverController = function(federation, resolver, types, $scope, $uibModalInstance, formUtils) {

	$scope.federation = federation;

	$scope.resolver = resolver;

	$scope.types = types;

	$scope.create = function(resolver) {
		$uibModalInstance.close(resolver);
	};

	$scope.cancel = function () {
		$uibModalInstance.dismiss('cancel');
	};

	$scope.hasError = formUtils.hasError;

	$scope.hasPropertyError = formUtils.hasPropertyError;

};

DialogFederationResolverController.$inject = ['federation', 'resolver', 'types', '$scope', '$uibModalInstance', 'formUtils'];
module.exports = DialogFederationResolverController;
