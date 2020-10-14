package util

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
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops_ort/t3c/config"
	"github.com/gofrs/flock"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const OneWeek = 604800

type FileLock struct {
	f_lock    *flock.Flock
	is_locked bool
}

type ServiceStatus int

const (
	SvcNotRunning ServiceStatus = 0
	SvcRunning    ServiceStatus = 1
	SvcUnknown    ServiceStatus = 2
)

func (s ServiceStatus) String() string {
	switch s {
	case 0:
		return "SvcNotRunning"
	case 1:
		return "SvcRunning"
	case 2:
		fallthrough
	default:
		return "SvcUnknown"
	}
}

// Try to get a file lock, non-blocking.
func (f *FileLock) GetLock(lockFile string) bool {
	f.f_lock = flock.New(lockFile)
	is_locked, err := f.f_lock.TryLock()
	f.is_locked = is_locked

	if err != nil { // some OS error attempting to obtain a file lock
		log.Errorf("unable to obtain a lock on %s\n", lockFile)
		return false
	}
	if !f.is_locked { // another process is running.
		log.Errorf("Another t3c process is already running, try again later\n")
		return false
	}

	return f.is_locked
}

// Releases a file lock and exits with the given status code.
func (f *FileLock) UnlockAndExit(code int) {
	if f.is_locked {
		f.f_lock.Unlock()
	}
	os.Exit(code)
}

func DirectoryExists(dir string) (bool, os.FileInfo) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	return info.IsDir(), info
}

func ExecCommand(fullCommand string, arg ...string) ([]byte, int, error) {
	var outbuf bytes.Buffer
	var errbuf bytes.Buffer
	cmd := exec.Command(fullCommand, arg...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()

	if err != nil {
		return outbuf.Bytes(), cmd.ProcessState.ExitCode(),
			errors.New("Error executing '" + fullCommand + "': " + errbuf.String())
	}
	return outbuf.Bytes(), cmd.ProcessState.ExitCode(), err
}

func FileExists(fn string) (bool, os.FileInfo) {
	info, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false, nil
	}
	return !info.IsDir(), info
}

func ReadFile(fn string) ([]byte, error) {
	var data []byte
	info, err := os.Stat(fn)
	if err != nil {
		return nil, err
	}
	size := info.Size()

	fd, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	data = make([]byte, size)
	c, err := fd.Read(data)
	if err != nil || int64(c) != size {
		return nil, errors.New("unable to completely read from '" + fn + "': " + err.Error())
	}
	fd.Close()

	return data, nil
}

func GetServiceStatus(name string) (ServiceStatus, int, error) {
	var pid int = -1
	var active bool = false

	output, rc, err := ExecCommand("/usr/sbin/service", name, "status")
	// service is down
	if rc == 3 {
		return SvcNotRunning, pid, nil
	} else if err != nil {
		return SvcUnknown, pid, errors.New("could not get status for service '" + name + "'\n")
	}
	lines := strings.Split(string(output), "\n")
	for ii := range lines {
		line := strings.TrimSpace(lines[ii])
		if strings.Contains(line, "Active: active") {
			active = true
		}
		if active && strings.Contains(line, "Main PID: ") {
			fmt.Sscanf(line, "Main PID: %d", &pid)
		}
	}

	if active {
		return SvcRunning, pid, nil
	} else {
		return SvcNotRunning, pid, nil
	}
}

// start or restart the service 'service'. cmd is 'start | restart'
func ServiceStart(service string, cmd string) (bool, error) {
	log.Infof("ServiceStart called for '%s'\n", service)
	svcStatus, pid, err := GetServiceStatus(service)
	if err != nil {
		return false, errors.New("Could not get status for '" + service + "' : " + err.Error())
	} else if svcStatus == SvcRunning && cmd == "start" {
		log.Infof("service '%s' is already running, pid: %d\n", service, pid)
	} else {
		_, rc, err := ExecCommand("/usr/sbin/service", service, cmd)
		if err != nil {
			return false, errors.New("Could not " + cmd + " the '" + service + "' service: " + err.Error())
		} else if rc == 0 {
			// service was sucessfully started
			return true, nil
		}
	}
	// not started, service is already running
	return false, nil
}

func WriteFile(fn string, data []byte, perm os.FileMode) (int, error) {
	return WriteFileWithOwner(fn, data, -1, -1, perm)
}

func WriteFileWithOwner(fn string, data []byte, uid int, gid int, perm os.FileMode) (int, error) {
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, errors.New("unable to open '" + fn + "' for writing: " + err.Error())
	}

	c, err := fd.Write(data)
	if err != nil {
		return 0, errors.New("error writing to '" + fn + "': " + err.Error())
	}
	fd.Close()

	if uid > -1 && gid > -1 {
		err = os.Chown(fn, uid, gid)
		if err != nil {
			return 0, errors.New("error changing ownership on '" + fn + "': " + err.Error())
		}
	}
	return c, nil
}

func PackageAction(cmdstr string, name string) (bool, error) {
	var rc int = -1
	var err error = nil
	var result bool = false

	switch cmdstr {
	case "info":
		_, rc, err = ExecCommand("/usr/bin/yum", "info", name)
	case "install":
		_, rc, err = ExecCommand("/usr/bin/yum", "install", "-y", name)
	case "remove":
		_, rc, err = ExecCommand("/usr/bin/yum", "remove", "-y", name)
	}

	if rc == 0 {
		result = true
		err = nil
	}
	return result, err
}

