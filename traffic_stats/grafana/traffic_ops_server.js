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
  from: "now-6h",
  to: "now-60s"
};

var which = 'argName';

if(!_.isUndefined(ARGS.which)) {
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
              "measurement": "bandwidth",
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
              "measurement": "bandwidth",
              "tags": {},
              "query": "SELECT mean(value)*1000 FROM \"monthly\".\"connections.1min\" WHERE hostname='" + which + "'  and $timeFilter GROUP BY time(60s)",
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
    }
  );
}
return dashboard;
