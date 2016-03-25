# Docker file for CameraFeed microservice

FROM debian

# Need to run this to see what's in the various sources (need this to get
# curl).
RUN apt-get -qq update

##############################
#    Libraries
##############################

# Get g++, make, other build tools
RUN apt-get install -y build-essential automake autoconf libtool git checkinstall wget pkg-config curl libcurl3 libcurl4-gnutls-dev
# For cmake, see
# http://askubuntu.com/questions/610291/how-to-install-cmake-3-2-on-ubuntu-14-04
RUN mkdir /usr/src/cmake
WORKDIR /usr/src/cmake
RUN wget https://cmake.org/files/v3.5/cmake-3.5.0.tar.gz
RUN tar xf cmake-3.5.0.tar.gz
WORKDIR cmake-3.5.0
RUN ./configure && make && checkinstall

# Packages for node.js service code
RUN curl -sL https://deb.nodesource.com/setup_5.x | /bin/bash -
RUN apt-get install -y nodejs

# MongoDB install
# See https://github.com/mongodb/mongo-cxx-driver/wiki/Quickstart-Guide-%28New-Driver%29

RUN mkdir /usr/src/mongodb
WORKDIR /usr/src/mongodb
RUN git clone -b r1.3 https://github.com/mongodb/mongo-c-driver
WORKDIR /usr/src/mongodb/mongo-c-driver
RUN ./autogen.sh && make && make install
WORKDIR /usr/src/mongodb
RUN git clone -b master https://github.com/mongodb/mongo-cxx-driver
WORKDIR /usr/src/mongodb/mongo-cxx-driver/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local ..
RUN make && make install

##############################
#    Custom Code
##############################

# Brind in CameraFeed source
RUN mkdir /usr/src/CameraFeed
WORKDIR /usr/src/CameraFeed
# Manual copy, excludes Dockerfile, emacs tmp files
COPY *js ./
COPY keys ./
COPY src src/
COPY *json ./

RUN npm install

# Build C++ code
WORKDIR src
RUN make
