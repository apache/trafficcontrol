package iso

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
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestISOS(t *testing.T) {
	cases := []struct {
		name      string                                           // name of testcase
		input     isoRequest                                       // input to isos function
		cmdMod    func(in *exec.Cmd) *exec.Cmd                     // optional modifier of the cmd that will be executed
		validator func(t *testing.T, w *httptest.ResponseRecorder) // validation function which should pass/fail the test
	}{
		{
			name: "cmd success",
			input: isoRequest{
				DHCP:          boolStr{true, false},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
			cmdMod: func(in *exec.Cmd) *exec.Cmd {
				// Modify command such that it will use the mock command with
				// a successful response. It will write to STDOUT the original
				// command+args, e.g. mkisofs arg1 arg2..., which should be copied
				// to the response's body.
				return mockISOCmd(in, false, "")
			},
			validator: func(t *testing.T, gotResp *httptest.ResponseRecorder) {
				// Validate response code
				if got, expected := gotResp.Code, http.StatusOK; got != expected {
					t.Errorf("response status = %d; expected %d", got, expected)
				} else {
					t.Logf("response status = %d", got)
				}

				// Validate Content-Disposition header
				if got, expectedPrefix := gotResp.Header().Get(httpHeaderContentDisposition), "attachment; filename="; !strings.HasPrefix(got, expectedPrefix) {
					t.Errorf("header %q = %q; expected prefix: %q", httpHeaderContentDisposition, got, expectedPrefix)
				} else {
					t.Logf("header %q = %q", httpHeaderContentDisposition, got)
				}

				// Validate Content-Type header
				if got, expected := gotResp.Header().Get(httpHeaderContentType), httpHeaderContentDownload; got != expected {
					t.Errorf("header %q = %q; expected: %q", httpHeaderContentType, got, expected)
				} else {
					t.Logf("header %q = %q", httpHeaderContentType, got)
				}

				// Because of the mocked command, the data written to the response body should be
				// the command and args that were originally set to be executed. In this case, it's
				// expected to be the default mkisofs command since there's no `generate` script present.
				if got, expectedPrefix := gotResp.Body.String(), mkisofsBin; !strings.HasPrefix(got, expectedPrefix) {
					t.Errorf("got command: %q; expected command prefix of: %q", got, expectedPrefix)
				} else {
					t.Logf("got command: %q", got)
				}
			},
		},

		{
			name: "cmd failure",
			input: isoRequest{
				DHCP:          boolStr{true, false},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
			cmdMod: func(in *exec.Cmd) *exec.Cmd {
				// Modify command such that it will return an error.
				return mockISOCmd(in, true, "")
			},
			validator: func(t *testing.T, gotResp *httptest.ResponseRecorder) {
				// TODO: Validate response code. Because of how api.HandleErr works, it's
				// not possible to see the actual response code in the response recorder.

				// Validate Content-Type header, which should be JSON for this error condition.
				if got, expected := gotResp.Header().Get(httpHeaderContentType), rfc.ApplicationJSON; got != expected {
					t.Errorf("header %q = %q; expected: %q", httpHeaderContentType, got, expected)
				} else {
					t.Logf("header %q = %q", httpHeaderContentType, got)
				}

				// Validate that the response body is a JSON-encoded tc.Alerts object.
				var expectedResp tc.Alerts
				if err := json.NewDecoder(gotResp.Body).Decode(&expectedResp); err != nil {
					t.Fatalf("unable to decode body into expected JSON structure: %v", err)
				}

				t.Logf("response: %#v", expectedResp)
			},
		},
	}

	tmpDirPrefix := t.Name()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Create the OS versions & ks_scripts directories within a temporary directory.
			// The tmpDir will be used, via a Parameter entry in the database, instead of the default
			// /var/www/files directory.

			tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
			if err != nil {
				t.Fatalf("error creating tempdir: %v", err)
			}
			// Clean up temp dir + file
			defer os.RemoveAll(tmpDir)
			if err = os.MkdirAll(filepath.Join(tmpDir, tc.input.OSVersionDir, ksCfgDir), 0777); err != nil {
				t.Fatal(err)
			}

			// Setup mock DB row such that the kickstarterDir function will
			// return tmpDir instead of the default /var/www/files directory.
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer mockDB.Close()
			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			dbCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
			defer cancel()

			// Setup mock DB to return row for SELECT query on parameter table.
			mock.ExpectBegin()
			cols := []string{"value"}
			rows := sqlmock.NewRows(cols)
			rows = rows.AddRow(tmpDir)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			tx, err := db.BeginTxx(dbCtx, nil)
			if err != nil {
				t.Fatalf("BeginTxx() err: %v", err)
			}
			defer tx.Commit()

			// END setup mock DB row

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/isos", nil) // The path doesn't matter here
			if err != nil {
				t.Fatal(err)
			}
			// If not nil, add the test case's cmdMod function to the
			// request's context. The isos function will then extract
			// the function from the request and apply it to the command.
			if tc.cmdMod != nil {
				ctx := context.WithValue(req.Context(), cmdOverwriteCtxKey, tc.cmdMod)
				req = req.WithContext(ctx)
			}

			var user auth.CurrentUser
			isos(tx, &user, w, req, tc.input)

			// pass or fail the test by inspecting the response recorder
			tc.validator(t, w)
		})
	}
}

