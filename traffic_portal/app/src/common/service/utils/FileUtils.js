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

};

FileUtils.$inject = [];
module.exports = FileUtils;
