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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	testutil "github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/util"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

const childCacheHostName = DefaultCacheHostName

func TestWaitForParentsTrue(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const parentCacheHostName = "atlanta-mid-16"

		// do an initial badass to get all configs

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "badass", util.StrPtr("false")); err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}

		fileNameToRemove := filepath.Join(TestConfigDir, "records.config")

		if !testutil.FileExists(fileNameToRemove) {
			t.Fatalf("expected: '%s' to exist after badass, actual: doesn't exist", fileNameToRemove)
		}

		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
		}

		// queue on both child and parent

		err := tcd.QueueUpdatesForServer(childCacheHostName, true)
		if err != nil {
			t.Fatalf("queue updates on child failed: %v", err)
		}
		err = tcd.QueueUpdatesForServer(parentCacheHostName, true)
		if err != nil {
			t.Fatalf("queue updates on parent failed: %v", err)
		}
		// verify child has parent-pending in TO status endpoint.

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if !status.ParentPending {
			t.Fatalf("expected: '%s' to have parent pending after queueing its parent '%s', actual: %+v", childCacheHostName, parentCacheHostName, status)
		} else {
			t.Logf("Update Status of child right before running t3c-apply: %+v", status)
		}

		// syncds, and wait for parents

		output, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", util.StrPtr("true"))
		if err != nil {
			t.Fatalf("t3c syncds failed: %v", err)
		}

		// should not have updated, because a parent is queued

		if !strings.Contains(output, "My parents still need an update, bailing.") {
			t.Fatalf("t3c wait-for-parents expected to wait for parents, actual: %s", output)
		}

		if testutil.FileExists(fileNameToRemove) {
			t.Fatalf("t3c wait-for-parents expected to wait for parents, actual: created file '%s'", fileNameToRemove)
		}

		// verify both child and parent are still queued

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if !status.UpdatePending {
			t.Errorf("expected: '%s' to still be queued after failed syncds run, actual: %+v", childCacheHostName, status)
		}

		parentStatus, err := getUpdateStatus(parentCacheHostName)
		if err != nil {
			t.Fatalf("checking '%s' queue status: %v", parentCacheHostName, err)
		} else if !parentStatus.UpdatePending {
			t.Errorf("expected: '%s' to still be queued after unrelated child failed syncds run, actual: %+v", parentCacheHostName, parentStatus)
		}

		// un-queue the parent

		if err = ExecTOUpdater(parentCacheHostName, parentStatus.ConfigUpdateTime, parentStatus.RevalidateUpdateTime, util.BoolPtr(false), util.BoolPtr(false)); err != nil {
			t.Fatalf("queue updates on child failed: %v", err)
		}

		// syncds, and wait for parents, with the parent not queued

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", util.StrPtr("true")); err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}

		if !testutil.FileExists(fileNameToRemove) {
			t.Errorf("expected: '%s' to exist after syncds wait-for-parents and no parents queued, actual: doesn't exist", fileNameToRemove)
		}

		// verify both child and parent are not queued now

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if status.UpdatePending {
			t.Errorf("expected: '%s' to not be queued after successful syncds run, actual: %+v", childCacheHostName, status)
		}

		if status, err := getUpdateStatus(parentCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", parentCacheHostName, err)
		} else if status.UpdatePending {
			t.Errorf("expected: '%s' to not be queued after successful syncds wait-for-parents run, actual: %+v", parentCacheHostName, status)
		}

	})
}

func TestWaitForParentsDefaultReval(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const parentCacheHostName = "atlanta-mid-16"

		// do an initial badass to get all configs

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "badass", util.StrPtr("false")); err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}

		fileNameToRemove := filepath.Join(TestConfigDir, "records.config")

		if !testutil.FileExists(fileNameToRemove) {
			t.Fatalf("expected: '%s' to exist after badass, actual: doesn't exist", fileNameToRemove)
		}

		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
		}

		// queue both child and parent

		err := tcd.QueueUpdatesForServer(childCacheHostName, true)
		if err != nil {
			t.Fatalf("queue updates on child failed: %v", err)
		}
		err = tcd.QueueUpdatesForServer(parentCacheHostName, true)
		if err != nil {
			t.Fatalf("queue updates on parent failed: %v", err)
		}
		// verify child has parent-pending in TO status endpoint.

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if !status.ParentPending {
			t.Fatalf("expected: '%s' to have parent pending after queueing its parent '%s', actual: %+v", childCacheHostName, parentCacheHostName, status)
		} else if !status.UseRevalPending {
			t.Fatalf("expected: Traffic Ops UseRevalPending must be true for this test, actual: false.")
		} else {
			t.Logf("Update Status of child right before running t3c-apply: %+v", status)
		}

		// syncds, and don't pass wait-for-parents, which should default to 'reval', which should *not* wait for parents since TO UseRevalPending is true.

		if output, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", nil); err != nil {
			t.Fatalf("t3c syncds failed: error '''%v''' output '''%v'''", err, output)
		}

		if !testutil.FileExists(fileNameToRemove) {
			t.Errorf("expected: '%s' to exist after syncds wait-for-parents=default=reval and parents queued, actual: doesn't exist", fileNameToRemove)
		}

		// verify child is not queued now

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if status.UpdatePending {
			t.Errorf("expected: '%s' to not be queued after successful syncds run, actual: %+v", childCacheHostName, status)
		}

	})
}

