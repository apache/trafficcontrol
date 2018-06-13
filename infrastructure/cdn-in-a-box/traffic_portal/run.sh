#!/usr/bin/env bash

NAME="Traffic Portal Application"
NODE_BIN_DIR="/usr/bin"
NODE_PATH="/opt/traffic_portal/node_modules"
FOREVER_BIN_DIR="/opt/traffic_portal/node_modules/forever/bin"
APPLICATION_PATH="/opt/traffic_portal/server.js"
PIDFILE="/var/run/traffic_portal.pid"
LOGFILE="/var/log/traffic_portal/traffic_portal.log"
MIN_UPTIME="5000"
SPIN_SLEEP_TIME="2000"

# Add node to the path for situations in which the environment is passed.
PATH=$FOREVER_BIN_DIR:$NODE_BIN_DIR:$PATH
forever \
    --pidFile $PIDFILE \
    -a \
    -l $LOGFILE \
    --minUptime $MIN_UPTIME \
    --spinSleepTime $SPIN_SLEEP_TIME \
    start $APPLICATION_PATH

tail -f /dev/null
