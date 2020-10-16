#
# Download and install Go
#
FROM centos:7 as install

RUN yum update && yum install -y wget
RUN wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz
RUN tar -C /usr/local/ -xzf go1.15.3.linux-amd64.tar.gz

#
# Copy Go and add it to the PATH
#
FROM centos:7

RUN yum update && yum install -y git && yum clean all
COPY --from=install /usr/local/go /usr/local/go
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin
ENV GOPATH=/root/go