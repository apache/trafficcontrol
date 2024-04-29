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
# RPM spec file for the Go version of Traffic Monitor (tm).
#
%define   debug_package %{nil}
Name:     traffic_monitor
Version:  %{traffic_control_version}
Release:  %{build_number}
Summary:  Monitor the caches
Vendor:   Apache Software Foundation
Group:    Applications/Communications
License:  Apache License, Version 2.0
URL:      https://github.com/apache/trafficcontrol
Source:   %{_sourcedir}/traffic_monitor-%{traffic_control_version}.tgz

%description
Installs traffic_monitor

%prep

%setup

%build
# copy traffic_monitor binary
godir=src/github.com/apache/trafficcontrol/traffic_monitor
( mkdir -p "$godir" && \
	cd "$godir" && \
	cp -r "$TC_DIR"/traffic_monitor/* .
) || { echo "Could not copy go program at $(pwd): $!"; exit 1; }

%install
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_monitor
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_monitor/bin
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_monitor/conf
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_monitor/backup
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_monitor/static
mkdir -p "${RPM_BUILD_ROOT}"/opt/traffic_monitor/var/run
mkdir -p "${RPM_BUILD_ROOT}"/var/log/traffic_monitor

# TODO: The /opt/traffic_monitor/var/log symlink is deprecated and should be removed for ATC 9.0.0.
ln -sfT /var/log/traffic_monitor "${RPM_BUILD_ROOT}"/opt/traffic_monitor/var/log
mkdir -p "${RPM_BUILD_ROOT}"/etc/init.d
mkdir -p "${RPM_BUILD_ROOT}"/etc/logrotate.d

src=src/github.com/apache/trafficcontrol/traffic_monitor
cp -p "$src"/traffic_monitor               "${RPM_BUILD_ROOT}"/opt/traffic_monitor/bin/traffic_monitor
cp "$src"/static/index.html                "${RPM_BUILD_ROOT}"/opt/traffic_monitor/static/index.html
cp "$src"/static/script.js                 "${RPM_BUILD_ROOT}"/opt/traffic_monitor/static/script.js
cp "$src"/static/style.css                 "${RPM_BUILD_ROOT}"/opt/traffic_monitor/static/style.css
cp "$src"/conf/traffic_ops.cfg             "${RPM_BUILD_ROOT}"/opt/traffic_monitor/conf/traffic_ops.cfg
cp "$src"/conf/traffic_monitor.cfg         "${RPM_BUILD_ROOT}"/opt/traffic_monitor/conf/traffic_monitor.cfg
cp "$src"/build/traffic_monitor.init       "${RPM_BUILD_ROOT}"/etc/init.d/traffic_monitor
cp "$src"/build/traffic_monitor.logrotate  "${RPM_BUILD_ROOT}"/etc/logrotate.d/traffic_monitor

%pre
old_log_dir=/opt/traffic_monitor/var/log
new_log_dir=/var/log/traffic_monitor
if [[ -d "$old_log_dir" ]]; then
	if [[ -d "$new_log_dir" ]]; then
		(
		# Include files starting with . in the * glob
		shopt -s dotglob
		mv "$old_log_dir"/* "$new_log_dir" || true
		)
		rmdir "$old_log_dir"
	else
		mv "$old_log_dir" "$new_log_dir"
	fi
	sync
fi

/usr/bin/getent group traffic_monitor >/dev/null
if [ $? -ne 0 ]; then
	/usr/sbin/groupadd -g 423 traffic_monitor
fi

/usr/bin/getent passwd traffic_monitor >/dev/null

if [ $? -ne 0 ]; then
	/usr/sbin/useradd -g traffic_monitor -u 423 -d /opt/traffic_monitor -M traffic_monitor
fi

/usr/bin/passwd -l traffic_monitor >/dev/null
/usr/bin/chage -E -1 -I -1 -m 0 -M 99999 -W 7 traffic_monitor

if [ -e /etc/init.d/traffic_monitor ]; then
	/sbin/service traffic_monitor stop
fi

#don't install over the top of java TM.  This is a workaround since yum doesn't respect the Conflicts tag.
if [[ $(rpm -q traffic_monitor --qf "%{VERSION}-%{RELEASE}") < 1.9.0 ]]
then
		echo -e "\n****************\n"
		echo "A java version of traffic_monitor is installed.  Please backup/remove that version before installing the golang version of traffic_monitor."
		echo -e "\n****************\n"
		exit 1
fi

%post

/sbin/chkconfig --add traffic_monitor
/sbin/chkconfig traffic_monitor on

%files
%license LICENSE
%defattr(644, traffic_monitor, traffic_monitor, 755)
%config(noreplace) %attr(600, traffic_monitor, traffic_monitor) /opt/traffic_monitor/conf/traffic_monitor.cfg
%config(noreplace) %attr(600, traffic_monitor, traffic_monitor) /opt/traffic_monitor/conf/traffic_ops.cfg
%config(noreplace) %attr(644, root, root) /etc/logrotate.d/traffic_monitor

%dir /opt/traffic_monitor
%dir /opt/traffic_monitor/bin
%dir /opt/traffic_monitor/conf
%dir /opt/traffic_monitor/backup
%dir /opt/traffic_monitor/static
%dir /opt/traffic_monitor/var
# TODO: The /opt/traffic_monitor/var/log symlink is deprecated and should be removed for ATC 9.0.0.
/opt/traffic_monitor/var/log
%dir /var/log/traffic_monitor
%dir /opt/traffic_monitor/var/run

%attr(600, traffic_monitor, traffic_monitor) /opt/traffic_monitor/static/index.html
%attr(600, traffic_monitor, traffic_monitor) /opt/traffic_monitor/static/script.js
%attr(600, traffic_monitor, traffic_monitor) /opt/traffic_monitor/static/style.css
%attr(755, traffic_monitor, traffic_monitor) /opt/traffic_monitor/bin/traffic_monitor
%attr(755, traffic_monitor, traffic_monitor) /etc/init.d/traffic_monitor

%preun
# args for hooks: https://www.ibm.com/developerworks/library/l-rpm2/
# if $1 = 0, this is an uninstallation, if $1 = 1, this is an upgrade (don't do anything)
if [ "$1" = "0" ]; then
	/sbin/chkconfig traffic_monitor off
	/etc/init.d/traffic_monitor stop
	/sbin/chkconfig --del traffic_monitor
fi
