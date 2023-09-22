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
	AddPlugin(10000, Funcs{onStartup: helloCtxStart, onRequest: helloCtxOnReq}, "example plugin for passing context data between hook functions", "1.0.0")
}

func helloCtxStart(d StartupData) {
	*d.Ctx = 42
	log.Debugf("Hello! This is a context plugin! Start set context: %+v\n", *d.Ctx)
}

func helloCtxOnReq(d OnRequestData) IsRequestHandled {
	ctx, ok := (*d.Ctx).(int)
	log.Debugf("Hello! This is a context plugin! On Request got context: %+v %+v\n", ok, ctx)
	return RequestUnhandled
}
