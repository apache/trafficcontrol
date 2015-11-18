.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
.. 
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
.. 
..     http://www.apache.org/licenses/LICENSE-2.0
.. 
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
.. 

.. _to-api-v12-influxdb:

InfluxDB
==========

.. Note:: The documentation needs a thorough review!

**GET /api/1.2/traffic_monitor/stats.json**

Authentication Required: Yes

Role(s) Required: None

**Response Properties**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
| ``aaData``           | array  |                                                |
+----------------------+--------+------------------------------------------------+

**Response Example**
::

  {
   "aaData": [
      [
         "0",
         "ALL",
         "ALL",
         "ALL",
         "true",
         "ALL",
         "142035",
         "172365661.85"
      ],
      [
         1,
         "EDGE1_TOP_421_PSPP",
         "odol-atsec-atl-03",
         "us-ga-atlanta",
         "1",
         "REPORTED",
         "596",
         "923510.04",
         "69.241.82.126"
      ]
   ],
  }

|

**GET /api/1.2/redis/stats.json**

Authentication Required: Yes

Role(s) Required: None

**Response Properties**

+----------------------+--------+------------------------------------------------+
| Parameter            | Type   | Description                                    |
+======================+========+================================================+
|``number``            | array  |                                                |
+----------------------+--------+------------------------------------------------+
|``what``              | string |                                                |
+----------------------+--------+------------------------------------------------+
|``which``             | string |                                                |
+----------------------+--------+------------------------------------------------+
|``interval``          | string |                                                |
+----------------------+--------+------------------------------------------------+
|``elapsed``           | string |                                                |
+----------------------+--------+------------------------------------------------+
|``end``               | string |                                                |
+----------------------+--------+------------------------------------------------+
|``start``             | string |                                                |
+----------------------+--------+------------------------------------------------+

**Response Example**
::

  {
   "number": -1,
   "what": null,
   "which": null,
   "interval": " 10 seconds ",
   "elapsed": "0.11271 (0.112065) ",
   "end": "Thu Jan  1 00:00:00 1970",
   "start": "Thu Jan  1 00:00:00 1970"
  }


|

**GET /api/1.2/redis/info/:host_name.json**

Authentication Required: Yes

Role(s) Required: None

**Request Route Parameters**

+--------------------------+--------+--------------------------------------------+
| Parameter                | Type   | Description                                |
+==========================+========+============================================+
|``host_name``             | string |                                            |
+--------------------------+--------+--------------------------------------------+

**Response Properties**

