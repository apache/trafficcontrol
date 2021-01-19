#!/bin/bash

set -o posix
set -o errexit
set -ex
echo "Test work is in progress. Test Suites are added gradually."
cd ./cdn-tests/trafficportal
./runtests.bash $URL
exit_code=$?
exit $exit_code
