/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


var myip = "127.0.0.1";
var myport = 80;
var config_url = "https://traffic-ops.com/api/4.0/cdns/cdn-name/snapshot";
var to_user = "";
var to_password = "";
var to_login_api = "/api/4.0/user/login";
var simulator_ua = "ATS Simulator/node.js " + process.version;

// first argument to follow node ats_sim.js
if (process.argv[2]) {
	config_url = process.argv[2];
}

// second argument to follow node ats_sim.js
if (process.argv[3]) {
	myport = process.argv[3];
}

if (process.argv[4]) {
	to_user = process.argv[4];
}

if (process.argv[5]) {
	to_password = process.argv[5];
}

var to_credentials = JSON.stringify({
	"u": to_user,
	"p": to_password,
});

var errorRateDenominator = 0;
var debug = 0;

// https://${tmHostname}/api/${apiVersion}/cdns/${cdnName}/snapshot
if (debug) console.log("point traffic_monitor::tm.crConfig.json.polling.url to: " + config_url);

var date = new Date();
var starttime = date.getTime();
var timeseed = 1383769987010;

console.log("Started " + date);

var protocol = (/^https/.test(config_url) == true ? "https" : "http");
var client_http = require(protocol);
var cr_config = '';
var tld = '';
var Url = require('url');
var stats = new Object();
var to_url = Url.parse(config_url, true);

process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
auth_get_config();

function auth_get_config() {
	var to_cookie = "";
	var options = {
		protocol: to_url.protocol,
		hostname: to_url.hostname,
		port: to_url.port,
		path: to_login_api,
		method: 'POST',
		headers: {
			"Content-Type": "application/x-www-form-urlencoded",
			"Content-Length": Buffer.byteLength(to_credentials),
			"User-Agent": simulator_ua
		}
	};

	request = client_http.request(options, function(res) {
		var cookie = res.headers["set-cookie"];

		if (cookie) {
			cookie.forEach(
				function (cookieMonster) {
					to_cookie += cookieMonster;
				}
			);
		}

		if (to_cookie) {
			if (debug) console.log("Cookie: " + to_cookie);
		}

		var body = '';

		res.on('data', function(chunk) {
			body += chunk;
		});

		res.on('end', function() {
			response = JSON.parse(body);

			if (res.statusCode == 200) {
				console.log("Successfully authenticated to Traffic Ops");
				get_config(options, to_cookie);
			} else {
				console.log("Authentication response: " + JSON.stringify(response));
				throw new Error("Authentication failed.");
			}
		});
	}).on('error', function(e) {
		console.log("Got error: ", e);
	});

	request.write(to_credentials);
	request.end();
}

function get_config(options, cookie) {
	options.method = "GET";
	options.path = to_url.path;
	options.headers = {
		"Cookie": cookie,
		"User-Agent": simulator_ua
	};

	request = client_http.request(options, function(res) {
		var body = '';

		res.on('data', function(chunk) {
			body += chunk;
		});

		res.on('end', function() {
			cr_config = JSON.parse(body);
			tld = cr_config.config['domain_name'];
			console.log("CDN TLD: " + tld);

			for(var cid in cr_config.contentServers) {
				if(cid.indexOf("atsec-sim")!=-1) {
					cr_config.contentServers[cid].queryIp = myip;
					cr_config.contentServers[cid].port = myport;
				}
			}

			cr_config
		});
	}).on('error', function(e) {
		console.log("Got error: ", e);
	});

	request.end();
}

var http = require('http');

http.createServer(function (request, response) {
	if (request.url == "/favicon.ico") {
		response.writeHead(404);
		response.end();
		return;
	}

	// https://${tmHostname}/api/${apiVersion}/cdns/${cdnName}/snapshot
	if (request.url.indexOf("/snapshot") != -1) {
		console.log("Delivering CRConfig.json");
		response.end(JSON.stringify(cr_config, null, 4));
		return;
	}

	if (errorRateDenominator && rand(errorRateDenominator) == 1) {
		response.writeHead(404);
		response.end();
		return;
	}

	var objToJson = { };
	objToJson.cr_config = cr_config;
	var str = JSON.stringify(getData(request));
	response.writeHead(200, {'Content-Type': 'application/json'});
	response.end(str);
}).listen(myport);

