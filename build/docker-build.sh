#!/bin/bash

export GITREPO=${GITREPO:-https://github.com/Comcast/traffic_control}
export BRANCH=${BRANCH:-master}

projects="traffic_ops traffic_monitor traffic_router traffic_stats"

# collect image names for later cleanup
local images=
createBuilders() {
	
	docker build -t traffic_control_gitter ./build
	images=traffic_control_gitter
	for p in $projects
	do
		docker build -t $p/build $p/build
		images="$images $p/build"
	done
}

runBuild() {
	docker run --name gitter -e GITREPO=$GITREPO -e BRANCH=$BRANCH traffic_control_gitter
	for p in $projects
	do
		docker run --rm --volumes-from gitter $p/build
	done
	docker cp gitter:/repo/traffic_control/dist .
	docker rm gitter
}

createBuilders
runBuild


# clean up...
docker rm -v gitter
docker rmi $images

