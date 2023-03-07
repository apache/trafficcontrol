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
 * LocationUtils provides methods for performing common actions that affect
 * routing.
 */
class LocationUtils {
	/**
	 * @param {import("angular").ILocationService} $location
	 * @param {import("./angular.ui.bootstrap").IModalService} $uibModal
	 */
	constructor($location, $uibModal) {
		this.$location = $location;
		this.$uibModal = $uibModal;
	}

	/**
	 * Navigates the current browsing context to the provided relative or
	 * absolute URL.
	 *
	 * **In almost all cases, you should be using an Anchor element with
	 * `ng-href` instead!**
	 *
	 * @param {string} path The path to which to navigate (may include query
	 * string and document fragment).
	 * @param {boolean} [unsavedChanges] If true, the user is asked if they want
	 * to navigate away from unsaved changes, allowing them to cancel the
	 * action.
	 * @returns {Promise<void>}
	 */
	async navigateToPath(path, unsavedChanges) {
		if (!unsavedChanges) {
			this.$location.url(path);
			return;
		}
		const params = {
			title: "Confirm Navigation",
			message: "You have unsaved changes that will be lost if you decide to continue.<br><br>Do you want to continue?"
		};
		const modalInstance = this.$uibModal.open({
			templateUrl: "common/modules/dialog/confirm/dialog.confirm.tpl.html",
			controller: "DialogConfirmController",
			size: "md",
			resolve: { params }
		});
		try {
			await modalInstance.result;
			this.$location.url(path);
		} catch {
			// this means the user cancelled
		}
	}
}

LocationUtils.$inject = ["$location", "$uibModal"];
module.exports = LocationUtils;
