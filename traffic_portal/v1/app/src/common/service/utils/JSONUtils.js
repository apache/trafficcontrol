/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var JSONUtils = function() {

    this.convertToCSV = function(JSONData, reportTitle, includedKeys) {
        // if JSONData is not an object then JSON.parse will parse the JSON string in an Object
        var arrData = typeof JSONData != 'object' ? JSON.parse(JSONData) : JSONData;

        var CSV = '';
        // set report title in first row or line
        CSV += reportTitle + '\r\n\n';

        // this loop will extract the labels from the first hash in the array
        var keys = [];
        for (var key in arrData[0]) {
            if (!includedKeys || _.contains(includedKeys, key)) {
                keys.push(key);
            }
        }
        keys.sort(); // alphabetically

        var row = "";
        for (var i = 0; i < keys.length; i++) {
            //Now convert each value to string and comma-separate
            row += keys[i] + ',';
        }
        row = row.slice(0, -1);

        //append Label row with line break
        CSV += row + '\r\n';

        // outer loop is to extract each row
        for (var j = 0; j < arrData.length; j++) {
            var row = "";

            // inner loop to extract each column by name and convert it to a comma-separated string
            for (var k = 0; k < keys.length; k++) {
                row += '"' + arrData[j][keys[k]] + '",';
            }

            row.slice(0, row.length - 1);

            // add a line break after each row
            CSV += row + '\r\n';
        }

        if (CSV == '') {
            alert("Invalid data");
            return;
        }

        // generate a file name
        var fileName = "";
        // this will remove the blank-spaces from the title and replace it with an underscore
        fileName += reportTitle.replace(/ /g,"_");

        // initialize file format to csv
        var uri = 'data:text/csv;charset=utf-8,' + escape(CSV);

        // Now the little tricky part.
        // you can use either>> window.open(uri);
        // but this will not work in some browsers
        // or you will not get the correct file extension

        // this trick will generate a temp <a /> tag
        var link = document.createElement("a");
        link.href = uri;

        // set the visibility hidden so it will not effect on your web-layout
        link.style = "visibility:hidden";
        link.download = fileName + ".csv";

        // this part will append the anchor tag and remove it after automatic click
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

};

JSONUtils.$inject = [];
module.exports = JSONUtils;
