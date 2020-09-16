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
	"strings"
)

func init() {
	AddPlugin(10000, Funcs{onRequest: hello}, "example HTTP plugin", "1.0.0")
}

const HelloPath = "/_hello"

func hello(d OnRequestData) IsRequestHandled {
	if !strings.HasPrefix(d.R.URL.Path, HelloPath) {
		return RequestUnhandled
	}
	d.W.Header().Set("Content-Type", "text/plain")
	d.W.Write([]byte("Hello, World!"))
	return RequestHandled
}
