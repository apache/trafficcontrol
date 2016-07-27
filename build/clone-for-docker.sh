#!/bin/sh

project=`basename $GITREPO .git`

[ ! -d /repo ] || ( echo "Expected to have volume mounted for /repo"; exit 1)

# Delete previous content if it exists
rm -rf /repo/$project

git clone --branch $BRANCH $GITREPO
