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

'use strict';
require('app-templates');

var App = function($urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
};

App.$inject = ['$urlRouterProvider'];

var trafficOps = angular.module('trafficOps', [
        'config',
        'ngAnimate',
        'ngResource',
        'ngSanitize',
        'ngRoute',
        'ngMap',
        'ui.router',
        'ui.bootstrap',
        'restangular',
        'app.templates',
        'angular-jwt',
        'angular-loading-bar',

        // public modules
        require('./modules/public').name,
        require('./modules/public/login').name,

        // private modules
        require('./modules/private').name,

        // current user
        require('./modules/private/user').name,
        require('./modules/private/user/edit').name,

        // admin
        require('./modules/private/admin').name,
        require('./modules/private/admin/asns').name,
        require('./modules/private/admin/asns/edit').name,
        require('./modules/private/admin/asns/list').name,
        require('./modules/private/admin/asns/new').name,
        require('./modules/private/admin/cdns').name,
        require('./modules/private/admin/cdns/deliveryServices').name,
        require('./modules/private/admin/cdns/edit').name,
        require('./modules/private/admin/cdns/list').name,
        require('./modules/private/admin/cdns/new').name,
        require('./modules/private/admin/cdns/servers').name,
        require('./modules/private/admin/divisions').name,
        require('./modules/private/admin/divisions/edit').name,
        require('./modules/private/admin/divisions/list').name,
        require('./modules/private/admin/divisions/new').name,
        require('./modules/private/admin/divisions/regions').name,
        require('./modules/private/admin/physLocations').name,
        require('./modules/private/admin/physLocations/edit').name,
        require('./modules/private/admin/physLocations/list').name,
        require('./modules/private/admin/physLocations/new').name,
        require('./modules/private/admin/physLocations/servers').name,
        require('./modules/private/admin/parameters').name,
        require('./modules/private/admin/parameters/edit').name,
        require('./modules/private/admin/parameters/list').name,
        require('./modules/private/admin/parameters/new').name,
        require('./modules/private/admin/parameters/profiles').name,
        require('./modules/private/admin/profiles').name,
        require('./modules/private/admin/profiles/edit').name,
        require('./modules/private/admin/profiles/list').name,
        require('./modules/private/admin/profiles/new').name,
        require('./modules/private/admin/profiles/parameters').name,
        require('./modules/private/admin/profiles/servers').name,
        require('./modules/private/admin/regions').name,
        require('./modules/private/admin/regions/edit').name,
        require('./modules/private/admin/regions/list').name,
        require('./modules/private/admin/regions/physLocations').name,
        require('./modules/private/admin/regions/new').name,
        require('./modules/private/admin/statuses').name,
        require('./modules/private/admin/statuses/edit').name,
        require('./modules/private/admin/statuses/list').name,
        require('./modules/private/admin/statuses/new').name,
        require('./modules/private/admin/statuses/servers').name,
        require('./modules/private/admin/tenants').name,
        require('./modules/private/admin/tenants/edit').name,
        require('./modules/private/admin/tenants/list').name,
        require('./modules/private/admin/tenants/new').name,
        require('./modules/private/admin/types').name,
        require('./modules/private/admin/types/edit').name,
        require('./modules/private/admin/types/list').name,
        require('./modules/private/admin/types/new').name,
        require('./modules/private/admin/types/servers').name,
        require('./modules/private/admin/types/cacheGroups').name,
        require('./modules/private/admin/types/deliveryServices').name,
        require('./modules/private/admin/users').name,
        require('./modules/private/admin/users/deliveryServices').name,
        require('./modules/private/admin/users/edit').name,
        require('./modules/private/admin/users/list').name,
        require('./modules/private/admin/users/new').name,

        // configure
        require('./modules/private/configure').name,
        require('./modules/private/configure/cacheGroups').name,
        require('./modules/private/configure/cacheGroups/edit').name,
        require('./modules/private/configure/cacheGroups/list').name,
        require('./modules/private/configure/cacheGroups/new').name,
        require('./modules/private/configure/cacheGroups/asns').name,
        require('./modules/private/configure/cacheGroups/parameters').name,
        require('./modules/private/configure/cacheGroups/servers').name,
        require('./modules/private/configure/deliveryServices').name,
        require('./modules/private/configure/deliveryServices/edit').name,
        require('./modules/private/configure/deliveryServices/list').name,
        require('./modules/private/configure/deliveryServices/new').name,
        require('./modules/private/configure/deliveryServices/regexes').name,
        require('./modules/private/configure/deliveryServices/servers').name,
        require('./modules/private/configure/deliveryServices/users').name,
        require('./modules/private/configure/servers').name,
        require('./modules/private/configure/servers/deliveryServices').name,
        require('./modules/private/configure/servers/edit').name,
        require('./modules/private/configure/servers/new').name,
        require('./modules/private/configure/servers/list').name,

        // monitor
        require('./modules/private/monitor').name,

        // dashboards
        require('./modules/private/monitor/dashboards').name,
        require('./modules/private/monitor/dashboards/one').name,
        require('./modules/private/monitor/dashboards/two').name,
        require('./modules/private/monitor/dashboards/three').name,

        // common modules
        require('./common/modules/dialog/confirm').name,
        require('./common/modules/dialog/delete').name,
        require('./common/modules/dialog/reset').name,
        require('./common/modules/header').name,
        require('./common/modules/message').name,
        require('./common/modules/navigation').name,
        require('./common/modules/release').name,

        // forms
        require('./common/modules/form/cacheGroup').name,
        require('./common/modules/form/cacheGroup/edit').name,
        require('./common/modules/form/cacheGroup/new').name,
        require('./common/modules/form/asn').name,
        require('./common/modules/form/asn/edit').name,
        require('./common/modules/form/asn/new').name,
        require('./common/modules/form/cdn').name,
        require('./common/modules/form/cdn/edit').name,
        require('./common/modules/form/cdn/new').name,
        require('./common/modules/form/deliveryService').name,
        require('./common/modules/form/deliveryService/edit').name,
        require('./common/modules/form/deliveryService/new').name,
        require('./common/modules/form/division').name,
        require('./common/modules/form/division/edit').name,
        require('./common/modules/form/division/new').name,
        require('./common/modules/form/physLocation').name,
        require('./common/modules/form/physLocation/edit').name,
        require('./common/modules/form/physLocation/new').name,
        require('./common/modules/form/parameter').name,
        require('./common/modules/form/parameter/edit').name,
        require('./common/modules/form/parameter/new').name,
        require('./common/modules/form/profile').name,
        require('./common/modules/form/profile/edit').name,
        require('./common/modules/form/profile/new').name,
        require('./common/modules/form/region').name,
        require('./common/modules/form/region/edit').name,
        require('./common/modules/form/region/new').name,
        require('./common/modules/form/server').name,
        require('./common/modules/form/server/edit').name,
        require('./common/modules/form/server/new').name,
        require('./common/modules/form/status').name,
        require('./common/modules/form/status/edit').name,
        require('./common/modules/form/status/new').name,
        require('./common/modules/form/tenant').name,
        require('./common/modules/form/tenant/edit').name,
        require('./common/modules/form/tenant/new').name,
        require('./common/modules/form/type').name,
        require('./common/modules/form/type/edit').name,
        require('./common/modules/form/type/new').name,
        require('./common/modules/form/user').name,
        require('./common/modules/form/user/edit').name,
        require('./common/modules/form/user/new').name,

        // tables
        require('./common/modules/table/cacheGroups').name,
        require('./common/modules/table/cacheGroupAsns').name,
        require('./common/modules/table/cacheGroupParameters').name,
        require('./common/modules/table/cacheGroupServers').name,
        require('./common/modules/table/asns').name,
        require('./common/modules/table/cdns').name,
        require('./common/modules/table/cdnDeliveryServices').name,
        require('./common/modules/table/cdnServers').name,
        require('./common/modules/table/deliveryServices').name,
        require('./common/modules/table/deliveryServiceRegexes').name,
        require('./common/modules/table/deliveryServiceServers').name,
        require('./common/modules/table/deliveryServiceUsers').name,
        require('./common/modules/table/divisions').name,
        require('./common/modules/table/divisionRegions').name,
        require('./common/modules/table/physLocations').name,
        require('./common/modules/table/physLocationServers').name,
        require('./common/modules/table/parameters').name,
        require('./common/modules/table/parameterProfiles').name,
        require('./common/modules/table/profileParameters').name,
        require('./common/modules/table/profileServers').name,
        require('./common/modules/table/profiles').name,
        require('./common/modules/table/regions').name,
        require('./common/modules/table/regionPhysLocations').name,
        require('./common/modules/table/servers').name,
        require('./common/modules/table/serverDeliveryServices').name,
        require('./common/modules/table/statuses').name,
        require('./common/modules/table/statusServers').name,
        require('./common/modules/table/types').name,
        require('./common/modules/table/typeCacheGroups').name,
        require('./common/modules/table/typeDeliveryServices').name,
        require('./common/modules/table/typeServers').name,
        require('./common/modules/table/users').name,
        require('./common/modules/table/userDeliveryServices').name,

        // models
        require('./common/models').name,
        require('./common/api').name,

        // directives
        require('./common/directives/match').name,

        // services
        require('./common/service/application').name,
        require('./common/service/utils').name,

        // filters
        require('./common/filters').name

    ], App)

        .config(function($stateProvider, $logProvider, $controllerProvider, RestangularProvider, ENV) {

            RestangularProvider.setBaseUrl(ENV.api['root']);

            RestangularProvider.setResponseInterceptor(function(data, operation, what) {
                if (angular.isDefined(data.response)) { // todo: this should not be needed. need better solution.
                    if (operation == 'getList') {
                        return data.response;
                    }
                    return data.response[0];
                } else {
                    return data;
                }
            });

            $controllerProvider.allowGlobals();
            $logProvider.debugEnabled(true);
            $stateProvider
                .state('trafficOps', {
                    url: '/',
                    abstract: true,
                    templateUrl: 'common/templates/master.tpl.html'
                });
        })

        .run(function($log, applicationService) {
            $log.debug("Application run...");
            applicationService.startup();
        })
    ;

