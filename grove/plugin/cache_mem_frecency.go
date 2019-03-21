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

	"github.com/apache/trafficcontrol/grove/cache/memfrecencycache"
	"github.com/apache/trafficcontrol/grove/icache"
	"github.com/apache/trafficcontrol/lib/go-log"
)

const MemFrecencyCacheName = "cache_mem_frecency"

func init() {
	AddPlugin(10000, Funcs{load: memFrecencyCacheLoadCfg, loadCache: memFrecencyCacheLoad})
}

type MemFrecencyCacheCfg struct {
	// SizeBytes is the size of the cache in bytes. Required.
	SizeBytes *uint64 `json:"size_bytes"`

	// HitWeight is the weight of hits against recency. Optional, defaults to memfrecencycache.DefaultHitWeight.
	HitWeight *float64 `json:"hit_weight"`
}

func memFrecencyCacheLoadCfg(b json.RawMessage) interface{} {
	cfg := MemFrecencyCacheCfg{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln(MemFrecencyCacheName + " cache plugin loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf(MemFrecencyCacheName+" cache plugin load success: %+v\n", cfg)
	return &cfg
}

func memFrecencyCacheLoad(icfg interface{}, d LoadCacheData) map[string]icache.Cache {
	log.Errorln("DEBUG memFrecencyCacheLoad started")
	if icfg == nil {
		log.Infoln(MemFrecencyCacheName + " cache has no config, not enabling.")
		return nil
	}

	cfg, ok := icfg.(*MemFrecencyCacheCfg)
	if !ok {
		// should never happen
		log.Errorf(MemFrecencyCacheName+" cache config '%v' type '%T' expected *MemFrecencyCacheCfg\n", icfg, icfg)
		return nil
	}

	if cfg.SizeBytes == nil {
		log.Errorln(MemFrecencyCacheName + " cache config missing required 'size_bytes', cannot enable")
		return nil
	}
	numBytes := *cfg.SizeBytes

	hitWeight := memfrecencycache.DefaultHitWeight
	if cfg.HitWeight != nil {
		hitWeight = *cfg.HitWeight
	}

	log.Errorln("DEBUG memFrecencyCacheLoad success returning!")
	return map[string]icache.Cache{"": memfrecencycache.New(numBytes, hitWeight)}
}
