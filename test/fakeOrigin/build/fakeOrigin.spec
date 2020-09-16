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
# RPM spec file for the test origin fakeOrigin
#

Summary: fakeOrigin CDN Origin
Name: fakeOrigin
Version: %{_version}
Release: %{_release}
Prefix: /usr/sbin/%{name}
Source: %{_sourcedir}/%{name}-%{_version}-%{_release}.tgz
Distribution: CentOS Linux
BuildRoot: %{buildroot}
License: Apache License, Version 2.0
URL: https://github.com/apache/trafficcontrol
Vendor:	Apache Software Foundation
Requires: initscripts

%description
A fake HTTP CDN Origin for testing

%prep

%build
tar -xvzf %{_sourcedir}/%{name}-%{_version}-%{_release}.tgz --directory %{_builddir}

%install
rm -rf %{buildroot}/opt/%{name}
mkdir -p %{buildroot}/opt/%{name}/example
cp -p %{name} %{buildroot}/opt/%{name}
cp -rp example/* %{buildroot}/opt/%{name}/example/

rm -rf %{buildroot}/etc/%{name}
mkdir -p -m 777 %{buildroot}/etc/%{name}
cp -p build/config.json %{buildroot}/etc/%{name}

rm -rf %{buildroot}/etc/logrotate.d/%{name}
mkdir -p -m 777 %{buildroot}/etc/logrotate.d/%{name}
cp -p build/%{name}.logrotate %{buildroot}/etc/logrotate.d/%{name}

rm -rf %{buildroot}/var/log/%{name}
mkdir -p -m 777 %{buildroot}/var/log/%{name}

mkdir -p -m 777 %{buildroot}/etc/init.d/
cp -p  build/%{name}.init %{buildroot}/etc/init.d/%{name}

%clean
echo "cleaning"
rm -r -f %{buildroot}

%files
/opt/%{name}/%{name}
/opt/%{name}/example
/var/log/%{name}
%config(noreplace) /etc/%{name}
%config(noreplace) /etc/logrotate.d/%{name}
/etc/init.d/%{name}
