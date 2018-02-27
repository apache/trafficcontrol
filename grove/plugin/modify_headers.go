package plugin

import (
	"encoding/json"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{load: modRespHdrLoad, beforeRespond: modRespHdr})
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
	if len(h) == 0 { // this happens on a dial tcp timeout
		log.Debugf("modifyheaders: Header is  a nil map")
		return
	}
	for _, hdr := range mh.Drop {
		log.Debugf("modifyheaders: Dropping header %s\n", hdr)
		h.Del(hdr)
	}
	for _, hdr := range mh.Set {
		log.Debugf("modifyheaders: Setting header %s: %s \n", hdr.Name, hdr.Value)
		h.Set(hdr.Name, hdr.Value)
	}
}

func modRespHdrLoad(b json.RawMessage) interface{} {
	cfg := ModHdrs{}
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
	cfg, ok := icfg.(*ModHdrs)
	if !ok {
		// should never happen
		log.Errorf("modifyheaders config '%v' type '%T' expected *ModHdrs\n", icfg, icfg)
		return
	}

	log.Debugf("modifyheaders config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	*d.Hdr = web.CopyHeader(*d.Hdr)
	cfg.Mod(*d.Hdr)
}
