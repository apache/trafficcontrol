const GOPATH = '/go',
    srcDir = `${GOPATH}/src/github.com/apache/trafficcontrol`,
    component = 'traffic_ops';
process.exit(
    require('child_process')
        .spawnSync('docker', ['run',
                '-e', `GOPATH=${GOPATH}`,
                '-v', `${process.env.GITHUB_WORKSPACE}:${srcDir}`,
                `trafficcontrol/${component}_builder`,
                `${srcDir}/build/build.sh`, component],
            {stdio: 'inherit'})
        .status
);
