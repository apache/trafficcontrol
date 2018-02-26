package beforeparentrequest

import (
	"encoding/json"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

const ModOrgReqHdrName = "mod_org_request_headers"

func init() {
	AddPlugin(10000, ModOrgReqHdrName, modOrgReqHdr, modOrgReqHdrLoad)
}

type Hdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ModHdrs struct {
	Set  []Hdr    `json:"set"`
	Drop []string `json:"drop"`
}

// Any returns whether any header modifications exist
func (mh *ModHdrs) Any() bool {
	return len(mh.Set) > 0 || len(mh.Drop) > 0
}

// Mod drops and sets the headers in h according to its rules.
func (mh *ModHdrs) Mod(h http.Header) {
	if h == nil || len(h) == 0 { // this happens on a dial tcp timeout
		log.Debugf(ModOrgReqHdrName + ": Header is  a nil map")
		return
	}
	for _, hdr := range mh.Drop {
		log.Debugf(ModOrgReqHdrName+": Dropping header %s\n", hdr)
		h.Del(hdr)
	}
	for _, hdr := range mh.Set {
		log.Debugf(ModOrgReqHdrName+": Setting header %s: %s \n", hdr.Name, hdr.Value)
		h.Set(hdr.Name, hdr.Value)
	}
}

func modOrgReqHdrLoad(b json.RawMessage) interface{} {
	cfg := ModHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln(ModOrgReqHdrName + " loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf(ModOrgReqHdrName+" load success: %+v\n", cfg)
	return &cfg
}

func modOrgReqHdr(icfg interface{}, d Data) {
	log.Debugf(ModOrgReqHdrName + " calling\n")
	if icfg == nil {
		log.Debugln(ModOrgReqHdrName + " has no config, returning.")
		return
	}
	cfg, ok := icfg.(*ModHdrs)
	if !ok {
		// should never happen
		log.Errorf(ModOrgReqHdrName+" config '%v' type '%T' expected *ModHdrs\n", icfg, icfg)
		return
	}

	log.Debugf(ModOrgReqHdrName+" config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	cfg.Mod(d.Req.Header)
}
