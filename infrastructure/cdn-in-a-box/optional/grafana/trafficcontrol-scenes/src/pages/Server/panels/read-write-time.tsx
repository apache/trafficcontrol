import { PanelBuilders, SceneQueryRunner } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getReadWriteTimePanel = () => {
  const defaultQueries = [
    {
      refId: 'A',
      query: `SELECT non_negative_derivative(sum("read_time"), 10s) AS "read_time" FROM "diskio" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
      rawQuery: true,
      resultFormat: 'time_series',
      alias: '$col',
    },
    {
      refId: 'B',
      query: `SELECT non_negative_derivative(sum("write_time"), 10s) AS "write_time" FROM "diskio" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
      rawQuery: true,
      resultFormat: 'time_series',
      alias: '$col',
    },
  ];

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.TELEGRAF,
    queries: [...defaultQueries],
  });

  return PanelBuilders.timeseries()
    .setTitle('Read/Write Time')
    .setData(qr)
    .setCustomFieldConfig('spanNulls', true)
    .setCustomFieldConfig('fillOpacity', 20)
    .setUnit('ns')
    .build();
};
