module.exports = angular.module('trafficOps.private.admin.types', [])
    .controller('TypesController', require('./TypesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.types', {
                url: '/types',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/types/types.tpl.html',
                        controller: 'TypesController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
