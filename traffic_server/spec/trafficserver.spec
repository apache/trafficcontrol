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

%global install_prefix "/opt"

Name:		trafficserver
Version:	%{traffic_control_version}
Release:	%{build_number}
Summary:	Apache Traffic Server
Vendor:		Apache Traffic Control
Group:		Applications/Communications
License:	Apache License, Version 2.0
URL:		http://trafficserver.apache.org/
Source0:        %{name}-%{version}.tgz
BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
Requires:	tcl, hwloc, pcre, openssl, libcap
BuildRequires:	autoconf, automake, libtool, gcc-c++, glibc-devel, openssl-devel, expat-devel, pcre, libcap-devel, pcre-devel, perl-ExtUtils-MakeMaker, tcl-devel, hwloc-devel

%description
Apache Traffic Server built via Traffic Control Docker Build Process

%prep
%setup
autoreconf -vfi

%build
./configure --prefix=%{install_prefix}/%{name} --with-user=ats --with-group=ats --with-build-number=%{release} --enable-experimental-plugins
make %{?_smp_mflags}

%install
make DESTDIR=$RPM_BUILD_ROOT install

mkdir -p $RPM_BUILD_ROOT%{install_prefix}/trafficserver/etc/trafficserver/snapshots
mkdir -p $RPM_BUILD_ROOT/etc/init.d
cp $RPM_BUILD_DIR/%{name}-%{version}/rc/trafficserver $RPM_BUILD_ROOT/etc/init.d/

%clean
rm -rf $RPM_BUILD_ROOT

%pre
getent group ats >/dev/null || groupadd -r ats
getent passwd ats >/dev/null || \
useradd -r  -g ats -d / -s /sbin/nologin \
    -c "Apache Traffic Server" ats &>/dev/null

%post
chkconfig --add %{name}

%preun
/etc/init.d/%{name} stop

# if 0 uninstall, if 1 upgrade
if [ "$1" = "0" ]; then
	chkconfig --del %{name}
fi

%postun
# Helpful in understanding order of operations in relation to install/uninstall/upgrade:
#     http://www.ibm.com/developerworks/library/l-rpm2/
# if 0 uninstall, if 1 upgrade
if [ "$1" = "0" ]; then
	id ats &>/dev/null && /usr/sbin/userdel ats
fi

%files
%defattr(-,root,root)
%attr(755,-,-) /etc/init.d/trafficserver
%dir /opt/trafficserver
/opt/trafficserver/bin
/opt/trafficserver/include
/opt/trafficserver/lib
/opt/trafficserver/libexec
/opt/trafficserver/share
%dir /opt/trafficserver/var
%attr(-,ats,ats) /opt/trafficserver/var/trafficserver
%dir /opt/trafficserver/var/log
%attr(-,ats,ats) /opt/trafficserver/var/log/trafficserver
%dir /opt/trafficserver/etc
%attr(-,ats,ats) %dir /opt/trafficserver/etc/trafficserver
%attr(-,ats,ats) %dir /opt/trafficserver/etc/trafficserver/snapshots
/opt/trafficserver/etc/trafficserver/body_factory
/opt/trafficserver/etc/trafficserver/trafficserver-release
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/cache.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/cluster.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/congestion.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/hosting.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/ip_allow.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/log_hosts.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/parent.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/plugin.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/records.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/remap.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/socks.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/splitdns.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/ssl_multicert.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/storage.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/vaddrs.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/volume.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/metrics.config
%config(noreplace) %attr(644,ats,ats) /opt/trafficserver/etc/trafficserver/logging.config
