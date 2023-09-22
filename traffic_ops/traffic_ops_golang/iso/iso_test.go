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

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	libtc "github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
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
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
				if got, expectedPrefix := gotResp.Header().Get(rfc.ContentDisposition), "attachment; filename="; !strings.HasPrefix(got, expectedPrefix) {
					t.Errorf("header %q = %q; expected prefix: %q", rfc.ContentDisposition, got, expectedPrefix)
				} else {
					t.Logf("header %q = %q", rfc.ContentType, got)
				}

				// Validate Content-Type header
				if got, expected := gotResp.Header().Get(rfc.ContentType), rfc.ApplicationOctetStream; got != expected {
					t.Errorf("header %q = %q; expected: %q", rfc.ContentType, got, expected)
				} else {
					t.Logf("header %q = %q", rfc.ContentType, got)
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
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
				if got, expected := gotResp.Header().Get(rfc.ContentType), rfc.ApplicationJSON; got != expected {
					t.Errorf("header %q = %q; expected: %q", rfc.ContentType, got, expected)
				} else {
					t.Logf("header %q = %q", rfc.ContentType, got)
				}

				// Validate that the response body is a JSON-encoded tc.Alerts object.
				var expectedResp libtc.Alerts
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
			// Clean up temp dir
			defer os.RemoveAll(tmpDir)
			// Create expected directory structure
			if err = os.MkdirAll(filepath.Join(tmpDir, tc.input.OSVersionDir, ksCfgDir), 0777); err != nil {
				t.Fatal(err)
			}

			// Create osversions.json file, which is read by validateOSDir method
			osVersions, err := json.Marshal(libtc.OSVersionsResponse{"test": tc.input.OSVersionDir})
			if err != nil {
				t.Fatal(err)
			}
			if err = ioutil.WriteFile(filepath.Join(tmpDir, cfgFilename), osVersions, 0600); err != nil {
				t.Fatal(err)
			}

			// Setup mock DB row such that the kickstarterDir and getOSVersions
			// functions will use tmpDir instead of the default /var/www/files
			// directory.
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
			// Mock 2 SELECT responses, 1 for the getOSVersions function
			// and another for the kickstarterDir function. Both expect
			// the same result, i.e. the temp directory.
			rows := sqlmock.NewRows(cols).AddRow(tmpDir)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			rows = sqlmock.NewRows(cols).AddRow(tmpDir)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)

			tx, err := db.BeginTxx(dbCtx, nil)
			if err != nil {
				t.Fatalf("BeginTxx() err: %v", err)
			}

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
			isos(w, req, tx, &user, tc.input)

			// pass or fail the test by inspecting the response recorder
			tc.validator(t, w)
		})
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
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
			"valid with CIDR prefix on IPv6 address",
			true,
			isoRequest{
				DHCP:          boolStr{true, false},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "2001:DB8::1/63",
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "testquest",
			},
		},
		{
			"valid with valid IPv6 address",
			true,
			isoRequest{
				DHCP:          boolStr{true, false},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "2001:DB8::1",
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "testquest",
			},
		},
		{
			"invalid with valid IPv4 address as an IPv6",
			false,
			isoRequest{
				DHCP:          boolStr{true, false},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "172.20.0.5",
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "testquest",
			},
		},
		{
			"invalid with invalid IPv6",
			false,
			isoRequest{
				DHCP:          boolStr{true, false},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "nonsense",
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "testquest",
			},
		},
		{
			"valid with dhcp true",
			true,
			isoRequest{
				DHCP:          boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    "",
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
				OSVersionDir:  "",
				HostName:      "",
				DomainName:    "",
				IPAddr:        net.IP{},
				InterfaceMTU:  0,
				InterfaceName: "",
				IP6Address:    "",
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
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			errs := tc.input.Validate(nil)
			if tc.expectedValidate != (errs == nil) {
				t.Fatalf("isoRequest.validate() = %v; expected errors = %v", errs, !tc.expectedValidate)
			}
			t.Logf("isoRequest.validate() = %+v", errs)
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

func TestISORequest_validateOSDir(t *testing.T) {
	const (
		validDir1  = "VALID-OS-DIR"
		validDir2  = "ANOTHER-VALID-OS-DIR"
		invalidDir = "INVALID-OS-DIR"
	)

	cases := []struct {
		input    isoRequest
		expected bool
	}{
		{
			isoRequest{OSVersionDir: validDir1},
			true,
		},
		{
			isoRequest{OSVersionDir: validDir2},
			true,
		},
		{
			isoRequest{OSVersionDir: invalidDir},
			false,
		},
	}

	tmpDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("error creating tempdir: %v", err)
	}
	// Clean up temp dir
	defer os.RemoveAll(tmpDir)

	// Create osversions.json file, which is read by validateOSDir method, with
	// valid directories
	osVersions, err := json.Marshal(libtc.OSVersionsResponse{
		"valid-1": validDir1,
		"valid-2": validDir2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = ioutil.WriteFile(filepath.Join(tmpDir, cfgFilename), osVersions, 0777); err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input.OSVersionDir, func(t *testing.T) {
			// Setup mock DB row such that the getOSVersions function will use tmpDir
			// instead of the default /var/www/files directory.
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
			rows := sqlmock.NewRows(cols).AddRow(tmpDir)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)

			tx, err := db.BeginTxx(dbCtx, nil)
			if err != nil {
				t.Fatalf("BeginTxx() err: %v", err)
			}

			// END setup mock DB row

			got, err := tc.input.validateOSDir(tx)
			if err != nil {
				t.Fatalf("unexpected error calling validateOSDir: %v", err)
			}

			if got != tc.expected {
				t.Fatalf("validateOSDir() got %v; expected %v", got, tc.expected)
			}
			t.Logf("validateOSDir() got %v", got)
		})
	}

}
