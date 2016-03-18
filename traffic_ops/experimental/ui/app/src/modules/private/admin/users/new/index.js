module.exports = angular.module('trafficOps.private.admin.users.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.users.new', {
                url: '/new',
                views: {
                    usersContent: {
                        templateUrl: 'common/modules/form/user/form.user.tpl.html',
                        controller: 'FormNewUserController',
                        resolve: {
                            user: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
