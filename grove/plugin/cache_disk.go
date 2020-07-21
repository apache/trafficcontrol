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

	"github.com/apache/trafficcontrol/grove/cache/diskcache"
	"github.com/apache/trafficcontrol/grove/cache/memlrucache"
	"github.com/apache/trafficcontrol/grove/cache/tiercache"
	"github.com/apache/trafficcontrol/grove/config"
	"github.com/apache/trafficcontrol/grove/icache"
	"github.com/apache/trafficcontrol/lib/go-log"
)

const DiskCacheName = "cache_disk"

func init() {
	AddPlugin(10000, Funcs{load: diskCacheLoadCfg, loadCache: diskCacheLoad})
}

type DiskCacheCfg struct {
	// MemSizeBytes is the size of the memory cache in front of the disk cache,  in bytes. Required.
	MemBytes uint64 `json:"mem_bytes"`

	// Files is the name of each disk cache, and the files to use for it.
	Files map[string][]config.CacheFile `json:"files"`
}

func diskCacheLoadCfg(b json.RawMessage) interface{} {
	cfg := DiskCacheCfg{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Errorln(DiskCacheName + " cache plugin loading config, unmarshalling JSON: " + err.Error())
		return nil
	}
	log.Debugf(DiskCacheName+" cache plugin load success: %+v\n", cfg)
	return &cfg
}

func diskCacheLoad(icfg interface{}, d LoadCacheData) map[string]icache.Cache {
	log.Errorln("disk cache starting")
	// TODO remove deprecated config, after the next major version.
	memBytes := d.Config.FileMemBytes
	files := d.Config.CacheFiles // map[string][]CacheFile `json:"cache_files"`
	if len(files) == 0 {
		if icfg == nil {
			log.Errorln("disk cache no config, returning")
			return nil
		}
		cfg, ok := icfg.(*DiskCacheCfg)
		if !ok {
			// should never happen
			log.Errorf(DiskCacheName+" cache config '%v' type '%T' expected *DiskCacheCfg\n", icfg, icfg)
			log.Errorln("disk cache bad config type, returning")
			return nil
		}

		memBytes = int(cfg.MemBytes)
		files = cfg.Files
	}

	if memBytes == 0 {
		log.Infoln(DiskCacheName + " disk cache config mem_bytes is 0, disabling.")
		log.Errorln("disk cache no mem bytes, returning")
		return nil
	}
	if len(files) == 0 {
		log.Infoln(DiskCacheName + " disk cache config has no files, disabling.")
		log.Errorln("disk cache no file, returning")
		return nil
	}

	caches := map[string]icache.Cache{}
	for name, files := range files {
		if name == "" {
			log.Errorln(DiskCacheName + " disk cache config has empty named cache, skipping.")
			log.Errorln("disk cache config empty name, returning")
			continue
		}
		if len(files) == 0 {
			log.Errorln("disk cache config no files, returning")
			log.Errorln(DiskCacheName + " disk cache name '" + name + "' has no files, skipping.")
			continue
		}

		multiDiskCache, err := diskcache.NewMulti(files)
		if err != nil {
			log.Errorln("disk cache err: " + err.Error())
			log.Errorln(DiskCacheName + " disk cache name '" + name + "' error, skipping: " + err.Error())
			continue
		}
		caches[name] = tiercache.New(memlrucache.New(uint64(memBytes)), multiDiskCache)
		log.Errorln("DEBUG diskCacheLoad adding '" + name + "'")
	}

	log.Errorln("DEBUG diskCacheLoad success returning!")
	return caches
}
