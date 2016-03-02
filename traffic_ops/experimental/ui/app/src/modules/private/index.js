module.exports = angular.module('trafficOps.private', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private', {
                url: '',
                abstract: true,
                views: {
                    navigation: {
                        templateUrl: 'common/modules/navigation/navigation.tpl.html',
                        controller: 'NavigationController'
                    },
                    header: {
                        templateUrl: 'common/modules/header/header.tpl.html',
                        controller: 'HeaderController'
                    },
                    message: {
                        templateUrl: 'common/modules/message/message.tpl.html',
                        controller: 'MessageController'
                    },
                    content: {
                        templateUrl: 'modules/private/private.tpl.html'
                    }
                },
                resolve: {
                    user: function($state, userService, userModel, ENV) {
                        if (userModel.user.loaded) {
                            return userModel.user;
                        } else {
                            return userService.getCurrentUser(ENV.apiEndpoint['base_url'] + 'tm_user/' + userModel.userId);
                        }
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
