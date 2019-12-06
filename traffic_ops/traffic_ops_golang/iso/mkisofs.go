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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/jmoiron/sqlx"
)

const (
	ksCfgNetwork     = "network.cfg"
	ksCfgMgmtNetwork = "mgmt_network.cfg"
	ksCfgPassword    = "password.cfg"
	ksCfgDisk        = "disk.cfg"
	ksStateOut       = "state.out"
)

func genISO(w io.Writer, tx *sqlx.Tx, r isoRequest, isoDest string) error {
	streamISO := r.Stream.val

	ksDir, err := kickstarterDir(tx, r.OSVersionDir)
	if err != nil {
		return err
	}
	cfgDir := filepath.Join(ksDir, ksCfgDir)
	log.Infof("cfg_dir: %s", cfgDir)

	if err = writeKSCfgs(cfgDir, r); err != nil {
		return err
	}

	if streamISO {
		cmd := newStreamISOCmd(ksDir, isoDest)
		log.Infof("Using %s ISO generation command: %s", cmd.cmdType, cmd.String())

		// Create state.out
		stateOut := fmt.Sprintf("Dir== %s\n%s\n", ksDir, cmd.String())
		if err = ioutil.WriteFile(filepath.Join(cfgDir, ksStateOut), []byte(stateOut), 0666); err != nil {
			return err
		}

		return cmd.stream(w)
	}

	cmd := newSaveISOCmd(ksDir, isoDest)
	log.Infof("Using %s ISO generation command: %s", cmd.cmdType, cmd.String())

	// Create state.out
	stateOut := fmt.Sprintf("Dir== %s\n%s\n", ksDir, cmd.String())
	if err = ioutil.WriteFile(filepath.Join(cfgDir, ksStateOut), []byte(stateOut), 0666); err != nil {
		return err
	}

	result, err := cmd.save()
	log.Infoln(result)

	return err
}

func newStreamISOCmd(ksDir, isoDest string) *streamISOCmd {
	s := streamISOCmd{
		isoDest: isoDest,
	}

	if customExec := customGenISOPath(ksDir); customExec != "" {
		s.cmdType = "custom"
		s.cmd = exec.Command(customExec, isoDest)
		s.isSavedToDisk = true
	} else {
		s.cmdType = "default"
		s.cmd = exec.Command(
			"mkisofs",
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
	}

	return &s
}

type streamISOCmd struct {
	cmd           *exec.Cmd
	cmdType       string
	isSavedToDisk bool
	isoDest       string
}

func (s *streamISOCmd) String() string {
	return strings.Join(s.cmd.Args, " ")
}

func (s *streamISOCmd) stream(w io.Writer) error {
	if s.isSavedToDisk {
		return s.streamFromFile(w)
	}

	return s.streamStdout(w)
}

func (s *streamISOCmd) streamStdout(w io.Writer) error {
	s.cmd.Stdout = w
	return s.cmd.Run()
}

func (s *streamISOCmd) streamFromFile(w io.Writer) error {
	if err := s.cmd.Run(); err != nil {
		return err
	}
	defer os.Remove(s.isoDest)

	isoFd, err := os.Open(s.isoDest)
	if err != nil {
		return err
	}
	defer isoFd.Close()

	_, err = io.Copy(w, bufio.NewReader(isoFd))
	return err
}

func newSaveISOCmd(ksDir, isoDest string) *saveISOCmd {
	var s saveISOCmd

	if customExec := customGenISOPath(ksDir); customExec != "" {
		s.cmdType = "custom"
		s.cmd = exec.Command(customExec, isoDest)
	} else {
		s.cmdType = "default"
		s.cmd = exec.Command(
			"mkisofs",
			"-o", isoDest, // output to file instead of STDOUT
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
	}

	return &s
}

type saveISOCmd struct {
	cmd     *exec.Cmd
	cmdType string
}

func (s *saveISOCmd) String() string {
	return strings.Join(s.cmd.Args, " ")
}

func (s *saveISOCmd) save() (string, error) {
	result, err := s.cmd.CombinedOutput()
	return string(result), err
}

func customGenISOPath(dir string) string {
	// Allow for a custom script to be used instead of the default command.
	// The script must:
	// - Be inside the ksDir and named "generate"
	// - Be executable (by somebody)
	// - Accept a single argument indicating where the resulting ISO image should be saved

	customPath := filepath.Join(dir, "generate")
	stat, err := os.Stat(customPath)

	// Check if file exists and is executable
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
func writeKSCfgs(dir string, r isoRequest) error {
	nameservers, err := readDefaultUnixResolve()
	if err != nil {
		return err
	}

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
	cfg.addOpt("MTU", strconv.Itoa(r.InterfaceMTU))
	cfg.addOpt("NAMESERVER", strings.Join(nameservers, ","))
	cfg.addOpt("HOSTNAME", r.fqdn())
	cfg.addOpt("NETWORKING_IPV6", "yes")
	cfg.addIP("IPV6ADDR", r.IP6Address)
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
//   OPT="VALUE"
//   OPT="VALUE"
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

// addIP adds an IP option to the config. It handles the
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
	if b.val {
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
