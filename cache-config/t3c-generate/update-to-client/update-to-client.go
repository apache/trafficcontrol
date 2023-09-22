package main

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
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	log.Init(os.Stderr, os.Stderr, os.Stderr, os.Stderr, os.Stderr)
	if len(os.Args) < 3 {
		log.Errorln(usageStr())
		os.Exit(0)
	}
	dir := os.Args[1]
	branch := os.Args[2]

	workingDir, err := os.Getwd()
	if err != nil {
		log.Errorf("Error getting working directory: %s\n", err.Error())
		os.Exit(1)
	}

	dir = filepath.Join(workingDir, dir) // make the given directory absolute

	if err := updateVendoredTOClient(dir, branch); err != nil {
		log.Errorf("Error updating vendored client: %s\n", err.Error())
		os.Exit(1)
	}
	if err := updateNewClientUsage(dir); err != nil {
		log.Errorf("Error updating new client usage: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func usageStr() string {
	return `usage: go run update-to-client.go /path/to/cache-config/t3c-generate branch-to-vendor

Example: go run update-to-client/update-to-client.go . 5.0.x

This script updates cache-config/t3c-generate after a Traffic Ops release is made.

Be aware! It can and will modify all .go files in the given directory! Back up any uncommitted changes!

Also be aware! This is very specific to the current code. If symbols or patterns are changed around how the master vs vendored client are used, this script will have to be updated.

Expecations:
- t3c-generate is at github.com/apache/trafficcontrol/v8/cache-config/t3c-generate
- The master TO client is at github.com/apache/trafficcontrol/v8/traffic_ops/v4-client
- The previous major version client is vendored at t3c-generate/toreq/vendor
- The master client wrapper for t3c-generate is at t3c-generate/toreqnew
- The clients are stored in config.TCCfg.TOClient and config.TCCfg.TOClientNew
- Every func in toreqnew.TOClient has a corresponding func with the same name in toreq.TOClient
- Every func in toreqnew.TOClient returns 3 variables: the object, boolean whether the request was unsupported, and an error.
- Every func in toreq.TOClient returns 2 variables: the object, and an error.
- Every usage of config.TOClientNew is immediately followed by a check 'if err == nil && unsupported', whose block calls the old client and sets defaults for the unsupported new feature.
- The script is running on a POSIX-like environment. Namely, cp and gofmt exist.

The arguments are the t3c-generate directory, and the name of the branch to vendor from.

This script should always be called from trafficcontrol/cache-config/t3c-generate.

It copies the traffic_ops/v4-client from that branch into toreq/vendor,
and then updates all references to cfg.TOClientNew to cfg.TOClient.

This must be done as soon as a release is made, before any new features are added to t3c-generate.
Thus, all existing toreqnew.TOClient function calls should exist in master.

If any new features were added after a release, before you ran this script,
you'll have to go back and fix compile errors for new client functions.

Further, if the new features were fields on existing functions, being moved
from toreqnew to toreq will make them still compile, but the fields will be nil!

It will also run gofmt on all go files in the directory (because it's much easier to manipulate code guaranteed to be gofmt'd).

Double check the results before you commit!
`
}

func updateNewClientUsage(appDir string) error {
	paths := []string{}
	err := filepath.Walk(appDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !strings.HasSuffix(path, `.go`) {
				return nil // skip anything not a go file.
			}
			if strings.Contains(path, `/vendor/`) {
				return nil // skip vendored code
			}
			if strings.HasPrefix(filepath.Base(path), `.`) {
				return nil // skip .files (usually created by editors)
			}
			if strings.Contains(path, `update-to-client.go`) {
				return nil // skip this file
			}

			paths = append(paths, path)
			return nil
		})
	if err != nil {
		return errors.New("reading directory: " + err.Error() + "\n")
	}

	toClientNewRe := regexp.MustCompile(`(\s+)([^\s]+), (.+), (.+) := (.+).TOClientNew.(.+)$`)

	for _, path := range paths {
		// gofmt the file, because the parsing relies on things being in that exact format.
		if err := exec.Command(`gofmt`, `-w`, path).Run(); err != nil {
			return errors.New("running gofmt on '" + path + ": " + err.Error() + "\n")
		}

		bts, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.New("reading file '" + path + "': " + err.Error() + "\n")
		}

		lines := strings.Split(string(bts), "\n")
		newLines := make([]string, 0, len(lines))

		const stateStart = 0
		const stateInFoundTOClientNew = 1
		const stateInTOClientNewUnsupportedBlock = 2

		state := stateStart
		spacePrefix := ""
		objVar := ""
		unsupportedVar := ""
		errVar := ""
		for _, line := range lines {
			switch state {
			case stateStart:
				if strings.Contains(line, `.TOClientNew.`) {
					matches := toClientNewRe.FindStringSubmatch(line)
					if len(matches) != 7 {
						// TODO is this really an error? Skip and continue?
						return fmt.Errorf("parsing file '"+path+"': line contains TOClientNew, but is unexpected format (only %v matches): '"+line+"'\n", len(matches))
					}
					spacePrefix = matches[1]
					objVar = matches[2]
					unsupportedVar = matches[3]
					errVar = matches[4]
					cfg := matches[5]
					funcStr := matches[6]
					newLine := spacePrefix + objVar + `, ` + errVar + ` := ` + cfg + `.TOClient.` + funcStr
					newLines = append(newLines, newLine)
					state = stateInFoundTOClientNew
					continue
				}
				newLines = append(newLines, line)
			case stateInFoundTOClientNew:
				// if there's a newline between the call and the unsupported check, skip it.
				if strings.TrimSpace(line) == "" {
					continue
				}
				if line == spacePrefix+`if `+errVar+` == nil && `+unsupportedVar+` {` ||
					line == spacePrefix+`if `+unsupportedVar+` && `+errVar+` == nil `+` {` ||
					line == spacePrefix+`if `+unsupportedVar+` {` {
					state = stateInTOClientNewUnsupportedBlock
					continue // continue without adding this line - we want to remove the check-and-fallback
				}
				return errors.New("parsing file '" + path + "': line contains TOClientNew, but is unexpected format: not followed by a check for the unsupported bool (or the parser failed to understand it)\n")
			case stateInTOClientNewUnsupportedBlock:
				// This is why we have to gofmt - if we didn't, we couldn't be guarnateed the block closing would be prefixed by exactly this many spaces. Then we'd have to parse the AST to find it. Ick.
				if line == spacePrefix+`}` {
					state = stateStart
				}
				// Whether or not we find the closing block, don't add it to the lines to output.
				// We want to remove the TOClientNew check-and-fallback block from the output.
			}
		}
		if state != stateStart {
			return errors.New("parsing modified file '" + path + "': appeared to be malformed (or maybe our parser is just broken, sorry)\n")
		}

		newFile := strings.Join(newLines, "\n")

		if err := ioutil.WriteFile(path, []byte(newFile), 0644); err != nil {
			return errors.New("writing modified file '" + path + "': " + err.Error() + "\n")
		}
	}
	return nil
}

