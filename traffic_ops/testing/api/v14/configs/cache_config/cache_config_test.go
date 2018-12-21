package cache_config

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// Note: the structs here are common to all config tests
// I could access them like config.NegativeTest

// Should mappings from error code to description replace the need for issue?
// No. An 'issue' can describe the error with more granularity.
type NegativeConfigTest struct {
	issue    string
	config   string
	expected int
}

// coverage should allow to to verify we could correctly make a config without an error
type PositiveConfigTest struct {
	coverage string // can be map corresponding "dest_host" and integer? sounds like an enum
	config   string
}

// does not verify that labels are in the correct position
// counts the number of occurances for each label
// this should be passed to a couple of helper functions to check validity
func getLabelMap(config string) map[string]int {

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

//func TestHelloWorld(t *testing.T) {}

// *testing.T
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

// PositiveNegative // multiline negative (only need on secondary line, but might as well include it on first line)
// PositivePositive // multiline positive
func TestPositiveCacheConfig(t *testing.T) {
	TestPrimaryPositive(t)
	TestSecondaryPositive(t)
	TestActionPositive(t)
}

func TestCacheConfig(t *testing.T) {

	// verbose prints error messages even if the test passes
	// I need to figure out how to want to fail a test in the TO environment
	// to be clear, these are unit tests
	verbose := false

	// note: a config with empty lines should pass..
	// the whole config shouldn't be empty though
	var tests = []NegativeConfigTest{
		{
			"empty",
			"",
			NOT_ENOUGH_ASSIGNMENTS,
		},
		{
			"short length",
			"foo1=foo",
			NOT_ENOUGH_ASSIGNMENTS,
		},
		{
			"empty assignment",
			"dest_domain= foo1=foo",
			BAD_ASSIGNMENT_MATCH,
		},
		{
			"missing equals in assignment",
			"dest_domain foo1=foo",
			BAD_ASSIGNMENT_MATCH,
		},

		// primary destination

		{
			"unknown primary field",
			"foo1=foo foo2=foo foo3=foo",
			INVALID_DESTINATION_LABEL,
		},
		{
			"bad regex",
			"url_regex=(example.com/articles/popular.* foo2=foo",
			INVALID_URL_REGEX,
		},
		{
			"bad hostname",
			"dest_domain=my%20bad%20domain.com foo1=foo",
			INVALID_HOST,
		},
		{
			"bad ip",
			"dest_ip=bad_ip foo1=foo",
			INVALID_IP, //host..
		},
	}

	for _, test := range tests {
		actual := parseCacheConfig(test.config)
		if verbose || actual.Code() != test.expected {
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
