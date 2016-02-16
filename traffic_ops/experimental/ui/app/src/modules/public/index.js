module.exports = angular.module('trafficOps.public', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.public', {
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
                        templateUrl: 'modules/public/public.tpl.html'
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
