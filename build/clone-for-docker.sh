#!/bin/sh

project=`basename $GITREPO .git`
mkdir -p /vol

# Delete previous content if it exists
rm -rf /vol/$project

cd /vol
git clone --branch $BRANCH $GITREPO
