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

import { Logger, LogLevel, logLevelToString, type LogStreams } from "./logging";

/**
 * TestingStreams is a Streams implementation that pushes each log line to a
 * publicly available array per stream, allowing for easy inspection by testing
 * routines afterward.
 */
class TestingStreams implements LogStreams {
	public readonly debugStream = new Array<string>();
	public readonly errorStream = new Array<string>();
	public readonly infoStream = new Array<string>();
	public readonly warnStream = new Array<string>();

	/**
	 * Logs to the debug stream.
	 *
	 * @param args anything
	 */
	public debug(...args: unknown[]): void {
		this.debugStream.push(args.join(" "));
	}
	/**
	 * Logs to the debug stream.
	 *
	 * @param args anything
	 */
	public error(...args: unknown[]): void {
		this.errorStream.push(args.join(" "));
	}
	/**
	 * Logs to the info stream.
	 *
	 * @param args anything
	 */
	public info(...args: unknown[]): void {
		this.infoStream.push(args.join(" "));
	}
	/**
	 * Logs to the warning stream.
	 *
	 * @param args anything
	 */
	public warn(...args: unknown[]): void {
		this.warnStream.push(args.join(" "));
	}
}

const timestampPattern = "\\d{4}-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\d\\.\\d+Z";

describe("logging utility functions", () => {
	it("converts debug level to a string", () => {
		expect(logLevelToString(LogLevel.DEBUG)).toBe("DEBUG");
	});
	it("converts error level to a string", () => {
		expect(logLevelToString(LogLevel.ERROR)).toBe("ERROR");
	});
	it("converts info level to a string", () => {
		expect(logLevelToString(LogLevel.INFO)).toBe("INFO");
	});
	it("converts warn level to a string", () => {
		expect(logLevelToString(LogLevel.WARN)).toBe("WARN");
	});
});

