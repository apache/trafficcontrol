module.exports = angular.module('trafficPortal.public.about', [])
    .controller('AboutController', require('./AboutController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.public.about', {
                url: 'about',
                views: {
                    publicContent: {
                        templateUrl: 'modules/public/about/about.tpl.html',
                        controller: 'AboutController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });