// Type definitions for jasmine-reporters
// Project: https://github.com/larrymyers/jasmine-reporters

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

interface AppVeyorOptions {
	batchSize?: number;
	color?: boolean;
	verbosity?: number;
}

interface JUnitXmlOptions {
	savePath?: string;
	consolidate?: boolean;
	consolidateAll?: boolean;
	useDotNotation?: boolean;
	useFullTestName?: boolean;
	filePrefix?: string;
	package?: string;
	stylesheetPath?: string;
	modifySuiteName?: CallableFunction;
	systemOut?: CallableFunction;
	captureStdout?: boolean;
}

interface FormattedSuiteData {
	disabled: number;
	failures: number;
	tests: number;
	time: number;
}

interface NUnitXmlOptions {
	savePath?: string;
	filename?: string;
	reportName?: string;
}

interface TeamCityOptions {
	modifySuiteName?: CallableFunction;
}

interface TerminalOptions {
	color?: boolean;
	showStack?: boolean;
	verbosity?: number;
}

declare module "jasmine-reporters" {
	class AppVeyorReporter {
		public api: {
			endpoint?: "/api/tests/batch";
			host?: string;
			port?: string;
		};
		public batchSize: number;
		public color: boolean;
		public options: AppVeyorOptions;
		public unreportedSpecs: Array<unknown>;
		public verbosity: number;
		constructor(options?: AppVeyorOptions);
		public jasmineDone(details?: jasmine.RunDetails): void;
		public specDone(spec: object): void;
		public specStarted(spec: object): void;
	}

	class JUnitXmlReporter {
		public started: boolean;
		public finished: boolean;
		public savePath: string;
		public consolidate: boolean;
		public consolidateAll: boolean;
		public useDotNotation: boolean;
		public useFullTestName: boolean;
		public filePrefix: string;
		public package: string | undefined;
		public stylesheetPath: string | undefined;
		public captureStdout: boolean;
		public logEntries: Array<string>;
		constructor(options?: JUnitXmlOptions);
		public removeStdoutWrapper: (callback: CallableFunction) => (string?: string) => void;
		public jasmineStarted(arg?: jasmine.SuiteInfo): void;
		public suiteStarted(suite: jasmine.CustomReporterResult): void;
		public specStarted(spec: jasmine.CustomReporterResult): void;
		public specDOne(spec: jasmine.Spec): void;
		public suiteDone(suite: jasmine.CustomReporterResult): void;
		public jasmineDone(details?: jasmine.RunDetails): void;
		public formatSuiteData(suite: jasmine.Suite): FormattedSuiteData;
		public getNestedSuiteData(suite: jasmine.Suite): FormattedSuiteData;
		public getOrWriteNestedOutput(suite: jasmine.Suite): string;
		public writeFile(filename: string, text: string): void;
	}

	class NUnitXmlReporter {
		public started: boolean;
		public finished: boolean;
		public savePath: string;
		public filename: string;
		public reportName: string;
		constructor(opts?: NUnitXmlOptions);
		public jasmineStarted(summary?: jasmine.SuiteInfo): void;
		public suiteStarted(suite: jasmine.CustomReporterResult): void;
		public specStarted(spec: jasmine.CustomReporterResult): void;
		public specDone(spec: jasmine.Spec): void;
		public suiteDone(suite: jasmine.CustomReporterResult): void;
		public jasmineDone(details?: jasmine.RunDetails): void;
		public writeFile(text: string): void;
	}

	class TapReporter {
		started: boolean;
		finished: boolean;
		constructor();
		public jasmineStarted(summary: jasmine.SuiteInfo): void;
		public suiteStarted(suite: jasmine.CustomReporterResult): void;
		public specStarted(): void;
		public specDone(spec: jasmine.Spec): void;
		public suiteDone(suite: jasmine.CustomReporterResult): void;
		public jasmineDone(details?: jasmine.RunDetails): void;
	}

	class TeamCityReporter {
		public started: boolean;
		public finished: boolean;
		constructor(opts?: TeamCityOptions)
		public jasmineStarted(): void;
		public suiteStarted(suite: jasmine.CustomReporterResult): void;
		public specStarted(spec: jasmine.Spec): void;
		public specDone(spec: jasmine.Spec): void;
		public suiteDone(suite: jasmine.CustomReporterResult): void;
		public jasmineDone(): void;
	}

	class TerminalReporter {
		public started: boolean;
		public finished: boolean;
		public color: boolean | undefined;
		public showStack: boolean | undefined;
		constructor(opts?: TerminalOptions);
		public jasmineStarted(suiteInfo: jasmine.SuiteInfo): void;
		public suiteStarted(suite: jasmine.CustomReporterResult): void;
		public specStarted(spec: jasmine.Spec): void;
		public specDone(spec: jasmine.Spec): void;
		public suiteDone(suite: jasmine.CustomReporterResult): void;
		public jasmineDone(): void;
	}
}
