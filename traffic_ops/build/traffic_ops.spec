#
#
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
#
#
# RPM spec file for Traffic Ops (tm).
#

%define TRAFFIC_OPS_USER trafops
%define TRAFFIC_OPS_GROUP trafops
%define TRAFFIC_OPS_LOG_DIR /var/log/traffic_ops
%define TRAFFIC_OPS_ROOT_CERTIFICATES_DIR /var/log/traffic_ops
%define debug_package %{nil}

Summary:          Traffic Ops
Name:             traffic_ops
Version:          %{traffic_control_version}
Release:          %{build_number}
License:          Apache License, Version 2.0
Group:            Base System/System Tools
Prefix:           /opt/traffic_ops
Source:           %{_sourcedir}/traffic_ops-%{version}.tgz
URL:              https://github.com/apache/trafficcontrol/
Vendor:           Apache Software Foundation
AutoReqProv:      no
Requires:         cpanminus, expat-devel, libcurl, libpcap-devel, mkisofs, tar
Requires:         openssl-devel, perl, perl-core, perl-DBD-Pg, perl-DBI, perl-Digest-SHA1
Requires:         libidn-devel, libcurl-devel, libcap
Requires:         postgresql13 >= 13.2
Requires:         perl-JSON, perl-libwww-perl, perl-Test-CPAN-Meta, perl-WWW-Curl, perl-TermReadKey, perl-Crypt-ScryptKDF
Requires:         python3
Requires(pre):    /usr/sbin/useradd, /usr/bin/getent
Requires(postun): /usr/sbin/userdel

%define PACKAGEDIR %{prefix}

%description
Traffic Ops is the tool for administration (configuration and monitoring) of all components in a Traffic Control CDN.

This package provides Traffic Ops with the following plugins:
%{getenv:PLUGINS}

Built: %(date) by %{getenv: USER}

%prep

%setup

%build

set -o nounset
# copy LICENSE
cp "${TC_DIR}/LICENSE" %{_builddir}
# avoid detecting LICENSE as an unpackaged RPM file
rm LICENSE

