module.exports = angular.module('trafficOps.public.login', [])
    .controller('LoginController', require('./LoginController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.public.login', {
                url: '',
                views: {
                    publicContent: {
                        templateUrl: 'modules/public/login/login.tpl.html',
                        controller: 'LoginController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });