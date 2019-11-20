package v14

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
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const EdgeHdrRwPrefix = "hdr_rw"
const MidHdrRwPrefix = "hdr_rw_mid"

func TestHdrRwDotConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		GetTestHdrRwDotConfig(t)
		GetTestHdrRwMidDotConfig(t)
		GetTestHdrRwDotConfigWithNewline(t)
	})
}

func getFirstDnsOrHttpDeliveryService(t *testing.T) *tc.DeliveryServiceNullable {
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Errorf("Cannot test hdr_rw_dot_config with no http or dns deliveryservices: %s", err)
		return nil
	}

	for _, ds := range dses {
		switch *ds.Type {
		case tc.DSTypeDNS:
		case tc.DSTypeDNSLive:
		case tc.DSTypeDNSLiveNational:
		case tc.DSTypeHTTP:
		case tc.DSTypeHTTPLive:
		case tc.DSTypeHTTPLiveNational:
		case tc.DSTypeHTTPNoCache:
		default:
			continue
		}

		return &ds
	}

	t.Errorf("Cannot test hdr_rw_dot_config with no http or dns deliveryservices: %s", err)
	return nil

}

func getExpectedLines(rwRules string) int {
	if rwRules == "" {
		return 1 // for the header comment
	}
	return 2 + strings.Count(rwRules, "__RETURN__") + strings.Count(rwRules, "\n")
}

func GetTestHdrRwDotConfigWithNewline(t *testing.T) {
	ds := getFirstDnsOrHttpDeliveryService(t)
	*ds.EdgeHeaderRewrite = "rw1\nrw2\nedge\nheader\nre-rewrite [L]"
	_, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*ds.ID), ds)
	if err != nil {
		t.Errorf("couldn't update delivery servie: %v", err)
	}

	filename := fmt.Sprintf("%s_%s.config", EdgeHdrRwPrefix, *ds.XMLID)
	config, _, _ := TOSession.GetATSCDNConfig(*ds.CDNID, filename)

	expectedLines := getExpectedLines(*ds.EdgeHeaderRewrite)
	count := strings.Count(config, "\n")
	if expectedLines != count {
		t.Errorf("expected %d lines in the config (actual = %d)", expectedLines, count)
	} else {
		log.Debugf("Tested %s sucessfully\n", filename)
	}
}

func GetTestHdrRwDotConfig(t *testing.T) {
	ds := getFirstDnsOrHttpDeliveryService(t)
	*ds.EdgeHeaderRewrite = "rw1__RETURN__rw2__RETURN__edge__RETURN__header__RETURN__re-rewrite [L]"
	_, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*ds.ID), ds)
	if err != nil {
		t.Errorf("couldn't update delivery servie: %v", err)
	}

	filename := fmt.Sprintf("%s_%s.config", EdgeHdrRwPrefix, *ds.XMLID)
	config, _, _ := TOSession.GetATSCDNConfig(*ds.CDNID, filename)

	expectedLines := getExpectedLines(*ds.EdgeHeaderRewrite)
	count := strings.Count(config, "\n")
	if expectedLines != count {
		t.Errorf("expected %d lines in the config (actual = %d)", expectedLines, count)
	} else {
		log.Debugf("Tested %s sucessfully\n", filename)
	}
}

func GetTestHdrRwMidDotConfig(t *testing.T) {
	ds := getFirstDnsOrHttpDeliveryService(t)
	*ds.MidHeaderRewrite = "rw1__RETURN__mid__RETURN__header__RETURN__re-rewrite [L]"
	_, err := TOSession.UpdateDeliveryServiceNullable(strconv.Itoa(*ds.ID), ds)
	if err != nil {
		t.Errorf("couldn't update delivery servie: %v", err)
	}

	filename := fmt.Sprintf("%s_%s.config", MidHdrRwPrefix, *ds.XMLID)
	config, _, _ := TOSession.GetATSCDNConfig(*ds.CDNID, filename)

	expectedLines := getExpectedLines(*ds.MidHeaderRewrite)
	count := strings.Count(config, "\n")
	if expectedLines != count {
		t.Errorf("expected %d lines in the config (actual = %d)", expectedLines, count)
	} else {
		log.Debugf("Tested %s sucessfully\n", filename)
	}
}
