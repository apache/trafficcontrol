<!--
	Licensed to the Apache Software Foundation (ASF) under one
	or more contributor license agreements. See the NOTICE file
	distributed with this work for additional information
	regarding copyright ownership. The ASF licenses this file
	to you under the Apache License, Version 2.0 (the
	"License"); you may not use this file except in compliance
	with the License. You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing,
	software distributed under the License is distributed on an
	"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
	KIND, either express or implied. See the License for the
	specific language governing permissions and limitations
	under the License.
-->

# Deprecated TP UI Tests
*Use /traffic_portal/integration for new tests*

The Traffic Portal UI tests use [Protractor](https://www.protractortest.org/#/tutorial), which thus must be installed prior to their execution. To run them, follow these steps:

1. Start up Selenium Server - typically done with `webdriver-manager start`
1. Make sure Traffic Portal is running (see [the official documentation](https://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_portal/installation.html))
1. Edit [conf.js](./conf.js) if necessary to match the environment (most notably ensure the port numbers match those in ([../conf/conf.js](../conf/conf.js) and that the login credentials are correct).
1. Run the tests - typically done with `protractor conf.js`

## Errors with webdriver
Most errors with webdriver can be remedied by running:
```shellsession
$ webdriver-manager clean
$ webdriver-manager update
```
