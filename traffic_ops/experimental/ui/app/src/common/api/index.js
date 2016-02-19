module.exports = angular.module('trafficOps.api', [])
    .service('authService', require('./AuthService'))
    .service('deliveryServiceService', require('./DeliveryServiceService'))
    .service('trafficOpsService', require('./TrafficOpsService'))
    .service('userService', require('./UserService'))
;