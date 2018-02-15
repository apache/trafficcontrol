package afterrespond

const RecordStatsName = "record_stats"

func init() {
	AddPlugin(10000, RecordStatsName, recordStats, nil)
}

func recordStats(icfg interface{}, d Data) {
	d.Stats.Write(d.W, d.Conn, d.Req.Host, d.Req.RemoteAddr, d.RespCode, d.BytesWritten, d.CacheHit)
}
