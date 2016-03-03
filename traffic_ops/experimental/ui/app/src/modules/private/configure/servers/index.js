module.exports = angular.module('trafficOps.private.configure.servers', [])
    .controller('ServersController', require('./ServersController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.servers', {
                url: '/servers',
                abstract: true,
                views: {
                    configureContent: {
                        templateUrl: 'modules/private/configure/servers/servers.tpl.html',
                        controller: 'ServersController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