console.log("HTTP Listener started on port " + myport);

var util = require('util');
var crypto = require('crypto')
var shasum = crypto.createHash('sha1');

function rand(n) {
	return Math.floor((Math.random() * n) + 1);
}

Number.prototype.zeroPad = function(numZeros) {
	var n = Math.abs(this);
	var zeros = Math.max(0, numZeros - n.toString().length);
	var zeroString = Math.pow(10,zeros).toString().substr(1);
	if( this < 0 ) {
		zeroString = '-' + zeroString;
	}

	return zeroString+n;
}

function getData(request) {
	var url_parts = Url.parse(request.url, true);
	var query_params = url_parts.query;

	if (!request.headers.host) {
		return;
	}

	if (debug) console.log(request.url);
	if (debug) console.log(query_params);
	if (debug) console.log(query_params.a);

	var ret = {};
	var d = new Date();
	var n = d.getTime();
	var unixtime = d.getTime();
	var basetime = unixtime - timeseed;
	var server_name = request.headers.host.split(':')[0];
	server_name = server_name.split(".")[0];
	var hex_string = server_name;
	var hex = crypto.createHash('sha1').update(server_name).digest('hex');
	hex = parseInt(hex, 16);
	ret.simulator = 1;
	ret.time = unixtime;
	ret.hex = hex;
	var int1  = hex % 4;
	var int2  = hex % 40;
	var int3  = hex % 400;
	var int4  = hex % 4000;
	var int5  = hex % 40000;
	var int6  = hex % 400000;
	var int7  = hex % 4000000;
	var int8  = hex % 40000000;
	var int9  = hex % 40000000;
	var int10 = hex % 400000000;

	var hit_fresh = unixtime / 2;
	var hit_fresh_process = hit_fresh;
	var hit_revalidated = (int10 / 2) + unixtime + int1;
	var miss_cold = (int10 / 2) + unixtime + int3;
	var miss_not_cacheable = (int7 / 3) + unixtime + int1;
	var miss_changed = (unixtime / 10000000) + int1;
	var miss_client_no_cache = 0;
	var aborts = (unixtime / 10000000) + int2 * 2;
	var possible_aborts = 0;
	var connect_failed = (unixtime / 1000000) + int2 * 4 + int1 * 2;
	var other = ((unixtime / 1000000) + int3) / 100;
	var unclassified = 0;
	var write_bytes = unixtime + int4 * int5;
	var current_client_connections = Math.round(int5 * Math.random());
	var bytes_used = write_bytes / 20 + int2 * int5;
	var bytes_total = bytes_used + int10;
	var v1_bytes_used = bytes_used;
	var v1_bytes_total = bytes_total;
	var load_avg_1 = (8 * Math.random()).toFixed(2);
	var load_avg_5 = (load_avg_1 / 2).toFixed(2);
	var load_avg_15 = (load_avg_1 / 3).toFixed(2);
	var running_procs = (15 * rand(int1));
	var total_procs = 2 * running_procs + (int1 + 2) * 2;
	var last_proc_id = int5;
	var proc_loadavg = util.format("%d %d %d %d/%d %d", load_avg_1, load_avg_5, load_avg_15, running_procs, total_procs, last_proc_id);
	var if_rbytes = basetime * 500;
	var if_rpackets = if_rbytes / 875;
	var if_rmcast = if_rpackets / int4;
	var if_tbytes = if_rbytes / 3;
	var if_tpackets = if_tbytes / 1500;
	var proc_net_dev = 'bond0';

	if (query_params['inf.name']) {
		proc_net_dev = query_params['inf.name'];
	}

	proc_net_dev += util.format(":%d %d 0 0 0 0 0 %d %d %d 0 0 0 0 0 0", if_rbytes.toFixed(0), if_rpackets.toFixed(0), if_rmcast.toFixed(0), if_tbytes.toFixed(0), if_tpackets.toFixed(0));
	var ds_bytes = {};
	var ds_bytes_in = {};
	var ds_200s = {};
	var ds_400s = {};
	var ds_500s = {};

	for (var i=1; i<=10; i++) {
		var rand1 = rand(int1);
		ds_bytes[i] = if_tbytes / 10;
		ds_bytes_in[i] = if_rbytes / 10;
		var tps = basetime / 10000;
		ds_200s[i] = tps * 0.95;
		ds_400s[i] = tps * 0.04;
		ds_500s[i] = tps * 0.01;
	}

	if (query_params.ds && query_params.p && query_params.v && query_params.m) {
		var these_stats = new Object();
		these_stats.delivery_service = query_params.ds;
		these_stats.parameter = query_params.p;
		these_stats.value = parseInt(query_params.v);
		these_stats.multiplier = parseFloat(query_params.m);
		stats[server_name] = these_stats;
	} else if (query_params.stats == 0) {
		delete stats[server_name];
	}

	if (debug) console.log(stats);

//	#proc.net.dev: "bond0:181566812618839 43321349767 0 0 0 0 0 4710035 517574148613675 34658736727 0 0 0 0 0 0"

	ret.ats = {};
	ret.ats["proxy.process.http.transaction_counts.hit_fresh"] = hit_fresh;
	ret.ats["proxy.process.http.transaction_counts.hit_fresh.process"] = hit_fresh_process;
	ret.ats["proxy.process.http.transaction_counts.hit_revalidated"] = hit_revalidated;
	ret.ats["proxy.process.http.transaction_counts.miss_cold"] = miss_cold;
	ret.ats["proxy.process.http.transaction_counts.miss_not_cacheable"] = miss_not_cacheable;
	ret.ats["proxy.process.http.transaction_counts.miss_changed"] = miss_changed;
	ret.ats["proxy.process.http.transaction_counts.miss_client_no_cache"] = miss_client_no_cache;
	ret.ats["proxy.process.http.transaction_counts.errors.aborts"] = aborts;
	ret.ats["proxy.process.http.transaction_counts.errors.possible_aborts"] = possible_aborts;
	ret.ats["proxy.process.http.transaction_counts.errors.connect_failed"] = connect_failed;
	ret.ats["proxy.process.http.transaction_counts.errors.other"] = other;
	ret.ats["proxy.process.http.transaction_counts.other.unclassified"] = unclassified;
	ret.ats["proxy.process.net.write_bytes"] = write_bytes;
	ret.ats["proxy.process.http.current_client_connections"] = current_client_connections;
	ret.ats["proxy.process.cache.bytes_used"] = bytes_used;
	ret.ats["proxy.process.cache.bytes_total"] = bytes_total;
	ret.ats["proxy.process.cache.volume_1.bytes_used"] = v1_bytes_used;
	ret.ats["proxy.process.cache.volume_1.bytes_total"] = v1_bytes_total;

	if (cr_config && cr_config.contentServers && server_name in cr_config.contentServers) {
		if (debug) console.log("Attempting to build stats for " + server_name);

		for (var delivery_service in cr_config.contentServers[server_name]['deliveryServices']) {
			var data = new Object();
			data['out_bytes'] = ds_bytes[rand(10)].toFixed(0);
			data['in_bytes'] = ds_bytes_in[rand(10)].toFixed(0);
			data['status_2xx'] = ds_200s[rand(10)].toFixed(0);
			data['status_4xx'] = ds_400s[rand(10)].toFixed(0);
			data['status_5xx'] = ds_500s[rand(10)].toFixed(0);

			var n = i.zeroPad(2);
			var ds_fqdn = cr_config.contentServers[server_name]['deliveryServices'][delivery_service][0];

			if (stats.hasOwnProperty(server_name) && stats[server_name].hasOwnProperty("delivery_service")) {
				var ds_tld = ds_fqdn.substring(ds_fqdn.indexOf(".") + 1);

				if (stats[server_name].delivery_service == ds_tld && stats[server_name].parameter in data) {
					var this_request = parseInt(d.getTime());

					if (stats[server_name].hasOwnProperty("last_request") && stats[server_name].hasOwnProperty("last_value")) {
						/*
						I'm requiring bignum here to ensure that it's only required when running bandwidth simulations
						if you need to install it:
						  npm install bignum
						..which will install it locally.. to install it globally add -g to the command
						if you have problems loading it due to where you're running ats_sim.js:
						  export NODE_PATH=/path/to/the/node_modules
						For example:
						  export NODE_PATH=/usr/local/lib/node_modules
						..and the bignum directory is under /usr/local/lib/node/modules
						See: https://npmjs.org/package/bignum
						..and: http://nodejs.org/api/modules.html#modules_loading_from_node_modules_folders
						*/
						var bignum = require('bignum');
						var max = bignum(18446744073709551615); // uint64 max
						var seconds_elapsed = Math.ceil((this_request - stats[server_name].last_request) / 1000); // convert to seconds and round up
						var last_value = bignum(stats[server_name].last_value);
						var current_value = bignum(stats[server_name].value * seconds_elapsed * stats[server_name].multiplier).add(last_value);

						if (debug) console.log(seconds_elapsed + " second(s) elapsed since the last request");

						if (current_value.cmp(max) > 0) {
							if (debug) console.log("Accounting for uint64 rollover; " + current_value.toString() + " > " + max.toString());
							current_value = current_value.sub(max);
						}

						if (debug) console.log("Overriding " + ds_tld + "'s " + stats[server_name].parameter + " with " + current_value.toString());
						data[stats[server_name].parameter] = current_value.toString(); // we'll reconvert this to a bignum next time through
					}

					stats[server_name].last_value = data[stats[server_name].parameter];
					stats[server_name].last_request = this_request;
				} else {
					if (debug) console.log("Unable to find " + stats[server_name].parameter + " for " + ds_tld);
				}
			}

			for (var key in data) {
				ret.ats["plugin.remap_stats." + ds_fqdn + "." + key] = parseInt(data[key]);
			}
		}
	} else {
		if (debug) console.log(server_name + " was not found in cr_config; returning default stats");

		for(var i=1; i<=10; i++) {
			var n = i.zeroPad(2);
			ret.ats["plugin.remap_stats." + server_name  + ".omg-" + n + tld+ ".out_bytes"] = parseInt(ds_bytes[i].toFixed(0));
			ret.ats["plugin.remap_stats." + server_name  + ".omg-" + n + tld+ ".in_bytes"] = parseInt(ds_bytes_in[i].toFixed(0));
			ret.ats["plugin.remap_stats." + server_name  + ".omg-" + n + tld+ ".status_2xx"] = parseInt(ds_200s[i].toFixed(0));
			ret.ats["plugin.remap_stats." + server_name  + ".omg-" + n + tld+ ".status_4xx"] = parseInt(ds_400s[i].toFixed(0));
			ret.ats["plugin.remap_stats." + server_name  + ".omg-" + n + tld+ ".status_5xx"] = parseInt(ds_500s[i].toFixed(0));
		}
	}

	ret.ats["server"] = "5.3.2-dev";

	ret.system = {};
	ret.system["inf.name"] = query_params["inf.name"];
	ret.system["inf.speed"] = 10000;
	ret.system["proc.net.dev"] = proc_net_dev;
	ret.system["proc.loadavg"] = proc_loadavg;
	ret.system["something"] = "here";

	return ret;
}
