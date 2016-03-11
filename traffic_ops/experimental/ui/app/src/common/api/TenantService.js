var TenantService = function(Restangular, messageModel) {

    this.getTenants = function() {
        return Restangular.all('tenant').getList();
    };

    this.getTenant = function(id) {
        return Restangular.one("tenant", id).get();
    };

    this.updateTenant = function(tenant) {
        return tenant.put()
            .then(
                function() {
                    messageModel.setMessages([ { level: 'success', text: 'Tenant updated' } ], false);
                },
                function() {
                    messageModel.setMessages([ { level: 'error', text: 'Tenant update failed' } ], false);
                }
            );
    };

};

TenantService.$inject = ['Restangular', 'messageModel'];
module.exports = TenantService;