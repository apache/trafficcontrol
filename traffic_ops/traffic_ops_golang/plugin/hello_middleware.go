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
	"net/http"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/routing/middleware"
)

func init() {
	AddPlugin(10000, Funcs{onRequest: helloMiddleware}, "example middleware plugin", "1.0.0")
}

const HelloMiddlewarePath = "/_hello_middleware"

func helloMiddleware(d OnRequestData) IsRequestHandled {
	if !strings.HasPrefix(d.R.URL.Path, HelloPath) {
		return RequestUnhandled
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!\n"))
	}
	requestTimeout := middleware.DefaultRequestTimeout
	if d.AppCfg.RequestTimeout != 0 {
		requestTimeout = time.Second * time.Duration(d.AppCfg.RequestTimeout)
	}
	mw := middleware.GetDefault(d.AppCfg.Secrets[0], requestTimeout)
	handler = middleware.Use(handler, mw)
	handler(d.W, d.R)

	return RequestHandled
}
