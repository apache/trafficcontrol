module.exports = angular.module('trafficOps.private.admin.profiles.new', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficOps.private.admin.profiles.new', {
                url: '/new',
                views: {
                    profilesContent: {
                        templateUrl: 'common/modules/form/profile/form.profile.tpl.html',
                        controller: 'FormNewProfileController',
                        resolve: {
                            profile: function() {
                                return {};
                            }
                        }
                    }
                }
            })
        ;
        $urlRouterProvider.otherwise('/');
    });
