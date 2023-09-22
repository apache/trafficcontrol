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

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
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
	ModifyFiles(d ModifyFilesData) []t3cutil.ATSConfigFile
}

func AddPlugin(priority uint64, funcs Funcs) {
	// print used here, because this is called in init, loggers don't exist yet
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Println(time.Now().Format(time.RFC3339Nano) + " Error plugin.AddPlugin: runtime.Caller failed, can't get plugin names")
		os.Exit(1)
	}

	pluginName := strings.TrimSuffix(path.Base(filename), ".go")
	fmt.Println("AddPlugin adding " + pluginName)
	initPlugins = append(initPlugins, pluginObj{funcs: funcs, priority: priority, name: pluginName})
}

type Funcs struct {
	onStartup   StartupFunc
	modifyFiles ModifyFilesFunc
}

type StartupData struct {
	Cfg config.Cfg
}

type ModifyFilesData struct {
	Cfg    config.Cfg
	TOData *t3cutil.ConfigData
	Files  []t3cutil.ATSConfigFile
}

type IsRequestHandled bool

const (
	RequestHandled   = IsRequestHandled(true)
	RequestUnhandled = IsRequestHandled(false)
)

type StartupFunc func(d StartupData)
type ModifyFilesFunc func(d ModifyFilesData) []t3cutil.ATSConfigFile

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

// ModifyFiles returns a slice of config files to use. May return d.Files unmodified, or may add, remove, or modify files in d.Files.
func (ps plugins) ModifyFiles(d ModifyFilesData) []t3cutil.ATSConfigFile {
	log.Infof("plugins.ModifyFiles calling %+v plugins\n", len(ps.slice))
	for _, p := range ps.slice {
		if p.funcs.modifyFiles == nil {
			log.Infoln("plugins.ModifyFiles plugging " + p.name + " - no modifyFiles func")
			continue
		}
		log.Infoln("plugins.ModifyFiles plugging " + p.name)
		d.Files = p.funcs.modifyFiles(d)
	}
	return d.Files
}
