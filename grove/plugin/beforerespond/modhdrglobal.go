package beforerespond

import (
	"encoding/json"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

const ModRespGlblHdrName = "mod_response_headers_global"

func init() {
	AddPlugin(10000, ModRespGlblHdrName, modRespGlblHdr, modRespGlblHdrLoad)
}

type GlblHdr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ModGlblHdrs struct {
	Set  []GlblHdr `json:"set"`
	Drop []string  `json:"drop"`
}

// Any returns whether any header modifications exist
func (mh *ModGlblHdrs) Any() bool {
	return len(mh.Set) > 0 || len(mh.Drop) > 0
}

// Mod drops and sets the headers in h according to its rules.
func (mh *ModGlblHdrs) Mod(h http.Header) {
	if h == nil || len(h) == 0 { // this happens on a dial tcp timeout
		log.Debugf(ModRespGlblHdrName + ": Header is  a nil map")
		return
	}
	for _, hdr := range mh.Drop {
		log.Debugf(ModRespGlblHdrName+": Dropping header %s\n", hdr)
		h.Del(hdr)
	}
	for _, hdr := range mh.Set {
		log.Debugf(ModRespGlblHdrName+": Setting header %s: %s \n", hdr.Name, hdr.Value)
		h.Set(hdr.Name, hdr.Value)
	}
}

func modRespGlblHdrLoad(b json.RawMessage) interface{} {
	cfg := ModGlblHdrs{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln(ModRespGlblHdrName + " loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf(ModRespGlblHdrName+" load success: %+v\n", cfg)
	return &cfg
}

func modRespGlblHdr(icfg interface{}, d Data) {
	log.Debugf(ModRespGlblHdrName + " calling\n")
	if icfg == nil {
		log.Debugln(ModRespGlblHdrName + " has no config, returning.")
		return
	}
	cfg, ok := icfg.(*ModGlblHdrs)
	if !ok {
		// should never happen
		log.Errorf(ModRespGlblHdrName+" config '%v' type '%T' expected *ModGlblHdrs\n", icfg, icfg)
		return
	}

	log.Debugf(ModRespGlblHdrName+" config len(set) %+v len(drop) %+v\n", cfg.Set, cfg.Drop)
	if !cfg.Any() {
		return
	}
	*d.Hdr = web.CopyHeader(*d.Hdr)
	cfg.Mod(*d.Hdr)
}
