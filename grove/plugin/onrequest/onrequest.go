package onrequest

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/stat"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

type Data struct {
	W             http.ResponseWriter
	R             *http.Request
	InterfaceName string
	Stats         stat.Stats
	StatRules     remapdata.RemapRulesStats
	HTTPConns     *web.ConnMap
	HTTPSConns    *web.ConnMap
}

// AddPlugin registers a request plugin, which will be called when a request is received.
//
// The most common use case for OnRequest plugins, are custom endpoints.
//
// It is highly recommended that custom endpoints be given a unique prefix, such as "_grove_plugin_foo" or even a UUID such as "_3c601922-c1f7-476a-aace-0e11020bd476_foo", to reduce the chance of conflicting with remap rules.
//
// Examples:
// * endpoint to get stats
// * endpoint to purge a cache entry
//
func AddPlugin(priority uint64, name string, f Func, loadF LoadFunc) {
	plugins = append(plugins, plugin{f: f, priority: priority, name: name, loadF: loadF})
}

// Func is called immediately when a request is received. It is given the config data loaded by its LoadFunc, and data including cache state and the HTTP Request and ResponseWriter.
// Returns whether the cache handler should return without processing the request. If true is returned, the cache handler immediately returns. Plugins should only ever return true if they have completely handled the request, and returned an appropriate response to the client.
type Func func(icfg interface{}, d Data) bool

// The LoadFunc is the function which loads any necessary configuration for the plugin. This config should be placed in the remap rules file, in the "plugins" object, under the key with the name of this plugin. Both keys within remap rules, and in the outer object will be passed to this function. As with all remap rules, if the object exists for a specific rule, it will be used, and in this case passed to the plugin call func; if the rule object doesn't exist, the outer object will be used.
//
// Each plugin's LoadFunc should take a json.RawMessage, build its config object, and return it. That object will be passed to the Func when called. If a particular plugin doesn't need any remap config data, it should return nil, or pass a nil LoadFunc.
type LoadFunc func(json.RawMessage) interface{}

type plugin struct {
	f        Func
	loadF    LoadFunc
	priority uint64
	name     string
}

type pluginSlice []plugin

func (p pluginSlice) Len() int           { return len(p) }
func (p pluginSlice) Less(i, j int) bool { return p[i].priority < p[j].priority }
func (p pluginSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var plugins = pluginSlice{}

// Get gets all the plugins. This must not be called in an init function, since plugins use init for registration. This must be called after initialization, after main has started executing.
func Get() Plugin {
	sort.Sort(plugins)
	return pluginSlice(plugins)
}

type Plugin interface {
	// Call executes this plugin. Note the header and body objects MUST NOT be modified. If this plugin modifies the response header or body, set the object pointed to to a new object, which must be copied from the original if necessary.
	Call(cfgs map[string]interface{}, d Data) bool
	LoadFuncs() map[string]LoadFunc
}

func (ps pluginSlice) LoadFuncs() map[string]LoadFunc {
	// TODO change to slice? Slice is faster, since we always iterate; but map is more intuitive
	lf := map[string]LoadFunc{}
	for _, p := range ps {
		if p.loadF != nil {
			lf[p.name] = p.loadF
		}
	}
	return lf
}

func (ps pluginSlice) Call(cfgs map[string]interface{}, d Data) bool {
	// TODO implement plugins signalling whether they'll modify, in order to only copy once here.
	log.Debugf("afterrespond.pluginSlice.Call looping over %+v cfgs %+v\n", len(ps), cfgs)
	if cfgs == nil {
		// easier and probably faster to make a map that returns nil for everything, than to check if cfgs is nil every time
		cfgs = map[string]interface{}{}
	}

	for _, p := range ps {
		log.Debugf("afterrespond.pluginSlice.Call calling %+v\n", p.name)
		if stop := p.f(cfgs[p.name], d); stop {
			return true
		}
	}
	return false
}
