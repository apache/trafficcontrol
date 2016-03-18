module.exports = angular.module('trafficOps.private.admin.locations.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.locations.new', {
                url: '/new',
                views: {
                    locationsContent: {
                        templateUrl: 'common/modules/form/location/form.location.tpl.html',
                        controller: 'FormNewLocationController',
                        resolve: {
                            location: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
