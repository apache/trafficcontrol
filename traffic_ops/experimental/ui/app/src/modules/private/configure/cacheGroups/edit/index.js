module.exports = angular.module('trafficOps.private.configure.cacheGroups.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.cacheGroups.edit', {
                url: '/{cacheGroupId}',
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
