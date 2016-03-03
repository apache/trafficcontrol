var CacheGroupService = function(httpService) {

    this.getCacheGroups = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.getCacheGroup = function(endpoint) {
        return httpService.get(endpoint);
    };

};

CacheGroupService.$inject = ['httpService'];
module.exports = CacheGroupService;