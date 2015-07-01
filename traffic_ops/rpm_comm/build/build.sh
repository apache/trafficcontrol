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

#-----------------------------------------------------------------------------
function usage () {
   echo "./build.sh [-b <branch>] [--build <build_number>] [-g <gitrepo>]"
   echo ""
   echo "Don't run this program ever."
   echo "   -b | --branch     Git branch"
   echo "                     default: master"
   echo "        --build      Build number"
   echo "                     default: commit hash"
   echo "   -g | --gitrepo    Git repository."
   echo "                     default: https://github.com/Comcast/traffic_control.git"
   echo "   -h | --help       Print this message"
   echo "   -w | --workspace  Working directory"
   echo "                     default: home directory"
   echo ""
}

function createWorkspace () {
   mkdir -p $workspace/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
   #echo '%_topdir %(echo $workspace)/rpmbuild' > ~/.rpmmacros
   cp $TOSRC/rpm_comm/rpmmacros ~/.rpmmacros

   VERSION=$(/bin/cat $UTILS_PM|/bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')

   echo "%traffic_ops_version $VERSION" >> ~/.rpmmacros
   echo "%traffic_ops_build $build" >> ~/.rpmmacros
   echo "%hosttype $HOSTTYPE" >> ~/.rpmmacros

}

function downloadRepo () {

   if [ ! -d  $workspace/repos ]; then
      mkdir -p $workspace/repos
   fi

   cd $workspace/repos

   if [ -d $TOSRC ]; then
      echo "cd to $TCSRC and then git pull"
      cd $TCSRC && /usr/bin/git pull
      echo "git checkout $BRANCH"
      /usr/bin/git checkout $BRANCH
      /usr/bin/git pull
   else
      echo "cd to $REPODIR and clone"
      cd $REPODIR && /usr/bin/git clone $gitrepo
      cd $TCSRC
      /usr/bin/git checkout $BRANCH
   fi

}

function runCarton () {
   sudo perl -MCPAN -e 'my $c = "CPAN::HandleConfig"; $c->load(doit => 1, autoconfig => 1); $c->edit(prerequisites_policy => "follow"); $c->edit(build_requires_install_policy => "yes"); $c->commit'
   sudo cpan -if MIYAGAWA/Carton-v1.0.15.tar.gz
}

function combine () {
   if [ -d $COMBINEDIR ]; then
      echo "removing $COMBINEDIR"
      rm -dfr $COMBINEDIR
   fi

   mkdir $COMBINEDIR

}

#-----------------------------------------------------------------------------
# MAIN
#-----------------------------------------------------------------------------

while [ "$1" != "" ]; do
   case $1 in
      -b | --branch )         shift
                              branch=$1
                              ;;
      --build )               shift
                              build=$1
                              ;;
      #-f | --file )           shift
      #                        filename=$1
      #                        ;;
      -g | --gitrepo )        shift
                              gitrepo=$1
                              ;;
      # Example of how to do command line without following var.
      #-i | --interactive )    interactive=1
      #                        ;;
      -h | --help )           usage
                              exit
                              ;;
      -w | --workspace )      shift
                              workspace=$1
                              ;;
      * )                     usage
                              exit 1
   esac
   shift
done

if [ "$branch" = "" ]; then
   branch="master"
   echo "Setting branch to $branch"
fi

if [ "$gitrepo" = "" ]; then
   gitrepo="https://github.com/Comcast/traffic_control.git"
   echo "Setting gitrepo to $gitrepo"
fi

if [ "$workspace" = "" ]; then
   workspace=~
   echo "Setting workspace to $workspace"
fi

if [ "$build" = "" ]; then
   echo "build not set. will use commit hash"
fi

# set vars
PACKAGE="traffic_ops"
REPODIR="$workspace/repos"
COMBINEDIR="$workspace/traffic_ops_combine"
RPMBUILDDIR="$workspace/rpmbuild/BUILD"
SOURCES="$workspace/rpmbuild/SOURCES"
SPECS="$workspace/rpmbuild/SPECS"
TCSRC="$REPODIR/traffic_control"
TOSRC="$TCSRC/$PACKAGE"
CARTONDIR="$workspace/carton"
UTILS_PM="$TOSRC/app/lib/UI/Utils.pm"

# Read in default profile
if [ -f /etc/profile ]; then
    . /etc/profile
fi

exit 0

createWorkspace
downloadRepo
runCarton
combine

exit 0


go get github.com/go-sql-driver/mysql
go get code.google.com/p/go.net/html
go get code.google.com/p/go.net/publicsuffix

echo "package: $PACKAGE"
echo "tcsrc: $TCSRC"
echo "tosrc: $TOSRC"

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


/bin/cp -R $TOSRC/rpm_comm/$PACKAGE.spec $SPECS

cd $TOSRC/install/bin
go build $TOSRC/install/go/src/comcast.com/dataload/dataload.go
go build $TOSRC/install/go/src/comcast.com/systemtest/systemtest.go

mkdir $COMBINEDIR/traffic_ops-$VERSION
cd $COMBINEDIR/traffic_ops-$VERSION

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

/bin/cp -R $CARTON/local/bin/* $COMBINEDIR/traffic_ops-$VERSION/bin
/bin/cp -R $CARTON/local/lib/perl5/* $COMBINEDIR/traffic_ops-$VERSION/lib/perl5
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

