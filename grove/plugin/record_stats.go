package plugin

func init() {
	AddPlugin(10000, Funcs{afterRespond: recordStats})
}

func recordStats(icfg interface{}, d AfterRespondData) {
	d.Stats.Write(d.W, d.Conn, d.Req.Host, d.Req.RemoteAddr, d.RespCode, d.BytesWritten, d.CacheHit)
}
