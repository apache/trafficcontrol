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

var PartitionFilter = function() {
    var cache = {};
    return function(arr, size) {
        if (!arr) { return; }
        var newArr = [];
        for (var i=0; i<arr.length; i+=size) {
            newArr.push(arr.slice(i, i+size));
        }
        var arrString = JSON.stringify(arr);
        var fromCache = cache[arrString+size];
        if (JSON.stringify(fromCache) === JSON.stringify(newArr)) {
            return fromCache;
        }
        cache[arrString+size] = newArr;
        return newArr;
    };
};

PartitionFilter.$inject = [];
module.exports = PartitionFilter;
