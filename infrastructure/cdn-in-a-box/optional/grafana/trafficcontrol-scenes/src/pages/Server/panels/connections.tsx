import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getConnectionsPanel = () => {
  const defaultConnectionsQuery = [
    {
      refId: 'A',
      query: `SELECT mean(value) FROM "monthly"."connections.1min" WHERE hostname='$hostname' AND $timeFilter GROUP BY time(60s)`,
      rawQuery: true,
      resultFormat: 'time_series',
      alias: 'connections',
    },
  ];

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    queries: [...defaultConnectionsQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Connections')
    .setData(qr)
    .setCustomFieldConfig('fillOpacity', 20)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setCustomFieldConfig('spanNulls', true)
    .build();
};
