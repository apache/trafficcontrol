module.exports = angular.module('trafficPortal.public.home.landing', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.public.home.landing', {
                url: ''
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
