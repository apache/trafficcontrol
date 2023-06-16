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
import "zone.js/node";

import { existsSync, readdirSync, readFileSync, statSync } from "fs";
import { createServer as createRedirectServer } from "http";
import { createServer, request, RequestOptions } from "https";
import { join } from "path";

import { APP_BASE_HREF } from "@angular/common";
import { ngExpressEngine } from "@nguniversal/express-engine";
import { ArgumentParser } from "argparse";
import * as express from "express";
import {
	defaultConfig,
	defaultConfigFile,
	getConfig,
	getVersion,
	ServerConfig,
	versionToString
} from "server.config";

import { AppServerModule } from "./src/main.server";

/**
 * StaticFile defines what compression files are available.
 */
interface StaticFile {
	compressions: Array<CompressionType>;
}

/**
 * CompressionType defines the different compression algorithms.
 */
interface CompressionType {
	fileExt: string;
	headerEncoding: string;
	name: string;
}

const gzip = {
	fileExt: "gz",
	headerEncoding: "gzip",
	name: "gzip"
};
const br = {
	fileExt: "br",
	headerEncoding: "br",
	name: "brotli"
};

/**
 * getFiles recursively gets all the files in a directory.
 *
 * @param path The path to get files from.
 * @returns Files found in the directory.
 */
function getFiles(path: string): string[] {
	const all = readdirSync(path)
		.map(file => join(path, file));
	const dirs = all
		.filter(file => statSync(file).isDirectory());
	let files = all
		.filter(file => !statSync(file).isDirectory());
	for (const dir of dirs) {
		files = files.concat(getFiles(dir));
	}
	return files;
}

let config: ServerConfig;

/**
 * The Express app is exported so that it can be used by serverless Functions.
 *
 * @param serverConfig Server configuration
 * @returns The Express.js application.
 */
export function app(serverConfig: ServerConfig): express.Express {
	const server = express();
	const indexHtml = join(serverConfig.browserFolder, "index.html");
	if (!existsSync(indexHtml)) {
		throw new Error(`Unable to start TP server, unable to find browser index.html at: ${indexHtml}`);
	}

	// Our Universal express-engine (found @ https://github.com/angular/universal/tree/master/modules/express-engine)
	server.engine("html", ngExpressEngine({
		bootstrap: AppServerModule
	}));

	server.set("view engine", "html");
	server.set("views", serverConfig.browserFolder);

	const allFiles = getFiles(serverConfig.browserFolder);
	const compressedFiles = new Map(allFiles
		.filter(file => file.match(/\.br|gz$/))
		.map(file => [file, undefined]));
	const foundFiles = new Map<string, StaticFile>(allFiles
		.filter(file => file.match(/\.js|css|tff|svg$/))
		.map(file => {
			const staticFile = {
				compressions: []
			} as StaticFile;
			if (compressedFiles.has(`${file}.${br.fileExt}`)) {
				staticFile.compressions.push(br);
			}
			if (compressedFiles.has(`${file}.${gzip.fileExt}`)) {
				staticFile.compressions.push(gzip);
			}
			return [file, staticFile];
		}));

	const typeMap = new Map([
		["js", "application/javascript"],
		["css", "text/css"],
		["ttf", "font/ttf"],
		["svg", "image/svg+xml"]
	]);
	// Could just use express compression `server.use(compression())` but that is calculated for each request
	server.get("*.(js|css|ttf|svg)", function(req, res, next) {
		const type = req.url.split(".").pop();
		if (type === undefined || !typeMap.has(type)) {
			return next();
		}
		const path = join(serverConfig.browserFolder, req.url.substring(1, req.url.length));
		const file = foundFiles.get(path);
		if(!file || file.compressions.length === 0) {
			return next();
		}
		const acceptedEncodings = req.acceptsEncodings();
		for(const compression of file.compressions) {
			if (acceptedEncodings.indexOf(compression.headerEncoding) === -1) {
				continue;
			}
			req.url = `${req.url}.${compression.fileExt}`;
			res.set("Content-Encoding", compression.headerEncoding);
			res.set("Content-Type", typeMap.get(type));
			console.log(`Serving ${compression.name} compressed file ${req.url}`);
			return next();
		}
		next();
	});
	// Example Express Rest API endpoints
	// server.get('/api/**', (req, res) => { });
	// Serve static files from /browser
	server.get("*.*", express.static(serverConfig.browserFolder, {
		maxAge: "1y"
	}));

	/**
	 * A handler for proxying the Traffic Ops API.
	 *
	 * @param req The client's request.
	 * @param res The server's response writer.
	 */
	function toProxyHandler(req: express.Request, res: express.Response): void {
		console.log(`Making TO API request to \`${req.originalUrl}\``);

		const fwdRequest: RequestOptions = {
			headers:            req.headers,
			host:               config.trafficOps.hostname,
			method:             req.method,
			path:               req.originalUrl,
			port:               config.trafficOps.port,
			rejectUnauthorized: !config.insecure,
		};

		try {
			const proxiedRequest = request(fwdRequest, (r) => {
				res.writeHead(r.statusCode ?? 502, r.headers);
				r.pipe(res);
			});
			req.pipe(proxiedRequest);
		} catch (e) {
			console.error("proxying request:", e);
		}
	}

	server.use("api/**", toProxyHandler);
	server.use("/api/**", toProxyHandler);

	// All regular routes use the Universal engine
	server.get("*", (req, res) => {
		res.render(indexHtml, { providers: [{ provide: APP_BASE_HREF, useValue: req.baseUrl }], req });
	});

	server.enable("trust proxy");
	return server;
}

