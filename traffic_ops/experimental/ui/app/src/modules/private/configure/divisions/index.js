module.exports = angular.module('trafficOps.private.configure.divisions', [])
    .controller('DivisionsController', require('./DivisionsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.divisions', {
                url: '/divisions',
                abstract: true,
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/divisions/divisions.tpl.html',
                        controller: 'DivisionsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