+-------------------------------------+--------+-------------+
|              Parameter              |  Type  | Description |
+=====================================+========+=============+
| ``Server``                          | hash   |             |
+-------------------------------------+--------+-------------+
| ``>redis_build_id``                 | string |             |
+-------------------------------------+--------+-------------+
| ``>config_file``                    | string |             |
+-------------------------------------+--------+-------------+
| ``>uptime_in_seconds``              | string |             |
+-------------------------------------+--------+-------------+
| ``>hz``                             | string |             |
+-------------------------------------+--------+-------------+
| ``>os``                             | string |             |
+-------------------------------------+--------+-------------+
| ``>redis_git_sha1``                 | string |             |
+-------------------------------------+--------+-------------+
| ``>redis_version``                  | string |             |
+-------------------------------------+--------+-------------+
| ``>tcp_port``                       | string |             |
+-------------------------------------+--------+-------------+
| ``>redis_git_dirty``                | string |             |
+-------------------------------------+--------+-------------+
| ``>redis_mode``                     | string |             |
+-------------------------------------+--------+-------------+
| ``>run_id``                         | string |             |
+-------------------------------------+--------+-------------+
| ``>uptime_in_days``                 | string |             |
+-------------------------------------+--------+-------------+
| ``>gcc_version``                    | string |             |
+-------------------------------------+--------+-------------+
| ``>arch_bits``                      | string |             |
+-------------------------------------+--------+-------------+
| ``>lru_clock``                      | string |             |
+-------------------------------------+--------+-------------+
| ``>multiplexing_api``               | string |             |
+-------------------------------------+--------+-------------+
| ``Keyspace``                        | string |             |
+-------------------------------------+--------+-------------+
| ``>db0``                            | string |             |
+-------------------------------------+--------+-------------+
| ``slowlog``                         | array  |             |
+-------------------------------------+--------+-------------+
| ``Persistence``                     | hash   |             |
+-------------------------------------+--------+-------------+
| ``>rdb_bgsave_in_progress``         | string |             |
+-------------------------------------+--------+-------------+
| ``>loading``                        | string |             |
+-------------------------------------+--------+-------------+
| ``>rdb_current_bgsave_time_sec``    | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_enabled``                    | string |             |
+-------------------------------------+--------+-------------+
| ``>rdb_last_bgsave_time_sec``       | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_last_rewrite_time_sec``      | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_last_write_status``          | string |             |
+-------------------------------------+--------+-------------+
| ``>rdb_last_bgsave_status``         | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_last_bgrewrite_status``      | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_current_rewrite_time_sec``   | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_rewrite_scheduled``          | string |             |
+-------------------------------------+--------+-------------+
| ``>aof_rewrite_in_progress``        | string |             |
+-------------------------------------+--------+-------------+
| ``>rdb_last_save_time``             | string |             |
+-------------------------------------+--------+-------------+
| ``>rdb_changes_since_last_save``    | string |             |
+-------------------------------------+--------+-------------+
| ``slowlen``                         | int    |             |
+-------------------------------------+--------+-------------+
| ``CPU``                             | hash   |             |
+-------------------------------------+--------+-------------+
| ``>used_cpu_user``                  | string |             |
+-------------------------------------+--------+-------------+
| ``>used_cpu_sys``                   | string |             |
+-------------------------------------+--------+-------------+
| ``>used_cpu_user_children``         | string |             |
+-------------------------------------+--------+-------------+
| ``>used_cpu_sys_children``          | string |             |
+-------------------------------------+--------+-------------+
| ``Memory``                          | string |             |
+-------------------------------------+--------+-------------+
| ``>used_memory_lua``                | string |             |
+-------------------------------------+--------+-------------+
| ``>mem_allocator``                  | string |             |
+-------------------------------------+--------+-------------+
| ``>used_memory_human``              | string |             |
+-------------------------------------+--------+-------------+
| ``>used_memory_peak_human``         | string |             |
+-------------------------------------+--------+-------------+
| ``>used_memory_peak``               | string |             |
+-------------------------------------+--------+-------------+
| ``>used_memory_rss``                | string |             |
+-------------------------------------+--------+-------------+
| ``>mem_fragmentation_ratio``        | string |             |
+-------------------------------------+--------+-------------+
| ``>used_memory``                    | string |             |
+-------------------------------------+--------+-------------+
| ``Replication``                     | hash   |             |
+-------------------------------------+--------+-------------+
| ``>repl_backlog_first_byte_offset`` | string |             |
+-------------------------------------+--------+-------------+
| ``>repl_backlog_active``            | string |             |
+-------------------------------------+--------+-------------+
| ``>repl_backlog_histlen``           | string |             |
+-------------------------------------+--------+-------------+
| ``>repl_backlog_size``              | string |             |
+-------------------------------------+--------+-------------+
| ``>role``                           | string |             |
+-------------------------------------+--------+-------------+
| ``>master_repl_offset``             | string |             |
+-------------------------------------+--------+-------------+
| ``>connected_slaves``               | string |             |
+-------------------------------------+--------+-------------+
| ``Clients``                         | hash   |             |
+-------------------------------------+--------+-------------+
| ``>client_biggest_input_buf``       | string |             |
+-------------------------------------+--------+-------------+
| ``>client_longest_output_list``     | string |             |
+-------------------------------------+--------+-------------+
| ``>blocked_clients``                | string |             |
+-------------------------------------+--------+-------------+
| ``>connected_clients``              | string |             |
+-------------------------------------+--------+-------------+
| ``Stats``                           | hash   |             |
+-------------------------------------+--------+-------------+
| ``>latest_fork_usec``               | string |             |
+-------------------------------------+--------+-------------+
| ``>rejected_connections``           | string |             |
+-------------------------------------+--------+-------------+
| ``>sync_partial_ok``                | string |             |
+-------------------------------------+--------+-------------+
| ``>pubsub_channels``                | string |             |
+-------------------------------------+--------+-------------+
| ``>instantaneous_ops_per_sec``      | string |             |
+-------------------------------------+--------+-------------+
| ``>total_connections_received``     | string |             |
+-------------------------------------+--------+-------------+
| ``>pubsub_patterns``                | string |             |
+-------------------------------------+--------+-------------+
| ``>sync_full``                      | string |             |
+-------------------------------------+--------+-------------+
| ``>keyspace_hits``                  | string |             |
+-------------------------------------+--------+-------------+
| ``>keyspace_misses``                | string |             |
+-------------------------------------+--------+-------------+
| ``>total_commands_processed``       | string |             |
+-------------------------------------+--------+-------------+
| ``>expired_keys``                   | string |             |
+-------------------------------------+--------+-------------+
| ``>sync_partial_err``               | string |             |
+-------------------------------------+--------+-------------+

