package v5

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
	"fmt"
	"os"
	"testing"
	"time"

	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"

	_ "github.com/lib/pq"
)

var (
	TOSession       *client.Session
	NoAuthTOSession *client.Session
)

func SetupSession(toReqTimeout time.Duration, toURL string, toUser string, toPass string) error {
	var err error

	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	NoAuthTOSession = client.NewNoAuthSession(toURL, true, "to-api-v5-client-tests", true, toReqTimeout)
	TOSession, _, err = client.LoginWithAgent(toURL, toUser, toPass, true, "to-api-v5-client-tests", true, toReqTimeout)
	return err
}

func TeardownSession(toReqTimeout time.Duration, toURL string, toUser string, toPass string) error {
	var err error
	toReqTimeout = time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	TOSession, _, err = client.LogoutWithAgent(toURL, toUser, toPass, true, "to-api-v5-client-tests", true, toReqTimeout)

	return err
}

func SwitchSession(toReqTimeout time.Duration, toURL string, toOldUser string, toOldPass string, toNewUser string, toNewPass string) error {
	err := TeardownSession(toReqTimeout, toURL, toOldUser, toOldPass)

	// intentionally skip errors so that we can continue with setup in the event of a 403

	err = SetupSession(toReqTimeout, toURL, toNewUser, toNewPass)
	return err
}

func TestLoginWithCert(t *testing.T) {
	if includeSystemTests {
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("couldn't get current working directory: %s", err.Error())
		}

		session, _, err := client.LoginWithCert(Config.TrafficOps.URL, true, time.Second*60,
			fmt.Sprintf("%s/client-intermediate-chain.crt.pem", pwd),
			fmt.Sprintf("%s/client.key.pem", pwd), "")

		if err != nil {
			t.Fatalf("expected no error while logging in with cert, but got %v", err)
		}
		if session == nil {
			t.Fatalf("expected a valid session, but got nothing")
		}

		_, _, err = session.GetAbout(client.RequestOptions{})
		if err != nil {
			t.Fatalf("expected no error while using the client cert session to hit an authenticated endpoint, but got %v", err)
		}
	}
}
