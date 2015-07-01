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
   echo "Don't run this script ever."
   echo "   -b | --branch     Git branch"
   echo "                     default: master"
   #echo "        --build      Build number"
   #echo "                     default: commit hash"
   echo "   -g | --gitrepo    Git repository."
   echo "                     default: https://github.com/Comcast/traffic_control.git"
   echo "   -h | --help       Print this message"
   echo "   -w | --workspace  Working directory"
   echo "                     default: home directory"
   echo ""
}

function createWorkspace () {
   echo ""
   echo ""
   echo "##################################################################"
   echo "# Creating Workspace"
   echo "##################################################################"

   mkdir -p $workspace/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
   #echo '%_topdir %(echo $workspace)/rpmbuild' > ~/.rpmmacros

}

function downloadRepo () {

   echo ""
   echo ""
   echo "##################################################################"
   echo "# Downloading Repo"
   echo "##################################################################"

   if [ ! -d  $workspace/repos ]; then
      mkdir -p $workspace/repos
   fi

   cd $workspace/repos

   if [ -d $TOSRC ]; then
      echo "cd to $TCSRC and then git pull"
      cd $TCSRC && /usr/bin/git pull
      echo "git checkout $branch"
      /usr/bin/git checkout $branch
      /usr/bin/git pull
   else
      echo "cd to $REPODIR and clone"
      cd $REPODIR && /usr/bin/git clone $gitrepo
      cd $TCSRC
      /usr/bin/git checkout $branch
   fi

   cp $TOSRC/rpm_comm/rpmmacros ~/.rpmmacros

   VERSION=$(/bin/cat $UTILS_PM|/bin/awk '/my \$version/{split($4,a,"\"");split(a[2],b,"-");printf("%s",b[1])}')

   echo "%traffic_ops_version $VERSION" >> ~/.rpmmacros
   echo "%hosttype $HOSTTYPE" >> ~/.rpmmacros

   #COMMITS=$(git rev-list HEAD --count)
   COMMITS=$(git shortlog | grep -E '^[ ]+\w+' | wc -l)
   SHA=$(git rev-parse --short=8 HEAD)
   BUILD_NUMBER="$COMMITS-$SHA"
   echo "%traffic_ops_build $BUILD_NUMBER" >> ~/.rpmmacros
   echo "%traffic_ops_release $COMMITS" >> ~/.rpmmacros
   echo "%traffic_ops_sha $SHA" >> ~/.rpmmacros
}

function getWebDeps () {

   echo ""
   echo ""
   echo "##################################################################"
   echo "# Getting web dependencies"
   echo "##################################################################"

   sudo cpan -if WWW::Curl::Easy
   sudo cpan -if IO::Uncompress::Unzip

   cd $TOSRC/install/bin
   ./download_web_deps

   rc=$?
   if [ $rc != 0 ]; then
      echo "download web dependencies failed."
      exit $rc
   fi

}

function runCarton () {

   echo ""
   echo ""
   echo "##################################################################"
   echo "# Running Carton"
   echo "##################################################################"

   sudo cpan -if MIYAGAWA/Carton-v1.0.15.tar.gz

   if [ ! -d $CARTONDIR ]; then
       /bin/mkdir $CARTONDIR
   fi

   cd $CARTONDIR

   /bin/cp $TOSRC/app/cpanfile .

   carton install

   rc=$?
   if [ $rc != 0 ]; then
      echo "carton download failed"
      exit $rc
   fi

}

