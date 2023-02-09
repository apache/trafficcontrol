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
Name:     traffic-portal
Version:  %{traffic_control_version}
Release:  %{build_number}
Summary:  Traffic Portal v2
Group:    Applications/Communications
License:  Apache License, Version 2.0
URL:      https://github.com/apache/trafficcontrol/
Source:   %{_sourcedir}/traffic-portal-%{traffic_control_version}.tgz
AutoReqProv: no
Requires: nodejs >= 2:16.0.0

%define traffic_portal_home /opt/traffic-portal
%description
Installs Traffic Portal

Built: @BUILT@

%prep
rm -rf $RPM_BUILD_DIR/traffic-portal-%{version}
tar -xzvf $RPM_SOURCE_DIR/traffic-portal-%{version}.tgz

%setup

%build
		npm ci
		npm run build:ssr

%install
		%__mkdir -p ${RPM_BUILD_ROOT}/etc/init.d
		%__mkdir -p ${RPM_BUILD_ROOT}/etc/logrotate.d
		%__mkdir -p ${RPM_BUILD_ROOT}/etc/traffic-portal
		%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_home}/browser
		%__mkdir -p ${RPM_BUILD_ROOT}%{traffic_portal_home}/server
		%__mkdir -p ${RPM_BUILD_ROOT}/var/log/traffic-portal

		%__cp ${RPM_BUILD_DIR}/traffic-portal-%{version}/build/config.json ${RPM_BUILD_ROOT}/etc/traffic-portal/.
		%__cp -r ${RPM_BUILD_DIR}/traffic-portal-%{version}/dist/traffic-portal/* ${RPM_BUILD_ROOT}/%{traffic_portal_home}/.
		%__cp -r ${RPM_BUILD_DIR}/traffic-portal-%{version}/package-lock.json ${RPM_BUILD_ROOT}%{traffic_portal_home}/.
		#%__cp ${RPM_BUILD_DIR}/traffic_portal-%{version}/build/etc/init.d/traffic-portal ${RPM_BUILD_ROOT}/etc/init.d/.


	# creates dynamic json file needed at runtime for traffic portal to display release info
	VERSION=%{version}-%{build_number}
	BUILD_DATE=$(date +'%Y-%m-%d %H:%M')
	VERSION="\"Version\":\"$VERSION\""
	BUILD_DATE="\"Build Date\":\"$BUILD_DATE\""
	JSON_VERSION="{\n$VERSION,\n$BUILD_DATE\n}"
	echo -e $JSON_VERSION > ${RPM_BUILD_ROOT}%{traffic_portal_home}/traffic-portal_release.json

%post
		echo "Successfully installed the traffic-portal assets to " %{traffic_portal_home}
		#%__chmod +x /etc/init.d/traffic-portal
		echo "Successfully installed the 'traffic-portal' service"
		#/sbin/chkconfig traffic-portal on
		echo ""
		echo "Start with 'service traffic-portal start'"

%files
%license LICENSE
%defattr(644,root,root,755)
#%attr(755,root,root) /etc/init.d/traffic-portal
%config(noreplace)/etc/traffic-portal/config.json
%config(noreplace)%{traffic_portal_home}/traffic-portal_release.json
%dir /var/log/traffic-portal
%{traffic_portal_home}/browser
%{traffic_portal_home}/server
%{traffic_portal_home}/package-lock.json

