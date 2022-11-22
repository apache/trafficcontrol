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
 * UnitsFilter is a factory for a filter that adds US-common commas to long
 * numeric values, and optionally also units.
 *
 * @param {import("../service/utils/NumberUtils")} numberUtils
 */
function UnitsFilter(numberUtils) {
	/**
	 * @param {number} kilobits
	 * @param {boolean} [prettify]
	 */
	return function(kilobits, prettify) {
		const [val, units] = numberUtils.shrink(kilobits);
		if (prettify) {
			return `${numberUtils.addCommas(val)} ${units}`;
		}
		return String(val);
	};
};

UnitsFilter.$inject = ['numberUtils'];
module.exports = UnitsFilter;
