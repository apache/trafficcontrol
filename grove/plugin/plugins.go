package plugin

import (
	"encoding/json"

	"github.com/apache/incubator-trafficcontrol/grove/plugin/afterrespond"
	"github.com/apache/incubator-trafficcontrol/grove/plugin/beforeparentrequest"
	"github.com/apache/incubator-trafficcontrol/grove/plugin/beforerespond"
	"github.com/apache/incubator-trafficcontrol/grove/plugin/onrequest"
	"github.com/apache/incubator-trafficcontrol/grove/plugin/startup"
)

// Plugins contains plugins functions. Plugins are sorted according to their registered priority, and should be called in the order they appear in the slice.
type Plugins struct {
	Startup             startup.Startup
	BeforeParentRequest beforeparentrequest.Plugin
	BeforeRespond       beforerespond.Plugin
	AfterRespond        afterrespond.Plugin
	OnRequest           onrequest.Plugin
}

// Get gets all the plugins. This must not be called in an init function, since plugins use init for registration. This must be called after initialization, after main has started executing.
func Get() Plugins {
	return Plugins{
		Startup:             startup.Get(),
		BeforeParentRequest: beforeparentrequest.Get(),
		BeforeRespond:       beforerespond.Get(),
		AfterRespond:        afterrespond.Get(),
		OnRequest:           onrequest.Get(),
	}
}

// LoadFuncs() returns the load functions of all plugins.
//
// Plugin loading acts like other remap configuration, with configuration for each plugin being used from each remap rule if it exists, and the global plugin configuration being used if no plugin value exists for a particular rule.
//
// Note there is currently no way to override only parts of a plugin config. For example, if the global plugins key for the `mod_request_headers` plugin has a `drop`, and a rule has an `add`, only the add will apply, because the rule's plugin is used in its entirety if it exists.
//
func (p Plugins) LoadFuncs() map[string]PluginLoadF {
	lf := map[string]PluginLoadF{}
	for name, f := range p.BeforeParentRequest.LoadFuncs() {
		lf[name] = PluginLoadF(f)
	}
	for name, f := range p.BeforeRespond.LoadFuncs() {
		lf[name] = PluginLoadF(f)
	}
	for name, f := range p.AfterRespond.LoadFuncs() {
		lf[name] = PluginLoadF(f)
	}
	for name, f := range p.OnRequest.LoadFuncs() {
		lf[name] = PluginLoadF(f)
	}
	return lf
}

type PluginLoadF func(msg json.RawMessage) interface{}
