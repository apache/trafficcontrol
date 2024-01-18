import { PanelBuilders, SceneQueryRunner } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getMemoryPanel = () => {
  const defaultQuery = {
    refId: 'A',
    query: `SELECT mean("used_percent") AS "mem_used" FROM "mem" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$col',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.TELEGRAF,
    queries: [defaultQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Memory Usage')
    .setData(qr)
    .setCustomFieldConfig('spanNulls', true)
    .setCustomFieldConfig('fillOpacity', 20)
    .setUnit('%')
    .build();
};
