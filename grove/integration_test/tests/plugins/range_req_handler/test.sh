#!/usr/bin/env bash -x

#curl -H'Host: mem-test.cdn.kabletown.net' -Lsv -r 50000-50009  http://localhost:8080/10Mb.txt
originurl="http://localhost/"
host="mem-test.cdn.kabletown.net"
cacheurl="http://localhost:8080/"
file="10Mb.txt"

result=0
testno=0

for host in "mem-test.cdn.kabletown.net" "disk-test.cdn.kabletown.net"
do
  for r in "0-0" "0-100" "5000-" "-100" "0-010-15" "0-100200-210" "33-9966-88" "-"
  do
    test="/compare_gets  --chdrs \"Host:$host Range:bytes=${r}\" --ohdrs \"Range:bytes=${r}\" --path \"10Mb.txt\" --ignorehdrs \"Server,Date\""
    testno=$(($testno+1))
    echo -n "Test $testno ($test): "

    /compare_gets  --chdrs "Host:$host Range:bytes=${r}" --ohdrs "Range:bytes=${r}" --path "10Mb.txt" --ignorehdrs "Server,Date"

    result=$(($result+$?))
  done
done


echo "$testno tests done, $result failed."

exit $result
