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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

// A lookup table, bool will always be true
var commonPasswords map[string]bool

// Expects a relative path from the traffic_ops directory
func LoadPasswordBlacklist(filePath string) error {

	if commonPasswords != nil {
		return errors.New("Password blacklist is already loaded")
	}

	commonPasswords = make(map[string]bool)

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
		commonPasswords[scanner.Text()] = true
	}

	in.Close()
	return scanner.Err()
}

func IsCommonPassword(pw string) bool {
	_, yes := commonPasswords[pw]
	return yes
}

func IsGoodLoginPair(username string, password string) (bool, error) {

	if username == "" {
		return false, errors.New("Your username cannot be blank.")
	}

	if username == password {
		return false, errors.New("Your password cannot be your username.")
	}

	return IsGoodPassword(password)
}

func IsGoodPassword(password string) (bool, error) {

	if len(password) < 8 {
		return false, errors.New("Password must be greater than 7 characters.")
	}

	if IsCommonPassword(password) {
		return false, errors.New("Password is too common.")
	}

	return true, nil
}
