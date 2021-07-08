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
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/testing/ort-tests/tcdata"
	testutil "github.com/apache/trafficcontrol/cache-config/testing/ort-tests/util"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestWaitForParentsTrue(t *testing.T) {
	testName := getFuncName()
	t.Logf("------------- Starting " + testName + " ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const childCacheHostName = "atlanta-edge-03"
		const parentCacheHostName = "atlanta-mid-16"

		// do an initial badass to get all configs

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "badass", util.StrPtr("false")); err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		fileNameToRemove := filepath.Join(test_config_dir, "records.config")

		if !testutil.FileExists(fileNameToRemove) {
			t.Fatalf("expected: '%v' to exist after badass, actual: doesn't exist", fileNameToRemove)
		}

		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
		}

		// queue on both child and parent

		if err := ExecTOUpdater(childCacheHostName, false, true); err != nil {
			t.Fatalf("queue updates on child failed: %v\n", err)
		}
		if err := ExecTOUpdater(parentCacheHostName, false, true); err != nil {
			t.Fatalf("queue updates on parent failed: %v\n", err)
		}

		// verify child has parent-pending in TO status endpoint.

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if !status.ParentPending {
			t.Fatalf("expected: '%v' to have parent pending after queueing its parent %v, actual: %+v", childCacheHostName, parentCacheHostName, status)
		} else {
			t.Logf("Update Status of child right before running t3c-apply: %+v\n", status)
		}

		// syncds, and wait for parents

		output, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", util.StrPtr("true"))
		if err != nil {
			t.Fatalf("ERROR: t3c syncds failed: %v\n", err)
		}

		// should not have updated, because a parent is queued

		if !strings.Contains(output, "My parents still need an update, bailing.") {
			t.Fatalf("t3c wait-for-parents expected to wait for parents, actual: '''%v'''\n", output)
		}

		if testutil.FileExists(fileNameToRemove) {
			t.Fatalf("t3c wait-for-parents expected to wait for parents, actual: created file '%v'\n", fileNameToRemove)
		}

		// verify both child and parent are still queued

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if !status.UpdatePending {
			t.Errorf("expected: '%v' to still be queued after failed syncds run, actual: %+v", childCacheHostName, status)
		}

		if status, err := getUpdateStatus(parentCacheHostName); err != nil {
			t.Fatalf("checking '" + parentCacheHostName + "' queue status: " + err.Error())
		} else if !status.UpdatePending {
			t.Errorf("expected: '%v' to still be queued after unrelated child failed syncds run, actual: %+v", parentCacheHostName, status)
		}

		// un-queue the parent

		if err = ExecTOUpdater(parentCacheHostName, false, false); err != nil {
			t.Fatalf("queue updates on child failed: %v\n", err)
		}

		// syncds, and wait for parents, with the parent not queued

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", util.StrPtr("true")); err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		if !testutil.FileExists(fileNameToRemove) {
			t.Errorf("expected: '%v' to exist after syncds wait-for-parents and no parents queued, actual: doesn't exist", fileNameToRemove)
		}

		// verify both child and parent are not queued now

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if status.UpdatePending {
			t.Errorf("expected: '%v' to not be queued after successful syncds run, actual: %+v", childCacheHostName, status)
		}

		if status, err := getUpdateStatus(parentCacheHostName); err != nil {
			t.Fatalf("checking '" + parentCacheHostName + "' queue status: " + err.Error())
		} else if status.UpdatePending {
			t.Errorf("expected: '%v' to still be queued after unrelated child successful syncds run, actual: %+v", parentCacheHostName, status)
		}

	})
	t.Logf("------------- End of " + testName + " ---------------")
}

func TestWaitForParentsDefaultReval(t *testing.T) {
	testName := getFuncName()
	t.Logf("------------- Starting " + testName + " ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const childCacheHostName = "atlanta-edge-03"
		const parentCacheHostName = "atlanta-mid-16"

		// do an initial badass to get all configs

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "badass", util.StrPtr("false")); err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		fileNameToRemove := filepath.Join(test_config_dir, "records.config")

		if !testutil.FileExists(fileNameToRemove) {
			t.Fatalf("expected: '%v' to exist after badass, actual: doesn't exist", fileNameToRemove)
		}

		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
		}

		// queue both child and parent

		if err := ExecTOUpdater(childCacheHostName, false, true); err != nil {
			t.Fatalf("queue updates on child failed: %v\n", err)
		}
		if err := ExecTOUpdater(parentCacheHostName, false, true); err != nil {
			t.Fatalf("queue updates on parent failed: %v\n", err)
		}

		// verify child has parent-pending in TO status endpoint.

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if !status.ParentPending {
			t.Fatalf("expected: '%v' to have parent pending after queueing its parent %v, actual: %+v", childCacheHostName, parentCacheHostName, status)
		} else if !status.UseRevalPending {
			t.Fatalf("expected: Traffic Ops UseRevalPending must be true for this test, actual: false.")
		} else {
			t.Logf("Update Status of child right before running t3c-apply: %+v\n", status)
		}

		// syncds, and don't pass wait-for-parents, which should default to 'reval', which should *not* wait for parents since TO UseRevalPending is true.

		if output, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", nil); err != nil {
			t.Fatalf("ERROR: t3c syncds failed: error '''%v''' output '''%v'''\n", err, output)
		}

		if !testutil.FileExists(fileNameToRemove) {
			t.Errorf("expected: '%v' to exist after syncds wait-for-parents=default=reval and parents queued, actual: doesn't exist", fileNameToRemove)
		}

		// verify child is not queued now

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if status.UpdatePending {
			t.Errorf("expected: '%v' to not be queued after successful syncds run, actual: %+v", childCacheHostName, status)
		}

	})
	t.Logf("------------- End of " + testName + " ---------------")
}

