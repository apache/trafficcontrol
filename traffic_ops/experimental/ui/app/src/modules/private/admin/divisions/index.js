module.exports = angular.module('trafficOps.private.admin.divisions', [])
    .controller('DivisionsController', require('./DivisionsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.divisions', {
                url: '/divisions',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/divisions/divisions.tpl.html',
                        controller: 'DivisionsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
