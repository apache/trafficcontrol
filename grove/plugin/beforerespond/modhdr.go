package beforerespond

import (
	"encoding/json"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

const ModRespHdrName = "mod_response_headers"

func init() {
	AddPlugin(10000, ModRespHdrName, modRespHdr, modRespHdrLoad)
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
		log.Debugf(ModRespHdrName + ": Header is  a nil map")
		return
	}
	for _, hdr := range mh.Drop {
		log.Debugf(ModRespHdrName+": Dropping header %s\n", hdr)
		h.Del(hdr)
	}
	for _, hdr := range mh.Set {
		log.Debugf(ModRespHdrName+": Setting header %s: %s \n", hdr.Name, hdr.Value)
		h.Set(hdr.Name, hdr.Value)
	}
}

func modRespHdrLoad(b json.RawMessage) interface{} {
	cfg := ModHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln(ModRespHdrName + " loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf(ModRespHdrName+" load success: %+v\n", cfg)
	return &cfg
}

func modRespHdr(icfg interface{}, d Data) {
	log.Debugf(ModRespHdrName + " calling\n")
	if icfg == nil {
		log.Debugln(ModRespHdrName + " has no config, returning.")
		return
	}
	cfg, ok := icfg.(*ModHdrs)
	if !ok {
		// should never happen
		log.Errorf(ModRespHdrName+" config '%v' type '%T' expected *ModHdrs\n", icfg, icfg)
		return
	}

	log.Debugf(ModRespHdrName+" config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	*d.Hdr = web.CopyHeader(*d.Hdr)
	cfg.Mod(*d.Hdr)
}
