package plugin

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"encoding/json"

	"github.com/apache/trafficcontrol/v8/grove/web"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
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
