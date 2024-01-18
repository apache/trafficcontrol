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
} from '@grafana/scenes';
import { getBandwidthPanel } from './panels/bandwidth';
import { getConnectionsPanel } from './panels/connections';
import { INFLUXDB_DATASOURCES_REF } from '../../constants';

export function getCacheGroupScene() {
  const timeRange = new SceneTimeRange({
    from: 'now-6h',
    to: 'now',
  });

  const cachegroup = new QueryVariable({
    name: 'cachegroup',
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    query: 'SHOW TAG VALUES ON "cache_stats" FROM "monthly"."bandwidth" with key = "cachegroup"',
  });

  return new EmbeddedScene({
    $timeRange: timeRange,
    $variables: new SceneVariableSet({
      variables: [cachegroup],
    }),
    body: new SceneFlexLayout({
      direction: 'column',
      children: [
        new SceneFlexItem({
          minHeight: 300,
          body: getBandwidthPanel(),
        }),
        new SceneFlexItem({
          minHeight: 300,
          body: getConnectionsPanel(),
        }),
      ],
    }),
    controls: [
      new VariableValueSelectors({}),
      new SceneControlsSpacer(),
      new SceneTimePicker({ isOnCanvas: true }),
      new SceneRefreshPicker({
        intervals: ['5s', '1m', '1h'],
        isOnCanvas: true,
      }),
    ],
  });
}
