module.exports = angular.module('trafficOps.private.configure.regions.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.regions.list', {
                url: '',
                views: {
                    regionsContent: {
                        templateUrl: 'common/modules/table/regions/table.regions.tpl.html',
                        controller: 'TableRegionsController',
                        resolve: {
                            regions: function(regionService) {
                                return regionService.getRegions();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
