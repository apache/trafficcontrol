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

	"github.com/apache/trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{onStartup: helloSharedConfigStartup}, "example plugin for loading and using config data for other plugins", "1.0.0")
}

func helloSharedConfigStartup(d StartupData) {
	if b, err := json.Marshal(d.SharedCfg); err == nil {
		log.Debugln("Hello! This is a shared plugin data config! Your shared plugin config is: " + string(b))
	} else {
		log.Debugf("Hello! This is a shared plugin data config! Your shared plugin config is: %+v\n", d.SharedCfg)
	}
}
