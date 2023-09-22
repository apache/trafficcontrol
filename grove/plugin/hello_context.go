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
	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{startup: helloCtxStart, afterRespond: helloCtxAfterResp})
}

func helloCtxStart(icfg interface{}, d StartupData) {
	*d.Context = 42
	log.Debugf("Hello World! Start set context: %+v\n", *d.Context)
}

func helloCtxAfterResp(icfg interface{}, d AfterRespondData) {
	ctx, ok := (*d.Context).(int)
	log.Debugf("Hello World! After Response got context: %+v %+v\n", ok, ctx)
}
