%define debug_package %{nil}

Name:       tomcat_tr
Version:    8.5
Release:    28
Summary:    Apache Tomcat Servlet/JSP Engine 8.5, RI for Servlet 3.1/JSP 2.3 API
License:    Apache Software License
URL:        https://github.com/apache/incubator-trafficcontrol/
Source:     %{_sourcedir}/apache-tomcat-%{version}.%{release}.tar.gz
Requires:   jdk >= 1.8

%define tomcat_home /opt/tomcat

%description
Tomcat8.5 is a minimal install of the Tomcat servlet container version 8.5. 
It contains no webapps of its own. To use it create your own CATALINA_BASE
directory and place your application specific webapps and server.xml there. 
You will also need your own systemd unit file for starting your application 
with the correct setting for CATALINA_BASE. 

Built:@BUILT@

%prep
%setup -q -n apache-tomcat-%{version}.%{release}

%build

%install
install -d -m 755 ${RPM_BUILD_ROOT}/%{tomcat_home}/
cp -R * ${RPM_BUILD_ROOT}/%{tomcat_home}/

# Remove all webapps. 
rm -rf ${RPM_BUILD_ROOT}/%{tomcat_home}/webapps/*

# Remove *.bat
rm -f ${RPM_BUILD_ROOT}/%{tomcat_home}/bin/*.bat

# Put logging in /var/log and link back.
rm -rf ${RPM_BUILD_ROOT}/%{tomcat_home}/logs
install -d -m 755 ${RPM_BUILD_ROOT}/var/log/tomcat
cd ${RPM_BUILD_ROOT}/%{tomcat_home}
ln -s /var/log/tomcat logs
cd -

# Drop sysd script
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sysconfdir}/systemd/system
install    -m 755 %_sourcedir/tomcat.service ${RPM_BUILD_ROOT}/%{_sysconfdir}/systemd/system/tomcat.service 

# Drop logrotate script
install -d -m 755 ${RPM_BUILD_ROOT}/%{_sysconfdir}/logrotate.d
install    -m 644 %_sourcedir/tomcat.logrotate ${RPM_BUILD_ROOT}/%{_sysconfdir}/logrotate.d/tomcat

%clean
rm -rf ${RPM_BUILD_ROOT}

%pre

%files
%defattr(-,root,root)
%{tomcat_home}
%dir /var/log/tomcat
/etc/logrotate.d/tomcat
%{_sysconfdir}/systemd/system/tomcat.service 

%post
systemctl daemon-reload
echo "Tomcat for Traffic Router installed successfully."
echo ""
echo "Start with 'sudo systemctl start tomcat'"

%preun

%postun

%changelog
