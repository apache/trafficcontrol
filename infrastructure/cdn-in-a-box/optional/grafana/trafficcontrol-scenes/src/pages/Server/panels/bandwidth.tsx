import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getBandwidthPanel = () => {
  const defaultBandwidthQuery = {
    refId: 'A',
    query: `SELECT mean(value) FROM "monthly"."bandwidth.1min" WHERE hostname='$hostname' AND $timeFilter GROUP BY time(60s)`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: 'bandwidth',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    queries: [defaultBandwidthQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Bandwidth')
    .setData(qr)
    .setCustomFieldConfig('fillOpacity', 20)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setUnit('kbps')
    .build();
};
