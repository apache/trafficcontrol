/**
 * @license Apache-2.0
 *
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
import { opendir } from "fs/promises";
import { join } from "path";

import type { NextFunction, Request, Response } from "express";

import { LogLevel, Logger } from "src/app/utils";
import { environment } from "src/environments/environment";

import type { ServerConfig } from "./server.config";

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

/**
 * TPResponseLocals are the express.Response.locals properties specific to a
 * response writer for the TP server.
 */
interface TPResponseLocals {
	config: ServerConfig;
	foundFiles: Map<string, StaticFile>;
	logger: Logger;
	/** The time at which the request was received. */
	startTime: Date;
	/**
	 * The time at which the response was finished being written (or
	 * `undefined` if not done yet).
	 */
	endTime?: Date | undefined;
}

/**
 * AuthenticatedResponse is a response writer for endpoints that require
 * authentication.
 */
export type TPResponseWriter = Response<unknown, TPResponseLocals>;

/**
 * An HTTP request handler for the TP server.
 */
export type TPHandler = (req: Request, resp: TPResponseWriter, next: NextFunction) => void | PromiseLike<void>;

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
async function getFiles(path: string): Promise<string[]> {
	const dir = await opendir(path);
	let dirEnt = await dir.read();
	let files = new Array<string>();

	while (dirEnt !== null) {
		const name = join(path, dirEnt.name);

		if (dirEnt.isDirectory()) {
			files = files.concat(await getFiles(name));
		} else {
			files.push(name);
		}

		dirEnt = await dir.read();
	}
	await dir.close();

	return files;
}

/**
 * loggingMiddleWare is a middleware factory for express.js that provides a
 * logger.
 * It does also provide a link to server configuration that can be used in
 * handlers, and a couple other niceties.
 *
 * @param config The server configuration.
 * @returns A middleware that adds a property `logger` to `resp.locals` for
 * logging purposes.
 */
export async function loggingMiddleWare(config: ServerConfig): Promise<TPHandler> {
	const allFiles = await getFiles(config.browserFolder);
	const compressedFiles = new Map(
		allFiles.filter(
			file => file.match(/\.(br|gz)$/)
		).map(
			file => [file, undefined]
		)
	);
	const foundFiles = new Map<string, StaticFile>(
		allFiles.filter(
			file => file.match(/\.(js|css|tff|svg)$/)
		).map(
			file => {
				const staticFile: StaticFile = {
					compressions: []
				};
				if (compressedFiles.has(`${file}.${br.fileExt}`)) {
					staticFile.compressions.push(br);
				}
				if (compressedFiles.has(`${file}.${gzip.fileExt}`)) {
					staticFile.compressions.push(gzip);
				}
				return [file, staticFile];
			}
		)
	);

	return async (req: Request, resp: TPResponseWriter, next: NextFunction): Promise<void> => {
		resp.locals.config = config;
		const prefix = `${req.ip} HTTP/${req.httpVersion} ${req.method} ${req.url} ${req.hostname}`;
		resp.locals.logger = new Logger(console, environment.production ? LogLevel.INFO : LogLevel.DEBUG, prefix);
		resp.locals.logger.debug("handling");
		resp.locals.startTime = new Date();
		resp.locals.foundFiles = foundFiles;

		next();
	};
}

/**
 * errorMiddleWare is a middleware for express.js that provides automatic
 * handling of errors that aren't caught in the endpoint handlers.
 *
 * @param err Any error passed along by other handlers.
 * @param _ The client request - unused.
 * @param resp The server's response-writer
 * @param next A function provided by Express which will call the next handler.
 */
export function errorMiddleWare(err: unknown, _: Request, resp: TPResponseWriter, next: NextFunction): void {
	if (err !== null && err !== undefined) {
		resp.locals.logger.error("unhandled error bubbled to routing:", String(err));
		if (!environment.production) {
			console.trace(err);
		}
		if (!resp.locals.endTime) {
			resp.status(502); // "Bad Gateway"
			resp.write('{"alerts":[{"level":"error","text":"Unknown Traffic Portal server error occurred"}]}\n');
			resp.end("\n");
			resp.locals.endTime = new Date();
			next(err);
		}
	}
}
