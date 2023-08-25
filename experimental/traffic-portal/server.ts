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

import { readFileSync } from "fs";
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
	type ServerConfig,
	versionToString
} from "server.config";

import { hasProperty, Logger, LogLevel } from "src/app/utils";
import { environment } from "src/environments/environment";
import { AppServerModule } from "src/main.server";

import { errorMiddleWare, loggingMiddleWare, type TPResponseWriter } from "./middleware";

const typeMap = new Map([
	["js", "application/javascript"],
	["css", "text/css"],
	["ttf", "font/ttf"],
	["svg", "image/svg+xml"]
]);

/**
 * A handler for serving files from compressed variants.
 *
 * @param req The client request.
 * @param res Response writer.
 * @param next A delegation to the next handler, to be called if this handler
 * determines it can't write a response (which is always because this handler
 * doesn't do that).
 * @returns nothing. This is just required because we're returning void function
 * calls. Not actually sure why, seems like a bug to me.
 */
function compressedFileHandler(req: express.Request, res: TPResponseWriter, next: express.NextFunction): void {
	const type = req.path.split(".").pop();
	if (type === undefined || !typeMap.has(type)) {
		res.locals.logger.debug("unrecognized/non-compress-able file extension:", type);
		return next();
	}
	const path = join(res.locals.config.browserFolder, req.path.substring(1));
	const file = res.locals.foundFiles.get(path);
	if(!file || file.compressions.length === 0) {
		res.locals.logger.debug("file", path, "doesn't have any available compression");
		return next();
	}
	const acceptedEncodings = req.acceptsEncodings();
	for(const compression of file.compressions) {
		if (!acceptedEncodings.includes(compression.headerEncoding)) {
			continue;
		}
		req.url = req.url.replace(`${req.path}`, `${req.path}.${compression.fileExt}`);
		res.set("Content-Encoding", compression.headerEncoding);
		res.set("Content-Type", typeMap.get(type));
		res.locals.logger.info("Serving", compression.name, "compressed file", req.path);
		return next();
	}

	res.locals.logger.debug("no file found that matches an encoding the client accepts - serving uncompressed");
	next();
}

/**
 * A handler for proxy-ing the Traffic Ops API.
 *
 * @param req The client's request.
 * @param res The server's response writer.
 */
function toProxyHandler(req: express.Request, res: TPResponseWriter): void {
	const {logger, config} = res.locals;

	logger.debug(`Making TO API request to \`${req.originalUrl}\``);

	const fwdRequest: RequestOptions = {
		headers: req.headers,
		host: config.trafficOps.hostname,
		method: req.method,
		path: req.originalUrl,
		port: config.trafficOps.port,
		rejectUnauthorized: !config.insecure,
	};

	try {
		const proxyRequest = request(fwdRequest, r => {
			res.writeHead(r.statusCode ?? 502, r.headers);
			r.pipe(res);
		});
		req.pipe(proxyRequest);
	} catch (e) {
		logger.error("proxy-ing request:", e);
	}
	res.locals.endTime = new Date();
}

/**
 * The Express app is exported so that it can be used by serverless Functions.
 *
 * @param serverConfig Server configuration.
 * @returns The Express.js application.
 */
