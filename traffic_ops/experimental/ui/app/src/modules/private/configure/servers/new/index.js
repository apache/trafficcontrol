module.exports = angular.module('trafficOps.private.configure.servers.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.configure.servers.new', {
                url: '/new',
                views: {
                    serversContent: {
                        templateUrl: 'common/modules/form/server/form.server.tpl.html',
                        controller: 'FormNewServerController',
                        resolve: {
                            server: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
