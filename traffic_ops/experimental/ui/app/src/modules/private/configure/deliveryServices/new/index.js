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

module.exports = angular.module('trafficOps.private.configure.deliveryServices.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.deliveryServices.new', {
                url: '/new',
                views: {
                    deliveryServicesContent: {
                        templateUrl: 'common/modules/form/deliveryService/form.deliveryService.tpl.html',
                        controller: 'FormNewDeliveryServiceController',
                        resolve: {
                            deliveryService: function() {
                                return {
                                    active: false,
                                    signed: false,
                                    qstringIgnore: "0",
                                    dscp: "0",
                                    geoLimit: "0",
                                    geoProvider: "0",
                                    initialDispersion: "1",
                                    ipv6RoutingEnabled: false,
                                    rangeRequestHandling: "0",
                                    multiSiteOrigin: false,
                                    regionalGeoBlocking: false,
                                    logsEnabled: false
                                };
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
