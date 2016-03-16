module.exports = angular.module('trafficOps.private.admin.users.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.users.list', {
                url: '',
                views: {
                    usersContent: {
                        templateUrl: 'common/modules/table/users/table.users.tpl.html',
                        controller: 'TableUsersController',
                        resolve: {
                            users: function(userService) {
                                return userService.getUsers();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
