module.exports = angular.module('trafficOps.private.admin.tenants.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.tenants.edit', {
                url: '/{tenantId}/edit',
                views: {
                    tenantsContent: {
                        templateUrl: 'common/modules/form/tenant/form.tenant.tpl.html',
                        controller: 'FormEditTenantController',
                        resolve: {
                            tenant: function($stateParams, tenantService) {
                                return tenantService.getTenant($stateParams.tenantId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
