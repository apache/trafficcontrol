module.exports = angular.module('trafficPortal.public.home', [])
    .controller('HomeController', require('./HomeController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.public.home', {
                url: '',
                abstract: true,
                views: {
                    publicContent: {
                        templateUrl: 'modules/public/home/home.tpl.html',
                        controller: 'HomeController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });