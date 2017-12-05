/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

'use strict';
require('app-templates');

var App = function($urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
};

App.$inject = ['$urlRouterProvider'];

var trafficPortal = angular.module('trafficPortal', [
        'config',
        'ngAnimate',
        'ngResource',
        'ngSanitize',
        'ngRoute',
        'ui.router',
        'ui.bootstrap',
        'ui.bootstrap.datetimepicker',
        'app.templates',
        'angular-loading-bar',

        // public modules
        require('./modules/public').name,
        require('./modules/public/home').name,
        require('./modules/public/home/landing').name,
        require('./modules/public/about').name,

        // private modules
        require('./modules/private').name,

        // collateral
        require('./modules/private/collateral').name,

        // dashboard
        require('./modules/private/dashboard').name,
        require('./modules/private/dashboard/overview').name,

        // delivery service
        require('./modules/private/deliveryService').name,
        require('./modules/private/deliveryService/new').name,
        require('./modules/private/deliveryService/view').name,
        require('./modules/private/deliveryService/view/overview').name,
        require('./modules/private/deliveryService/view/overview/detail').name,

        // delivery service charts
        require('./modules/private/deliveryService/view/charts').name,
        require('./modules/private/deliveryService/view/charts/bandwidthPerSecond').name,
        require('./modules/private/deliveryService/view/charts/httpStatus').name,
        require('./modules/private/deliveryService/view/charts/transactionsPerSecond').name,

        // user
        require('./modules/private/user').name,
        require('./modules/private/user/edit').name,
        require('./modules/private/user/register').name,
        require('./modules/private/user/reset').name,

        // common modules
        require('./common/modules/cacheGroups').name,
        require('./common/modules/chart/bandwidthPerSecond').name,
        require('./common/modules/chart/capacity').name,
        require('./common/modules/chart/dates').name,
        require('./common/modules/chart/httpStatus').name,
        require('./common/modules/chart/routing').name,
        require('./common/modules/chart/transactionsPerSecond').name,
        require('./common/modules/dates').name,
        require('./common/modules/deliveryService/config/edit').name,
        require('./common/modules/dialog/confirm').name,
        require('./common/modules/dialog/reset').name,
        require('./common/modules/footer').name,
        require('./common/modules/header').name,
        require('./common/modules/message').name,
        require('./common/modules/release/version').name,
        require('./common/modules/tools/purge').name,

        require('./common/models').name,
        require('./common/api').name,

        //directives
        require('./common/directives/enter').name,
        require('./common/directives/formattedDate').name,
        require('./common/directives/match').name,
        require('./common/directives/rcSubmit').name,
        require('./common/directives/rcVerifySet').name,
        require('./common/directives/selectOnClick').name,

        // services
        require('./common/service/application').name,
        require('./common/service/utils').name,
        require('./common/service/utils/date').name,

        //filters
        require('./common/filters').name

    ], App)

        .controller('AppController', require('./AppController'))

        .config(function($stateProvider, $logProvider, $controllerProvider) {
            $controllerProvider.allowGlobals();
            $logProvider.debugEnabled(true);
            $stateProvider
                .state('trafficPortal', {
                    url: '/',
                    abstract: true,
                    templateUrl: 'common/templates/master.tpl.html',
                    controller: 'AppController',
                    resolve: {
                        properties: function(portalService, propertiesModel) {
                            return portalService.getProperties()
                                .then(function(result) {
                                    propertiesModel.setProperties(result);
                                });
                        }
                    }
                });
        })

        .run(function(applicationService) {
            applicationService.startup();
        })
    ;

trafficPortal.factory('authInterceptor', function ($q, $location, $timeout, dateUtils, messageModel, userModel) {
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
                        // forward the to the home page with ?redirect=page/they/were/trying/to/reach
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
                    messageModel.setMessages([ { level: 'error', text: rejection.status.toString() + ': ' + rejection.statusText + ' (' + dateUtils.dateFormat(new Date(), "UTC:dd/mmm/yyyy:HH:MM:ss o") + ')'  } ], false);
                }, 200);
            }

            return $q.reject(rejection);
        }
    };
});

trafficPortal.config(function ($httpProvider) {
    $httpProvider.interceptors.push('authInterceptor');
});


