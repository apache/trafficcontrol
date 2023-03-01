/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import { execSync } from "child_process";
import { existsSync, readFileSync } from "fs";

/**
 * ServerVersion contains versioning information for the server,
 * consistent with what other components provide, even if some
 * of it doesn't really make sense for a Node server.
 */
export interface ServerVersion {
	/** The ATC version of the server e.g. 5.0.0 */
	version: string;
	/** Datetime when the current version was built*/
	date?: string;
	/**
	 * The number of commits in the development branch that produced this
	 * version of ATC - if known.
	 */
	commits?: string;
	/**
	 * The hash of the commit at which this ATC instance was built - if known.
	 */
	hash?: string;
	/**
	 * The "Enterprise Linux" release for which this version of ATC was built
	 * (which should never actually matter for Traffic Portal) - if known.
	 */
	elRelease?: string;
	/**
	 * The CPU architecture for which this version of ATC was built (which
	 * should REALLY never matter for Traffic Portal) - if known.
	 */
	arch?: string;
}

/**
 * Converts the given version to a string.
 *
 * @param v The server version to convert.
 * @returns A string representation of `v`.
 */
export function versionToString(v: ServerVersion): string {
	let ret = `traffic-portal-${v.version}`;

	// Can't allow either-or, because since git hashes could
	// possibly be entirely numeric that would make the case
	// where only one is present confusing - is it commit
	// count or hash?
	if (v.commits) {
		ret += `-${v.commits}`;
		if (v.hash) {
			ret += `.${v.hash}`;
		}
	}

	if (v.elRelease) {
		ret += `.${v.elRelease}`;
	}

	if (v.arch) {
		ret += `.${v.arch}`;
	}

	return ret;
}

/**
 * Checks if some unknown Javascript value is a valid {@link ServerVersion}.
 *
 * @param v The object to check.
 * @returns Whether `v` is a valid {@link ServerVersion}.
 */
function isServerVersion(v: unknown): v is ServerVersion {
	if (typeof(v) !== "object") {
		console.error("version does not represent a server version");
		return false;
	}
	if (!v) {
		console.error("'null' is not a valid server version");
		return false;
	}

	if (!Object.prototype.hasOwnProperty.call(v, "version")) {
		console.error("version missing required field 'version'");
		return false;
	}
	if (typeof((v as {version: unknown}).version) !== "string") {
		return false;
	}

	if (Object.prototype.hasOwnProperty.call(v, "commits") && (typeof((v as {commits: unknown}).commits)) !== "string") {
		console.error(`version property 'commits' has incorrect type; want: string, got: ${typeof((v as {commits: unknown}).commits)}`);
		return false;
	}
	if (Object.prototype.hasOwnProperty.call(v, "hash") && (typeof((v as {hash: unknown}).hash)) !== "string") {
		console.error(`version property 'hash' has incorrect type; want: string, got: ${typeof((v as {hash: unknown}).hash)}`);
		return false;
	}
	if (Object.prototype.hasOwnProperty.call(v, "elRelease") && (typeof((v as {elRelease: unknown}).elRelease)) !== "string") {
		console.error(
			`version property 'elRelease' has incorrect type; want: string, got: ${typeof (v as {elRelease: unknown}).elRelease}`
		);
		return false;
	}
	if (Object.prototype.hasOwnProperty.call(v, "arch") && (typeof((v as {arch: unknown}).arch)) !== "string") {
		console.error(`version property 'arch' has incorrect type; want: string, got: ${typeof((v as {arch: unknown}).arch)}`);
		return false;
	}
	return true;
}

/**
 * The base properties shared by all ServerConfigs.
 */
interface BaseConfig {
	/** Whether or not SSL certificate errors from Traffic Ops will be ignored. */
	insecure?: boolean;
	/** The port on which Traffic Portal listens. */
	port: number;
	/** The URL of the Traffic Ops API. */
	trafficOps: URL;
	/** Whether or not to serve HTTPS */
	useSSL?: boolean;
	/** Contains all of the versioning information. */
	version: ServerVersion;
}

