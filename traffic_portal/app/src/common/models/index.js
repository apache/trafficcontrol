/**
 * Define the application models
 */
module.exports = angular.module('trafficPortal.models', [])
    .service('authModel', require('./AuthModel'))
    .service('chartModel', require('./ChartModel'))
    .service('deliveryServicesModel', require('./DeliveryServicesModel'))
    .service('messageModel', require('./MessageModel'))
    .service('propertiesModel', require('./PropertiesModel'))
    .service('userModel', require('./UserModel'));
