package licenses

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"golang.org/x/pkgsite/internal/licenses"
)

var licensePath, licenseName, packageName string

const (
	licensePathFlag  = "licensepath"
	packageNameFlag  = "packagename"
	minimumCoverage  = 80.0
	godocSite        = "pkg.go.dev"
	documentationURL = "https://traffic-control-cdn.readthedocs.io/en/latest/development/godocs.html"
)

func TestMain(m *testing.M) {
	flag.StringVar(&licensePath, licensePathFlag, "", "The config file path")
	flag.StringVar(&packageName, packageNameFlag, "", "The name of the package whose license is being tested")
	flag.Parse()
	if licensePath == "" {
		fmt.Printf("Flag '%s' is required.\n", licensePathFlag)
		os.Exit(1)
	}
	if packageName == "" {
		dirName := path.Dir(licensePath)
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			goPath = build.Default.GOPATH
		}
		re := regexp.MustCompile("^" + goPath + "/src/(.*)")
		matches := re.FindStringSubmatch(dirName)
		if len(matches) < 2 {
			fmt.Printf("-licensepath should refer to a file within the GOPATH.")
		}
		packageName = matches[1]
	}
	os.Exit(m.Run())
}

// TestCoverage, in order to be run, needs this source file to be placed in a directory that can
// access the internal package golang.org/x/pkgsite/internal/licenses. Putting it in directory
// golang.org/x/pkgsite/internal/licenses/atc_coverage_test/ would work.
func TestCoverage(t *testing.T) {
	licenseData, err := ioutil.ReadFile(licensePath)
	if err != nil {
		t.Fatalf("Reading License file %s: %s", licensePath, err.Error())
	}
	_, cov := licenses.DetectFile(licenseData, "Apache-2.0", nil)
	coverageMessage := fmt.Sprintf("%.2f%% of the content within the %s License was recognized as valid licenses", cov.Percent, packageName)
	if cov.Percent >= minimumCoverage {
		t.Logf("%s (minimum is %.2f%%) .", coverageMessage, minimumCoverage)
		return
	}
	message := fmt.Sprintf(
		"%s, but at least %.2f%% of the content must be recognized in order to "+
			"ensure that the godocs for %s will display on %s. See %s for info on "+
			"how to safely add content to the %s License.\n",
		coverageMessage, minimumCoverage, packageName, godocSite, documentationURL, packageName)
	fmt.Printf("::error file=LICENSE::%s", message)
	t.Fail()
}
