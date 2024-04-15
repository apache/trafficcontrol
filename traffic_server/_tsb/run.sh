#!/usr/bin/env bash

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
# SPDX-License-Identifier: Apache-2.0

die() {
	{ test -n "$@" && echo "$@"; exit 1; } >&2
}

setowner() {
	own="$(stat -c%u:%g "$1")"
	shift
	chown -R "${own}" "$@"
}
trap 'exit_code=$?; setowner /rpmbuilddir/RPMS/x86_64 /rpmbuilddir/RPMS/x86_64; exit $exit_code' EXIT;

mkdir /opt/build
cp -fa /opt/{src,build}/jansson
cp -fa /opt/{src,build}/cjose

if [ "$1" == "--with_openssl" ]; then
	cp -fa /opt/{src,build}/openssl

	# Build OpenSSL
	(
		cd /opt/build/openssl && \
		./config --prefix=/opt/trafficserver/openssl --openssldir=/opt/trafficserver/openssl zlib && \
		make -j`nproc` && \
		make install_sw
	) || die "Failed to build OpenSSL"
	cjose_openssl='--with-openssl=/opt/trafficserver/openssl'
	rpmbuild_openssl='--with openssl_included'
else
	cjose_openssl=''
	rpmbuild_openssl='--without openssl_included'
fi


(cd /opt/build/jansson && patch -p1 < /opt/src/jansson.pic.patch && autoreconf -i && ./configure --enable-shared=no && make -j`nproc` && make install) || die "Failed to install jansson from source."
(cd /opt/build/cjose && patch -p1 < /opt/src/cjose.pic.patch && autoreconf -i && ./configure --enable-shared=no ${cjose_openssl} && make -j`nproc` && make install) || die "Failed to install cjose from source."

# Patch astats in so that it builds in-tree.
cp -far /opt/src/astats_over_http /rpmbuilddir/SOURCES/src/plugins/astats_over_http
cat > /rpmbuilddir/SOURCES/src/plugins/astats_over_http/Makefile.inc <<MAKEFILE
pkglib_LTLIBRARIES += astats_over_http/astats_over_http.la
astats_over_http_astats_over_http_la_SOURCES = astats_over_http/astats_over_http.c
MAKEFILE
(ed /rpmbuilddir/SOURCES/src/plugins/Makefile.am <<ED
/stats_over_http/
t
s/stats/astats/g
w
ED
) || die "Failed to patch plugins makefile to include astats."

# Patch trafficserver systemd service
# This includes adding udev-settle to wait for disks
(sed -i 's/After=syslog.target network.target/Wants=systemd-udev-settle.service \nAfter=syslog.target network.target systemd-udev-settle.service/g' /rpmbuilddir/SOURCES/src/rc/trafficserver.service.in)
rpmbuild -bb ${rpmbuild_openssl} --define "_topdir /rpmbuilddir" /rpmbuilddir/SPECS/trafficserver.spec || die "Failed to build rpm."
