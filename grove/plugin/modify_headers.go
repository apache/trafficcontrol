package plugin

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{load: modRespHdrLoad, beforeRespond: modRespHdr})
}

func modRespHdrLoad(b json.RawMessage) interface{} {
	cfg := web.ModHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln("modifyheaders loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf("modifyheaders load success: %+v\n", cfg)
	return &cfg
}

func modRespHdr(icfg interface{}, d BeforeRespondData) {
	log.Debugf("modifyheaders calling\n")
	if icfg == nil {
		log.Debugln("modifyheaders has no config, returning.")
		return
	}
	cfg, ok := icfg.(*web.ModHdrs)
	if !ok {
		// should never happen
		log.Errorf("modifyheaders config '%v' type '%T' expected *web.ModHdrs\n", icfg, icfg)
		return
	}

	log.Debugf("modifyheaders config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	*d.Hdr = web.CopyHeader(*d.Hdr)
	cfg.Mod(*d.Hdr)
}
