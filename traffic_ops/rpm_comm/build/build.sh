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
   echo "If running from a Vagrant VM WORKSPACE should have been set to "
   echo " /home/vagrant."
   echo "example: export WORKSPACE=/home/vagrant"
	exit 1
fi

if [ ! -d $WORKSPACE ]; then
   echo "$WORKSPACE does not exist."
   exit 1
fi

mkdir -p $WORKSPACE/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
echo '%_topdir %(echo $WORKSPACE)/rpmbuild' > ~/.rpmmacros

mkdir -p $WORKSPACE/repos

if [ -d $WORKSPACE/traffic_ops_combine ]; then
   echo "removing $WORKSPACE/traffic_ops_combine"
   rm -dfr $WORKSPACE/traffic_ops_combine
fi

mkdir $WORKSPACE/traffic_ops_combine

#sudo perl -MCPAN -e 'my $c = "CPAN::HandleConfig"; $c->load(doit => 1, autoconfig => 1); $c->edit(prerequisites_policy => "follow"); $c->edit(build_requires_install_policy => "yes"); $c->commit'
#sudo cpan -if MIYAGAWA/Carton-v1.0.15.tar.gz

go get github.com/go-sql-driver/mysql
go get code.google.com/p/go.net/html
go get code.google.com/p/go.net/publicsuffix

GITREPO="https://github.com/Comcast/traffic_control.git"
PACKAGE="traffic_ops"
REPODIR="$WORKSPACE/repos"
COMBINEDIR="$WORKSPACE/traffic_ops_combine"
RPMBUILDDIR="$WORKSPACE/rpmbuild/BUILD"
SOURCES="$WORKSPACE/rpmbuild/SOURCES"
SPECS="$WORKSPACE/rpmbuild/SPECS"
TCSRC="$REPODIR/traffic_control"
TOSRC="$TCSRC/$PACKAGE"
CARTON="$WORKSPACE/carton"
UTILS_PM="$TOSRC/app/lib/UI/Utils.pm"

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

if [ ! -z $3 ]; then
	GITREPO=$3
fi

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
   echo "cd to $REPODIR and clone"
   cd $REPODIR && /usr/bin/git clone $GITREPO
   cd $TCSRC
   /usr/bin/git checkout $BRANCH
fi

cp $TOSRC/rpm_comm/rpmmacros ~/.rpmmacros

VERSION=$(/bin/cat $UTILS_PM|/bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')

echo "%traffic_ops_version $VERSION" >> ~/.rpmmacros
echo "%traffic_ops_build $BUILD_NUMBER" >> ~/.rpmmacros
echo "%hosttype $HOSTTYPE" >> ~/.rpmmacros

#RPM="${PACKAGE}-${VERSION}-${BUILD_NUMBER}.x86_64.rpm"

echo
echo "=================================================================="
echo "Building Traffic Ops rpm traffic_ops-${VERSION}-$BUILD_NUMBER"
echo
echo "GOPATH=$GOPATH"
echo "RPMBUILDDIR=$RPMBUILDDIR"
echo "TOSRC=$TOSRC"
echo "CARTON=$CARTON"
echo "UTILS_PM=$TOSRC/app/lib/UI/Utils.pm"
echo "COMBINEDIR=$COMBINEDIR"
echo "=================================================================="
echo


#cd $COMBINEDIR

# TODO check to see what Comcast was doing with this
#if [ ! -d dist ]; then
#    /bin/mkdir dist
#fi

/bin/cp -R $TOSRC/rpm_comm/$PACKAGE.spec $SPECS
#echo "build.number=$BUILD_NUMBER" > $BUILDDIR/build.number

cd $TOSRC/install/bin
go build $TOSRC/install/go/src/comcast.com/dataload/dataload.go
go build $TOSRC/install/go/src/comcast.com/systemtest/systemtest.go

if [ ! -d $CARTON ]; then
    /bin/mkdir $CARTON
fi

cd $CARTON

/bin/cp $TOSRC/app/cpanfile .

#carton install

mkdir $COMBINEDIR/traffic_ops-$VERSION
cd $COMBINEDIR/traffic_ops-$VERSION

#for link in etc app install doc; do
#   if [ ! -s $link ]; then
#      ln -s $TOSRC/$link $link
#   fi
#done

if [ -d lib ]; then
    /bin/rm -rf lib
fi

/bin/mkdir -p lib/perl5

if [ -d bin ]; then
    /bin/rm -rf bin
fi

/bin/mkdir bin

/bin/cp -R $CARTON/local/bin/* bin
/bin/cp -R $CARTON/local/lib/perl5/* lib/perl5
for directory in etc app install doc; do
   cp -r $TOSRC/$directory $COMBINEDIR/traffic_ops-$VERSION
done

cd $COMBINEDIR

tar -czf $SOURCES/$PACKAGE-$VERSION-$BUILD_NUMBER.$HOSTTYPE.tar.gz ./*

cd $SPECS
rpmbuild -ba $PACKAGE.spec

#
# Ant builds the rpm, perl modules should have been built
# by carton already and placed in the lib/perl5 directory.
#echo -e "arch=x86_64\nto_version=$VERSION" > $BUILDDIR/traffic_ops.properties
#cd rpm && ant

#if [ $? != 0 ]; then
#    echo -e "\nRPM BUILD FAILED.\n\n"
#else
#    echo
#    echo "========================================================================================"
#    echo "RPM BUILD SUCCEEDED, See $BUILDDIR/dist/$RPM for the newly built rpm."
#    echo "========================================================================================"
#    echo
#    #if [ $BRANCH != "master" ]; then
#	 #   /usr/bin/git checkout master
#    #fi
#    #/bin/cp $BUILDDIR/rpm/dist/*.rpm $BUILDDIR/dist
#fi

