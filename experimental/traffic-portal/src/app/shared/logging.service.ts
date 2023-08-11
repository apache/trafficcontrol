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

import { Injectable } from "@angular/core";

import { LogLevel, Logger } from "../utils";

/**
 * LoggingService is for logging things in a consistent way across the UI.
 *
 * It's basically just a thin wrapper around a {@link Logger} so that only one
 * instance needs to exist and injection makes its setup consistent across all
 * usages.
 */
@Injectable({
	providedIn: "root"
})
export class LoggingService {

	public logger: Logger;

	constructor() {
		this.logger = new Logger(console, LogLevel.DEBUG, "", false);
	}

	/**
	 * Logs a debugging message.
	 *
	 * @param args Anything you want to log.
	 */
	public debug(...args: unknown[]): void {
		this.logger.debug(...args);
	}

	/**
	 * Logs an error message.
	 *
	 * @param args Anything you want to log.
	 */
	public error(...args: unknown[]): void {
		this.logger.error(...args);
	}

	/**
	 * Logs an informational message.
	 *
	 * @param args Anything you want to log.
	 */
	public info(...args: unknown[]): void {
		this.logger.info(...args);
	}

	/**
	 * Logs a warning message.
	 *
	 * @param args Anything you want to log.
	 */
	public warn(...args: unknown[]): void {
		this.logger.warn(...args);
	}
}
