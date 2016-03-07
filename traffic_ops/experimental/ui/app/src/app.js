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
        'app.templates',
        'angular-jwt',
        'angular-loading-bar',

        // public modules
        require('./modules/public').name,
        require('./modules/public/login').name,

        // private modules
        require('./modules/private').name,

        // users
        require('./modules/private/user').name,
        require('./modules/private/user/edit').name,

        // administer
        require('./modules/private/administer').name,
        require('./modules/private/administer/users').name,
        require('./modules/private/administer/users/edit').name,
        require('./modules/private/administer/users/list').name,

        // configure
        require('./modules/private/configure').name,
        require('./modules/private/configure/cacheGroups').name,
        require('./modules/private/configure/cacheGroups/list').name,
        require('./modules/private/configure/deliveryServices').name,
        require('./modules/private/configure/deliveryServices/list').name,
        require('./modules/private/configure/deliveryServices/edit').name,
        require('./modules/private/configure/divisions').name,
        require('./modules/private/configure/divisions/list').name,
        require('./modules/private/configure/locations').name,
        require('./modules/private/configure/locations/list').name,
        require('./modules/private/configure/servers').name,
        require('./modules/private/configure/servers/list').name,
        require('./modules/private/configure/regions').name,
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
        require('./common/modules/form/user').name,
        require('./common/modules/header').name,
        require('./common/modules/message').name,
        require('./common/modules/navigation').name,
        require('./common/modules/release').name,

        // tables
        require('./common/modules/table/cacheGroups').name,
        require('./common/modules/table/deliveryServices').name,
        require('./common/modules/table/divisions').name,
        require('./common/modules/table/locations').name,
        require('./common/modules/table/regions').name,
        require('./common/modules/table/servers').name,
        require('./common/modules/table/tenants').name,
        require('./common/modules/table/users').name,

        // common models
        require('./common/models').name,
        require('./common/api').name,

        // common directives
        require('./common/directives/match').name,

        // common services
        require('./common/service/application').name,

        //filters
        require('./common/filters').name

    ], App)

        .config(function($stateProvider, $logProvider, $controllerProvider) {
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


