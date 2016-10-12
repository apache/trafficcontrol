#!/bin/bash
#
# docker-build.sh creates containers for building each component of traffic_control with
# all needed dependencies.  Once the build is complete, all rpms are copied into the "dist"
# directory in the current directory.
#
# Usage: docker-build.sh [<options>]
# Options:
#    -r <gitrepo> git repository to clone from (defaults to value of GITREPO env variable or
#		  `https://github.com/Comcast/traffic_control').  Can be a URI or local directory.
#    -b <branch>  branch (or tag) in repository to checkout (defaults to value of BRANCH env variable or `master')
#    -d <dir>     directory to copy build artifacts (default is ./dist)

export GITREPO="${GITREPO:-https://github.com/Comcast/traffic_control}"
export BRANCH="${BRANCH:-master}"

Usage() {
	echo "Usage:"
	echo "	$0 [<option>...] [<project name>...]"
	echo "	One of -a or list of projects must be provided."
	echo "	Options:"
	echo "		-a 			build all subprojects"
	echo "		-h			show usage"
	echo "		-r <repository path>:	repository (local directory or https) to clone from"
	echo "		-b <branch name>:	branch within repository"
	echo ""
}

while getopts :hacr:b:d: opt
do
	case $opt in
		h)	Usage
			exit 1;
			;;
		a)	buildall=1
			;;
		r)
			GITREPO="$OPTARG"
			;;
		b)
			BRANCH="$OPTARG"
			;;
		*) 
			echo "Invalid option: $opt"
			Usage
			exit 1;
			;;
	esac
done
shift $((OPTIND-1))

# anything remaining is list of projects to build
if [[ -n $buildall ]]
then
	projects="traffic_ops traffic_monitor traffic_router traffic_stats traffic_portal"
else
	projects="$@"
fi

if [[ -z $projects ]]
then
	echo "One of -a or list of project names must be provided"
	Usage
	exit 1
fi


# if repo is local directory, get absolute path
if [[ -d $GITREPO ]]
then
	GITREPO=$(cd $GITREPO && pwd)
fi

DIR="$( cd "$(dirname $( dirname "${BASH_SOURCE[0]}" ))" && pwd )"

cd $DIR/infrastructure/docker/build
dist=$(pwd)/artifacts

cat <<-ENDMSG
	********************************************************
	
	Building from git repository '$GITREPO' branch '$BRANCH'
	Artifacts will be delivered to '$dist'

	Projects to build: $projects
	********************************************************

ENDMSG

# GITREPO and BRANCH are exported, so this will pick them up..
docker-compose up

echo "rpms created in $dist: "
ls -l "$dist/."
