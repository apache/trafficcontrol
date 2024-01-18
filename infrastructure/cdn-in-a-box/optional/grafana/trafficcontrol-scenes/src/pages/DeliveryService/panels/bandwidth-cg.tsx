import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getBandwidthByCGPanel = () => {
  const bandwidthByCacheGroupQuery = {
    refId: 'A',
    query: `SELECT mean(value) FROM "monthly"."kbps.cg.1min" WHERE deliveryservice='$deliveryservice' AND cachegroup != 'all' and $timeFilter GROUP BY time(60s), cachegroup`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$tag_cachegroup',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.DELIVERYSERVICE_STATS,
    queries: [bandwidthByCacheGroupQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Bandwidth by CacheGroup')
    .setData(qr)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setCustomFieldConfig('axisCenteredZero', true)
    .setCustomFieldConfig('spanNulls', true)
    .setUnit('bps')
    .build();
};
