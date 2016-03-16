# Docker file for CameraFeed microservice

FROM debian

# Need to run this to see what's in the various sources (need this to get
# curl).
RUN apt-get -qq update

# Get g++, make
RUN apt-get install -y build-essential

# Packages for node.js service code
RUN apt-get install -y curl
RUN curl -sL https://deb.nodesource.com/setup_5.x | /bin/bash -
RUN apt-get install -y nodejs

# Packages used by C++ code to read camera stream
RUN apt-get install -y libcurl3 libcurl4-gnutls-dev

# Brind in CameraFeed source
RUN mkdir /usr/src/CameraFeed
WORKDIR /usr/src/CameraFeed
# Manual copy, excludes Dockerfile, emacs tmp files
COPY *js .
COPY keys .
COPY src .
COPY *json .

RUN npm install

# Build C++ code
WORKDIR src
RUN make
