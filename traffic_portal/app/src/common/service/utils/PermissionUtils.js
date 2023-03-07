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
 * PermissionUtils provides methods for checking the user's Permissions.
 */
class PermissionUtils {
	/**
	 * @param {import("../../models/UserModel")} userModel
	 */
	constructor(userModel) {
		this.userModel = userModel;
	}

	/**
	 * Checks if the user has the given "Capability" (now called a
	 * "Permission").
	 *
	 * @deprecated "Capabilities" have been (more or less) renamed to
	 * Permissions, so further checks should use `hasPermission` instead. Note
	 * also that this doesn't check for the special "admin" Role that is
	 * afforded every Permission.
	 *
	 * @param {string} cap
	 * @returns {boolean}
	 */
	hasCapability(cap) {
		return this.userModel.hasCapability(cap);
	}

	/**
	 * Checks if the user has the given Permission.
	 *
	 * @param {string} permission
	 * @returns {boolean}
	 */
	hasPermission(permission) {
		return this.userModel.user.role === "admin" || this.userModel.hasCapability(permission);
	}
}

PermissionUtils.$inject = ["userModel"];
module.exports = PermissionUtils;
