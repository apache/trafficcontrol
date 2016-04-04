module.exports = angular.module('trafficOps.private.admin.types.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.types.list', {
                url: '',
                views: {
                    typesContent: {
                        templateUrl: 'common/modules/table/types/table.types.tpl.html',
                        controller: 'TableTypesController',
                        resolve: {
                            types: function(typeService) {
                                return typeService.getTypes();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
