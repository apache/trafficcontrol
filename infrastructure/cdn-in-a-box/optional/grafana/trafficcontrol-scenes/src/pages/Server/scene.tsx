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
import {
  getBandwidthPanel,
  getConnectionsPanel,
  getCPUPanel,
  getMemoryPanel,
  getLoadAveragePanel,
  getReadWriteTimePanel,
  getWrapCountPanel,
  getNetstatPanel,
} from './panels';
import { INFLUXDB_DATASOURCES_REF } from '../../constants';

export function getServerScene() {
  const timeRange = new SceneTimeRange({
    from: 'now-6h',
    to: 'now',
  });

  const hostname = new QueryVariable({
    name: 'hostname',
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    query: 'SHOW TAG VALUES ON "cache_stats" FROM "monthly"."bandwidth" with key = "hostname"',
  });

  return new EmbeddedScene({
    $timeRange: timeRange,
    $variables: new SceneVariableSet({
      variables: [hostname],
    }),
    body: new SceneFlexLayout({
      direction: 'column',
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
          direction: 'row',
          height: 250,
          children: [
            new SceneFlexItem({
              width: '50%',
              body: getCPUPanel(),
            }),
            new SceneFlexItem({
              width: '50%',
              body: getMemoryPanel(),
            }),
          ],
        }),
        new SceneFlexLayout({
          direction: 'row',
          height: 250,
          children: [
            new SceneFlexItem({
              width: '50%',
              body: getLoadAveragePanel(),
            }),
            new SceneFlexItem({
              width: '50%',
              body: getReadWriteTimePanel(),
            }),
          ],
        }),
        new SceneFlexLayout({
          direction: 'row',
          height: 250,
          children: [
            new SceneFlexItem({
              width: '50%',
              body: getWrapCountPanel(),
            }),
            new SceneFlexItem({
              width: '50%',
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
        intervals: ['5s', '1m', '1h'],
        isOnCanvas: true,
      }),
    ],
  });
}
