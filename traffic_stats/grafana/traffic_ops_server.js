/*
     Copyright 2015 Comcast Cable Communications Management, LLC

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

/* global _ */

/*
 * Scripted dashboard for traffic ops.
 *
 * Based on the grafana scripted.js script (which is ASF 2.0 licensed).
 */



// accessible variables in this scope
var window, document, ARGS, $, jQuery, moment, kbn;

// Setup some variables
var dashboard;

// All url parameters are available via the ARGS object
var ARGS;

// Intialize a skeleton with nothing but a rows array and service object
dashboard = {
  rows : [],
};


// Set default time
// time can be overriden in the url using from/to parameters, but this is
// handled automatically in grafana core during dashboard initialization
dashboard.time = {
  from: "now-24h",
  to: "now"
};

var which = 'argName';

if(!_.isUndefined(ARGS.which)) {
  which = ARGS.which;
}

// Set a title
dashboard.title = which;
//set refresh interval
dashboard.refresh = "30s";


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
          "x-axis": true,
          "y-axis": true,
          "y_formats": [
            "bps",
            "short"
          ],
          "grid": {
            "leftLogBase": 1,
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "rightLogBase": 1,
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
            "shared": true
          },
          "timeFrom": null,
          "timeShift": null,
          "targets": [
            {
              "measurement": "bandwidth.1min",
              "tags": {},
              "query": "SELECT mean(value)*1000 FROM \"monthly\".\"bandwidth.1min\" WHERE hostname='" + which + "'  and $timeFilter GROUP BY time(60s)",
              "rawQuery": true
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
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
          "x-axis": true,
          "y-axis": true,
          "y_formats": [
            "short"
          ],
          "grid": {
            "leftLogBase": 1,
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "rightLogBase": 1,
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
            "shared": true
          },
          "timeFrom": null,
          "timeShift": null,
          "targets": [
            {
              "measurement": "connections.1min",
              "tags": {},
              "query": "SELECT mean(value) FROM \"monthly\".\"connections.1min\" WHERE hostname='" + which + "'  and $timeFilter GROUP BY time(60s)",
              "rawQuery": true
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
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
                  "key": "host",
                  "operator": "=~",
                  "value": "/" + which + "/"
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
                ]
              ],
              "measurement": "cpu"
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
            "msResolution": true
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
                  "key": "host",
                  "operator": "=~",
                  "value": "/" + which + "/"
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
              "measurement": "mem"
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
            "msResolution": true
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
                  "key": "host",
                  "operator": "=~",
                  "value": "/" + which + "/"
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
              "measurement": "system"
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
            "msResolution": true
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
                  "key": "host",
                  "operator": "=~",
                  "value": "/" + which + "/"
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
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "read_time"
                    ]
                  }
                ],
                [
                  {
                    "type": "field",
                    "params": [
                      "write_time"
                    ]
                  },
                  {
                    "type": "mean",
                    "params": []
                  },
                  {
                    "type": "alias",
                    "params": [
                      "write_time"
                    ]
                  }
                ]
              ],
              "measurement": "diskio"
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
