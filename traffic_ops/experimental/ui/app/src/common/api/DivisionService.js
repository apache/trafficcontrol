var DivisionService = function(httpService) {

    this.getDivisions = function(endpoint) {
        return httpService.get(endpoint);
    };

    this.getDivision = function(endpoint) {
        return httpService.get(endpoint);
    };

};

DivisionService.$inject = ['httpService'];
module.exports = DivisionService;