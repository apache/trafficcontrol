package cache_config

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
)

// does not verify that labels are in the correct position
// counts the number of occurances for each label
// this should be passed to a couple of helper functions to check validity
//
// This is used to verify that in the positive test we get to test validity
// on at least one of each label.
func getLabelMap(config string) map[string]int {

	// Rewrite this.
	// Why?
	// I don't need the count.
	// What I need is this:
	// config -> array of labels
	// This is to check that each label
	// is used ONCE at the very least.
	// I do not need to check here for
	// the validity ie. if the illegal
	// action of having two ports for
	// a single config was done.

	var labels = map[string]int{
		"dest_domain": 0,
		"dest_host":   0,
		"dest_ip":     0,
		"host_regex":  0,
		"url_regex":   0,
		"port":        0,
		"scheme":      0,
		"prefix":      0,
		"suffix":      0,
		"method":      0,
		"time":        0,
		"src_ip":      0,
		"internal":    0,
	}

	assignment := regexp.MustCompile(`([a-z_\d]+)=(\S+)`)
	for _, elem := range strings.Split(config, " ") {
		if match := assignment.FindStringSubmatch(elem); match != nil {
			labels[match[1]]++
		}
		return nil
	}

	return labels
}

// This seems to make sense to do, but I don't want to now because it isn't clear
// how to split the tests up yet

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
	var tests = []config.NegativeTest{
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
			InvalidDestinationLabel,
		},
		{
			"bad regex",
			"url_regex=(example.com/articles/popular.* foo2=foo",
			InvalidURLRegex,
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
