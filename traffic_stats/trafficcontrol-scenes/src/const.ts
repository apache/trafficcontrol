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

import pluginJson from "./plugin.json";

export const PLUGIN_BASE_URL = `/a/${pluginJson.id}`;

export const ROUTES = {
	cacheGroup: "cache-group",
	deliveryService: "delivery-service",
	server: "server",
};

export const PROMETHEUS_DATASOURCE_REF = {
	type: "prometheus",
	uid: "prometheus",
};

export const INFLUXDB_DATASOURCES_REF = {
	cacheStats: {
		type: "influxdb",
		uid: "cache_stats",
	},
	dailyStats: {
		type: "influxdb",
		uid: "daily_stats",
	},
	deliveryServiceStats: {
		type: "influxdb",
		uid: "deliveryservice_stats",
	},
	telegraf: {
		type: "influxdb",
		uid: "telegraf",
	},
};
