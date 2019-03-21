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
	"net/http"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/grove/cachedata"
	"github.com/apache/trafficcontrol/grove/cacheobj"
	"github.com/apache/trafficcontrol/grove/config"
	"github.com/apache/trafficcontrol/grove/icache"
	"github.com/apache/trafficcontrol/grove/remapdata"
	"github.com/apache/trafficcontrol/grove/stat"
	"github.com/apache/trafficcontrol/grove/web"

	"github.com/apache/trafficcontrol/lib/go-log"
)

func AddPlugin(priority uint64, funcs Funcs) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error plugin.AddPlugin: runtime.Caller failed, can't get plugin names") // print, because this is called in init, loggers don't exist yet
		os.Exit(1)
	}
	pluginName := strings.TrimSuffix(path.Base(filename), ".go")
	plugins = append(plugins, pluginObj{funcs: funcs, priority: priority, name: pluginName})
}

type Funcs struct {
	load                LoadFunc
	startup             StartupFunc
	loadCache           LoadCacheFunc
	onRequest           OnRequestFunc
	beforeCacheLookUp   BeforeCacheLookupFunc
	beforeParentRequest BeforeParentRequestFunc
	beforeRespond       BeforeRespondFunc
	afterRespond        AfterRespondFunc
}

type StartupData struct {
	Config  config.Config
	Context *interface{}
	// Shared is the "plugins_shared" data for all rules. This is a `map[ruleName][key]value`. Keys and values are arbitrary data. This allows plugins to do pre-processing on the config, and store computed data in the context, to save processing during requests.
	Shared map[string]map[string]json.RawMessage
}

type LoadCacheData struct {
	Config  config.Config
	Context *interface{}
	Shared  map[string]json.RawMessage
}

type OnRequestData struct {
	W             http.ResponseWriter
	R             *http.Request
	InterfaceName string
	Stats         stat.Stats
	StatRules     remapdata.RemapRulesStats
	HTTPConns     *web.ConnMap
	HTTPSConns    *web.ConnMap
	RequestID     uint64
	Context       *interface{}
	cachedata.SrvrData
}

type BeforeParentRequestData struct {
	Req       *http.Request
	RemapRule string
	Context   *interface{}
}

// BeforeRespondData holds the data passed to plugins. The objects pointed to MAY NOT be modified, however, the location pointed to may be changed for the Code, Hdr, and Body. That iss, `*d.Hdr = myHdr` is ok, but `d.Hdr.Add("a", "b") is not.
// If that's confusing, recall `http.Header` is a map, therefore Hdr and Body are both pointers-to-pointers.
type BeforeRespondData struct {
	Req *http.Request
	// CacheObj is the object to be cached, containing information about the origin request. The code, headers, and body should not be considered authoritative. Look at Code, Hdr, and Body instead, as the actual values about to be sent. Note CacheObj may be nil, if an error occurred (e.g. the Origin failed to respond).
	CacheObj  *cacheobj.CacheObj
	Code      *int
	Hdr       *http.Header
	Body      *[]byte
	RemapRule string
	Context   *interface{}
}

type BeforeCacheLookUpData struct {
	Req                  *http.Request
	CacheKeyOverrideFunc func(string)
	DefaultCacheKey      string
	Context              *interface{}
}

type AfterRespondData struct {
	W         http.ResponseWriter
	Stats     stat.Stats
	RequestID uint64
	cachedata.ReqData
	cachedata.SrvrData
	cachedata.ParentRespData
	cachedata.RespData
	Context *interface{}
}

type LoadFunc func(json.RawMessage) interface{}
type StartupFunc func(icfg interface{}, d StartupData)

// LoadCacheFunc returns a map of new caches. On error, plugins should log and return a nil cache.
// If the cache only returns one cache with a key of "", it will be given the name of the plugin.
type LoadCacheFunc func(icfg interface{}, d LoadCacheData) map[string]icache.Cache

type OnRequestFunc func(icfg interface{}, d OnRequestData) bool
type BeforeCacheLookupFunc func(icfg interface{}, d BeforeCacheLookUpData)
type BeforeParentRequestFunc func(icfg interface{}, d BeforeParentRequestData)
type BeforeRespondFunc func(icfg interface{}, d BeforeRespondData)
type AfterRespondFunc func(icfg interface{}, d AfterRespondData)

type pluginObj struct {
	funcs    Funcs
	priority uint64
	name     string
}

type pluginsSlice []pluginObj

