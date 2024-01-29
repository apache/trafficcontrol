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

import { PanelBuilders, SceneQueryRunner, VizPanel } from "@grafana/scenes";
import { INFLUXDB_DATASOURCES_REF } from "const";

export const getCPUPanel = (): VizPanel => {
	const defaultQuery = {
		alias: "$col",
		query: "SELECT mean(\"usage_system\") AS \"cpu_system\", mean(\"usage_iowait\") AS \"cpu_iowait\"," +
			" mean(\"usage_user\") AS \"cpu_user\", mean(\"usage_guest\") AS \"cpu_guest\"," +
			" mean(\"usage_steal\") AS \"cpu_steal\" FROM \"cpu\"" +
			" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)",
		rawQuery: true,
		refId: "A",
		resultFormat: "time_series",
	};

	const qr = new SceneQueryRunner({
		datasource: INFLUXDB_DATASOURCES_REF.telegraf,
		queries: [defaultQuery],
	});

	return PanelBuilders.timeseries()
		.setTitle("CPU Usage")
		.setData(qr)
		.setUnit("%")
		.setCustomFieldConfig("spanNulls", true)
		.setCustomFieldConfig("fillOpacity", 20)
		.build();
};
