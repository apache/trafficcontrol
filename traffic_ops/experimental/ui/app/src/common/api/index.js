module.exports = angular.module('trafficOps.api', [])
    .service('authService', require('./AuthService'))
    .service('trafficOpsService', require('./TrafficOpsService'))
    .service('userService', require('./UserService'))
;