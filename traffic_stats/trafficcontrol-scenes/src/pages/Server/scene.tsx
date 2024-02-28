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
	EmbeddedScene,
	QueryVariable,
	SceneControlsSpacer,
	SceneFlexItem,
	SceneFlexLayout,
	SceneRefreshPicker,
	SceneTimePicker,
	SceneTimeRange,
	SceneVariableSet,
	VariableValueSelectors,
} from "@grafana/scenes";
import { INFLUXDB_DATASOURCES_REF } from "const";

import {
	getBandwidthPanel,
	getConnectionsPanel,
	getCPUPanel,
	getLoadAveragePanel,
	getMemoryPanel,
	getNetstatPanel,
	getReadWriteTimePanel,
	getWrapCountPanel,
} from "./panels";

/**
 * Returns an EmbeddedScene with a specific time range and variables, consisting of multiple
 * SceneFlexLayout and SceneFlexItem components for displaying various panels and controls.
 *
 * @returns The EmbeddedScene with the specified time range, variables, body, and controls.
 */
export function getServerScene(): EmbeddedScene {
	const timeRange = new SceneTimeRange({
		from: "now-6h",
		to: "now",
	});

	const hostname = new QueryVariable({
		datasource: INFLUXDB_DATASOURCES_REF.cacheStats,
		name: "hostname",
		query: 'SHOW TAG VALUES ON "cache_stats" FROM "monthly"."bandwidth" with key = "hostname"',
	});

	return new EmbeddedScene({
		$timeRange: timeRange,
		$variables: new SceneVariableSet({
			variables: [hostname],
		}),
		body: new SceneFlexLayout({
			children: [
				new SceneFlexItem({
					body: getBandwidthPanel(),
					height: 250,
				}),
				new SceneFlexItem({
					body: getConnectionsPanel(),
					height: 250,
				}),
				new SceneFlexLayout({
					children: [
						new SceneFlexItem({
							body: getCPUPanel(),
							width: "50%",
						}),
						new SceneFlexItem({
							body: getMemoryPanel(),
							width: "50%",
						}),
					],
					direction: "row",
					height: 250,
				}),
				new SceneFlexLayout({
					children: [
						new SceneFlexItem({
							body: getLoadAveragePanel(),
							width: "50%",
						}),
						new SceneFlexItem({
							body: getReadWriteTimePanel(),
							width: "50%",
						}),
					],
					direction: "row",
					height: 250,
				}),
				new SceneFlexLayout({
					children: [
						new SceneFlexItem({
							body: getWrapCountPanel(),
							width: "50%",
						}),
						new SceneFlexItem({
							body: getNetstatPanel(),
							width: "50%",
						}),
					],
					direction: "row",
					height: 250,
				}),
			],
			direction: "column",
		}),
		controls: [
			new VariableValueSelectors({}),
			new SceneControlsSpacer(),
			new SceneTimePicker({isOnCanvas: true}),
			new SceneRefreshPicker({
				intervals: ["5s", "1m", "1h"],
				isOnCanvas: true,
			}),
		],
	});
}
