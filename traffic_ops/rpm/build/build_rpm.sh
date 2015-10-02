#!/bin/bash

#
# Copyright 2015 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#----------------------------------------
function buildRpm () {
	echo "Building the rpm."

	cd $RPMBUILD && \
		rpmbuild --define "_topdir $(pwd)" \
			 --define "traffic_ops_version $TM_VERSION" \
			 --define "build_number $BUILD_NUMBER" -ba SPECS/traffic_ops.spec

	if [ $? != 0 ]; then
		echo -e "\nRPM BUILD FAILED.\n\n"
		exit 1
	else
	echo
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $DIST/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo
	if [ $BRANCH != "master" ]; then
		/usr/bin/git checkout master
	fi

	if [ ! -d $DIST ]; then
		mkdir $DIST
	fi

	/bin/cp $RPMBUILD/RPMS/*/*.rpm $DIST/.
	fi
}

#----------------------------------------
function copyToReleases() {
	cp -v $DIST/$RPM $RELEASES
	echo "The new release should be here: <a href='http://tm-ci.cdnlab.comcast.net:8888'>Release URL</a>"
}

#----------------------------------------
function checkEnvironment() {
	echo "Verifying the build configuration environment."
	# 
	# Verify the build project configuration.
	# The Jenkins configuration for this project should have the 
	# BRANCH and HOTFIX_BRANCH variables in the build parameters section.
	#
	# get traffic_control src path -- relative to build_rpm.sh script
	local srcroot=$(readlink -f $(dirname $0)/../../..)

	# set reasonable defaults for things not explicitly set in environment
	WORKSPACE=${WORKSPACE:-$srcroot}
	BRANCH=${BRANCH:-master}
	HOTFIX_BRANCH=${HOTFIX_BRANCH:-hotfix}
	BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}

	GITREPO=$WORKSPACE/traffic_ops   # WORKSPACE is the local GIT repository.
	DIST="$WORKSPACE/dist"
	PACKAGE="traffic_ops"
	RELEASES="/var/www/releases"
	RPMBUILD="$WORKSPACE/rpmbuild"
	TRAFFIC_OPS_USER="trafops"

	# set the TM_VERSION environment variable.
	TM_VERSION=$(/bin/cat $GITREPO/app/lib/UI/Utils.pm | /bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')
	RPM="${PACKAGE}-${TM_VERSION}-${BUILD_NUMBER}.x86_64.rpm"

	# verify required tools available in path
	for pgm in go ; do
		type $pgm 2>/dev/null || { echo "$pgm not found in PATH"; exit 1; }
	done
	echo "Build environment has been verified."
}

# ---------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	cd $WORKSPACE 
	#/bin/mv $RPMBUILD/carton /vol1/tmp
	/bin/rm -rf $RPMBUILD && mkdir $RPMBUILD

	/bin/mkdir $RPMBUILD/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT}
	/bin/cp -r $GITREPO/rpm/* $RPMBUILD/.
	# build the go scripts for database initialization and tm testing.

	cd $RPMBUILD

	/bin/cp traffic_ops.spec SPECS/. || { echo "Could not copy $RPMBUILD/traffic_ops.spec to $(pwd)/SPECS"; exit 1; }

	# tar/gzip the source
	local target=traffic_ops-$TM_VERSION
	local srcpath=$(pwd)/SOURCES
	local targetpath=$srcpath/$target
	/bin/mkdir -p $targetpath || { echo "Could not create $targetpath"; exit 1; }

	cd $GITREPO
	git ls-files etc app install doc | xargs /bin/cp --target=$targetpath/. --parents || \
		{ echo "Could not copy source files to $targetpath: $!"; exit 1; }

	# compile go executables used during postinstall
	cd $targetpath/install/go
	export GOPATH=$(pwd)
	export GOBIN=$targetpath/install/bin

	echo "Compiling go executables"
	for d in src/comcast.com/*; do
		if [ ! -d "$d" ]; then
			echo "Could not find $d"
			exit 1
		fi
		(cd $d && go get || { echo "Could not compile $d"; exit 1; } )
	done

   	cd $srcpath
	tar czvf $target.tgz $target || { echo "Could not create tar archive $target.tgz from $(pwd)/$target"; exit 1; }

	echo "The build area has been initialized."
}

# ---------------------------------------
function initLocalGitRepo() {
	# Allow skipping init
	if [ -n "$SKIP_INITGITREPO" ]; then
		echo "Skipping init of the local git repository."
		return
	fi
	echo "Initializing the local git repository."
	cd $GITREPO 
	/usr/bin/git checkout master && /usr/bin/git pull
	# checkout the specified BRANCH
	/usr/bin/git checkout $BRANCH && /usr/bin/git pull
	echo "Local repository is initialized, using branch $BRANCH"
}

# ---------------------------------------
function installRpm() {
	sudo /usr/bin/yum install -y $DIST/$RPM
	runCarton
	echo "Restarting traffic_ops."
	/usr/bin/sudo service traffic_ops start
}

# ---------------------------------------
function getBranch() {
	# Now update the build.number with the new branch
	PRIOR_BRANCH=$(grep branch.name= $BUILD_NUMBER_FILE|cut -d "=" -f 2)
	echo "Prior Branch: $PRIOR_BRANCH"
	# Keep the existing branch name from the prior release
	BRANCH=$(grep branch.name= $BUILD_NUMBER_FILE|cut -d "=" -f 2)
	echo "BRANCH: $BRANCH"
}

# ---------------------------------------
function getBuildNumber() {
	local commits=$(git rev-list HEAD | wc -l)
	local sha=$(git rev-parse --short=8 HEAD)
	echo "$commits.$sha"
}

# ---------------------------------------
function moveAndPushBranch() {
	cd $GITREPO
	# In case the branch already existed.
	/usr/bin/git branch -D $BRANCH
	echo "Creating new branch: $BRANCH"
	/usr/bin/git checkout -b $BRANCH

	# Update git with the new branch (if this is a release)
	git push -u origin $BRANCH
}

# ---------------------------------------
function runCarton() {
	echo ""
	echo ""
	echo "##################################################################"
	echo "# Running Carton"
	echo "##################################################################"

	if [ ! -f /usr/local/bin/carton ]; then
		sudo perl -MCPAN -e 'my $c = "CPAN::HandleConfig"; $c->load(doit => 1, autoconfig => 1); $c->edit(prerequisites_policy => "follow"); $c->edit(build_requires_install_policy => "yes"); $c->commit'
		sudo cpan -i MIYAGAWA/Carton-v1.0.15.tar.gz
	fi

	sudo -u $TRAFFIC_OPS_USER /bin/bash -c "cd /opt/traffic_ops/app && /usr/local/bin/carton install"
}

# ---------------------------------------
function runGooseUp() {
	echo "Executing Goose Up."
	cd $WORKSPACE
	./install/bin/goose up
}

# ---------------------------------------
function saveBranch() {
	BRANCH=$1
	# Now update the build.number with the new branch
	PRIOR_BRANCH=$(grep branch.name= $BUILD_NUMBER_FILE|cut -d "=" -f 2)
	echo "Prior Branch: $PRIOR_BRANCH"

	#The branch that was passed in from Jenkins is kept
	sed -i "s/\(branch.name*=*\).*/\1$BRANCH/" $BUILD_NUMBER_FILE
	#echo "Saved Branch: $BRANCH"
}

