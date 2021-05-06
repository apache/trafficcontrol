#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
# RPM spec file for Traffic Stats (tm).
#
%define debug_package %{nil}
Name:     trafficcontrol-cache-config
Summary:  Installs Traffic Control cache configuration tools
Version:  %{traffic_control_version}
Release:  %{build_number}
License:  Apache License, Version 2.0
Group:    Applications/Communications
Source0:  trafficcontrol-cache-config-%{version}.tgz
URL:      https://github.com/apache/trafficcontrol/
Vendor:   Apache Software Foundation
Packager: dev at trafficcontrol dot Apache dot org
%{?el6:Requires: git, perl}
%{?el7:Requires: git, perl}
%{?el8:Requires: git, perl}


%description
Installs Traffic Control Cache Configuration utilities. See the `t3c` application.

%prep
tar xvf %{SOURCE0} -C $RPM_SOURCE_DIR


%build
set -o nounset
# copy license
cp "${TC_DIR}/LICENSE" %{_builddir}

ccdir="cache-config"
ccpath="src/github.com/apache/trafficcontrol/${ccdir}/"

# copy atstccfg binary
godir="$ccpath"/atstccfg
( mkdir -p "$godir" && \
	cd "$godir" && \
	cp "$TC_DIR"/"$ccdir"/atstccfg/atstccfg .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy t3c-apply binary
got3cdir="$ccpath"/t3c-apply
( mkdir -p "$got3cdir" && \
	cd "$got3cdir" && \
	cp "$TC_DIR"/"$ccdir"/t3c-apply/t3c-apply .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy to_requester binary
go_toreq_dir="$ccpath"/to_requester
( mkdir -p "$go_toreq_dir" && \
	cd "$go_toreq_dir" && \
	cp "$TC_DIR"/"$ccdir"/to_requester/to_requester .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy t3c-update binary
go_toupd_dir="$ccpath"/t3c-update
( mkdir -p "$go_toupd_dir" && \
	cd "$go_toupd_dir" && \
	cp "$TC_DIR"/"$ccdir"/t3c-update/t3c-update .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy plugin_verifier binary
go_plugin_dir="$ccpath"/plugin_verifier
( mkdir -p "$go_plugin_dir" && \
	cd "$go_plugin_dir" && \
	cp "$TC_DIR"/"$ccdir"/plugin_verifier/plugin_verifier .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

%install
ccdir="cache-config/"
installdir="/usr/bin"

mkdir -p ${RPM_BUILD_ROOT}/"$installdir"
mkdir -p "${RPM_BUILD_ROOT}"/etc/logrotate.d
mkdir -p "${RPM_BUILD_ROOT}"/var/log/trafficcontrol-cache-config

cp -p ${RPM_SOURCE_DIR}/trafficcontrol-cache-config-%{version}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/"$installdir"
cp -p ${RPM_SOURCE_DIR}/trafficcontrol-cache-config-%{version}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/"$installdir"

src=src/github.com/apache/trafficcontrol/cache-config
cp -p ${RPM_SOURCE_DIR}/trafficcontrol-cache-config-%{version}/build/atstccfg.logrotate "${RPM_BUILD_ROOT}"/etc/logrotate.d/atstccfg
touch ${RPM_BUILD_ROOT}/var/log/trafficcontrol-cache-config/atstccfg.log
cp -p "$src"/atstccfg/atstccfg ${RPM_BUILD_ROOT}/"$installdir"

t3csrc=src/github.com/apache/trafficcontrol/"$ccdir"/t3c-apply
cp -p "$t3csrc"/t3c-apply ${RPM_BUILD_ROOT}/"$installdir"

to_req_src=src/github.com/apache/trafficcontrol/"$ccdir"/to_requester
cp -p "$to_req_src"/to_requester ${RPM_BUILD_ROOT}/"$installdir"

to_upd_src=src/github.com/apache/trafficcontrol/"$ccdir"/t3c-update
cp -p "$to_upd_src"/t3c-update ${RPM_BUILD_ROOT}/"$installdir"

plugin_vfy_src=src/github.com/apache/trafficcontrol/"$ccdir"/plugin_verifier
cp -p "$plugin_vfy_src"/plugin_verifier ${RPM_BUILD_ROOT}/"$installdir"

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%license LICENSE
%attr(755, root, root)
/usr/bin/traffic_ops_ort.pl
/usr/bin/supermicro_udev_mapper.pl
/usr/bin/atstccfg
/usr/bin/t3c-apply
/usr/bin/to_requester
/usr/bin/t3c-update
/usr/bin/plugin_verifier

%config(noreplace) /etc/logrotate.d/atstccfg
%config(noreplace) /var/log/trafficcontrol-cache-config/atstccfg.log

%changelog
