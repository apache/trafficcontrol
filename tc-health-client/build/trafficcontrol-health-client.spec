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
# RPM spec file for Trafficcontrol health client.
#
%define debug_package %{nil}
Name:     trafficcontrol-health-client
Summary:  Installs the Traffic Control cache health client tool
Version:  %{traffic_control_version}
Release:  %{build_number}
License:  Apache License, Version 2.0
Group:    Applications/Communications
Source0:  trafficcontrol-health-client-%{version}.tgz
URL:      https://github.com/apache/trafficcontrol/
Vendor:   Apache Software Foundation
Packager: dev at trafficcontrol dot Apache dot org
Requires: git 

%description
Installs Traffic Control cache health client. See the `tc-health-client` application.

%prep
tar xvf %{SOURCE0} -C $RPM_SOURCE_DIR


%build
set -o nounset
# copy license
cp ${TC_DIR}/LICENSE %{_builddir}

hcdir="tc-health-client"
hcpath="src/github.com/apache/trafficcontrol/${hcdir}/"
was_active="unk"

%pre

set -x
echo "###### in pre section ######"

if [[ -f /etc/trafficcontrol/tc-health-client.json ]]; then
  active=`systemctl status tc-health-client | awk '/Active:/ {print $2}'`
  if [[ ${active} == "active" ]]; then
    systemctl stop tc-health-client
    touch /run/tc-health-client.pid
  fi
fi

%install

set -x
echo "###### in install section ######"

hcdir="tc-health-client"
installdir="/usr/bin"
mandir="/usr/share/man"
man1dir="man1"

mkdir -p ${RPM_BUILD_ROOT}/${installdir}
mkdir -p ${RPM_BUILD_ROOT}/etc/logrotate.d
mkdir -p ${RPM_BUILD_ROOT}/${mandir}/${man1dir}
mkdir -p ${RPM_BUILD_ROOT}/etc/trafficcontrol
mkdir -p ${RPM_BUILD_ROOT}/usr/lib/systemd/system

cp -p ${RPM_SOURCE_DIR}/trafficcontrol-health-client-%{version}/tc-health-client.logrotate ${RPM_BUILD_ROOT}/etc/logrotate.d/tc-health-client.logrotate

src="trafficcontrol-health-client-%{version}"
cp -p ${RPM_SOURCE_DIR}/${src}/tc-health-client ${RPM_BUILD_ROOT}/${installdir}
cp -p ${RPM_SOURCE_DIR}/${src}/tc-health-client.sample.json ${RPM_BUILD_ROOT}/etc/trafficcontrol
cp -p ${RPM_SOURCE_DIR}/${src}/tc-health-client.logrotate ${RPM_BUILD_ROOT}/etc/logrotate.d
cp -p ${RPM_SOURCE_DIR}/${src}/tc-health-client.service ${RPM_BUILD_ROOT}/usr/lib/systemd/system
gzip -c -9 ${RPM_SOURCE_DIR}/${src}/tc-health-client.1 > ${RPM_BUILD_ROOT}/${mandir}/${man1dir}/tc-health-client.1.gz

ls ${RPM_BUILD_ROOT}/${mandir}/${man1dir}/

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

set -x
echo "###### in post section ######"

# we want all the cache logs under /var/log/trafficcontrol
if [[ ! -d /var/log/trafficcontrol ]]; then
  mkdir -p /var/log/trafficcontrol
  touch /var/log/trafficcontrol/tc-health-client.log
fi

# make sure the service unit file is loaded
systemctl daemon-reload

if [[ -f /run/tc-health-client.pid ]]; then
  systemctl enable tc-health-client
  systemctl start tc-health-client
fi

# update mandb to put man pages in the whatis database, so apps like 'whatis' and 'apropos' get the new pages
mandb_out="$(mandb 2>&1)"
mandb_ret=$?
if [ $mandb_ret -eq 0 ]; then
	printf "%s\n" "Updated mandb"
else
	printf "Failed to update mandb: code %s\n%s\n" ${mandb_ret} ${mandb_out}
fi

%preun

set -x
echo "###### in preun section ######"

if [[ -f /etc/trafficcontrol/tc-health-client.json ]]; then
  active=`systemctl status tc-health-client | awk '/Active:/ {print $2}'`
  if [[ ${active} == "active" ]]; then
    systemctl stop tc-health-client
    touch /run/tc-health-client.pid
  fi
fi

%postun

set -x
echo "###### in postun section ######"

# update whatis database, to remove tc-health-client data
mandb_out="$(mandb 2>&1)"
mandb_ret=$?
if [ $mandb_ret -eq 0 ]; then
	printf "%s\n" "Updated mandb"
else
	printf "Failed to update mandb: code %s\n%s\n" ${mandb_ret} ${mandb_out}
fi

%files
%license LICENSE
%attr(755, root, root)
/usr/bin/tc-health-client
/usr/share/man/man1/tc-health-client.1.gz
/etc/trafficcontrol/tc-health-client.sample.json
/etc/logrotate.d/tc-health-client.logrotate
/usr/lib/systemd/system/tc-health-client.service

%changelog
