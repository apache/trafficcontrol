/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

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
