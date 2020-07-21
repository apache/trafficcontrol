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

	"github.com/apache/trafficcontrol/grove/cache/memlrucache"
	"github.com/apache/trafficcontrol/grove/icache"
	"github.com/apache/trafficcontrol/lib/go-log"
)

const MemLRUCacheName = "cache_mem_lru"

func init() {
	AddPlugin(10000, Funcs{load: memlruCacheLoadCfg, loadCache: memlruCacheLoad})
}

type MemLRUCacheCfg struct {
	// SizeBytes is the size of the cache in bytes. Required.
	SizeBytes *uint64 `json:"size_bytes"`
}

func memlruCacheLoadCfg(b json.RawMessage) interface{} {
	cfg := MemLRUCacheCfg{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln(MemLRUCacheName + " cache plugin loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf(MemLRUCacheName+" cache plugin load success: %+v\n", cfg)
	return &cfg
}

func memlruCacheLoad(icfg interface{}, d LoadCacheData) map[string]icache.Cache {
	log.Errorln("DEBUG memlruCacheLoad started")
	// TODO: remove old deprecated config.Config values, after the next major release
	numBytes := uint64(d.Config.CacheSizeBytes)
	if numBytes == 0 {
		log.Errorln("DEBUG memlruCacheLoad no bytes, returning!")
		return nil
	}
	if icfg != nil {
		cfg, ok := icfg.(*MemLRUCacheCfg)
		if !ok {
			// should never happen
			log.Errorf(MemLRUCacheName+" cache config '%v' type '%T' expected *MemLRUCacheCfg\n", icfg, icfg)
			return nil
		}
		if cfg.SizeBytes == nil {
			log.Errorln(MemLRUCacheName + " cache config missing required 'size_bytes', cannot enable")
			return nil
		}
		numBytes = *cfg.SizeBytes
	}
	log.Errorln("DEBUG memlruCacheLoad success returning!")
	return map[string]icache.Cache{"": memlrucache.New(numBytes)}
}
