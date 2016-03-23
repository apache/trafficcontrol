module.exports = angular.module('trafficOps.private.admin.cdns.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.cdns.list', {
                url: '',
                views: {
                    cdnsContent: {
                        templateUrl: 'common/modules/table/cdns/table.cdns.tpl.html',
                        controller: 'TableCDNsController',
                        resolve: {
                            cdns: function(cdnService) {
                                return cdnService.getCDNs();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
