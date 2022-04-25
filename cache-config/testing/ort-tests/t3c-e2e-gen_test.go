package orttest

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

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/cache-config/t3c-apply/util"
	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
)

func TestT3cCertGen(t *testing.T) {
	t.Log("------------- Starting TestT3cCertGen tests ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {
		doTestT3cCertGen(t)
	})
	t.Log("------------- End of TestT3cCertGen tests ---------------")
}

func doTestT3cCertGen(t *testing.T) {
	// verifies that when the CA is changed,
	// 1. the old CA is copied from ca.new to ca.old
	// 2. the new CA is copied from ca to ca.new
	// 3. the updated ca.new and ca.old are concatenated and copied to etc/trafficserver

	ensureDirs := []string{
		config.VarDir,
		config.EtcDir,
		"/opt/trafficserver/etc/trafficserver/ssl",
	}
	for _, dir := range ensureDirs {
		if err := os.MkdirAll(dir, 0644); err != nil {
			t.Fatalf("creating required dir '" + dir + "' expected no error, actual: " + err.Error())
		}
	}

	const cacheHostName = `atlanta-edge-03`

	// TODO make configurable? const?
	atsCAPath := "/opt/trafficserver/etc/trafficserver/ssl/e2e-ssl-ca.cert"

	// write the very first cert and key

	if err := ioutil.WriteFile(config.DefaultE2ESSLCACertPath, []byte("first-cert\n"), 0600); err != nil {
		t.Fatalf("writing first ca file '" + config.DefaultE2ESSLCACertPath + "' expected no error, actual: " + err.Error())
	}
	if err := ioutil.WriteFile(config.DefaultE2ESSLCAKeyPath, []byte("first-key"), 0600); err != nil {
		t.Fatalf("writing first ca key file '" + config.DefaultE2ESSLCAKeyPath + "' expected no error, actual: " + err.Error())
	}

	t.Logf("calling t3c-apply")
	if stdOut, stdErr, exitCode := t3cApplyE2ESSCA(cacheHostName, "badass", "", ""); exitCode != 0 {
		t.Fatalf("t3c-apply failed: code '%v' output '%v' stderr '%v'\n", exitCode, stdOut, stdErr)
	}

	{
		oldFile, err := ioutil.ReadFile(util.E2ESSLOldCAPath)
		if err != nil {
			t.Fatalf("after running t3c with CA, expected old CA file to be created and read without error, actual err: " + err.Error())
		}
		if strings.TrimSpace(string(oldFile)) != "" {
			t.Fatalf("after running t3c for the first time with a CA, expected old CA file '" + util.E2ESSLOldCAPath + "' to be created as empty, actual '" + string(oldFile))
		}
	}
	{
		newFile, err := ioutil.ReadFile(util.E2ESSLNewCAPath)
		if err != nil {
			t.Fatalf("after running t3c with CA, expected new CA file '" + util.E2ESSLNewCAPath + "' to be created and read without error, actual err: " + err.Error())
		}
		if strings.TrimSpace(string(newFile)) != "first-cert" {
			t.Fatalf("after running t3c for the first time with a CA, expected new CA file '" + util.E2ESSLNewCAPath + "' to be created as a copy of the new CA 'first-cert', actual '" + string(newFile) + "'")
		}
	}
	{
		atsCAFile, err := ioutil.ReadFile(atsCAPath)
		if err != nil {
			// lsOut, stdErr, exitCode := t3cutil.Do("ls", "-lah", "/opt/trafficserver/etc/trafficserver/ssl/")
			t.Fatalf("after running t3c with CA, expected ATS CA file '" + atsCAPath + "' to be created and readable, actual: " + err.Error())
		}
		if strings.TrimSpace(string(atsCAFile)) != "first-cert" {
			t.Fatalf("after running t3c for the first time with a CA, expected ATS CA file '" + atsCAPath + "' to be created as CA '" + "first-cert" + "', actual '" + string(atsCAFile))
		}
	}

	// write a new cert (but don't change the key)

	if err := ioutil.WriteFile(config.DefaultE2ESSLCACertPath, []byte("second-cert\n"), 0600); err != nil {
		t.Fatalf("writing first ca file '" + config.DefaultE2ESSLCACertPath + "' expected no error, actual: " + err.Error())
	}

	t.Logf("calling t3c-apply")
	if stdOut, stdErr, exitCode := t3cApplyE2ESSCA(cacheHostName, "badass", "", ""); exitCode != 0 {
		t.Fatalf("t3c-apply failed: code '%v' output '%v' stderr '%v'\n", exitCode, stdOut, stdErr)
	}

	{
		oldFile, err := ioutil.ReadFile(util.E2ESSLOldCAPath)
		if err != nil {
			t.Fatalf("after running t3c with CA, expected old CA file to be created and read without error, actual err: " + err.Error())
		}
		if strings.TrimSpace(string(oldFile)) != "first-cert" {
			t.Fatalf("after running t3c for the second time with a CA, expected old CA file '" + util.E2ESSLOldCAPath + "' to be set to the previous 'first-cert', actual '" + string(oldFile))
		}
	}
	{
		newFile, err := ioutil.ReadFile(util.E2ESSLNewCAPath)
		if err != nil {
			t.Fatalf("after running t3c with CA, expected new CA file '" + util.E2ESSLNewCAPath + "' to be created and read without error, actual err: " + err.Error())
		}
		if strings.TrimSpace(string(newFile)) != "second-cert" {
			t.Fatalf("after running t3c for the second time with a CA, expected new CA file '" + util.E2ESSLNewCAPath + "' to be created as a copy of the new CA 'second-cert', actual '" + string(newFile) + "'")
		}
	}
	{
		atsCAFile, err := ioutil.ReadFile(atsCAPath)
		if err != nil {
			// lsOut, stdErr, exitCode := t3cutil.Do("ls", "-lah", "/opt/trafficserver/etc/trafficserver/ssl/")
			t.Fatalf("after running t3c with CA, expected ATS CA file '" + atsCAPath + "' to be created and readable, actual: " + err.Error())
		}
		if strings.TrimSpace(string(atsCAFile)) != "second-cert\nfirst-cert" {
			t.Fatalf("after running t3c for the first time with a CA, expected ATS CA file '" + atsCAPath + "' to be set to the concatenated new and old '" + "second-cert\nfirst-cert" + "', actual '" + string(atsCAFile))
		} else {
			t.Log("TestT3cCertGen verified concatenated cert was written")
		}
	}

}

func t3cApplyE2ESSCA(host string, runMode string, e2eCAPath string, e2eCAKeyPath string) (string, string, int) {
	args := []string{
		"apply",
		"--traffic-ops-insecure=true",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
		"--omit-via-string-release=true",
		"--git=no",
		"--run-mode=" + runMode,
	}
	if e2eCAPath != "" {
		args = append(args, "--e2e-ca-cert="+e2eCAPath)
	}
	if e2eCAKeyPath != "" {
		args = append(args, "--e2e-ca-key="+e2eCAKeyPath)
	}

	stdOut, stdErr, exitCode := t3cutil.Do("t3c", args...) // should be no stderr, we told it to log to stdout
	return string(stdOut), string(stdErr), exitCode
}
