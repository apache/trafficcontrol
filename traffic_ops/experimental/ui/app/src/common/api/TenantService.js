var TenantService = function(Restangular, messageModel) {

    this.getTenants = function() {
        return Restangular.all('tenant').getList();
    };

    this.getTenant = function(id) {
        return Restangular.one("tenant", id).get();
    };

    this.createTenant = function(tenant) {
        return Restangular.service('tenant').post(tenant)
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant created' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
        );
    };

    this.updateTenant = function(tenant) {
        return tenant.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant updated' } ], false);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

    this.deleteTenant = function(id) {
        return Restangular.one("tenant", id).remove()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant deleted' } ], true);
                },
                function(fault) {
                    messageModel.setMessages(fault.data.alerts, false);
                }
            );
    };

};

TenantService.$inject = ['Restangular', 'messageModel'];
module.exports = TenantService;