#!/bin/sh

project=`basename $GITREPO .git`
mkdir -p /repo

# Delete previous content if it exists
rm -rf /repo/$project

cd /repo
git clone --branch $BRANCH $GITREPO
