# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
############################################################
# Dockerfile to build Traffic Server container images
#   as Edges for Traffic Control 1.4
# Based on CentOS 6.6
############################################################

# For cache, you may either use (RAM or disk) block devices or disk directories
# To use RAM block devices, pass them as /dev/ram0 and /dev/ram1 via `docker run --device`
# To use disk directories, simply don't pass devices, and the container will configure Traffic Server for directories

# Block devices may be created on the native machine with, for example, `modprobe brd`.
# The recommended minimum size for each block devices is 1G.
# For example, `sudo modprobe brd rd_size=1048576 rd_nr=2`

# Example Build and Run:
#
# docker build --rm --tag traffic_server_edge:1.4 Traffic_Server_Edge
#
# docker run --name my-edge-0 --hostname my-edge-0 --net cdnet --device /dev/ram0:/dev/ram0 --device /dev/ram1:/dev/ram1 --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --detach traffic_server_edge:1.4
#
# OR
#
# docker run --name my-edge-0 --hostname my-edge-0 --net cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --detach traffic_server_edge:1.4

FROM centos:6.6
MAINTAINER dev@trafficcontrol.apache.org

RUN yum install -y perl-JSON

RUN curl -O http://traffic-control-cdn.net/downloads/trafficserver-5.3.2-599.089d585.el6.x86_64.rpm
RUN yum install -y trafficserver-5.3.2-599.089d585.el6.x86_64.rpm

RUN mkdir /opt/ort
RUN cd /opt/ort && curl -LO https://github.com/apache/trafficcontrol/raw/RELEASE-1.4.0-RC0/traffic_ops/bin/traffic_ops_ort.pl
RUN chmod 777 /opt/ort/traffic_ops_ort.pl

RUN curl -O http://traffic-control-cdn.net/downloads/astats_over_http-1.2-8.el6.x86_64.rpm
RUN yum install -y astats_over_http-1.2-8.el6.x86_64.rpm

RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_cop
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_crashlog
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_ctl
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_layout
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_line
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_logcat
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_logstats
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_manager
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_sac
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/trafficserver
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_server
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_top
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_via
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/tspush
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/tsxs

# \todo move Heka to its own container, sharing the ATS log file via --volume
RUN curl -LO https://github.com/mozilla-services/heka/releases/download/v0.10.0/heka-0_10_0-linux-amd64.rpm
RUN yum install -y heka-0_10_0-linux-amd64.rpm
RUN mkdir etc/hekad
RUN printf '[ats_traffic_logs] \n\
type = "LogstreamerInput" \n\
splitter = "TokenSplitter" \n\
decoder = "ATS_transform_decoder" \n\
log_directory = "/opt/trafficserver/var/log/trafficserver" \n\
file_match = "custom_ats_2.log" \n\
[ATS_transform_decoder] \n\
type = "PayloadRegexDecoder" \n\
match_regex = '"'^(?P<UnixTimestamp>[\d]+\.[\d]+) chi=(?P<chi>\S+) phn=(?P<phn>\S+) shn=(?P<shn>\S+) url=(?P<url>\S+) cqhm=(?P<cqhm>\w+) cqhv=(?P<cqhv>\S+) pssc=(?P<pssc>\d+) ttms=(?P<ttms>\d+) b=(?P<b>\d+) sssc=(?P<sssc>\d+) sscl=(?P<sscl>\d+)  cfsc=(?P<cfsc>\S+) pfsc=(?P<pfsc>\S+) crc=(?P<crc>\S+) phr=(?P<phr>\S+) uas=(?P<uas>\S+) xmt=(?P<xmt>\S+)'"' \n\n\
[ATS_transform_decoder.message_fields] \n\
Type = "ats_traffic" \n\
timestamp = "%%UnixTimestamp%%" \n\
clientip = "%%chi%%" \n\
host = "%%phn%%" \n\
shn = "%%shn%%" \n\
url = "%%url%%" \n\
method = "%%cqhm%%" \n\
version = "%%cqhv%%" \n\
status = "%%pssc%%" \n\
request_duration = "%%ttms%%" \n\
bytes = "%%b%%" \n\
response_code = "%%sssc%%" \n\
response_length = "%%sscl%%" \n\
client_status = "%%cfsc%%" \n\
proxy_code = "%%pfsc%%" \n\
crc = "%%crc%%" \n\
phr = "%%phr%%" \n\
useragent = "%%uas%%" \n\
money_trace = "%%xmt%%" \n\
[PayloadEncoder] \n\
type = "PayloadEncoder" \n\
[FxaKafkaOutput] \n\
type = "KafkaOutput" \n\
topic = "ipcdn" \n\
message_matcher = "TRUE" \n\
encoder = "PayloadEncoder" \n\
addrs = ["{{.KafkaUri}}"] \n\
[Message_Counter] \n\
type = "CounterFilter" \n\
message_matcher = "Type != '"'heka.counter-output'"'" \n\
encoder = "CounterLogEncoder" \n\
[CounterLogEncoder] \n\
type="PayloadEncoder" \n\
append_newlines = true \n\
prefix_ts = true \n\
ts_format = "Mon Jan _2 15:04:05 MST 2006" \n\
[CounterLogOutput] \n\
type = "FileOutput" \n\
message_matcher = "Type == '"'heka.counter-output'"'" \n\
encoder = "CounterLogEncoder" \n\
path = "/tmp/hekad_counter.log"' > /etc/hekad/heka.toml

