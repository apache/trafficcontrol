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
BuildArch:  noarch
Release:    %{build_number}
Summary:    Apache Tomcat Servlet/JSP Engine 9.0+, RI for Servlet 3.1/JSP 2.3 API
License:    Apache Software License
URL:        https://github.com/apache/trafficcontrol/
Source:     %{_sourcedir}/apache-tomcat-%{version}.tar.gz
Requires:   java-11-openjdk-headless

%define tomcat_home /opt/tomcat

%description
This rpm is a minimal install of the Tomcat servlet container version 9.0.
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
rmdir logs
mkdir -p "${RPM_BUILD_ROOT}"/var/log/tomcat
cp -R * ${RPM_BUILD_ROOT}/%{tomcat_home}/
ln -sfT /var/log/tomcat "${RPM_BUILD_ROOT}"%{tomcat_home}/logs

# Remove all webapps.
rm -rf ${RPM_BUILD_ROOT}/%{tomcat_home}/webapps/*

# Remove *.bat
rm -f ${RPM_BUILD_ROOT}/%{tomcat_home}/bin/*.bat

%clean
rm -rf ${RPM_BUILD_ROOT}

# This here takes care of stopping and removing tomcat before installing new files
%pretrans
if [[ -e "/etc/init.d/tomcat" ]]; then
	echo "Disabling and stopping SysV tomcat service..."
	chkconfig tomcat off
	service stop tomcat
fi

if [ -d /opt/apache-tomcat-* ]; then
	echo "Deleting unmanaged Tomcat install from < 2.3 version of Traffic Router"
	rm -rf /opt/apache-tomcat-*
	rm -rf /opt/tomcat
fi

%pre
old_log_dir=/opt/tomcat/logs
new_log_dir=/var/log/tomcat
if [[ -d "$old_log_dir" ]]; then
	if [[ -d "$new_log_dir" ]]; then
		(
		# Include files starting with . in the * glob
		shopt -s dotglob
		mv "$old_log_dir"/* "$new_log_dir" || true
		)
		rmdir "$old_log_dir"
	else
		mv "$old_log_dir" "$new_log_dir"
	fi
	sync
fi

%files
%license LICENSE
%defattr(-,root,root)
%{tomcat_home}
%dir /var/log/tomcat

%post

%preun

%postun

%changelog
* Tue Nov 13 2018 Steve Malenfant <smalenfant@apache.org>
- Remove old installation of tomcat
- Removed systemd service for tomcat
- Requires now leaves java choice to operator
