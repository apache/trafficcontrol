#!/bin/bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

function die() {
  { test -n "$@" && echo "$@"; exit 1; } >&2 
}

OS_VERSION="${OS_VERSION}"
OS_DISTRO=${OS_DISTRO}
ATS_VERSION="${ATS_VERSION}"
ATS_MAJOR_VERSION=${ATS_VERSION::1}
echo "OS_DISTRO:${OS_DISTRO}"
echo "OS_VERSION:${OS_VERSION}"
echo "ATS_VERSION:${ATS_VERSION}"
echo "ATS_MAJOR_VERSION:${ATS_MAJOR_VERSION}"

mkdir -p /opt/build
cd /opt/build

# build openssl 1.1.1 if OS_VERSION is not 8 or greater.
if [ ${OS_VERSION%%.*} -le 7 ]; then
  git clone $OPENSSL_URL --branch $OPENSSL_TAG || die "Failed to fetch the OpenSSL Source."
  (
    cd /opt/build/openssl && 
    ./config --prefix=/opt/trafficserver/openssl --openssldir=/opt/trafficserver/openssl zlib && 
    make -j$(nproc) && make install_sw
  ) || die "Failed to build OpenSSL"
  cjose_openssl='--with-openssl=/opt/trafficserver/openssl'
  rpmbuild_openssl='--with openssl_included'
else
  cjose_openssl=''
  rpmbuild_openssl='--without openssl_included'
fi

# Build jansson
(
  git clone $JANSSON_URL --branch $JANSSON_TAG 
  cd /opt/build/jansson && patch -p1 < /jansson.pic.patch && 
    autoreconf -i && ./configure --enable-shared=no && make -j &&
    make install
) || die "Failed to install jansson from source."

# Build and install cjose
(
  git clone $CJOSE_URL --branch $CJOSE_TAG 
  cd /opt/build/cjose && patch -p1 < /cjose.pic.patch && 
    autoreconf -i && ./configure --enable-shared=no \
    ${cjose_openssl} && make -j$(nproc) && make install
) || die "Falled to build cjose from source."

# prep build environment
cd /root
[ -e rpmbuild ] && rm -rf rpmbuild
[ ! -e rpmbuild ] || { echo "Failed to clean up rpm build directory 'rpmbuild': $?" >&2; exit 1; }
mkdir -p rpmbuild/{BUILD,BUILDROOT,RPMS,SPECS,SOURCES,SRPMS} || die "Failed to initialize the build environment"

echo "Building a RPM for ATS version: $ATS_VERSION and OS version: $OS_VERSION"

# add the 'ats' user
id ats &>/dev/null || /usr/sbin/useradd -u 176 -r ats -s /sbin/nologin -d /

# setup the environment to use the devtoolset-11 tools.
if [ "${OS_VERSION%%.*}" -le 7 ]; then 
  source scl_source enable devtoolset-11
else
  source scl_source enable gcc-toolset-11
fi

cd /root
# prep build environment
[ -e rpmbuild ] && rm -rf rpmbuild
[ ! -e rpmbuild ] || { echo "Failed to clean up rpm build directory 'rpmbuild': $?" >&2; exit 1; }
mkdir -p rpmbuild/{BUILD,BUILDROOT,RPMS,SPECS,SOURCES,SRPMS} || die "Failed to create build directory '$RPMBUILD': $?"

cd /root/rpmbuild/SOURCES
# clone the trafficserver repo
git clone https://github.com/apache/trafficserver.git --branch $ATS_VERSION || die "Failed to fetch the ATS Source"
cp /traffic_server_jemalloc .
cp /trafficserver.env .

# patch in the astats plugin
(cp -fa /astats_over_http /root/rpmbuild/SOURCES/trafficserver/plugins/astats_over_http

cat > /root/rpmbuild/SOURCES/trafficserver/plugins/astats_over_http/Makefile.inc <<MAKEFILE
pkglib_LTLIBRARIES += astats_over_http/astats_over_http.la
astats_over_http_astats_over_http_la_SOURCES = astats_over_http/astats_over_http.c
MAKEFILE

ex /root/rpmbuild/SOURCES/trafficserver/plugins/Makefile.am << ED
/stats_over_http/
a
include astats_over_http/Makefile.inc
.
wq
ED
) || die "Failed to patch in astats_over_http"

arch="$(rpm --eval %_arch)"

# build a trafficserver RPM
rm -f /root/rpmbuild/RPMS/${arch}/trafficserver-*.rpm
cd trafficserver

if [[ ${RUN_ATS_UNIT_TESTS} == true ]]; then
  rpmbuild --define "ats_version $ATS_MAJOR_VERSION" -bb ${rpmbuild_openssl} /trafficserver.spec --define 'run_unit_tests 1' || die "Failed to build the ATS RPM."
else
  rpmbuild --define "ats_version $ATS_MAJOR_VERSION" -bb ${rpmbuild_openssl} /trafficserver.spec || die "Failed to build the ATS RPM."
fi

echo "Build completed"

if [[ ! -d /trafficcontrol/dist ]]; then
  mkdir /trafficcontrol/dist
fi

cp /root/rpmbuild/RPMS/${arch}/trafficserver*.rpm /trafficcontrol/dist ||
    die "Failed to copy the ATS RPM to the dist directory"

echo "trafficserver RPM has been copied"
