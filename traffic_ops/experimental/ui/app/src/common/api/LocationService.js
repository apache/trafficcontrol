var LocationService = function(httpService) {

    this.getLocations = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.getLocation = function(endpoint) {
        return httpService.get(endpoint);
    };

};

LocationService.$inject = ['httpService'];
module.exports = LocationService;