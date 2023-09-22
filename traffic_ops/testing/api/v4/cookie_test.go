package v4

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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func TestCookies(t *testing.T) {
	WithObjs(t, []TCObj{CDNs}, func() {
		CookiesTest(t)
	})
}

func CookiesTest(t *testing.T) {
	s, _, err := toclient.LoginWithAgent(Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword, true, "to-api-v4-client-tests", false, toReqTimeout)
	credentials := tc.UserCredentials{
		Username: Config.TrafficOps.Users.Admin,
		Password: Config.TrafficOps.UserPassword,
	}

	js, err := json.Marshal(credentials)
	if err != nil {
		t.Fatal("unable to json marshal login credentials")
	}
	path := TestAPIBase + "/user/login"
	loginResp, _, err := s.RawRequest(http.MethodPost, path, js)
	if err != nil {
		t.Fatal("unable to request POST /user/login")
	}
	defer loginResp.Body.Close()
	_, readErr := ioutil.ReadAll(loginResp.Body)
	if readErr != nil {
		t.Fatal("unable to read response body from POST /user/login")
	}
	ensureCookie(loginResp, t)

	cdnResp, _, err := s.RawRequest(http.MethodGet, TestAPIBase+"/cdns", nil)
	if err != nil {
		t.Fatal("unable to request GET /cdns")
	}
	defer cdnResp.Body.Close()
	_, readErr = ioutil.ReadAll(cdnResp.Body)
	if readErr != nil {
		t.Fatal("unable to read response body from GET /cdns")
	}
	ensureCookie(cdnResp, t)
}

func ensureCookie(r *http.Response, t *testing.T) {
	cookies := r.Cookies()
	if len(cookies) < 1 {
		t.Fatal("expected at least one cookie in response, actual: zero")
	}
	if cookies[0].MaxAge < 1 {
		t.Errorf("expected auth cookie Max-Age > 0, actual: %v", *cookies[0])
	}
	if cookies[0].Expires.IsZero() {
		t.Errorf("expected auth cookie with non-zero Expires, actual: %v", *cookies[0])
	}
}
