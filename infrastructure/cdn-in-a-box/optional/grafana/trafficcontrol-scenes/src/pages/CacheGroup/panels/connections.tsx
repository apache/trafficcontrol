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
