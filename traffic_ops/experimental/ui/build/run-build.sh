#!/usr/bin/env bash

target=$1
[[ -z $target ]] && echo "No target specified"
echo "Building $target"

echo "GITREPO=${GITREPO:=https://github.com/apache/incubator-trafficcontrol}"
echo "BRANCH=${BRANCH:=master}"

set -x
git clone $GITREPO -b $BRANCH traffic_control
distdir=$(pwd)/traffic_control/dist

cd traffic_control/traffic_ops/experimental/ui
./build/build_rpm.sh
mkdir -p /artifacts
cp $distdir/* /artifacts/.

# Clean up for next build
cd -
rm -r traffic_control
