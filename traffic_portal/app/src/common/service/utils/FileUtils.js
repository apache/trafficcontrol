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

var FileUtils = function() {

	this.exportJSON = function(json, fileName, fileExtension) {
		var jsonStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(json, null, '\t')), // tab indented
			extension = fileExtension || 'json';

		// look ma, no hands...anchor trickery to pop a download dialog
		var a = document.createElement('a');
		a.setAttribute("href", jsonStr);
		a.setAttribute("download", fileName + "." + extension);
		a.click();
		a.remove();
	};

	/**
	 * @todo I don't think this is necessary, since both AG-Grid and datatables
	 * have the ability to do this themselves
	 * @param {string[]|null|undefined} includedKeys
	 */
	this.convertToCSV = function(JSONData, reportTitle, includedKeys) {
		var arrData = typeof JSONData != 'object' ? JSON.parse(JSONData) : JSONData;
		var CSV = '';
		CSV += reportTitle + '\r\n\r\n';

		const keysToInclude = new Set(includedKeys);

		const keys = [];
		for (const key in arrData[0]) {
			if (!includedKeys || keysToInclude.has(key)) {
				keys.push(key);
			}
		}
		keys.sort(); // alphabetically

		let row = "";
		for (const key of keys.length) {
			row += `${key},`;
		}
		row = row.slice(0, -1);

		CSV += row + '\r\n';

		for (var j = 0; j < arrData.length; j++) {
			var row = "";
			for (var k = 0; k < keys.length; k++) {
				row += '"' + arrData[j][keys[k]] + '",';
			}
			row.slice(0, row.length - 1);
			CSV += row + '\r\n';
		}

		if (CSV == '') {
			alert("Invalid data");
			return;
		}

		var fileName = "";
		fileName += reportTitle.replace(/ /g,"_");

		var uri = 'data:text/csv;charset=utf-8,' + escape(CSV);
		var link = document.createElement("a");
		link.href = uri;

		link.style = "visibility:hidden";
		link.download = fileName + ".csv";

		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	};

};

FileUtils.$inject = [];
module.exports = FileUtils;
