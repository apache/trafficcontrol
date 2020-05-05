const child_process = require("child_process");

const GOPATH = "/go";
const srcDir = `${GOPATH}/src/github.com/apache/trafficcontrol`,
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

for (const component of components) {
	const proc = child_process.spawnSync(
		"docker",
		dockerArgs.concat([
			`trafficcontrol/${component}_builder`,
			`${srcDir}/build/build.sh`,
			component
		]),
		spawnArgs
	);

	if (proc.status !== 0) {
		console.error(`Build for ${component} failed; exiting.`);
		process.exit(proc.status);
	}
}

process.exit(0);
