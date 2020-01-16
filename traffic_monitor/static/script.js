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

const millisecondsInSecond = 1000;
const kilobitsInGigabit = 1000000;
const kilobitsInMegabit = 1000;

/**
 * This is the index of the latest event already catalogued by the UI. TM doesn't
 * provide a way to fetch event logs "older than x" etc., so this is how we keep
 * track of what we've seen.
*/
var lastEvent = 0;

/**
 * Adds a comma for every 3rd power of ten in a number represented as a string.
 * e.g. numberStrWithCommas("100000") outputs "100,000".
 * (source: http://stackoverflow.com/a/2901298/292623)
 * @param x A string containing a number which will be formatted.
*/
function numberStrWithCommas(x) {
	return x.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

/**
 * ajax performs an XMLHttpRequest and passes the result to a function - *only if the response code is EXACTLY 200*.
 * @param endpoint An API endpoint relative to the root of the TM webserver e.g. /publish/DsStats
 * @param f A function that takes a single argument which will be the entire response body as a string.
 */
function ajax(endpoint, f) {
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = () => {
		if (xhttp.readyState === 4 && xhttp.status === 200) {
			f(xhttp.responseText);
		}
	};
	xhttp.open("GET", endpoint, true);
	xhttp.send();
}

function getCacheCount() {
	ajax("/api/cache-count", function(r) {
		document.getElementById("cache-count").innerHTML = r;
	});
}

function getCacheAvailableCount() {
	ajax("/api/cache-available-count", function(r) {
		document.getElementById("cache-available").innerHTML = r;
	});
}

function getBandwidth() {
	ajax("/api/bandwidth-kbps", function(r) {
		document.getElementById("bandwidth").innerHTML = numberStrWithCommas((parseFloat(r) / kilobitsInGigabit).toFixed(2));
	});
}

function getBandwidthCapacity() {
	ajax("/api/bandwidth-capacity-kbps", function(r) {
		document.getElementById("bandwidth-capacity").innerHTML = numberStrWithCommas((r / kilobitsInGigabit).toString());
	});
}

function getCacheDownCount() {
	ajax("/api/cache-down-count", function(r) {
		document.getElementById("cache-down").innerHTML = r;
	});
}

function getVersion() {
	ajax("/api/version", function(r) {
		document.getElementById("version").innerHTML = r;
	});
}

function getTrafficOpsUri() {
	ajax("/api/traffic-ops-uri", function(r) {
		document.getElementById("source-uri").innerHTML = "<a href='" + r + "'>" + r + "</a>";
	});
}


function getTrafficOpsCdn() {
	ajax("/publish/ConfigDoc", function(r) {
		var j = JSON.parse(r);
		document.getElementById("cdn-name").innerHTML = j.cdnName;
	});
}

/**
 * Fetches the event log from TM and updates the "Event Log" table with the new
 * results.
*/
function getEvents() {
	/// \todo add /api/events-since/{index} (and change Traffic Monitor to keep latest
	ajax("/publish/EventLog", function(r) {
		const events = JSON.parse(r).events || [];
		for (const event of events.slice(lastEvent+1)) {
			lastEvent = event.index
			const row = document.getElementById("event-log").insertRow(0);

			row.insertCell(0).textContent = event.name;
			row.insertCell(1).textContent = event.type;

			const cell = row.insertCell(2);
			if(event.isAvailable) {
				cell.textContent = "available";
			} else {
				cell.textContent = "offline";
				row.classList.add("error");
			}

			row.insertCell(3).textContent = event.description;
			row.insertCell(4).textContent = new Date(event.time * 1000).toISOString();
		}
	});
}

/**
 * Fetches the current cache server states and statistics from TM and updates
 * the "Cache States" table with the results - replacing the current content.
*/
function getCacheStates() {
	ajax("/api/cache-statuses", function(r) {
		const servers = new Map(Object.entries(JSON.parse(r)));
		const table = document.createElement('TBODY');
		table.id = "cache-states"

		for (const [serverName, server] of servers) {
			const row = table.insertRow(0);

			row.insertCell(0).textContent = serverName;
			row.insertCell(1).textContent = server.type || "UNKNOWN";
			row.insertCell(2).textContent = server.status.indexOf("ONLINE") !== 0 ? server.ipv4_available : "N/A";
			row.insertCell(3).textContent = server.status.indexOf("ONLINE") !== 0 ?  server.ipv6_available : "N/A";
			row.insertCell(4).textContent = server.status || "";
			if (Object.prototype.hasOwnProperty.call(server, "status")) {
				if (server.status.indexOf("ADMIN_DOWN") !== -1 || server.status.indexOf("OFFLINE") !== -1) {
					row.classList.add("warning");
				} else if (!server.combined_available && server.status.indexOf("ONLINE") !== 0) {
					row.classList.add("error");
				} else if (server.status.indexOf(" availableBandwidth") !== -1) {
					row.classList.add("error");
				}
			}

			row.insertCell(5).textContent = server.load_average || "";
			row.insertCell(6).textContent = server.query_time_ms || "";
			row.insertCell(7).textContent = server.health_time_ms || "";
			row.insertCell(8).textContent = server.stat_time_ms || "";
			row.insertCell(9).textContent = server.health_span_ms || "";
			row.insertCell(10).textContent = server.stat_span_ms || "";

			if (Object.prototype.hasOwnProperty.call(server, "bandwidth_kbps")) {
				const kbps = (server.bandwidth_kbps / kilobitsInMegabit).toFixed(2);
				const max = numberStrWithCommas((server.bandwidth_capacity_kbps / kilobitsInMegabit).toFixed(0));
				row.insertCell(11).textContent = `${kbps} / ${max}`;
			} else {
				row.insertCell(11).textContent = "N/A";
			}

			row.insertCell(12).textContent = server.connection_count || "N/A";
		}


		const oldtable = document.getElementById("cache-states");
		oldtable.parentNode.replaceChild(table, oldtable);

	});
}

/**
 * dsDisplayFloat takes a float, and returns the string to display. For nonzero values, it returns two decimal places.
 * For zero values, it returns an empty string, to make nonzero values more visible.
 * @param f The floating point number to format
*/
const dsDisplayFloat = (f) => { return f === 0 ? "" : f.toFixed(2); }

/**
 * Attempts to extract data from a deliveryService object, but falls back on "N/A" if it doesn't have that
 * property.
 * @param ds The Delivery Service from which to extract data
 * @param prop The property being extracted. Technically, the extracted property is dsDisplayFloat(parseFloat(ds[prop][0].value)).
*/
function getDSProperty(ds, prop) {
	try {
		return dsDisplayFloat(parseFloat(ds[prop][0].value));
	} catch (e) {
		console.error(e);
	}
	return "N/A";
}

/**
 * Fetches the current Delivery Service stats from TM and updates the "Delivery Service States"
 * table with the results - replacing the current content.
*/
function getDsStats() {
	var now = Date.now();

	/// \todo add /api/delivery-service-stats which only returns the data needed by the UI, for efficiency
	ajax("/publish/DsStats", function(r) {
		const deliveryServices = new Map(Object.entries(JSON.parse(r)));
		const table = document.createElement('TBODY');
		table.id = "deliveryservice-stats";

		for (const [dsName, deliveryService] of deliveryServices) {
			const row = table.insertRow(0);
			const available = !deliveryService.isAvailable || !deliveryService.isAvailable[0] || !deliveryService.isAvailable[0].value === "true";
			if (available) {
				row.classList.add("error");
			}

			row.insertCell(0).textContent = dsName;
			row.insertCell(1).textContent = available ? "available" : `unavailable - ${deliveryService["error-string"][0].value}`;
			row.insertCell(2).textContent = (Object.prototype.hasOwnProperty.call(deliveryService, "caches-reporting") &&
			                                 Object.prototype.hasOwnProperty.call(deliveryService, "caches-available") &&
			                                 Object.prototype.hasOwnProperty.call(deliveryService, "caches-configured")) ?
			                                 `${deliveryService['caches-reporting'][0].value} / ${deliveryService['caches-available'][0].value} / ${deliveryService['caches-configured'][0].value}`;

			row.insertCell(3).textContent = Object.prototype.hasOwnProperty.call(deliveryService, "total.kbps") ? (jds[deliveryService]['total.kbps'][0].value / kilobitsInMegabit).toFixed(2) : "N/A";
			row.insertCell(4).textContent = getDSProperty(deliveryService, "total.tps_total");
			row.insertCell(5).textContent = getDSProperty(deliveryService, "total.tps_2xx");
			row.insertCell(6).textContent = getDSProperty(deliveryService, "total.tps_3xx");
			row.insertCell(7).textContent = getDSProperty(deliveryService, "total.tps_4xx");
			row.insertCell(8).textContent = getDSProperty(deliveryService, "total.tps_5xx");
			row.insertCell(9);
			// \todo implement disabled locations
		}

		const oldtable = document.getElementById("deliveryservice-stats");
		oldtable.parentNode.replaceChild(table, oldtable);
	});
}

/**
 * Fetches not only the "Cache States" but also the aggregate cache server statistics used in the
 * informational section at the top of the page.
 */
function getCacheStatuses() {
	getCacheCount();
	getCacheAvailableCount();
	getCacheDownCount();
	getCacheStates();
}

/**
 * Fetches the metadata information used at the very top of the page.
 */
function getTopBar() {
	getVersion();
	getTrafficOpsUri();
	getTrafficOpsCdn();
	getCacheStatuses();
}

/**
 * Runs immediately after content is loaded, fetching initial information and setting intervals for
 * for gathering other data.
 */
function init() {
	getTopBar();
	setInterval(getCacheCount, 4755);
	setInterval(getCacheAvailableCount, 4800);
	setInterval(getBandwidth, 4621);
	setInterval(getBandwidthCapacity, 4591);
	setInterval(getCacheDownCount, 4832);
	setInterval(getVersion, 10007); // change to retry on failure, and only do on startup
	setInterval(getTrafficOpsUri, 10019); // change to retry on failure, and only do on startup
	setInterval(getTrafficOpsCdn, 10500); // change to retry on failure, and only do on startup
	setInterval(getEvents, 2004); // change to retry on failure, and only do on startup
	setInterval(getCacheStatuses, 5009);
	setInterval(getDsStats, 4003);
}
