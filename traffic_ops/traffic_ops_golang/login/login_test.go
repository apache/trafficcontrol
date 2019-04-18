package login

import "testing"

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

func TestVerifyUrlOnWhiteList(t *testing.T) {
	url := "test.right"
	whitelistedUrls := [][]string{[]string{}, []string{""}, []string{"*"}, []string{"test.wrong"}, []string{"test.right"}, []string{"*.right"}, []string{"test.wrong", "test.right"}}
	expected := []bool{false, false, true, false, true, true, true}

	for i, urlList := range whitelistedUrls {
		if VerifyUrlOnWhiteList(url, urlList) != expected[i] {
			t.Errorf("expected: %v, actual: %v", expected[i], VerifyUrlOnWhiteList(url, urlList))
		}
	}
}
