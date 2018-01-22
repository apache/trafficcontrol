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

package towrap

import (
	"net"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/config"
)

func SetupSession(cfg config.Config, toURL string, toUser string, toPass string) (*client.Session, net.Addr, error) {
	var err error
	var session *client.Session
	var netAddr net.Addr
	toReqTimeout := time.Second * time.Duration(cfg.Default.Session.TimeoutInSecs)
	session, netAddr, err = to.LoginWithAgent(toURL, toUser, toPass, true, "to-api-client-tests", true, toReqTimeout)
	if err != nil {
		return nil, nil, err
	}

	return session, netAddr, err
}
