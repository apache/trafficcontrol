import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getBandwidthPanel = () => {
  const defaultBandwidthQuery = {
    refId: 'A',
    query: `SELECT mean(value) FROM "monthly"."kbps.ds.1min" WHERE deliveryservice='$deliveryservice' AND cachegroup = 'total'  and $timeFilter GROUP BY time(60s), deliveryservice ORDER BY asc`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$tag_deliveryservice',
    measurement: 'bw',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.DELIVERYSERVICE_STATS,
    queries: [defaultBandwidthQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Bandwidth')
    .setData(qr)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setCustomFieldConfig('axisCenteredZero', true)
    .setUnit('bps')
    .build();
};
