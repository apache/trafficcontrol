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
const fs = require("fs");
const path = require("path");
const spawnOptions = {
	stdio: "inherit",
	stderr: "inherit"
};
const dockerCompose = ["docker", "compose", "-f", "docker-compose.yml", "-f", "docker-compose.readiness.yml"];
process.env.DOCKER_BUILDKIT = 1;
process.env.COMPOSE_DOCKER_CLI_BUILD = 1;

function moveRPMs() {
	process.chdir(`${process.env.GITHUB_WORKSPACE}/dist`);
	fs.readdirSync(".") // read contents of the dist directory
		.filter(item => fs.lstatSync(item).isDirectory()) // get a list of directories within dist
		.flatMap(directory => fs.readdirSync(directory).map(item => path.join(directory, item))) // list files within those directories
		.filter(item => item.endsWith(".rpm")) // get a list of RPMs
		.forEach(rpm => fs.renameSync(rpm, path.basename(rpm))); // move the RPMs to the dist folder
}

function runProcess(...commandArguments) {
	console.info(...commandArguments);
	const proc = child_process.spawnSync(
		commandArguments[0],
		commandArguments.slice(1),
		spawnOptions
	);
	if (proc.status === 0) {
		return;
	}
	console.error("Child process", ...commandArguments, "exited with status code", proc.status, "!");
	process.exit(proc.status);
}

moveRPMs();
process.chdir(`${process.env.GITHUB_WORKSPACE}/infrastructure/cdn-in-a-box`);
runProcess("make"); // Place the RPMs for docker compose build. All RPMs should have already been built.
runProcess(...dockerCompose, "build", "--parallel");
