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

import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getConnectionsPanel = () => {
  const connectionQuery = {
    refId: 'A',
    query:
      'SELECT mean("value") FROM "monthly"."connections.1min" WHERE ("cachegroup" = \'$cachegroup\') AND $timeFilter GROUP BY time($interval), "hostname" fill(null)',
    rawQuery: true,
    resultFormat: 'time_series',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    queries: [connectionQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Connections (stacked)')
    .setCustomFieldConfig('fillOpacity', 20)
    .setData(qr)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setCustomFieldConfig('spanNulls', true)
    .build();
};
