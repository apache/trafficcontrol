#!/usr/bin/env bash -x

#curl -H'Host: mem-test.cdn.kabletown.net' -Lsv -r 50000-50009  http://localhost:8080/10Mb.txt
originurl="http://localhost/"
host="mem-test.cdn.kabletown.net"
cacheurl="http://localhost:8080/"
file="10Mb.txt"

result=0
testno=0

#curl -s -r 50000-50009 ${originurl}${file} > /tmp/out1 && echo FAIL test ${testno}
#result=$(($result+$?))
#testno=$(($testno+1))
#
#curl -s -r 50000-50009  -H"Host: ${host}" ${cacheurl}/${file} > /tmp/out2 && echo FAIL test ${testno}
#
#result=$(($result+$?))
#testno=$(($testno+1))
#
#diff /tmp/out1 /tmp/out2 && echo FAIL test ${testno}
#result=$(($result+$?))
#testno=$(($testno+1))

for host in "mem-test.cdn.kabletown.net", "disk1-test.cdn.kabletown.net"
do
  for r in "0-0", "0-100", "5000-", "-100", "0-0,10-15", "0-100,200-210", "33-99,66-88" "-"
  do
    test="/compare_gets  --chdrs \"Host:$host,Range:bytes=\r${r}\" --ohdrs \"Range:bytes=${r}\" --path \"10Mb.txt\" --ignorehdrs \"Server,Date\""
    testno=$(($testno+1))
    echo -n "Test $testno ($test): "

    /compare_gets  --chdrs "Host:$host,Range:bytes=${r}" --ohdrs "Range:bytes=${r}" --path "10Mb.txt" --ignorehdrs "Server,Date"

    result=$(($result+$?))
  done
done


echo "$testno tests done, $result failed."

exit $result
