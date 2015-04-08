Requires: redis >= 2.8.15
Requires: spdb_gateway >= 1.18

%define debug_package %{nil}
Name:		tmredis
Version:	@VERSION@
Release:	@RELEASE@
Summary:	Tools to pull data from traffic monitor and store in Redis
Packager:	jeffrey_elsloo at Cable dot Comcast dot com
Vendor:		Comcast Cable NETO PAS VSS CDNENG
Group:		Applications/Communications
License:	N/A
URL:		https://gitlab.sys.comcast.net/cdneng/tm/
Source:		$RPM_SOURCE_DIR/tmredis-@VERSION@.tar.gz

%description
Installs tmredis tools.

%prep

%setup

%build

%install
mkdir -p ${RPM_BUILD_ROOT}/opt/tmredis
mkdir -p ${RPM_BUILD_ROOT}/opt/tmredis/bin
mkdir -p ${RPM_BUILD_ROOT}/opt/tmredis/conf
mkdir -p ${RPM_BUILD_ROOT}/opt/tmredis/backup
mkdir -p ${RPM_BUILD_ROOT}/opt/tmredis/var/log/tmredis
mkdir -p ${RPM_BUILD_ROOT}/etc/cron.d
mkdir -p ${RPM_BUILD_ROOT}/etc/init.d
mkdir -p ${RPM_BUILD_ROOT}/etc/logrotate.d
mkdir -p ${RPM_BUILD_ROOT}/etc/spdb_gateway.d

cp $GOPATH/src/redis_stats/rascal_2_redis ${RPM_BUILD_ROOT}/opt/tmredis/bin
cp $GOPATH/src/redis_stats/redis_2_spdb ${RPM_BUILD_ROOT}/opt/tmredis/bin
cp $GOPATH/src/redis_stats/backup_redis_daily ${RPM_BUILD_ROOT}/opt/tmredis/bin
cp $GOPATH/src/redis_stats/bck.sh ${RPM_BUILD_ROOT}/opt/tmredis/bin
cp $GOPATH/src/redis_stats/r2r_local.config ${RPM_BUILD_ROOT}/opt/tmredis/conf/r2r.cfg
cp $GOPATH/src/redis_stats/r2s_local.config ${RPM_BUILD_ROOT}/opt/tmredis/conf/r2s.cfg
cp $GOPATH/src/redis_stats/seelog.xml ${RPM_BUILD_ROOT}/opt/tmredis/conf
cp $GOPATH/src/redis_stats/tmredis.cron ${RPM_BUILD_ROOT}/etc/cron.d
cp $GOPATH/src/redis_stats/ipcdn_redis.cfg ${RPM_BUILD_ROOT}/etc/spdb_gateway.d
cp $GOPATH/src/redis_stats/tmredis.init ${RPM_BUILD_ROOT}/etc/init.d/tmredis
cp $GOPATH/src/redis_stats/tmredis.logrotate ${RPM_BUILD_ROOT}/etc/logrotate.d/tmredis

%pre
/usr/bin/getent group tmredis >/dev/null

if [ $? -ne 0 ]; then
	/usr/sbin/groupadd -g 422 tmredis
fi

/usr/bin/getent passwd tmredis >/dev/null

if [ $? -ne 0 ]; then
	/usr/sbin/useradd -g tmredis -u 422 -d /opt/tmredis -M tmredis
fi

/usr/bin/passwd -l tmredis >/dev/null
/usr/bin/chage -E -1 -I -1 -m 0 -M 99999 -W 7 tmredis

if [ -e /etc/init.d/tmredis ]; then
	/sbin/service tmredis stop
fi

%post
/sbin/chkconfig --add tmredis
/sbin/chkconfig tmredis on

%files
%defattr(644, tmredis, tmredis, 755)

%config(noreplace) /opt/tmredis/conf/r2r.cfg
%config(noreplace) /opt/tmredis/conf/r2s.cfg
%config(noreplace) /opt/tmredis/conf/seelog.xml
%config(noreplace) /etc/cron.d/tmredis.cron
%config(noreplace) /etc/logrotate.d/tmredis
%config(noreplace) /etc/spdb_gateway.d/ipcdn_redis.cfg

%dir /opt/tmredis
%dir /opt/tmredis/bin
%dir /opt/tmredis/conf
%dir /opt/tmredis/backup
%dir /opt/tmredis/var
%dir /opt/tmredis/var/log
%dir /opt/tmredis/var/log/tmredis

%attr(600, tmredis, tmredis) /opt/tmredis/conf/*
%attr(755, tmredis, tmredis) /opt/tmredis/bin/*
%attr(755, tmredis, tmredis) /etc/init.d/tmredis
%attr(644, root, root) /etc/cron.d/tmredis.cron

%preun
# args for hooks: http://www.ibm.com/developerworks/library/l-rpm2/
# if $1 = 0, this is an uninstallation, if $1 = 1, this is an upgrade (don't do anything)
if [ "$1" = "0" ]; then
	/sbin/chkconfig --del tmredis
	/etc/init.d/tmredis stop
fi
