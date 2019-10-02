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
	"testing"
)

func TestDeriveGoodPassword(t *testing.T) {

	pass := "password"
	derivedPassword, err := DerivePassword(pass)
	err = VerifySCRYPTPassword(pass, derivedPassword)
	if err != nil {
		t.Errorf("password should be valid")
	}
}

func TestDeriveBadPassword(t *testing.T) {

	pass := "password"
	derivedPassword, err := DerivePassword(pass)
	err = VerifySCRYPTPassword("badpassword", derivedPassword)
	if err == nil {
		t.Errorf("password should be invalid")
	}
}

func TestScryptPasswordIsRequired(t *testing.T) {

	err := VerifySCRYPTPassword("password", "")
	if err == nil {
		t.Errorf("scrypt password should be required")
	}
}

// The purpose of this test is to show that all password requirements are being met
func TestUsernamePassword(t *testing.T) {
	if commonPasswords == nil {
		defer func() { commonPasswords = nil }() // global variable, reset after this test
	}

	passwords := []string{"username", "password", "pa$$word", "", "red"}
	expected := []bool{false, false, true, false, false}
	if err := LoadPasswordBlacklist("app/conf/invalid_passwords.txt"); err != nil {
		t.Fatalf("LoadPasswordBlacklist err expected: nil, actual: %v", err)
	}

	for i, password := range passwords {
		if ok, err := IsGoodLoginPair("username", password); ok != expected[i] {
			if expected[i] {
				t.Errorf("\"%s\" should have been marked as an invalid password", password)
			} else {
				t.Errorf("\"%s\" should be an ok password, but got the error: %v", password, err)
			}
		}
	}

	if ok, _ := IsGoodLoginPair("", "GoOdPa$$woRd"); ok {
		t.Errorf("An empty username should not pass")
	}
}

// The purpose of this test is to show that the file is being read, and we can tell if a password is in the file
func TestCommonPassword(t *testing.T) {
	if commonPasswords == nil {
		defer func() { commonPasswords = nil }() // global variable, reset after this test
	}

	passwords := []string{"password", "pa$$word"}
	expected := []bool{true, false}

	if err := LoadPasswordBlacklist("app/conf/invalid_passwords.txt"); err != nil {
		t.Fatalf("LoadPasswordBlacklist err expected: nil, actual: %v", err)
	}

	for i, password := range passwords {
		if IsCommonPassword(password) != expected[i] {
			if expected[i] {
				t.Errorf("\"%s\" should have been marked as an invalid password", password)
			} else {
				t.Errorf("\"%s\" should be an ok password", password)
			}
		}
	}

}
