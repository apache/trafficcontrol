%global install_prefix "/opt"

Name:		astats_over_http
Version:	1.1.0
Release:	1%{?dist}
Summary:	Apache Traffic Server %{name} plugin
Vendor:		Comcast
Group:		Applications/Communications
License:	Apache License, Version 2.0
URL:		https://github.com/Comcast/traffic_control/tree/master/traffic_server/plugins/astats_over_http
Source0:	%{name}.tar.gz
BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
Requires:	trafficserver = 5.2.0
BuildRequires:	trafficserver = 5.2.0

%description
Apache Traffic Server plugin

%prep
%setup -n %{name}

%build
%{install_prefix}/trafficserver/bin/tsxs -v -c %{name}.c -o %{name}.so

%install
mkdir -p $RPM_BUILD_ROOT%{install_prefix}/trafficserver/libexec/trafficserver
DESTDIR=$RPM_BUILD_ROOT %{install_prefix}/trafficserver/bin/tsxs -v -o %{name}.so -i

%clean
rm -rf $RPM_BUILD_ROOT

%post

%postun

%files
%defattr(-,root,root)
/opt/trafficserver/libexec/trafficserver/%{name}.so