trafficOps.factory('authInterceptor', function ($q, $window, $location, $timeout, messageModel, userModel) {
    return {
        responseError: function (rejection) {
            var url = $location.url(),
                alerts = [];

            try { alerts = rejection.data.alerts; }
            catch(e) {}

            // 401, 403, 404 and 5xx errors handled globally; all others handled in fault handler
            if (rejection.status === 401) {
                userModel.resetUser();
                if (url == '/' || $location.search().redirect) {
                    messageModel.setMessages(alerts, false);
                } else {
                    $timeout(function () {
                        messageModel.setMessages(alerts, true);
                        // forward the to the login page with ?redirect=page/they/were/trying/to/reach
                        $location.url('/').search({ redirect: encodeURIComponent(url) });
                    }, 200);
                }
            } else if (rejection.status === 403 || rejection.status === 404) {
                $timeout(function () {
                    messageModel.setMessages(alerts, false);
                }, 200);
            } else if (rejection.status.toString().match(/^5\d[01356789]$/)) {
                // matches 5xx EXCEPT for 502's and 504's which indicate a timeout and will be handled by each service call accordingly
                $timeout(function () {
                    messageModel.setMessages([ { level: 'error', text: rejection.status.toString() + ': ' + rejection.statusText } ], false);
                }, 200);
            }

            return $q.reject(rejection);
        }
    };
});

trafficOps.config(function ($httpProvider) {
    $httpProvider.interceptors.push('authInterceptor');
});


