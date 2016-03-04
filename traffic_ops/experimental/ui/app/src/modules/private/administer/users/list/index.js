module.exports = angular.module('trafficOps.private.administer.users.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.administer.users.list', {
                url: '',
                views: {
                    usersContent: {
                        templateUrl: 'common/modules/table/users/table.users.tpl.html',
                        controller: 'TableUsersController',
                        resolve: {
                            users: function(userService, ENV) {
                                return userService.getUsers(ENV.api['base_url'] + 'tm_user');
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
