package plugin

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{load: modOrgReqHdrLoad, beforeParentRequest: modOrgReqHdr})
}

func modOrgReqHdrLoad(b json.RawMessage) interface{} {
	cfg := web.ModHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln("modify_parent_request_headers loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf("modify_parent_request_headers load success: %+v\n", cfg)
	return &cfg
}

func modOrgReqHdr(icfg interface{}, d BeforeParentRequestData) {
	log.Debugf("modify_parent_request_headers calling\n")
	if icfg == nil {
		log.Debugln("modify_parent_request_headers has no config, returning.")
		return
	}
	cfg, ok := icfg.(*web.ModHdrs)
	if !ok {
		// should never happen
		log.Errorf("modify_parent_request_headers config '%v' type '%T' expected *ModHdrs\n", icfg, icfg)
		return
	}

	log.Debugf("modify_parent_request_headers config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	cfg.Mod(d.Req.Header)
}
