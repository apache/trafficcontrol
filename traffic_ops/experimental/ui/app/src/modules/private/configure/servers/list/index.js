module.exports = angular.module('trafficOps.private.configure.servers.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.servers.list', {
                url: '',
                views: {
                    serversContent: {
                        templateUrl: 'common/modules/table/servers/table.servers.tpl.html',
                        controller: 'TableServersController',
                        resolve: {
                            servers: function(serverService) {
                                return serverService.getServers();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
