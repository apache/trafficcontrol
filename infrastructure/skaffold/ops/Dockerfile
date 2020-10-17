FROM tcgo:latest
WORKDIR /root/go/src/github.com/apache/trafficcontrol

# Copy source and deps
COPY vendor vendor
COPY lib lib
COPY traffic_ops traffic_ops

# Whatever go does to make vendor available to our executable
WORKDIR /root/go/src/github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang
RUN go get && go install
CMD [ "traffic_ops_golang" ]