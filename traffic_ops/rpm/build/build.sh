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

if [ -z $WORKSPACE ]; then
	echo "Error: the 'WORKSPACE' environment variable is not set."
	exit 1
fi

BRANCH="master"
GOPATH="/var/lib/jenkins/go"; export GOPATH
BUILDDIR=$WORKSPACE/build
TMSRC="$WORKSPACE/tm"
CARTON="$WORKSPACE/carton"
UTILS_PM="$TMSRC/app/lib/UI/Utils.pm"
PACKAGE="traffic_ops"

if [ -f /etc/profile ]; then
    . /etc/profile
fi

if [ -z $1 ]; then
    echo "The BRANCH variable is not set."
    exit 1
else
    BRANCH=$1
fi

if [ -z $2 ]; then
	echo "The BUILD_NUMBER variable is not set."
	exit 2
else
	BUILD_NUMBER=$2
fi

if [ -z $3 ]; then
	echo "The GIT variable is not set."
	exit 3
else
	GIT=$3
fi

if [ -d $TMSRC ]; then
    cd $TMSRC && /usr/bin/git pull
    /usr/bin/git checkout $BRANCH
    /usr/bin/git pull
else
    if [ ! -d $WORKSPACE/tm ]; then
	/bin/mkdir $WORKSPACE/tm
    fi

    cd $WORKSPACE && /usr/bin/git clone $GIT
    /usr/bin/git checkout $BRANCH
fi  

VERSION=$(/bin/cat $UTILS_PM|/bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')
RPM="${PACKAGE}-${VERSION}-${BUILD_NUMBER}.x86_64.rpm"

echo
echo "=================================================================="
echo "Building Traffic Ops rpm traffic_ops-${VERSION}-$BUILD_NUMBER"
echo
echo "GOPATH=$GOPATH"
echo "BUILDDIR=$WORKSPACE/build"
echo "TMSRC=$WORKSPACE/traffic_ops"
echo "CARTON=$WORKSPACE/carton"
echo "UTILS_PM=$TMSRC/app/lib/UI/Utils.pm"
echo "=================================================================="
echo

cd $BUILDDIR

if [ ! -d dist ]; then
    /bin/mkdir dist
fi

if [ -d rpm ]; then
    /bin/rm -rf rpm
fi

/bin/cp -R $TMSRC/rpm .
echo "build.number=$BUILD_NUMBER" > rpm/build.number

cd $TMSRC/install/bin
/usr/local/go/bin/go build $TMSRC/install/go/src/comcast.com/dataload/dataload.go
/usr/local/go/bin/go build $TMSRC/install/go/src/comcast.com/systemtest/systemtest.go

cd $BUILDDIR

for link in etc app install db doc; do
    if [ ! -s $link ]; then
	ln -s $TMSRC/$link $link
    fi
done

if [ -d lib ]; then
    /bin/rm -rf lib
fi

/bin/mkdir -p lib/perl5

if [ -d bin ]; then
    /bin/rm -rf bin
fi

/bin/mkdir bin

if [ ! -d $CARTON ]; then
    /bin/mkdir $CARTON
fi

cd $CARTON

/bin/cp $TMSRC/app/cpanfile .

/usr/local/bin/carton install

cd $BUILDDIR

/bin/cp -R $CARTON/local/bin/* bin
/bin/cp -R $CARTON/local/lib/perl5/* lib/perl5

#
# Ant builds the rpm, perl modules should have been built
# by carton already and placed in the lib/perl5 directory.
echo -e "arch=x86_64\ntm_version=$VERSION" > rpm/traffic_ops.properties
cd rpm && /usr/local/ant/bin/ant

if [ $? != 0 ]; then
    echo -e "\nRPM BUILD FAILED.\n\n"
else
    echo
    echo "========================================================================================"
    echo "RPM BUILD SUCCEEDED, See $BUILDDIR/dist/$RPM for the newly built rpm."
    echo "========================================================================================"
    echo
    if [ $BRANCH != "master" ]; then
	    /usr/bin/git checkout master
    fi
    /bin/cp $BUILDDIR/rpm/dist/*.rpm $BUILDDIR/dist
fi

