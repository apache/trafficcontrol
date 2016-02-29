module.exports = angular.module('trafficOps.private.administer.user', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.administer.user', {
                url: '/users/{userId}',
                views: {
                    administerContent: {
                        templateUrl: 'common/modules/form/user/form.user.tpl.html',
                        controller: 'FormUserController',
                        resolve: {
                            foo: function($stateParams, userService) {
                                return userService.getUser($stateParams.userId, false);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
