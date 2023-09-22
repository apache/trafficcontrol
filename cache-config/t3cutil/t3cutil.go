package t3cutil

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
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

type ATSConfigFile struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	ContentType string   `json:"content_type"`
	LineComment string   `json:"line_comment"`
	Secure      bool     `json:"secure"`
	Text        string   `json:"text"`
	Warnings    []string `json:"warnings"`
}

var installdir string

// initializes InstallDir to executable dir
// If error, returns "/usr/bin" as default.
func InstallDir() string {
	if installdir == "" {
		execpath, err := os.Executable()
		if err != nil {
			installdir = `/usr/bin`
			log.Infof("InstallDir setting to fallback: '%s', %v\n", installdir, err)
		} else {
			log.Infof("Executable path is %s", execpath)
			installdir = filepath.Dir(execpath)
		}
	}
	log.Infof("Return Installdir '%s'", installdir)
	return installdir
}

// ATSConfigFiles implements sort.Interface and sorts by the Location and then FileNameOnDisk, i.e. the full file path.
type ATSConfigFiles []ATSConfigFile

func (fs ATSConfigFiles) Len() int { return len(fs) }
func (fs ATSConfigFiles) Less(i, j int) bool {
	if fs[i].Path != fs[j].Path {
		return fs[i].Path < fs[j].Path
	}
	return fs[i].Name < fs[j].Name
}
func (fs ATSConfigFiles) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

// CommentsFilter is used to remove comment
// lines from config files while making
// comparisons.
func CommentsFilter(body []string, lineComment string) []string {
	var newlines []string

	newlines = make([]string, 0)

	for ii := range body {
		line := body[ii]
		if strings.HasPrefix(line, lineComment) {
			continue
		}
		newlines = append(newlines, line)
	}

	return newlines
}

// PermCk will compare file permissions against existing file and octal permission provided.
func PermCk(path string, perm int) bool {
	mode := os.FileMode(perm)
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error getting file status", path)
	}
	if file.Mode() != mode.Perm() {
		return true
	}
	return false
}

// OwnershipCk will compare owner and group settings against existing file and owner/group settings provided.
func OwnershipCk(path string, uid int, gid int) bool {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println("error getting file status", path)
	}
	stat := file.Sys().(*syscall.Stat_t)
	if uid != int(stat.Uid) || gid != int(stat.Gid) {
		return true
	}
	return false
}

// NewLineFilter removes carriage returns
// from config files while making comparisons.
func NewLineFilter(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\n")
	return strings.TrimSpace(str)
}

// ReadFile reads a file and returns the
// file contents.
func ReadFile(f string) []byte {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println("Error reading file ", f)
		os.Exit(1)
	}
	return data
}

// UnencodeFilter translates HTML escape
// sequences while making config file comparisons.
func UnencodeFilter(body []string) []string {
	var newlines []string

	newlines = make([]string, 0)
	sp := regexp.MustCompile(`\s+`)
	el := regexp.MustCompile(`^\s+|\s+$`)

	for ii := range body {
		s := body[ii]
		s = sp.ReplaceAllString(s, " ")
		s = el.ReplaceAllString(s, "")
		s = html.UnescapeString(s)
		s = strings.TrimSpace(s)
		newlines = append(newlines, s)
	}

	return newlines
}

// Do executes the given command and returns the stdout, stderr, and exit code.

// This is a convenience wrapper around os/exec.
// Since t3c only needs to make simple calls and get the stdout, stderr, and code, this provides a simpler and terser interface.
//
// If you need anything more complex, or don't find this simpler, you should probably use os/exec directly.
//
// Each arg must be passed as its own string. Unfortunately, Go doesn't have a way to pass multiple args as a single string, and splitting on spaces would require complex quote parsing.
//
// Note each arg must be passed without quotes. Go calls the app with args as if they were quoted. if you add quotes, they'll be passed to the command literally, as if you called 'mycommand "\"escaped-quotes\""`.
//
// Note if Go fails to run the command, the error from Go will be returned as the stderr and the code -1,
// which will differ from what would have been returned by a command line.
func Do(cmdStr string, args ...string) ([]byte, []byte, int) {
	cmd := exec.Command(cmdStr, args...)

	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	code := 0
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return nil, []byte(err.Error()), -1
		} else {
			code = exitErr.ExitCode()
		}
	}

	return outbuf.Bytes(), errbuf.Bytes(), code
}

