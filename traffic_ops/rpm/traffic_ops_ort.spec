#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
%define debug_package %{nil}
Name:		traffic_ops_ort
Version:	0.50c
Release:	1%{?dist}
Summary:	Installs ORT script for Traffic Control caches
Packager:	mark_torluemke at Cable dot Comcast dot com
Vendor:		Comcast
Group:		Applications/Communications
License:	GNU LGPL Version 2.1
Requires:	perl-JSON
URL:		https://evc.io.comcast.net/video/cdneng/configs/twelve_monkeys/tags/current/
Source0:	traffic_ops_ort.tgz


%description
Installs ORT script for Traffic Ops caches

%prep
rm -f $RPM_SOURCE_DIR/traffic_ops_ort.pl
rm -f $RPM_SOURCE_DIR/supermicro_udev_mapper.pl
tar xvf $RPM_SOURCE_DIR/traffic_ops_ort.tgz -C $RPM_SOURCE_DIR


%build


%install
mkdir -p ${RPM_BUILD_ROOT}/opt/ort
cp -r ${RPM_SOURCE_DIR}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/opt/ort
cp -r ${RPM_SOURCE_DIR}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/opt/ort

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%attr(755, root, root)
/opt/ort/traffic_ops_ort.pl
/opt/ort/supermicro_udev_mapper.pl

%changelog
* Thu Jan 29 2015 Mark Torluemke <mark_torluemke(at)cable.comcast.com>
- 0.42b Initial RPM build based on git hash 688ac64
- 0.43d RPM rebuild based on git hash c7e829e
- 0.44a RPM rebuild based on git hash 9171d76
- 0.45a RPM rebuild based on git hash 8fab555 - Thu Feb 19 2015
- 0.45b RPM rebuild based on git hash f12afeb - Thu Feb 19 2015
- 0.46a RPM rebuild based on git hash 43b6535 - Thu Feb 19 2015
- 0.46b RPM rebuild based on git hash ea8e57d - Thu Mar 04 2015
- 0.47a RPM rebuild based on git hash 9c0f864 - Mon Mar 09 2015
- 0.48a RPM rebuild based on git hash a4f0615 - Mon Mar 09 2015
- 0.49a RPM rebuild based on git hash 60754d2 - Tue May 12 2015
- 0.50c RPM rebuild based on git hash d61a85f - Thu Jun 11 2015
