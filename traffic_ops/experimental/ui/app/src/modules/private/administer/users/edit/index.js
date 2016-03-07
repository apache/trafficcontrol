module.exports = angular.module('trafficOps.private.administer.users.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.administer.users.edit', {
                url: '/{userId}',
                views: {
                    usersContent: {
                        templateUrl: 'common/modules/form/user/form.user.tpl.html',
                        controller: 'FormUserController',
                        resolve: {
                            user: function($stateParams, userService, ENV) {
                                return userService.getUser(ENV.api['root'] + 'tm_user/' + $stateParams.userId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