// DoInput is like Do but takes the stdin to pass to the command.
func DoInput(input []byte, cmdStr string, args ...string) ([]byte, []byte, int) {
	cmd := exec.Command(cmdStr, args...)

	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Stdin = bytes.NewBuffer(input)

	code := 0
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return nil, []byte(err.Error()), -1
		} else {
			code = exitErr.ExitCode()
		}
	}

	return outbuf.Bytes(), errbuf.Bytes(), code
}

func ValidateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme expected 'http' or 'https', actual '" + u.Scheme + "'")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("no host")
	}
	return nil
}

// VersionStr returns a common version string format for all t3c apps.
// The appName is the command itself, e.g. t3c-apply.
// The versionNum is the version number from the build system. It should include the major, minor, git revision, and a monotonically increasing number, e.g. '4.2.1234.abc123'.
func VersionStr(appName string, versionNum string, gitRevision string) string {
	if len(gitRevision) > 8 {
		gitRevision = gitRevision[:8]
	}
	return appName + " " + versionNum + ".." + gitRevision
}

func UserAgentStr(appName string, versionNum string, gitRevision string) string {
	if len(gitRevision) > 8 {
		gitRevision = gitRevision[:8]
	}
	return appName + "/" + versionNum + ".." + gitRevision
}

// NewApplyMetaData creates a new, empty ApplyMetaData object.
func NewApplyMetaData() *ApplyMetaData {
	return &ApplyMetaData{
		Version:           MetaDataVersion,
		InstalledPackages: []string{},              // construct a slice, so JSON serializes '[]' not 'null'.
		OwnedFilePaths:    []string{},              // construct a slice, so JSON serializes '[]' not 'null'.
		Actions:           []ApplyMetaDataAction{}, // construct a map, so JSON serializes '{}' not 'null'.
	}
}

// MetaDataVersion is the version of the metadata file.
// This should update the major version with breaking changes,
// and t3c versions should strive to maintain compatibility
// at least one major version back, so features like tracking
// t3c-owned files continue to work through upgrades.
const MetaDataVersion = "1.0"

// ApplyMetaData is metadata about a t3c-apply run.
// Always use NewApplyMetaData, don't use a literal to construct a new object.
type ApplyMetaData struct {
	// Version is the metadata version of this metadata object or file. See MetaDataVersion.
	Version string `json:"version"`

	// ServerFQDN is the FQDN of this server.
	// The primary purpose of this field is to allow distinguishing
	// metadata files from different servers.
	ServerHostName string `json:"server-hostname"`

	// ReloadedATS is whether this run restarted ATS.
	// Note this is whether ATS was actually restarted, not whether it would have been,
	// e.g. because of --report-only or --service-action.
	ReloadedATS bool `json:"reloaded-ats"`

	// RestartedATS is whether this run restarted ATS.
	// Note this is whether ATS was actually restarted, not whether it would have been,
	// e.g. because of --report-only or --service-action.
	RestartedATS bool `json:"restarted-ats"`

	// UnsetUpdateFlag is whether this t3c-apply run unset the update flag for this server.
	// Note this is whether the flag was actually unset, not whether it would have been e.g.
	// because of --no-unset-update-flag or --report-only.
	UnsetUpdateFlag bool `json:"unset-update-flag"`

	// UnsetRevalFlag is whether this t3c-apply run unset the revalidate flag for this server.
	// Note this is whether the flag was actually unset, not whether it would have been e.g.
	// because of --no-unset-reval-flag or --report-only.
	UnsetRevalFlag bool `json:"unset-reval-flag"`

	// InstalledPackages is which yum packages are installed.
	// Note this packages currently installed, not what would have been e.g.
	// because of --install-packages=false or --report-only.
	InstalledPackages []string `json:"installed-packages"`

	// OwnedFilePaths is the list of files t3c-apply produced in this run.
	//
	// This can be used to know which files in the ATS config directory were produced by t3c,
	// and which were produced by some other means.
	//
	// Note this is all files produced, not necessarily all files written to disk. This
	// will include files generated, but not changed on disk because they had no
	// semantic diff from the existing file.
	//
	// This may be used in the future for t3c-apply to delete files produced by a previous
	// run which no longer exist (for example, Header Rewrites from a Delivery Service
	// no longer assigned to this server).
	//
	// Files are the full path and file name.
	OwnedFilePaths []string `json:"owned-files"`

	// Succeeded is whether this t3c-apply run generally succeeded.
	//
	// Note not all scenarios are black or white success-or-fail.
	// For example, files may be successfully created, but reloading ATS may fail.
	// In these scenarios, t3c-apply will attempt to set Succeeded to false,
	// but also attempt to set other metadata about what was actually performed.
	//
	// But when partial failure occurrs, nothing is guaranteed in the metadata.
	// Operators should consider the logs authoritative over the metadata.
	Succeeded bool `json:"succeeded"`

	// PartialSuccess indicates that some actions were successful, but
	// later actions failed.
	//
	// This is a bad place to be, because it means some things were changed,
	// but not everything that needed to be. This is often not fatal, because,
	// for example, if config files were changed by ATS failed to reload,
	// those config files typically needed placed anyway.
	//
	// But nevertheless, partial success is potentially catastrophic, and operators
	// are strongly encouraged to set alarms and read logs in the event it occurs,
	// to determine what was changed, what failed, and what actions need taken.
	PartialSuccess bool `json:"partial-success"`

	Actions []ApplyMetaDataAction `json:"actions"`
}
type ApplyMetaDataAction struct {
	Action string `json:"action"`
	Status string `json:"status"`
}

