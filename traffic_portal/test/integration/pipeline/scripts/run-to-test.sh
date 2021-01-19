#!/bin/bash

set -o posix
set -o errexit
set -ex
echo "Test work is in progress. Test Suites are added gradually."
cd /go/src/github.comcast.com/cdn/cdn-tests/trafficops
./runtests.bash $API_VERSION $ENVIRONMENT
exit_code=$?
exit $exit_code
