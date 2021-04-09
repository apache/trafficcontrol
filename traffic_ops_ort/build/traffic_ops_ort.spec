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
Name:     traffic_ops_ort
Summary:  Installs ORT script for Traffic Control caches
Version:  %{traffic_control_version}
Release:  %{build_number}
License:  Apache License, Version 2.0
Group:    Applications/Communications
Source0:  traffic_ops_ort-%{version}.tgz
URL:      https://github.com/apache/trafficcontrol/
Vendor:   Apache Software Foundation
Packager: dev at trafficcontrol dot Apache dot org
%{?el6:Requires: git, perl-JSON, perl-libwww-perl, perl-Crypt-SSLeay, perl-Digest-SHA}
%{?el7:Requires: git, perl-JSON, perl-libwww-perl, perl-Crypt-SSLeay, perl-LWP-Protocol-https, perl-Digest-SHA}
%{?el8:Requires: git, perl-JSON, perl-libwww-perl, perl-Net-SSLeay, perl-LWP-Protocol-https, perl-Digest-SHA}


%description
Installs ORT script for Traffic Ops caches

%prep
tar xvf %{SOURCE0} -C $RPM_SOURCE_DIR


%build
set -o nounset
# copy license
cp "${TC_DIR}/LICENSE" %{_builddir}

# copy atstccfg binary
godir=src/github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg
( mkdir -p "$godir" && \
	cd "$godir" && \
	cp "$TC_DIR"/traffic_ops_ort/atstccfg/atstccfg .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy t3c binary
got3cdir=src/github.com/apache/trafficcontrol/traffic_ops_ort/t3c
( mkdir -p "$got3cdir" && \
	cd "$got3cdir" && \
	cp "$TC_DIR"/traffic_ops_ort/t3c/t3c .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy to_requester binary
go_toreq_dir=src/github.com/apache/trafficcontrol/traffic_ops_ort/to_requester
( mkdir -p "$go_toreq_dir" && \
	cd "$go_toreq_dir" && \
	cp "$TC_DIR"/traffic_ops_ort/to_requester/to_requester .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy to_updater binary
go_toupd_dir=src/github.com/apache/trafficcontrol/traffic_ops_ort/to_updater
( mkdir -p "$go_toupd_dir" && \
	cd "$go_toupd_dir" && \
	cp "$TC_DIR"/traffic_ops_ort/to_updater/to_updater .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy plugin_verifier binary
go_plugin_dir=src/github.com/apache/trafficcontrol/traffic_ops_ort/plugin_verifier
( mkdir -p "$go_plugin_dir" && \
	cd "$go_plugin_dir" && \
	cp "$TC_DIR"/traffic_ops_ort/plugin_verifier/plugin_verifier .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

%install
mkdir -p ${RPM_BUILD_ROOT}/opt/ort
mkdir -p "${RPM_BUILD_ROOT}"/etc/logrotate.d
mkdir -p "${RPM_BUILD_ROOT}"/var/log/ort

cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/opt/ort

src=src/github.com/apache/trafficcontrol/traffic_ops_ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/build/atstccfg.logrotate "${RPM_BUILD_ROOT}"/etc/logrotate.d/atstccfg
touch ${RPM_BUILD_ROOT}/var/log/ort/atstccfg.log
cp -p "$src"/atstccfg/atstccfg ${RPM_BUILD_ROOT}/opt/ort

t3csrc=src/github.com/apache/trafficcontrol/traffic_ops_ort/t3c
cp -p "$t3csrc"/t3c ${RPM_BUILD_ROOT}/opt/ort

to_req_src=src/github.com/apache/trafficcontrol/traffic_ops_ort/to_requester
cp -p "$to_req_src"/to_requester ${RPM_BUILD_ROOT}/opt/ort

to_upd_src=src/github.com/apache/trafficcontrol/traffic_ops_ort/to_updater
cp -p "$to_upd_src"/to_updater ${RPM_BUILD_ROOT}/opt/ort

plugin_vfy_src=src/github.com/apache/trafficcontrol/traffic_ops_ort/plugin_verifier
cp -p "$plugin_vfy_src"/plugin_verifier ${RPM_BUILD_ROOT}/opt/ort

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%license LICENSE
%attr(755, root, root)
/opt/ort/traffic_ops_ort.pl
/opt/ort/supermicro_udev_mapper.pl
/opt/ort/atstccfg
/opt/ort/t3c
/opt/ort/to_requester
/opt/ort/to_updater
/opt/ort/plugin_verifier

%config(noreplace) /etc/logrotate.d/atstccfg
%config(noreplace) /var/log/ort/atstccfg.log

%changelog
