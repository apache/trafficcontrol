module.exports = angular.module('trafficOps.private.admin.cdns.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.cdns.edit', {
                url: '/{cdnId}/edit',
                views: {
                    cdnsContent: {
                        templateUrl: 'common/modules/form/cdn/form.cdn.tpl.html',
                        controller: 'FormEditCDNController',
                        resolve: {
                            cdn: function($stateParams, cdnService) {
                                return cdnService.getCDN($stateParams.cdnId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
