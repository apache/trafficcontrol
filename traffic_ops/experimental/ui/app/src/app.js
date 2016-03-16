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
        require('./modules/private/admin/users').name,
        require('./modules/private/admin/users/edit').name,
        require('./modules/private/admin/users/list').name,

        // configure
        require('./modules/private/configure').name,
        require('./modules/private/configure/cacheGroups').name,
        require('./modules/private/configure/cacheGroups/edit').name,
        require('./modules/private/configure/cacheGroups/list').name,
        require('./modules/private/configure/deliveryServices').name,
        require('./modules/private/configure/deliveryServices/edit').name,
        require('./modules/private/configure/deliveryServices/list').name,
        require('./modules/private/configure/divisions').name,
        require('./modules/private/configure/divisions/edit').name,
        require('./modules/private/configure/divisions/list').name,
        require('./modules/private/configure/locations').name,
        require('./modules/private/configure/locations/edit').name,
        require('./modules/private/configure/locations/list').name,
        require('./modules/private/configure/servers').name,
        require('./modules/private/configure/servers/edit').name,
        require('./modules/private/configure/servers/list').name,
        require('./modules/private/configure/regions').name,
        require('./modules/private/configure/regions/edit').name,
        require('./modules/private/configure/regions/list').name,

        // monitor
        require('./modules/private/monitor').name,

        // dashboards
        require('./modules/private/monitor/dashboards').name,
        require('./modules/private/monitor/dashboards/one').name,
        require('./modules/private/monitor/dashboards/two').name,
        require('./modules/private/monitor/dashboards/three').name,

        // common modules
        require('./common/modules/dialog/confirm').name,
        require('./common/modules/dialog/reset').name,
        require('./common/modules/header').name,
        require('./common/modules/message').name,
        require('./common/modules/navigation').name,
        require('./common/modules/release').name,

        // forms
        require('./common/modules/form/cacheGroup').name,
        require('./common/modules/form/deliveryService').name,
        require('./common/modules/form/division').name,
        require('./common/modules/form/location').name,
        require('./common/modules/form/region').name,
        require('./common/modules/form/server').name,
        require('./common/modules/form/tenant').name,
        require('./common/modules/form/user').name,

        // tables
        require('./common/modules/table/cacheGroups').name,
        require('./common/modules/table/deliveryServices').name,
        require('./common/modules/table/divisions').name,
        require('./common/modules/table/locations').name,
        require('./common/modules/table/regions').name,
        require('./common/modules/table/servers').name,
        require('./common/modules/table/tenants').name,
        require('./common/modules/table/users').name,

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
        request: function (config) {
            config.headers = config.headers || {};
            if ($window.sessionStorage.token) {
                config.headers.Authorization = $window.sessionStorage.token;
            }
            return config;
        },
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


