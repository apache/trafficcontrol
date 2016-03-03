var ServerService = function(httpService) {

    this.getServers = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.getServer = function(endpoint) {
        return httpService.get(endpoint);
    };

};

ServerService.$inject = ['httpService'];
module.exports = ServerService;