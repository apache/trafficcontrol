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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{load: helloConfigLoad, onStartup: helloConfigStartup}, "example plugin for loading and using config file data", "1.0.0")
}

type HelloConfig struct {
	Hello string `json:"hello"`
}

func helloConfigLoad(b json.RawMessage) interface{} {
	cfg := HelloConfig{}
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		log.Debugln(`Hello! This is a config plugin! Unfortunately, your config JSON is not properly formatted. Config should look like: {"plugin_config": {"hello_config":{"hello": "anything can go here"}}}`)
		return nil
	}
	log.Debugln("Hello! This is a config plugin! Successfully loaded config!")
	return &cfg
}

func helloConfigStartup(d StartupData) {
	if d.Cfg == nil {
		log.Debugln("Hello! This is a config plugin! Unfortunately, your config is not set properly.")
	}
	cfg, ok := d.Cfg.(*HelloConfig)
	if !ok {
		// should never happen
		log.Debugf("helloLoadConfig config '%v' type '%T' expected *HelloConfig\n", d.Cfg, d.Cfg)
		return
	}
	log.Debugf("Hello! This is a config plugin! Your config is: %+v\n", cfg)
}
