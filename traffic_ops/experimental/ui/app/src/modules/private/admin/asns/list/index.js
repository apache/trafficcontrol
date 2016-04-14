module.exports = angular.module('trafficOps.private.admin.asns.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.asns.list', {
                url: '',
                views: {
                    asnsContent: {
                        templateUrl: 'common/modules/table/asns/table.asns.tpl.html',
                        controller: 'TableASNsController',
                        resolve: {
                            asns: function(asnService) {
                                return asnService.getASNs();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
