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

%define debug_package %{nil}

Name:       tomcat
Version:    %{tomcat_version}
Release:    %{build_number}
Summary:    Apache Tomcat Servlet/JSP Engine 8.5+, RI for Servlet 3.1/JSP 2.3 API
License:    Apache Software License
URL:        https://github.com/apache/incubator-trafficcontrol/
Source:     %{_sourcedir}/apache-tomcat-%{version}.tar.gz
Requires:   jdk >= 1.8

%define startup_script %{_sysconfdir}/systemd/system/tomcat.service
%define tomcat_home /opt/tomcat

%description
This rpm is a minimal install of the Tomcat servlet container version 8.5.
It gets installed to /opt/tomcat and contains no webapps of its own.
To use it create your own CATALINA_BASE directory and place your application
specific webapps and server.xml there.
You will also need your own systemd unit file for starting your application
with the correct setting for CATALINA_BASE.

Built:@BUILT@

%prep
%setup -q -n apache-tomcat-%{version}

%build

%install
install -d -m 755 ${RPM_BUILD_ROOT}/%{tomcat_home}/
cp -R * ${RPM_BUILD_ROOT}/%{tomcat_home}/

# Remove all webapps.
rm -rf ${RPM_BUILD_ROOT}/%{tomcat_home}/webapps/*

# Remove *.bat
rm -f ${RPM_BUILD_ROOT}/%{tomcat_home}/bin/*.bat

# install sysd script
install -d -m 755 ${RPM_BUILD_ROOT}%{_sysconfdir}/systemd/system
install    -m 755 %_sourcedir/tomcat.service ${RPM_BUILD_ROOT}%{startup_script}

%clean
rm -rf ${RPM_BUILD_ROOT}

%pre
if [[ -e "/etc/init.d/tomcat" ]]; then
  echo "Disabling tomcat service..."
  chkconfig tomcat off
fi

if [ -d /opt/apache-tomcat-* ]; then
  echo "Deleting unmanaged Tomcat install from < 2.3 version of Traffic Router"
  rm -rf /opt/apache-tomcat-*
  rm -rf /opt/tomcat
fi

%files
%defattr(-,root,root)
%{tomcat_home}
%{startup_script}

%post
systemctl daemon-reload

echo "Tomcat for Traffic Router installed successfully."
echo ""
echo "Start with 'systemctl start traffic_router'"

%preun

%postun

%changelog