**Response Example**
::

  {
   "Server": {
      "redis_build_id": "606641459177bc09",
      "config_file": "\/etc\/redis\/redis.conf",
      "uptime_in_seconds": "1113787",
      "hz": "10",
      "os": "Linux 2.6.32-220.el6.x86_64 x86_64",
      "redis_git_sha1": "00000000",
      "redis_version": "2.8.15",
      "process_id": "14607",
      "tcp_port": "6379",
      "redis_git_dirty": "0",
      "redis_mode": "standalone",
      "run_id": "43c5d003453b96e38ad3eae54026d8e1b078a7fd",
      "uptime_in_days": "12",
      "gcc_version": "4.4.6",
      "arch_bits": "64",
      "lru_clock": "16050046",
      "multiplexing_api": "epoll"
   },
   "Keyspace": {
      "db0": "keys=26319,expires=0,avg_ttl=0"
   },
   "slowlog": [
      [
         "32656",
         "1425336191",
         "18539",
         [
            "keys",
            "*"
         ]
      ]
   ],
   "Persistence": {
      "rdb_bgsave_in_progress": "0",
      "loading": "0",
      "rdb_current_bgsave_time_sec": "-1",
      "aof_enabled": "0",
      "rdb_last_bgsave_time_sec": "-1",
      "aof_last_rewrite_time_sec": "-1",
      "aof_last_write_status": "ok",
      "rdb_last_bgsave_status": "ok",
      "aof_last_bgrewrite_status": "ok",
      "aof_current_rewrite_time_sec": "-1",
      "aof_rewrite_scheduled": "0",
      "aof_rewrite_in_progress": "0",
      "rdb_last_save_time": "1424222403",
      "rdb_changes_since_last_save": "2595831724"
   },
   "slowlen": 128,
   "CPU": {
      "used_cpu_user": "45252.98",
      "used_cpu_sys": "154718.84",
      "used_cpu_user_children": "0.00",
      "used_cpu_sys_children": "0.00"
   },
   "Memory": {
      "used_memory_lua": "33792",
      "mem_allocator": "jemalloc-3.6.0",
      "used_memory_human": "5.25G",
      "used_memory_peak_human": "8.08G",
      "used_memory_peak": "8675798632",
      "used_memory_rss": "8870088704",
      "mem_fragmentation_ratio": "1.57",
      "used_memory": "5633381640"
   },
   "Replication": {
      "repl_backlog_first_byte_offset": "0",
      "repl_backlog_active": "0",
      "repl_backlog_histlen": "0",
      "repl_backlog_size": "1048576",
      "role": "master",
      "master_repl_offset": "0",
      "connected_slaves": "0"
   },
   "Clients": {
      "client_biggest_input_buf": "0",
      "client_longest_output_list": "0",
      "blocked_clients": "0",
      "connected_clients": "16"
   },
   "Stats": {
      "latest_fork_usec": "0",
      "rejected_connections": "0",
      "sync_partial_ok": "0",
      "pubsub_channels": "0",
      "instantaneous_ops_per_sec": "2238",
      "total_connections_received": "2502657",
      "evicted_keys": "0",
      "pubsub_patterns": "0",
      "sync_full": "0",
      "keyspace_hits": "49388626",
      "keyspace_misses": "780",
      "total_commands_processed": "2645272238",
      "expired_keys": "0",
      "sync_partial_err": "0"
   }
  }


|

**GET /api/1.2/redis/match/#match/start_date/:start_date/end_date/:end_date/interval/:interval.json**

Authentication Required: Yes

Role(s) Required: None

**Request Route Parameters**

+--------------------------+--------+--------------------------------------------+
| Parameter                | Type   | Description                                |
+==========================+========+============================================+
|``start_date``            | string |                                            |
+--------------------------+--------+--------------------------------------------+
|``end_date``              | string |                                            |
+--------------------------+--------+--------------------------------------------+
|``interval``              | string |                                            |
+--------------------------+--------+--------------------------------------------+

**Response Properties**

+-------------+--------+-------------+
|  Parameter  |  Type  | Description |
+=============+========+=============+
| ``alerts``  | array  |             |
+-------------+--------+-------------+
| ``>level``  | string |             |
+-------------+--------+-------------+
| ``>text``   | string |             |
+-------------+--------+-------------+
| ``version`` | string |             |
+-------------+--------+-------------+

**Response Example**

TBD



