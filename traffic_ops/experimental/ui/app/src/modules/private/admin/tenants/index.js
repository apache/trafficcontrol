module.exports = angular.module('trafficOps.private.admin.tenants', [])
    .controller('TenantsController', require('./TenantsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.tenants', {
                url: '/tenants',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/tenants/tenants.tpl.html',
                        controller: 'TenantsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
