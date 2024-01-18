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
import { getTpsPanel } from './panels/tps';
import { getBandwidthByCGPanel } from './panels/bandwidth-cg';
import { INFLUXDB_DATASOURCES_REF } from '../../constants';

export function getDeliveryServiceScene() {
  const timeRange = new SceneTimeRange({
    from: 'now-6h',
    to: 'now',
  });

  const deliveryService = new QueryVariable({
    name: 'deliveryservice',
    datasource: INFLUXDB_DATASOURCES_REF.DELIVERYSERVICE_STATS,
    query: 'SHOW TAG VALUES ON "deliveryservice_stats" FROM "monthly"."kbps" with key = "deliveryservice"',
  });

  return new EmbeddedScene({
    $timeRange: timeRange,
    $variables: new SceneVariableSet({
      variables: [deliveryService],
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
          body: getTpsPanel(),
        }),
        new SceneFlexItem({
          minHeight: 300,
          body: getBandwidthByCGPanel(),
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
