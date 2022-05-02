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

// Refer to the online docs for more details: https://nightwatchjs.org/gettingstarted/configuration/
const Services = {};

module.exports = {
	// See https://nightwatchjs.org/guide/extending-nightwatch/#writing-custom-assertions
	custom_assertions_path: "",

	// See https://nightwatchjs.org/guide/extending-nightwatch/#writing-custom-commands
	custom_commands_path: "",

	// See https://nightwatchjs.org/guide/#external-globals
	globals_path: "../out-tsc/nightwatch/globals/globals.js",

	// See https://nightwatchjs.org/guide/working-with-page-objects/
	page_objects_path: "./out-tsc/nightwatch/page_objects",
	// An array of folders (excluding subfolders) where your tests are located;
	// if this is not specified, the test source must be passed as the second argument to the test runner.
	src_folders: ["./out-tsc/nightwatch/tests"],

	test_settings: {
		chrome: {
			desiredCapabilities: {
				browserName: "chrome",
				"goog:chromeOptions": {
					args: [
						//'--no-sandbox',
						//'--ignore-certificate-errors',
						//'--allow-insecure-localhost',
						//'--headless'
					],
					// More info on Chromedriver: https://sites.google.com/a/chromium.org/chromedriver/
					//
					// w3c:false tells Chromedriver to run using the legacy JSONWire protocol (not required in Chrome 78)
					w3c: true
				}
			},

			webdriver: {
				cli_args: [
					// --verbose
				],
				server_path: "",
				start_process: true
			}
		},

		chrome_headless: {
			desiredCapabilities: {
				"goog:chromeOptions": {
					args: [
						"--headless",
						"--window-size=1920,1080"
					]
				}
			},

			extends: "chrome"
		},

		default: {
			desiredCapabilities: {
				browserName: "chrome"
			},
			disable_error_log: false,
			launch_url: "http://localhost:4200",
			output_folder: "nightwatch/junit",
			screenshots: {
				enabled: true,
				on_failure: true,
				path: "nightwatch/screens"
			},
			webdriver: {
				server_path: "",
				start_process: true
			}
		},

		edge: {
			desiredCapabilities: {
				browserName: "MicrosoftEdge",
				"ms:edgeOptions": {
					// More info on EdgeDriver: https://docs.microsoft.com/en-us/microsoft-edge/webdriver-chromium/capabilities-edge-options
					args: [
						//'--headless'
					],
					w3c: true
				}
			},

			webdriver: {
				cli_args: [
					// --verbose
				],
				// Download msedgedriver from https://docs.microsoft.com/en-us/microsoft-edge/webdriver-chromium/
				//  and set the location below:
				server_path: "",
				start_process: true
			}
		},

		firefox: {
			desiredCapabilities: {
				alwaysMatch: {
					acceptInsecureCerts: true,
					"moz:firefoxOptions": {
						args: [
							// '-headless',
							// '-verbose'
						]
					}
				},
				browserName: "firefox",
			},
			webdriver: {
				cli_args: [
					// very verbose geckodriver logs
					// '-vv'
				],
				server_path: "",
				start_process: true
			}
		},
		safari: {
			desiredCapabilities: {
				alwaysMatch: {
					acceptInsecureCerts: false
				},
				browserName: "safari"
			},
			webdriver: {
				server_path: "",
				start_process: true
			}
		},

	},

	webdriver: {}
};

/**
 * Loads browser drivers if available
 *
 */
function loadServices() {
	try {
		// eslint-disable-next-line @typescript-eslint/no-require-imports
		Services.chromedriver = require("chromedriver");
	} catch (err) {
		console.log("Unable to load chromedriver");
	}

	try {
		// eslint-disable-next-line @typescript-eslint/no-require-imports
		Services.geckodriver = require("geckodriver");
	} catch (err) {
		console.log("Unable to load geckodriver");
	}
}

loadServices();
