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
%global commit %(git blame --incremental astats_over_http.c | head -1 | awk '{print substr($1,0,7)}' )
%global no_commits %(git log astats_over_http.c | grep '^commit' | wc -l)
%global _find_debuginfo_dwz_opts %{nil}

Name:       astats_over_http
Version:    1.3.0
Release:    %{no_commits}.%{commit}%{?dist}
Epoch:      434
Summary:    Apache Traffic Server %{name} plugin
Vendor:     Comcast
Group:      Applications/Communications
License:    Apache License, Version 2.0
URL:        https://github.com/apache/trafficcontrol/tree/master/traffic_server/plugins/astats_over_http
BuildRoot:  %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
Requires: trafficserver >= 6011
BuildRequires:  trafficserver >= 6011

%description
Apache Traffic Server plugin

%prep
rm -rf %{name}
git clone https://git-wip-us.apache.org/repos/asf/trafficcontrol.git
cd trafficcontrol
git checkout master
git checkout %{commit} .
cd ..
mv trafficcontrol/traffic_server/plugins/astats_over_http %{name}
# copy license
cp trafficcontrol/LICENSE %{_builddir}
rm -rf trafficcontrol
%setup -D -n %{name} -T

%build
%{install_prefix}/trafficserver/bin/tsxs -v -c %{name}.c -o %{name}.so

%install
mkdir -p $RPM_BUILD_ROOT%{install_prefix}/trafficserver/libexec/trafficserver
DESTDIR=$RPM_BUILD_ROOT %{install_prefix}/trafficserver/bin/tsxs -v -o %{name}.so -i

%clean
rm -rf $RPM_BUILD_ROOT

%post

%postun

%files
%defattr(-,root,root)
%license LICENSE
/opt/trafficserver/libexec/trafficserver/%{name}.so
