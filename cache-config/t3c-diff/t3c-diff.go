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
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/kylelemons/godebug/diff"
	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-diff"

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

func main() {
	help := getopt.BoolLong("help", 'h', "Print usage info and exit")
	version := getopt.BoolLong("version", 'V', "Print version information and exit")
	lineComment := getopt.StringLong("line_comment", 'l', "#", "Comment symbol")
	mode := getopt.IntLong("file-mode", 'm', 0644, "file mode default is 644")
	uid := getopt.IntLong("file-uid", 'u', 0, "User id the file being checked should have, default is running process's uid")
	gid := getopt.IntLong("file-gid", 'g', 0, "Group id the file being checked should have, default is running process's gid")
	fa := getopt.StringLong("file-a", 'a', "", "first diff file")
	fb := getopt.StringLong("file-b", 'b', "", "second diff file")
	getopt.ParseV2()

	log.Init(os.Stderr, os.Stderr, os.Stderr, os.Stderr, os.Stderr)

	if *help {
		log.Errorln(usageStr)
		os.Exit(0)
	} else if *version {
		fmt.Println(t3cutil.VersionStr(AppName, Version, GitRevision))
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		log.Errorln(usageStr)
		os.Exit(3)
	}

	fileNameA := strings.TrimSpace(*fa)
	fileNameB := strings.TrimSpace(*fb)

	if len(fileNameA) == 0 || len(fileNameB) == 0 {
		log.Errorln(usageStr)
		os.Exit(4)
	}

	if *uid == 0 {
		*uid = os.Geteuid()
	}

	if *gid == 0 {
		*gid = os.Getgid()
	}

	fileA, fileAExisted, err := readFileOrStdin(fileNameA)
	if err != nil {
		log.Errorf("error reading first: %s\n", err.Error())
		os.Exit(5)
	}
	fileB, fileBExisted, err := readFileOrStdin(fileNameB)
	if err != nil {
		log.Errorf("error reading second: %s\n", err.Error())
		os.Exit(6)
	}

	fileALines := strings.Split(string(fileA), "\n")
	fileALines = t3cutil.UnencodeFilter(fileALines)
	fileALines = t3cutil.CommentsFilter(fileALines, *lineComment)
	fileA = strings.Join(fileALines, "\n")
	fileA = t3cutil.NewLineFilter(fileA)

	fileBLines := strings.Split(string(fileB), "\n")
	fileBLines = t3cutil.UnencodeFilter(fileBLines)
	fileBLines = t3cutil.CommentsFilter(fileBLines, *lineComment)
	fileB = strings.Join(fileBLines, "\n")
	fileB = t3cutil.NewLineFilter(fileB)

	if fileA != fileB {
		match := regexp.MustCompile(`(?m)^\+.*|^-.*`)
		changes := diff.Diff(fileA, fileB)
		for _, change := range match.FindAllString(changes, -1) {
			fmt.Println(change)
		}
		os.Exit(1)
	}
	if fileAExisted != fileBExisted {
		os.Exit(1)
	}
	switch {
	case fileNameA != "stdin":
		if t3cutil.PermCk(fileNameA, *mode) {
			log.Infoln("File permissions are incorrect, should be ", fmt.Sprintf("%#o", *mode))
			os.Exit(1)
		}
		if t3cutil.OwnershipCk(fileNameA, *uid, *gid) {
			log.Infoln("user or group ownership are incorrect, should be ", fmt.Sprintf("Uid:%d Gid:%d", *uid, *gid))
			os.Exit(1)
		}
	case fileNameB != "stdin":
		if t3cutil.PermCk(fileNameB, *mode) {
			log.Infoln("File permissions are incorrect, should be ", fmt.Sprintf("%#o", *mode))
			os.Exit(1)
		}
		if t3cutil.OwnershipCk(fileNameB, *uid, *gid) {
			log.Infoln("user or group ownership are incorrect, should be ", fmt.Sprintf("Uid:%d Gid:%d", *uid, *gid))
			os.Exit(1)
		}
	}
	os.Exit(0)

}

const usageStr = `usage: t3c-diff [--help]
        -a <file-a> -b <file-b> -l <line comment> -m <file mode> -u <file uid> -g <file gid>

Either file may be 'stdin', in which case that file is read from stdin.
Either file may not exist.

Prints the diff to stdout, and returns the exit code 0 if there was no diff, 1 if there was a diff.
If one file exists but the other doesn't, it will always be a diff.

Mode is file permissions in octal format, default is 0644.
Line comment is a character that signals the line is a comment, default is #

Uid is the User id the file being checked should have, default is running process's uid.
Gid is the Group id the file being checked should have, default is running process's gid.

Note this means there may be no diff text printed to stdout but still exit 1 indicating a diff
if the file being created or deleted is semantically empty.`

// readFileOrStdin reads the file, or if fileOrStdin is 'stdin', reads from stdin.
// Returns the file, whether it existed, and any error.
func readFileOrStdin(fileOrStdin string) (string, bool, error) {
	if strings.ToLower(fileOrStdin) == "stdin" {
		bts, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", false, errors.New("reading stdin: " + err.Error())
		}
		return string(bts), true, nil
	}
	bts, err := ioutil.ReadFile(fileOrStdin)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, errors.New("reading file: " + err.Error())
	}
	return string(bts), true, nil
}
