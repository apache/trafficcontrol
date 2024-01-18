import { PanelBuilders, SceneQueryRunner } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getLoadAveragePanel = () => {
  const defaultQuery = {
    refId: 'A',
    query: `SELECT mean("load1") AS "load1", mean("load5") AS "load5", mean("load15") AS "load15" FROM "system" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$col',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.TELEGRAF,
    queries: [defaultQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Load Average')
    .setData(qr)
    .setCustomFieldConfig('fillOpacity', 20)
    .setCustomFieldConfig('spanNulls', true)
    .build();
};
