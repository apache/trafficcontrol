FROM tcgo:latest
WORKDIR /go/root/src/github.com/apache/trafficcontrol

COPY lib lib
COPY traffic_ops traffic_ops

