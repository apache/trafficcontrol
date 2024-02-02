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

import { getBandwidthPanel } from "./panels/bandwidth";
import { getBandwidthByCGPanel } from "./panels/bandwidth-cg";
import { getTpsPanel } from "./panels/tps";

/**
 * Returns an EmbeddedScene representing the delivery service scene.
 *
 * @returns EmbeddedScene representing the delivery service scene
 */
export function getDeliveryServiceScene(): EmbeddedScene {
	const timeRange = new SceneTimeRange({
		from: "now-6h",
		to: "now",
	});

	const deliveryService = new QueryVariable({
		datasource: INFLUXDB_DATASOURCES_REF.deliveryServiceStats,
		name: "deliveryservice",
		query: 'SHOW TAG VALUES ON "deliveryservice_stats" FROM "monthly"."kbps" with key = "deliveryservice"',
	});

	return new EmbeddedScene({
		$timeRange: timeRange,
		$variables: new SceneVariableSet({
			variables: [deliveryService],
		}),
		body: new SceneFlexLayout({
			children: [
				new SceneFlexItem({
					body: getBandwidthPanel(),
					minHeight: 300,
				}),
				new SceneFlexItem({
					body: getTpsPanel(),
					minHeight: 300,
				}),
				new SceneFlexItem({
					body: getBandwidthByCGPanel(),
					minHeight: 300,
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
