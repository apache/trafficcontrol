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

"use strict";

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
		document.getElementById("cache-count").textContent = r;
	});
}

function getCacheAvailableCount() {
	ajax("/api/cache-available-count", function(r) {
		document.getElementById("cache-available").textContent = r;
	});
}

function getBandwidth() {
	ajax("/api/bandwidth-kbps", function(r) {
		document.getElementById("bandwidth").textContent = numberStrWithCommas((parseFloat(r) / kilobitsInGigabit).toFixed(2));
	});
}

function getBandwidthCapacity() {
	ajax("/api/bandwidth-capacity-kbps", function(r) {
		document.getElementById("bandwidth-capacity").textContent = numberStrWithCommas((r / kilobitsInGigabit).toString());
	});
}

function getCacheDownCount() {
	ajax("/api/cache-down-count", function(r) {
		document.getElementById("cache-down").textContent = r;
	});
}

function getVersion() {
	ajax("/api/version", function(r) {
		document.getElementById("version").textContent = r;
	});
}

function getTrafficOpsUri() {
	ajax("/api/traffic-ops-uri", function(r) {
		// This used to be done by setting the element's `innerHTML`, but that doesn't remove
		// the child nodes. They're orphaned, but continue to take up memory.
		const link = document.createElement('A');
		link.href = r;
		link.textContent = r;
		const sourceURISpan = document.getElementById('source-uri');
		while (sourceURISpan.lastChild) {
			sourceURISpan.removeChild(sourceURISpan.lastChild);
		}
		sourceURISpan.appendChild(link);
	});
}

