package v13

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

	"github.com/apache/incubator-trafficcontrol/traffic_ops/client/v13"
	_ "github.com/lib/pq"
)

var (
	TOSession       *v13.Session
	NoAuthTOSession *v13.Session
)

func SetupSession(toReqTimeout time.Duration, toURL string, toUser string, toPass string) error {
	var err error

	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	NoAuthTOSession = v13.NewNoAuthSession(toURL, true, "to-api-v13-client-tests", true, toReqTimeout)
	TOSession, _, err = v13.LoginWithAgent(toURL, toUser, toPass, true, "to-api-v13-client-tests", true, toReqTimeout)
	return err
}

func TeardownSession(toReqTimeout time.Duration, toURL string, toUser string, toPass string) error {
	var err error
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	TOSession, _, err = v13.LogoutWithAgent(toURL, toUser, toPass, true, "to-api-v13-client-tests", true, toReqTimeout)

	return err
}

func SwitchSession(toReqTimeout time.Duration, toURL string, toOldUser string, toOldPass string, toNewUser string, toNewPass string) error {
	err := TeardownSession(toReqTimeout, toURL, toOldUser, toOldPass)

	// intentially skip errors so that we can continue with setup in the event of a 403

	err = SetupSession(toReqTimeout, toURL, toNewUser, toNewPass)
	return err
}
