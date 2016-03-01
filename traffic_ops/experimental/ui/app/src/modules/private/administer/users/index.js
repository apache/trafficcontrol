module.exports = angular.module('trafficOps.private.administer.users', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.administer.users', {
                url: '/users',
                views: {
                    administerContent: {
                        templateUrl: 'common/modules/table/users/table.users.tpl.html',
                        controller: 'TableUsersController',
                        resolve: {
                            users: function(userService, ENV) {
                                return userService.getUsers(ENV.apiEndpoint['get_users']);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
