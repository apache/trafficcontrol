module.exports = angular.module('trafficOps.private.admin.parameters.list', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.parameters.list', {
                url: '',
                views: {
                    parametersContent: {
                        templateUrl: 'common/modules/table/parameters/table.parameters.tpl.html',
                        controller: 'TableParametersController',
                        resolve: {
                            parameters: function(parameterService) {
                                return parameterService.getParameters();
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
