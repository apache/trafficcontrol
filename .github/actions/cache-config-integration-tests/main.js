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

"use strict";
const child_process = require("child_process");
const spawnOptions = {
	stdio: "inherit",
	stderr: "inherit",
};

const dockerCompose = ["docker", "compose", "-f", `${process.env.GITHUB_WORKSPACE}/cache-config/testing/docker/docker-compose.yml`];

function runProcess(...commandArguments) {
	console.info(...commandArguments);
	let status = 1;
	try {
		status = child_process.spawnSync(commandArguments[0], commandArguments.slice(1), spawnOptions).status;
		if (status !== 0) {
			console.error('Child process "', ...commandArguments, '" exited with status code', status, "!");
		}
	} catch (e) {
		const msg = e instanceof Error ? e.message : e;
		console.error('Failed to spawn child process "', ...commandArguments, '":', msg);
	}
	process.exit(status);
}

runProcess(...dockerCompose, "run", "ort_test");