/**
 * The properties a ServerConfig has if SSL is in use.
 */
interface ConfigWithSSL {
	/** The path to the SSL certificate Traffic Portal will use. */
	certPath: string;
	/** The path to the SSL private key Traffic Portal will use. */
	keyPath: string;
	/** Whether or not to serve HTTPS */
	useSSL: true;
}

/**
 * The properties a ServerConfig has if SSL is not in use.
 */
interface ConfigWithoutSSL {
	/** Whether or not to serve HTTPS */
	useSSL?: false;
}

/** ServerConfig holds server configuration. */
export type ServerConfig = BaseConfig & (ConfigWithSSL | ConfigWithoutSSL);

/**
 * isConfig checks specifically the contents of configuration files
 * passed through JSON.parse, so it doesn't validate the 'version', since
 * that doesn't need to be in the configuration file.
 *
 * @param c The ostensibly configuration object to check.
 * @returns Whether c is a valid server configuration.
 */
function isConfig(c: unknown): c is ServerConfig {
	if (typeof(c) !== "object") {
		throw new Error("configuration is not an object");
	}
	if (!c) {
		throw new Error("'null' is not a valid configuration");
	}

	if (Object.prototype.hasOwnProperty.call(c, "insecure")) {
		if (typeof((c as {insecure: unknown}).insecure) !== "boolean") {
			throw new Error("'insecure' must be a boolean");
		}
	} else {
		(c as {insecure: boolean}).insecure = false;
	}
	if (!Object.prototype.hasOwnProperty.call(c, "port")) {
		throw new Error("'port' is required");
	}
	if (typeof((c as {port: unknown}).port) !== "number") {
		throw new Error("'port' must be a number");
	}
	if (!Object.prototype.hasOwnProperty.call(c, "trafficOps")) {
		throw new Error("'trafficOps' is required");
	}
	if (typeof((c as {trafficOps: unknown}).trafficOps) !== "string") {
		throw new Error("'trafficOps' must be a string");
	}

	try {
		(c as {trafficOps: URL}).trafficOps = new URL((c as {trafficOps: string}).trafficOps);
	} catch (e) {
		throw new Error(`'trafficOps' is not a valid URL: ${e}`);
	}

	if (Object.prototype.hasOwnProperty.call(c, "useSSL")) {
		if (typeof((c as {useSSL: unknown}).useSSL) !== "boolean") {
			throw new Error("'useSSL' must be a boolean");
		}
		if ((c as {useSSL: boolean}).useSSL) {
			if (!Object.prototype.hasOwnProperty.call(c, "certPath")) {
				throw new Error("'certPath' is required to use SSL");
			}
			if (typeof((c as {certPath: unknown}).certPath) !== "string") {
				throw new Error("'certPath' must be a string");
			}
			if (!Object.prototype.hasOwnProperty.call(c, "keyPath")) {
				throw new Error("'keyPath' is required to use SSL");
			}
			if (typeof((c as {keyPath: unknown}).keyPath) !== "string") {
				throw new Error("'keyPath' must be a string");
			}
		}
	}

	return true;
}

const defaultVersionFile = "/etc/traffic-portal/version.json";

/**
 * Retrieves the server's version from the file path provided.
 *
 * @param path The path to a version file containing a ServerVersion object.
 * Defaults to /etc/traffic-portal/version.json. If this file doesn't exist,
 * the version may be deduced from the execution environment using git and
 * looking for the ATC VERSION file.
 * @returns The parsed server version.
 */
