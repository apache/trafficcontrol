package beforerespond

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/apache/incubator-trafficcontrol/grove/cacheobj"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)


// Data holds the data passed to plugins. The objects pointed to MAY NOT be modified, however, the location pointed to may be changed for the Code, Hdr, and Body. That iss, `*d.Hdr = myHdr` is ok, but `d.Hdr.Add("a", "b") is not.
//
// If that's confusing, recall `http.Header` is a map, therefore Hdr and Body are both pointers-to-pointers.
type Data struct {
	Req      *http.Request
	// CacheObj is the object to be cached, containing information about the origin request. The code, headers, and body should not be considered authoritative. Look at Code, Hdr, and Body instead, as the actual values about to be sent. Note CacheObj may be nil, if an error occurred (e.g. the Origin failed to respond).
	CacheObj *cacheobj.CacheObj
	Code     *int
	Hdr      *http.Header
	Body     *[]byte
}

// AddPlugin registers a request plugin, which will be called when a request is received.
//
// Request plugins are called immediately after a request is received, and before remappings are processed.
//
// Request plugins may manipulate the request, and any manipulations will be used by all further processing.
//
// Examples:
// * Add or delete request headers
// * Set the request RemoteAddr to X-Real-IP or X-Forwarded-For headers
//
func AddPlugin(priority uint64, name string, f Func, loadF LoadFunc) {
	plugins = append(plugins, plugin{f: f, priority: priority, name: name, loadF: loadF})
}

type Func func(icfg interface{}, d Data)

// The LoadFunc is the function which loads any necessary configuration for the plugin. This config  should be placed in the remap rules file, in the "plugins" object, under the key with the name of this plugin. Both keys within remap rules, and in the outer object will be passed to this function. As with all remap rules, if the object exists for a specific rule, it will be used, and in this case passed to the plugin call func; if the rule object doesn't exist, the outer object will be used.
//
// Each plugin's LoadFunc should take a json.RawMessage, build its config object, and return it. That object will be passed to the Func when called. If a particular plugin doesn't need any remap config data, it should return nil.
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
	Call(cfgs map[string]interface{}, d Data)
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

func (ps pluginSlice) Call(cfgs map[string]interface{}, d Data) {
	// TODO implement plugins signalling whether they'll modify, in order to only copy once here.
	log.Debugf("beforerespond.pluginSlice.Call looping over %+v cfgs %+v\n", len(ps), cfgs)
	for _, p := range ps {
		log.Debugf("beforerespond.pluginSlice.Call calling %+v\n", p.name)
		p.f(cfgs[p.name], d)
	}
}
