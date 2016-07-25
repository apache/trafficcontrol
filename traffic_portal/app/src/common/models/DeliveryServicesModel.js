var DeliveryServicesModel = function($location) {

    var model = this;

    var deliveryServices = [];
    var loaded = false;

    this.deliveryServices = deliveryServices;
    this.loaded = loaded;

    this.getDeliveryService = function(dsId) {
        return _.find(model.deliveryServices, function(ds){ return ds.id === dsId });
    };

    this.setDeliveryServices = function(deliveryServicesData) {
        this.deliveryServices = deliveryServicesData;
        this.loaded = true;
    };

    this.resetDeliveryServices = function() {
        this.deliveryServices = [];
        this.loaded = false;
    };

};

DeliveryServicesModel.$inject = ['$location'];
module.exports = DeliveryServicesModel;