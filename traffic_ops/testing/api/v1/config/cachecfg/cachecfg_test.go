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
		"too few assignments",
		"foo1=foo",
		config.NotEnoughAssignments,
	},
	{
		"empty assignment",
		"dest_domain= foo1=foo",
		config.BadAssignmentMatch,
	},
	{
		"missing equals in assignment",
		"dest_domain foo1=foo",
		config.BadAssignmentMatch,
	},
	{
		"more than one primary destination",
		"dest_domain=example dest_domain=example",
		config.ExcessLabel,
	},
	{
		"using a single secondary specifier twice",
		"dest_domain=example scheme=http scheme=https",
		config.ExcessLabel,
	},
	{
		"missing primary destination label",
		"port=80 method=get time=08:00-14:00",
		config.MissingLabel,
	},
	{
		"unknown primary field",
		"foo1=foo foo2=foo foo3=foo",
		config.InvalidLabel,
	},
	{
		"good first line, but bad second",
		fmt.Sprintf("%s\n\n%s",
			"dest_domain=example.com suffix=js revalidate=1d",
			"foo1=foo foo2=foo foo3=foo"),
		config.InvalidLabel,
	},
}

var primaryDestinationNegativeTests = []config.NegativeTest{
	{
		"bad url regex",
		"url_regex=(example.com/.* foo2=foo",
		config.InvalidRegex,
	},
	{
		"bad host regex",
		"host_regex=(example.* foo2=foo",
		config.InvalidRegex,
	},
	{
		"bad hostname",
		"dest_domain=my%20bad%20domain.com foo1=foo",
		config.InvalidHost,
	},
	{
		"bad ip",
		"dest_ip=bad_ip foo1=foo",
		config.InvalidIP,
	},
}

var secondarySpecifierNegativeTests = []config.NegativeTest{
	{
		"bad port",
		"dest_domain=example port=90009000",
		config.InvalidPort,
	},
	{
		"bad scheme",
		"dest_domain=example scheme=httpz",
		config.InvalidHTTPScheme,
	},
	{
		"bad method",
		"dest_domain=example method=xxx",
		config.UnknownMethod,
	},
	{
		"bad time range",
		"dest_domain=example time=16:00",
		config.InvalidTimeRange24Hr,
	},
	{
		"bad src ip",
		"dest_domain=example src_ip=bad_ip",
		config.InvalidIP,
	},
	{
		"bad boolean value",
		"dest_domain=example internal=xxx",
		config.InvalidBool,
	},
}

var actionNegativeTests = []config.NegativeTest{
	{
		"bad action value",
		"dest_domain=example action=xxx",
		config.InvalidAction,
	},
	{
		"bad cache-responses-to-cookies",
		"dest_domain=example cache-responses-to-cookies=42",
		config.InvalidCacheCookieResponse,
	},
	{
		"bad time format",
		"dest_domain=example pin-in-cache=xxx",
		config.InvalidTimeFormatDHMS,
	},
	{
		"missing action label",
		"dest_domain=example scheme=http",
		config.MissingLabel,
	},
}

var positiveTests = []config.PositiveTest{
	{
		"empty config",
		"",
	},
	{
		"empty multiline config",
		"\n",
	},
	{
		"empty config with whitespace",
		"\t ",
	},
	{
		"normal config returned from traffic ops",
		fmt.Sprintf("%s\n%s",
			"# DO NOT EDIT - Generated for ATS_EDGE_TIER_CACHE by Traffic Ops on Fri Feb 15 22:01:53 UTC 2019",
			"dest_domain=origin.infra.ciab.test port=80 scheme=http action=never-cache"),
	},
	{
		"tab-delimitted config",
		"dest_domain=origin.infra.ciab.test\tport=80\tscheme=http\taction=never-cache",
	},
	{
		"multi-space delimitted config",
		"dest_domain=origin.infra.ciab.test  port=80  scheme=http  action=never-cache",
	},
	{
		"multiline config with empty line",
		fmt.Sprintf("%s\n\n%s",
			"dest_domain=example.com suffix=js revalidate=1d",
			"dest_domain=example.com prefix=foo suffix=js revalidate=7d"),
	},
	{
		"many empty lines in config",
		fmt.Sprintf("\n%s\n\n%s",
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
