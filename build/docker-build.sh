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
dist="./dist"
cleanup=

Usage() {
	echo "Usage:"
	echo "	$0 [<option>...] [<project name>...]"
	echo "	One of -a or list of projects must be provided."
	echo "	Options:"
	echo "		-a 			build all subprojects"
	echo "		-h			show usage"
	echo "		-r <repository path>:	repository (local directory or https) to clone from"
	echo "		-b <branch name>:	branch within repository"
	echo "		-c			start clean: remove all traffic_control docker images prior to building"
	echo "		-d <dist dir>:		local directory to copy built rpms"
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
		c)
			cleanup=1
			;;
		r)
			GITREPO="$OPTARG"
			;;
		b)
			BRANCH="$OPTARG"
			;;
		d)
			dist="$OPTARG"
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

# Get absolute path to dist dir
mkdir -p $dist || exit 1
dist=$(cd $dist && pwd)

cleanmsg=$([[ $cleanup ]] && echo "be cleaned up" || echo "not be cleaned up")
cat <<-ENDMSG
	********************************************************
	
	Building from git repository '$GITREPO' branch '$BRANCH'
	Artifacts will be delivered to '$dist'

	Projects to build: $projects
	********************************************************

ENDMSG

# collect image names for later cleanup
createBuilders() {
	# topdir=.../traffic_control
	local topdir=$(cd "$( echo "${BASH_SOURCE[0]%/*}" )/.."; pwd)

	echo -n "** Create Builders: "; date
	for p in $projects
	do 
		local image=$p/build
		echo -n "**   $image: "; date
		docker build --tag $image "$topdir/$p/build"
	done
}

runBuild() {
	echo -n "** Run Build: "; date

	# Check if gitrepo is a local directory to be provided as a volume
	if [[ -d $GITREPO ]]
	then
		vol="-v $GITREPO:$GITREPO"
	fi
	mkdir -p dist
	for p in $projects
	do
		echo -n "**   building $p: "; date
		docker run --rm --env "GITREPO=$GITREPO" --env "BRANCH=$BRANCH" $vol -v $dist:/dist $p/build
	done
	echo -n "** End Build: "; date
}

createBuilders
runBuild

echo "rpms created: "
ls -l "$dist/."
