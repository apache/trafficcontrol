#!/usr/bin/bash

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

echo Configuring Grove Integration Test Environment.

#perl -v
#ls -l /var/www/html
cd /var/www/html
echo Generating origin test files...
perl -e 'foreach $i ( 0 ... 1024*1024-1 ) { printf "%09d\n", $i*10 }' > 10Mb.txt
for i in {1..1000} ; do  dd if=/dev/urandom of=${i}k.bin bs=${i}k count=1 > /dev/null 2>&1 ; done
httpd

cd /
echo Setting up go enviroment...
export GOPATH=~/go
go mod vendor -v

mkdir -p $GOPATH/src/github.com/apache/
cd $GOPATH/src/github.com/apache/
#git clone https://github.com/apache/trafficcontrol
git clone $REPO
cd $GOPATH/src/github.com/apache/trafficcontrol/grove
git checkout $BRANCH
go build

cd /
openssl req -newkey rsa:2048 -new -nodes -x509 -days 365 -keyout key.pem -out cert.pem -subj "/C=US/ST=CO/L=Denver/O=.../OU=.../CN=.../emailAddress=..."

cp /remap-base-test.json /remap.json
ls -l
${GOPATH}/src/github.com/apache/trafficcontrol/grove/grove -cfg grove.cfg &


sleep 3
curl -H'Host: mem-test.cdn.kabletown.net' -Lsv -r 50000-50009  http://localhost:8080/10Mb.txt

#cd $GOPATH/src/github.com/apache/trafficcontrol/grove/integration_test
go build compare_gets.go

#function run_test () {
#  command=${1}
#  output=${2:-/dev/null}
#  echo -n "Test ${testno} (${command}): "
#  ${command} > ${output} 2>&1
#  thisresult=$?
#
#  result=$(($result+thisresult))
#  if [ $thisresult -eq 0 ]
#  then
#    echo "PASS"
#  else
#    echo "FAIL"
#  fi
#
#  testno=$(($testno+1))
#}


function run_test () {
  "$@" > /tmp/run_test.out 2>&1
  thisresult=$?

   echo -n "Test ${testno} ($@): "
  result=$(($result+thisresult))
  if [ $thisresult -eq 0 ]
  then
    echo "PASS"
  else
    echo "FAIL"
  fi

  testno=$(($testno+1))
}
export -f run_test

cp  $GOPATH/src/github.com/apache/trafficcontrol/grove/integration_test/tests/plugins/modify_headers/remap.json /remap.json
pkill -HUP grove
bash $GOPATH/src/github.com/apache/trafficcontrol/grove/integration_test/tests/plugins/modify_headers/test.sh

cp $GOPATH/src/github.com/apache/trafficcontrol/grove/integration_test/tests/plugins/range_req_handler/remap.json /remap.json
pkill -HUP grove
bash $GOPATH/src/github.com/apache/trafficcontrol/grove/integration_test/tests/plugins/range_req_handler/test.sh


