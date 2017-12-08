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
                url: '/new?type',
                views: {
                    deliveryServicesContent: {
                        templateUrl: function ($stateParams) {
                            var type = $stateParams.type,
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
                            geoMissLat: function(parameterService) {
                                return parameterService.getParameters({ name: 'default_geo_miss_latitude', configFile: 'global' });
                            },
                            geoMissLong: function(parameterService) {
                                return parameterService.getParameters({ name: 'default_geo_miss_longitude', configFile: 'global' });
                            },
                            deliveryService: function(geoMissLat, geoMissLong, $stateParams) {
                                var type = $stateParams.type;

                                var anyMapDefaults = {
                                    dscp: 0, // any map ds's don't use dscp but it's required so we'll just send it to make the api/db happy
                                    regionalGeoBlocking: false,
                                    logsEnabled: false,
                                    geoProvider: 0,
                                    geoLimit: 0
                                };

                                var dnsDefaults = {
                                    routingName: 'cdn',
                                    dscp: 0,
                                    ipv6RoutingEnabled: false,
                                    rangeRequestHandling: 0,
                                    qstringIgnore: 0,
                                    multiSiteOrigin: false,
                                    logsEnabled: false,
                                    geoProvider: 0,
                                    geoLimit: 0,
                                    missLat: (geoMissLat[0]) ? parseFloat(geoMissLat[0].value) : null,
                                    missLong: (geoMissLong[0]) ? parseFloat(geoMissLong[0].value) : null,
                                    signingAlgorithm: null,
                                    regionalGeoBlocking: false // dns ds's don't use regionalGeoBlocking but it's required so we'll just send it to make the api/db happy
                                };

                                var httpDefaults = {
                                    routingName: 'cdn',
                                    dscp: 0,
                                    ipv6RoutingEnabled: false,
                                    rangeRequestHandling: 0,
                                    qstringIgnore: 0,
                                    multiSiteOrigin: false,
                                    logsEnabled: false,
                                    initialDispersion: 0,
                                    regionalGeoBlocking: false,
                                    geoProvider: 0,
                                    geoLimit: 0,
                                    missLat: (geoMissLat[0]) ? parseFloat(geoMissLat[0].value) : null,
                                    missLong: (geoMissLong[0]) ? parseFloat(geoMissLong[0].value) : null,
                                    signingAlgorithm: null
                                };

                                var steeringDefaults = {
                                    routingName: 'cdn',
                                    ipv6RoutingEnabled: false,
                                    logsEnabled: false,
                                    geoProvider: 0,
                                    geoLimit: 0,
                                    dscp: 0, // steering ds's don't use dscp but it's required so we'll just send 0 to make the api/db happy
                                    regionalGeoBlocking: false // steering ds's don't use regionalGeoBlocking but it's required so we'll just send 0 to make the api/db happy
                                };

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
                            type: function($stateParams) {
                                return $stateParams.type;
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
