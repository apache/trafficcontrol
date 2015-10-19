%define debug_package %{nil}
Name:		traffic_ops_ort
Version:	0.54a
Release:	1%{?dist}
Summary:	Installs ORT script for Traffic Control caches
Packager:	mark_torluemke at Cable dot Comcast dot com
Vendor:		Comcast
Group:		Applications/Communications
License:	Apache License, Version 2.0
Requires:	perl-JSON
URL:		https://github.com/Comcast/traffic_control/
Source0:	traffic_ops_ort.tgz


%description
Installs ORT script for Traffic Ops caches

%prep
rm -f $RPM_SOURCE_DIR/traffic_ops_ort.pl
rm -f $RPM_SOURCE_DIR/supermicro_udev_mapper.pl
tar xvf $RPM_SOURCE_DIR/traffic_ops_ort.tgz -C $RPM_SOURCE_DIR


%build


%install
mkdir -p ${RPM_BUILD_ROOT}/opt/ort
cp -r ${RPM_SOURCE_DIR}/traffic_ops_ort.pl ${RPM_BUILD_ROOT}/opt/ort
cp -r ${RPM_SOURCE_DIR}/supermicro_udev_mapper.pl ${RPM_BUILD_ROOT}/opt/ort

%clean
rm -rf ${RPM_BUILD_ROOT}

%post

%files
%attr(755, root, root)
/opt/ort/traffic_ops_ort.pl
/opt/ort/supermicro_udev_mapper.pl

%changelog
