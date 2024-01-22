<!--
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Running
From the `traffic-portal` directory, all you have to do is run `ng e2e`.
`npm run e2e` is functionally equivalent to `ng e2e`. This will run the tests in
interactive mode, allowing you to specify the test you want to run as well as
the browser in which to run it. In this mode, the tests will recompile and
re-run as they are edited in the source file. To run all of the tests once and
then exit - headlessly - use `ng run traffic-portal:cypress-run` or equivalently
`npm run e2e:ci`.

## Configuration
This suite assumes that you have a working Traffic Ops instance. The file
`cypress/fixtures/to.config.json` defines how the tests will connect to this
instance, including  the admin user/password needed for testing as well as the
TO API URL. The admin user defined in this file will be used for data setup, but
the user as whom the tests will run is given in `cypress/fixtures/login.json`.

## Testing Data
On each run, some mock data is set up (and **not** cleaned up afterwards) by the
testing suite and written to a fixtures data file at
`cypress/fixtures/test.data.json`. Modifying this file after the tests have
started running will not affect them in any way.
