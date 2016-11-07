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
# RPM spec file for the Traffic Ops v2
#
%define		debug_package %{nil}
Name:		traffic_ops_v2
Version:	%{traffic_control_version}
Release:	%{build_number}
Summary:	Traffic Ops v2
Group:		Applications/Communications
License:	Apache License, Version 2.0
URL:		https://github.com/Comcast/traffic_control/
Source:		%{_sourcedir}/traffic_ops_v2-%{traffic_control_version}.tgz
AutoReqProv: no
Requires: nodejs

%define traffic_ops_v2_home /opt/traffic_ops_v2
%description
Installs Traffic Ops v2

Built: @BUILT@

%prep
rm -rf $RPM_BUILD_DIR/traffic_ops_v2-%{version}
tar -xzvf $RPM_SOURCE_DIR/traffic_ops_v2-%{version}.tgz

%setup

%build
    /usr/bin/npm install
    /usr/bin/bower install
    /usr/bin/grunt dist

%install
    %__mkdir -p ${RPM_BUILD_ROOT}/etc/init.d
    %__mkdir -p ${RPM_BUILD_ROOT}/etc/logrotate.d
    %__mkdir -p ${RPM_BUILD_ROOT}/etc/traffic_ops_v2
    %__mkdir -p ${RPM_BUILD_ROOT}%{traffic_ops_v2_home}/public
    %__mkdir -p ${RPM_BUILD_ROOT}%{traffic_ops_v2_home}/server
    %__mkdir -p ${RPM_BUILD_ROOT}/var/log/traffic_ops_v2

    # creates dynamic json file needed at runtime for traffic ops v2 to display release info
    BUILD_DATE=$(date +'%Y-%m-%d %H:%M:%S')
    VERSION="\"Version\":\"$VERSION\""
    BUILD_NUMBER="\"Build Number\":\"$BUILD_NUMBER\""
    BUILD_DATE="\"Build Date\":\"$BUILD_DATE\""
    JSON_VERSION="{\n$VERSION,\n$BUILD_NUMBER,\n$BUILD_DATE\n}"
    echo -e $JSON_VERSION > ${RPM_BUILD_ROOT}%{traffic_ops_v2_home}/public/traffic_ops_release.json

    %__cp ${RPM_BUILD_DIR}/traffic_ops_v2-%{version}/server/server.js ${RPM_BUILD_ROOT}%{traffic_ops_v2_home}/server/.
    %__cp -r ${RPM_BUILD_DIR}/traffic_ops_v2-%{version}/conf ${RPM_BUILD_ROOT}/etc/traffic_ops_v2/.
    %__cp ${RPM_BUILD_DIR}/traffic_ops_v2-%{version}/build/etc/init.d/traffic_ops_v2 ${RPM_BUILD_ROOT}/etc/init.d/.
    %__cp ${RPM_BUILD_DIR}/traffic_ops_v2-%{version}/build/etc/logrotate.d/traffic_ops_v2 ${RPM_BUILD_ROOT}/etc/logrotate.d/.
    %__cp ${RPM_BUILD_DIR}/traffic_ops_v2-%{version}/build/etc/logrotate.d/traffic_ops_v2-access ${RPM_BUILD_ROOT}/etc/logrotate.d/.
    %__cp -r ${RPM_BUILD_DIR}/traffic_ops_v2-%{version}/app/dist/* ${RPM_BUILD_ROOT}%{traffic_ops_v2_home}/.

%post
    echo "Successfully installed the traffic_ops_v2 assets to " %{traffic_ops_v2_home}
    %__chmod +x %{traffic_ops_v2_home}/node_modules/forever/bin/forever
    %__chmod +x /etc/init.d/traffic_ops_v2
    echo "Successfully installed the 'traffic_ops_v2' service"
    /sbin/chkconfig traffic_ops_v2 on
    echo ""
    echo "Start with 'service traffic_ops_v2 start'"

%files
%defattr(644,root,root,755)
%attr(755,root,root) /etc/init.d/traffic_ops_v2
%attr(755,root,root) %{traffic_ops_v2_home}/node_modules/forever/bin/*
%config(noreplace)/etc/traffic_ops_v2/conf/config.js
%dir /var/log/traffic_ops_v2
/etc/traffic_ops_v2/conf/config-template.js
%{traffic_ops_v2_home}/*
/etc/logrotate.d/traffic_ops_v2
/etc/logrotate.d/traffic_ops_v2-access
/etc/init.d/traffic_ops_v2
