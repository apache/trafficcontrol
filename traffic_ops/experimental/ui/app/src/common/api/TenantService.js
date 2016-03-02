var TenantService = function(httpService) {

    this.getTenants = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.getTenant = function(endpoint) {
        return httpService.get(endpoint);
    };

};

TenantService.$inject = ['httpService'];
module.exports = TenantService;