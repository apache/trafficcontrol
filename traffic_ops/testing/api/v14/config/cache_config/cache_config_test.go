package cache_config

import (
	"fmt"
	"testing"

	. "github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
)

func TestCommonNegative(t *testing.T) {

}

func TestPrimaryNegative(t *testing.T) {

}

func TestSecondaryNegative(t *testing.T) {

}

func TestActionNegative(t *testing.T) {

}

func TestNegativeCacheConfig(t *testing.T) {
	TestCommonNegative(t)
	TestPrimaryNegative(t)
	TestSecondaryNegative(t)
	TestActionNegative(t)
}

func TestPrimaryPositive(t *testing.T) {

}

func TestSecondaryPositive(t *testing.T) {

}

func TestActionPositive(t *testing.T) {

}

// PositiveNegative // multiline negative
// PositivePositive // multiline positive
func TestPositiveCacheConfig(t *testing.T) {
	TestPrimaryPositive(t)
	TestSecondaryPositive(t)
	TestActionPositive(t)
}

func TestCacheConfig(t *testing.T) {

	// verbose prints error messages even if the test passes
	// I need to figure out how to want to fail a test in the TO environment
	verbose := false

	// note: a config with empty lines should pass..
	// the whole config shouldn't be empty though
	var tests = []NegativeTest{
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

		// primary destination

		{
			"unknown primary field",
			"foo1=foo foo2=foo foo3=foo",
			InvalidLabel,
		},
		{
			"bad regex",
			"url_regex=(example.com/articles/popular.* foo2=foo",
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
			InvalidIP, //host..
		},
	}

	for _, test := range tests {
		actual := parseCacheConfig(test.Config)
		if verbose || actual.Code() != test.Expected {
			fmt.Println(actual)
			//fmt.Printf("config: \"%v\"\nexpected: %v\nactual: %v\n\n", test.config, test.expected, actual)
		}
	}

}

/*
	var positiveTests = []struct {
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

	for _, negTest := range positiveTests {
		_ = parseCacheConfig(negTest.config)
		//fmt.Printf("config: \n%v\nerror: %v\n\n", test.config, err)
	}*/
//badSecondaryField := "foo1=foo foo2=foo foo3=foo"
//badActionField :=
//noSecondary := "foo1=foo foo2=foo" // should be made as valid
//ex3 := `url_regex=example.com/game/.* pin-in-cache=1h`
