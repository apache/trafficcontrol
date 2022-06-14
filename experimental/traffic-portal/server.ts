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


import { existsSync, readFileSync } from "fs";
import { createServer as createRedirectServer } from "http";
import { createServer, request } from "https";
import { join } from "path";

import { APP_BASE_HREF } from "@angular/common";
import { ngExpressEngine } from "@nguniversal/express-engine";
import { ArgumentParser } from "argparse";
import * as express from "express";
import { getConfig, getVersion, ServerConfig, versionToString } from "server.config";

import { AppServerModule } from "./src/main.server";

let config: ServerConfig;

/**
 * The Express app is exported so that it can be used by serverless Functions.
 *
 * @returns The Express.js application.
 */
export function app(): express.Express {
	const server = express();
	const distFolder = join(process.cwd(), "dist/traffic-portal/browser");
	const indexHtml = existsSync(join(distFolder, "index.original.html")) ? "index.original.html" : "index";

	// Our Universal express-engine (found @ https://github.com/angular/universal/tree/master/modules/express-engine)
	server.engine("html", ngExpressEngine({
		bootstrap: AppServerModule,
	}));

	server.set("view engine", "html");
	server.set("views", distFolder);

	// Example Express Rest API endpoints
	// server.get('/api/**', (req, res) => { });
	// Serve static files from /browser
	server.get("*.*", express.static(distFolder, {
		maxAge: "1y"
	}));

	// All regular routes use the Universal engine
	server.get("*", (req, res) => {
		res.render(indexHtml, { providers: [{ provide: APP_BASE_HREF, useValue: req.baseUrl }], req });
	});

	server.use("/api/**", (req, res) => {
		console.log(`Making TO API request to \`${req.originalUrl}\``);

		let origURL: URL;
		try {
			origURL = new URL(req.originalUrl);
		} catch (err) {
			console.error(`Failed to parse request URL ${req.originalUrl} as a URL: ${err}`);
			res.statusCode = 502;
			res.setHeader("Content-Type", "application/json");
			res.write('{"alerts":[{"level":"error","text":"Traffic Ops is unreachable"}]}');
			return;
		}

		const fwdRequest = {
			headers: req.headers,
			host:    config.trafficOps.hostname,
			method:  req.method,
			path:    origURL.pathname+origURL.search,
			port:    config.trafficOps.port,
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
		res.end();
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
		action: "storeTrue",
		help: "Skip Traffic Ops server certificate validation. This affects requests from Traffic Portal to Traffic Ops AND signature" +
			" verification of any passed SSL keys/certificates"
	});
	parser.add_argument("-p", "--port", {
		default: 4200,
		help: "Specify the port on which Traffic Portal will listen (Default: 4200)",
		type: "int"
	});
	parser.add_argument("-c", "--cert-path", {
		dest: "certPath",
		help: "Specify a location for an SSL certificate to be used by Traffic Portal. (Requires `-K`/`--key-path`. If both are omitted," +
			" will serve using HTTP)",
		type: "str"
	});
	parser.add_argument("-K", "--key-path", {
		dest: "keyPath",
		help: "Specify a location for an SSL certificate to be used by Traffic Portal. (Requires `-c`/`--cert-path`. If both are omitted," +
			" will serve using HTTP)",
		type: "str"
	});
	parser.add_argument("-C", "--config-file", {
		default: "/etc/traffic-portal/config.js",
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
	const server = app();

	if (config.useSSL) {
		let cert: string;
		let key: string;
		try {
			cert = readFileSync(config.certPath, {encoding: "utf8"});
			key = readFileSync(config.keyPath, {encoding: "utf8"});
		} catch (e) {
			console.error("reading SSL key/cert:", e);
			return 1;
		}
		createServer(
			{
				cert,
				key,
				rejectUnauthorized: !config.insecure
			},
			server
		).listen(config.port, ()=> {
			console.log(`Node Express server listening on port ${config.port}`);
		});
		try {
			createRedirectServer(
				(req, res) => {
					if (!req.url) {
						res.statusCode = 500;
						console.error("got HTTP request for redirect that had no URL");
						res.end();
						return;
					}
					res.statusCode = 308;
					res.setHeader("Location", req.url.replace(/^[hH][tT][tT][pP]:/, "https"));
					res.end();
				}
			).listen(80);
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
/* eslint-enable no-underscore-dangle */
const mainModule = __non_webpack_require__.main;
const moduleFilename = mainModule && mainModule.filename || "";
if (moduleFilename === __filename || moduleFilename.includes("iisnode")) {
	process.exit(run());
}

export * from "./src/main.server";
