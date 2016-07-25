module.exports = angular.module('trafficPortal.private', [])
    .config(function($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('trafficPortal.private', {
                url: '',
                abstract: true,
                views: {
                    header: {
                        templateUrl: 'common/modules/header/header.tpl.html',
                        controller: 'HeaderController'
                    },
                    message: {
                        templateUrl: 'common/modules/message/message.tpl.html',
                        controller: 'MessageController'
                    },
                    content: {
                        templateUrl: 'modules/private/private.tpl.html'
                    },
                    footer: {
                        templateUrl: 'common/modules/footer/footer.tpl.html',
                        controller: 'FooterController'
                    }
                },
                resolve: {
                    user: function($state, userService, deliveryServiceService, userModel, deliveryServicesModel) {
                        if (userModel.user.loaded) {
                            return userModel.user;
                        } else {
                            return userService.getCurrentUser()
                                .then(function() {
                                    if (!deliveryServicesModel.loaded) {
                                        return deliveryServiceService.getDeliveryServices();
                                    }
                                });
                        }
                    }
                }
            });
        $urlRouterProvider.otherwise('/');
    });
