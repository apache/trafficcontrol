#!/bin/bash

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

VARNISHD_EXECUTABLE="/usr/sbin/varnishd"
HITCH_EXECUTABLE="/usr/sbin/hitch"

is_varnishd_running() {
  pgrep -x "$(basename "$VARNISHD_EXECUTABLE")" >/dev/null
}

start_varnishd() {
  if is_varnishd_running; then
    echo "varnishd is already running."
  else
    echo "Starting varnishd..."
    "$VARNISHD_EXECUTABLE" -f /opt/cache/etc/varnish/default.vcl -a :80 -a :6081,PROXY
    echo "varnishd is now running."
  fi
}

stop_varnishd() {
  if is_varnishd_running; then
    echo "Stopping varnishd..."

    # Send SIGTERM signal to varnishd to terminate gracefully
    pkill -x "$(basename "$VARNISHD_EXECUTABLE")"

    # Wait for varnishd to stop, giving it a timeout of 10 seconds
    timeout=10
    while is_varnishd_running; do
      if ((timeout-- == 0)); then
        echo "Timed out waiting for varnishd to stop. Sending SIGKILL..."
        pkill -9 -x "$(basename "$VARNISHD_EXECUTABLE")"
        break
      fi
      sleep 1
    done

    if is_varnishd_running; then
      echo "Failed to stop varnishd."
    else
      echo "varnishd is stopped."
    fi
  else
    echo "varnishd is not running."
  fi
}

restart_varnishd() {
  echo "Restarting varnishd..."
  stop_varnishd
  start_varnishd
}

is_hitch_running() {
  pgrep -x "$(basename "$HITCH_EXECUTABLE")" >/dev/null
}


start_hitch() {
  if is_hitch_running; then
    echo "hitch is already running."
  else
    echo "Starting hitch..."
    "$HITCH_EXECUTABLE" --config /opt/cache/etc/varnish/hitch.conf --daemon
    echo "hitch is now running."
  fi

}

reload_hitch() {
  if is_hitch_running; then
    pkill -HUP "$(basename "$HITCH_EXECUTABLE")"
  else
    echo "hitch is not running"
    exit 1
  fi
}

handle_varnish_service() {
  case "$1" in
    enable)
      exit 0
      ;;
    start)
      start_varnishd
      ;;
    stop)
      stop_varnishd
      ;;
    restart)
      restart_varnishd
      ;;
    status)
      if is_varnishd_running; then
        # t3c-apply looks for this specific string
        echo "Active: active"
        exit 0
      fi
      exit 3
      ;;
    *)
      echo "Usage: $0 {start|stop|restart|enable|status}"
      exit 1
  esac
}

handle_hitch_service() {
  case "$1" in
    enable)
      exit 0
      ;;
    start)
      start_hitch
      ;;
    reload)
      reload_hitch
      ;;
    status)
      if is_hitch_running; then
        # t3c-apply looks for this specific string
        echo "Active: active"
        exit 0
      fi
      exit 3
      ;;
    *)
      echo "Usage: $0 {start|stop|restart|reload|enable|status}"
      exit 1
  esac
}

if [[ $2 == "varnish.service" ]]
then
  handle_varnish_service $1
elif [[ $2 == "hitch.service" ]]
then
  handle_hitch_service $1
fi

exit 0
