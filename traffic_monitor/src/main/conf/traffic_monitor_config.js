{
	"traffic_monitor_config": {
		"health.polling.interval": "5000",
		"tm.polling.interval": "10000",
		"tm.hostname": "",
		"tm.healthParams.polling.url": "https://${tmHostname}/health/${cdnName}",
		"hack.ttl": "30",
		"cdnName": "",
		"peers.polling.url": "http://${hostname}/publish/CrStates?raw",
		"health.timepad": "20",
		"health.event-count": "200",
		"tm.dataServer.polling.url": "https://${tmHostname}/dataserver/orderby/id",
		"tm.auth.url": "https://${tmHostname}/login",
		"tm.auth.username": "",
		"tm.auth.password": "",
		"tm.crConfig.json.polling.url": "https://${tmHostname}/CRConfig-Snapshots/${cdnName}/CRConfig.json"
	}
}
