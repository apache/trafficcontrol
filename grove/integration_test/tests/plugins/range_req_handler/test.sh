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
CMP_TOOL="${CMP_TOOL:-/compare_gets}"

#curl -H'Host: mem-test.cdn.kabletown.net' -Lsv -r 50000-50009  http://localhost:8080/10Mb.txt
originurl="http://localhost/"
host="mem-test.cdn.kabletown.net"
cacheurl="http://localhost:8080/"
file="10Mb.txt"

result=0
testno=0

# some basic and overlapping ranges
for host in "store-ranges.cdn.kabletown.net" "get-full.cdn.kabletown.net"
do
  for r in "0-0" "0-100" "5000-" "-100" "8-10,9-15,100-200" "0-300,200-250" "-33,66-99,50-150" "8192-20000" "8191-8192" "2000-3012" "1021-"
  do
    test="${CMP_TOOL}  --chdrs \"Host:$host Range:bytes=${r}\" --ohdrs \"Range:bytes=${r}\" --path \"10Mb.txt\" --ignorehdrs \"Server,Date\""
    testno=$(($testno+1))
    echo -n "Test $testno ($test): "

    ${CMP_TOOL}  --chdrs "Host:$host Range:bytes=${r}" --ohdrs "Range:bytes=${r}" --path "10Mb.txt" --ignorehdrs "Server,Date"

    result=$(($result+$?))
  done
done

# multipart
for host in "store-ranges.cdn.kabletown.net" "get-full.cdn.kabletown.net"
do
  for r in "0-0,10-15" "0-100,200-210" "33-99,101-188" "8-10,9-15,100-200" "0-300,200-250" "-33,66-99,50-150" "300-304,500-,600-700"
  do
    test="${CMP_TOOL}  --chdrs \"Host:$host Range:bytes=${r}\" --ohdrs \"Range:bytes=${r}\" --path \"10Mb.txt\" --ignorehdrs \"Server,Date\" --ignorempb"
    testno=$(($testno+1))
    echo -n "Test $testno ($test): "

    ${CMP_TOOL}  --chdrs "Host:$host Range:bytes=${r}" --ohdrs "Range:bytes=${r}" --path "10Mb.txt" --ignorehdrs "Server,Date" --ignorempb

    result=$(($result+$?))
  done
done

host="slice.cdn.kabletown.net"
for r in "0-1048567" "22-2000000" "0-0" "8388608-9437183" "2000-3000" "1099-3033" "50001-61111" "121212-121212" "121212-121215" "001-313" "4096-8191" "4096-8192" "4000-9000" "121200-121901" "0-5000" "6000-7000" "0-100" "5000-" "-100" "8191-8192" "2000-3012" "1021-"
  do
  test="${CMP_TOOL}  --chdrs \"Host:$host Range:bytes=${r}\" --ohdrs \"Range:bytes=${r}\" --path \"10Mb.txt\" --ignorehdrs \"Server,Date\""
  testno=$(($testno+1))
  echo -n "Test $testno ($test): "

  ${CMP_TOOL}  --chdrs "Host:$host Range:bytes=${r}" --ohdrs "Range:bytes=${r}" --path "10Mb.txt" --ignorehdrs "Server,Date"

  result=$(($result+$?))
done

# "normal" GET
for host in "store-ranges.cdn.kabletown.net" "get-full.cdn.kabletown.net" "slice.cdn.kabletown.net"
do
  test="${CMP_TOOL}  --chdrs \"Host:$host\" --path \"10Mb.txt\" --ignorehdrs \"Server,Date\""
  testno=$(($testno+1))
  echo -n "Test $testno ($test): "

  ${CMP_TOOL}  --chdrs "Host:$host Range:bytes=${r}" --ohdrs "Range:bytes=${r}" --path "10Mb.txt" --ignorehdrs "Server,Date"

  result=$(($result+$?))
done

echo "plugin/range_req_handler: $testno tests done, $result failed."

exit $result
