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

const child_process = require("child_process");

const GOPATH = "/go";
const srcDir = `${GOPATH}/src/github.com/apache/trafficcontrol`;
const components = [
	"traffic_monitor",
	"traffic_ops",
	"traffic_router",
	"traffic_stats"
];

const dockerArgs = [
	"run",
	"-e",
	`GOPATH=${GOPATH}`,
	"-v",
	`${process.env.GITHUB_WORKSPACE}:${srcDir}`
];

const spawnArgs = {stdio: "inherit"};

const timers = new Map();

for (const component of components) {
	const then = new Date();
	const proc = child_process.spawnSync(
		"docker",
		dockerArgs.concat([
			`trafficcontrol/${component}_builder`,
			`${srcDir}/build/build.sh`,
			component
		]),
		spawnArgs
	);
	const diff = new Date() - then;

	console.log(`Finished build for ${component} in ${diff}ms total`);

	if (proc.status !== 0) {
		console.error(`Build for ${component} failed; exiting.`);
		process.exit(proc.status);
	}

	timers.set(component, diff);
}

console.log(`Finished building RPMS`)
for (const [component, time] of timers) {
	console.log(`Built ${component} in ${time}ms`);
}

process.exit(0);
