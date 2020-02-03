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
Name:		traffic_ops_ort
Summary:	Installs ORT script for Traffic Control caches
Version:	%{traffic_control_version}
Release:	%{build_number}
License:	Apache License, Version 2.0
Group:		Applications/Communications
Source0:	traffic_ops_ort-%{version}.tgz
URL:		https://github.com/apache/trafficcontrol/
Vendor:		Apache Software Foundation
Packager:	daniel_kirkwood at Cable dot Comcast dot com
%{?el6:Requires: perl-JSON, perl-libwww-perl, perl-Crypt-SSLeay, perl-Digest-SHA}
%{?el7:Requires: perl-JSON, perl-libwww-perl, perl-Crypt-SSLeay, perl-LWP-Protocol-https, perl-Digest-SHA}


%description
Installs ORT script for Traffic Ops caches

%prep
tar xvf %{SOURCE0} -C $RPM_SOURCE_DIR


%build
# copy atstccfg binary
godir=src/github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg
( mkdir -p "$godir" && \
  cd "$godir" && \
  cp "$TC_DIR"/traffic_ops/ort/atstccfg/atstccfg .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }


%install
mkdir -p ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/opt/ort

src=src/github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg
cp -p "$src"/atstccfg ${RPM_BUILD_ROOT}/opt/ort

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%attr(755, root, root)
/opt/ort/traffic_ops_ort.pl
/opt/ort/supermicro_udev_mapper.pl
/opt/ort/atstccfg

%changelog
