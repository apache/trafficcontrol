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

// source: http://stackoverflow.com/a/2901298/292623
function numberStrWithCommas(x) {
	return x.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
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
 * This is the index of the latest event already catalogued by the UI. TM doesn't
 * provide a way to fetch event logs "older than x" etc., so this is how we keep
 * track of what we've seen.
*/
var lastEvent = 0;

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

const millisecondsInSecond = 1000;
const kilobitsInGigabit = 1000000;
const kilobitsInMegabit = 1000;

/**
 * dsDisplayFloat takes a float, and returns the string to display. For nonzero values, it returns two decimal places.
 * For zero values, it returns an empty string, to make nonzero values more visible.
 * @param f The floating point number to format
*/
const dsDisplayFloat = (f) => { return f === 0 ? "" : f.toFixed(2); }

function getDsStats() {
	var now = Date.now();

	/// \todo add /api/delivery-service-stats which only returns the data needed by the UI, for efficiency
	ajax("/publish/DsStats", function(r) {
		var j = JSON.parse(r);
		var jds = j.deliveryService
		var deliveryServiceNames = Object.keys(jds); //debug
		//decrementing for loop so DsNames are alphabetical A-Z
		//TODO allow for filtering of columns so this isn't necessary
			for (var i = deliveryServiceNames.length - 1; i >= 0; i--) {
			var deliveryService = deliveryServiceNames[i];

			if (!document.getElementById("deliveryservice-stats-" + deliveryService)) {
				var row = document.getElementById("deliveryservice-stats").insertRow(0); //document.createElement("tr");
				row.id = "deliveryservice-stats-" + deliveryService
				row.insertCell(0).id = row.id + "-delivery-service";
				row.insertCell(1).id = row.id + "-status";
				row.insertCell(2).id = row.id + "-caches-reporting";
				row.insertCell(3).id = row.id + "-bandwidth";
				row.insertCell(4).id = row.id + "-tps";
				row.insertCell(5).id = row.id + "-2xx";
				row.insertCell(6).id = row.id + "-3xx";
				row.insertCell(7).id = row.id + "-4xx";
				row.insertCell(8).id = row.id + "-5xx";
				row.insertCell(9).id = row.id + "-disabled-locations";
				document.getElementById(row.id + "-delivery-service").textContent = deliveryService;
				document.getElementById(row.id + "-delivery-service").style.whiteSpace = "nowrap";
				document.getElementById(row.id + "-caches-reporting").style.textAlign = "right";
				document.getElementById(row.id + "-bandwidth").style.textAlign = "right";
				document.getElementById(row.id + "-tps").style.textAlign = "right";
				document.getElementById(row.id + "-2xx").style.textAlign = "right";
				document.getElementById(row.id + "-3xx").style.textAlign = "right";
				document.getElementById(row.id + "-4xx").style.textAlign = "right";
				document.getElementById(row.id + "-5xx").style.textAlign = "right";
			}

			// \todo check that array has a member before dereferencing [0]
			if (jds[deliveryService].hasOwnProperty("isAvailable")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-status").textContent = jds[deliveryService]["isAvailable"][0].value == "true" ? "available" : "unavailable - " + jds[deliveryService]["error-string"][0].value;
			}
			if (jds[deliveryService].hasOwnProperty("caches-reporting") && jds[deliveryService].hasOwnProperty("caches-available") && jds[deliveryService].hasOwnProperty("caches-configured")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-caches-reporting").textContent = jds[deliveryService]['caches-reporting'][0].value + " / " + jds[deliveryService]['caches-available'][0].value + " / " + jds[deliveryService]['caches-configured'][0].value;
			}
			if (jds[deliveryService].hasOwnProperty("total.kbps")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-bandwidth").textContent = (jds[deliveryService]['total.kbps'][0].value / kilobitsInMegabit).toFixed(2);
			} else {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-bandwidth").textContent = "N/A";
			}
			if (jds[deliveryService].hasOwnProperty("total.tps_total")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-tps").textContent = dsDisplayFloat(parseFloat(jds[deliveryService]['total.tps_total'][0].value));
			} else {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-tps").textContent = "N/A";
			}
			if (jds[deliveryService].hasOwnProperty("total.tps_2xx")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-2xx").textContent = dsDisplayFloat(parseFloat(jds[deliveryService]['total.tps_2xx'][0].value));
			} else {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-2xx").textContent = "N/A";
			}
			if (jds[deliveryService].hasOwnProperty("total.tps_3xx")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-3xx").textContent = dsDisplayFloat(parseFloat(jds[deliveryService]['total.tps_3xx'][0].value));
			} else {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-3xx").textContent = "N/A";
			}
			if (jds[deliveryService].hasOwnProperty("total.tps_4xx")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-4xx").textContent = dsDisplayFloat(parseFloat(jds[deliveryService]['total.tps_4xx'][0].value));
			} else {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-4xx").textContent = "N/A";
			}
			if (jds[deliveryService].hasOwnProperty("total.tps_5xx")) {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-5xx").textContent = dsDisplayFloat(parseFloat(jds[deliveryService]['total.tps_5xx'][0].value));
			} else {
				document.getElementById("deliveryservice-stats-" + deliveryService + "-5xx").textContent = "N/A";
			}

			// \todo implement disabled locations

			var row = document.getElementById("deliveryservice-stats-" + deliveryService);
			if (jds[deliveryService]["isAvailable"][0].value == "true") {
				row.classList.add("stripes");
				row.classList.remove("error");
			} else {
				row.classList.add("error");
				row.classList.remove("stripes");
			}
		}
	})
}

function getCacheStatuses() {
	getCacheCount();
	getCacheAvailableCount();
	getCacheDownCount();
	getCacheStates();
}

function getTopBar() {
	getVersion();
	getTrafficOpsUri();
	getTrafficOpsCdn();
	getCacheStatuses();
}

function ajax(endpoint, f) {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (xhttp.readyState == 4 && xhttp.status == 200) {
			f(xhttp.responseText);
		}
	};
	xhttp.open("GET", endpoint, true);
	xhttp.send();
}