export async function app(serverConfig: ServerConfig): Promise<express.Express> {
	const server = express();
	const indexHtml = join(serverConfig.browserFolder, "index.html");

	// Our Universal express-engine (found @ https://github.com/angular/universal/tree/master/modules/express-engine)
	server.engine("html", ngExpressEngine({
		bootstrap: AppServerModule
	}));

	server.set("view engine", "html");
	server.set("views", "./");

	// Express 4.x doesn't handle Promise rejections (need to be manually
	// propagated with `next`), so it's not technically accurate to say that
	// void Promises are the same as void. Using `async`, though, is so much
	// easier than not doing that, so we're gonna go ahead and pretend that
	// `Promise<void>` is the same as `void`, in this one case.
	//
	// Note: Express 5.x fully supports async handlers - including seamless
	// rejections - but it's still in beta at the time of this writing.
	const loggingMW: express.RequestHandler = await loggingMiddleWare(serverConfig) as express.RequestHandler;
	server.use(loggingMW);

	// Could just use express compression `server.use(compression())` but that is calculated for each request
	server.get("*.(js|css|ttf|svg)", compressedFileHandler);

	server.get(
		"*.*",
		(req, res: TPResponseWriter, next) => {
			express.static(res.locals.config.browserFolder, {maxAge: "1y"})(req, res, next);
			// Express's static handler doesn't call `next` and calling it
			// yourself will break it for some reason, so we need to do this by
			// hand here.
			const elapsed = (new Date()).valueOf() - res.locals.startTime.valueOf();
			res.locals.logger.info("handled in", elapsed, "milliseconds with code", res.statusCode);
		}
	);

	server.use("api/**", toProxyHandler);
	server.use("/api/**", toProxyHandler);

	// All regular routes use the Universal engine
	server.get("*", (req, res: TPResponseWriter) => {
		res.render(
			indexHtml,
			{
				providers: [
					{provide: APP_BASE_HREF, useValue: req.baseUrl},
					{provide: "TP_V1_URL", useValue: res.locals.config.tpv1Url},
				],
				req
			},
		);
		res.locals.endTime = new Date();
	});

	server.use(errorMiddleWare);
	server.use((_, resp: TPResponseWriter) => {
		if (!resp.locals.endTime) {
			resp.locals.endTime = new Date();
		}
		const elapsed = resp.locals.endTime.valueOf() - resp.locals.startTime.valueOf();
		resp.locals.logger.info("handled in", elapsed, "milliseconds with code", resp.statusCode);
	});

	server.enable("trust proxy");
	return server;
}

/**
 * Runs the server.
 *
 * @returns An exit code for the process.
 */
async function run(): Promise<number> {
	const version = await getVersion();
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
	parser.add_argument("-u", "--tpv1-url", {
		dest: "tpv1Url",
		help: "Specify the Traffic Portal v1 URL. (Default: uses the `TP_V1_URL` environment variable)",
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

	let config: ServerConfig;
	try {
		config = await getConfig(parser.parse_args(), version);
	} catch (e) {
		// Logger cannot be initialized before reading server configuration
		// eslint-disable-next-line no-console
		console.error(`Failed to initialize server configuration: ${e}`);
		return 1;
	}

	const logger = new Logger(console, environment.production ? LogLevel.INFO : LogLevel.DEBUG);

	// Start up the Node server
	const server = await app(config);

	if (config.useSSL) {
		let cert: string;
		let key: string;
		let ca: Array<string>;
		try {
			cert = readFileSync(config.certPath, {encoding: "utf8"});
			key = readFileSync(config.keyPath, {encoding: "utf8"});
			ca = config.certificateAuthPaths.map(c => readFileSync(c, {encoding: "utf8"}));
		} catch (e) {
			logger.error("reading SSL key/cert:", String(e));
			if (!environment.production) {
				console.trace(e);
			}
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
		).listen(config.port, () => {
			logger.debug(`Node Express server listening on port ${config.port}`);
		});
		try {
			const redirectServer = createRedirectServer(
				(req, res) => {
					if (!req.url) {
						res.statusCode = 500;
						logger.error("got HTTP request for redirect that had no URL");
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
				logger.error("redirect server encountered error:", String(e));
				if (hasProperty(e, "code", "string") && e.code === "EACCES") {
					logger.warn("access to port 80 not allowed; closing redirect server");
					redirectServer.close();
				}
			});
		} catch (e) {
			logger.warn("Failed to initialize HTTP-to-HTTPS redirect listener:", e);
		}
	} else {
		server.listen(config.port, () => {
			logger.debug(`Node Express server listening on port ${config.port}`);
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
		run().then(
			code => {
				if (code) {
					process.exit(code);
				}
			}
		);
	}
} catch (e) {
	// Logger cannot be initialized before reading server configuration
	// eslint-disable-next-line no-console
	console.error("Encountered error while running server:", e);
	process.exit(1);
}
export * from "./src/main.server";
