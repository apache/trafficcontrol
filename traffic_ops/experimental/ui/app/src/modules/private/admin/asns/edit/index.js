module.exports = angular.module('trafficOps.private.admin.asns.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.asns.edit', {
                url: '/{asnId}/edit',
                views: {
                    asnsContent: {
                        templateUrl: 'common/modules/form/asn/form.asn.tpl.html',
                        controller: 'FormEditASNController',
                        resolve: {
                            asn: function($stateParams, asnService) {
                                return asnService.getASN($stateParams.asnId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
