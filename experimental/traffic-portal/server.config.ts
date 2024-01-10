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

// Logging cannot be initialized until after the job of the routines in this
// file are complete.
/* eslint-disable no-console */

import { execSync } from "child_process";
import { access, constants, readFile, readdir, realpath } from "fs/promises";
import { join, sep } from "path";

import { hasProperty } from "src/app/utils";

/**
 * A Node system error. I don't know why but this isn't exposed by Node - it's a
 * class but you won't be able to use `instanceof` - and isn't present in Node
 * typings. I copied the properties and their descriptions from the NodeJS
 * documentation.
 */
type SystemError = Error & {
	/** If present, the address to which a network connection failed. */
	readonly address?: string;
	/** The string error code. */
	readonly code: string;
	/** If present, the file path destination when reporting a file system error. */
	readonly dest?: string;
	/** The system-provided error number. */
	readonly errno: number;
	/** If present, extra details about the error condition. */
	readonly info?: unknown;
	/** A system-provided human-readable description of the error. */
	readonly message: string;
	/** If present, the file path when reporting a file system error. */
	readonly path?: string;
	/** If present, the network connection port that is not available. */
	readonly port?: number;
	/** The name of the system call that triggered the error. */
	readonly syscall: string;
};

/**
 * Checks if an {@link Error} is a {@link SystemError}.
 *
 * @param e The {@link Error} to check.
 * @returns Whether `e` is a {@link SystemError}.
 */
function isSystemError(e: Error): e is SystemError {
	return hasProperty(e, "code");
}

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

	if (!hasProperty(v, "version", "string")) {
		console.error("version required field 'version' missing or invalid");
		return false;
	}

	if (hasProperty(v, "commits") && typeof(v.commits) !== "string") {
		console.error(`version property 'commits' has incorrect type; want: string, got: ${typeof(v.commits)}`);
		return false;
	}
	if (hasProperty(v, "hash") && typeof(v.hash) !== "string") {
		console.error(`version property 'hash' has incorrect type; want: string, got: ${typeof(v.hash)}`);
		return false;
	}
	if (hasProperty(v, "elRelease") && typeof(v.elRelease) !== "string") {
		console.error(`version property 'elRelease' has incorrect type; want: string, got: ${typeof(v.elRelease)}`);
		return false;
	}
	if (hasProperty(v, "arch") && typeof(v.arch) !== "string") {
		console.error(`version property 'arch' has incorrect type; want: string, got: ${typeof(v.arch)}`);
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
	/** Path to the folder containing browser files. **/
	browserFolder: string;
	/** The URL of the Traffic Portal V1. */
	tpv1Url: URL;
}

/**
 * The properties a ServerConfig has if SSL is in use.
 */
interface ConfigWithSSL {
	/** The path to the SSL certificate Traffic Portal will use. */
	certPath: string;
	/** The path to the SSL private key Traffic Portal will use. */
	keyPath: string;
	/** The paths to trusted root certificates, setting this is equivalent
	 * to the path to the environment variable NODE_EXTRA_CA_CERTS */
	certificateAuthPaths: Array<string>;
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

	if (hasProperty(c, "insecure")) {
		if (typeof(c.insecure) !== "boolean") {
			throw new Error("'insecure' must be a boolean");
		}
	} else {
		(c as {insecure: boolean}).insecure = false;
	}
	if (!hasProperty(c, "port", "number")) {
		throw new Error("required configuration for 'port' is missing or not a valid number");
	}
	if (!hasProperty(c, "trafficOps", "string")) {
		throw new Error("required configuration for 'trafficOps' is missing or not a string");
	}
	if (!hasProperty(c, "browserFolder", "string")) {
		throw new Error("required configuration for 'browserFolder' is missing or not a string");
	}
	if(!hasProperty(c, "tpv1Url", "string")){
		throw new Error("required configuration for 'tpv1Url' is missing or not a string");
	}

	try {
		(c as {trafficOps: URL | string}).trafficOps = new URL(c.trafficOps);
	} catch (e) {
		throw new Error(`'trafficOps' is not a valid URL: ${e}`);
	}

	try {
		(c as {tpv1Url: URL | string}).tpv1Url = new URL(c.tpv1Url);
	} catch (e) {
		throw new Error(`'tpv1Url' is not a valid URL: ${e}`);
	}

	if (hasProperty(c, "useSSL")) {
		if (typeof(c.useSSL) !== "boolean") {
			throw new Error("'useSSL' must be a boolean");
		}
		if (c.useSSL) {
			if (!hasProperty(c, "certPath", "string")) {
				throw new Error("missing or invalid 'certPath' - required to use SSL");
			}
			if (!hasProperty(c, "keyPath", "string")) {
				throw new Error("missing or invalid 'keyPath' - required to use SSL");
			}
		}
	}

	return true;
}

const defaultVersionFile = "/etc/traffic-portal/version.json";

/**
 * Searches recursively upward through the filesystem to find a file named
 * "VERSION" and returns the real, absolute path to that file.
 *
 * @param path The path from which to begin the search.
 * @returns The path to the VERSION file, assuming it was found.
 * @throws {Error} If no VERSION file could be found in `path` or any of its
 * ancestor directories.
 * @throws {SystemError} If the given path isn't a directory, or directory
 * traversal fails for some reason.
 */