describe("Logger", () => {
	let streams: TestingStreams;

	beforeEach(() => {
		streams = new TestingStreams();
		expect(streams.debugStream).toHaveSize(0);
		expect(streams.errorStream).toHaveSize(0);
		expect(streams.infoStream).toHaveSize(0);
		expect(streams.warnStream).toHaveSize(0);
	});

	describe("prefix-less logging", () => {
		let logger: Logger;
		const msg = "testquest";
		beforeEach(() => {
			logger = new Logger(streams, LogLevel.DEBUG, "", false, false);
		});

		it("logs to the debug stream", () => {
			logger.debug(msg);
			expect(streams.debugStream).toHaveSize(1);
			expect(streams.debugStream).toContain(msg);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the error stream", () => {
			logger.error(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(1);
			expect(streams.errorStream).toContain(msg);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the info stream", () => {
			logger.info(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(1);
			expect(streams.infoStream).toContain(msg);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the warn stream", () => {
			logger.warn(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(1);
			expect(streams.warnStream).toContain(msg);
		});
	});

	describe("static prefixed logging", () => {
		let logger: Logger;
		const prefix = "test";
		const msg = "quest";
		beforeEach(() => {
			logger = new Logger(streams, LogLevel.DEBUG, prefix, false, false);
		});

		it("logs to the debug stream", () => {
			logger.debug(msg);
			expect(streams.debugStream).toHaveSize(1);
			expect(streams.debugStream).toContain(`${prefix}: ${msg}`);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the error stream", () => {
			logger.error(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(1);
			expect(streams.errorStream).toContain(`${prefix}: ${msg}`);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the info stream", () => {
			logger.info(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(1);
			expect(streams.infoStream).toContain(`${prefix}: ${msg}`);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the warn stream", () => {
			logger.warn(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(1);
			expect(streams.warnStream).toContain(`${prefix}: ${msg}`);
		});
	});

	describe("timestamp-prefixed logging", () => {
		let logger: Logger;
		const msg = "testquest";
		beforeEach(() => {
			logger = new Logger(streams, LogLevel.DEBUG, "", false, true);
		});

		it("logs to the debug stream", () => {
			logger.debug(msg);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
			if (streams.debugStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.debugStream.length}`);
			}
			expect(streams.debugStream[0]).toMatch(`^${timestampPattern}: ${msg}$`);
		});

		it("logs to the error stream", () => {
			logger.error(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
			if (streams.errorStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.errorStream.length}`);
			}
			expect(streams.errorStream[0]).toMatch(`^${timestampPattern}: ${msg}$`);
		});

		it("logs to the info stream", () => {
			logger.info(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
			if (streams.infoStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.infoStream.length}`);
			}
			expect(streams.infoStream[0]).toMatch(`^${timestampPattern}: ${msg}$`);
		});

		it("logs to the warn stream", () => {
			logger.warn(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			if (streams.warnStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.warnStream.length}`);
			}
			expect(streams.warnStream[0]).toMatch(`^${timestampPattern}: ${msg}$`);
		});
	});

	describe("log-level-prefixed logging", () => {
		let logger: Logger;
		const msg = "testquest";
		beforeEach(() => {
			logger = new Logger(streams, LogLevel.DEBUG, "", true, false);
		});

		it("logs to the debug stream", () => {
			logger.debug(msg);
			expect(streams.debugStream).toHaveSize(1);
			expect(streams.debugStream).toContain(`${logLevelToString(LogLevel.DEBUG)}: ${msg}`);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the error stream", () => {
			logger.error(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(1);
			expect(streams.errorStream).toContain(`${logLevelToString(LogLevel.ERROR)}: ${msg}`);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the info stream", () => {
			logger.info(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(1);
			expect(streams.infoStream).toContain(`${logLevelToString(LogLevel.INFO)}: ${msg}`);
			expect(streams.warnStream).toHaveSize(0);
		});

		it("logs to the warn stream", () => {
			logger.warn(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(1);
			expect(streams.warnStream).toContain(`${logLevelToString(LogLevel.WARN)}: ${msg}`);
		});
	});

	describe("fully-prefixed logging", () => {
		let logger: Logger;
		const prefix = "test";
		const msg = "quest";
		beforeEach(() => {
			logger = new Logger(streams, LogLevel.DEBUG, prefix);
		});

		it("logs to the debug stream", () => {
			logger.debug(msg);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
			if (streams.debugStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.debugStream.length}`);
			}
			expect(streams.debugStream[0]).toMatch(`^${logLevelToString(LogLevel.DEBUG)} ${timestampPattern} ${prefix}: ${msg}$`);
		});

		it("logs to the error stream", () => {
			logger.error(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
			if (streams.errorStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.errorStream.length}`);
			}
			expect(streams.errorStream[0]).toMatch(`^${logLevelToString(LogLevel.ERROR)} ${timestampPattern} ${prefix}: ${msg}$`);
		});

		it("logs to the info stream", () => {
			logger.info(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
			if (streams.infoStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.infoStream.length}`);
			}
			expect(streams.infoStream[0]).toMatch(`^${logLevelToString(LogLevel.INFO)} ${timestampPattern} ${prefix}: ${msg}$`);
		});

		it("logs to the warn stream", () => {
			logger.warn(msg);
			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(0);
			expect(streams.infoStream).toHaveSize(0);
			if (streams.warnStream.length !== 1) {
				return fail(`incorrect stream size after logging; want: 1, got: ${streams.warnStream.length}`);
			}
			expect(streams.warnStream[0]).toMatch(`^${logLevelToString(LogLevel.WARN)} ${timestampPattern} ${prefix}: ${msg}$`);
		});
	});

	describe("log-level specification", () => {
		it("won't log above INFO if set to INFO", () => {
			const logger = new Logger(streams, LogLevel.INFO);

			logger.debug("anything");
			logger.error("anything");
			logger.info("anything");
			logger.warn("anything");

			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(1);
			expect(streams.infoStream).toHaveSize(1);
			expect(streams.warnStream).toHaveSize(1);
		});

		it("won't log above WARN if set to WARN", () => {
			const logger = new Logger(streams, LogLevel.WARN);

			logger.debug("anything");
			logger.error("anything");
			logger.info("anything");
			logger.warn("anything");

			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(1);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(1);
		});

		it("won't log above ERROR if set to ERROR", () => {
			const logger = new Logger(streams, LogLevel.ERROR);

			logger.debug("anything");
			logger.error("anything");
			logger.info("anything");
			logger.warn("anything");

			expect(streams.debugStream).toHaveSize(0);
			expect(streams.errorStream).toHaveSize(1);
			expect(streams.infoStream).toHaveSize(0);
			expect(streams.warnStream).toHaveSize(0);
		});
	});
});
