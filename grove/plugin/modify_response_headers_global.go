package plugin

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{load: modRespGlblHdrLoad, beforeRespond: modRespGlblHdr})
}

func modRespGlblHdrLoad(b json.RawMessage) interface{} {
	cfg := web.ModHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln("mod_response_headers_global  loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf("mod_response_headers_global: load success: %+v\n", cfg)
	return &cfg
}

func modRespGlblHdr(icfg interface{}, d BeforeRespondData) {
	log.Debugf("mod_response_headers_global calling\n")
	if icfg == nil {
		log.Debugln("mod_response_headers_global has no config, returning.")
		return
	}
	cfg, ok := icfg.(*web.ModHdrs)
	if !ok {
		// should never happen
		log.Errorf("mod_response_headers_global config '%v' type '%T' expected *ModGlblHdrs\n", icfg, icfg)
		return
	}

	log.Debugf("mod_response_headers_global config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	*d.Hdr = web.CopyHeader(*d.Hdr)
	cfg.Mod(*d.Hdr)
}
