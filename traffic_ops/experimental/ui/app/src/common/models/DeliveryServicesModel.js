var DeliveryServicesModel = function($log) {

    var model = this;

    var deliveryServices = [];
    var tenants = [];
    var loaded = false;

    this.deliveryServices = deliveryServices;
    this.tenants = tenants;
    this.loaded = loaded;

    this.showTenant = true;

    this.getDeliveryService = function(dsId) {
        return _.find(model.deliveryServices, function(ds){ return ds.xmlId === dsId });
    };

    this.getDeliveryServiceByTenant = function(tenantId) {
        return _.filter(model.deliveryServices, function(ds) { return ds.tenantId === tenantId });
    };

    this.setDeliveryServices = function(deliveryServicesData) {
        this.deliveryServices = deliveryServicesData;

        // the api currently doesn't support tenants but will
        var fakeTenants = [];
        fakeTenants[0] = 'Apache';
        fakeTenants[1] = 'FooBar Inc.';
        fakeTenants[2] = 'XYZ Corporation';

        for (var i = 0; i < this.deliveryServices.length; i++) {
            var ds = this.deliveryServices[i];
            ds['tenantId'] = ds.id;
            ds['tenantName'] = fakeTenants[Math.round(Math.random()*(fakeTenants.length - 1))];
        }

        model.setTenants(this.deliveryServices);

        this.loaded = true;
    };

    this.resetDeliveryServices = function() {
        this.deliveryServices = [];
        this.tenants = [];
        this.loaded = false;
    };

    this.getTenant = function(tenantId) {
        var tenant = _.find(model.tenants, function(tenant){ return tenant.id === tenantId });
        return tenant;
    };

    this.setTenants = function(deliveryServices) {
        this.tenants = _.map(deliveryServices, function(ds) {
            return { id: ds.tenantId, name: ds.tenantName }
        });
    };

};

DeliveryServicesModel.$inject = ['$log'];
module.exports = DeliveryServicesModel;