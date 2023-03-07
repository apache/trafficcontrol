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
 * FormUtils contains helper methods that make it easier to interact with
 * angular forms.
 */
class FormUtils {

	/**
	 * Checks if the given controller has any errors.
	 *
	 * @param {import("angular").IFormController} input
	 * @returns {boolean}
	 */
	hasError(input) {
		return input && !input.$focused && input.$invalid;
	}

	/**
	 * Checks if the given controller has a specific error.
	 *
	 * @param {import("angular").IFormController} input
	 * @param {string} property The name of the error for which to check.
	 * @returns {boolean}
	 */
	hasPropertyError(input, property) {
		return input && !input.$focused && !!input.$error[property];
	}
}

FormUtils.$inject = [];
module.exports = FormUtils;
