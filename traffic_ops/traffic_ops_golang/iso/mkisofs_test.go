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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGenISO(t *testing.T) {
	t.Run("standard-to-file", func(t *testing.T) {
		m := newMkiso("/var/ks", "/tmp/image.iso", false)
		m.cmd = mockCmd(m.cmd)

		var out bytes.Buffer
		if err := m.do(&out); err != nil {
			t.Fatal(err)
		}

		if l := out.Len(); l > 0 {
			t.Fatalf("mkiso.do(w) wrote %d bytes to w; expected 0", l)
		}
	})

	t.Run("standard-stream", func(t *testing.T) {
		m := newMkiso("/var/ks", "/tmp/image.iso", true)
		m.cmd = mockCmd(m.cmd)

		var out bytes.Buffer
		if err := m.do(&out); err != nil {
			t.Fatal(err)
		}

		if l := out.Len(); l < 0 {
			t.Fatalf("mkiso.do(w) wrote %d bytes to w; expected > 0", l)
		}
		t.Logf("command: %q", out.String())
	})

	t.Run("custom-to-file", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "TestGenISO")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		// Create custom executable that should be used instead of mkisofs
		fd, err := os.OpenFile(filepath.Join(dir, "generate"), os.O_CREATE|os.O_EXCL, 0777)
		if err != nil {
			t.Fatal(err)
		}
		defer fd.Close()

		m := newMkiso(dir, "/tmp/image.iso", false)
		m.cmd = mockCmd(m.cmd)

		var out bytes.Buffer
		if err := m.do(&out); err != nil {
			t.Fatal(err)
		}

		if l := out.Len(); l > 0 {
			t.Fatalf("mkiso.do(w) wrote %d bytes to w; expected 0", l)
		}
	})

	t.Run("custom-stream", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "TestGenISO")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		// Create custom executable that should be used instead of mkisofs
		fd, err := os.OpenFile(filepath.Join(dir, "generate"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		if err != nil {
			t.Fatal(err)
		}
		defer fd.Close()

		if _, err = fmt.Fprint(fd, t.Name()); err != nil {
			t.Fatal(err)
		}

		m := newMkiso(dir, fd.Name(), true)
		m.cmd = mockCmd(m.cmd)

		var out bytes.Buffer
		if err := m.do(&out); err != nil {
			t.Fatal(err)
		}

		if l := out.Len(); l < 0 {
			t.Fatalf("mkiso.do(w) wrote %d bytes to w; expected > 0", l)
		}
		t.Logf("command: %q", out.String())
	})
}

// mockCmd returns a modified version of the given Cmd
// so that when run, the command actually invokes the
// TestHelperMockCmd test. See TestHelperMockCmd for
// more details on its behavior.
func mockCmd(cmd *exec.Cmd) *exec.Cmd {
	args := []string{
		"-test.run=TestHelperMockCmd",
		"--",
	}
	args = append(args, cmd.Args...)

	// os.Args[0] is the invokation of this test binary
	mocked := exec.Command(os.Args[0], args...)

	env := cmd.Env
	env = append(cmd.Env, "GO_HELPER_PROCESS=1")
	mocked.Env = env

	return mocked
}

func TestHelperMockCmd(t *testing.T) {
	if os.Getenv("GO_HELPER_PROCESS") != "1" {
		return
	}

	var respCode int
	if os.Getenv("GO_FORCE_ERROR") == "1" {
		respCode = 1
	}

	// Set args to all arguments past '--'.
	var args []string
	for i, v := range os.Args {
		if v == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	out := os.Stdout
	if respCode != 0 {
		out = os.Stderr
	}

	fmt.Fprintf(out, strings.Join(args, " "))
	os.Exit(respCode)
}

func TestKickstarterDir(t *testing.T) {
	cases := []struct {
		name           string
		parameterValue string
		osVersionsDir  string
		expected       string
	}{
		{
			"default",
			"", // No DB parameter entry
			"templeOS",
			cfgDefaultDir + "/templeOS",
		},
		{
			"param-override",
			"/var/override/dir",
			"anotherOS/dir",
			"/var/override/dir/anotherOS/dir",
		},
	}

	// Mock DB setup ...

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() err: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// <mock DB row>

			dbCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
			defer cancel()

			// Setup mock DB to return rows for SELECT query on parameter table.
			// If parameterValue is empty, no rows will be returned.
			mock.ExpectBegin()
			cols := []string{"value"}
			rows := sqlmock.NewRows(cols)
			if tc.parameterValue != "" {
				rows = rows.AddRow(tc.parameterValue)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			tx, err := db.BeginTxx(dbCtx, nil)
			if err != nil {
				t.Fatalf("BeginTxx() err: %v", err)
			}
			defer tx.Commit()

			// </mock DB row>

			got, err := kickstarterDir(tx, tc.osVersionsDir)
			if err != nil {
				t.Fatalf("kickstarterDir(tx, %q) error: %v", tc.osVersionsDir, err)
			}

			if got != tc.expected {
				t.Fatalf("kickstarterDir(tx, %q) = %q; expected %q", tc.osVersionsDir, got, tc.expected)
			}
			t.Logf("kickstarterDir(tx, %q) = %q", tc.osVersionsDir, got)
		})
	}
}

