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
	"bytes"
	"errors"
	"os/exec"
	"testing"
)

func t3c_check_refs_exec(filename string, t *testing.T) (int, error) {
	if !fileExists("./t3c-check-refs") {
		t.Fatalf("You must first build t3c-check-refs before running tests")
	}
	args := []string{
		"--trafficserver-config-dir=./test-files/etc",
		"--trafficserver-plugin-dir=./test-files/libexec",
	}
	args = append(args, filename)
	cmd := exec.Command("./t3c-check-refs", args...)
	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		return -1, errors.New("error from t3c-check-refs: " + err.Error() + ": " + errbuf.String())
	}

	return cmd.ProcessState.ExitCode(), nil
}

func TestRemapConfig(t *testing.T) {
	rc, err := t3c_check_refs_exec("./test-files/etc/remap.config", t)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if rc != 0 {
		t.Errorf("expected 0 errors got %d errors\n", rc)
	}
}

func TestBadRemapConfig(t *testing.T) {
	rc, _ := t3c_check_refs_exec("./test-files/etc/bad-remap.config", t)
	if rc != -1 {
		t.Errorf("expected 2 errors got %d errors\n", rc)
	}
}

func TestMultiRemapConfig(t *testing.T) {
	rc, err := t3c_check_refs_exec("./test-files/etc/remap-multiline.config", t)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if rc != 0 {
		t.Errorf("expected 0 errors got %d errors\n", rc)
	}
}

func TestBadRemapMultilineConfig(t *testing.T) {
	rc, _ := t3c_check_refs_exec("./test-files/etc/bad-remap-multiline.config", t)
	if rc != -1 {
		t.Errorf("expected 0 errors got %d errors\n", rc)
	}
}

func TestPluConfig(t *testing.T) {
	rc, err := t3c_check_refs_exec("./test-files/etc/plugin.config", t)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if rc != 0 {
		t.Errorf("expected 0 errors got %d errors\n", rc)
	}
}

func TestBadPluConfig(t *testing.T) {
	rc, _ := t3c_check_refs_exec("./test-files/etc/bad-plugin.config", t)
	if rc != -1 {
		t.Errorf("expected 0 errors got %d errors\n", rc)
	}
}
