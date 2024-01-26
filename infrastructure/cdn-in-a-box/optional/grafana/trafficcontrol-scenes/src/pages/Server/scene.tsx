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

import {
	SceneTimeRange,
	EmbeddedScene,
	SceneFlexLayout,
	SceneFlexItem,
	SceneControlsSpacer,
	SceneRefreshPicker,
	SceneTimePicker,
	QueryVariable,
	SceneVariableSet,
	VariableValueSelectors,
} from "@grafana/scenes";
import {
	getBandwidthPanel,
	getConnectionsPanel,
	getCPUPanel,
	getMemoryPanel,
	getLoadAveragePanel,
	getReadWriteTimePanel,
	getWrapCountPanel,
	getNetstatPanel,
} from "./panels";
import { INFLUXDB_DATASOURCES_REF } from "../../constants";

export function getServerScene() {
	const timeRange = new SceneTimeRange({
		from: "now-6h",
		to: "now",
	});

	const hostname = new QueryVariable({
		datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
		name: "hostname",
		query: 'SHOW TAG VALUES ON "cache_stats" FROM "monthly"."bandwidth" with key = "hostname"',
	});

	return new EmbeddedScene({
		$timeRange: timeRange,
		$variables: new SceneVariableSet({
			variables: [hostname],
		}),
		body: new SceneFlexLayout({
			direction: "column",
			children: [
				new SceneFlexItem({
					height: 250,
					body: getBandwidthPanel(),
				}),
				new SceneFlexItem({
					height: 250,
					body: getConnectionsPanel(),
				}),
				new SceneFlexLayout({
					direction: "row",
					height: 250,
					children: [
						new SceneFlexItem({
							width: "50%",
							body: getCPUPanel(),
						}),
						new SceneFlexItem({
							width: "50%",
							body: getMemoryPanel(),
						}),
					],
				}),
				new SceneFlexLayout({
					direction: "row",
					height: 250,
					children: [
						new SceneFlexItem({
							width: "50%",
							body: getLoadAveragePanel(),
						}),
						new SceneFlexItem({
							width: "50%",
							body: getReadWriteTimePanel(),
						}),
					],
				}),
				new SceneFlexLayout({
					direction: "row",
					height: 250,
					children: [
						new SceneFlexItem({
							width: "50%",
							body: getWrapCountPanel(),
						}),
						new SceneFlexItem({
							width: "50%",
							body: getNetstatPanel(),
						}),
					],
				}),
			],
		}),
		controls: [
			new VariableValueSelectors({}),
			new SceneControlsSpacer(),
			new SceneTimePicker({ isOnCanvas: true }),
			new SceneRefreshPicker({
				intervals: ["5s", "1m", "1h"],
				isOnCanvas: true,
			}),
		],
	});
}
