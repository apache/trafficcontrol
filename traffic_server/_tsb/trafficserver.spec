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
# SPDX-License-Identifier: Apache-2.0

%global src %{_topdir}/SOURCES/src
%global git_args --git-dir="%{src}/.git" --work-tree="%{src}"
%global tag %(git %{git_args} describe --long --tags --match='*[0-9.][0-9.][0-9.]' |      sed 's/^\\\(.*\\\)-\\\([0-9]\\\+\\\)-g\\\([0-9a-f]\\\+\\\)$/\\\1/' | sed 's/-/_/')
%global distance %(git %{git_args} describe --long --tags --match='*[0-9.][0-9.][0-9.]' | sed 's/^\\\(.*\\\)-\\\([0-9]\\\+\\\)-g\\\([0-9a-f]\\\+\\\)$/\\\2/')
%global commit %(git %{git_args} describe --long --tags --match='*[0-9.][0-9.][0-9.]' |   sed 's/^\\\(.*\\\)-\\\([0-9]\\\+\\\)-g\\\([0-9a-f]\\\+\\\)$/\\\3/')
%global git_serial %(git %{git_args} rev-list HEAD | wc -l)
%global install_prefix "/opt"
%global api_stats "4096"
%global _find_debuginfo_dwz_opts %{nil}
%{!?_with_openssl_included: %{!?_without_openssl_included: %define _without_openssl_included --without-openssl_included}}
%{?_with_openssl_included: %{?_without_openssl_included: %{error: both _with_openssl_included and _without_openssl_included}}}
%{!?_with_openssl_included: %{!?_without_openssl_included: %{error: neither _with_openssl_included nor _without_openssl_included}}}
%{?_without_openssl_included:BuildRequires: openssl-devel}

Name:		trafficserver
Version:	%{tag}
Epoch:		%{git_serial}
Release:	%{distance}.%{commit}%{?dist}
Summary:	Apache Traffic Server
Vendor:		Apache
Group:		Applications/Communications
License:	Apache License, Version 2.0
URL:		https://github.com/apache/trafficserver
BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
Requires:	tcl, hwloc, pcre, libcap
BuildRequires:	autoconf, automake, libtool, gcc-c++, glibc-devel, expat-devel, pcre, libcap-devel, pcre-devel, perl-ExtUtils-MakeMaker, tcl-devel, hwloc-devel
Source: src

%description
Apache Traffic Server with Apache Traffic Control modifications and environment specific modifications

%prep
%setup -c -T
cp -far %{src}/. .
cp -far %{src}/../traffic_server_jemalloc ..
cp -far %{src}/../trafficserver.env ..
autoreconf -vfi

%build
%if %{?_with_openssl_included:1}%{!?_with_openssl_included:0}
./configure --with-openssl=/opt/trafficserver/openssl --prefix=%{install_prefix}/%{name} --with-user=ats --with-group=ats --with-build-number=%{release} --enable-experimental-plugins --with-max-api-stats=%{api_stats} --disable-unwind
%else
./configure --prefix=%{install_prefix}/%{name} --with-user=ats --with-group=ats --with-build-number=%{release} --enable-experimental-plugins --with-max-api-stats=%{api_stats} --disable-unwind
%endif
make %{?_smp_mflags}
%if %{?_with_openssl_included:1}%{!?_with_openssl_included:0}
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/trafficserver/openssl/lib:/usr/local/lib
%else
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
%endif
make %{?_smp_mflags} check || ( cat ./test-suite.log; exit 1 )

%install
make DESTDIR=$RPM_BUILD_ROOT install

mkdir -p $RPM_BUILD_ROOT/opt/trafficserver/etc/trafficserver/snapshots
mkdir -p $RPM_BUILD_ROOT/usr/lib/systemd/system
mkdir -p $RPM_BUILD_ROOT/etc/sysconfig
cp rc/trafficserver.service $RPM_BUILD_ROOT/usr/lib/systemd/system/
cp ../traffic_server_jemalloc $RPM_BUILD_ROOT/opt/trafficserver/bin/
touch $RPM_BUILD_ROOT/etc/sysconfig/trafficserver
cp ../trafficserver.env $RPM_BUILD_ROOT/etc/sysconfig/trafficserver
mkdir -p $RPM_BUILD_ROOT/var/log/trafficserver

