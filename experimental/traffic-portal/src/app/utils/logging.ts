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

/**
 * LogStreams are the underlying raw event writers used by {@link Logger}s. The
 * simplest and most useful example of a LogStreams implementation is `console`.
 */
export interface LogStreams {
	debug(...args: unknown[]): void;
	info(...args: unknown[]): void;
	error(...args: unknown[]): void;
	warn(...args: unknown[]): void;
}

/**
 * A LogLevel describes the verbosity of logging. Each level is cumulative,
 * meaning that a logger set to some level will also log all of the levels above
 * it.
 */
export const enum LogLevel {
	/** Log only errors. */
	ERROR,
	/** Log warnings and errors. */
	WARN,
	/** Log informational messages, warnings, and errors. */
	INFO,
	/** Log debugging messages, informational messages, warnings, and errors. */
	DEBUG,
}

/**
 * Converts a log level to a human-readable string.
 *
 * @example
 * console.log(logLevelToString(LogLevel.DEBUG));
 * // Output:
 * // DEBUG
 *
 * @param level The level to convert.
 * @returns A string representation of `level`.
 */
export function logLevelToString(level: LogLevel): string {
	switch(level) {
		case LogLevel.DEBUG:
			return "DEBUG";
		case LogLevel.ERROR:
			return "ERROR";
		case LogLevel.INFO:
			return "INFO";
		case LogLevel.WARN:
			return "WARN";
	}
}

/**
 * A Logger logs things. The output streams are customizable, mostly for testing
 * but also in case we want to write directly to a file handle someday.
 *
 * The output format is a bit customizable, it allows for messages to be
 * prefixed in a number of ways:
 * - With the level at which the message was logged
 * - With a timestamp for the time at which logging occurred (ISO format)
 * - With some static string
 *
 * in that order. For example, if all of them are specified:
 *
 * @example
 * (new Logger(console, LogLevel.DEBUG, "test", true, true)).info("quest");
 * // Output (example date is UNIX epoch):
 * // INFO 1970-01-01T00:00:00.000Z test: quest
 */
export class Logger {
	private readonly prefix: string;

	/**
	 * Constructor.
	 *
	 * @param streams The output stream abstractions.
	 * @param level The level at which the logger operates. Any level higher
	 * than the one specified will not be logged.
	 * @param prefix If given, prepends a prefix to each message.
	 * @param useLevelPrefixes If true, log lines will be prefixed with the name
	 * of the level at which they were logged (useful if all streams point to
	 * the same file descriptor).
	 * @param timestamps If true, each log line will be accompanied by a
	 * timestamp prefix (note that the time is determined when logging occurs,
	 * not necessarily when the logging method is called).
	 */
	constructor(
		private readonly streams: LogStreams,
		level: LogLevel,
		prefix: string = "",
		private readonly useLevelPrefixes: boolean = true,
		private readonly timestamps: boolean = true,
	) {
		if (prefix) {
			prefix = prefix.trim().replace(/:$/, "").trimEnd();
		}

		this.prefix = prefix;

		const doNothing = (): void => { /* Do nothing */ };
		switch (level) {
			case LogLevel.ERROR:
				this.warn = doNothing;
			case LogLevel.WARN:
				this.info = doNothing;
			case LogLevel.INFO:
				this.debug = doNothing;
		}

		// saves time later; getPrefix will make these same checks and return
		// the same value if they all go the same way.
		if (!this.timestamps && !this.useLevelPrefixes && !this.prefix) {
			this.getPrefix = (): string => "";
		}
	}

	/**
	 * Constructs a prefix for logging at a given level based on the Logger's
	 * configuration.
	 *
	 * @param level The level at which a message is being logged.
	 * @returns A prefix, or an empty string if no prefix is to be used.
	 */
	private getPrefix(level: LogLevel): string {
		const parts = new Array<string>();

		if (this.timestamps) {
			parts.push(new Date().toISOString());
		}

		if (this.useLevelPrefixes) {
			parts.unshift(logLevelToString(level));
		}

		if (this.prefix) {
			parts.push(this.prefix);
		}

		// This colon isn't a problem, because if none of the above checks to
		// add content to `parts` passed, the constructor would have optimized
		// this whole method away.
		return `${parts.join(" ")}:`;
	}

	/**
	 * Logs a message at the DEBUG level.
	 *
	 * @param args Anything representable as text. Be careful passing objects;
	 * while technically allowed, this will probably cause multi-line log
	 * messages which are not easy to parse. Similarly, please don't use
	 * newlines.
	 */
	public debug(...args: unknown[]): void {
		const prefix = this.getPrefix(LogLevel.DEBUG);
		if (prefix) {
			this.streams.debug(prefix, ...args);
			return;
		}
		this.streams.debug(...args);
	}

	/**
	 * Logs a message at the ERROR level.
	 *
	 * @param args Anything representable as text. Be careful passing objects;
	 * while technically allowed, this will probably cause multi-line log
	 * messages which are not easy to parse. Similarly, please don't use
	 * newlines.
	 */
	public error(...args: unknown[]): void {
		const prefix = this.getPrefix(LogLevel.ERROR);
		if (prefix) {
			this.streams.error(prefix, ...args);
			return;
		}
		this.streams.error(...args);
	}

	/**
	 * Logs a message at the INFO level.
	 *
	 * @param args Anything representable as text. Be careful passing objects;
	 * while technically allowed, this will probably cause multi-line log
	 * messages which are not easy to parse. Similarly, please don't use
	 * newlines.
	 */
	public info(...args: unknown[]): void {
		const prefix = this.getPrefix(LogLevel.INFO);
		if (prefix) {
			this.streams.info(prefix, ...args);
			return;
		}
		this.streams.info(...args);
	}

	/**
	 * Logs a message at the WARN level.
	 *
	 * @param args Anything representable as text. Be careful passing objects;
	 * while technically allowed, this will probably cause multi-line log
	 * messages which are not easy to parse. Similarly, please don't use
	 * newlines.
	 */
	public warn(...args: unknown[]): void {
		const prefix = this.getPrefix(LogLevel.WARN);
		if (prefix) {
			this.streams.warn(prefix, ...args);
			return;
		}
		this.streams.warn(...args);
	}
}
