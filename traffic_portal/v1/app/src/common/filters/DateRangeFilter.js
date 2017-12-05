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

var DateRangeFilter = function() {
    return function(items, dateRange) {
        if (_.isUndefined(dateRange)) return items;
        var filteredItems = _.filter(items, function(item) {
            return (
                ( moment(item.date).isAfter(dateRange.start, 'day') || moment(item.date).isSame(dateRange.start, 'day') ) &&
                ( moment(item.date).isBefore(dateRange.end, 'day') || moment(item.date).isSame(dateRange.end, 'day') )
            );
        });
        return filteredItems;
    };
};

DateRangeFilter.$inject = [];
module.exports = DateRangeFilter;
