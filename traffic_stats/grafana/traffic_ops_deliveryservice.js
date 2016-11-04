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
  from: "now-6h",
  to: "now-60s"
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
    dashboard.rows.push( {
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "aliasColors": {
            "bw {deliveryservice: xb-dlassets}": "#7EB26D"
          },
          "bars": false,
          "datasource": "deliveryservice_stats",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "leftLogBase": 1,
            "leftMax": null,
            "leftMin": null,
            "rightLogBase": 1,
            "rightMax": null,
            "rightMin": null,
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "id": 1,
          "legend": {
            "avg": false,
            "current": true,
            "max": true,
            "min": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "connected",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 12,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "measurement": "bw",
              "query": "SELECT mean(value)*1000 FROM \"monthly\".\"kbps.ds.1min\" WHERE deliveryservice='" + which + "' and cachegroup = 'total'  and $timeFilter GROUP BY time(60s), deliveryservice ORDER BY asc",
              "rawQuery": true,
              "tags": {
                "deliveryservice": which
              },
              "alias": "$tag_deliveryservice"
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "bandwidth",
          "tooltip": {
            "shared": true,
            "value_type": "cumulative"
          },
          "type": "graph",
          "x-axis": true,
          "y-axis": true,
          "y_formats": [
            "bps",
            "short"
          ]
        }
      ],
      "title": "DsGbps"
    },
  {
      "title": "TPS!",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "title": "tps",
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 3,
          "datasource": "deliveryservice_stats",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "y_formats": [
            "short",
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
            "shared": true
          },
          "timeFrom": null,
          "timeShift": null,
          "targets": [
            {
              "measurement": "tps_2xx",
              "tags": {
                "deliveryservice": which
              },
              "query": "SELECT mean(value) FROM \"monthly\".\"tps_2xx.ds.1min\" WHERE $timeFilter AND deliveryservice='" + which + "' GROUP BY time(60s) ORDER BY asc",
              "hide": false,
              "rawQuery": true
            },
            {
              "target": "",
              "rawQuery": true,
              "query": "SELECT mean(value) FROM \"monthly\".\"tps_3xx.ds.1min\" WHERE $timeFilter AND deliveryservice='" + which + "' GROUP BY time(60s) ORDER BY asc"
            },
            {
              "target": "",
              "rawQuery": true,
              "query": "SELECT mean(value) FROM \"monthly\".\"tps_4xx.ds.1min\" WHERE $timeFilter AND deliveryservice='" + which + "' GROUP BY time(60s) ORDER BY asc"
            },
            {
              "target": "",
              "rawQuery": true,
              "query": "SELECT mean(value) FROM \"monthly\".\"tps_5xx.ds.1min\" WHERE $timeFilter AND deliveryservice='" + which + "' GROUP BY time(60s) ORDER BY asc"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "links": []
        }
      ]
    },
    {
      "collapse": false,
      "editable": true,
      "height": "250px",
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "datasource": "deliveryservice_stats",
          "editable": true,
          "error": false,
          "fill": 1,
          "grid": {
            "leftLogBase": 1,
            "leftMax": null,
            "leftMin": null,
            "rightLogBase": 1,
            "rightMax": null,
            "rightMin": null,
            "threshold1": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2": null,
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "id": 2,
          "legend": {
            "avg": false,
            "current": true,
            "hideEmpty": false,
            "max": true,
            "min": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "nullPointMode": "connected",
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "span": 12,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "query": "SELECT mean(value)*1000 FROM \"monthly\".\"kbps.cg.1min\" WHERE deliveryservice='" + which + "' and cachegroup != 'all' and $timeFilter GROUP BY time(60s), cachegroup",
              "rawQuery": true,
              "alias": "$tag_cachegroup"
            }
          ],
          "timeFrom": null,
          "timeShift": null,
          "title": "bandwidth by cachegroup",
          "tooltip": {
            "shared": true,
            "value_type": "individual"
          },
          "type": "graph",
          "x-axis": true,
          "y-axis": true,
          "y_formats": [
            "bps"
          ]
        }
      ],
      "title": "bwByCg"
  }
  );
}
return dashboard;
