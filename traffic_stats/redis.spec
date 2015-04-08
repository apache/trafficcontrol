%define pid_dir %{_localstatedir}/run/redis
%define pid_file %{pid_dir}/redis.pid

Summary: redis
Name: redis
Version: 2.8.15
Release: 1
License: BSD
Group: Applications/Multimedia
URL: http://code.google.com/p/redis/

Source0: redis-%{version}.tar.gz
Source1: redis.conf

BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
BuildRequires: gcc, make
Requires(post): /sbin/chkconfig /usr/sbin/useradd
Requires(preun): /sbin/chkconfig, /sbin/service
Requires(postun): /sbin/service
Provides: redis

#Packager: Jason Priebe <jpriebe@cbcnewmedia.com>
# Snagged from: https://github.com/causes/redis-centos/blob/master/spec/redis.spec
Packager: Jeff Elsloo <jeffrey_elsloo@cable.comcast.com>

%description
Redis is a key-value database. It is similar to memcached but the dataset is
not volatile, and values can be strings, exactly like in memcached, but also
lists and sets with atomic operations to push/pop elements.

In order to be very fast but at the same time persistent the whole dataset is
taken in memory and from time to time and/or when a number of changes to the
dataset are performed it is written asynchronously on disk. You may lose the
last few queries that is acceptable in many applications but it is as fast
as an in memory DB (beta 6 of Redis includes initial support for master-slave
replication in order to solve this problem by redundancy).

Compression and other interesting features are a work in progress. Redis is
written in ANSI C and works in most POSIX systems like Linux, *BSD, Mac OS X,
and so on. Redis is free software released under the very liberal BSD license.


%prep
%setup

