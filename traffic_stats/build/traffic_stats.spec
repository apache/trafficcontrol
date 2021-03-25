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
# RPM spec file for Traffic Stats (tm).
#
%define   debug_package %{nil}
Name:     traffic_stats
Version:  %{traffic_control_version}
Release:  %{build_number}
Summary:  Tool to pull data from traffic monitor and store in Influxdb
Packager: david_neuman2 at Cable dot Comcast dot com
Vendor:   Apache Software Foundation
Group:    Applications/Communications
License:  Apache License, Version 2.0
URL:      https://github.com/apache/trafficcontrol
Source:   %{_sourcedir}/traffic_stats-%{traffic_control_version}.tgz

%description
Installs traffic_stats which performs the follwing functions:
	1. Gets data from Traffic Monitor via a RESTful API and stores the data in InfluxDb
	2. Calculates Daily Summary stats from the raw data and stores it in Traffic Ops as well as InfluxDb

%prep

%setup

%build
# copy traffic_stats client
godir=src/github.com/apache/trafficcontrol/traffic_stats
( mkdir -p "$godir" && \
	cd "$godir" && \
	cp -LR "$TC_DIR"/traffic_stats/* .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

# copy influxdb_tools
godir=src/github.com/apache/trafficcontrol/traffic_stats/influxdb_tools
( mkdir -p "$godir" && \
	cd "$godir" && \
	cp -R "$TC_DIR"/traffic_stats/influxdb_tools/* .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

%install
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats/bin
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats/conf
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats/backup
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats/influxdb_tools
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats/var/run
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_stats/var/log/traffic_stats
mkdir -p "${RPM_BUILD_ROOT}"/etc/init.d
mkdir -p "${RPM_BUILD_ROOT}"/etc/logrotate.d
mkdir -p "${RPM_BUILD_ROOT}"/usr/share/grafana/public/dashboards/

src=src/github.com/apache/trafficcontrol/traffic_stats
cp -p "$src"/traffic_stats         "${RPM_BUILD_ROOT}"/opt/traffic_stats/bin/traffic_stats
cp "$src"/traffic_stats.cfg        "${RPM_BUILD_ROOT}"/opt/traffic_stats/conf/traffic_stats.cfg
cp "$src"/traffic_stats_seelog.xml "${RPM_BUILD_ROOT}"/opt/traffic_stats/conf/traffic_stats_seelog.xml
cp "$src"/traffic_stats.init       "${RPM_BUILD_ROOT}"/etc/init.d/traffic_stats
cp "$src"/traffic_stats.logrotate  "${RPM_BUILD_ROOT}"/etc/logrotate.d/traffic_stats
cp "$src"/grafana/*.js             "${RPM_BUILD_ROOT}"/usr/share/grafana/public/dashboards/
cp "$src"/influxdb_tools/sync_ts_databases  "${RPM_BUILD_ROOT}"/opt/traffic_stats/influxdb_tools/
cp "$src"/influxdb_tools/create_ts_databases  "${RPM_BUILD_ROOT}"/opt/traffic_stats/influxdb_tools/


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
%license LICENSE
%defattr(644, traffic_stats, traffic_stats, 755)

%config(noreplace) %attr(600, traffic_stats, traffic_stats) /opt/traffic_stats/conf/traffic_stats.cfg
%config(noreplace) %attr(600, traffic_stats, traffic_stats) /opt/traffic_stats/conf/traffic_stats_seelog.xml
%config(noreplace) %attr(644, root, root) /etc/logrotate.d/traffic_stats

%dir /opt/traffic_stats
%dir /opt/traffic_stats/bin
%dir /opt/traffic_stats/conf
%dir /opt/traffic_stats/backup
%dir /opt/traffic_stats/var
%dir /opt/traffic_stats/var/log
%dir /opt/traffic_stats/var/run
%dir /opt/traffic_stats/var/log/traffic_stats
%dir /usr/share/grafana/public/dashboards
%dir /opt/traffic_stats/influxdb_tools

%attr(755, traffic_stats, traffic_stats) /opt/traffic_stats/bin/traffic_stats
%attr(755, traffic_stats, traffic_stats) /etc/init.d/traffic_stats
%attr(644, traffic_stats, traffic_stats) /usr/share/grafana/public/dashboards/traffic_ops_cachegroup.js
%attr(644, traffic_stats, traffic_stats) /usr/share/grafana/public/dashboards/traffic_ops_deliveryservice.js
%attr(644, traffic_stats, traffic_stats) /usr/share/grafana/public/dashboards/traffic_ops_scripted.js
%attr(644, traffic_stats, traffic_stats) /usr/share/grafana/public/dashboards/traffic_ops_server.js
%attr(755, traffic_stats, traffic_stats) /opt/traffic_stats/influxdb_tools/create_ts_databases
%attr(755, traffic_stats, traffic_stats) /opt/traffic_stats/influxdb_tools/sync_ts_databases

%preun
# args for hooks: https://www.ibm.com/developerworks/library/l-rpm2/
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