func (p pluginsSlice) Len() int           { return len(p) }
func (p pluginsSlice) Less(i, j int) bool { return p[i].priority < p[j].priority }
func (p pluginsSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var plugins = pluginsSlice{}

func Get(enabled []string) Plugins {
	enabledM := map[string]struct{}{}
	for _, name := range enabled {
		enabledM[name] = struct{}{}
	}
	enabledPlugins := pluginsSlice{}
	for _, plugin := range plugins {
		if _, ok := enabledM[plugin.name]; !ok {
			continue
		}
		enabledPlugins = append(enabledPlugins, plugin)
	}
	sort.Sort(enabledPlugins)
	return enabledPlugins
}

type Plugins interface {
	LoadFuncs() map[string]LoadFunc
	OnStartup(cfgs map[string]interface{}, context map[string]*interface{}, d StartupData)
	LoadCaches(cfgs map[string]interface{}, context map[string]*interface{}, d LoadCacheData) map[string]icache.Cache
	OnRequest(cfgs map[string]interface{}, context map[string]*interface{}, d OnRequestData) bool
	OnBeforeCacheLookup(cfgs map[string]interface{}, context map[string]*interface{}, d BeforeCacheLookUpData)
	OnBeforeParentRequest(cfgs map[string]interface{}, context map[string]*interface{}, d BeforeParentRequestData)
	OnBeforeRespond(cfgs map[string]interface{}, context map[string]*interface{}, d BeforeRespondData)
	OnAfterRespond(cfgs map[string]interface{}, context map[string]*interface{}, d AfterRespondData)
}

func (plugins pluginsSlice) LoadFuncs() map[string]LoadFunc {
	lf := map[string]LoadFunc{}
	for _, plugin := range plugins {
		if plugin.funcs.load == nil {
			continue
		}
		lf[plugin.name] = LoadFunc(plugin.funcs.load)
	}
	return lf
}

func (ps pluginsSlice) OnStartup(cfgs map[string]interface{}, context map[string]*interface{}, d StartupData) {
	for _, p := range ps {
		ictx := interface{}(nil)
		context[p.name] = &ictx

		if p.funcs.startup == nil {
			continue
		}
		d.Context = context[p.name]
		p.funcs.startup(cfgs[p.name], d)
	}
}

// LoadCaches returns the map of cache names to caches.
func (ps pluginsSlice) LoadCaches(cfgs map[string]interface{}, context map[string]*interface{}, d LoadCacheData) map[string]icache.Cache {
	log.Errorln("DEBUG pluginsSlice.LoadCaches started")
	log.Errorf("DEBUG pluginsSlice.LoadCaches len(ps) %v\n", len(ps))
	caches := map[string]icache.Cache{}
	for _, p := range ps {
		log.Errorln("DEBUG pluginsSlice.LoadCaches iterating '" + p.name + "'")
		if p.funcs.loadCache == nil {
			continue
		}
		d.Context = context[p.name]
		newCaches := p.funcs.loadCache(cfgs[p.name], d)
		for newCacheName, newCache := range newCaches {
			if newCacheName == "" && len(newCaches) == 1 {
				newCacheName = p.name
			}
			if _, ok := caches[newCacheName]; ok {
				log.Errorln("plugin loading caches: multiple caches with the name'" + newCacheName + "' - disabling both for safety! Please rename the plugin or config to have unique names!")
				delete(caches, newCacheName)
				continue
			}
			caches[newCacheName] = newCache
		}
	}
	return caches
}

// OnRequest returns a boolean whether to immediately stop processing the request. If a plugin returns true, this is immediately returned with no further plugins processed.
func (ps pluginsSlice) OnRequest(cfgs map[string]interface{}, context map[string]*interface{}, d OnRequestData) bool {
	for _, p := range ps {
		if p.funcs.onRequest == nil {
			continue
		}
		d.Context = context[p.name]
		if stop := p.funcs.onRequest(cfgs[p.name], d); stop {
			return true
		}
	}
	return false
}

func (ps pluginsSlice) OnBeforeCacheLookup(cfgs map[string]interface{}, context map[string]*interface{}, d BeforeCacheLookUpData) {
	for _, p := range ps {
		if p.funcs.beforeCacheLookUp == nil {
			continue
		}
		d.Context = context[p.name]
		p.funcs.beforeCacheLookUp(cfgs[p.name], d)
	}
}

func (ps pluginsSlice) OnBeforeParentRequest(cfgs map[string]interface{}, context map[string]*interface{}, d BeforeParentRequestData) {
	for _, p := range ps {
		if p.funcs.beforeParentRequest == nil {
			continue
		}
		d.Context = context[p.name]
		p.funcs.beforeParentRequest(cfgs[p.name], d)
	}
}

func (ps pluginsSlice) OnBeforeRespond(cfgs map[string]interface{}, context map[string]*interface{}, d BeforeRespondData) {
	for _, p := range ps {
		if p.funcs.beforeRespond == nil {
			continue
		}
		d.Context = context[p.name]
		p.funcs.beforeRespond(cfgs[p.name], d)
	}
}

func (ps pluginsSlice) OnAfterRespond(cfgs map[string]interface{}, context map[string]*interface{}, d AfterRespondData) {
	for _, p := range ps {
		if p.funcs.afterRespond == nil {
			continue
		}
		d.Context = context[p.name]
		p.funcs.afterRespond(cfgs[p.name], d)
	}
}