%if %{?_with_openssl_included:1}%{!?_with_openssl_included:0}
mkdir -p $RPM_BUILD_ROOT/opt/trafficserver/openssl
cp -r /opt/trafficserver/openssl/lib $RPM_BUILD_ROOT/opt/trafficserver/openssl/lib
%endif

%clean
rm -rf $RPM_BUILD_ROOT

%pre
id ats &>/dev/null || /usr/sbin/useradd -u 176 -r ats -s /sbin/nologin -d /

%post
/bin/systemctl daemon-reload
/bin/systemctl enable trafficserver

%preun
/bin/systemctl stop trafficserver

# if 0 uninstall, if 1 upgrade
if [ "$1" = "0" ]; then
	/bin/systemctl disable trafficserver
fi

%postun
# Helpful in understanding order of operations in relation to install/uninstall/upgrade:
#     https://fedoraproject.org/wiki/Packaging:Scriptlets

# Always do this because the service file may have been updated.
/bin/systemctl daemon-reload

# if 0 uninstall, if 1 upgrade
if [ "$1" = "0" ]; then
	id ats &>/dev/null && /usr/sbin/userdel ats
fi

%files
%license LICENSE
%defattr(-,root,root)
%attr(644,-,-) /usr/lib/systemd/system/trafficserver.service
%attr(644,-,-) /etc/sysconfig/trafficserver
%dir /opt/trafficserver
%if %{?_with_openssl_included:1}%{!?_with_openssl_included:0}
/opt/trafficserver/openssl
%endif
/opt/trafficserver/bin
/opt/trafficserver/include
/opt/trafficserver/lib
/opt/trafficserver/libexec
/opt/trafficserver/share
%dir /opt/trafficserver/var
%attr(-,ats,ats) /opt/trafficserver/var/trafficserver
%dir /var/log/trafficserver
%attr(-,ats,ats) /var/log/trafficserver
%dir /opt/trafficserver/etc
%attr(-,ats,ats) %dir /opt/trafficserver/etc/trafficserver
%attr(-,ats,ats) %dir /opt/trafficserver/etc/trafficserver/snapshots
/opt/trafficserver/etc/trafficserver/body_factory
/opt/trafficserver/etc/trafficserver/trafficserver-release
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/cache.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/hosting.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/logging.yaml
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/parent.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/plugin.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/records.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/remap.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/socks.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/splitdns.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/ssl_multicert.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/storage.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/volume.config
%if "%{tag}" >= "9"
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/ip_allow.yaml
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/sni.yaml
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/strategies.yaml
%else
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/ip_allow.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/ssl_server_name.yaml
%endif

%changelog
* Thu Mar 7 2024 The Anh Nguyen <ntheanh201(at)gmail.com>
- Modified to support ATS 9.2.x
* Wed Mar 10 2021 Jonathan Gray <jhg03a(at)apache.org>
- Modified to support stop bundling openssl with ats
* Wed Aug 26 2020 Chris Lemmons <alficles(at)gmail.com>
- Updated to incorporate new tooling and Apache Traffic Control patches
* Wed Jun 8 2016 John Rushford <john_rushford(at)cable.comcast.com>
- Added tools/rc_admin.pl to complete rpm tasks under both Enterprise Linux 6 or 7 using either chkconfig or systemd commands.
- Modified this spec file to use rc_admin.pl
* Wed Aug 7 2013 Jeff Elsloo <jeffrey_elsloo(at)cable.comcast.com>
- Modified to support building 3.3.x
- Modified to support upgrades
* Sun Aug 12 2012 John Benton <john_benton(at)cable.comcast.com>
- Initial RPM build based on SVN version 2376
- Rev for ATS 3.2.0 based on SVN version 2470
- Rev for ATS 3.2.0 based on SVN version 2555
- Rev for ATS 3.2.0 based on SVN version 4812
