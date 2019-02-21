package cache_config

import (
	"testing"

	. "github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
)

var commonNegativeTests = []NegativeTest{
	{
		"empty",
		"",
		NotEnoughAssignments,
	},
	{
		"short length",
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
		"using a secondary specifier twice",
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
}

//{
//	"bad regex",
//	"url_regex=(example.com/articles/popular.* foo2=foo",
//	InvalidRegex,
//},
//{
//	"bad hostname",
//	"dest_domain=my%20bad%20domain.com foo1=foo",
//	InvalidHost,
//},
//{
//	"bad ip",
//	"dest_ip=bad_ip foo1=foo",
//	InvalidIP, //host..
//},

func NegativeCacheConfigTests(t *testing.T) {

	tests := commonNegativeTests

	for _, test := range tests {
		actual := parseCacheConfig(test.Config)
		if actual == nil || actual.Code() != test.Expected {
			t.Errorf("config: \"%v\"\nexpected: %v\nactual: %v\n\n", test.Config, test.Expected, actual)
		}
	}
}

// multiline negative
// multiline positive
// note: a config with empty lines should pass..
// the whole config shouldn't be empty though
func PositiveCacheConfigTests(t *testing.T) {

}

func TestCacheConfig(t *testing.T) {
	NegativeCacheConfigTests(t)
	PositiveCacheConfigTests(t)
}

/*
	var positivesTest = []struct {
		coverage string
		config   string
	}{
		{
			"dest_domain (primary) suffix (secondary) revalidate (action)",
			fmt.Sprintf("%s\n%s",
				"dest_domain=example.com suffix=js revalidate=1d",
				"dest_domain=example.com prefix=foo suffix=js revalidate=7d"),
		},
	}

	for _, negTest := range positivesTest {
		_ = parseCacheConfig(negTest.config)
		//fmt.Printf("config: \n%v\nerror: %v\n\n", test.config, err)
	}*/
//badSecondaryField := "foo1=foo foo2=foo foo3=foo"
//badActionField :=
//noSecondary := "foo1=foo foo2=foo" // should be made as valid
//ex3 := `url_regex=example.com/game/.* pin-in-cache=1h`
