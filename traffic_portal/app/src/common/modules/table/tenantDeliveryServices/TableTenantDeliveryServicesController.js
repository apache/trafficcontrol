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

var TableTenantDeliveryServicesController = function(tenant, deliveryServices, steeringTargets, filter, $controller, $scope) {

	// extends the TableDeliveryServicesController to inherit common methods
	angular.extend(this, $controller('TableDeliveryServicesController', { tableName: 'tenantDS', deliveryServices: deliveryServices, steeringTargets: steeringTargets, filter: filter, $scope: $scope }));

	$scope.tenant = tenant;
	$scope.breadCrumbs = [
		{
			href: "#!/tenants",
			text: "Tenants"
		},
		{
			getText: () => tenant.name,
			getHref: () => `#!/tenants/${tenant.id}`
		},
		{
			text: "Delivery Services"
		}
	];
};

TableTenantDeliveryServicesController.$inject = ['tenant', 'deliveryServices', 'steeringTargets', 'filter', '$controller', '$scope'];
module.exports = TableTenantDeliveryServicesController;
