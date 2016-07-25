/**
 * Define the application filters
 */
module.exports = angular.module('trafficPortal.filters', [])
    //filters
    .filter('dateRangeFilter', require('./DateRangeFilter'))
    .filter('offsetFilter', require('./OffsetFilter'))
    .filter('partitionFilter', require('./PartitionFilter'))
    .filter('percentFilter', require('./PercentFilter'))
    .filter('unitsFilter', require('./UnitsFilter'))
;