export function getVersion(path?: string): ServerVersion {
	if (!path) {
		path = defaultVersionFile;
	}

	if (existsSync(path)) {
		const v = JSON.parse(readFileSync(path, {encoding: "utf8"}));
		if (isServerVersion(v)) {
			return v;
		}
		throw new Error(`contents of version file '${path}' does not represent an ATC version`);
	}

	if (!existsSync("../../../../VERSION")) {
		throw new Error(`'${path}' doesn't exist and '../../../../VERSION' doesn't exist`);
	}
	const ver: ServerVersion = {
		version: readFileSync("../../../../VERSION", {encoding: "utf8"}).trimEnd()
	};

	try {
		ver.commits = String(execSync("git rev-list HEAD", {encoding: "utf8"}).split("\n").length);
		ver.hash = execSync("git rev-parse --short=8 HEAD", {encoding: "utf8"}).trimEnd();
	} catch (e) {
		console.warn("getting git parts of version:", e);
	}

	try {
		const releaseNo = execSync("rpm -q --qf '%{version}' -f /etc/redhat-release", {encoding: "utf8"}).trimEnd();
		ver.elRelease = `el${releaseNo}`;
	} catch (e) {
		ver.elRelease = "el7";
		console.warn(`getting RHEL version: ${e}`);
	}

	try {
		ver.arch = execSync("uname -m", {encoding: "utf8"});
	} catch (e) {
		console.warn(`getting system architecture: ${e}`);
	}

	return ver;
}

/** The type of command line arguments to Traffic Portal. */
interface Args {
	trafficOps?: URL;
	insecure?: boolean;
	port?: number;
	certPath?: string;
	keyPath?: string;
	configFile: string;
}

export const defaultConfigFile = "/etc/traffic-portal/config.js";

/**
 * Gets the configuration for the Traffic Portal server.
 *
 * @param args The arguments passed to Traffic Portal.
 * @param ver The version to use for the server.
 * @returns A full configuration for the server.
 */
export function getConfig(args: Args, ver: ServerVersion): ServerConfig {
	let cfg: ServerConfig = {
		insecure: false,
		port: 4200,
		trafficOps: new URL("https://example.com"),
		useSSL: false,
		version: ver
	};

	let readFromFile = false;
	if (existsSync(args.configFile)) {
		const cfgFromFile = JSON.parse(readFileSync(args.configFile, {encoding: "utf8"}));
		try {
			if (isConfig(cfgFromFile)) {
				cfg = cfgFromFile;
				cfg.version = ver;
			}
		} catch (err) {
			throw new Error(`invalid configuration file at '${args.configFile}': ${err}`);
		}
		readFromFile = true;
	} else if (args.configFile !== defaultConfigFile) {
		throw new Error(`no such configuration file: ${args.configFile}`);
	}

	if (args.port) {
		cfg.port = args.port;
	}
	if (isNaN(cfg.port) || cfg.port <= 0 || cfg.port > 65535) {
		throw new Error(`invalid port: ${cfg.port}`);
	}

	if (args.trafficOps) {
		cfg.trafficOps = args.trafficOps;
	} else if (!readFromFile) {
		const envURL = process.env.TO_URL;
		if (!envURL) {
			throw new Error("Traffic Ops URL must be specified");
		}
		try {
			cfg.trafficOps = new URL(envURL);
		} catch (e) {
			throw new Error(`invalid Traffic Ops URL from environment: ${envURL}`);
		}
	}

	if (readFromFile && cfg.useSSL) {
		if (args.certPath) {
			cfg.certPath = args.certPath;
		}
		if (args.keyPath) {
			cfg.keyPath = args.keyPath;
		}
	} else if (!readFromFile || cfg.useSSL === undefined) {
		if (args.certPath) {
			if (!args.keyPath) {
				throw new Error("must specify either both a key path and a cert path, or neither");
			}
			cfg = {
				certPath: args.certPath,
				insecure: cfg.insecure,
				keyPath: args.keyPath,
				port: cfg.port,
				trafficOps: cfg.trafficOps,
				useSSL: true,
				version: ver
			};
		} else if (args.keyPath) {
			throw new Error("must specify either both a key path and a cert path, or neither");
		}
	}

	if(args.insecure === true) {
		cfg.insecure = args.insecure;
	}

	if (cfg.useSSL) {
		if (!existsSync(cfg.certPath)) {
			throw new Error(`no such certificate file: ${cfg.certPath}`);
		}
		if (!existsSync(cfg.keyPath)) {
			throw new Error(`no such key file: ${cfg.keyPath}`);
		}
	}

	return cfg;
}
