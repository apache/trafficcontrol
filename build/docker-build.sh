#!/bin/bash

# TODO: add repo and branch as cmd line options
export GITREPO=${GITREPO:-https://github.com/Comcast/traffic_control}
export BRANCH=${BRANCH:-master}

# TODO: add cmd line option to clean up images
projects="traffic_ops traffic_monitor traffic_router traffic_stats traffic_portal"

# collect image names for later cleanup
images=
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
docker cp gitter:/repo/traffic_control/dist .
docker rm -v gitter
# TODO: remove images only if requested by cmd line option
#   docker rmi $images
echo "These images were created: $images"
