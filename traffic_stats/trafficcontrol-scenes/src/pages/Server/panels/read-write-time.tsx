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

export const getReadWriteTimePanel = (): VizPanel => {
	const defaultQueries = [
		{
			alias: "$col",
			query: "SELECT non_negative_derivative(sum(\"read_time\"), 10s) AS \"read_time\" FROM \"diskio\" " +
				" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)",
			rawQuery: true,
			refId: "A",
			resultFormat: "time_series",
		},
		{
			alias: "$col",
			query: "SELECT non_negative_derivative(sum(\"write_time\"), 10s) AS \"write_time\" FROM \"diskio\" " +
				" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)",
			rawQuery: true,
			refId: "B",
			resultFormat: "time_series",
		},
	];

	const qr = new SceneQueryRunner({
		datasource: INFLUXDB_DATASOURCES_REF.telegraf,
		queries: [...defaultQueries],
	});

	return PanelBuilders.timeseries()
		.setTitle("Read/Write Time")
		.setData(qr)
		.setCustomFieldConfig("spanNulls", true)
		.setCustomFieldConfig("fillOpacity", 20)
		.setUnit("ns")
		.build();
};
