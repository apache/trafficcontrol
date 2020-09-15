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

module.exports = function() {
	// selectDropdownbyNum - pass in the <SELECT> element and a option number, typically 1
	this.selectDropdownbyNum = function ( element, optionNum ) {
		if (optionNum){
			var options = element.all(by.tagName('option'))
				.then(function(options){
					options[optionNum].click();
				});
		}
	};

	this.urlPath = function ( url ) {
		return '/' + String(url).split('/').slice(3).join('/');
	};

	this.shuffle = (string) => {
		return [...string].sort(
			(a, b) => {
				return Math.floor(Math.random() * 3) - 1;
			}
		).join("");
	}

	this.random = (max) => {
		return Math.round(Math.random() * max);   
	}

};
