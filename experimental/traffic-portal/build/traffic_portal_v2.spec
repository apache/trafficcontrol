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
# RPM spec file for the Traffic Portal
#
%define   debug_package %{nil}
Name:     traffic_portal_v2
Version:  %{traffic_control_version}
Release:  %{build_number}
Summary:  Traffic Portal v2
Group:    Applications/Communications
License:  Apache License, Version 2.0
URL:      https://github.com/apache/trafficcontrol/
Source:   %{_sourcedir}/traffic-portal-%{traffic_control_version}.tgz
AutoReqProv: no
Requires: nodejs >= 2:20.0.0
Requires(pre): /usr/sbin/useradd, /usr/bin/getent

%define traffic_portal_home /opt/traffic-portal
%define traffic_portal_log /var/log/traffic-portal
%define traffic_portal_conf /etc/traffic-portal
%define traffic_portal traffic-portal-%{version}
%define traffic_portal_user trafficportal
%description
Installs Traffic Portal

Built: @BUILT@

%prep
%__rm -rf $RPM_BUILD_DIR/%{traffic_portal}
tar -xzvf $RPM_SOURCE_DIR/%{traffic_portal}.tgz

%pre
/usr/bin/getent group %{traffic_portal_user} || /usr/sbin/groupadd -r %{traffic_portal_user}
/usr/bin/getent passwd %{traffic_portal_user} || /usr/sbin/useradd -r -d %{traffic_portal_home} -s /sbin/nologin %{traffic_portal_user} -g %{traffic_portal_user}

%build
cd ${RPM_BUILD_DIR}/%{traffic_portal}
npm run build:ssr

%install
%__mkdir -p ${RPM_BUILD_ROOT}/etc/init.d
%__mkdir -p ${RPM_BUILD_ROOT}/etc/logrotate.d
%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_conf}
%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_home}/browser
%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_home}/server
%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_home}/node_modules
%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_log}

%__cp ${RPM_BUILD_DIR}/%{traffic_portal}/build/config.json ${RPM_BUILD_ROOT}%{traffic_portal_conf}/.
%__cp -r ${RPM_BUILD_DIR}/%{traffic_portal}/dist/traffic-portal/* ${RPM_BUILD_ROOT}%{traffic_portal_home}/.
%__cp -r ${RPM_BUILD_DIR}/%{traffic_portal}/build/node_modules ${RPM_BUILD_ROOT}%{traffic_portal_home}/.
%__cp ${RPM_BUILD_DIR}/%{traffic_portal}/build/etc/init.d/traffic-portal ${RPM_BUILD_ROOT}/etc/init.d/.
%__cp ${RPM_BUILD_DIR}/%{traffic_portal}/build/etc/logrotate.d/traffic-portal ${RPM_BUILD_ROOT}/etc/logrotate.d/.
%__cp ${RPM_BUILD_DIR}/%{traffic_portal}/LICENSE ${RPM_BUILD_DIR}/.

# creates dynamic json file needed at runtime for traffic portal to display release info
echo "{
	\"date\": \"$(date +'%Y-%m-%d %H:%M')\",
	\"elRelease\": \"%{rhel_vers}\",
	\"hash\": \"%{build_number}\",
	\"version\": \"%{version}\"
}" > ${RPM_BUILD_ROOT}%{traffic_portal_conf}/version.json

%post
%__chmod +x %{traffic_portal_home}/node_modules/pm2/bin/pm2
echo "Successfully installed traffic-portal to " %{traffic_portal_home}
/sbin/chkconfig traffic-portal on
echo ""
echo "Start with 'systemctl start traffic-portal' or by running '%{traffic_portal_conf}/traffic-portal'"


%files
%license LICENSE
%defattr(644,%{traffic_portal_user},%{traffic_portal_user},755)
%attr(755,%{traffic_portal_user},%{traffic_portal_user}) /etc/init.d/traffic-portal
%config(noreplace)%{traffic_portal_conf}/config.json
%config(noreplace)%{traffic_portal_conf}/version.json
%config(noreplace)/etc/logrotate.d/traffic-portal
%dir %{traffic_portal_log}
%{traffic_portal_home}