# ---------------------------------------
function saveAntBuildNumber() {
	BUILD_NUMBER=$1

	touch $ANT_BUILD_NUMBER_FILE
	echo "New Ant build.number: $BUILD_NUMBER in $ANT_BUILD_NUMBER_FILE"
	#The branch that was passed in from Jenkins is kept
	sed -i "s/\(build.number*=*\).*/\1$BUILD_NUMBER/" $ANT_BUILD_NUMBER_FILE
	#echo "Saved Build Number: $BUILD_NUMBER"
}

# ---------------------------------------
function tagRelease() {
	cd $GITREPO
	echo `pwd`
	TAG=$1
	echo "RELEASE TAG: $TAG"

	#set tag
	git tag -f $TAG

	#show tags
	git tag

	git push origin --tags
}

# ---------------------------------------
function downloadWebDeps() {
  sudo -u $TRAFFIC_OPS_USER /bin/bash -c "export PERL5LIB=/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5 && cd /opt/traffic_ops/install/bin && ./download_web_deps"
}


# --------------------------------------
# MAIN
# --------------------------------------

if [ -f /etc/profile ]; then
	. /etc/profile
fi

checkEnvironment


echo "=================================================="
echo "BRANCH: $BRANCH"
echo "HOTFIX_BRANCH: $HOTFIX_BRANCH"
echo "BUILD_NUMBER: $BUILD_NUMBER"
echo "RPM: $RPM"
echo "--------------------------------------------------"

# setup the local git repo.
initLocalGitRepo

# setup the build directory.
initBuildArea

# Build the required tm perl modules and copy them to the
# rpm build directory.
# cd $WORKSPACE/app
# RESULT=`/usr/bin/git rev-list --since="1 days ago" HEAD cpanfile`
# echo "RESULT: $RESULT"
# if [ "$RESULT" != "" ]; then
#runCarton
# fi
#runGooseUp

if [ "$BRANCH" != "master" ]; then
	echo "Executing RELEASE Flow"
	moveAndPushBranch $BRANCH
	tagRelease traffic_ops-release-${BRANCH}
	buildRpm 
	installRpm
	copyToReleases
elif [ "$HOTFIX_BRANCH" != "hotfix" ]; then
	echo "Executing HOTFIX Flow"
	tagRelease traffic_ops-hotfix-${HOTFIX_BRANCH}
	buildRpm 
	installRpm
	copyToReleases
else
	echo "Executing CI Flow"
	buildRpm 
	installRpm
	downloadWebDeps
	echo "New rpm created here: $DIST/$RPM"
fi

