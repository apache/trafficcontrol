Summary: Grove HTTP Caching Proxy Traffic Control config generator
Name: grovetccfg
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
A Traffic Control config generator for the Grove HTTP Caching Proxy

%prep

%build
tar -xvzf %{_sourcedir}/%{name}-%{version}.tgz --directory %{_builddir}

%install
rm -rf %{buildroot}/opt/%{name}
mkdir -p %{buildroot}/opt/%{name}/
cp -p %{name} %{buildroot}/opt/%{name}/

%clean
echo "cleaning"
rm -r -f %{buildroot}

%files
/opt/%{name}/%{name}
