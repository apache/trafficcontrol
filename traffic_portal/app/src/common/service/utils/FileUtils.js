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

/**
 * FileUtils provides methods that allow transforming data into file downloads
 * for the user.
 */
class FileUtils {

	/**
	 * Downloads an arbitrary data blob in JSON format (indented with tab
	 * characters).
	 *
	 * @param {object} json Any object that can be stringified - which means
	 * **no** circular references and **no** functions.
	 * @param {string} fileName The name the browser will suggest for the file
	 * download.
	 * @param {string} [fileExtension] The file name "extension" to use (with
	 * the `.` omitted!). If not provided, the default is "json".
	 */
	exportJSON(json, fileName, fileExtension) {
		const jsonStr = `data:text/json;charset=utf-8,${encodeURIComponent(JSON.stringify(json, null, "\t"))}`;
		const extension = fileExtension ?? "json";
		const a = document.createElement("a");
		a.setAttribute("href", jsonStr);
		a.setAttribute("download", `${fileName}.${extension}`);
		a.click();
		a.remove();
	};

	/**
	 * Downloads a set of arbitrary data rows as a CSV. The first row is
	 * expected to be a header. If all rows are not the same length, the output
	 * will likely be considered invalid by most applications - but this is
	 * **not** checked by this utility.
	 *
	 * @deprecated Tables should all be ag-grid-based now, which has this
	 * functionality built-in so that we don't need to maintain it ourselves.
	 *
	 * @param {readonly [readonly string[], ...(readonly unknown[] | Record<PropertyKey, unknown>)[]]} arrData The data to download.
	 * @param {string} reportTitle The name that the browser will suggest for
	 * the download. ".csv" will be appended, and therefore need not be
	 * included.
	 * @param {Set<string>|string[]} [includedKeys] The keys of objects to
	 * include. If not given, all keys are included. Otherwise, for each
	 * non-header row, the keys of each object that are found in this collection
	 * will be added as row elements in alphabetical order of key name. The
	 * header row will automatically be trimmed to only include headers for
	 * included keys.
	 */
	convertToCSV(arrData, reportTitle, includedKeys) {
		let CSV = reportTitle + '\r\n\r\n';

		const keysToInclude = Array.isArray(includedKeys) ? new Set(includedKeys) : (includedKeys ?? false);
		const keys = arrData[0].filter(key => !keysToInclude || keysToInclude.has(key))
		keys.sort(); // alphabetically
		CSV += keys.join(",") + '\r\n';

		for (const rowData of arrData) {
			const row = [];
			for (const k of keys) {
				row.push(`"${rowData[k]}"`);
			}
			CSV += row.join(",") + '\r\n';
		}

		if (!CSV) {
			alert("Invalid data");
			return;
		}

		const fileName = reportTitle.replace(/ /g, "_");

		const uri = `data:text/csv;charset=utf-8,${encodeURIComponent(CSV)}`;
		const link = document.createElement("a");
		link.href = uri;
		link.style.visibility = "hidden";
		link.download = `${fileName}.csv`;

		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
}

FileUtils.$inject = [];
module.exports = FileUtils;
