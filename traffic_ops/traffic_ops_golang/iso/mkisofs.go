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
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	mkisofsBin = "mkisofs" // name of the binary that's used in the default case to generate an ISO
)

// newStreamISOCmd returns a instantiated streamISOCmd. The given ksDir
// is expected to be the root directory containing the kickstarter files
// of the desired OS. It will detect a custom `generate` executable if present,
// otherwise will use the default `mkisofs` command.
func newStreamISOCmd(ksDir string) (*streamISOCmd, error) {
	var s streamISOCmd

	if customExec := customGenISOPath(ksDir); customExec != "" {
		// The custom script must accept a single argument: The path
		// where it will write the ISO. Here we create a temporary
		// directory for this purpose. The cleanup method is responsible
		// for removing it.
		tmpDir, err := ioutil.TempDir("", "genISO")
		if err != nil {
			return nil, err
		}
		s.isoDest = filepath.Join(tmpDir, "tmp.iso")

		s.cmdType = "custom"
		s.cmd = exec.Command(customExec, s.isoDest)

		return &s, nil
	}

	s.cmdType = "default"
	s.cmd = exec.Command(
		mkisofsBin,
		"-joliet-long",
		"-input-charset", "utf-8",
		"-b", "isolinux/isolinux.bin",
		"-c", "isolinux/boot.cat",
		"-no-emul-boot",
		"-boot-load-size", "4",
		"-boot-info-table",
		"-R",
		"-J",
		"-v",
		"-T",
		ksDir,
	)

	return &s, nil
}

// streamISOCmd encapsulate the logic for executing the ISO
// generation command.
type streamISOCmd struct {
	cmd     *exec.Cmd
	cmdType string // Description of command: "default" or "custom"

	// If empty, then cmd writes to STDOUT. Othewrise, cmd writes the
	// ISO to this path.
	isoDest string
}

// String returns the command that the stream method
// will execute.
func (s *streamISOCmd) String() string {
	// Note: Go 1.13 adds exec.Cmd#String method
	return strings.Join(s.cmd.Args, " ")
}

// cleanup should be defered after calling newStreamISOCmd.
// It removes any temporary resources created. If the command
// doesn't need any cleanup, this is a no-op.
func (s *streamISOCmd) cleanup() error {
	if s.isoDest == "" {
		return nil
	}
	return os.RemoveAll(filepath.Dir(s.isoDest))
}

// stream writes to w the ISO data. Callers should
// always use this method and not the other more
// specific stream methods.
func (s *streamISOCmd) stream(w io.Writer) error {
	if s.isoDest != "" {
		return s.streamFromFile(w)
	}

	return s.streamStdout(w)
}

// streamStdout invokes the command and pipes its STDOUT
// to w.
func (s *streamISOCmd) streamStdout(w io.Writer) error {
	var stderr bytes.Buffer
	s.cmd.Stdout = w
	s.cmd.Stderr = &stderr
	if err := s.cmd.Run(); err != nil {
		return fmt.Errorf("%v: %s", err, &stderr)
	}
	return nil
}

// streamFromFile invokes the command and expects the ISO
// to be written to isoDest. It then copies the contents
// of that file to w.
func (s *streamISOCmd) streamFromFile(w io.Writer) error {
	var stderr bytes.Buffer
	s.cmd.Stderr = &stderr
	if err := s.cmd.Run(); err != nil {
		return fmt.Errorf("%v: %s", err, &stderr)
	}

	isoFd, err := os.Open(s.isoDest)
	if err != nil {
		return err
	}
	defer isoFd.Close()

	_, err = io.Copy(w, bufio.NewReader(isoFd))
	return err
}

// customGenISOPath returns the complete path to an alternative executable
// for generating the ISO. If not found, an empty string is returned.
// In order to be valid, the script/executable must:
//   - Be inside the ksDir and named "generate"
//   - Be executable (by somebody)
//   - Accept a single argument indicating where the resulting ISO image should be saved (not
//     verified by this function)
func customGenISOPath(dir string) string {
	customPath := filepath.Join(dir, ksAltCommand)
	stat, err := os.Stat(customPath)

	// Check if file exists and is executable.
	const executablePermBits = 0111
	if err == nil && stat.Mode().Perm()&executablePermBits != 0 {
		return customPath
	}
	return ""
}

// kickstarterDir returns the directory containing the kickstarter files for
// the given OS. This is the directory passed to the mkisofs command.
// The root part of the directory can be overriden with a Parameter database entry.
func kickstarterDir(tx *sqlx.Tx, osVersionDir string) (string, error) {
	var baseDir string
	err := tx.QueryRow(
		`SELECT value FROM parameter WHERE name = $1 AND config_file = $2 LIMIT 1`,
		ksFilesParamName,
		ksFilesParamConfigFile,
	).Scan(&baseDir)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if baseDir == "" {
		baseDir = cfgDefaultDir
	}

	return filepath.Join(baseDir, osVersionDir), nil
}

