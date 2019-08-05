package login

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

import "testing"

func TestVerifyUrlOnWhiteList(t *testing.T) {
	type TestResult struct {
		Whitelist      []string
		ExpectedResult bool
	}

	completeTestResults := struct {
		Results []TestResult
	}{}

	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{}, ExpectedResult: false})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{""}, ExpectedResult: false})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"*"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong"}, ExpectedResult: false})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.right.com"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"*.right.com"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "test.right.com"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "*.right.*"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "*right*"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "*right"}, ExpectedResult: false})

	url := "https://test.right.com/other/parts"

	for _, result := range completeTestResults.Results {
		if matched, _ := VerifyUrlOnWhiteList(url, result.Whitelist); matched != result.ExpectedResult {
			t.Errorf("for whitelist: %v, expected: %v, actual: %v", result.Whitelist, result.ExpectedResult, matched)
		}
	}
}