async function findVersionFile(path: string = "."): Promise<string> {
	for (const ent of await readdir(path)) {
		if (ent === "VERSION") {
			return realpath(join(path, ent));
		}
	}
	path = await realpath(join(path, ".."));
	if (path === sep) {
		throw new Error("VERSION file not found");
	}
	return findVersionFile(path);
}

/**
 * Retrieves the server's version from the file path provided.
 *
 * @param path The path to a version file containing a ServerVersion object.
 * Defaults to /etc/traffic-portal/version.json. If this file doesn't exist,
 * the version may be deduced from the execution environment using git and
 * looking for the ATC VERSION file.
 * @returns The parsed server version.
 */
export async function getVersion(path?: string): Promise<ServerVersion> {
	if (!path) {
		path = defaultVersionFile;
	}

	try {
		const v = JSON.parse(await readFile(path, {encoding: "utf8"}));
		if (isServerVersion(v)) {
			return v;
		}
		throw new Error(`contents of version file '${path}' does not represent an ATC version`);
	} catch (e) {
		if (e instanceof Error && isSystemError(e)) {
			if (e.code !== "ENOENT") {
				throw new Error(`file at "${path}" could not be read: ${e.message}`);
			}
		} else {
			throw new Error(`file at "${path}" could not be read: ${e}`);
		}
	}

	let versionFilePath: string;
	try {
		versionFilePath = await findVersionFile();
	} catch (e) {
		throw new Error(`'${path}' doesn't exist and couldn't find a VERSION file from which to read a server version: ${e}`);
	}

	const ver: ServerVersion = {
		version: (await readFile(versionFilePath, {encoding: "utf8"})).trimEnd()
	};

	try {
		ver.hash = execSync("git rev-parse --short=8 HEAD", {encoding: "utf8"}).trimEnd();
		ver.commits = String(execSync("git describe --long --tags " +
			"--match=RELEASE-[0-9].[0-9].[0-9] --match=RELEASE-[0-9][0-9].[0-9][0-9].[0-9][0-9] " +
			"--match=v[0-9].[0-9].[0-9] --match=v[0-9][0-9].[0-9][0-9].[0-9][0-9]", {encoding: "utf8"}).split("-").slice(-2)[0]);
	} catch (e) {
		ver.commits = "0";
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
	tpv1Url?: URL;
	insecure: boolean;
	port: number;
	certPath?: string;
	keyPath?: string;
	configFile: string;
	browserFolder: string;
}

export const defaultConfigFile = "/etc/traffic-portal/config.json";

export const defaultConfig: ServerConfig = {
	browserFolder: "/opt/traffic-portal/browser",
	insecure: false,
	port: 4200,
	tpv1Url: new URL("https://example.com"),
	trafficOps: new URL("https://example.com"),
	version: { version: "" }
};
/**
 * Gets the configuration for the Traffic Portal server.
 *
 * @param args The arguments passed to Traffic Portal.
 * @param ver The version to use for the server.
 * @returns A full configuration for the server.
 */
export async function getConfig(args: Args, ver: ServerVersion): Promise<ServerConfig> {
	let cfg = defaultConfig;
	cfg.version = ver;

	let readFromFile = false;
	try {
		const cfgFromFile = JSON.parse(await readFile(args.configFile, {encoding: "utf8"}));
		if (isConfig(cfgFromFile)) {
			cfg = cfgFromFile;
			cfg.version = ver;
		} else {
			throw new Error("bad contents; doesn't represent a configuration file");
		}
		readFromFile = true;
	} catch (err) {
		const msg = `invalid configuration file at '${args.configFile}'`;
		if (err instanceof Error) {
			if (!isSystemError(err) || (err.code !== "ENOENT" || args.configFile !== defaultConfigFile)) {
				throw new Error(`${msg}: ${err.message}`);
			}
		} else {
			throw new Error(`${msg}: ${err}`);
		}
	}

	if(args.browserFolder !== defaultConfig.browserFolder) {
		cfg.browserFolder = args.browserFolder;
	}

	try {
		if (!(await readdir(cfg.browserFolder)).includes("index.html")) {
			throw new Error("directory doesn't include an 'index.html' file");
		}
	} catch (e) {
		throw new Error(`setting browser directory: ${e instanceof Error ? e.message : e}`);
	}

	if(args.port !== defaultConfig.port) {
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

	if (args.tpv1Url) {
		cfg.tpv1Url = args.tpv1Url;
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
				browserFolder: cfg.browserFolder,
				certPath: args.certPath,
				certificateAuthPaths: [],
				insecure: cfg.insecure,
				keyPath: args.keyPath,
				port: cfg.port,
				tpv1Url: cfg.tpv1Url,
				trafficOps: cfg.trafficOps,
				useSSL: true,
				version: ver
			};
		} else if (args.keyPath) {
			throw new Error("must specify either both a key path and a cert path, or neither");
		}
	}

	if(args.insecure) {
		cfg.insecure = args.insecure;
	}

	if (cfg.useSSL) {
		try {
			await access(cfg.certPath, constants.R_OK);
		} catch (e) {
			throw new Error(`checking certificate file "${cfg.certPath}": ${e instanceof Error ? e.message : e}`);
		}
		try {
			await access(cfg.keyPath, constants.R_OK);
		} catch (e) {
			throw new Error(`checking key file "${cfg.keyPath}": ${e instanceof Error ? e.message : e}`);
		}
	}

	return cfg;
}
