%define debug_package %{nil}
Name:		traffic_ops_ort
Summary:	Installs ORT script for Traffic Control caches
Version:	%{traffic_control_version}
Release:	%{build_number}
License:	Apache License, Version 2.0
Group:		Applications/Communications
Source0:	traffic_ops_ort-%{version}.tgz
URL:		https://github.com/Comcast/traffic_control/
Vendor:		Comcast
Packager:	daniel_kirkwood at Cable dot Comcast dot com
Requires:	perl-JSON


%description
Installs ORT script for Traffic Ops caches

%prep
tar xvf %{SOURCE0} -C $RPM_SOURCE_DIR


%build


%install
mkdir -p ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/opt/ort
cp -p ${RPM_SOURCE_DIR}/traffic_ops_ort-%{version}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/opt/ort

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%attr(755, root, root)
/opt/ort/traffic_ops_ort.pl
/opt/ort/supermicro_udev_mapper.pl

%changelog
