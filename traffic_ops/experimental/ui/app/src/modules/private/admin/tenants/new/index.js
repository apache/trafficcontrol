module.exports = angular.module('trafficOps.private.admin.tenants.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.tenants.new', {
                url: '/new',
                views: {
                    tenantsContent: {
                        templateUrl: 'common/modules/form/tenant/form.tenant.tpl.html',
                        controller: 'FormNewTenantController',
                        resolve: {
                            tenant: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
