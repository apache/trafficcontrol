module.exports = angular.module('trafficOps.private.configure.cacheGroups', [])
    .controller('CacheGroupsController', require('./CacheGroupsController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.cacheGroups', {
                url: '/cache-groups',
                abstract: true,
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/cacheGroups/cacheGroups.tpl.html',
                        controller: 'CacheGroupsController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
