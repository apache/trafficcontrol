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

module.exports = angular.module('trafficOps.api', [])
    .service('authService', require('./AuthService'))
    .service('asnService', require('./ASNService'))
    .service('cacheGroupService', require('./CacheGroupService'))
    .service('cacheGroupParameterService', require('./CacheGroupParameterService'))
    .service('cdnService', require('./CDNService'))
    .service('changeLogService', require('./ChangeLogService'))
    .service('deliveryServiceService', require('./DeliveryServiceService'))
	.service('deliveryServiceRegexService', require('./DeliveryServiceRegexService'))
	.service('divisionService', require('./DivisionService'))
    .service('httpService', require('./HttpService'))
    .service('jobService', require('./JobService'))
    .service('physLocationService', require('./PhysLocationService'))
    .service('parameterService', require('./ParameterService'))
    .service('profileService', require('./ProfileService'))
    .service('profileParameterService', require('./ProfileParameterService'))
    .service('roleService', require('./RoleService'))
    .service('regionService', require('./RegionService'))
    .service('serverService', require('./ServerService'))
    .service('staticDnsEntryService', require('./StaticDnsEntryService'))
    .service('statusService', require('./StatusService'))
    .service('tenantService', require('./TenantService'))
    .service('typeService', require('./TypeService'))
    .service('trafficOpsService', require('./TrafficOpsService'))
    .service('userService', require('./UserService'))
;
