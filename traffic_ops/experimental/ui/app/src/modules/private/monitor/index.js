module.exports = angular.module('trafficOps.private.monitor', [])
    .controller('MonitorController', require('./MonitorController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.monitor', {
                url: 'monitor',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/monitor/monitor.tpl.html',
                        controller: 'MonitorController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