// Format prints the ApplyMetaData in a format designed to be written to a file,
// and structured but pretty-printed to work well with line-based diffs (e.g. in git).
func (md *ApplyMetaData) Format() ([]byte, error) {
	bts, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return nil, errors.New("marshalling metadata file: " + err.Error())
	}
	bts = append(bts, '\n') // newline at the end of the file, so it's a valid POSIX text file

	return bts, nil
}

func PackagesToMetaData(pkg map[string]bool) []string {
	pkgs := []string{}
	for k, v := range pkg {
		if v {
			pkgs = append(pkgs, k)
		}
	}
	sort.Strings(pkgs)
	return pkgs
}

// CombineOwnedFilePaths combines the owned file paths of two metadata objects.
//
// This is primarily useful when a config run, such as revalidate, adds owned files, but not
// all owned files, but we don't want to write metadata incidating we don't own existing files,
// so this can be used to combine the new files with the previous metadata.
//
// Both am and bm are may be nil, in which case the files from the non-nil object is returned,
// or an empty array if both are nil.
func CombineOwnedFilePaths(am *ApplyMetaData, bm *ApplyMetaData) []string {
	if am == nil && bm == nil {
		return []string{}
	} else if am == nil {
		sort.Strings(bm.OwnedFilePaths) // the func guarantees the returned array will always be sorted
		return bm.OwnedFilePaths
	} else if bm == nil {
		sort.Strings(am.OwnedFilePaths) // the func guarantees the returned array will always be sorted
		return am.OwnedFilePaths
	}
	return sortAndCombineStrs(am.OwnedFilePaths, bm.OwnedFilePaths)
}

// sortAndCombineStrs sorts as and bs, and then returns an array containing
// the unique strings in each, without duplicates.
func sortAndCombineStrs(as []string, bs []string) []string {
	sort.Strings(as)
	sort.Strings(bs)
	combined := []string{}
	ai := 0
	bi := 0
	for ai < len(as) && bi < len(bs) {
		if as[ai] == bs[bi] {
			combined = append(combined, as[ai])
			ai++
			bi++
			continue
		}
		// at this point we know they don't match
		// so add the lesser, increment it, and loop (but don't add or increment the greater)
		if as[ai] < bs[bi] {
			combined = append(combined, as[ai])
			ai++
			continue
		}
		combined = append(combined, bs[bi])
		bi++
	}

	// at this point, we added everything up to the end of one of the arrays,
	// but potentially not the other. So add the remaining strings in the other

	for ai < len(as) {
		combined = append(combined, as[ai])
		ai++
	}
	for bi < len(bs) {
		combined = append(combined, bs[bi])
		bi++
	}
	return combined
}

// CheckRefsInputFileAndAdding is the input (stdin or file) for t3c-check-refs if
// --files-adding=input. If not, the input is simply the raw file to check.
type CheckRefsInputFileAndAdding struct {
	File   []byte   `json:"file"`
	Adding []string `json:"adding"`
}