/**
 * Runs the server.
 *
 * @returns An exit code for the process.
 */
function run(): number {
	const version = getVersion();
	const parser = new ArgumentParser({
		// Nothing I can do about this, library specifies its interface.
		/* eslint-disable @typescript-eslint/naming-convention */
		add_help: true,
		/* eslint-enable @typescript-eslint/naming-convention */
		description: "Traffic Portal re-written in modern Angular"
	});
	parser.add_argument("-t", "--traffic-ops", {
		dest: "trafficOps",
		help: "Specify the Traffic Ops host/URL, including port. (Default: uses the `TO_URL` environment variable)",
		type: (arg: string) => {
			try {
				return new URL(arg);
			} catch (e) {
				if (e instanceof TypeError) {
					return new URL(`https://${arg}`);
				}
				throw e;
			}
		}
	});
	parser.add_argument("-k", "--insecure", {
		action: "store_true",
		help: "Skip Traffic Ops server certificate validation. This affects requests from Traffic Portal to Traffic Ops AND signature" +
			" verification of any passed SSL keys/certificates"
	});
	parser.add_argument("-p", "--port", {
		default: defaultConfig.port,
		help: "Specify the port on which Traffic Portal will listen (Default: 4200)",
		type: "int"
	});
	parser.add_argument("-c", "--cert-path", {
		dest: "certPath",
		help: "Specify a location for an SSL certificate to be used by Traffic Portal. (Requires `-K`/`--key-path`. If both are omitted," +
			" will serve using HTTP)",
		type: "str"
	});
	parser.add_argument("-d", "--browser-folder", {
		default: defaultConfig.browserFolder,
		dest: "browserFolder",
		help: "Specify location for the folder that holds the browser files",
		type: "str"
	});
	parser.add_argument("-K", "--key-path", {
		dest: "keyPath",
		help: "Specify a location for an SSL certificate to be used by Traffic Portal. (Requires `-c`/`--cert-path`. If both are omitted," +
			" will serve using HTTP)",
		type: "str"
	});
	parser.add_argument("-C", "--config-file", {
		default: defaultConfigFile,
		dest: "configFile",
		help: "Specify a path to a configuration file - options are overridden by command-line flags.",
		type: "str"
	});
	parser.add_argument("-v", "--version", {
		action: "version",
		version: versionToString(version)
	});

	try {
		config = getConfig(parser.parse_args(), version);
	} catch (e) {
		console.error(`Failed to initialize server configuration: ${e}`);
		return 1;
	}

	// Start up the Node server
	const server = app(config);

	if (config.useSSL) {
		let cert: string;
		let key: string;
		let ca: Array<string>;
		try {
			cert = readFileSync(config.certPath, {encoding: "utf8"});
			key = readFileSync(config.keyPath, {encoding: "utf8"});
			ca = config.certificateAuthPaths.map(c => readFileSync(c, {encoding: "utf8"}));
		} catch (e) {
			console.error("reading SSL key/cert:", e);
			return 1;
		}
		createServer(
			{
				ca,
				cert,
				key,
				rejectUnauthorized: !config.insecure,
			},
			server
		).listen(config.port, ()=> {
			console.log(`Node Express server listening on port ${config.port}`);
		});
		try {
			const redirectServer = createRedirectServer(
				(req, res) => {
					if (!req.url) {
						res.statusCode = 500;
						console.error("got HTTP request for redirect that had no URL");
						res.end();
						return;
					}
					res.statusCode = 308;
					res.setHeader("Location", req.url.replace(/^[hH][tT][tT][pP]:/, "https:"));
					res.end();
				}
			);
			redirectServer.listen(80);
			redirectServer.on("error", e => {
				console.error(`redirect server encountered error: ${e}`);
				if (Object.prototype.hasOwnProperty.call(e, "code") && (e as typeof e & {code: unknown}).code === "EACCES") {
					console.warn("access to port 80 not allowed; closing redirect server");
					redirectServer.close();
				}
			});
		} catch (e) {
			console.warn("Failed to initialize HTTP-to-HTTPS redirect listener:", e);
		}
	} else {
		server.listen(config.port, () => {
			console.log(`Node Express server listening on port ${config.port}`);
		});
	}
	return 0;
}

// Webpack will replace 'require' with '__webpack_require__'
// '__non_webpack_require__' is a proxy to Node 'require'
// The below code is to ensure that the server is run only when not requiring the bundle.
/* eslint-disable no-underscore-dangle */
// eslint-disable-next-line @typescript-eslint/naming-convention
declare const __non_webpack_require__: NodeRequire;
try {
	/* eslint-enable no-underscore-dangle */
	const mainModule = __non_webpack_require__.main;
	const moduleFilename = mainModule && mainModule.filename || "";
	if (moduleFilename === __filename || moduleFilename.includes("iisnode")) {
		const code = run();
		if (code) {
			process.exit(code);
		}
	}
} catch(e) {
	console.error("Encountered error while running server:", e);
	process.exit(1);
}
export * from "./src/main.server";
