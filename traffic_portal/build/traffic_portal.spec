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
#
# RPM spec file for the Traffic Portal
#
%define debug_package %{nil}
Name:		traffic_portal
Version:        %{traffic_control_version}
Release:        %{build_number}
Summary:        Traffic Portal
Group:		Applications/Communications
License:	Apache License, Version 2.0
URL:		https://github.com/Comcast/traffic_control/
Source:		%{_sourcedir}/traffic_portal-%{traffic_control_version}.tgz
AutoReqProv: no
Requires: nodejs

%define traffic_portal_home /opt/traffic_portal
%description
Installs Traffic Portal

Built: @BUILT@

%prep

%setup

%install
    if [ -d $RPM_BUILD_ROOT ]; then
	    %__rm -rf $RPM_BUILD_ROOT
    fi

    if [ ! -d $RPM_BUILD_ROOT ]; then
        %__mkdir -p $RPM_BUILD_ROOT
    fi

    %__cp -R $RPM_BUILD_DIR/traffic_portal-%{version}/* $RPM_BUILD_ROOT

%post
    echo "Successfully installed the traffic_portal assets to " %{traffic_portal_home}
    %__mkdir -p /var/log/traffic_portal
    %__chmod +x %{traffic_portal_home}/node_modules/forever/bin/forever
    %__chmod +x /etc/init.d/traffic_portal
    echo "Successfully installed the 'traffic_portal' service"
    /sbin/chkconfig traffic_portal on
    echo ""
    echo "Start with 'service traffic_portal start'"

%files
%defattr(644,root,root,755)
%config(noreplace)/etc/traffic_portal/conf/config.js
/etc/traffic_portal/conf/config-template.js
%{traffic_portal_home}/*
%{traffic_portal_home}/server/server.js
/etc/logrotate.d/traffic_portal
/etc/logrotate.d/traffic_portal-access
/etc/init.d/traffic_portal
