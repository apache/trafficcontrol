#!/bin/sh
git clone ${GITREPO}
cd `basename ${GITREPO} .git`
git checkout ${BRANCH}
