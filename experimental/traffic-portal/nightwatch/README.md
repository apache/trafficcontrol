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
From the `traffic-portal` directory, all npm run scripts with `e2e` can be used to run e2e tests. 
`ng e2e` also works and is functionally equivalent to `npm run e2e`.

## Requirements
This suite assumes that you have a working Traffic Ops/Traffic Portal v2 instance. The file `globals/globals.ts` defines 
the admin user/password needed for testing as well as the TO API URL.

## Browsers
It is possible to run tests using your browser of choice (see `nightwatch.conf.js`), this repo has the necessary packages
to run using Firefox and Chrome with Chrome being the default. To run under a specific environment pass
the `--env {browserName}` flag to the npm run script e.g. `npm run e2e -- --env firefox`.
