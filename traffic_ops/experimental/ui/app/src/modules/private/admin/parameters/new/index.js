module.exports = angular.module('trafficOps.private.admin.parameters.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.parameters.new', {
                url: '/new',
                views: {
                    parametersContent: {
                        templateUrl: 'common/modules/form/parameter/form.parameter.tpl.html',
                        controller: 'FormNewParameterController',
                        resolve: {
                            parameter: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