func TestWriteKSCfgs(t *testing.T) {
	tmpPrefix := t.Name()

	cases := []struct {
		name  string
		input isoRequest
	}{
		{
			"empty",
			isoRequest{},
		},
		{
			"complete",
			isoRequest{
				RootPass:      "password",
				Disk:          "sda1",
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "bond0",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "cdn.example",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
				MgmtIPAddress: net.IP{10, 10, 0, 1},
				MgmtIPGateway: net.IP{10, 10, 0, 255},
				MgmtIPNetmask: net.IP{10, 10, 255, 255},
				MgmtInterface: "eth0",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", tmpPrefix)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if err := writeKSCfgs(dir, tc.input); err != nil {
				t.Fatalf("writeKSCfgs(%q, isoRequest) failed: %v", dir, err)
			}
			t.Logf("writeKSCfgs(%q, isoRequest):", dir)

			expectedFiles := []string{
				ksCfgNetwork,
				ksCfgMgmtNetwork,
				ksCfgPassword,
				ksCfgDisk,
			}

			// Ensure that the expected files exist and are not empty.

			for _, expectedFile := range expectedFiles {
				p := filepath.Join(dir, expectedFile)

				gotB, err := ioutil.ReadFile(p)
				if err != nil {
					t.Errorf("error reading %q: %v", p, err)
					continue
				}
				got := string(gotB)

				if strings.TrimSpace(got) == "" {
					t.Errorf("%q file is empty; expected non-empty", p)
					continue
				}

				t.Logf("%s:\n%s", expectedFile, got)
			}
		})
	}
}

func TestWriteNetworkCfg(t *testing.T) {
	cases := []struct {
		name        string
		input       isoRequest
		nameservers []string
		expected    string
	}{
		{
			"empty",
			isoRequest{},
			nil,

			`
IPADDR=""
NETMASK=""
GATEWAY=""
DEVICE=""
MTU="0"
NAMESERVER=""
HOSTNAME=""
NETWORKING_IPV6="yes"
IPV6ADDR=""
IPV6_DEFAULTGW=""
DHCP="no"
`,
		},

		{
			"no domain",
			isoRequest{
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "eth0",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
			},
			[]string{"8.8.8.8", "1.1.1.1"},

			`
IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
DEVICE="eth0"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
DHCP="yes"
`,
		},

		{
			"non-bonded",
			isoRequest{
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "eth0",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "example.com",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
			},
			[]string{"8.8.8.8", "1.1.1.1"},

			`
IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
DEVICE="eth0"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server.example.com"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
DHCP="yes"
`,
		},

		{
			"bonded",
			isoRequest{
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "bond01",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
			},
			[]string{"8.8.8.8", "1.1.1.1"},

			`
IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
BOND_DEVICE="bond01"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
BONDING_OPTS="miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4"
DHCP="yes"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writeNetworkCfg(&w, tc.input, tc.nameservers); err != nil {
				t.Fatalf("writeNetworkCfg() err = %v", err)
			}
			got := w.String()
			expected := strings.TrimSpace(tc.expected)

			if got != expected {
				t.Fatalf("writeNetworkCfg() got != expected\n got:\n%s\n expected:\n%s", got, expected)
			}
			t.Logf("writeNetworkCfg():\n%s", got)
		})
	}
}

func TestWriteMgmtNetworkCfg(t *testing.T) {
	cases := []struct {
		name     string
		input    isoRequest
		expected string
	}{
		{
			"empty",
			isoRequest{},

			`
IPADDR=""
NETMASK=""
GATEWAY=""
DEVICE=""
`,
		},

		{
			"IPv4",
			isoRequest{
				MgmtIPAddress: net.IP{192, 168, 2, 3},
				MgmtIPNetmask: net.IP{255, 255, 255, 255},
				MgmtIPGateway: net.IP{192, 168, 1, 255},
				MgmtInterface: "eth0",
			},

			`
IPADDR="192.168.2.3"
NETMASK="255.255.255.255"
GATEWAY="192.168.1.255"
DEVICE="eth0"
`,
		},

		{
			"IPv6",
			isoRequest{
				MgmtIPAddress: net.ParseIP("beef::1"),
				MgmtIPNetmask: net.IP{255, 255, 255, 255},
				MgmtIPGateway: net.IP{192, 168, 1, 255},
				MgmtInterface: "eth0",
			},

			`
IPV6ADDR="beef::1"
NETMASK="255.255.255.255"
GATEWAY="192.168.1.255"
DEVICE="eth0"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writeMgmtNetworkCfg(&w, tc.input); err != nil {
				t.Fatalf("writeMgmtNetworkCfg() err = %v", err)
			}
			got := w.String()
			expected := strings.TrimSpace(tc.expected)

			if got != expected {
				t.Fatalf("writeMgmtNetworkCfg() got != expected\n got:\n%s\n expected:\n%s", got, expected)
			}
			t.Logf("writeMgmtNetworkCfg():\n%s", got)
		})
	}
}

