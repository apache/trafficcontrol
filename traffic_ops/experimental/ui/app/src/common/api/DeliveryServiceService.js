var DeliveryServiceService = function(Restangular, messageModel) {

    this.getDeliveryServices = function() {
        return Restangular.all('deliveryservice').getList();
    };

    this.getDeliveryService = function(id) {
        return Restangular.one("deliveryservice", id).get();
    };

    this.updateDeliveryService = function(deliveryService) {
        return deliveryService.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Delivery service updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Delivery service update failed' } ], false);
                }
            );
    };

};

DeliveryServiceService.$inject = ['Restangular', 'messageModel'];
module.exports = DeliveryServiceService;