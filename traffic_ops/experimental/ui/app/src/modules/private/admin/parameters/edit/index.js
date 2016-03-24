module.exports = angular.module('trafficOps.private.admin.parameters.edit', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.parameters.edit', {
                url: '/{parameterId}/edit',
                views: {
                    parametersContent: {
                        templateUrl: 'common/modules/form/parameter/form.parameter.tpl.html',
                        controller: 'FormEditParameterController',
                        resolve: {
                            parameter: function($stateParams, parameterService) {
                                return parameterService.getParameter($stateParams.parameterId);
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
