package plugin

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

	"github.com/apache/incubator-trafficcontrol/grove/cachedata"
	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"
	"github.com/apache/incubator-trafficcontrol/grove/config"
	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/stat"
	"github.com/apache/incubator-trafficcontrol/grove/web"
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
	onRequest           OnRequestFunc
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

type OnRequestData struct {
	W             http.ResponseWriter
	R             *http.Request
	InterfaceName string
	Stats         stat.Stats
	StatRules     remapdata.RemapRulesStats
	HTTPConns     *web.ConnMap
	HTTPSConns    *web.ConnMap
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

type AfterRespondData struct {
	W     http.ResponseWriter
	Stats stat.Stats

	cachedata.ReqData
	cachedata.SrvrData
	cachedata.ParentRespData
	cachedata.RespData
	Context *interface{}
}

type LoadFunc func(json.RawMessage) interface{}
type StartupFunc func(icfg interface{}, d StartupData)
type OnRequestFunc func(icfg interface{}, d OnRequestData) bool
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

func Get() Plugins {
	sort.Sort(plugins)
	return pluginsSlice(plugins)
}

type Plugins interface {
	LoadFuncs() map[string]LoadFunc
	OnStartup(cfgs map[string]interface{}, context map[string]*interface{}, d StartupData)
	OnRequest(cfgs map[string]interface{}, context map[string]*interface{}, d OnRequestData) bool
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