function getTrafficOpsCdn() {
	const cdnName = document.getElementById("cdn-name");
	const discIconContainer = document.getElementById("icon-disc-holder");
	ajax("/publish/ConfigDoc", function(r) {
		let opsConfig = JSON.parse(r);
		cdnName.textContent = opsConfig.cdnName || "unknown";
		if (opsConfig.usingDummyTO === false) {
		    discIconContainer.hidden = true;
		} else {
			discIconContainer.hidden = false;
		}
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
			lastEvent = event.index;
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
	function parseIPAvailable(server, ipField) {
		return server.status.indexOf("ONLINE") !== 0 ? server[ipField] : "N/A";
	}

	function parseBandwidth(server) {
		if (Object.prototype.hasOwnProperty.call(server, "bandwidth_kbps") &&
				Object.prototype.hasOwnProperty.call(server, "bandwidth_capacity_kbps")) {
			const kbps = (server.bandwidth_kbps / kilobitsInMegabit).toFixed(2);
			const max = numberStrWithCommas((server.bandwidth_capacity_kbps / kilobitsInMegabit).toFixed(0));
			return `${kbps} / ${max}`;
		} else {
			return "N/A";
		}
	}

	ajax("/api/cache-statuses", function(r) {
		let serversArray = Object.entries(JSON.parse(r)).sort((serverTupleA, serverTupleB) => {
			return -1*serverTupleA[0].localeCompare(serverTupleB[0]);
		});
		const servers = new Map(serversArray);

		const oldtable = document.getElementById("cache-states");
		const table = document.createElement('TBODY');
		const interfaceTableTemplate = document.getElementById("interface-template").content.children[0];
		const interfaceRowTemplate = document.getElementById("interface-row-template").content.children[0];
		const cacheStatusRowTemplate = document.getElementById("cache-status-row-template").content.children[0];
		table.id = oldtable.id;

		// Match visibility of interface tables based on previous table
		const interfaceRows = oldtable.querySelectorAll(".encompassing-row");
		let openCachesByName = new Set();
		for(const row of interfaceRows) {
			if(row.classList.contains("visible")){
				openCachesByName.add(row.querySelector(".sub-table").getAttribute("server-name"));
			}
		}

		for (const [serverName, server] of servers) {
			const row = cacheStatusRowTemplate.cloneNode(true);
			const cacheRowChildren = row.children;
			const indicatorDiv = cacheRowChildren[0].children[0];

			if (Object.prototype.hasOwnProperty.call(server, "status") &&
					Object.prototype.hasOwnProperty.call(server, "combined_available")) {
				if (server.status.indexOf("ADMIN_DOWN") !== -1 || server.status.indexOf("OFFLINE") !== -1) {
					row.classList.add("warning");
				} else if (!server.combined_available && server.status.indexOf("ONLINE") !== 0) {
					row.classList.add("error");
				} else if (server.status.indexOf(" availableBandwidth") !== -1) {
					row.classList.add("error");
				}
			}

			cacheRowChildren[1].textContent = serverName;
			cacheRowChildren[2].textContent = server.type || "UNKNOWN";
			cacheRowChildren[3].textContent = parseIPAvailable(server, "ipv4_available");
			cacheRowChildren[4].textContent = parseIPAvailable(server, "ipv6_available");
			cacheRowChildren[5].textContent = server.status || "";
			cacheRowChildren[6].textContent = server.load_average || "";
			cacheRowChildren[7].textContent = server.query_time_ms || "";
			cacheRowChildren[8].textContent = server.health_time_ms || "";
			cacheRowChildren[9].textContent = server.stat_time_ms || "";
			cacheRowChildren[10].textContent = server.health_span_ms || "";
			cacheRowChildren[11].textContent = server.stat_span_ms || "";
			cacheRowChildren[12].textContent = parseBandwidth(server);
			cacheRowChildren[13].textContent = server.connection_count || "N/A";
			table.prepend(row);

			const encompassingRow = table.insertRow(1);
			encompassingRow.classList.add("encompassing-row");
			const encompassingCell = encompassingRow.insertCell(0);
			const interfaceTable = interfaceTableTemplate.cloneNode(true);
			encompassingCell.colSpan = 14;
			// Add interfaces
			if (Object.prototype.hasOwnProperty.call(server, "interfaces")) {
				let interfacesArray = Object.entries(server.interfaces).sort((interfaceTupleA, interfaceTupleB) => {
					return -1 * interfaceTupleA[0].localeCompare(interfaceTupleB[0]);
				});
				server.interfaces = new Map(interfacesArray);
				interfaceTable.removeAttribute("id");
				// To determine what cache this interface table belongs to
				// used to ensure servers that were expanded remain expanded when refreshing the data.
				interfaceTable.setAttribute("server-name", serverName);
				const interfaceBody = interfaceTable.querySelector(".interface-content");
				for (const [interfaceName, stat] of server.interfaces) {
					const interfaceRow = interfaceRowTemplate.cloneNode(true);
					const cells = interfaceRow.children;

					cells[0].textContent = interfaceName;
					cells[1].textContent = parseIPAvailable(stat, "ipv4_available");
					cells[2].textContent = parseIPAvailable(stat, "ipv6_available");
					cells[3].textContent = stat.status || "";
					cells[4].textContent = parseBandwidth(stat);
					cells[5].textContent = stat.connection_count || "N/A";

					if (Object.prototype.hasOwnProperty.call(stat, "available")) {
					    if (stat.available === false) {
							interfaceRow.classList.add("error");
						}
					}

					interfaceBody.prepend(interfaceRow);
				}
				row.onclick = function() {
					if(encompassingRow.classList.contains("visible")) {
						encompassingRow.classList.remove("visible");
						indicatorDiv.classList.remove("down");
					}
					else {
						encompassingRow.classList.add("visible");
						indicatorDiv.classList.add("down");
					}
				};
			}
			else {
				indicatorDiv.classList.add("hidden");
			}

			encompassingCell.appendChild(interfaceTable);
			// Row was unhidden previously
			if(openCachesByName.has(serverName)) {
				row.click();
			}
		}

		oldtable.parentNode.replaceChild(table, oldtable);
	});
}

/**
 * dsDisplayFloat takes a float, and returns the string to display. For nonzero values, it returns two decimal places.
 * For zero values, it returns an empty string, to make nonzero values more visible.
 * @param f The floating point number to format
*/
const dsDisplayFloat = (f) => { return f === 0 ? "" : f.toFixed(2); };

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
	/// \todo add /api/delivery-service-stats which only returns the data needed by the UI, for efficiency
	ajax("/publish/DsStats", function(r) {
		const deliveryServices = new Map(Object.entries(JSON.parse(r).deliveryService).sort((dsTupleA, dsTupleB) => {
			return -1 * dsTupleA[0].localeCompare(dsTupleB[0]);
		}));

		const table = document.createElement('TBODY');
		table.id = "deliveryservice-stats";

		for (const [dsName, deliveryService] of deliveryServices) {
			const row = table.insertRow(0);
			const available = deliveryService.isAvailable && deliveryService.isAvailable[0] && deliveryService.isAvailable[0].value === "true";
			if (!available) {
				row.classList.add("error");
			}

			row.insertCell(0).textContent = dsName;
			row.insertCell(1).textContent = available ? "available" : `unavailable - ${deliveryService["error-string"][0].value}`;
			row.insertCell(2).textContent = (Object.prototype.hasOwnProperty.call(deliveryService, "caches-reporting") &&
											 Object.prototype.hasOwnProperty.call(deliveryService, "caches-available") &&
											 Object.prototype.hasOwnProperty.call(deliveryService, "caches-configured")) ?
											 `${deliveryService['caches-reporting'][0].value} / ${deliveryService['caches-available'][0].value} / ${deliveryService['caches-configured'][0].value}` : "N/A";

			row.insertCell(3).textContent = Object.prototype.hasOwnProperty.call(deliveryService, "total.kbps") ? (deliveryService['total.kbps'][0].value / kilobitsInMegabit).toFixed(2) : "N/A";
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

window.addEventListener('load', init);