// runs the rpm command.
// if the return code from rpm == 0, then a valid package list is returned.
//
// if the return code is 1, the the 'name' queried for is not part of a
//   package or is not installed.
//
// otherwise, if the return code is not 0 or 1 and error is set, a general
// rpm command execution error is assumed and the error is returned.
func PackageInfo(cmdstr string, name string) ([]string, error) {
	var result []string
	switch cmdstr {
	case "cfg-files": // returns a list of the package configuration files.
		output, rc, err := ExecCommand("/bin/rpm", "-q", "-c", name)
		if rc == 1 { // rpm package for 'name' was not found.
			return nil, nil
		} else if rc == 0 { // add the package name the file belongs to.
			log.Debugf("output from cfg-files query: %s\n", string(output))
			files := strings.Split(string(output), "\n")
			for ii := range files {
				result = append(result, strings.TrimSpace(files[ii]))
			}
			log.Debugf("result length: %d, result: %s\n", len(result), string(output))
		} else if err != nil {
			return nil, err
		}
	case "file-query": // returns the rpm name that owns the file 'name'
		output, rc, err := ExecCommand("/bin/rpm", "-q", "-f", name)
		if rc == 1 { // file is not part of any package.
			return nil, nil
		} else if rc == 0 { // add the package name the file belongs to.
			log.Debugf("output from file-query: %s\n", string(output))
			result = append(result, string(strings.TrimSpace(string(output))))
			log.Debugf("result length: %d, result: %s\n", len(result), string(output))
		} else if err != nil {
			return nil, err
		}
	case "pkg-provides": // returns the package name that provides 'name'
		output, rc, err := ExecCommand("/bin/rpm", "-q", "--whatprovides", name)
		log.Debugf("pkg-provides - name: %s, output: %s\n", name, output)
		if rc == 1 { // no package provides 'name'
			return nil, nil
		} else if rc == 0 {
			pkgs := strings.Split(string(output), "\n")
			for ii := range pkgs {
				result = append(result, strings.TrimSpace(pkgs[ii]))
			}
		} else if err != nil {
			return nil, err
		}
	case "pkg-query": // returns the package name for 'name'.
		output, rc, err := ExecCommand("/bin/rpm", "-q", name)
		if rc == 1 { // the package is not installed.
			return nil, nil
		} else if rc == 0 { // add the rpm name
			result = append(result, string(strings.TrimSpace(string(output))))
		} else if err != nil {
			return nil, err
		}
	case "pkg-requires": // returns a list of packages that requires package 'name'
		output, rc, err := ExecCommand("/bin/rpm", "-q", "--whatrequires", name)
		if rc == 1 { // no package reuires package 'name'
			return nil, nil
		} else if rc == 0 {
			pkgs := strings.Split(string(output), "\n")
			for ii := range pkgs {
				result = append(result, strings.TrimSpace(pkgs[ii]))
			}
		} else if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func RandomDuration(max time.Duration) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Int63n(int64(max)))
}

func CheckUser(cfg config.Cfg) bool {
	result := true
	userInfo, err := user.Current()

	if err != nil {
		log.Errorf("could not obtain the current user info: %s\n", err.Error())
		return false
	}

	switch cfg.RunMode {
	case config.BadAss:
		fallthrough
	case config.SyncDS:
		if userInfo.Username != "root" {
			log.Errorf("Only the root user may run in BadAss, or SyncDS mode, current user: %s\n",
				userInfo.Username)
			result = false
		}
	default:
		log.Infof("current mode: %s, run user: %s\n", cfg.RunMode, userInfo.Username)
	}
	return result
}

func CleanTmpDir() bool {
	now := time.Now().Unix()

	if len(config.TmpBase) == 0 || !strings.HasPrefix(config.TmpBase, "/") {
		log.Errorf("config.TmpBase is misconfigured: '%s', refusing to remove any files or directories.", config.TmpBase)
		return false
	}

	files, err := ioutil.ReadDir(config.TmpBase)
	if err != nil {
		log.Errorf("opening %s: %v\n", config.TmpBase, err)
		return false
	}
	for _, f := range files {
		// remove any directory and its contents under the config.TmpBase (/tmp/ort) that
		// is more than a week old.
		if f.IsDir() && (now-f.ModTime().Unix()) > OneWeek {
			dir := filepath.Join(config.TmpBase, f.Name())
			if dir == "/" {
				log.Errorf("config.TmpBase is incorrectly configured, refusing to remove '/'. check config.TmpBase.")
				return false
			}
			if f.Name() == "." || f.Name() == ".." {
				continue
			}
			err = os.RemoveAll(dir)
			if err != nil {
				log.Errorf("could not remove '%s': %v\n", dir, err)
				return false
			}
		}
	}
	return true
}

func MkDir(name string, cfg config.Cfg) bool {
	fileInfo, err := os.Stat(name)
	if err == nil && fileInfo.Mode().IsDir() {
		log.Debugf("the directory '%s' already exists", name)
		return true
	}
	if err != nil {
		if cfg.RunMode != config.Report {
			if err != nil { // the path does not exist.
				err = os.Mkdir(name, 0755)
				if err != nil {
					log.Errorf("unable to create the directory '%s': %v", name, err)
					return false
				}
			} else if fileInfo.Mode().IsDir() {
				log.Debugf("the directory: %s, already exists\n", name)
			} else {
				log.Errorf("there is a file named, '%s' that is not a directory: %v", name, err)
				return false
			}
		} else {
			log.Infof("the directory %s does not exist and was not created, runMode: %s", name, cfg.RunMode)
			return true
		}
	}
	return false
}

func Touch(fn string) error {
	myfile, err := os.Create(fn)
	if err != nil {
		return err
	}
	myfile.Close()
	return nil
}
