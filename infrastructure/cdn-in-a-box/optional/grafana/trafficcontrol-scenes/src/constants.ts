import pluginJson from './plugin.json';

export const PLUGIN_BASE_URL = `/a/${pluginJson.id}`;

export enum ROUTES {
  CacheGroup = 'cache-group',
  DeliveryService = 'delivery-service',
  Server = 'server',
}

export const PROMETHEUS_DATASOURCE_REF = {
  uid: 'prometheus',
  type: 'prometheus',
};

export const INFLUXDB_DATASOURCES_REF = {
  CACHE_STATS: {
    uid: 'cache_stats',
    type: 'influxdb',
  },
  DELIVERYSERVICE_STATS: {
    uid: 'deliveryservice_stats',
    type: 'influxdb',
  },
  DAILY_STATS: {
    uid: 'daily_stats',
    type: 'influxdb',
  },
  TELEGRAF: {
    uid: 'telegraf',
    type: 'influxdb',
  },
};
