#!/usr/bin/env bash

target=$1
[[ -z $target ]] && echo "No target specified"
echo "Building $target"

echo "GITREPO=${GITREPO:=https://github.com/apache/incubator-trafficcontrol}"
echo "BRANCH=${BRANCH:=master}"

dir=$(basename $GITREPO)
set -x
git clone "$GITREPO" -b "$BRANCH" $dir || echo "Clone failed: $!"

cd $dir/$target
./build/build_rpm.sh
mkdir -p /artifacts
cp ../dist/* /artifacts/.

# Clean up for next build
cd -
rm -r $dir
