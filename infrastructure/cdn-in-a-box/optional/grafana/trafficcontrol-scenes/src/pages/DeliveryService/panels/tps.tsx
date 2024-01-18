import { SceneQueryRunner, PanelBuilders } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getTpsPanel = () => {
  const tpsQueries = [
    {
      refId: 'A',
      query: `SELECT mean(value) FROM "monthly"."tps_2xx.ds.1min" WHERE $timeFilter AND deliveryservice='$deliveryservice' GROUP BY time(60s) ORDER BY asc`,
      rawQuery: true,
      resultFormat: 'time_series',
      measurement: 'tps_2xx',
      hide: false,
      tags: {
        deliveryservice: `$deliveryservice`,
      },
    },
    {
      refId: 'B',
      query: `SELECT mean(value) FROM \"monthly\".\"tps_3xx.ds.1min\" WHERE $timeFilter AND deliveryservice='$deliveryservice' GROUP BY time(60s) ORDER BY asc`,
      rawQuery: true,
      resultFormat: 'time_series',
    },
    {
      refId: 'C',
      query: `SELECT mean(value) FROM \"monthly\".\"tps_4xx.ds.1min\" WHERE $timeFilter AND deliveryservice='$deliveryservice' GROUP BY time(60s) ORDER BY asc`,
      rawQuery: true,
      resultFormat: 'time_series',
    },
    {
      refId: 'D',
      query: `SELECT mean(value) FROM \"monthly\".\"tps_5xx.ds.1min\" WHERE $timeFilter AND deliveryservice='$deliveryservice' GROUP BY time(60s) ORDER BY asc`,
      rawQuery: true,
      resultFormat: 'time_series',
    },
  ];

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.DELIVERYSERVICE_STATS,
    queries: [...tpsQueries],
  });

  return PanelBuilders.timeseries()
    .setTitle('TPS')
    .setData(qr)
    .setOption('legend', { showLegend: true, calcs: ['max'] })
    .setCustomFieldConfig('axisCenteredZero', true)
    .setCustomFieldConfig('spanNulls', true)
    .build();
};
