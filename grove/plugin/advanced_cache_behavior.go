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
	"fmt"
	"strings"

	log "github.com/apache/trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{load: loadCacheBehaviorConfig, beforeCacheLookUp: customCacheLookUp})
}

type AdvancedCacheBehaviorConfig struct {
	Whitelist []string `json:"whitelist"`
}

func customCacheLookUp(icfg interface{}, d BeforeCacheLookUpData) {
	acbConfig := icfg.(*AdvancedCacheBehaviorConfig)
	headers := make([]string, 0, 16)
	for _, headerName := range acbConfig.Whitelist {
		if headerValue := d.Req.Header.Get(headerName); headerValue != "" {
			headers = append(headers, fmt.Sprintf("%v=%v", headerName, headerValue))
		}
	}
	if len(headers) > 0 {
		d.CacheKeyOverrideFunc("headers=" + strings.Join(headers, ",") + ":" + d.DefaultCacheKey)
	}
}

func loadCacheBehaviorConfig(b json.RawMessage) interface{} {
	acbConfig := AdvancedCacheBehaviorConfig{}
	err := json.Unmarshal(b, &acbConfig)
	if err != nil {
		log.Errorln("Advanced Cache Behavior config loading error: " + err.Error())
		return nil
	}
	log.Debugf("Advanced Cache Behavior config loaded")
	return &acbConfig
}
