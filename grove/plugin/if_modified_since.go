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

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
)

func init() {
	AddPlugin(5000, Funcs{beforeRespond: ifModifiedSince})
}

func ifModifiedSince(icfg interface{}, d BeforeRespondData) {
	if d.CacheObj == nil {
		return // if we don't have a cacheobj from the origin, there's no object to have been modified.
	}
	modifiedSince, ok := rfc.GetHTTPDate(d.Req.Header, "If-Modified-Since")
	if !ok {
		return
	}
	if d.CacheObj.LastModified.After(modifiedSince) {
		return
	}
	*d.Code, *d.Hdr, *d.Body = http.StatusNotModified, nil, nil
}
