var ChartUtils = function() {

    this.formatData = function(dataPoints) {
        if (!angular.isDefined(dataPoints)) { return [] };

        var formattedData = [];

        var dataPointsArray = dataPoints.split(',');
        for (var i = 0; i < dataPointsArray.length; i++) {
            dataPointsArray[i] = +dataPointsArray[i];
        }
        for (var j = 0; j < dataPointsArray.length; j+=2) {
            formattedData.push(dataPointsArray.slice(j, j+2))
        }
        return formattedData;
    };

};

ChartUtils.$inject = [];
module.exports = ChartUtils;
