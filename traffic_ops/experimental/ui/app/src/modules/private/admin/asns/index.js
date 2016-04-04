module.exports = angular.module('trafficOps.private.admin.asns', [])
    .controller('AsnsController', require('./AsnsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.asns', {
                url: '/asns',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/asns/asns.tpl.html',
                        controller: 'AsnsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
