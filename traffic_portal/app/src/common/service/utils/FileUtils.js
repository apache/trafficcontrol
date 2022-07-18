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
	 * @param {any[]|string} JSONData
	 * @param {string} reportTitle
	 * @param {string[]|null|undefined} includedKeys
	 */
	this.convertToCSV = function(JSONData, reportTitle, includedKeys) {
		const arrData = typeof JSONData != 'object' ? JSON.parse(JSONData) : JSONData;
		if (!Array.isArray(arrData)) {
			throw new TypeError("data for CSV must be an array or a JSON-encoded array")
		}
		let CSV = reportTitle + '\r\n\r\n';

		const keysToInclude = new Set(includedKeys);

		const keys = [];
		for (const key in arrData[0]) {
			if (!includedKeys || keysToInclude.has(key)) {
				keys.push(key);
			}
		}
		keys.sort(); // alphabetically

		CSV += keys.join(",") + '\r\n';

		for (const data of arrData) {
			CSV += keys.map(key=>`"${data[key]}"`).join(",");
			CSV += '\r\n';
		}

		if (!CSV) {
			alert("Invalid data");
			return;
		}

		const fileName = reportTitle.replace(/ /g,"_");
		const uri = `data:text/csv;charset=utf-8,${(CSV)}`;
		const link = document.createElement("a");
		link.href = uri;

		link.style.display = "none";
		link.style.visibility = "hidden";
		link.download = `${fileName}.csv`;

		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	};

};

FileUtils.$inject = [];
module.exports = FileUtils;
