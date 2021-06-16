// Type definitions for protractor-beautiful-reporter
// Project: https://github.com/Evilweed/protractor-beautiful-reporter

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
interface Options {
	baseDirectory: string;
	clientDefaults?: {
		showTotalDurationIn?: string;
		totalDurationFormat?: string;
	};
	jsonsSubfolder?: string;
	screenshotsSubfolder?: string;
	takeScreenShotsOnlyForFailedSpecs?: boolean;
	docTitle?: string;
}
/**
 * @todo This type definition is only complete enough to allow using the module
 * as it was being used at the time of this writing. Any changes to how the
 * reporter is being used MUST be reflected in updates to this type definition
 * as necessary!
 */
declare module "protractor-beautiful-reporter" {
	class HtmlReporter {
		constructor(opts: Options);
		public getJasmine2Reporter(): jasmine.Reporter | jasmine.CustomReporter;
	}
	export = HtmlReporter;
}
