module.exports = angular.module('trafficOps.private.admin.types.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.types.edit', {
                url: '/{typeId}/edit',
                views: {
                    typesContent: {
                        templateUrl: 'common/modules/form/type/form.type.tpl.html',
                        controller: 'FormEditTypeController',
                        resolve: {
                            type: function($stateParams, typeService) {
                                return typeService.getType($stateParams.typeId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
