module.exports = angular.module('trafficOps.private.admin.users.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.users.edit', {
                url: '/{userId}/edit',
                views: {
                    usersContent: {
                        templateUrl: 'common/modules/form/user/form.user.tpl.html',
                        controller: 'FormEditUserController',
                        resolve: {
                            user: function($stateParams, userService) {
                                return userService.getUser($stateParams.userId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
