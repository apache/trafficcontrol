module.exports = angular.module('trafficOps.private.admin.types.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.types.new', {
                url: '/new',
                views: {
                    typesContent: {
                        templateUrl: 'common/modules/form/type/form.type.tpl.html',
                        controller: 'FormNewTypeController',
                        resolve: {
                            type: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
