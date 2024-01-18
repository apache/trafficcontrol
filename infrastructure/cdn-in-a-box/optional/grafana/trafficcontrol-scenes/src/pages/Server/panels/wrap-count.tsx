import { PanelBuilders, SceneQueryRunner } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getWrapCountPanel = () => {
  const defaultQuery = {
    refId: 'A',
    query: `SELECT mean("vol1_wrap_count") AS "vol1", mean("vol2_wrap_count") AS "vol2" FROM "monthly"."wrap_count.1min" WHERE hostname='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$col',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.CACHE_STATS,
    queries: [defaultQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Wrap Count')
    .setData(qr)
    .setCustomFieldConfig('spanNulls', true)
    .setCustomFieldConfig('fillOpacity', 20)
    .build();
};
