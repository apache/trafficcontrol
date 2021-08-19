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
import { hasProperty } from "./CommonUtils/utils";

/** Possible levels of TO Alerts */
export type AlertLevel = "success" | "info" | "warning" | "error";

/**
 * Checks whether an object is a valid Alert level.
 *
 * @param s The object to check.
 * @returns Whether or not `s` is an AlertLevel.
 */
export function isAlertLevel(s: unknown): s is AlertLevel {
    if (typeof(s) !== "string") {
        return false;
    }
    switch (s) {
        case "success":
        case "info":
        case "warning":
        case "error":
            return true;
    }
    return false;
}

/** TO API Alerts */
export interface Alert {
    level: AlertLevel;
    text: string;
}

/**
 * Checks whether an object is an Alert like the ones TO normally returns.
 *
 * @param a The object to check.
 * @returns Whether or not `a` is an Alert.
 */
export function isAlert(a: unknown): a is Alert {
    if (typeof(a) !== "object" || a === null) {
        return false;
    }
    if (!hasProperty(a, "level") || !hasProperty(a, "text", "string")) {
        return false;
    }
    return isAlertLevel(a.level);
}

/**
 * Logs an alert to the appropriate console stream based on its `level`.
 *
 * @param a The Alert to log.
 * @param prefix Optional prefix for the log message
 */
export function logAlert(a: Alert, prefix?: string): void {
    let logfunc;
    let pre = (prefix ?? "").trimStart();
    switch (a.level) {
        case "success":
            logfunc = console.log;
            pre = `SUCCESS: ${pre}`;
            break;
        case "info":
            logfunc = console.info
            pre = `INFO: ${pre}`;
            break;
        case "warning":
            logfunc = console.warn
            pre = `WARN: ${pre}`;
            break;
        case "error":
            logfunc = console.error
            pre = `ERROR: ${pre}`;
            break;
    }
    logfunc(pre, a.text);
}

/** TestingConfig is the type of a testing configuration. */
export interface TestingConfig {
	/** This is login information for a user with admin-level permissions. */
	readonly login: {
		readonly password: string;
		readonly username: string;
	};
	/** The URL at which the Traffic Ops API can be accessed. */
	readonly apiUrl: string;
	/** The URL at which Traffic Portal is served - root path. */
	readonly baseUrl: string;
	/** Logging alert levels that are enabled. */
	readonly alertLevels?: Array<AlertLevel>;
}

/**
 * Checks if a passed object is a valid testing configuration.
 *
 * @param c The object to check.
 * @returns `true`, always. When the check fails, it throws an error that
 * explains why.
 */
export function isTestingConfig(c: unknown): c is TestingConfig {
	if (typeof(c) !== "object") {
		throw new Error(`testing configuration must be an object, not a '${typeof(c)}'`);
	}
	if (c === null) {
		throw new Error("testing configuration must be an object, not 'null'");
	}

	if (!hasProperty(c, "login") || typeof(c.login) !== "object" || c.login === null) {
		throw new Error("missing or invalid 'login' property");
	}
	if (!hasProperty(c.login, "password", "string") || !hasProperty(c.login, "username", "string")) {
		throw new Error("'login' property has missing and/or invalid 'password' and/or 'username' property(ies)");
	}
	if (c.login.username === "" || c.login.password === "") {
		throw new Error("neither 'login.username' nor 'login.password' may be blank");
	}
	if (!hasProperty(c, "apiUrl", "string")) {
		throw new Error("missing or invalid 'apiUrl' property");
	}
	try {
		new URL(c.apiUrl);
	} catch (e) {
		throw new Error(`'apiUrl' is not a valid URL: ${c.apiUrl}`);
	}
	let baseURL;
	if (!hasProperty(c, "baseUrl", "string")) {
		throw new Error("missing or invalid 'baseUrl' property");
	}
	try {
		baseURL = new URL(c.baseUrl);
	} catch (e) {
		throw new Error(`'baseUrl' is not a valid URL: ${c.baseUrl}`);
	}
	if (baseURL.pathname !== "/") {
		throw new Error("'baseUrl' must be a root path");
	}
	if (!hasProperty(c, "alertLevels")) {
		return true;
	}
	if (!(c.alertLevels instanceof Array)) {
		throw new Error("'alertLevels' must be an array");
	}
	if (!c.alertLevels.every(isAlertLevel)) {
		throw new Error(`invalid alert levels: ${c.alertLevels}`);
	}
	return true;
}
