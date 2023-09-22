package tcdata

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
	"time"

	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"

	_ "github.com/lib/pq"
)

var (
	TOSession       *client.Session
	NoAuthTOSession *client.Session
)

func (r *TCData) SetupSession(toReqTimeout time.Duration, toURL string, toUser string, toPass string) error {
	var err error

	toReqTimeout = time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	NoAuthTOSession = client.NewNoAuthSession(toURL, true, "to-api-v4-client-tests", true, toReqTimeout)
	TOSession, _, err = client.LoginWithAgent(toURL, toUser, toPass, true, "to-api-v4-client-tests", true, toReqTimeout)
	return err
}

func (r *TCData) TeardownSession(toReqTimeout time.Duration, toURL string, toUser string, toPass string) error {
	var err error
	toReqTimeout = time.Second * time.Duration(r.Config.Default.Session.TimeoutInSecs)
	TOSession, _, err = client.LogoutWithAgent(toURL, toUser, toPass, true, "to-api-v4-client-tests", true, toReqTimeout)

	return err
}

func (r *TCData) SwitchSession(toReqTimeout time.Duration, toURL string, toOldUser string, toOldPass string, toNewUser string, toNewPass string) error {
	err := r.TeardownSession(toReqTimeout, toURL, toOldUser, toOldPass)

	// intentially skip errors so that we can continue with setup in the event of a 403

	err = r.SetupSession(toReqTimeout, toURL, toNewUser, toNewPass)
	return err
}
