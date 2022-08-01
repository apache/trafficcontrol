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

module.exports = angular.module('trafficPortal.private.deliveryServices.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private.deliveryServices.new', {
                url: '/new?dsType',
                views: {
                    deliveryServicesContent: {
                        templateUrl: function ($stateParams) {
                            var type = $stateParams.dsType,
                                template;

                            if (type.indexOf('ANY_MAP') != -1) {
                                template = 'common/modules/form/deliveryService/form.deliveryService.anyMap.tpl.html'
                            } else if (type.indexOf('DNS') != -1) {
                                template = 'common/modules/form/deliveryService/form.deliveryService.DNS.tpl.html'
                            } else if (type.indexOf('HTTP') != -1) {
                                template = 'common/modules/form/deliveryService/form.deliveryService.HTTP.tpl.html'
                            } else if (type.indexOf('STEERING') != -1) {
                                template = 'common/modules/form/deliveryService/form.deliveryService.Steering.tpl.html'
                            }

                            return template;
                        },
                        controller: 'FormNewDeliveryServiceController',
                        resolve: {
                            deliveryService: function($stateParams, propertiesModel) {
                                var type = $stateParams.dsType,
                                    anyMapDefaults = angular.copy(propertiesModel.properties.deliveryServices.defaults.ANY_MAP),
                                    dnsDefaults = angular.copy(propertiesModel.properties.deliveryServices.defaults.DNS),
                                    httpDefaults = angular.copy(propertiesModel.properties.deliveryServices.defaults.HTTP),
                                    steeringDefaults = angular.copy(propertiesModel.properties.deliveryServices.defaults.STEERING);

                                if (type.indexOf('ANY_MAP') != -1) {
                                    return anyMapDefaults;
                                } else if (type.indexOf('DNS') != -1) {
                                    return dnsDefaults;
                                } else if (type.indexOf('HTTP') != -1) {
                                    return httpDefaults;
                                } else if (type.indexOf('STEERING') != -1) {
                                    return steeringDefaults;
                                } else {
                                    return {};
                                }
                            },
                            origin: function() {
                                return [{}];
                            },
                            topologies: function(topologyService) {
                                return topologyService.getTopologies();
                            },
                            type: function($stateParams) {
                                return $stateParams.dsType;
                            },
                            types: function(typeService) {
                                return typeService.getTypes({ useInTable: 'deliveryservice' });
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