# copy traffic_ops_golang binary
godir=src/github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang
( mkdir -p "$godir" && \
	cd "$godir" && \
	cp "$TC_DIR"/traffic_ops/traffic_ops_golang/traffic_ops_golang .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy TO DB admin
db_admin_dir=src/github.com/apache/trafficcontrol/traffic_ops/app/db
( mkdir -p "$db_admin_dir" && \
	cd "$db_admin_dir" && \
	cp "$TC_DIR"/traffic_ops/app/db/admin .
) || { echo "Could not copy go db admin at $(pwd): $!"; exit 1; };

# copy ToDnssecRefresh
to_dnssec_refresh_dir=src/github.com/apache/trafficcontrol/traffic_ops/app/bin/checks/DnssecRefresh
( mkdir -p "$to_dnssec_refresh_dir" && \
	cd "$to_dnssec_refresh_dir" && \
	cp "$TC_DIR"/traffic_ops/app/bin/checks/DnssecRefresh/ToDnssecRefresh .
) || { echo "Could not copy ToDnssecRefresh at $(pwd): $!"; exit 1; };

# copy TV DB reencrypt
reencrypt_dir=src/github.com/apache/trafficcontrol/traffic_ops/app/db/reencrypt
( mkdir -p "$reencrypt_dir" && \
	cd "$reencrypt_dir" && \
	cp "$TC_DIR"/traffic_ops/app/db/reencrypt/reencrypt .
) || { echo "Could not copy go db reencrypt at $(pwd): $!"; exit 1; };

# copy TV migrate
tvm_dir=src/github.com/apache/trafficcontrol/traffic_ops/app/db/traffic_vault_migrate
( mkdir -p "$tvm_dir" && \
	cd "$tvm_dir" && \
	cp "$TC_DIR"/traffic_ops/app/db/traffic_vault_migrate/traffic_vault_migrate .
) || { echo "Could not copy go db traffic_vault_migrate at $(pwd): $!"; exit 1; };

# copy TO profile converter
convert_dir=src/github.com/apache/trafficcontrol/traffic_ops/install/bin/convert_profile
( mkdir -p "$convert_dir" && \
	cd "$convert_dir" && \
	cp "$TC_DIR"/traffic_ops/install/bin/convert_profile/convert_profile .
) || { echo "Could not copy go profile converter at $(pwd): $!"; exit 1; };


%install

if [ -d $RPM_BUILD_ROOT ]; then
	%__rm -rf $RPM_BUILD_ROOT
fi

if [ ! -d $RPM_BUILD_ROOT/%{PACKAGEDIR} ]; then
	%__mkdir -p $RPM_BUILD_ROOT/%{PACKAGEDIR}
fi

%__cp -R $RPM_BUILD_DIR/traffic_ops-%{version}/* $RPM_BUILD_ROOT/%{PACKAGEDIR}
echo "go rming $RPM_BUILD_ROOT/%{PACKAGEDIR}/{pkg,src,bin}"
%__rm -rf $RPM_BUILD_ROOT/%{PACKAGEDIR}/{pkg,src,bin}

%__mkdir -p $RPM_BUILD_ROOT/var/www/files
%__cp install/data/json/osversions.json $RPM_BUILD_ROOT/var/www/files/.

# install traffic_ops_golang binary
if [ ! -d $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/bin ]; then
	%__mkdir -p $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/bin
fi

src=src/github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang
%__cp -p  "$src"/traffic_ops_golang        "${RPM_BUILD_ROOT}"/opt/traffic_ops/app/bin/traffic_ops_golang

db_admin_src=src/github.com/apache/trafficcontrol/traffic_ops/app/db
%__cp -p  "$db_admin_src"/admin           "${RPM_BUILD_ROOT}"/opt/traffic_ops/app/db/admin
%__rm $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/db/*.go
%__rm -r $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/db/trafficvault/test

to_dnssec_refresh_src=src/github.com/apache/trafficcontrol/traffic_ops/app/bin/checks/DnssecRefresh
%__cp -p  "$to_dnssec_refresh_src"/ToDnssecRefresh           "${RPM_BUILD_ROOT}"/opt/traffic_ops/app/bin/checks/DnssecRefresh/ToDnssecRefresh
%__rm $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/bin/checks/DnssecRefresh/*.go
%__rm -r $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/bin/checks/DnssecRefresh/config

reencrypt_src=src/github.com/apache/trafficcontrol/traffic_ops/app/db/reencrypt
%__cp -p  "$reencrypt_src"/reencrypt           "${RPM_BUILD_ROOT}"/opt/traffic_ops/app/db/reencrypt/reencrypt
%__rm $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/db/reencrypt/*.go

tv_migrate_src=src/github.com/apache/trafficcontrol/traffic_ops/app/db/traffic_vault_migrate
%__cp -p  "$tv_migrate_src"/traffic_vault_migrate           "${RPM_BUILD_ROOT}"/opt/traffic_ops/app/db/traffic_vault_migrate/traffic_vault_migrate
%__rm $RPM_BUILD_ROOT/%{PACKAGEDIR}/app/db/traffic_vault_migrate/*.go

convert_profile_src=src/github.com/apache/trafficcontrol/traffic_ops/install/bin/convert_profile
%__cp -p  "$convert_profile_src"/convert_profile           "${RPM_BUILD_ROOT}"/opt/traffic_ops/install/bin/convert_profile
%__rm $RPM_BUILD_ROOT/%{PACKAGEDIR}/install/bin/convert_profile/*.go

%pre
/usr/bin/getent group %{TRAFFIC_OPS_GROUP} || /usr/sbin/groupadd -r %{TRAFFIC_OPS_GROUP}
/usr/bin/getent passwd %{TRAFFIC_OPS_USER} || /usr/sbin/useradd -r -d %{PACKAGEDIR} -s /sbin/nologin %{TRAFFIC_OPS_USER} -g %{TRAFFIC_OPS_GROUP}
if [ -d %{PACKAGEDIR}/app/conf ]; then
echo -e "\nBacking up config files.\n"
if [ -f /var/tmp/traffic_ops-backup.tar ]; then
	%__rm /var/tmp/traffic_ops-backup.tar
fi
cd %{PACKAGEDIR} && tar cf /var/tmp/traffic_ops-backup.tar app/conf app/db/dbconf.yml app/local app/cpanfile.snapshot
fi

# upgrade
if [ "$1" == "2" ]; then
	systemctl stop traffic_ops
fi

%post
%__cp %{PACKAGEDIR}/etc/init.d/traffic_ops /etc/init.d/traffic_ops
%__mkdir -p /var/www/files
%__mkdir -p /etc/cron.d/
%__cp %{PACKAGEDIR}/etc/cron.d/trafops_dnssec_refresh /etc/cron.d/trafops_dnssec_refresh
%__cp %{PACKAGEDIR}/etc/cron.d/trafops_clean_isos /etc/cron.d/trafops_clean_isos
%__cp %{PACKAGEDIR}/etc/cron.d/autorenew_certs /etc/cron.d/autorenew_certs
%__cp %{PACKAGEDIR}/etc/logrotate.d/traffic_ops /etc/logrotate.d/traffic_ops
%__cp %{PACKAGEDIR}/etc/logrotate.d/traffic_ops_golang /etc/logrotate.d/traffic_ops_golang
%__cp %{PACKAGEDIR}/etc/logrotate.d/traffic_ops_access /etc/logrotate.d/traffic_ops_access
%__cp %{PACKAGEDIR}/etc/profile.d/traffic_ops.sh /etc/profile.d/traffic_ops.sh
%__chown root:root /etc/init.d/traffic_ops
%__chown root:root /etc/cron.d/trafops_dnssec_refresh
%__chown root:root /etc/cron.d/autorenew_certs
%__chown root:root /etc/cron.d/trafops_clean_isos
%__chown root:root /etc/logrotate.d/traffic_ops
%__chown root:root /etc/logrotate.d/traffic_ops_golang
%__chown root:root /etc/logrotate.d/traffic_ops_access
%__chmod +x /etc/init.d/traffic_ops
%__chmod +x %{PACKAGEDIR}/install/bin/*
/sbin/chkconfig --add traffic_ops

%__mkdir -p %{TRAFFIC_OPS_LOG_DIR} TRAFFIC_OPS_ROOT_CERTIFICATES_DIR

if [ -f /var/tmp/traffic_ops-backup.tar ]; then
	echo -e "\nRestoring config files.\n"
	cd %{PACKAGEDIR} && tar xf /var/tmp/traffic_ops-backup.tar
fi

# install
if [ "$1" = "1" ]; then
	# see postinstall, the .reconfigure file triggers init().
	echo -e "\nRun /opt/traffic_ops/install/bin/postinstall from the root home directory to complete the install.\n"
fi

# upgrade
if [ "$1" == "2" ]; then
	echo -e "\n\nTo complete the update, perform the following steps:\n"
	echo -e "1. If any *.rpmnew files are in /opt/traffic_ops/...,  reconcile with any local changes\n"
	echo -e "2. Run './db/admin --env production upgrade'\n"
	echo -e "   from the /opt/traffic_ops/app directory.\n"
	echo -e "To start Traffic Ops:  systemctl start traffic_ops\n";
	echo -e "To stop Traffic Ops:   systemctl stop traffic_ops\n\n";
fi
/bin/chown -R %{TRAFFIC_OPS_USER}:%{TRAFFIC_OPS_GROUP} %{PACKAGEDIR}
/bin/chown -R %{TRAFFIC_OPS_USER}:%{TRAFFIC_OPS_GROUP} %{TRAFFIC_OPS_LOG_DIR} TRAFFIC_OPS_ROOT_CERTIFICATES_DIR
setcap cap_net_bind_service=+ep %{PACKAGEDIR}/app/bin/traffic_ops_golang

%preun

if [ "$1" = "0" ]; then
	# stop service before starting the uninstall
	systemctl stop traffic_ops
fi

%postun

if [ "$1" = "0" ]; then
	# this is an uninstall
	%__rm -rf %{PACKAGEDIR}
	%__rm /etc/init.d/traffic_ops
	/usr/bin/getent passwd %{TRAFFIC_OPS_USER} || /usr/sbin/userdel %{TRAFFIC_OPS_USER}
	/usr/bin/getent group %{TRAFFIC_OPS_GROUP} || /usr/sbin/groupdel %{TRAFFIC_OPS_GROUP}
fi

%files
%license ../LICENSE
%defattr(644,root,root,755)
%attr(755,%{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/bin/traffic_ops_golang
%attr(755,%{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/script/detect10ginterfaces.pl
%attr(755,%{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/script/generate_raid0_files.pl
%config(noreplace) %attr(750,%{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) /opt/traffic_ops/app/conf
%config(noreplace) %attr(750,%{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) /opt/traffic_ops/app/db/dbconf.yml
%config(noreplace)/var/www/files/osversions.json
%attr(755, %{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/db/admin
%exclude %{PACKAGEDIR}/app/db/SQUASH.md
%exclude %{PACKAGEDIR}/app/db/squash_migrations.sh
%attr(755, %{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/install/bin/convert_profile/convert_profile
%attr(755, %{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/bin/checks/DnssecRefresh/ToDnssecRefresh
%attr(755, %{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/db/reencrypt/reencrypt
%attr(755, %{TRAFFIC_OPS_USER},%{TRAFFIC_OPS_GROUP}) %{PACKAGEDIR}/app/db/traffic_vault_migrate/traffic_vault_migrate
%{PACKAGEDIR}/etc
%{PACKAGEDIR}/app/bin/checks
%{PACKAGEDIR}/app/bin/tests
%{PACKAGEDIR}/app/db
%{PACKAGEDIR}/app/templates
%{PACKAGEDIR}/install
