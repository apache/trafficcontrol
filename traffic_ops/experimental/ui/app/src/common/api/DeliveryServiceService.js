var DeliveryServiceService = function(Restangular, locationUtils, messageModel) {

    this.getDeliveryServices = function() {
        return Restangular.all('deliveryservices').getList();
    };

    this.getDeliveryService = function(id) {
        return Restangular.one("deliveryservices", id).get();
    };

    this.createDeliveryService = function(deliveryService) {
        return Restangular.service('deliveryservices').post(deliveryService)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'DeliveryService created' } ], true);
                    locationUtils.navigateToPath('/configure/delivery-services');
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.updateDeliveryService = function(deliveryService) {
        return deliveryService.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteDeliveryService = function(id) {
        return Restangular.one("deliveryservices", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service deleted' } ], true);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Delivery service delete failed' } ], false);
                }
            );
    };

};

DeliveryServiceService.$inject = ['Restangular', 'locationUtils', 'messageModel'];
module.exports = DeliveryServiceService;