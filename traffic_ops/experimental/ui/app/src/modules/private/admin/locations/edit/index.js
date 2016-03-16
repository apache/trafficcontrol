module.exports = angular.module('trafficOps.private.admin.locations.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.locations.edit', {
                url: '/{locationId}',
                views: {
                    locationsContent: {
                        templateUrl: 'common/modules/form/location/form.location.tpl.html',
                        controller: 'FormLocationController',
                        resolve: {
                            location: function($stateParams, locationService) {
                                return locationService.getLocation($stateParams.locationId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
