%define debug_package %{nil}

Name:       tomcat
Version:    %{tomcat_version}
Release:    %{build_number}
Summary:    Apache Tomcat Servlet/JSP Engine 8.5+, RI for Servlet 3.1/JSP 2.3 API
License:    Apache Software License
URL:        https://github.com/apache/incubator-trafficcontrol/
Source:     %{_sourcedir}/apache-tomcat-%{version}.tar.gz
Requires:   jdk >= 1.8
Obsoletes:  traffic_router < 2.3

# Set OS version specific variables here
%define startup_script %{_sysconfdir}/systemd/system/tomcat.service

# include the common .spec for tomcat
%include %{_sourcedir}/tomcat.inc