func TestWriteDiskCfg(t *testing.T) {
	cases := []struct {
		name     string
		input    isoRequest
		expected string
	}{
		{
			"empty",
			isoRequest{},
			`
boot_drives=""
`,
		},

		{
			"non-empty",
			isoRequest{
				Disk: "sda1",
			},
			`
boot_drives="sda1"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writeDiskCfg(&w, tc.input); err != nil {
				t.Fatalf("writeDiskCfg() err = %v", err)
			}
			got := w.String()
			expected := strings.TrimSpace(tc.expected)

			if got != expected {
				t.Fatalf("writeDiskCfg() got != expected\n got:\n%s\n expected:\n%s", got, expected)
			}
			t.Logf("writeDiskCfg():\n%s", got)
		})
	}
}

func TestWritePasswordCfg(t *testing.T) {
	cases := []struct {
		name     string
		input    isoRequest
		salt     string
		expected string
	}{
		{
			"empty",
			isoRequest{},
			"salt",
			"rootpw --iscrypted $1$salt$UsdFqFVB.FsuinRDK5eE..\n",
		},

		{
			"non-empty",
			isoRequest{
				RootPass: "Traffic Ops",
			},
			"salt",
			"rootpw --iscrypted $1$salt$17HeaymOIi.65dl76MkK01\n",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writePasswordCfg(&w, tc.input, tc.salt); err != nil {
				t.Fatalf("writePasswordCfg() err = %v", err)
			}
			got := w.String()

			if got != tc.expected {
				t.Fatalf("writePasswordCfg() got != expected\n got:\n%s\n expected:\n%s", got, tc.expected)
			}
			t.Logf("writePasswordCfg():\n%q", got)
		})
	}
}

func TestWritePasswordCfg_rndSalt(t *testing.T) {
	cases := []struct {
		name  string
		input isoRequest
	}{
		{
			"empty",
			isoRequest{},
		},

		{
			"non-empty",
			isoRequest{
				RootPass: "Traffic Ops",
			},
		},
		{
			"long",
			isoRequest{
				RootPass: "this is a long password made longer even now",
			},
		},
	}

	const (
		expectedPrefix = "rootpw --iscrypted $1$"
		expectedPWLen  = 32
	)

	// Ensure use of random salt generates correct looking passwords.

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writePasswordCfg(&w, tc.input, ""); err != nil {
				t.Fatalf("writePasswordCfg() err = %v", err)
			}
			got := w.String()

			if !strings.HasPrefix(got, expectedPrefix) {
				t.Fatalf("writePasswordCfg() got: %q\nexpected prefix of: %q", got, expectedPrefix)
			}
			if pwLen := len(got) - len(expectedPrefix); pwLen != expectedPWLen {
				t.Fatalf("writePasswordCfg() got: %q with password length %d\nexpected password length of at least: %d", got, pwLen, expectedPWLen)
			}

			t.Logf("writePasswordCfg():\n%q", got)
		})
	}
}
