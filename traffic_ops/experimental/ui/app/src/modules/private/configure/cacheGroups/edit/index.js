module.exports = angular.module('trafficOps.private.configure.cacheGroups.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.cacheGroups.edit', {
                url: '/{cacheGroupId}/edit',
                views: {
                    cacheGroupsContent: {
                        templateUrl: 'common/modules/form/cacheGroup/form.cacheGroup.tpl.html',
                        controller: 'FormEditCacheGroupController',
                        resolve: {
                            cacheGroup: function($stateParams, cacheGroupService) {
                                return cacheGroupService.getCacheGroup($stateParams.cacheGroupId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