func TestWaitForParentsFalse(t *testing.T) {
	tcd.WithObjs(t, []tcdata.TCObj{
		tcdata.CDNs, tcdata.Types, tcdata.Tenants, tcdata.Parameters,
		tcdata.Profiles, tcdata.ProfileParameters,
		tcdata.Divisions, tcdata.Regions, tcdata.PhysLocations,
		tcdata.CacheGroups, tcdata.Servers, tcdata.Topologies,
		tcdata.DeliveryServices}, func() {

		const parentCacheHostName = "atlanta-mid-16"

		// do an initial badass to get all configs

		if _, err := t3cUpdateWaitForParents(childCacheHostName, "badass", util.StrPtr("false")); err != nil {
			t.Fatalf("t3c badass failed: %v", err)
		}

		fileNameToRemove := filepath.Join(TestConfigDir, "records.config")

		if !testutil.FileExists(fileNameToRemove) {
			t.Fatalf("expected: '%s' to exist after badass, actual: doesn't exist", fileNameToRemove)
		}

		if err := os.Remove(fileNameToRemove); err != nil {
			t.Fatalf("failed to remove file '%s': %v", fileNameToRemove, err)
		}

		// queue both child and parent

		err := tcd.QueueUpdatesForServer(childCacheHostName, true)
		if err != nil {
			t.Fatalf("queue updates on child failed: %v", err)
		}
		err = tcd.QueueUpdatesForServer(parentCacheHostName, true)
		if err != nil {
			t.Fatalf("queue updates on parent failed: %v", err)
		}

		// delete use_reval_pending parameter, because for a syncds run, wait-for-parents 'false' behaves like 'reval' if it exists,
		// so we want to make sure it doesn't, to make sure wait-for-parents isn't actually executing 'reval'.

		params, _, err := tcdata.TOSession.GetParameters(toclient.RequestOptions{})
		if err != nil {
			t.Fatalf("getting parameters: %v", err)
		}

		useRevalPendingParamID := -1
		for _, param := range params.Response {
			if tc.ConfigFileName(param.ConfigFile) != tc.GlobalConfigFileName || tc.ParameterName(param.Name) != tc.UseRevalPendingParameterName {
				continue
			}
			useRevalPendingParamID = param.ID
			break
		}
		if useRevalPendingParamID != -1 {
			if _, _, err := tcdata.TOSession.DeleteParameter(useRevalPendingParamID, toclient.RequestOptions{}); err != nil {
				t.Fatalf("deleting useReval param: queue status: %v", err)
			}
		} else {
			t.Fatalf("expected '%s' '%s' Param to exist in test data, actually: missing", tc.GlobalConfigFileName, tc.UseRevalPendingParameterName)
		}

		// verify child has parent-pending in TO status endpoint.

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if !status.ParentPending {
			t.Fatalf("expected: '%s' to have parent pending after queueing its parent '%s', actual: %+v", childCacheHostName, parentCacheHostName, status)
		} else if status.UseRevalPending {
			t.Fatal("expected: Traffic Ops UseRevalPending to be false after deleting Parameter, actual: true")
		} else {
			t.Logf("Update Status of child right before running t3c-apply: %+v", status)
		}

		// syncds, and pass wait-for-parents=false

		if output, err := t3cUpdateWaitForParents(childCacheHostName, "syncds", util.StrPtr("false")); err != nil {
			t.Fatalf("t3c syncds failed: error '''%v''' output '''%v'''", err, output)
		}

		if !testutil.FileExists(fileNameToRemove) {
			t.Errorf("expected: '%s' to exist after syncds wait-for-parents=false and parents queued, actual: doesn't exist", fileNameToRemove)
		}

		// verify child is not queued now

		if status, err := getUpdateStatus(childCacheHostName); err != nil {
			t.Fatalf("checking '%s' queue status: %v", childCacheHostName, err)
		} else if status.UpdatePending {
			t.Errorf("expected: '%s' to not be queued after successful syncds run, actual: %+v", childCacheHostName, status)
		}

	})
}

func getUpdateStatus(hostName string) (atscfg.ServerUpdateStatus, error) {
	st := atscfg.ServerUpdateStatus{}
	if output, err := runRequest(hostName, CMDUpdateStatus); err != nil {
		return atscfg.ServerUpdateStatus{}, errors.New("t3c-request failed: " + err.Error())
	} else if err = json.Unmarshal([]byte(output), &st); err != nil {
		return atscfg.ServerUpdateStatus{}, errors.New("unmarshalling t3c-request json output: " + err.Error())
	}
	return st, nil
}

func t3cUpdateWaitForParents(host string, runMode string, waitForParents *string) (string, error) {
	args := []string{
		"apply",
		"--no-confirm-service-action",
		"--traffic-ops-insecure=true",
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
