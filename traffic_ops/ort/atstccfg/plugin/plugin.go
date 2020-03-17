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
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

// List returns the list of plugins compiled into the calling executable.
func List() []string {
	l := []string{}
	for _, p := range initPlugins {
		l = append(l, p.name)
	}
	return l
}

func Get(appCfg config.Cfg) Plugins {
	pluginSlice := getAll()
	return plugins{slice: pluginSlice}
}

func getAll() pluginsSlice {
	enabledPlugins := pluginsSlice{}
	for _, plugin := range initPlugins {
		enabledPlugins = append(enabledPlugins, plugin)
	}
	sort.Sort(enabledPlugins)
	return enabledPlugins
}

type Plugins interface {
	OnStartup(d StartupData)
	OnRequest(d OnRequestData) bool
}

func AddPlugin(priority uint64, funcs Funcs) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error plugin.AddPlugin: runtime.Caller failed, can't get plugin names") // print, because this is called in init, loggers don't exist yet
		os.Exit(1)
	}

	pluginName := strings.TrimSuffix(path.Base(filename), ".go")
	log.Infoln("AddPlugin adding " + pluginName)
	initPlugins = append(initPlugins, pluginObj{funcs: funcs, priority: priority, name: pluginName})
}

type Funcs struct {
	onStartup StartupFunc
	onRequest OnRequestFunc
}

type StartupData struct {
	Cfg config.Cfg
}

type OnRequestData struct {
	Cfg config.TCCfg
}

type IsRequestHandled bool

const (
	RequestHandled   = IsRequestHandled(true)
	RequestUnhandled = IsRequestHandled(false)
)

type StartupFunc func(d StartupData)
type OnRequestFunc func(d OnRequestData) IsRequestHandled

type pluginObj struct {
	funcs    Funcs
	priority uint64
	name     string
}

type plugins struct {
	slice pluginsSlice
	cfg   map[string]interface{}
	ctx   map[string]*interface{}
}

type pluginsSlice []pluginObj

func (p pluginsSlice) Len() int           { return len(p) }
func (p pluginsSlice) Less(i, j int) bool { return p[i].priority < p[j].priority }
func (p pluginsSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// initPlugins is where plugins are registered via their init functions.
var initPlugins = pluginsSlice{}

func (ps plugins) OnStartup(d StartupData) {
	for _, p := range ps.slice {
		if p.funcs.onStartup == nil {
			continue
		}
		p.funcs.onStartup(d)
	}
}

// OnRequest returns a boolean whether to immediately stop processing the request. If a plugin returns true, this is immediately returned with no further plugins processed.
func (ps plugins) OnRequest(d OnRequestData) bool {
	log.Infof("plugins.OnRequest calling %+v plugins\n", len(ps.slice))
	for _, p := range ps.slice {
		if p.funcs.onRequest == nil {
			log.Infoln("plugins.OnRequest plugging " + p.name + " - no onRequest func")
			continue
		}
		log.Infoln("plugins.OnRequest plugging " + p.name)
		if stop := p.funcs.onRequest(d); stop {
			return true
		}
	}
	return false
}
