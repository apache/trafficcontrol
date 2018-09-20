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
	"github.com/apache/trafficcontrol/lib/go-log"
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

func TestPasswordStrength(t *testing.T) {

	passwords := []string{"password", "pa$$word"}
	expected := []bool{true, false}
	LoadPasswordBlacklist("app/conf/invalid_passwords.txt")

	for i, password := range passwords {
		if IsInvalidPassword(password) != expected[i] {
			if expected[i] {
				t.Errorf("%s should have been marked as an invalid password", password)
			} else {
				t.Errorf("%s should be an ok password", password)
			}
		}
	}

}
