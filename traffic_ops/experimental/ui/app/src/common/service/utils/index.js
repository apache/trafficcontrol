module.exports = angular.module('trafficOps.utils', [])
    .service('formUtils', require('./FormUtils'))
    .service('locationUtils', require('./LocationUtils'))
    .service('stringUtils', require('./StringUtils'));