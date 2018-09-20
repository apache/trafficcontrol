package auth

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
	"fmt"
	"os"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
)

// A lookup table, bool will always be true
var invalidPasswords map[string]bool

// Expects a relative path from the traffic_ops directory
func LoadPasswordBlacklist(filePath string) error {

	if invalidPasswords == nil {
		invalidPasswords = make(map[string]bool)
	} else {
		return fmt.Errorf("Password blacklist is already loaded")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath = fmt.Sprintf("%straffic_ops/%s", pwd[:strings.Index(pwd, "traffic_ops")], filePath)

	log.Infof("full path to password blacklist: %s\n", filePath)
	in, err := os.Open(filePath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		invalidPasswords[scanner.Text()] = true
	}

	in.Close()
	return scanner.Err()
}

func IsInvalidPassword(pw string) bool {
	_, bad := invalidPasswords[pw]
	return bad
}

func IsGoodPassword(pw string, confirmPw string, username string) error {

	if pw == "" {
		return nil
	}

	if pw != confirmPw {
		return fmt.Errorf("Passwords do not match.")
	}

	if pw == username {
		return fmt.Errorf("Your password cannot be the same as your username.")
	}

	if len(pw) < 8 {
		return fmt.Errorf("Password must be greater than 7 chars.")
	}

	if IsInvalidPassword(pw) {
		return fmt.Errorf("Password is too common.")
	}

	// At this point we're happy with the password
	return nil
}