%{__cat} <<EOF >redis.logrotate
%{_localstatedir}/log/redis/*log {
    missingok
}
EOF

%{__cat} <<'EOF' >redis.sysv
#!/bin/bash
#
# Init file for redis
#
# Written by Jason Priebe <jpriebe@cbcnewmedia.com>
#
# chkconfig: - 80 12
# description: A persistent key-value database with network interface
# processname: redis-server
# config: /etc/redis/redis.conf
# pidfile: %{pidfile}

source %{_sysconfdir}/init.d/functions

RETVAL=0
prog="redis-server"

start() {
  echo -n $"Starting $prog: "
  daemon --user redis --pidfile %{pid_file} %{_sbindir}/$prog /etc/redis/redis.conf
  RETVAL=$?
  echo
  [ $RETVAL -eq 0 ] && touch %{_localstatedir}/lock/subsys/$prog
  return $RETVAL
}

stop() {
    PID=`cat %{pid_file} 2>/dev/null`
    if [ -n "$PID" ]; then
        echo "Shutdown may take a while; redis needs to save the entire database";
        echo -n $"Shutting down $prog: "
        /usr/bin/redis-cli shutdown
        if checkpid $PID 2>&1; then
            echo_failure
            RETVAL=1
        else
            rm -f /var/lib/redis/temp*rdb
            rm -f /var/lock/subsys/$prog
            echo_success
            RETVAL=0
        fi
    else
        echo -n $"$prog is not running"
        echo_failure
        RETVAL=1
    fi

    echo
    return $RETVAL
}

restart() {
  stop
  start
}

condrestart() {
    [-e /var/lock/subsys/$prog] && restart || :
}

case "$1" in
  start)
  start
  ;;
  stop)
  stop
  ;;
  status)
  status -p %{pid_file} $prog
  RETVAL=$?
  ;;
  restart)
  restart
  ;;
  condrestart|try-restart)
  condrestart
  ;;
   *)
  echo $"Usage: $0 {start|stop|status|restart|condrestart}"
  RETVAL=1
esac

exit $RETVAL
EOF


%build
%{__make}

%install
%{__rm} -rf %{buildroot}
mkdir -p %{buildroot}%{_bindir}
%{__install} -Dp -m 0755 src/redis-server %{buildroot}%{_sbindir}/redis-server
%{__install} -Dp -m 0755 src/redis-benchmark %{buildroot}%{_bindir}/redis-benchmark
%{__install} -Dp -m 0755 src/redis-cli %{buildroot}%{_bindir}/redis-cli

%{__install} -Dp -m 0755 redis.logrotate %{buildroot}%{_sysconfdir}/logrotate.d/redis
%{__install} -Dp -m 0755 redis.sysv %{buildroot}%{_sysconfdir}/init.d/redis
%{__install} -Dp -m 0644 %{SOURCE1} %{buildroot}%{_sysconfdir}/redis/redis.conf
%{__install} -p -d -m 0755 %{buildroot}%{_localstatedir}/lib/redis
%{__install} -p -d -m 0755 %{buildroot}%{_localstatedir}/log/redis
%{__install} -p -d -m 0755 %{buildroot}%{pid_dir}

%pre
/usr/sbin/useradd -c 'Redis' -u 421 -s /bin/false -r -d %{_localstatedir}/lib/redis redis 2> /dev/null || :

%preun
if [ $1 = 0 ]; then
    # make sure redis service is not running before uninstalling

    # when the preun section is run, we've got stdin attached.  If we
    # call stop() in the redis init script, it will pass stdin along to
    # the redis-cli script; this will cause redis-cli to read an extraneous
    # argument, and the redis-cli shutdown will fail due to the wrong number
    # of arguments.  So we do this little bit of magic to reconnect stdin
    # to the terminal
    term="/dev/$(ps -p$$ --no-heading | awk '{print $2}')"
    exec < $term

    /sbin/service redis stop > /dev/null 2>&1 || :
    /sbin/chkconfig --del redis
fi

%post
/sbin/chkconfig --add redis

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{_sbindir}/redis-server
%{_bindir}/redis-benchmark
%{_bindir}/redis-cli
%{_sysconfdir}/init.d/redis
%config(noreplace) %{_sysconfdir}/redis/redis.conf
%{_sysconfdir}/logrotate.d/redis
%dir %attr(0770,redis,redis) %{_localstatedir}/lib/redis
%dir %attr(0755,redis,redis) %{_localstatedir}/log/redis
%dir %attr(0755,redis,redis) %{_localstatedir}/run/redis

%changelog
* Mon Sep 15 2014 - jse 2.8.15
- Upped version to 2.8.15
- Made changes related to our environment

* Tue Jul 13 2010 - jay at causes dot com 2.0.0-rc2
- upped to 2.0.0-rc2

* Mon May 24 2010 - jay at causes dot com 1.3.9-2
- moved pidfile back to /var/run/redis/redis.pid, so the redis
  user can write to the pidfile.
- Factored it out into %{pid_dir} (/var/run/redis), and
  %{pid_file} (%{pid_dir}/redis.pid)


* Wed May 05 2010 - brad at causes dot com 1.3.9-1
- redis updated to version 1.3.9 (development release from GitHub)
- extract config file from spec file
- move pid file from /var/run/redis/redis.pid to just /var/run/redis.pid
- move init file to /etc/init.d/ instead of /etc/rc.d/init.d/

* Fri Sep 11 2009 - jpriebe at cbcnewmedia dot com 1.0-1
- redis updated to version 1.0 stable

* Mon Jun 01 2009 - jpriebe at cbcnewmedia dot com 0.100-1
- Massive redis changes in moving from 0.09x to 0.100
- removed log timestamp patch; this feature is now part of standard release

* Tue May 12 2009 - jpriebe at cbcnewmedia dot com 0.096-1
- A memory leak when passing more than 16 arguments to a command (including
  itself).
- A memory leak when loading compressed objects from disk is now fixed.

* Mon May 04 2009 - jpriebe at cbcnewmedia dot com 0.094-2
- Patch: applied patch to add timestamp to the log messages
- moved redis-server to /usr/sbin
- set %config(noreplace) on redis.conf to prevent config file overwrites
  on upgrade

* Fri May 01 2009 - jpriebe at cbcnewmedia dot com 0.094-1
- Bugfix: 32bit integer overflow bug; there was a problem with datasets
  consisting of more than 20,000,000 keys resulting in a lot of CPU usage
  for iterated hash table resizing.

* Wed Apr 29 2009 - jpriebe at cbcnewmedia dot com 0.093-2
- added message to init.d script to warn user that shutdown may take a while

* Wed Apr 29 2009 - jpriebe at cbcnewmedia dot com 0.093-1
- version 0.093: fixed bug in save that would cause a crash
- version 0.092: fix for bug in RANDOMKEY command

* Fri Apr 24 2009 - jpriebe at cbcnewmedia dot com 0.091-3
- change permissions on /var/log/redis and /var/run/redis to 755; this allows
  non-root users to check the service status and to read the logs

* Wed Apr 22 2009 - jpriebe at cbcnewmedia dot com 0.091-2
- cleanup of temp*rdb files in /var/lib/redis after shutdown
- better handling of pid file, especially with status

* Tue Apr 14 2009 - jpriebe at cbcnewmedia dot com 0.091-1
- Initial release.
