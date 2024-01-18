import { PanelBuilders, SceneQueryRunner } from '@grafana/scenes';
import { INFLUXDB_DATASOURCES_REF } from '../../../constants';

export const getNetstatPanel = () => {
  const defaultQuery = {
    refId: 'A',
    query: `SELECT mean("tcp_close") AS "tcp_close", mean("tcp_close_wait") AS "tcp_close_wait", mean("tcp_established") AS "tcp_established", mean("tcp_time_wait") AS "tcp_time_wait", mean("tcp_closing") AS "tcp_closing", mean("tcp_fin_wait1") AS "tcp_fin_wait1", mean("tcp_fin_wait2") AS "tcp_fin_wait2", mean("tcp_last_ack") AS "tcp_last_ack", mean("tcp_syn_recv") AS "tcp_syn_recv", mean("tcp_syn_sent") AS "tcp_syn_sent" FROM "netstat" WHERE host='$hostname' AND $timeFilter GROUP BY time($interval) fill(null)`,
    rawQuery: true,
    resultFormat: 'time_series',
    alias: '$col',
  };

  const qr = new SceneQueryRunner({
    datasource: INFLUXDB_DATASOURCES_REF.TELEGRAF,
    queries: [defaultQuery],
  });

  return PanelBuilders.timeseries()
    .setTitle('Netstat')
    .setData(qr)
    .setCustomFieldConfig('spanNulls', true)
    .setCustomFieldConfig('fillOpacity', 20)
    .build();
};
