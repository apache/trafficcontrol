Summary: Grove HTTP Caching Proxy
Name: grove
Version: %{version}
Release: 1
License: Apache License, Version 2.0
Group: Base System/System Tools
Prefix: /usr/sbin/%{name}
Source: %{_sourcedir}/%{name}-%{version}.tgz
URL: https://github.com/apache/incubator-trafficcontrol%{name}
Distribution: CentOS Linux
Vendor: Apache Software Foundation
BuildRoot: %{buildroot}

# %define PACKAGEDIR %{prefix}

%description
An HTTP Caching Proxy

%prep

%build
tar -xvzf %{_sourcedir}/%{name}-%{version}.tgz --directory %{_builddir}

%install
rm -rf %{buildroot}/usr/sbin/%{name}
mkdir -p %{buildroot}/usr/sbin/
cp -p %{name} %{buildroot}/usr/sbin/

rm -rf %{buildroot}/etc/%{name}
mkdir -p -m 777 %{buildroot}/etc/%{name}
cp -p conf/%{name}.cfg %{buildroot}/etc/%{name}

rm -rf %{buildroot}/var/log/%{name}
mkdir -p -m 777 %{buildroot}/var/log/%{name}

mkdir -p -m 777 %{buildroot}/etc/init.d/
cp -p  build/%{name}.init %{buildroot}/etc/init.d/%{name}

%clean
echo "cleaning"
rm -r -f %{buildroot}

%files
/usr/sbin/%{name}
/var/log/%{name}
%config(noreplace) /etc/%{name}
/etc/init.d/%{name}
