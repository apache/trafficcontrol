module.exports = angular.module('trafficOps.private.admin.tenants.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.tenants.list', {
                url: '',
                views: {
                    tenantsContent: {
                        templateUrl: 'common/modules/table/tenants/table.tenants.tpl.html',
                        controller: 'TableTenantsController',
                        resolve: {
                            tenants: function(tenantService) {
                                return tenantService.getTenants();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
