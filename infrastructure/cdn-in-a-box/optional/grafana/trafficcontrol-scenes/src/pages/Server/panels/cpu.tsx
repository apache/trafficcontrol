import { PanelBuilders, SceneQueryRunner } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getCPUPanel = () => {
  const defaultQuery = {
    refId: 'A',
    query: `SELECT mean("usage_system") AS "cpu_system", mean("usage_iowait") AS "cpu_iowait", mean("usage_user") AS "cpu_user", mean("usage_guest") AS "cpu_guest", mean("usage_steal") AS "cpu_steal" FROM "cpu" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$col',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.TELEGRAF,
    queries: [defaultQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('CPU Usage')
    .setData(qr)
    .setUnit('%')
    .setCustomFieldConfig('spanNulls', true)
    .setCustomFieldConfig('fillOpacity', 20)
    .build();
};
