{"traffic_monitor_config": {
  "tm.healthParams.polling.url": "https://${tmHostname}/health/${cdnName}",
  "hack.ttl": "30",
  "tm.auth.url": "https://${tmHostname}/login",
  "tm.auth.username": "admin",
  "tm.auth.password": 'password',
  "health.polling.interval": "5000",
  "peers.polling.url": "http://${hostname}/publish/CrStates?raw",
  "cdnName": "cdnname",
  "health.event-count": "200",
  "health.timepad": "20",
  "tm.dataServer.polling.url": "https://${tmHostname}/dataserver/orderby/id",
  "tm.crConfig.json.polling.url": "https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json",
  "tm.polling.interval": "10000",
  "tm.hostname": "traffic-ops.company.net"
}}
