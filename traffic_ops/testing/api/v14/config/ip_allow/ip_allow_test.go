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

package ip_allow

import (
	"fmt"
	"testing"

	. "github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
)

var commonNegativeTests = []NegativeTest{
	{
		"too few assignments",
		"foo1=foo",
		NotEnoughAssignments,
	},
	{
		"empty assignment",
		"dest_domain= foo1=foo",
		BadAssignmentMatch,
	},
	{
		"missing equals in assignment",
		"dest_domain foo1=foo",
		BadAssignmentMatch,
	},
	{
		"using a single secondary specifier twice",
		"dest_domain=example scheme=http scheme=https",
		ExcessLabel,
	},
	{
		"missing primary destination label",
		"port=80 method=get time=08:00-14:00",
		MissingLabel,
	},
	{
		"unknown primary field",
		"foo1=foo foo2=foo foo3=foo",
		InvalidLabel,
	},
	{
		"good first line, but bad second",
		fmt.Sprintf("%s\n\n%s",
			"dest_domain=example.com suffix=js revalidate=1d",
			"foo1=foo foo2=foo foo3=foo"),
		InvalidLabel,
	},
}

var primaryDestinationNegativeTests = []NegativeTest{
	{
		"bad url regex",
		"url_regex=(example.com/.* foo2=foo",
		InvalidRegex,
	},
	{
		"bad host regex",
		"host_regex=(example.* foo2=foo",
		InvalidRegex,
	},
	{
		"bad hostname",
		"dest_domain=my%20bad%20domain.com foo1=foo",
		InvalidHost,
	},
	{
		"bad ip",
		"dest_ip=bad_ip foo1=foo",
		InvalidIP,
	},
}

/*
var positiveTests = []PositiveTest{
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
}*/

var positiveTests = []PositiveTest{
	{
		"cfg1",
		"src_ip=127.0.0.1 action=ip_allow method=ALL",
	},
	{
		"cfg2",
		"src_ip=::1 action=ip_allow method=ALL",
	},
	{
		"cfg3",
		"src_ip=0.0.0.0-255.255.255.255 action=ip_deny method=PUSH|PURGE|DELETE",
	},
	{
		"cfg4",
		"src_ip=::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff action=ip_deny method=PUSH|PURGE|DELETE",
	},
	{
		"cfg5",
		"src_ip=127.0.0.1 action=ip_allow method=ALL",
	},
	{
		"cfg6",
		"src_ip=0/0 action=ip_deny method=PUSH|PURGE|DELETE",
	},
	{
		"cfg7",
		"src_ip=::/0 action=ip_deny method=PUSH|PURGE|DELETE",
	},
	{
		"cfg8",
		"src_ip=0.0.0.0-255.255.255.255 action=ip_allow",
	},
	{
		"cfg9",
		"src_ip=123.12.3.000-123.12.3.123 action=ip_allow",
	},
	{
		"cfg10",
		"src_ip=123.45.6.0-123.45.6.123 action=ip_deny",
	},
	{
		"cfg11",
		"src_ip=123.45.6.0/24 action=ip_deny",
	},
	{
		"cfg12",
		"dest_ip=0.0.0.0-255.255.255.255 action=ip_allow",
	},
	{
		"cfg13",
		"dest_ip=0/0 action=ip_allow",
	},
	{
		"cfg14",
		"dest_ip=10.0.0.0-10.0.255.255 action=ip_deny",
	},
	{
		"cfg15",
		"dest_ip=10.0.0.0/16 action=ip_deny",
	},
	{
		"cfg15",
		"dest_ip=10.0.0.0/16 action=ip_allow method=GET method=HEAD",
	},
	{
		"cfg15",
		"dest_ip=10.0.0.0/16 action=ip_allow method=GET|HEAD",
	},
}

func negativeTestDriver(tests []NegativeTest, t *testing.T) {
	for _, test := range tests {
		actual := parseConfig(test.Config)
		if actual == nil || actual.Code() != test.Expected {
			t.Errorf(`
  config: "%v"
  returned error: "%v"
  error should be related to: %v`, test.Config, actual, test.Description)
		}
	}
}

func TestIPAllowConfig(t *testing.T) {

	// Negative Tests
	//negativeTestDriver(commonNegativeTests, t)
	//negativeTestDriver(primaryDestinationNegativeTests, t)

	// Positive Tests
	for _, test := range positiveTests {
		actual := parseConfig(test.Config)
		if actual != nil {
			t.Errorf(`
  config: "%v"
  returned error: "%v"
  error should be nil
  description: %v`, test.Config, actual, test.Description)
		}

	}
}
