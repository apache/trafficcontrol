Summary: Grove HTTP Caching Proxy
Name: grove
Version: 0.1
Release: 1
License: Apache License, Version 2.0
Group: Base System/System Tools
Prefix: /opt/%{name}
Source: %{_sourcedir}/%{name}-%{version}.tgz
URL: https://github.com/apache/incubator-trafficcontrol/%{name}
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
rm -rf %{buildroot}/opt/%{name}
mkdir -p %{buildroot}/opt/%{name}
cp -p %{name} %{buildroot}/opt/%{name}

rm -rf %{buildroot}/etc/%{name}
mkdir -p -m 777 %{buildroot}/etc/%{name}
cp -p  %{name}.cfg %{buildroot}/etc/%{name}

rm -rf %{buildroot}/var/log/%{name}
mkdir -p -m 777 %{buildroot}/var/log/%{name}

mkdir -p %{buildroot}/usr/lib/systemd/system/
cp -p  %{name}.service %{buildroot}/usr/lib/systemd/system/

%clean
echo "cleaning"
rm -r -f %{buildroot}

%files
/opt/%{name}
/var/log/%{name}
%config(noreplace) /etc/%{name}
/usr/lib/systemd/system/
