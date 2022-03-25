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
loadServices();

//  _   _  _         _      _                     _          _
// | \ | |(_)       | |    | |                   | |        | |
// |  \| | _   __ _ | |__  | |_ __      __  __ _ | |_   ___ | |__
// | . ` || | / _` || '_ \ | __|\ \ /\ / / / _` || __| / __|| '_ \
// | |\  || || (_| || | | || |_  \ V  V / | (_| || |_ | (__ | | | |
// \_| \_/|_| \__, ||_| |_| \__|  \_/\_/   \__,_| \__| \___||_| |_|
//             __/ |
//            |___/

module.exports = {
	// An array of folders (excluding subfolders) where your tests are located;
	// if this is not specified, the test source must be passed as the second argument to the test runner.
	src_folders: ['./out-tsc/nightwatch/tests'],

	// See https://nightwatchjs.org/guide/working-with-page-objects/
	page_objects_path: './out-tsc/nightwatch/page_objects',

	// See https://nightwatchjs.org/guide/extending-nightwatch/#writing-custom-commands
	custom_commands_path: '',

	// See https://nightwatchjs.org/guide/extending-nightwatch/#writing-custom-assertions
	custom_assertions_path: '',

	// See https://nightwatchjs.org/guide/#external-globals
	globals_path: "../out-tsc/nightwatch/globals/globals.js",

	webdriver: {},

	test_settings: {
		default: {
			disable_error_log: false,
			launch_url: 'http://localhost:4200',

			screenshots: {
				enabled: true,
				path: 'nightwatch/screens',
				on_failure: true
			},

			output_folder: "nightwatch/junit",

			desiredCapabilities: {
				browserName: 'chrome'
			},

			webdriver: {
				start_process: true,
				server_path: ''
			}
		},

		safari: {
			desiredCapabilities: {
				browserName: 'safari',
				alwaysMatch: {
					acceptInsecureCerts: false
				}
			},
			webdriver: {
				start_process: true,
				server_path: ''
			}
		},

		firefox: {
			desiredCapabilities: {
				browserName: 'firefox',
				alwaysMatch: {
					acceptInsecureCerts: true,
					'moz:firefoxOptions': {
						args: [
							// '-headless',
							// '-verbose'
						]
					}
				}
			},
			webdriver: {
				start_process: true,
				server_path: '',
				cli_args: [
					// very verbose geckodriver logs
					// '-vv'
				]
			}
		},

		chrome: {
			desiredCapabilities: {
				browserName: 'chrome',
				'goog:chromeOptions': {
					// More info on Chromedriver: https://sites.google.com/a/chromium.org/chromedriver/
					//
					// w3c:false tells Chromedriver to run using the legacy JSONWire protocol (not required in Chrome 78)
					w3c: true,
					args: [
						//'--no-sandbox',
						//'--ignore-certificate-errors',
						//'--allow-insecure-localhost',
						//'--headless'
					]
				}
			},

			webdriver: {
				start_process: true,
				server_path: '',
				cli_args: [
					// --verbose
				]
			}
		},

		chrome_headless: {
			extends: "chrome",

			desiredCapabilities: {
				'goog:chromeOptions': {
					args: [
						'--headless'
					]
				}
			},
		},

		edge: {
			desiredCapabilities: {
				browserName: 'MicrosoftEdge',
				'ms:edgeOptions': {
					w3c: true,
					// More info on EdgeDriver: https://docs.microsoft.com/en-us/microsoft-edge/webdriver-chromium/capabilities-edge-options
					args: [
						//'--headless'
					]
				}
			},

			webdriver: {
				start_process: true,
				// Download msedgedriver from https://docs.microsoft.com/en-us/microsoft-edge/webdriver-chromium/
				//  and set the location below:
				server_path: '',
				cli_args: [
					// --verbose
				]
			}
		},

		//////////////////////////////////////////////////////////////////////////////////
		// Configuration for when using the Selenium service, either locally or remote,  |
		//  like Selenium Grid                                                           |
		//////////////////////////////////////////////////////////////////////////////////
		selenium_server: {
			// Selenium Server is running locally and is managed by Nightwatch
			selenium: {
				start_process: true,
				port: 4444,
				cli_args: {
					'webdriver.gecko.driver': (Services.geckodriver ? Services.geckodriver.path : ''),
					'webdriver.chrome.driver': (Services.chromedriver ? Services.chromedriver.path : '')
				}
			}
		},

		'selenium.chrome': {
			extends: 'selenium_server',
			desiredCapabilities: {
				browserName: 'chrome',
				chromeOptions: {
					w3c: true,
					args: [
						"--headless"
					]
				}
			}
		},

		'selenium.firefox': {
			extends: 'selenium_server',
			desiredCapabilities: {
				browserName: 'firefox',
				'moz:firefoxOptions': {
					args: [
						// '-headless',
						// '-verbose'
					]
				}
			}
		}
	}
};

function loadServices() {
	try {
		Services.seleniumServer = require('selenium-server');
	} catch (err) {
	}

	try {
		Services.chromedriver = require('chromedriver');
	} catch (err) {
	}

	try {
		Services.geckodriver = require('geckodriver');
	} catch (err) {
	}
}
