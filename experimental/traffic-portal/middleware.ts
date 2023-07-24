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

import type { NextFunction, Request, Response } from "express";

import { LogLevel, Logger } from "src/app/utils";
import { environment } from "src/environments/environment";

import type { ServerConfig } from "./server.config";

/**
 * TPResponseLocals are the express.Response.locals properties specific to a
 * response writer for the TP server.
 */
interface TPResponseLocals {
	config: ServerConfig;
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
export function loggingMiddleWare(config: ServerConfig): TPHandler {
	return async (req: Request, resp: TPResponseWriter, next: NextFunction): Promise<void> => {
		resp.locals.config = config;
		const prefix = `${req.ip} HTTP/${req.httpVersion} ${req.method} ${req.url} ${req.hostname}`;
		resp.locals.logger = new Logger(console, environment.production ? LogLevel.INFO : LogLevel.DEBUG, prefix);
		resp.locals.startTime = new Date();
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
		resp.status(502); // "Bad Gateway"
		resp.write('{"alerts":[{"level":"error","text":"Unknown Traffic Portal server error occurred"}]}');
		resp.end("\n");
		resp.locals.endTime = new Date();
		next(err);
	}

	if (!resp.locals.endTime) {
		resp.locals.endTime = new Date();
	}
	const elapsed = resp.locals.endTime.valueOf() - resp.locals.startTime.valueOf();
	resp.locals.logger.info("handled in", elapsed, "milliseconds with code", resp.statusCode);
}
