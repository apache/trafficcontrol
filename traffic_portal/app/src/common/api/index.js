/**
 * Define the remote services
 */
module.exports = angular.module('trafficPortal.api', [])
    .service('authService', require('./AuthService'))
    .service('deliveryServiceService', require('./DeliveryServiceService'))
    .service('healthService', require('./HealthService'))
    .service('portalService', require('./PortalService'))
    .service('statsService', require('./StatsService'))
    .service('userService', require('./UserService'))
;