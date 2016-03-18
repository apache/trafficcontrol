module.exports = angular.module('trafficOps.private.admin.divisions.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.divisions.new', {
                url: '/new',
                views: {
                    divisionsContent: {
                        templateUrl: 'common/modules/form/division/form.division.tpl.html',
                        controller: 'FormNewDivisionController',
                        resolve: {
                            division: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
