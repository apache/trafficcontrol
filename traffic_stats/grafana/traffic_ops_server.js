/*

     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.
*/

/*
 * Scripted dashboard for traffic ops.
 *
 * Based on the grafana scripted.js script (which is ASF 2.0 licensed).
 */

'use strict';

// Setup some variables
var dashboard;

// All URL parameters are available via the ARGS object
var ARGS;

// Intialize a skeleton with nothing but a rows array and service object,
// and setting default time and refresh interval.
dashboard = {
  refresh: "30s",
  rows: [],
  // time can be overridden in the URL using from/to parameters, but this is
  // handled automatically in grafana core during dashboard initialization
  time: {
    from: "now-24h",
    to: "now"
  }
};

let which = 'argName';

if (ARGS.which !== undefined) {
  which = ARGS.which;
}

// Set a title
dashboard.title = which;

{
  dashboard.rows.push(
    {
      "height": "250px",
      "panels": [
        {
          "title": "bandwidth",
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 1,
          "datasource": "cache_stats",
          "renderer": "flot",
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "percentage": false,
          "legend": {
            "show": true,
            "values": true,
            "min": false,
            "max": true,
            "current": true,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "shared": true,
            "sort": 0,
            "msResolution": false
          },
          "timeFrom": null,
          "timeShift": null,
          "targets": [
            {
              "measurement": "bandwidth.1min",
              "tags": {},
              "query": `SELECT mean(value) FROM "monthly"."bandwidth.1min" WHERE hostname= '${which}' and $timeFilter GROUP BY time(60s)`,
              "rawQuery": true,
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "value"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  }
                ]
              ],
              "alias": "bandwidth"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "links": [],
          "yaxes": [
            {
              "show": true,
              "min": null,
              "max": null,
              "logBase": 1,
              "format": "Kbits"
            },
            {
              "show": true,
              "min": null,
              "max": null,
              "logBase": 1,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          }
        }
      ],
      "title": "Row",
      "collapse": false,
      "editable": true
    },
    {
      "height": "250px",
      "panels": [
        {
          "title": "conns",
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 2,
          "datasource": "cache_stats",
          "renderer": "flot",
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "percentage": false,
          "legend": {
            "show": true,
            "values": true,
            "min": false,
            "max": true,
            "current": true,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "shared": true,
            "sort": 0,
            "msResolution": false
          },
          "timeFrom": null,
          "timeShift": null,
          "targets": [
            {
              "measurement": "connections.1min",
              "tags": {},
              "query": `SELECT mean(value) FROM "monthly"."connections.1min" WHERE hostname= '${which}' and $timeFilter GROUP BY time(60s)`,
              "rawQuery": true,
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "value"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  }
                ]
              ],
              "alias": "connections"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "links": [],
          "yaxes": [
            {
              "show": true,
              "min": null,
              "max": null,
              "logBase": 1,
              "format": "short"
            },
            {
              "show": true,
              "min": null,
              "max": null,
              "logBase": 1
            }
          ],
          "xaxis": {
            "show": true
          }
        }
      ],
      "title": "Row",
      "collapse": false,
      "editable": true
    },
    {
      "title": "cpu and mem",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "title": "CPU Usage",
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "isNew": true,
          "id": 3,
          "targets": [
            {
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "host",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "usage_system"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "cpu_system"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "usage_iowait"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "cpu_iowait"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "usage_user"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "cpu_user"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "usage_guest"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "cpu_guest"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "usage_steal"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "cpu_steal"
                    ]
                  }
                ]
              ],
              "measurement": "cpu",
              "alias": "$col"
            }
          ],
          "datasource": "telegraf",
          "renderer": "flot",
          "yaxes": [
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "percent"
            },
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          },
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "percentage": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "individual",
            "shared": true,
            "msResolution": true,
            "sort": 2
          },
          "timeFrom": null,
          "timeShift": null,
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        },
        {
          "title": "Memory Usage",
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "isNew": true,
          "id": 4,
          "targets": [
            {
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "host",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "used_percent"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "mem_used"
                    ]
                  }
                ]
              ],
              "measurement": "mem",
              "alias": "$col"
            }
          ],
          "datasource": "telegraf",
          "renderer": "flot",
          "yaxes": [
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "percent"
            },
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          },
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "percentage": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "individual",
            "shared": true,
            "msResolution": true,
            "sort": 0
          },
          "timeFrom": null,
          "timeShift": null,
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        }
      ]
    },
    {
      "title": "load avg and diskio",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "title": "Load Average",
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "isNew": true,
          "id": 5,
          "targets": [
            {
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "host",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "load1"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "load1"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "load5"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "load5"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "load15"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "load15"
                    ]
                  }
                ]
              ],
              "measurement": "system",
              "alias": "$col"
            }
          ],
          "datasource": "telegraf",
          "renderer": "flot",
          "yaxes": [
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            },
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          },
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "percentage": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "shared": true,
            "msResolution": true,
            "sort": 0
          },
          "timeFrom": null,
          "timeShift": null,
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        },
        {
          "title": "Read/Write Time",
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "isNew": true,
          "id": 6,
          "targets": [
            {
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "host",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "read_time"
                    ]
                  },
                  {
                    "type": "sum",
                    "params": []
                  },
                  {
                    "type": "non_negative_derivative",
                    "params": [
                      "10s"
                    ]
                  },
                  {
                    "type": "alias",
                    "params": [
                      "read_time"
                    ]
                  }
                ]
              ],
              "measurement": "diskio",
              "alias": "$col"
            },
            {
              "refId": "B",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "host",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "write_time"
                    ]
                  },
                  {
                    "type": "sum",
                    "params": []
                  },
                  {
                    "type": "non_negative_derivative",
                    "params": [
                      "10s"
                    ]
                  },
                  {
                    "type": "alias",
                    "params": [
                      "write_time"
                    ]
                  }
                ]
              ],
              "measurement": "diskio",
              "alias": "$col"
            }
          ],
          "datasource": "telegraf",
          "renderer": "flot",
          "yaxes": [
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "ns"
            },
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          },
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "percentage": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "shared": true,
            "msResolution": true,
            "sort": 0
          },
          "timeFrom": null,
          "timeShift": null,
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        }
      ]
    },
    {
      "title": "Wrap Count and netstat",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "title": "wrap count",
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "isNew": true,
          "id": 7,
          "targets": [
            {
              "refId": "A",
              "policy": "monthly",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "hostname",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "vol1_wrap_count"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "vol1"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "vol2_wrap_count"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "vol2"
                    ]
                  }
                ]
              ],
              "measurement": "wrap_count.1min",
              "alias": "$col"
            }
          ],
          "datasource": "cache_stats",
          "renderer": "flot",
          "yaxes": [
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            },
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          },
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "percentage": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "shared": true,
            "sort": 0,
            "msResolution": true
          },
          "timeFrom": null,
          "timeShift": null,
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        },
        {
          "title": "netstat",
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "isNew": true,
          "id": 8,
          "targets": [
            {
              "refId": "A",
              "policy": "default",
              "dsType": "influxdb",
              "resultFormat": "time_series",
              "tags": [
                {
                  key: "host",
                  operator: "=",
                  value: which
                }
              ],
              "groupBy": [
                {
                  "type": "time",
                  "params": [
                    "$interval"
                  ]
                },
                {
                  "type": "fill",
                  "params": [
                    "null"
                  ]
                }
              ],
              "select": [
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_close"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_close"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_close_wait"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_close_wait"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_established"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_established"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_time_wait"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_time_wait"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_closing"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_closing"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_fin_wait1"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_fin_wait1"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_fin_wait2"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_fin_wait2"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_last_ack"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_last_ack"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_syn_recv"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_syn_recv"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "tcp_syn_sent"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "tcp_syn_sent"
                    ]
                  }
                ]
              ],
              "measurement": "netstat",
              "alias": "$col"
            }
          ],
          "datasource": "telegraf",
          "renderer": "flot",
          "yaxes": [
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            },
            {
              "label": null,
              "show": true,
              "logBase": 1,
              "min": null,
              "max": null,
              "format": "short"
            }
          ],
          "xaxis": {
            "show": true
          },
          "grid": {
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "percentage": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false,
            "hideEmpty": true,
            "hideZero": true
          },
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "shared": true,
            "sort": 2,
            "msResolution": true
          },
          "timeFrom": null,
          "timeShift": null,
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        }
      ]
    }
  );
}
return dashboard;
