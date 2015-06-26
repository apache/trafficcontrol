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
   echo "If running from a Vagrant VM set WORKSPACE should have been set to "
   echo " /home/vagrant/rpmbuild."
   echo "example: export WORKSPACE=/home/vagrant/"
	exit 1
fi

if [ ! -d $WORKSPACE ]; then
   echo "$WORKSPACE does not exist."
   exit 1
fi

#if [ -z $TCSRC ]; then
#	echo "Error: the 'TCSRC' environment variable is not set."
#   echo "If running from a Vagrant VM TCSRC environment should be /home/vagrant/traffic_control."
#   echo "example: export TCSRC=/vagrant"
#	exit 1
#fi

sudo perl -MCPAN -e 'my $c = "CPAN::HandleConfig"; $c->load(doit => 1, autoconfig => 1); $c->edit(prerequisites_policy => "follow"); $c->edit(build_requires_install_policy => "yes"); $c->commit'
sudo cpan -if MIYAGAWA/Carton-v1.0.15.tar.gz

go get github.com/go-sql-driver/mysql
go get code.google.com/p/go.net/html
go get code.google.com/p/go.net/publicsuffix

PACKAGE="traffic_ops"
BUILDDIR=$WORKSPACE/build
TCSRC="$WORKSPACE/traffic_control"
TOSRC="$TCSRC/$PACKAGE"
CARTON="$WORKSPACE/carton"
UTILS_PM="$TCSRC/$PACKAGE/app/lib/UI/Utils.pm"

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

#if [ -z $3 ]; then
#	echo "The GIT variable is not set."
#	exit 3
#else
#	GIT=$3
#fi

echo "package: $PACKAGE"
echo "tcsrc: $TCSRC"
echo "tosrc: $TOSRC"

if [ -d $TOSRC ]; then
   echo "cd to $TCSRC and then git pull"
   cd $TCSRC && /usr/bin/git pull
   echo "git checkout $BRANCH"
   /usr/bin/git checkout $BRANCH
   /usr/bin/git pull
else
    #if [ ! -d $WORKSPACE/traffic_ops ]; then
	 #   /bin/mkdir $WORKSPACE/traffic_ops
    #fi
   echo "cd to $WORKSPACE and clone"
    cd $WORKSPACE && /usr/bin/git clone https://github.com/Comcast/traffic_control.git
    cd traffic_ops
    /usr/bin/git checkout $BRANCH
fi  

VERSION=$(/bin/cat $UTILS_PM|/bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')
RPM="${PACKAGE}-${VERSION}-${BUILD_NUMBER}.x86_64.rpm"

echo
echo "=================================================================="
echo "Building Traffic Ops rpm traffic_ops-${VERSION}-$BUILD_NUMBER"
echo
#echo "GOPATH=$GOPATH"
echo "BUILDDIR=$BUILDDIR"
echo "TOSRC=$TOSRC"
echo "CARTON=$CARTON"
echo "UTILS_PM=$TOSRC/app/lib/UI/Utils.pm"
echo "=================================================================="
echo


if [ ! -d $BUILDDIR ]; then
   #echo "The build dir $BUILDDIR does not exist. Please create it."
   #exit 1
   mkdir $BUILDDIR
   cd $BUILDDIR
else
   cd $BUILDDIR
fi

if [ ! -d dist ]; then
    /bin/mkdir dist
fi

if [ -d rpm ]; then
    /bin/rm -rf rpm
fi

/bin/cp -R $TOSRC/rpm .
echo "build.number=$BUILD_NUMBER" > rpm/build.number

cd $TOSRC/install/bin
go build $TOSRC/install/go/src/comcast.com/dataload/dataload.go
go build $TOSRC/install/go/src/comcast.com/systemtest/systemtest.go

cd $BUILDDIR/rpm

for link in etc app install doc; do
    if [ ! -s $link ]; then
	ln -s $TOSRC/$link $link
    fi
done

cd ..

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

/bin/cp $TOSRC/app/cpanfile .

carton install

cd $BUILDDIR

/bin/cp -R $CARTON/local/bin/* bin
/bin/cp -R $CARTON/local/lib/perl5/* lib/perl5

#
# Ant builds the rpm, perl modules should have been built
# by carton already and placed in the lib/perl5 directory.
echo -e "arch=x86_64\nto_version=$VERSION" > rpm/traffic_ops.properties
cd rpm && ant

if [ $? != 0 ]; then
    echo -e "\nRPM BUILD FAILED.\n\n"
else
    echo
    echo "========================================================================================"
    echo "RPM BUILD SUCCEEDED, See $BUILDDIR/dist/$RPM for the newly built rpm."
    echo "========================================================================================"
    echo
    #if [ $BRANCH != "master" ]; then
	 #   /usr/bin/git checkout master
    #fi
    #/bin/cp $BUILDDIR/rpm/dist/*.rpm $BUILDDIR/dist
fi

