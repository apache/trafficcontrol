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

package orttest

// ORT Integration test functions

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

func runPluginVerifier(config_file string) error {
	args := []string{
		"--log-location-debug=test.log",
		config_file,
	}
	cmd := exec.Command("/opt/ort/plugin_verifier", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}
	return nil
}

func runTORequester(host string, data_req string) (string, error) {
	args := []string{
		"--traffic-ops-insecure=true",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--log-location-error=test.log",
		"--log-location-info=test.log",
		"--log-location-debug=test.log",
		"--get-data=" + data_req,
	}
	cmd := exec.Command("/opt/ort/to_requester", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return "", errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}

	// capture the last line of JSON in the 'Stdout' buffer 'out'
	output := strings.Split(strings.TrimSpace(strings.Replace(out.String(), "\r\n", "\n", -1)), "\n")
	lastLine := output[len(output)-1]

	return lastLine, nil
}

func runTOUpdater(host string, reval_status bool, update_status bool) error {
	args := []string{
		"--traffic-ops-insecure=true",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--log-location-error=test.log",
		"--log-location-info=test.log",
		"--log-location-debug=test.log",
		"--set-reval-status=" + strconv.FormatBool(reval_status),
		"--set-update-status=" + strconv.FormatBool(update_status),
	}
	cmd := exec.Command("/opt/ort/to_updater", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}

	return nil
}

func runT3cUpdate(host string, run_mode string) error {
	args := []string{
		"--traffic-ops-insecure=true",
		"--dispersion=0",
		"--login-dispersion=0",
		"--traffic-ops-timeout-milliseconds=3000",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host,
		"--log-location-error=test.log",
		"--log-location-info=test.log",
		"--log-location-debug=test.log",
		"--run-mode=" + run_mode,
	}
	cmd := exec.Command("/opt/ort/t3c", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}
	return nil
}

func setQueueUpdateStatus(host_name string, update string) error {
	args := []string{
		"--dir=/opt/trafficserver/etc/traffficserver",
		"--traffic-ops-insecure",
		"--traffic-ops-timeout-milliseconds=30000",
		"--traffic-ops-disable-proxy=true",
		"--traffic-ops-user=" + tcd.Config.TrafficOps.Users.Admin,
		"--traffic-ops-password=" + tcd.Config.TrafficOps.UserPassword,
		"--traffic-ops-url=" + tcd.Config.TrafficOps.URL,
		"--cache-host-name=" + host_name,
		"--log-location-error=stdout",
		"--log-location-info=stdout",
		"--log-location-warning=stdout",
		"--set-queue-status=" + update,
		"--set-reval-status=false",
	}
	cmd := exec.Command("/opt/ort/atstccfg", args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errors.New(err.Error() + ": " + "stdout: " + out.String() + " stderr: " + errOut.String())
	}
	return nil
}
