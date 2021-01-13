ARG RHEL_VERSION=8
FROM golang:1.15.5 as builder

WORKDIR /go/src/github.com/apache
RUN git clone https://github.com/apache/trafficcontrol.git --branch 4.1.x --single-branch
WORKDIR trafficcontrol/traffic_stats
RUN go get ./...
RUN go build .


FROM centos:$RHEL_VERSION
ARG RHEL_VERSION=8

RUN /bin/echo "RHEL VERSION: $RHEL_VERSION"

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

#WORKDIR /root/
COPY --from=builder /go/bin/traffic_stats /opt/traffic_stats/bin/

RUN mkdir -p /opt/traffic_stats/var/run && mkdir -p /opt/traffic_stats/var/log

ADD run.sh /
ENTRYPOINT /run.sh
#CMD ["./traffic_stats","-cfg","/opt/traffic_stats/conf/traffic_stats.cfg"]