func TestWaitForParentsFalse(t *testing.T) {
	testName := getFuncName()
	t.Logf("------------- Starting " + testName + " ---------------")
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters, tcdata.Statuses,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const childCacheHostName = "atlanta-edge-03"
		const parentCacheHostName = "atlanta-mid-16"

		// do an initial badass to get all configs

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "badass", util.StrPtr("false")); err != nil {
			t.Fatalf("ERROR: t3c badass failed: %v\n", err)
		}

		fileNameToRemove := filepath.Join(test_config_dir, "records.config")

		if !testutil.FileExists(fileNameToRemove) {
			t.Fatalf("expected: '%v' to exist after badass, actual: doesn't exist", fileNameToRemove)
		}

		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '" + fileNameToRemove + "': " + err.Error())
		}

		// queue both child and parent

		if err := ExecTOUpdater(childCacheHostName, false, true); err != nil {
			t.Fatalf("queue updates on child failed: %v\n", err)
		}
		if err := ExecTOUpdater(parentCacheHostName, false, true); err != nil {
			t.Fatalf("queue updates on parent failed: %v\n", err)
		}

		// delete use_reval_pending parameter, because for a syncds run, wait-for-parents 'false' behaves like 'reval' if it exists,
		// so we want to make sure it doesn't, to make sure wait-for-parents isn't actually executing 'reval'.

		params, _, err := tcdata.TOSession.GetParametersWithHdr(nil)
		if err != nil {
			t.Fatalf("getting parameters: " + err.Error())
		}

		useRevalPendingParamID := -1
		for _, param := range params {
			if tc.ConfigFileName(param.ConfigFile) != tc.GlobalConfigFileName || tc.ParameterName(param.Name) != tc.UseRevalPendingParameterName {
				continue
			}
			useRevalPendingParamID = param.ID
			break
		}
		if useRevalPendingParamID != -1 {
			if _, _, err := tcdata.TOSession.DeleteParameterByID(useRevalPendingParamID); err != nil {
				t.Fatalf("deleting useReval param: queue status: " + err.Error())
			}
		} else {
			t.Fatalf("expected '%v' '%v' Param to exist in test data, actually: missing", tc.GlobalConfigFileName, tc.UseRevalPendingParameterName)
		}

		// verify child has parent-pending in TO status endpoint.

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if !status.ParentPending {
			t.Fatalf("expected: '%v' to have parent pending after queueing its parent %v, actual: %+v", childCacheHostName, parentCacheHostName, status)
		} else if status.UseRevalPending {
			t.Fatalf("expected: Traffic Ops UseRevalPending to be false after deleting Parameter, actual: true")
		} else {
			t.Logf("Update Status of child right before running t3c-apply: %+v\n", status)
		}

		// syncds, and pass wait-for-parents=false

		if output, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", util.StrPtr("false")); err != nil {
			t.Fatalf("ERROR: t3c syncds failed: error '''%v''' output '''%v'''\n", err, output)
		}

		if !testutil.FileExists(fileNameToRemove) {
			t.Errorf("expected: '%v' to exist after syncds wait-for-parents=false and parents queued, actual: doesn't exist", fileNameToRemove)
		}

		// verify child is not queued now

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '" + childCacheHostName + "' queue status: " + err.Error())
		} else if status.UpdatePending {
			t.Errorf("expected: '%v' to not be queued after successful syncds run, actual: %+v", childCacheHostName, status)
		}

	})
	t.Logf("------------- End of " + testName + " ---------------")
}

func getUpdateStatus(hostName string) (tc.ServerUpdateStatus, error) {
	st := tc.ServerUpdateStatus{}
	if output, err := runRequest(hostName, "update-status"); err != nil {
		return tc.ServerUpdateStatus{}, errors.New("t3c-request failed: " + err.Error())
	} else if err = json.Unmarshal([]byte(output), &st); err != nil {
		return tc.ServerUpdateStatus{}, errors.New("unmarshalling t3c-request json output: " + err.Error())
	}
	return st, nil
}

func t3cUpdateWaitForParents(host string, runMode string, waitForParents *string) (string, error) {
	args := []string{
		"apply",
		"--traffic-ops-insecure=true",
		"--dispersion=0",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"-vv",
		"--run-mode=" + runMode,
		"--git=no",
		"--dns-local-bind",
	}
	if waitForParents != nil {
		args = append(args, "--wait-for-parents="+*waitForParents)
	}
	cmd := exec.Command("t3c", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return "", errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}
	return out.String() + "\n" + errOut.String(), nil
}

// getFuncName() returns the function name of the calling function.
func getFuncName() string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
