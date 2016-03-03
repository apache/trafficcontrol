var RegionService = function(httpService) {

    this.getRegions = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.getRegion = function(endpoint) {
        return httpService.get(endpoint);
    };

};

RegionService.$inject = ['httpService'];
module.exports = RegionService;