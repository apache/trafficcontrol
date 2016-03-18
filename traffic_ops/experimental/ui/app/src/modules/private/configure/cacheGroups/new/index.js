module.exports = angular.module('trafficOps.private.configure.cacheGroups.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.cacheGroups.new', {
                url: '/new',
                views: {
                    cacheGroupsContent: {
                        templateUrl: 'common/modules/form/cacheGroup/form.cacheGroup.tpl.html',
                        controller: 'FormNewCacheGroupController',
                        resolve: {
                            cacheGroup: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
