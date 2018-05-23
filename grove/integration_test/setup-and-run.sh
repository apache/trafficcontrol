#!/usr/bin/bash
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
go get golang.org/x/text
go get golang.org/x/sys/unix
go get golang.org/x/net/http2
go get golang.org/x/net/ipv4
go get golang.org/x/net/ipv6

mkdir -p $GOPATH/src/github.com/apache/
cd $GOPATH/src/github.com/apache/
#git clone https://github.com/apache/incubator-trafficcontrol
git clone $REPO
cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove
git checkout $BRANCH
go build

cd /
openssl req -newkey rsa:2048 -new -nodes -x509 -days 365 -keyout key.pem -out cert.pem -subj "/C=US/ST=CO/L=Denver/O=.../OU=.../CN=.../emailAddress=..."

cp /remap-base-test.json /remap.json
ls -l
${GOPATH}/src/github.com/apache/incubator-trafficcontrol/grove/grove -cfg grove.cfg &


sleep 3
curl -H'Host: mem-test.cdn.kabletown.net' -Lsv -r 50000-50009  http://localhost:8080/10Mb.txt

#cd $GOPATH/src/github.com/apache/incubator-trafficcontrol/grove/integration_test
go build compare_gets.go




