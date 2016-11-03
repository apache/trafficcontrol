package cache

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestAstats(t *testing.T) {
	t.Log("Running Astats Tests")

	text, err := ioutil.ReadFile("astats.json")
	if err != nil {
		t.Log(err)
	}
	aStats, err := Unmarshal(text)
	fmt.Printf("aStats ---> %v\n", aStats)
	if err != nil {
		t.Log(err)
	}
	fmt.Printf("Found %v key/val pairs in ats\n", len(aStats.Ats))
}
