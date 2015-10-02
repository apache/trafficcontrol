%define debug_package %{nil}
Name:		traffic_stats
Version:	@VERSION@
Release:	@RELEASE@
Summary:	Tool to pull data from traffic monitor and store in Influxdb
Packager:	david_neuman2 at Cable dot Comcast dot com
Vendor:		Comcast Cable
Group:		Applications/Communications
License:	N/A
URL:		https://github.com/comcast/traffic_control/
Source:		~/rpmbuild/SOURCES/traffic_stats-@VERSION@.tar.gz

%description
Installs traffic_stats which is comprised of two seperate rpms:
	- traffic_stats: gets data from traffic monitor and stores in InfluxDB and also calculates daily summary of the data from InfluxDB and stores in traffic_ops

%prep

%setup

%build

%install
mkdir -p ${RPM_BUILD_ROOT}/opt/traffic_stats
mkdir -p ${RPM_BUILD_ROOT}/opt/traffic_stats/bin
mkdir -p ${RPM_BUILD_ROOT}/opt/traffic_stats/conf
mkdir -p ${RPM_BUILD_ROOT}/opt/traffic_stats/backup
mkdir -p ${RPM_BUILD_ROOT}/opt/traffic_stats/var/run
mkdir -p ${RPM_BUILD_ROOT}/opt/traffic_stats/var/log/traffic_stats
mkdir -p ${RPM_BUILD_ROOT}/etc/init.d
mkdir -p ${RPM_BUILD_ROOT}/etc/logrotate.d

cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats ${RPM_BUILD_ROOT}/opt/traffic_stats/bin/traffic_stats
cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats.cfg ${RPM_BUILD_ROOT}/opt/traffic_stats/conf/traffic_stats.cfg
cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats_seelog.xml ${RPM_BUILD_ROOT}/opt/traffic_stats/conf/traffic_stats_seelog.xml
cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats.init ${RPM_BUILD_ROOT}/etc/init.d/traffic_stats
cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats.logrotate ${RPM_BUILD_ROOT}/etc/logrotate.d/traffic_stats

%pre
/usr/bin/getent group traffic_stats >/dev/null

if [ $? -ne 0 ]; then

	/usr/sbin/groupadd -g 422 traffic_stats
fi

/usr/bin/getent passwd traffic_stats >/dev/null

if [ $? -ne 0 ]; then

	/usr/sbin/useradd -g traffic_stats -u 422 -d /opt/traffic_stats -M traffic_stats

fi

/usr/bin/passwd -l traffic_stats >/dev/null
/usr/bin/chage -E -1 -I -1 -m 0 -M 99999 -W 7 traffic_stats

if [ -e /etc/init.d/write_traffic_stats ]; then
	/sbin/service write_traffic_stats stop
fi

if [ -e /etc/init.d/ts_daily_summary ]; then
	/sbin/service ts_daily_summary stop
fi

if [ -e /etc/init.d/traffic_stats ]; then
	/sbin/service traffic_stats stop
fi

%post

/sbin/chkconfig --add traffic_stats
/sbin/chkconfig traffic_stats on

%files
%defattr(644, traffic_stats, traffic_stats, 755)

%config(noreplace) /opt/traffic_stats/conf/traffic_stats.cfg
%config(noreplace) /opt/traffic_stats/conf/traffic_stats_seelog.xml
%config(noreplace) /etc/logrotate.d/traffic_stats

%dir /opt/traffic_stats
%dir /opt/traffic_stats/bin
%dir /opt/traffic_stats/conf
%dir /opt/traffic_stats/backup
%dir /opt/traffic_stats/var
%dir /opt/traffic_stats/var/log
%dir /opt/traffic_stats/var/log/traffic_stats

%attr(600, traffic_stats, traffic_stats) /opt/traffic_stats/conf/*
%attr(755, traffic_stats, traffic_stats) /opt/traffic_stats/bin/*
%attr(755, traffic_stats, traffic_stats) /etc/init.d/traffic_stats

%preun
# args for hooks: http://www.ibm.com/developerworks/library/l-rpm2/
# if $1 = 0, this is an uninstallation, if $1 = 1, this is an upgrade (don't do anything)
if [ "$1" = "0" ]; then
	/sbin/chkconfig traffic_stats off
	/etc/init.d/traffic_stats stop
	/sbin/chkconfig --del traffic_stats
fi

if [ -e /etc/init.d/write_traffic_stats ]; then
	/sbin/chkconfig write_traffic_stats off
	/etc/init.d/write_traffic_stats stop
	/sbin/chkconfig --del write_traffic_stats
fi

if [ -e /etc/init.d/ts_daily_summary ]; then
	/sbin/chkconfig ts_daily_summary off
	/etc/init.d/ts_daily_summary stop
	/sbin/chkconfig --del ts_daily_summary
fi
