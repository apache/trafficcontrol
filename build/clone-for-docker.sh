#!/bin/sh

set -x
if [ $# -ne 2 ]
then
	echo "Usage: $0 <gitrepo> <branch>"
	exit 1
fi

gitrepo=$1
branch=$2

project=`basename $gitrepo .git`
mkdir -p /vol

# Delete previous content if it exists
rm -rf /vol/$project

cd /vol
git clone --branch $branch $gitrepo
