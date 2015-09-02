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
    echo -e "arch=x86_64\ntm_version=$TM_VERSION" > $RPMBUILD/traffic_ops.properties
    cd $RPMBUILD && ant

    if [ $? != 0 ]; then
        echo -e "\nRPM BUILD FAILED.\n\n"
        exit $?
    else
	echo
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $BUILDDIR/dist/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo
	if [ $BRANCH != "master" ]; then
	    /usr/bin/git checkout master
	fi

	if [ ! -d $DIST ]; then
	    mkdir $DIST
	fi

	/bin/cp $RPMBUILD/dist/*.rpm $DIST
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
	WORKSPACE=${WORKSPACE:-$HOME/workspace}
	BRANCH=${BRANCH:-master}
	HOTFIX_BRANCH=${HOTFIX_BRANCH:-hotfix}
	BUILD_NUMBER=${BUILD_NUMBER:-0}

    # set the TM_VERSION environment variable.
    TM_VERSION=$(/bin/cat $GITREPO/app/lib/UI/Utils.pm | /bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')
    RPM="${PACKAGE}-${TM_VERSION}-${BUILD_NUMBER}.x86_64.rpm"

    echo "Build environment has been verified."
}

# ---------------------------------------
function initBuildArea() {
    echo "Initializing the build area."
    cd $JOB_DIRECTORY 
    #/bin/mv $RPMBUILD/carton /vol1/tmp
    /bin/rm -rf $RPMBUILD && mkdir $RPMBUILD
    #/bin/mv /vol1/tmp/carton $RPMBUILD


    /bin/cp -R $GITREPO/rpm/* $RPMBUILD
    /bin/cp -R $GITREPO/install $RPMBUILD
    /bin/cp -R $GITREPO/doc $RPMBUILD
    # build the go scripts for database initialization and tm testing.
    export GOBIN=$RPMBUILD/install/bin
    echo "Compiling go"
    for d in $RPMBUILD/install/go/src/comcast.com/*; do
	(cd $d; go get; go install)
    done
    
    cd $RPMBUILD
    # write the build.number file required by ant
    echo "build.number=$BUILD_NUMBER" > build.number

    # setup the links to the source files in the GITREPO
    for d in etc app; do
	mkdir $d
	/bin/cp -R $GITREPO/$d/* $d
    done

    echo "The build area has been initialized."
}

# ---------------------------------------
function initLocalGitRepo() {
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
    # Keep the existing branch name from the prior build
	# TODO: base the build number on something from the repo
	if [ -f "$BUILD_NUMBER_FILE" ]; then
	    BUILD_NUMBER=$(grep build.number= $BUILD_NUMBER_FILE|cut -d "=" -f 2)
	else
	    BUILD_NUMBER=0
	fi
	
    echo "BUILD_NUMBER: $BUILD_NUMBER"
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
    cd $JOB_DIRECTORY
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

# Environment Constants
GITREPO=$WORKSPACE/traffic_ops   # WORKSPACE is the local GIT repository.
JOB_DIRECTORY=$WORKSPACE
DIST="$JOB_DIRECTORY/dist"
PACKAGE="traffic_ops"
RELEASES="/var/www/releases"
RPMBUILD="$JOB_DIRECTORY/rpmbuild"
TRAFFIC_OPS_USER="trafops"


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
fi