RUN printf '#!/bin/sh \n\n\
# \n\
# hekad <summary> \n\
# \n\
# chkconfig:   2345 80 20 \n\
# description: Starts and stops a single heka instance on this system \n\
# \n\
### BEGIN INIT INFO \n\
# Provides: Heka \n\
# Required-Start: $network $named \n\
# Required-Stop: $network $named \n\
# Default-Start: 2 3 4 5 \n\
# Default-Stop: 0 1 6 \n\
# Short-Description: This service manages the hekad daemon \n\
# Description: Heka is a high performance general purpose data acquisition and processing engine. \n\
### END INIT INFO \n\
# \n\
# init.d / servicectl compatibility (openSUSE) \n\
# \n\
if [ -f /etc/rc.status ]; then \n\
    . /etc/rc.status \n\
    rc_reset \n\
fi \n\
# \n\
# Source function library. \n\
# \n\
if [ -f /etc/rc.d/init.d/functions ]; then \n\
    . /etc/rc.d/init.d/functions \n\
fi \n\
. /etc/init.d/functions \n\
name="hekad" \n\
exec="/usr/bin/hekad" \n\
prog="hekad" \n\
user="root" \n\
group="root" \n\
pidfile=/var/run/${prog}.pid \n\
conf=/etc/hekad/heka.toml \n\
log=/var/log/heka.log \n\
DAEMON_ARGS=${DAEMON_ARGS---user root} \n\
nice=19 \n\
args=" --config $conf" \n\
[ -e /etc/sysconfig/$prog ] && . /etc/sysconfig/$prog \n\
lockfile=/var/lock/subsys/$prog \n\
HEKA_USER=root \n\
start() { \n\
    [ -x $exec ] || exit 5 \n\
    [ -f $CONF_FILE ] || exit 6 \n\
    # if not running, start it up here, usually something like "daemon $exec" \n\
    # Run the program! \n\
    nice -n ${nice} chroot --userspec $user:$group / sh -c " exec \"$prog\" $args " > ${log} 2>&1 & \n\
     # Generate the pidfile from here. If we instead made the forked process \n\
  # generate it there will be a race condition between the pidfile writing \n\
  # and a process possibly asking for status. \n\
  echo $! > $pidfile \n\
  echo "$name started." \n\
  return 0 \n\
} \n\
stop() { \n\
      # Try a few times to kill TERM the program \n\
  if status ; then \n\
    pid=`cat "$pidfile"` \n\
    echo "Killing $name (pid $pid) with SIGTERM" \n\
  kill -9 $pid \n\
    # Wait for it to exit. \n\
    for i in 1 2 3 4 5 ; do \n\
      echo "Waiting $name (pid $pid) to die..." \n\
      status || break \n\
      sleep 1 \n\
    done \n\
    if status ; then \n\
      echo "$name stop failed; still running." \n\
    else \n\
      echo "$name stopped." \n\
    fi \n\
  fi \n\
} \n\
restart() { \n\
    stop \n\
    start \n\
} \n\
reload() { \n\
    restart \n\
} \n\
force_reload() { \n\
    restart \n\
} \n\
status(){ \n\
  if [ -f "$pidfile" ] ; then \n\
    pid=`cat "$pidfile"` \n\
    if kill -0 $pid > /dev/null 2> /dev/null ; then \n\
      # process by this pid is running. \n\
      # It may not be our pid, but thats what you get with just pidfiles. \n\
      # TODO(sissel): Check if this process seems to be the same as the one we \n\
      # expect. Itd be nice to use flock here, but flock uses fork, not exec, \n\
      # so it makes it quite awkward to use in this case. \n\
      return 0 \n\
    else \n\
      return 2 # program is dead but pid file exists \n\
    fi \n\
  else \n\
    return 3 # program is not running \n\
  fi \n\
} \n\
rh_status() { \n\
    # run checks to determine if the service is running or use generic status \n\
    status -p $pidfile $prog \n\
} \n\
rh_status_q() { \n\
    rh_status >/dev/null 2>&1 \n\
} \n\
case "$1" in \n\
    start) \n\
        rh_status_q && exit 0 \n\
        $1 \n\
        ;; \n\
    stop) \n\
        rh_status_q || exit 0 \n\
        $1 \n\
        ;; \n\
    restart) \n\
        $1 \n\
        ;; \n\
    reload) \n\
        rh_status_q || exit 7 \n\
        $1 \n\
        ;; \n\
    force-reload) \n\
        force_reload \n\
        ;; \n\
    status) \n\
     status \n\
    code=$? \n\
    if [ $code -eq 0 ] ; then \n\
      echo "$prog is running" \n\
    else \n\
      echo "$prog is not running" \n\
    fi \n\
    exit $code \n\
        ;; \n\
    condrestart|try-restart) \n\
        rh_status_q || exit 0 \n\
        restart \n\
        ;; \n\
    *) \n\
        echo $"Usage: $0 {start|stop|status|restart|condrestart|try-restart|reload|force-reload}" \n\
        exit 2 \n\
esac' > /etc/init.d/hekad
RUN chmod +x /etc/init.d/hekad

EXPOSE 80 443
ADD run.sh /
ENTRYPOINT /run.sh
