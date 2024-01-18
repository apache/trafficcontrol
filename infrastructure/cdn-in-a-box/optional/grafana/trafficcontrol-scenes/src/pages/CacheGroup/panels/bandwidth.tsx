import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getBandwidthPanel = () => {
  const cacheGroupBandwidthQuery = {
    refId: 'A',
    query:
      'SELECT sum(value) FROM "monthly"."bandwidth.1min" WHERE "cachegroup" = \'$cachegroup\' AND $timeFilter GROUP BY time(60s), cachegroup',
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$tag_cachegroup',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    queries: [cacheGroupBandwidthQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Total bandwidth (stacked)')
    .setData(qr)
    .setCustomFieldConfig('fillOpacity', 20)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setUnit('Kbits')
    .build();
};
