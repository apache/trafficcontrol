package startup

import (
	"sort"

	"github.com/apache/incubator-trafficcontrol/grove/config"
)

// AddPlugin registers a startup plugin, which will be called when the application starts.
//
// Startup plugins may do anything that needs to happen once on startup.
//
// Examples:
// * Starting a long-running goroutine which polls something
// * Initializing data for another plugin
//
// TODO document exactly when Startup plugins are called
func AddPlugin(priority uint64, f Func) {
	plugins = append(plugins, plugin{f: f, priority: priority})
}

type Func func(cfg config.Config)

type plugin struct {
	f        Func
	priority uint64
}

type pluginSlice []plugin

func (p pluginSlice) Len() int           { return len(p) }
func (p pluginSlice) Less(i, j int) bool { return p[i].priority < p[j].priority }
func (p pluginSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var plugins = pluginSlice{}

// Get gets all the plugins. This must not be called in an init function, since plugins use init for registration. This must be called after main has started executing.
func Get() Startup {
	sort.Sort(plugins)
	ps := []Func{}
	for _, p := range plugins {
		ps = append(ps, p.f)
	}
	return funcs(ps)
}

type Startup interface {
	Call(cfg config.Config)
}

type funcs []Func

func (fs funcs) Call(cfg config.Config) {
	for _, f := range fs {
		f(cfg)
	}
}
