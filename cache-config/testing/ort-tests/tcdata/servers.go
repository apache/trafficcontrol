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
	"fmt"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func (r *TCData) QueueUpdatesForServer(hostname string, queue bool) error {
	respServer, _, err := toreq.GetServerByHostName(TOSession, hostname)
	if err != nil {
		return fmt.Errorf("cannot GET Server by hostname '%s': %v", hostname, err)
	}
	if respServer.ID == 0 {
		return fmt.Errorf("server '%s' had nil ID", hostname)
	}
	if _, _, err := TOSession.SetServerQueueUpdate(respServer.ID, queue, toclient.RequestOptions{}); err != nil {
		return err
	}
	return nil
}
