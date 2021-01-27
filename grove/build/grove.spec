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

Summary:      Grove HTTP Caching Proxy
Name:         grove
Version:      %{version}
Release:      %{build_number}
License:      Apache License, Version 2.0
Group:        Base System/System Tools
Prefix:       /usr/sbin/%{name}
Source:       %{_sourcedir}/%{name}-%{version}.tgz
URL:          https://github.com/apache/trafficcontrol/%{name}
Distribution: CentOS Linux
Vendor:       Apache Software Foundation
BuildRoot:    %{buildroot}

# %define PACKAGEDIR %{prefix}

%description
An HTTP Caching Proxy

%prep

%build
set -o nounset
# copy license
cp "${TC_DIR}/LICENSE" %{_builddir}

tar -xvzf %{_sourcedir}/%{name}-%{version}.tgz --directory %{_builddir}

%install
rm -rf %{buildroot}/usr/sbin/%{name}
mkdir -p %{buildroot}/usr/sbin/
cp -p %{name} %{buildroot}/usr/sbin/

rm -rf %{buildroot}/etc/%{name}
mkdir -p -m 777 %{buildroot}/etc/%{name}
cp -p conf/%{name}.cfg %{buildroot}/etc/%{name}

rm -rf %{buildroot}/var/log/%{name}
mkdir -p -m 777 %{buildroot}/var/log/%{name}

mkdir -p -m 777 %{buildroot}/etc/init.d/
cp -p  build/%{name}.init %{buildroot}/etc/init.d/%{name}

mkdir -p -m 777 %{buildroot}/etc/logrotate.d/
cp -p build/%{name}.logrotate %{buildroot}/etc/logrotate.d/%{name}

%clean
echo "cleaning"
rm -r -f %{buildroot}

%files
%license LICENSE
/usr/sbin/%{name}
/var/log/%{name}
%config(noreplace) /etc/%{name}
%config(noreplace) /etc/logrotate.d/%{name}
/etc/init.d/%{name}
