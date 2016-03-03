module.exports = angular.module('trafficOps.private.configure', [])
    .controller('ConfigureController', require('./ConfigureController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure', {
                url: 'configure',
                abstract: true,
                views: {
                    privateContent: {
                        templateUrl: 'modules/private/configure/configure.tpl.html',
                        controller: 'ConfigureController'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
