module.exports = angular.module('trafficPortal.utils', [])
    .service('chartUtils', require('./ChartUtils'))
    .service('formUtils', require('./FormUtils'))
    .service('jsonUtils', require('./JSONUtils'))
    .service('numberUtils', require('./NumberUtils'));