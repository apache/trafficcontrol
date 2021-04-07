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

package cachecfg_test

import (
	"fmt"
	"testing"

	"github.com/apache/trafficcontrol/traffic_ops/testing/api/v1/config"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/v1/config/cachecfg"
)

var commonNegativeTests = []config.NegativeTest{
	{
		Description: "too few assignments",
		Config:      "foo1=foo",
		Expected:    config.NotEnoughAssignments,
	},
	{
		Description: "empty assignment",
		Config:      "dest_domain= foo1=foo",
		Expected:    config.BadAssignmentMatch,
	},
	{
		Description: "missing equals in assignment",
		Config:      "dest_domain foo1=foo",
		Expected:    config.BadAssignmentMatch,
	},
	{
		Description: "more than one primary destination",
		Config:      "dest_domain=example dest_domain=example",
		Expected:    config.ExcessLabel,
	},
	{
		Description: "using a single secondary specifier twice",
		Config:      "dest_domain=example scheme=http scheme=https",
		Expected:    config.ExcessLabel,
	},
	{
		Description: "missing primary destination label",
		Config:      "port=80 method=get time=08:00-14:00",
		Expected:    config.MissingLabel,
	},
	{
		Description: "unknown primary field",
		Config:      "foo1=foo foo2=foo foo3=foo",
		Expected:    config.InvalidLabel,
	},
	{
		Description: "good first line, but bad second",
		Config: fmt.Sprintf("%s\n\n%s",
			"dest_domain=example.com suffix=js revalidate=1d",
			"foo1=foo foo2=foo foo3=foo"),
		Expected: config.InvalidLabel,
	},
}

var primaryDestinationNegativeTests = []config.NegativeTest{
	{
		Description: "bad url regex",
		Config:      "url_regex=(example.com/.* foo2=foo",
		Expected:    config.InvalidRegex,
	},
	{
		Description: "bad host regex",
		Config:      "host_regex=(example.* foo2=foo",
		Expected:    config.InvalidRegex,
	},
	{
		Description: "bad hostname",
		Config:      "dest_domain=my%20bad%20domain.com foo1=foo",
		Expected:    config.InvalidHost,
	},
	{
		Description: "bad ip",
		Config:      "dest_ip=bad_ip foo1=foo",
		Expected:    config.InvalidIP,
	},
}

var secondarySpecifierNegativeTests = []config.NegativeTest{
	{
		Description: "bad port",
		Config:      "dest_domain=example port=90009000",
		Expected:    config.InvalidPort,
	},
	{
		Description: "bad scheme",
		Config:      "dest_domain=example scheme=httpz",
		Expected:    config.InvalidHTTPScheme,
	},
	{
		Description: "bad method",
		Config:      "dest_domain=example method=xxx",
		Expected:    config.UnknownMethod,
	},
	{
		Description: "bad time range",
		Config:      "dest_domain=example time=16:00",
		Expected:    config.InvalidTimeRange24Hr,
	},
	{
		Description: "bad src ip",
		Config:      "dest_domain=example src_ip=bad_ip",
		Expected:    config.InvalidIP,
	},
	{
		Description: "bad boolean value",
		Config:      "dest_domain=example internal=xxx",
		Expected:    config.InvalidBool,
	},
}

var actionNegativeTests = []config.NegativeTest{
	{
		Description: "bad action value",
		Config:      "dest_domain=example action=xxx",
		Expected:    config.InvalidAction,
	},
	{
		Description: "bad cache-responses-to-cookies",
		Config:      "dest_domain=example cache-responses-to-cookies=42",
		Expected:    config.InvalidCacheCookieResponse,
	},
	{
		Description: "bad time format",
		Config:      "dest_domain=example pin-in-cache=xxx",
		Expected:    config.InvalidTimeFormatDHMS,
	},
	{
		Description: "missing action label",
		Config:      "dest_domain=example scheme=http",
		Expected:    config.MissingLabel,
	},
}

var positiveTests = []config.PositiveTest{
	{
		Description: "empty config",
	},
	{
		Description: "empty multiline config",
		Config:      "\n",
	},
	{
		Description: "empty config with whitespace",
		Config:      "\t ",
	},
	{
		Description: "normal config returned from traffic ops",
		Config: fmt.Sprintf("%s\n%s",
			"# DO NOT EDIT - Generated for ATS_EDGE_TIER_CACHE by Traffic Ops on Fri Feb 15 22:01:53 UTC 2019",
			"dest_domain=origin.infra.ciab.test port=80 scheme=http action=never-cache"),
	},
	{
		Description: "tab-delimitted config",
		Config:      "dest_domain=origin.infra.ciab.test\tport=80\tscheme=http\taction=never-cache",
	},
	{
		Description: "multi-space delimitted config",
		Config:      "dest_domain=origin.infra.ciab.test  port=80  scheme=http  action=never-cache",
	},
	{
		Description: "multiline config with empty line",
		Config: fmt.Sprintf("%s\n\n%s",
			"dest_domain=example.com suffix=js revalidate=1d",
			"dest_domain=example.com prefix=foo suffix=js revalidate=7d"),
	},
	{
		Description: "many empty lines in config",
		Config: fmt.Sprintf("\n%s\n\n%s",
			"dest_domain=example.com suffix=js revalidate=1d\n",
			"dest_domain=example.com prefix=foo suffix=js revalidate=7d\n"),
	},
}

func negativeTestDriver(tests []config.NegativeTest, t *testing.T) {
	for _, test := range tests {
		actual := cachecfg.Parse(test.Config)
		if actual == nil || actual.Code() != test.Expected {
			t.Errorf(`
  config: "%v"
  returned error: "%v"
  error should be related to: %v`, test.Config, actual, test.Description)
		}
	}
}

func TestCacheConfig(t *testing.T) {

	// Negative Tests
	negativeTestDriver(commonNegativeTests, t)
	negativeTestDriver(primaryDestinationNegativeTests, t)
	negativeTestDriver(secondarySpecifierNegativeTests, t)
	negativeTestDriver(actionNegativeTests, t)

	// Positive Tests
	for _, test := range positiveTests {
		actual := cachecfg.Parse(test.Config)
		if actual != nil {
			t.Errorf(`
  config: "%v"
  returned error: "%v"
  error should be nil
  description: %v`, test.Config, actual, test.Description)
		}

	}
}