function combine () {


   echo ""
   echo ""
   echo "##################################################################"
   echo "# Combining things together"
   echo "##################################################################"

   if [ -d $COMBINEDIR ]; then
      echo "removing $COMBINEDIR"
      rm -dfr $COMBINEDIR
   fi

   mkdir $COMBINEDIR
   mkdir $COMBINEDIR/traffic_ops-$VERSION-$BUILD_NUMBER
   cd $COMBINEDIR/traffic_ops-$VERSION-$BUILD_NUMBER

   if [ -d lib ]; then
       /bin/rm -rf lib
   fi

   /bin/mkdir -p lib/perl5

   if [ -d bin ]; then
       /bin/rm -rf bin
   fi

   /bin/mkdir bin

   # copy carton files
   /bin/cp -R $CARTONDIR/local/bin/* $COMBINEDIR/traffic_ops-$VERSION-$BUILD_NUMBER/bin
   /bin/cp -R $CARTONDIR/local/lib/perl5/* $COMBINEDIR/traffic_ops-$VERSION-$BUILD_NUMBER/lib/perl5

   # copy Traffic Ops source
   cp -r $TOSRC/* $COMBINEDIR/traffic_ops-$VERSION-$BUILD_NUMBER/

   cd $COMBINEDIR

   tar -czf $SOURCES/$PACKAGE-$VERSION-$BUILD_NUMBER.$HOSTTYPE.tar.gz ./*


}

function buildRpm () {

   echo ""
   echo ""
   echo "##################################################################"
   echo "# Building the RPM"
   echo "##################################################################"

   /bin/cp -R $TOSRC/rpm_comm/$PACKAGE.spec $SPECS
   go get github.com/go-sql-driver/mysql
   go get code.google.com/p/go.net/html
   go get code.google.com/p/go.net/publicsuffix

   cd $TOSRC/install/bin
   go build $TOSRC/install/go/src/comcast.com/dataload/dataload.go
   go build $TOSRC/install/go/src/comcast.com/systemtest/systemtest.go

   echo
   echo "=================================================================="
   echo "Building Traffic Ops rpm traffic_ops-$VERSION-$BUILD_NUMBER"
   echo
   echo "GOPATH=$GOPATH"
   echo "RPMBUILDDIR=$RPMBUILDDIR"
   echo "TOSRC=$TOSRC"
   echo "CARTON=$CARTON"
   echo "UTILS_PM=$TOSRC/app/lib/UI/Utils.pm"
   echo "COMBINEDIR=$COMBINEDIR"
   echo "VERSION=$VERSION"
   echo "BUILD_NUMBER=$BUILD_NUMBER"
   echo "=================================================================="
   echo

   cd $SPECS
   rpmbuild -ba $PACKAGE.spec

   rc=$?
   if [ $rc != 0 ]; then
      echo "rpm build failed"
      exit $rc
   else
      echo
      echo "================================================================================"
      echo "RPM BUILD SUCCEEDED"
      echo "See $RPMS/$HOSTTYPE/traffic_ops-$VERSION-$BUILD_NUMBER.rpm for the newly built rpm."
      echo "================================================================================"
      echo
   fi
}

#-----------------------------------------------------------------------------
# MAIN
#-----------------------------------------------------------------------------

while [ "$1" != "" ]; do
   case $1 in
      -b | --branch )         shift
                              branch=$1
                              ;;
      #--build )               shift
      #                        build=$1
      #                        ;;
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
   echo "branch not set"
   exit 1
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

#if [ "$build" = "" ]; then
#   echo "build not set. will use commit hash"
#fi

# set vars
PACKAGE="traffic_ops"
REPODIR="$workspace/repos"
COMBINEDIR="$workspace/traffic_ops_combine"
RPMBUILDDIR="$workspace/rpmbuild/BUILD"
SOURCES="$workspace/rpmbuild/SOURCES"
SPECS="$workspace/rpmbuild/SPECS"
RPMS="$workspace/rpmbuild/RPMS"
TCSRC="$REPODIR/traffic_control"
TOSRC="$TCSRC/$PACKAGE"
CARTONDIR="$workspace/carton"
UTILS_PM="$TOSRC/app/lib/UI/Utils.pm"
BUILD_NUMBER=""
VERSION=""

# Read in default profile
if [ -f /etc/profile ]; then
    . /etc/profile
fi

# Tell cpan to answer yes
sudo perl -MCPAN -e 'my $c = "CPAN::HandleConfig"; $c->load(doit => 1, autoconfig => 1); $c->edit(prerequisites_policy => "follow"); $c->edit(build_requires_install_policy => "yes"); $c->commit'

createWorkspace
downloadRepo
runCarton
getWebDeps
combine
buildRpm

exit 0


