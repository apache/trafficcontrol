module.exports = angular.module('trafficOps.private.configure.cacheGroups.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.cacheGroups.list', {
                url: '',
                views: {
                    cacheGroupsContent: {
                        templateUrl: 'common/modules/table/cacheGroups/table.cacheGroups.tpl.html',
                        controller: 'TableCacheGroupsController',
                        resolve: {
                            cacheGroups: function(cacheGroupService) {
                                return cacheGroupService.getCacheGroups();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
