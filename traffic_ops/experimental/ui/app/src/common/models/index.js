module.exports = angular.module('trafficOps.models', [])
    .service('messageModel', require('./MessageModel'))
    .service('userModel', require('./UserModel'));
