module.exports = angular.module('trafficOps.private.admin.profiles', [])
    .controller('ProfilesController', require('./ProfilesController'))
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.profiles', {
                url: '/profiles',
                abstract: true,
                views: {
                    adminContent: {
                        templateUrl: 'modules/private/admin/profiles/profiles.tpl.html',
                        controller: 'ProfilesController'
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