func TestWriteRespErrorAlerts(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/isos", nil) // The path doesn't matter here
	if err != nil {
		t.Fatal(err)
	}

	errMsgs := []string{
		"error 1",
		"another error",
	}

	writeRespErrorAlerts(w, req, errMsgs)

	if got, expected := w.Code, http.StatusBadRequest; got != expected {
		t.Errorf("got response code %d; expected %d", got, expected)
	}

	var gotResp tc.Alerts
	if err := json.NewDecoder(w.Body).Decode(&gotResp); err != nil {
		t.Fatalf("unable to decode response body into expected JSON structure: %v", err)
	}

	t.Logf("got response: %v", gotResp)

	if got, expected := len(gotResp.Alerts), len(errMsgs); got != expected {
		t.Fatalf("got %d error messages; expected %d", got, expected)
	}

	for i, v := range errMsgs {
		if got, expected := gotResp.Alerts[i].Level, tc.ErrorLevel.String(); got != expected {
			t.Errorf("got response with alerts[%d].Level = %s; expected %s", i, got, expected)
		}
		if got, expected := gotResp.Alerts[i].Text, v; got != expected {
			t.Errorf("got response with alerts[%d].Text = %s; expected %s", i, got, expected)
		}
	}
}

func TestISORequest_validate(t *testing.T) {
	cases := []struct {
		name             string
		expectedValidate bool
		input            isoRequest
	}{
		{
			"valid with dhcp false",
			true,
			isoRequest{
				DHCP:          boolStr{true, false},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"valid with dhcp true",
			true,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"invalid with dhcp false",
			false,
			isoRequest{
				DHCP:          boolStr{true, false},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},

		{
			"valid with mgmt addr",
			true,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "not empty",
				MgmtIPAddress: net.IP{192, 168, 0, 1},
				MgmtIPGateway: net.IP{192, 168, 0, 2},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"invalid with mgmt addr",
			false,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{192, 168, 0, 1},
				MgmtIPGateway: net.IP{192, 168, 0, 2},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},

		{
			"invalid with zero values",
			false,
			isoRequest{
				DHCP:          boolStr{false, false},
				Stream:        boolStr{false, false},
				OSVersionDir:  "",
				HostName:      "",
				DomainName:    "",
				IPAddr:        net.IP{},
				InterfaceMTU:  0,
				InterfaceName: "",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "",
				RootPass:      "",
			},
		},

		{
			"invalid with stream true",
			false,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, false},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "not empty",
				MgmtIPAddress: net.IP{192, 168, 0, 1},
				MgmtIPGateway: net.IP{192, 168, 0, 2},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"valid with stream unset",
			true,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{false, false},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "not empty",
				MgmtIPAddress: net.IP{192, 168, 0, 1},
				MgmtIPGateway: net.IP{192, 168, 0, 2},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotErrs := tc.input.validate()
			if tc.expectedValidate != (len(gotErrs) == 0) {
				t.Fatalf("isoRequest.validate() = %v; expected errors = %v", gotErrs, !tc.expectedValidate)
			}
			t.Logf("isoRequest.validate() = %+v", gotErrs)
		})
	}
}

func TestBoolStr_UnmarshalText(t *testing.T) {
	cases := []struct {
		input    string
		expected boolStr
	}{
		{
			`no`,
			boolStr{isSet: true, v: false},
		},
		{
			`No`,
			boolStr{isSet: true, v: false},
		},
		{
			`YES`,
			boolStr{isSet: true, v: true},
		},
		{
			`other`,
			boolStr{isSet: false, v: false},
		},
		{
			``,
			boolStr{isSet: false, v: false},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			var got boolStr
			if err := got.UnmarshalText([]byte(tc.input)); err != nil {
				t.Fatal(err)
			}

			if got != tc.expected {
				t.Fatalf("got %+v; expected %+v", got, tc.expected)
			}
			t.Logf("got %+v", got)
		})
	}
}