func updateVendoredTOClient(appDir string, branch string) error {
	vendorDir := filepath.Join(appDir, `toreq`, `vendor`)
	vendorTCDir := filepath.Join(vendorDir, `github.com`, `apache`, `trafficcontrol`)
	vendorClientDir := filepath.Join(vendorTCDir, `traffic_ops`, `client`)

	vendorFileInfo, err := os.Stat(vendorClientDir)
	if err != nil {
		return errors.New("getting vendor dir '" + vendorDir + "' info: " + err.Error())
	}
	if !vendorFileInfo.IsDir() {
		return errors.New("getting vendor dir '" + vendorDir + "' info: not a directory")
	}
	if err := os.RemoveAll(vendorClientDir); err != nil {
		return errors.New("removing vendor dir '" + vendorDir + "': " + err.Error())
	}
	if err := os.Mkdir(vendorClientDir, 0755); err != nil {
		return errors.New("creating vendor dir '" + vendorDir + "': " + err.Error())
	}

	cmd := exec.Command(`git`, `show`, branch+`:../../client`)
	cmd.Dir = appDir
	clientFileListBts, err := cmd.Output()
	if err != nil {
		return errors.New("getting files from git: " + err.Error())
	}
	clientFileListStr := string(clientFileListBts)
	clientFileList := strings.Split(clientFileListStr, "\n")
	if len(clientFileList) < 2 {
		return errors.New("getting files from git: got no files")
	}
	clientFileList = clientFileList[1:] // first line is a header, remove it.
	// fmt.Printf("DEBUG got client file list '''%+v'''\n", clientFileListStr)

	for _, clientFile := range clientFileList {
		clientFile = strings.TrimSpace(clientFile)
		if clientFile == "" {
			continue
		}
		cmd := exec.Command(`git`, `show`, branch+`:../../client`+`/`+clientFile)
		cmd.Dir = appDir
		fileBts, err := cmd.Output()
		if err != nil {
			return errors.New("getting client file '" + clientFile + "' from git: " + err.Error())
		}
		if err := ioutil.WriteFile(vendorClientDir+`/`+clientFile, fileBts, 0644); err != nil {
			return errors.New("Error writing vendored file '" + clientFile + "': " + err.Error())
		}
	}

	// update VERSION in vendored dir
	cmd = exec.Command(`git`, `show`, branch+`:VERSION`)
	cmd.Dir = appDir
	versionBts, err := cmd.Output()
	if err != nil {
		return errors.New("getting VERSION file from git: " + err.Error())
	}

	versionPath := filepath.Join(vendorTCDir, `VERSION`)
	if err := ioutil.WriteFile(versionPath, versionBts, 0644); err != nil {
		return errors.New("Error writing vendored VERSION file '" + versionPath + "': " + err.Error())
	}

	cmd = exec.Command(`git`, `rev-parse`, branch)
	cmd.Dir = appDir
	changesetBts, err := cmd.Output()
	if err != nil {
		return errors.New("getting VERSION file from git: " + err.Error())
	}

	changesetTxt := branch + "\n" + string(changesetBts)

	changesetTxtPath := filepath.Join(vendorTCDir, `changeset.txt`)
	if err := ioutil.WriteFile(changesetTxtPath, []byte(changesetTxt), 0644); err != nil {
		return errors.New("Error writing vendored changeset.txt file '" + changesetTxtPath + "': " + err.Error())
	}

	return nil
}
