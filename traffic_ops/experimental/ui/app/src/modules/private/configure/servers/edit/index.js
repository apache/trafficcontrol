module.exports = angular.module('trafficOps.private.configure.servers.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.servers.edit', {
                url: '/{serverId}/edit',
                views: {
                    serversContent: {
                        templateUrl: 'common/modules/form/server/form.server.tpl.html',
                        controller: 'FormEditServerController',
                        resolve: {
                            server: function($stateParams, serverService) {
                                return serverService.getServer($stateParams.serverId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