// writeKSCfgs writes to the given directory the various Kickstart
// configuration files using the data from the given isoRequest.
// The cmd string is used to log to a file the command that will
// be executed to create the ISO.
func writeKSCfgs(dir string, r isoRequest, cmd string) error {
	nameservers, err := readDefaultUnixResolve()
	if err != nil {
		return err
	}

	// Create state.out

	stateFd, err := os.Create(filepath.Join(dir, ksStateOut))
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(stateFd, "Dir== %s\n%s\n", dir, cmd); err != nil {
		return err
	}
	defer stateFd.Close()

	// Create network.cfg

	networkCfgFd, err := os.Create(filepath.Join(dir, ksCfgNetwork))
	if err != nil {
		return err
	}
	defer networkCfgFd.Close()
	if err = writeNetworkCfg(networkCfgFd, r, nameservers); err != nil {
		return err
	}

	// Create mgmt_network.cfg

	mgmtNetworkCfgFd, err := os.Create(filepath.Join(dir, ksCfgMgmtNetwork))
	if err != nil {
		return err
	}
	defer mgmtNetworkCfgFd.Close()
	if err = writeMgmtNetworkCfg(mgmtNetworkCfgFd, r); err != nil {
		return err
	}

	// Create password.cfg

	passwordCfgFd, err := os.Create(filepath.Join(dir, ksCfgPassword))
	if err != nil {
		return err
	}
	defer passwordCfgFd.Close()
	// Empty salt parameter causes a random salt to be generated,
	// which is the desired behavior.
	if err = writePasswordCfg(passwordCfgFd, r, ""); err != nil {
		return err
	}

	// Create disk.cfg

	diskCfgFd, err := os.Create(filepath.Join(dir, ksCfgDisk))
	if err != nil {
		return err
	}
	defer diskCfgFd.Close()
	if err = writeDiskCfg(diskCfgFd, r); err != nil {
		return err
	}

	return nil
}

// bondedRegex matches a bonded device interface name.
var bondedRegex = regexp.MustCompile(`^bond\d+`)

// writeNetworkCfg writes the network.cfg config to w.
func writeNetworkCfg(w io.Writer, r isoRequest, nameservers []string) error {
	var cfg configWriter

	cfg.addIP("IPADDR", r.IPAddr)
	cfg.addIP("NETMASK", r.IPNetmask)
	cfg.addIP("GATEWAY", r.IPGateway)
	isBonded := bondedRegex.MatchString(r.InterfaceName)
	if isBonded {
		cfg.addOpt("BOND_DEVICE", r.InterfaceName)
	} else {
		cfg.addOpt("DEVICE", r.InterfaceName)
	}
	cfg.addOpt("MTU", r.InterfaceMTU.String())
	cfg.addOpt("NAMESERVER", strings.Join(nameservers, ","))
	cfg.addOpt("HOSTNAME", r.fqdn())
	cfg.addOpt("NETWORKING_IPV6", "yes")
	cfg.addOpt("IPV6ADDR", r.IP6Address)
	cfg.addIP("IPV6_DEFAULTGW", r.IP6Gateway)
	if isBonded {
		cfg.addOpt("BONDING_OPTS", "miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4")
	}
	cfg.addBoolStr("DHCP", r.DHCP)

	_, err := io.Copy(w, &cfg)
	return err
}

// writeMgmtNetworkCfg writes the mgmt_network.cfg config to w.
func writeMgmtNetworkCfg(w io.Writer, r isoRequest) error {
	var cfg configWriter

	// Test if management IP is IPv6
	if r.MgmtIPAddress.To16() != nil && r.MgmtIPAddress.To4() == nil {
		cfg.addIP("IPV6ADDR", r.MgmtIPAddress)
	} else {
		cfg.addIP("IPADDR", r.MgmtIPAddress)
	}
	cfg.addIP("NETMASK", r.MgmtIPNetmask)
	cfg.addIP("GATEWAY", r.MgmtIPGateway)
	cfg.addOpt("DEVICE", r.MgmtInterface)

	_, err := io.Copy(w, &cfg)
	return err
}

// writeDiskCfg writes the disk.cfg config to w.
func writeDiskCfg(w io.Writer, r isoRequest) error {
	var cfg configWriter

	cfg.addOpt("boot_drives", r.Disk)

	_, err := io.Copy(w, &cfg)
	return err
}

// writePasswordCfg writes the password.cfg config to w.
// The salt parameter is optional. If salt is blank, then a
// random 8-character salt will be used.
func writePasswordCfg(w io.Writer, r isoRequest, salt string) error {
	if salt == "" {
		salt = rndSalt(8)
	}

	cryptedPw, err := crypt(r.RootPass, salt)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "rootpw --iscrypted %s\n", cryptedPw)
	return err
}

// configWriter is a helper type to create config files
// of format:
//
//	OPT="VALUE"
//	OPT="VALUE"
type configWriter struct {
	b    bytes.Buffer
	line int
}

// addOpt adds an option to the config.
func (c *configWriter) addOpt(name, value string) {
	if c.line > 0 {
		fmt.Fprintln(&c.b)
	}
	c.line++
	fmt.Fprintf(&c.b, "%s=%q", name, value)
}

// addIP adds an IPv4 option to the config. It handles the
// case where the given IP is empty/nil.
func (c *configWriter) addIP(name string, ip net.IP) {
	// Avoid using `<nil>`, i.e. net.IP{}.String() = <nil>
	var v string
	if len(ip) > 0 {
		v = ip.String()
	}
	c.addOpt(name, v)
}

// addBoolStr adds a BoolStr option to the config, using
// "yes" / "no" values.
func (c *configWriter) addBoolStr(name string, b boolStr) {
	var v string
	if bv, _ := b.val(); bv {
		v = "yes"
	} else {
		v = "no"
	}
	c.addOpt(name, v)
}

// Read satisfies the io.Reader interface, and allows for
// using the configWriter by the io.Copy() function.
func (c *configWriter) Read(p []byte) (int, error) {
	return c.b.Read(p)
}
