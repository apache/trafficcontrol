package cache_config

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
		"more than one primary destination",
		"dest_domain=example dest_domain=example",
		ExcessLabel,
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

var secondarySpecifierNegativeTests = []NegativeTest{
	{
		"bad port",
		"dest_domain=example port=90009000",
		InvalidPort,
	},
	{
		"bad scheme",
		"dest_domain=example scheme=httpz",
		InvalidHTTPScheme,
	},
	{
		"bad method",
		"dest_domain=example method=xxx",
		InvalidMethod,
	},
	{
		"bad time range",
		"dest_domain=example time=16:00",
		InvalidTimeRange24Hr,
	},
	{
		"bad src ip",
		"dest_domain=example src_ip=bad_ip",
		InvalidIP,
	},
	{
		"bad boolean value",
		"dest_domain=example internal=xxx",
		InvalidBool,
	},
}

var actionNegativeTests = []NegativeTest{
	{
		"bad action value",
		"dest_domain=example action=xxx",
		InvalidAction,
	},
	{
		"bad cache-responses-to-cookies",
		"dest_domain=example cache-responses-to-cookies=42",
		InvalidCacheCookieResponse,
	},
	{
		"bad time format",
		"dest_domain=example pin-in-cache=xxx",
		InvalidTimeFormatDHMS,
	},
	{
		"missing action label",
		"dest_domain=example scheme=http",
		MissingLabel,
	},
}

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
		"dest_domain=origin.infra.ciab.test port=80 scheme=http action=never-cache",
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

func negativeTestDriver(tests []NegativeTest, t *testing.T) {
	for _, test := range tests {
		actual := parseCacheConfig(test.Config)
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
		actual := parseCacheConfig(test.Config)
		if actual != nil {
			t.Errorf(`
  config: "%v"
  returned error: "%v"
  error should be nil
  description: %v`, test.Config, actual, test.Description)
		}

	}
}
