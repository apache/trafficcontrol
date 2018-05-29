#!/usr/bin/env bash -x

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


originurl="http://localhost/"
host="mem-test.cdn.kabletown.net"
cacheurl="http://localhost:8080/"
file="10Mb.txt"

result=0
testno=0


# test global header set
run_test curl -Lkvs -HHost:mem-test.cdn.kabletown.net -o /tmp/10k.bin http://localhost:8080/10k.bin
cp /tmp/run_test.out /tmp/hdrs.out
run_test diff /tmp/10k.bin /var/www/html/10k.bin
run_test grep "< Server: Grove/0.39999999"  /tmp/hdrs.out

# test per remap header set up and down
run_test curl -XTRACE -H'Host: disk-test.cdn.kabletown.net' http://localhost:8080/10Mb.txt -Lkvs -o /tmp/out
cp /tmp/run_test.out /tmp/hdrs.out
run_test grep "X-From-Cdn: Traffic-Control" /tmp/out
run_test grep "< X-Cdn-Name: GroverCDN" /tmp/hdrs.out

echo "plugin/modify_headers: $testno tests done, $result failed."

exit $result
