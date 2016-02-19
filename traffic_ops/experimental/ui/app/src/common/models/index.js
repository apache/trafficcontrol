module.exports = angular.module('trafficOps.models', [])
    .service('deliveryServicesModel', require('./DeliveryServicesModel'))
    .service('messageModel', require('./MessageModel'))
    .service('userModel', require('./UserModel'));
